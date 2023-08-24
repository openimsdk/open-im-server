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

package msgtransfer

import (
	"context"

	"github.com/OpenIMSDK/protocol/constant"
	pbmsg "github.com/OpenIMSDK/protocol/msg"
	"github.com/OpenIMSDK/tools/log"
	"github.com/OpenIMSDK/tools/utils"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/controller"
	kfk "github.com/OpenIMSDK/Open-IM-Server/pkg/common/kafka"

	"github.com/Shopify/sarama"
	"google.golang.org/protobuf/proto"
)

type PersistentConsumerHandler struct {
	persistentConsumerGroup *kfk.MConsumerGroup
	chatLogDatabase         controller.ChatLogDatabase
}

func NewPersistentConsumerHandler(database controller.ChatLogDatabase) *PersistentConsumerHandler {
	return &PersistentConsumerHandler{
		persistentConsumerGroup: kfk.NewMConsumerGroup(&kfk.MConsumerGroupConfig{
			KafkaVersion:   sarama.V2_0_0_0,
			OffsetsInitial: sarama.OffsetNewest, IsReturnErr: false,
		}, []string{config.Config.Kafka.LatestMsgToRedis.Topic},
			config.Config.Kafka.Addr, config.Config.Kafka.ConsumerGroupID.MsgToMySql),
		chatLogDatabase: database,
	}
}

func (pc *PersistentConsumerHandler) handleChatWs2Mysql(
	ctx context.Context,
	cMsg *sarama.ConsumerMessage,
	msgKey string,
	_ sarama.ConsumerGroupSession,
) {
	msg := cMsg.Value
	var tag bool
	msgFromMQ := pbmsg.MsgDataToMQ{}
	err := proto.Unmarshal(msg, &msgFromMQ)
	if err != nil {
		log.ZError(ctx, "msg_transfer Unmarshal msg err", err)
		return
	}
	return
	log.ZDebug(ctx, "handleChatWs2Mysql", "msg", msgFromMQ.MsgData)
	// Control whether to store history messages (mysql)
	isPersist := utils.GetSwitchFromOptions(msgFromMQ.MsgData.Options, constant.IsPersistent)
	// Only process receiver data
	if isPersist {
		switch msgFromMQ.MsgData.SessionType {
		case constant.SingleChatType, constant.NotificationChatType:
			if msgKey == msgFromMQ.MsgData.RecvID {
				tag = true
			}
		case constant.GroupChatType:
			if msgKey == msgFromMQ.MsgData.SendID {
				tag = true
			}
		case constant.SuperGroupChatType:
			tag = true
		}
		if tag {
			log.ZInfo(ctx, "msg_transfer msg persisting", "msg", string(msg))
			if err = pc.chatLogDatabase.CreateChatLog(&msgFromMQ); err != nil {
				log.ZError(ctx, "Message insert failed", err, "msg", msgFromMQ.String())
				return
			}
		}
	}
}
func (PersistentConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (PersistentConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (pc *PersistentConsumerHandler) ConsumeClaim(
	sess sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) error {
	for msg := range claim.Messages() {
		ctx := pc.persistentConsumerGroup.GetContextFromMsg(msg)
		log.ZDebug(
			ctx,
			"kafka get info to mysql",
			"msgTopic",
			msg.Topic,
			"msgPartition",
			msg.Partition,
			"msg",
			string(msg.Value),
			"key",
			string(msg.Key),
		)
		if len(msg.Value) != 0 {
			pc.handleChatWs2Mysql(ctx, msg, string(msg.Key), sess)
		} else {
			log.ZError(ctx, "msg get from kafka but is nil", nil, "key", msg.Key)
		}
		sess.MarkMessage(msg, "")
	}
	return nil
}
