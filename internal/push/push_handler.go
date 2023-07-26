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

	"github.com/Shopify/sarama"
	"google.golang.org/protobuf/proto"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	kfk "github.com/OpenIMSDK/Open-IM-Server/pkg/common/kafka"
	"github.com/OpenIMSDK/protocol/constant"
	pbChat "github.com/OpenIMSDK/protocol/msg"
	pbPush "github.com/OpenIMSDK/protocol/push"
	"github.com/OpenIMSDK/tools/log"
	"github.com/OpenIMSDK/tools/utils"
)

type ConsumerHandler struct {
	pushConsumerGroup *kfk.MConsumerGroup
	pusher            *Pusher
}

func NewConsumerHandler(pusher *Pusher) *ConsumerHandler {
	var consumerHandler ConsumerHandler
	consumerHandler.pusher = pusher
	consumerHandler.pushConsumerGroup = kfk.NewMConsumerGroup(&kfk.MConsumerGroupConfig{
		KafkaVersion:   sarama.V2_0_0_0,
		OffsetsInitial: sarama.OffsetNewest, IsReturnErr: false,
	}, []string{config.Config.Kafka.MsgToPush.Topic}, config.Config.Kafka.Addr,
		config.Config.Kafka.ConsumerGroupID.MsgToPush)
	return &consumerHandler
}

func (c *ConsumerHandler) handleMs2PsChat(ctx context.Context, msg []byte) {
	msgFromMQ := pbChat.PushMsgDataToMQ{}
	if err := proto.Unmarshal(msg, &msgFromMQ); err != nil {
		log.ZError(ctx, "push Unmarshal msg err", err, "msg", string(msg))
		return
	}
	pbData := &pbPush.PushMsgReq{
		MsgData:        msgFromMQ.MsgData,
		ConversationID: msgFromMQ.ConversationID,
	}
	sec := msgFromMQ.MsgData.SendTime / 1000
	nowSec := utils.GetCurrentTimestampBySecond()
	if nowSec-sec > 10 {
		return
	}
	var err error
	switch msgFromMQ.MsgData.SessionType {
	case constant.SuperGroupChatType:
		err = c.pusher.Push2SuperGroup(ctx, pbData.MsgData.GroupID, pbData.MsgData)
	default:
		var pushUserIDs []string
		if pbData.MsgData.SendID != pbData.MsgData.RecvID {
			pushUserIDs = []string{pbData.MsgData.SendID, pbData.MsgData.RecvID}
		} else {
			pushUserIDs = []string{pbData.MsgData.SendID}
		}
		err = c.pusher.Push2User(ctx, pushUserIDs, pbData.MsgData)
	}
	if err != nil {
		if err == errNoOfflinePusher {
			log.ZWarn(ctx, "offline push failed", err, "msg", pbData.String())
		} else {
			log.ZError(ctx, "push failed", err, "msg", pbData.String())
		}
	}
}
func (ConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (ConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (c *ConsumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) error {
	for msg := range claim.Messages() {
		ctx := c.pushConsumerGroup.GetContextFromMsg(msg)
		c.handleMs2PsChat(ctx, msg.Value)
		sess.MarkMessage(msg, "")
	}
	return nil
}
