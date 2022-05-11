package logic

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	kfk "Open_IM/pkg/common/kafka"
	"Open_IM/pkg/common/log"
	pbMsg "Open_IM/pkg/proto/chat"
	"Open_IM/pkg/utils"
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	"time"
)

type OfflineHistoryConsumerHandler struct {
	msgHandle            map[string]fcb
	cmdCh                chan Cmd2Value
	historyConsumerGroup *kfk.MConsumerGroup
}

func (mc *OfflineHistoryConsumerHandler) Init(cmdCh chan Cmd2Value) {
	mc.msgHandle = make(map[string]fcb)
	mc.cmdCh = cmdCh
	mc.msgHandle[config.Config.Kafka.Ws2mschatOffline.Topic] = mc.handleChatWs2Mongo
	mc.historyConsumerGroup = kfk.NewMConsumerGroup(&kfk.MConsumerGroupConfig{KafkaVersion: sarama.V0_10_2_0,
		OffsetsInitial: sarama.OffsetNewest, IsReturnErr: false}, []string{config.Config.Kafka.Ws2mschatOffline.Topic},
		config.Config.Kafka.Ws2mschatOffline.Addr, config.Config.Kafka.ConsumerGroupID.MsgToMongoOffline)

}

func (mc *OfflineHistoryConsumerHandler) handleChatWs2Mongo(msg []byte, msgKey string) {
	now := time.Now()
	msgFromMQ := pbMsg.MsgDataToMQ{}
	err := proto.Unmarshal(msg, &msgFromMQ)
	if err != nil {
		log.Error("msg_transfer Unmarshal msg err", "", "msg", string(msg), "err", err.Error())
		return
	}
	operationID := msgFromMQ.OperationID
	log.NewInfo(operationID, "msg come mongo!!!", "", "msg", string(msg))
	//Control whether to store offline messages (mongo)
	isHistory := utils.GetSwitchFromOptions(msgFromMQ.MsgData.Options, constant.IsHistory)
	//Control whether to store history messages (mysql)
	isPersist := utils.GetSwitchFromOptions(msgFromMQ.MsgData.Options, constant.IsPersistent)
	isSenderSync := utils.GetSwitchFromOptions(msgFromMQ.MsgData.Options, constant.IsSenderSync)
	switch msgFromMQ.MsgData.SessionType {
	case constant.SingleChatType:
		log.NewDebug(msgFromMQ.OperationID, "msg_transfer msg type = SingleChatType", isHistory, isPersist)
		if isHistory {
			err := saveUserChat(msgKey, &msgFromMQ)
			if err != nil {
				singleMsgFailedCount++
				log.NewError(operationID, "single data insert to mongo err", err.Error(), msgFromMQ.String())
				return
			}
			singleMsgSuccessCount++
			log.NewDebug(msgFromMQ.OperationID, "sendMessageToPush cost time ", time.Since(now))
		}
		if !isSenderSync && msgKey == msgFromMQ.MsgData.SendID {
		} else {
			go sendMessageToPush(&msgFromMQ, msgKey)
		}
		log.NewDebug(operationID, "saveUserChat cost time ", time.Since(now))
	case constant.GroupChatType:
		log.NewDebug(msgFromMQ.OperationID, "msg_transfer msg type = GroupChatType", isHistory, isPersist)
		if isHistory {
			err := saveUserChat(msgFromMQ.MsgData.RecvID, &msgFromMQ)
			if err != nil {
				log.NewError(operationID, "group data insert to mongo err", msgFromMQ.String(), msgFromMQ.MsgData.RecvID, err.Error())
				return
			}
			groupMsgCount++
		}
		go sendMessageToPush(&msgFromMQ, msgFromMQ.MsgData.RecvID)
	case constant.NotificationChatType:
		log.NewDebug(msgFromMQ.OperationID, "msg_transfer msg type = NotificationChatType", isHistory, isPersist)
		if isHistory {
			err := saveUserChat(msgKey, &msgFromMQ)
			if err != nil {
				log.NewError(operationID, "single data insert to mongo err", err.Error(), msgFromMQ.String())
				return
			}
			log.NewDebug(msgFromMQ.OperationID, "sendMessageToPush cost time ", time.Since(now))
		}
		if !isSenderSync && msgKey == msgFromMQ.MsgData.SendID {
		} else {
			go sendMessageToPush(&msgFromMQ, msgKey)
		}
		log.NewDebug(operationID, "saveUserChat cost time ", time.Since(now))
	default:
		log.NewError(msgFromMQ.OperationID, "SessionType error", msgFromMQ.String())
		return
	}
	log.NewDebug(msgFromMQ.OperationID, "msg_transfer handle topic data to database success...", msgFromMQ.String())
}

func (OfflineHistoryConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (OfflineHistoryConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (mc *OfflineHistoryConsumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim) error { // a instance in the consumer group
	//log.NewDebug("", "offline new session msg come", claim.HighWaterMarkOffset(), claim.Topic(), claim.Partition())
	//for msg := range claim.Messages() {
	//	log.NewDebug("", "kafka get info to delay mongo", "msgTopic", msg.Topic, "msgPartition", msg.Partition, "offline")
	//	//mc.msgHandle[msg.Topic](msg.Value, string(msg.Key))
	//}
	for msg := range claim.Messages() {
		if GetOnlineTopicStatus() == OnlineTopicVacancy {
			log.NewDebug("", "vacancy offline kafka get info to mongo", "msgTopic", msg.Topic, "msgPartition", msg.Partition, "msg", string(msg.Value))
			mc.msgHandle[msg.Topic](msg.Value, string(msg.Key))
			sess.MarkMessage(msg, "")
		} else {
			select {
			case <-mc.cmdCh:
				log.NewDebug("", "cmd offline kafka get info to mongo", "msgTopic", msg.Topic, "msgPartition", msg.Partition, "msg", string(msg.Value))
			case <-time.After(time.Millisecond * time.Duration(100)):
				log.NewDebug("", "timeout offline kafka get info to mongo", "msgTopic", msg.Topic, "msgPartition", msg.Partition, "msg", string(msg.Value))
			}
			mc.msgHandle[msg.Topic](msg.Value, string(msg.Key))
			sess.MarkMessage(msg, "")
		}
	}

	return nil
}
