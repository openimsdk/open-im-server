package msgtransfer

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db/cache"
	"Open_IM/pkg/common/db/controller"
	kfk "Open_IM/pkg/common/kafka"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/tracelog"
	pbMsg "Open_IM/pkg/proto/msg"
	"Open_IM/pkg/proto/sdkws"
	"Open_IM/pkg/utils"
	"context"
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
)

type OnlineHistoryMongoConsumerHandler struct {
	historyConsumerGroup *kfk.MConsumerGroup
	msgDatabase          controller.MsgDatabase
	cache                cache.Cache
}

func (mc *OnlineHistoryMongoConsumerHandler) Init() {
	mc.historyConsumerGroup = kfk.NewMConsumerGroup(&kfk.MConsumerGroupConfig{KafkaVersion: sarama.V2_0_0_0,
		OffsetsInitial: sarama.OffsetNewest, IsReturnErr: false}, []string{config.Config.Kafka.MsgToMongo.Topic},
		config.Config.Kafka.Ws2mschat.Addr, config.Config.Kafka.ConsumerGroupID.MsgToMongo)
}
func (mc *OnlineHistoryMongoConsumerHandler) handleChatWs2Mongo(cMsg *sarama.ConsumerMessage, msgKey string, _ sarama.ConsumerGroupSession) {
	msg := cMsg.Value
	msgFromMQ := pbMsg.MsgDataToMongoByMQ{}
	err := proto.Unmarshal(msg, &msgFromMQ)
	if err != nil {
		log.Error("msg_transfer Unmarshal msg err", "", "msg", string(msg), "err", err.Error())
		return
	}
	log.Info(msgFromMQ.TriggerID, "BatchInsertChat2DB userID: ", msgFromMQ.AggregationID, "msgFromMQ.LastSeq: ", msgFromMQ.LastSeq)
	ctx := context.Background()
	tracelog.SetOperationID(ctx, msgFromMQ.TriggerID)
	//err = db.DB.BatchInsertChat2DB(msgFromMQ.AggregationID, msgFromMQ.MessageList, msgFromMQ.TriggerID, msgFromMQ.LastSeq)
	err = mc.msgDatabase.BatchInsertChat2DB(ctx, msgFromMQ.AggregationID, msgFromMQ.MessageList, msgFromMQ.LastSeq)
	if err != nil {
		log.NewError(msgFromMQ.TriggerID, "single data insert to mongo err", err.Error(), msgFromMQ.MessageList, msgFromMQ.AggregationID, msgFromMQ.TriggerID)
	}
	//err = db.DB.DeleteMessageFromCache(msgFromMQ.MessageList, msgFromMQ.AggregationID, msgFromMQ.GetTriggerID())
	err = mc.msgDatabase.DeleteMessageFromCache(ctx, msgFromMQ.AggregationID, msgFromMQ.MessageList)
	if err != nil {
		log.NewError(msgFromMQ.TriggerID, "remove cache msg from redis err", err.Error(), msgFromMQ.MessageList, msgFromMQ.AggregationID, msgFromMQ.TriggerID)
	}
	for _, v := range msgFromMQ.MessageList {
		if v.MsgData.ContentType == constant.DeleteMessageNotification {
			tips := sdkws.TipsComm{}
			DeleteMessageTips := sdkws.DeleteMessageTips{}
			err := proto.Unmarshal(v.MsgData.Content, &tips)
			if err != nil {
				log.NewError(msgFromMQ.TriggerID, "tips unmarshal err:", err.Error(), v.String())
				continue
			}
			err = proto.Unmarshal(tips.Detail, &DeleteMessageTips)
			if err != nil {
				log.NewError(msgFromMQ.TriggerID, "deleteMessageTips unmarshal err:", err.Error(), v.String())
				continue
			}
			if totalUnExistSeqs, err := mc.msgDatabase.DelMsgBySeqs(ctx, DeleteMessageTips.UserID, DeleteMessageTips.Seqs); err != nil {
				log.NewError(v.OperationID, utils.GetSelfFuncName(), "DelMsgBySeqs args: ", DeleteMessageTips.UserID, DeleteMessageTips.Seqs, "error:", err.Error(), "totalUnExistSeqs: ", totalUnExistSeqs)
			}

		}
	}
}

func (OnlineHistoryMongoConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (OnlineHistoryMongoConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (mc *OnlineHistoryMongoConsumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim) error { // a instance in the consumer group
	log.NewDebug("", "online new session msg come", claim.HighWaterMarkOffset(), claim.Topic(), claim.Partition())
	for msg := range claim.Messages() {
		log.NewDebug("", "kafka get info to mongo", "msgTopic", msg.Topic, "msgPartition", msg.Partition, "msg", string(msg.Value), "key", string(msg.Key))
		if len(msg.Value) != 0 {
			mc.handleChatWs2Mongo(msg, string(msg.Key), sess)
		} else {
			log.Error("", "mongo msg get from kafka but is nil", msg.Key)
		}
		sess.MarkMessage(msg, "")
	}
	return nil
}
