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
	"github.com/openimsdk/open-im-server/v3/pkg/common/cmd"
	"github.com/openimsdk/tools/db/mongoutil"
	"github.com/openimsdk/tools/db/redisutil"

	"github.com/openimsdk/open-im-server/v3/pkg/common/db/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/controller"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/mgo"
	"github.com/openimsdk/open-im-server/v3/pkg/rpccache"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/conversation"
	"github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/tools/discovery"
	"google.golang.org/grpc"
)

type (
	// MessageInterceptorChain defines a chain of message interceptor functions.
	MessageInterceptorChain []MessageInterceptorFunc

	// MsgServer encapsulates dependencies required for message handling.
	msgServer struct {
		RegisterCenter         discovery.SvcDiscoveryRegistry   // Service discovery registry for service registration.
		MsgDatabase            controller.CommonMsgDatabase     // Interface for message database operations.
		Conversation           *rpcclient.ConversationRpcClient // RPC client for conversation service.
		UserLocalCache         *rpccache.UserLocalCache         // Local cache for user data.
		FriendLocalCache       *rpccache.FriendLocalCache       // Local cache for friend data.
		GroupLocalCache        *rpccache.GroupLocalCache        // Local cache for group data.
		ConversationLocalCache *rpccache.ConversationLocalCache // Local cache for conversation data.
		Handlers               MessageInterceptorChain          // Chain of handlers for processing messages.
		notificationSender     *rpcclient.NotificationSender    // RPC client for sending notifications.
		config                 *cmd.MsgConfig                   // Global configuration settings.
	}
)

func (m *msgServer) addInterceptorHandler(interceptorFunc ...MessageInterceptorFunc) {
	m.Handlers = append(m.Handlers, interceptorFunc...)

}

func Start(ctx context.Context, config *cmd.MsgConfig, client discovery.SvcDiscoveryRegistry, server *grpc.Server) error {
	mgocli, err := mongoutil.NewMongoDB(ctx, config.MongodbConfig.Build())
	if err != nil {
		return err
	}
	rdb, err := redisutil.NewRedisClient(ctx, config.RedisConfig.Build())
	if err != nil {
		return err
	}
	msgDocModel, err := mgo.NewMsgMongo(mgocli.GetDB())
	if err != nil {
		return err
	}
	//todo MsgCacheTimeout
	msgModel := cache.NewMsgCache(rdb, 86400, config.RedisConfig.EnablePipeline)
	seqModel := cache.NewSeqCache(rdb)
	conversationClient := rpcclient.NewConversationRpcClient(client, config.Share.RpcRegisterName.Conversation)
	userRpcClient := rpcclient.NewUserRpcClient(client, config.Share.RpcRegisterName.User, &config.Share.IMAdmin)
	groupRpcClient := rpcclient.NewGroupRpcClient(client, config.Share.RpcRegisterName.Group)
	friendRpcClient := rpcclient.NewFriendRpcClient(client, config.Share.RpcRegisterName.Friend)
	msgDatabase, err := controller.NewCommonMsgDatabase(msgDocModel, msgModel, seqModel, &config.KafkaConfig)
	if err != nil {
		return err
	}
	s := &msgServer{
		Conversation:           &conversationClient,
		MsgDatabase:            msgDatabase,
		RegisterCenter:         client,
		UserLocalCache:         rpccache.NewUserLocalCache(userRpcClient, rdb),
		GroupLocalCache:        rpccache.NewGroupLocalCache(groupRpcClient, rdb),
		ConversationLocalCache: rpccache.NewConversationLocalCache(conversationClient, rdb),
		FriendLocalCache:       rpccache.NewFriendLocalCache(friendRpcClient, rdb),
		config:                 config,
	}

	s.notificationSender = rpcclient.NewNotificationSender(&config.NotificationConfig, rpcclient.WithLocalSendMsg(s.SendMsg))
	s.addInterceptorHandler(MessageHasReadEnabled)
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
	} else if conversation.ConversationType == constant.SuperGroupChatType {
		return conversation.GroupID
	}
	return ""
}
