// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package msg

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/mcache"
	"github.com/openimsdk/open-im-server/v3/pkg/dbbuild"
	"github.com/openimsdk/open-im-server/v3/pkg/mqbuild"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcli"

	"google.golang.org/grpc"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/redis"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/controller"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database/mgo"
	"github.com/openimsdk/open-im-server/v3/pkg/common/webhook"
	"github.com/openimsdk/open-im-server/v3/pkg/notification"
	"github.com/openimsdk/open-im-server/v3/pkg/rpccache"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/conversation"
	"github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/discovery"
)

type MessageInterceptorFunc func(ctx context.Context, globalConfig *Config, req *msg.SendMsgReq) (*sdkws.MsgData, error)

// MessageInterceptorChain defines a chain of message interceptor functions.
type MessageInterceptorChain []MessageInterceptorFunc

type Config struct {
	RpcConfig          config.Msg
	RedisConfig        config.Redis
	MongodbConfig      config.Mongo
	KafkaConfig        config.Kafka
	NotificationConfig config.Notification
	Share              config.Share
	WebhooksConfig     config.Webhooks
	LocalCacheConfig   config.LocalCache
	Discovery          config.Discovery
}

// MsgServer encapsulates dependencies required for message handling.
type msgServer struct {
	msg.UnimplementedMsgServer
	RegisterCenter         discovery.Conn                   // Service discovery registry for service registration.
	MsgDatabase            controller.CommonMsgDatabase     // Interface for message database operations.
	UserLocalCache         *rpccache.UserLocalCache         // Local cache for user data.
	FriendLocalCache       *rpccache.FriendLocalCache       // Local cache for friend data.
	GroupLocalCache        *rpccache.GroupLocalCache        // Local cache for group data.
	ConversationLocalCache *rpccache.ConversationLocalCache // Local cache for conversation data.
	Handlers               MessageInterceptorChain          // Chain of handlers for processing messages.
	notificationSender     *notification.NotificationSender // RPC client for sending notifications.
	msgNotificationSender  *MsgNotificationSender           // RPC client for sending msg notifications.
	config                 *Config                          // Global configuration settings.
	webhookClient          *webhook.Client
	conversationClient     *rpcli.ConversationClient
}

func (m *msgServer) addInterceptorHandler(interceptorFunc ...MessageInterceptorFunc) {
	m.Handlers = append(m.Handlers, interceptorFunc...)

}

func Start(ctx context.Context, config *Config, client discovery.SvcDiscoveryRegistry, server grpc.ServiceRegistrar) error {
	builder := mqbuild.NewBuilder(&config.KafkaConfig)
	redisProducer, err := builder.GetTopicProducer(ctx, config.KafkaConfig.ToRedisTopic)
	if err != nil {
		return err
	}
	dbb := dbbuild.NewBuilder(&config.MongodbConfig, &config.RedisConfig)
	mgocli, err := dbb.Mongo(ctx)
	if err != nil {
		return err
	}
	rdb, err := dbb.Redis(ctx)
	if err != nil {
		return err
	}
	msgDocModel, err := mgo.NewMsgMongo(mgocli.GetDB())
	if err != nil {
		return err
	}
	var msgModel cache.MsgCache
	if rdb == nil {
		cm, err := mgo.NewCacheMgo(mgocli.GetDB())
		if err != nil {
			return err
		}
		msgModel = mcache.NewMsgCache(cm, msgDocModel)
	} else {
		msgModel = redis.NewMsgCache(rdb, msgDocModel)
	}
	seqConversation, err := mgo.NewSeqConversationMongo(mgocli.GetDB())
	if err != nil {
		return err
	}
	seqConversationCache := redis.NewSeqConversationCacheRedis(rdb, seqConversation)
	seqUser, err := mgo.NewSeqUserMongo(mgocli.GetDB())
	if err != nil {
		return err
	}
	seqUserCache := redis.NewSeqUserCacheRedis(rdb, seqUser)
	userConn, err := client.GetConn(ctx, config.Discovery.RpcService.User)
	if err != nil {
		return err
	}
	groupConn, err := client.GetConn(ctx, config.Discovery.RpcService.Group)
	if err != nil {
		return err
	}
	friendConn, err := client.GetConn(ctx, config.Discovery.RpcService.Friend)
	if err != nil {
		return err
	}
	conversationConn, err := client.GetConn(ctx, config.Discovery.RpcService.Conversation)
	if err != nil {
		return err
	}
	conversationClient := rpcli.NewConversationClient(conversationConn)
	msgDatabase := controller.NewCommonMsgDatabase(msgDocModel, msgModel, seqUserCache, seqConversationCache, redisProducer)
	s := &msgServer{
		MsgDatabase:            msgDatabase,
		RegisterCenter:         client,
		UserLocalCache:         rpccache.NewUserLocalCache(rpcli.NewUserClient(userConn), &config.LocalCacheConfig, rdb),
		GroupLocalCache:        rpccache.NewGroupLocalCache(rpcli.NewGroupClient(groupConn), &config.LocalCacheConfig, rdb),
		ConversationLocalCache: rpccache.NewConversationLocalCache(conversationClient, &config.LocalCacheConfig, rdb),
		FriendLocalCache:       rpccache.NewFriendLocalCache(rpcli.NewRelationClient(friendConn), &config.LocalCacheConfig, rdb),
		config:                 config,
		webhookClient:          webhook.NewWebhookClient(config.WebhooksConfig.URL),
		conversationClient:     conversationClient,
	}

	s.notificationSender = notification.NewNotificationSender(&config.NotificationConfig, notification.WithLocalSendMsg(s.SendMsg))
	s.msgNotificationSender = NewMsgNotificationSender(config, notification.WithLocalSendMsg(s.SendMsg))

	msg.RegisterMsgServer(server, s)

	return nil
}

func (m *msgServer) conversationAndGetRecvID(conversation *conversation.Conversation, userID string) string {
	if conversation.ConversationType == constant.SingleChatType ||
		conversation.ConversationType == constant.NotificationChatType {
		if userID == conversation.OwnerUserID {
			return conversation.UserID
		} else {
			return conversation.OwnerUserID
		}
	} else if conversation.ConversationType == constant.ReadGroupChatType {
		return conversation.GroupID
	}
	return ""
}
