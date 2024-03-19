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
	"github.com/OpenIMSDK/protocol/constant"
	"github.com/OpenIMSDK/protocol/conversation"
	"github.com/OpenIMSDK/protocol/msg"
	"github.com/OpenIMSDK/tools/discoveryregistry"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/controller"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/unrelation"
	"github.com/openimsdk/open-im-server/v3/pkg/rpccache"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient"
	"google.golang.org/grpc"
)

type (
	// MessageInterceptorChain defines a chain of message interceptor functions.
	MessageInterceptorChain []MessageInterceptorFunc

	// MsgServer encapsulates dependencies required for message handling.
	msgServer struct {
		RegisterCenter         discoveryregistry.SvcDiscoveryRegistry // Service discovery registry for service registration.
		MsgDatabase            controller.CommonMsgDatabase           // Interface for message database operations.
		Conversation           *rpcclient.ConversationRpcClient       // RPC client for conversation service.
		UserLocalCache         *rpccache.UserLocalCache               // Local cache for user data.
		FriendLocalCache       *rpccache.FriendLocalCache             // Local cache for friend data.
		GroupLocalCache        *rpccache.GroupLocalCache              // Local cache for group data.
		ConversationLocalCache *rpccache.ConversationLocalCache       // Local cache for conversation data.
		Handlers               MessageInterceptorChain                // Chain of handlers for processing messages.
		notificationSender     *rpcclient.NotificationSender          // RPC client for sending notifications.
		config                 *config.GlobalConfig                   // Global configuration settings.
	}
)

func (m *msgServer) addInterceptorHandler(interceptorFunc ...MessageInterceptorFunc) {
	m.Handlers = append(m.Handlers, interceptorFunc...)
}

//func (m *msgServer) execInterceptorHandler(ctx context.Context, config *config.GlobalConfig, req *msg.SendMsgReq) error {
//	for _, handler := range m.Handlers {
//		msgData, err := handler(ctx, config, req)
//		if err != nil {
//			return err
//		}
//		req.MsgData = msgData
//	}
//	return nil
//}

func Start(config *config.GlobalConfig, client discoveryregistry.SvcDiscoveryRegistry, server *grpc.Server) error {
	rdb, err := cache.NewRedis(&config.Redis)
	if err != nil {
		return err
	}
	mongo, err := unrelation.NewMongo(&config.Mongo)
	if err != nil {
		return err
	}
	if err := mongo.CreateMsgIndex(); err != nil {
		return err
	}
	cacheModel := cache.NewMsgCacheModel(rdb, config.MsgCacheTimeout, &config.Redis)
	msgDocModel := unrelation.NewMsgMongoDriver(mongo.GetDatabase(config.Mongo.Database))
	conversationClient := rpcclient.NewConversationRpcClient(client, config.RpcRegisterName.OpenImConversationName)
	userRpcClient := rpcclient.NewUserRpcClient(client, config.RpcRegisterName.OpenImUserName, &config.Manager, &config.IMAdmin)
	groupRpcClient := rpcclient.NewGroupRpcClient(client, config.RpcRegisterName.OpenImGroupName)
	friendRpcClient := rpcclient.NewFriendRpcClient(client, config.RpcRegisterName.OpenImFriendName)
	msgDatabase, err := controller.NewCommonMsgDatabase(msgDocModel, cacheModel, &config.Kafka)
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

	s.notificationSender = rpcclient.NewNotificationSender(&config.Notification, rpcclient.WithLocalSendMsg(s.SendMsg))
	s.addInterceptorHandler(MessageHasReadEnabled)
	msg.RegisterMsgServer(server, s)
	return nil
}

func (m *msgServer) conversationAndGetRecvID(conversation *conversation.Conversation, userID string) (recvID string) {
	if conversation.ConversationType == constant.SingleChatType ||
		conversation.ConversationType == constant.NotificationChatType {
		if userID == conversation.OwnerUserID {
			recvID = conversation.UserID
		} else {
			recvID = conversation.OwnerUserID
		}
	} else if conversation.ConversationType == constant.SuperGroupChatType {
		recvID = conversation.GroupID
	}
	return
}
