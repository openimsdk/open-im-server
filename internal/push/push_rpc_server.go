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

package push

import (
	"context"
	"sync"

	"google.golang.org/grpc"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/cache"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/controller"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/localcache"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/rpcclient"
	"github.com/OpenIMSDK/protocol/constant"
	pbPush "github.com/OpenIMSDK/protocol/push"
	"github.com/OpenIMSDK/tools/discoveryregistry"
	"github.com/OpenIMSDK/tools/log"
)

type pushServer struct {
	pusher *Pusher
}

func Start(client discoveryregistry.SvcDiscoveryRegistry, server *grpc.Server) error {
	rdb, err := cache.NewRedis()
	if err != nil {
		return err
	}
	cacheModel := cache.NewMsgCacheModel(rdb)
	offlinePusher := NewOfflinePusher(cacheModel)
	database := controller.NewPushDatabase(cacheModel)
	groupRpcClient := rpcclient.NewGroupRpcClient(client)
	conversationRpcClient := rpcclient.NewConversationRpcClient(client)
	msgRpcClient := rpcclient.NewMessageRpcClient(client)
	pusher := NewPusher(
		client,
		offlinePusher,
		database,
		localcache.NewGroupLocalCache(&groupRpcClient),
		localcache.NewConversationLocalCache(&conversationRpcClient),
		&conversationRpcClient,
		&groupRpcClient,
		&msgRpcClient,
	)
	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		pbPush.RegisterPushMsgServiceServer(server, &pushServer{
			pusher: pusher,
		})
	}()
	go func() {
		defer wg.Done()
		consumer := NewConsumer(pusher)
		consumer.initPrometheus()
		consumer.Start()
	}()
	wg.Wait()
	return nil
}

func (r *pushServer) PushMsg(ctx context.Context, pbData *pbPush.PushMsgReq) (resp *pbPush.PushMsgResp, err error) {
	switch pbData.MsgData.SessionType {
	case constant.SuperGroupChatType:
		err = r.pusher.Push2SuperGroup(ctx, pbData.MsgData.GroupID, pbData.MsgData)
	default:
		err = r.pusher.Push2User(ctx, []string{pbData.MsgData.RecvID, pbData.MsgData.SendID}, pbData.MsgData)
	}
	if err != nil {
		if err != errNoOfflinePusher {
			return nil, err
		} else {
			log.ZWarn(ctx, "offline push failed", err, "msg", pbData.String())
		}
	}
	return &pbPush.PushMsgResp{}, nil
}

func (r *pushServer) DelUserPushToken(
	ctx context.Context,
	req *pbPush.DelUserPushTokenReq,
) (resp *pbPush.DelUserPushTokenResp, err error) {
	if err = r.pusher.database.DelFcmToken(ctx, req.UserID, int(req.PlatformID)); err != nil {
		return nil, err
	}
	return &pbPush.DelUserPushTokenResp{}, nil
}
