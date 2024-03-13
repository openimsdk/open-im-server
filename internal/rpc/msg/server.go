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
	MessageInterceptorChain []MessageInterceptorFunc
	msgServer               struct {
		RegisterCenter         discoveryregistry.SvcDiscoveryRegistry
		MsgDatabase            controller.CommonMsgDatabase
		Conversation           *rpcclient.ConversationRpcClient
		UserLocalCache         *rpccache.UserLocalCache
		FriendLocalCache       *rpccache.FriendLocalCache
		GroupLocalCache        *rpccache.GroupLocalCache
		ConversationLocalCache *rpccache.ConversationLocalCache
		Handlers               MessageInterceptorChain
		notificationSender     *rpcclient.NotificationSender
		config                 *config.GlobalConfig
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
	rdb, err := cache.NewRedis(config)
	if err != nil {
		return err
	}
	mongo, err := unrelation.NewMongo(config)
	if err != nil {
		return err
	}
	if err := mongo.CreateMsgIndex(); err != nil {
		return err
	}
	cacheModel := cache.NewMsgCacheModel(rdb, config)
	msgDocModel := unrelation.NewMsgMongoDriver(mongo.GetDatabase(config.Mongo.Database))
	conversationClient := rpcclient.NewConversationRpcClient(client, config)
	userRpcClient := rpcclient.NewUserRpcClient(client, config)
	groupRpcClient := rpcclient.NewGroupRpcClient(client, config)
	friendRpcClient := rpcclient.NewFriendRpcClient(client, config)
	msgDatabase, err := controller.NewCommonMsgDatabase(msgDocModel, cacheModel, config)
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
	s.notificationSender = rpcclient.NewNotificationSender(config, rpcclient.WithLocalSendMsg(s.SendMsg))
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
