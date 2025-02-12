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
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/prommetrics"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/controller"
	pbmsg "github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mq/kafka"
	"google.golang.org/protobuf/proto"
)

type OnlineHistoryMongoConsumerHandler struct {
	historyConsumerGroup *kafka.MConsumerGroup
	msgTransferDatabase  controller.MsgTransferDatabase
}

func NewOnlineHistoryMongoConsumerHandler(kafkaConf *config.Kafka, database controller.MsgTransferDatabase) (*OnlineHistoryMongoConsumerHandler, error) {
	historyConsumerGroup, err := kafka.NewMConsumerGroup(kafkaConf.Build(), kafkaConf.ToMongoGroupID, []string{kafkaConf.ToMongoTopic}, true)
	if err != nil {
		return nil, err
	}

	mc := &OnlineHistoryMongoConsumerHandler{
		historyConsumerGroup: historyConsumerGroup,
		msgTransferDatabase:  database,
	}
	return mc, nil
}

func (mc *OnlineHistoryMongoConsumerHandler) handleChatWs2Mongo(ctx context.Context, cMsg *sarama.ConsumerMessage, key string, session sarama.ConsumerGroupSession) {
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
	log.ZDebug(ctx, "mongo consumer recv msg", "msgs", msgFromMQ.String())
	err = mc.msgTransferDatabase.BatchInsertChat2DB(ctx, msgFromMQ.ConversationID, msgFromMQ.MsgData, msgFromMQ.LastSeq)
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
	//var seqs []int64
	//for _, msg := range msgFromMQ.MsgData {
	//	seqs = append(seqs, msg.Seq)
	//}
	//if err := mc.msgTransferDatabase.DeleteMessagesFromCache(ctx, msgFromMQ.ConversationID, seqs); err != nil {
	//	log.ZError(ctx, "remove cache msg from redis err", err, "msg",
	//		msgFromMQ.MsgData, "conversationID", msgFromMQ.ConversationID)
	//}
}

func (*OnlineHistoryMongoConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error { return nil }

func (*OnlineHistoryMongoConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (mc *OnlineHistoryMongoConsumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error { // an instance in the consumer group
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
