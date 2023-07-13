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
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/relation"

	"google.golang.org/grpc"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/cache"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/controller"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/localcache"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/unrelation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/prome"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/conversation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/msg"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/rpcclient"
)

type MessageInterceptorChain []MessageInterceptorFunc
type msgServer struct {
	RegisterCenter         discoveryregistry.SvcDiscoveryRegistry
	MsgDatabase            controller.CommonMsgDatabase
	Group                  *rpcclient.GroupRpcClient
	User                   *rpcclient.UserRpcClient
	Conversation           *rpcclient.ConversationRpcClient
	friend                 *rpcclient.FriendRpcClient
	GroupLocalCache        *localcache.GroupLocalCache
	ConversationLocalCache *localcache.ConversationLocalCache
	MessageLocker          MessageLocker
	Handlers               MessageInterceptorChain
	notificationSender     *rpcclient.NotificationSender
}

func (m *msgServer) addInterceptorHandler(interceptorFunc ...MessageInterceptorFunc) {
	m.Handlers = append(m.Handlers, interceptorFunc...)
}

func (m *msgServer) execInterceptorHandler(ctx context.Context, req *msg.SendMsgReq) error {
	for _, handler := range m.Handlers {
		msgData, err := handler(ctx, req)
		if err != nil {
			return err
		}
		req.MsgData = msgData
	}
	return nil
}

func Start(client discoveryregistry.SvcDiscoveryRegistry, server *grpc.Server) error {
	rdb, err := cache.NewRedis()
	if err != nil {
		return err
	}
	mongo, err := unrelation.NewMongo()
	if err != nil {
		return err
	}
	if err := mongo.CreateMsgIndex(); err != nil {
		return err
	}
	cacheModel := cache.NewMsgCacheModel(rdb)
	msgDocModel := unrelation.NewMsgMongoDriver(mongo.GetDatabase())
	conversationClient := rpcclient.NewConversationRpcClient(client)
	userRpcClient := rpcclient.NewUserRpcClient(client)
	groupRpcClient := rpcclient.NewGroupRpcClient(client)
	friendRpcClient := rpcclient.NewFriendRpcClient(client)
	mysql, err := relation.NewGormDB()
	if err != nil {
		return err
	}
	msgMysModel := relation.NewChatLogGorm(mysql)
	msgDatabase := controller.NewCommonMsgDatabase(msgDocModel, cacheModel, msgMysModel)
	s := &msgServer{
		Conversation:           &conversationClient,
		User:                   &userRpcClient,
		Group:                  &groupRpcClient,
		MsgDatabase:            msgDatabase,
		RegisterCenter:         client,
		GroupLocalCache:        localcache.NewGroupLocalCache(&groupRpcClient),
		ConversationLocalCache: localcache.NewConversationLocalCache(&conversationClient),
		friend:                 &friendRpcClient,
		MessageLocker:          NewLockerMessage(cacheModel),
	}
	s.notificationSender = rpcclient.NewNotificationSender(rpcclient.WithLocalSendMsg(s.SendMsg))
	s.addInterceptorHandler(MessageHasReadEnabled)
	s.initPrometheus()
	msg.RegisterMsgServer(server, s)
	return nil
}

func (m *msgServer) initPrometheus() {
	prome.NewMsgPullFromRedisSuccessCounter()
	prome.NewMsgPullFromRedisFailedCounter()
	prome.NewMsgPullFromMongoSuccessCounter()
	prome.NewMsgPullFromMongoFailedCounter()
	prome.NewSingleChatMsgRecvSuccessCounter()
	prome.NewGroupChatMsgRecvSuccessCounter()
	prome.NewWorkSuperGroupChatMsgRecvSuccessCounter()
	prome.NewSingleChatMsgProcessSuccessCounter()
	prome.NewSingleChatMsgProcessFailedCounter()
	prome.NewGroupChatMsgProcessSuccessCounter()
	prome.NewGroupChatMsgProcessFailedCounter()
	prome.NewWorkSuperGroupChatMsgProcessSuccessCounter()
	prome.NewWorkSuperGroupChatMsgProcessFailedCounter()
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
