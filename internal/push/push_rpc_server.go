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
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/tools/db/redisutil"

	"github.com/openimsdk/open-im-server/v3/pkg/common/db/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/controller"
	"github.com/openimsdk/open-im-server/v3/pkg/rpccache"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient"
	"github.com/openimsdk/protocol/constant"
	pbpush "github.com/openimsdk/protocol/push"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/datautil"
	"google.golang.org/grpc"
)

type pushServer struct {
	pusher *Pusher
}

type Config struct {
	RpcConfig          config.Push
	RedisConfig        config.Redis
	MongodbConfig      config.Mongo
	KafkaConfig        config.Kafka
	ZookeeperConfig    config.ZooKeeper
	NotificationConfig config.Notification
	Share              config.Share
	WebhooksConfig     config.Webhooks
	LocalCacheConfig   config.LocalCache
}

func Start(ctx context.Context, config *Config, client discovery.SvcDiscoveryRegistry, server *grpc.Server) error {
	rdb, err := redisutil.NewRedisClient(ctx, config.RedisConfig.Build())
	if err != nil {
		return err
	}
	cacheModel := cache.NewThirdCache(rdb)
	offlinePusher, err := NewOfflinePusher(&config.RpcConfig, cacheModel)
	if err != nil {
		return err
	}
	database := controller.NewPushDatabase(cacheModel)
	groupRpcClient := rpcclient.NewGroupRpcClient(client, config.Share.RpcRegisterName.Group)
	conversationRpcClient := rpcclient.NewConversationRpcClient(client, config.Share.RpcRegisterName.Conversation)
	msgRpcClient := rpcclient.NewMessageRpcClient(client, config.Share.RpcRegisterName.Msg)
	pusher := NewPusher(
		config,
		client,
		offlinePusher,
		database,
		rpccache.NewGroupLocalCache(groupRpcClient, &config.LocalCacheConfig, rdb),
		rpccache.NewConversationLocalCache(conversationRpcClient, &config.LocalCacheConfig, rdb),
		&conversationRpcClient,
		&groupRpcClient,
		&msgRpcClient,
	)

	pbpush.RegisterPushMsgServiceServer(server, &pushServer{
		pusher: pusher,
	})

	consumer, err := NewConsumer(&config.KafkaConfig, pusher)
	if err != nil {
		return err
	}

	consumer.Start()

	return nil
}

func (r *pushServer) PushMsg(ctx context.Context, pbData *pbpush.PushMsgReq) (resp *pbpush.PushMsgResp, err error) {
	switch pbData.MsgData.SessionType {
	case constant.SuperGroupChatType:
		err = r.pusher.Push2SuperGroup(ctx, pbData.MsgData.GroupID, pbData.MsgData)
	default:
		var pushUserIDList []string
		isSenderSync := datautil.GetSwitchFromOptions(pbData.MsgData.Options, constant.IsSenderSync)
		if !isSenderSync {
			pushUserIDList = append(pushUserIDList, pbData.MsgData.RecvID)
		} else {
			pushUserIDList = append(pushUserIDList, pbData.MsgData.RecvID, pbData.MsgData.SendID)
		}
		err = r.pusher.Push2User(ctx, pushUserIDList, pbData.MsgData)
	}
	if err != nil {
		if err != errNoOfflinePusher {
			return nil, err
		}
		log.ZWarn(ctx, "offline push failed", err, "msg", pbData.String())
	}
	return &pbpush.PushMsgResp{}, nil
}

func (r *pushServer) DelUserPushToken(
	ctx context.Context,
	req *pbpush.DelUserPushTokenReq,
) (resp *pbpush.DelUserPushTokenResp, err error) {
	if err = r.pusher.database.DelFcmToken(ctx, req.UserID, int(req.PlatformID)); err != nil {
		return nil, err
	}
	return &pbpush.DelUserPushTokenResp{}, nil
}
