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

package controller

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache"
	"github.com/openimsdk/protocol/push"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mq"
	"google.golang.org/protobuf/proto"
)

type PushDatabase interface {
	DelFcmToken(ctx context.Context, userID string, platformID int) error
	MsgToOfflinePushMQ(ctx context.Context, key string, userIDs []string, msg2mq *sdkws.MsgData) error
}

type pushDataBase struct {
	cache                 cache.ThirdCache
	producerToOfflinePush mq.Producer
}

func NewPushDatabase(cache cache.ThirdCache, offlinePushProducer mq.Producer) PushDatabase {
	return &pushDataBase{
		cache:                 cache,
		producerToOfflinePush: offlinePushProducer,
	}
}

func (p *pushDataBase) DelFcmToken(ctx context.Context, userID string, platformID int) error {
	return p.cache.DelFcmToken(ctx, userID, platformID)
}

func (p *pushDataBase) MsgToOfflinePushMQ(ctx context.Context, key string, userIDs []string, msg2mq *sdkws.MsgData) error {
	data, err := proto.Marshal(&push.PushMsgReq{MsgData: msg2mq, UserIDs: userIDs})
	if err != nil {
		return err
	}
	if err := p.producerToOfflinePush.SendMessage(ctx, key, data); err != nil {
		log.ZError(ctx, "message is push to offlinePush topic", err, "key", key, "userIDs", userIDs, "msg", msg2mq.String())
	}
	return err
}
