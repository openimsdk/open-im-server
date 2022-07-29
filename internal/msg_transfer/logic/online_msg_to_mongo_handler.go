package logic

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	kfk "Open_IM/pkg/common/kafka"
	"Open_IM/pkg/common/log"
	pbMsg "Open_IM/pkg/proto/msg"
	server_api_params "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
)

type OnlineHistoryMongoConsumerHandler struct {
	msgHandle            map[string]fcb
	historyConsumerGroup *kfk.MConsumerGroup
}

func (och *OnlineHistoryMongoConsumerHandler) Init() {
	och.msgHandle = make(map[string]fcb)
	och.msgHandle[config.Config.Kafka.MsgToMongo.Topic] = och.handleChatWs2Mongo
	och.historyConsumerGroup = kfk.NewMConsumerGroup(&kfk.MConsumerGroupConfig{KafkaVersion: sarama.V2_0_0_0,
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
	err = db.DB.BatchInsertChat2DB(msgFromMQ.AggregationID, msgFromMQ.MessageList, msgFromMQ.TriggerID, msgFromMQ.LastSeq)
	if err != nil {
		log.NewError(msgFromMQ.TriggerID, "single data insert to mongo err", err.Error(), msgFromMQ.MessageList, msgFromMQ.AggregationID, msgFromMQ.TriggerID)
	} else {
		err = db.DB.DeleteMessageFromCache(msgFromMQ.MessageList, msgFromMQ.AggregationID, msgFromMQ.GetTriggerID())
		if err != nil {
			log.NewError(msgFromMQ.TriggerID, "remove cache msg from redis  err", err.Error(), msgFromMQ.MessageList, msgFromMQ.AggregationID, msgFromMQ.TriggerID)
		}
	}
	for _, v := range msgFromMQ.MessageList {
		if v.MsgData.ContentType == constant.DeleteMessageNotification {
			tips := server_api_params.TipsComm{}
			DeleteMessageTips := server_api_params.DeleteMessageTips{}
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
			if unexistSeqList, err := db.DB.DelMsgBySeqList(DeleteMessageTips.UserID, DeleteMessageTips.SeqList, v.OperationID); err != nil {
				log.NewError(v.OperationID, utils.GetSelfFuncName(), "DelMsgBySeqList args: ", DeleteMessageTips.UserID, DeleteMessageTips.SeqList, v.OperationID, err.Error(), unexistSeqList)
			}

		}
	}
}

func (OnlineHistoryMongoConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (OnlineHistoryMongoConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (och *OnlineHistoryMongoConsumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim) error { // a instance in the consumer group
	log.NewDebug("", "online new session msg come", claim.HighWaterMarkOffset(), claim.Topic(), claim.Partition())
	for msg := range claim.Messages() {
		log.NewDebug("", "kafka get info to mongo", "msgTopic", msg.Topic, "msgPartition", msg.Partition, "msg", string(msg.Value), "key", string(msg.Key))
		if len(msg.Value) != 0 {
			och.msgHandle[msg.Topic](msg, string(msg.Key), sess)
		} else {
			log.Error("", "mongo msg get from kafka but is nil", msg.Key)
		}
		sess.MarkMessage(msg, "")
	}
	return nil
}
