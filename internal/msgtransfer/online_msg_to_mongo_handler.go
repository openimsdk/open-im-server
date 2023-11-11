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

	"github.com/IBM/sarama"
	"google.golang.org/protobuf/proto"

	pbmsg "github.com/OpenIMSDK/protocol/msg"
	"github.com/OpenIMSDK/tools/log"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/controller"
	kfk "github.com/openimsdk/open-im-server/v3/pkg/common/kafka"
	"github.com/openimsdk/open-im-server/v3/pkg/common/prommetrics"
)

type OnlineHistoryMongoConsumerHandler struct {
	historyConsumerGroup *kfk.MConsumerGroup
	msgDatabase          controller.CommonMsgDatabase
}

func NewOnlineHistoryMongoConsumerHandler(database controller.CommonMsgDatabase) *OnlineHistoryMongoConsumerHandler {
	mc := &OnlineHistoryMongoConsumerHandler{
		historyConsumerGroup: kfk.NewMConsumerGroup(&kfk.MConsumerGroupConfig{
			KafkaVersion:   sarama.V2_0_0_0,
			OffsetsInitial: sarama.OffsetNewest, IsReturnErr: false,
		}, []string{config.Config.Kafka.MsgToMongo.Topic},
			config.Config.Kafka.Addr, config.Config.Kafka.ConsumerGroupID.MsgToMongo),
		msgDatabase: database,
	}
	return mc
}

func (mc *OnlineHistoryMongoConsumerHandler) handleChatWs2Mongo(
	ctx context.Context,
	cMsg *sarama.ConsumerMessage,
	key string,
	session sarama.ConsumerGroupSession,
) {
	msg := cMsg.Value
	msgFromMQ := pbmsg.MsgDataToMongoByMQ{}
	err := proto.Unmarshal(msg, &msgFromMQ)
	if err != nil {
		log.ZError(ctx, "unmarshall failed", err, "key", key, "len", len(msg))
		return
	}
	if len(msgFromMQ.MsgData) == 0 {
		log.ZError(ctx, "msgFromMQ.MsgData is empty", nil, "cMsg", cMsg)
		return
	}
	log.ZInfo(ctx, "mongo consumer recv msg", "msgs", msgFromMQ.String())
	err = mc.msgDatabase.BatchInsertChat2DB(ctx, msgFromMQ.ConversationID, msgFromMQ.MsgData, msgFromMQ.LastSeq)
	if err != nil {
		log.ZError(
			ctx,
			"single data insert to mongo err",
			err,
			"msg",
			msgFromMQ.MsgData,
			"conversationID",
			msgFromMQ.ConversationID,
		)
		prommetrics.MsgInsertMongoFailedCounter.Inc()
	} else {
		prommetrics.MsgInsertMongoSuccessCounter.Inc()
	}
	var seqs []int64
	for _, msg := range msgFromMQ.MsgData {
		seqs = append(seqs, msg.Seq)
	}
	err = mc.msgDatabase.DeleteMessagesFromCache(ctx, msgFromMQ.ConversationID, seqs)
	if err != nil {
		log.ZError(
			ctx,
			"remove cache msg from redis err",
			err,
			"msg",
			msgFromMQ.MsgData,
			"conversationID",
			msgFromMQ.ConversationID,
		)
	}
	mc.msgDatabase.DelUserDeleteMsgsList(ctx, msgFromMQ.ConversationID, seqs)
}

func (OnlineHistoryMongoConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (OnlineHistoryMongoConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (mc *OnlineHistoryMongoConsumerHandler) ConsumeClaim(
	sess sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim,
) error { // a instance in the consumer group
	log.ZDebug(context.Background(), "online new session msg come", "highWaterMarkOffset",
		claim.HighWaterMarkOffset(), "topic", claim.Topic(), "partition", claim.Partition())
	for msg := range claim.Messages() {
		ctx := mc.historyConsumerGroup.GetContextFromMsg(msg)
		if len(msg.Value) != 0 {
			mc.handleChatWs2Mongo(ctx, msg, string(msg.Key), sess)
		} else {
			log.ZError(ctx, "mongo msg get from kafka but is nil", nil, "conversationID", msg.Key)
		}
		sess.MarkMessage(msg, "")
	}
	return nil
}
