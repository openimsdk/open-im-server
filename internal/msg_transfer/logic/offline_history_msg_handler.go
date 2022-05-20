package logic

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
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
	historyConsumerGroup *kfk.MConsumerGroup
	cmdCh                chan Cmd2Value
	msgCh                chan Cmd2Value
	chArrays             [ChannelNum]chan Cmd2Value
	msgDistributionCh    chan Cmd2Value
}

func (mc *OfflineHistoryConsumerHandler) Init(cmdCh chan Cmd2Value) {
	mc.msgHandle = make(map[string]fcb)
	mc.msgDistributionCh = make(chan Cmd2Value) //no buffer channel
	go mc.MessagesDistributionHandle()
	mc.cmdCh = cmdCh
	mc.msgCh = make(chan Cmd2Value, 1000)
	if config.Config.ReliableStorage {
		mc.msgHandle[config.Config.Kafka.Ws2mschat.Topic] = mc.handleChatWs2Mongo
	} else {
		mc.msgHandle[config.Config.Kafka.Ws2mschat.Topic] = mc.handleChatWs2MongoLowReliability
		for i := 0; i < ChannelNum; i++ {
			mc.chArrays[i] = make(chan Cmd2Value, 1000)
			go mc.Run(i)
		}
	}
	mc.historyConsumerGroup = kfk.NewMConsumerGroup(&kfk.MConsumerGroupConfig{KafkaVersion: sarama.V0_10_2_0,
		OffsetsInitial: sarama.OffsetNewest, IsReturnErr: false}, []string{config.Config.Kafka.Ws2mschatOffline.Topic},
		config.Config.Kafka.Ws2mschatOffline.Addr, config.Config.Kafka.ConsumerGroupID.MsgToMongoOffline)

}
func (och *OfflineHistoryConsumerHandler) Run(channelID int) {
	for {
		select {
		case cmd := <-och.chArrays[channelID]:
			switch cmd.Cmd {
			case UserMessages:
				msgChannelValue := cmd.Value.(MsgChannelValue)
				msgList := msgChannelValue.msgList
				triggerID := msgChannelValue.triggerID
				storageMsgList := make([]*pbMsg.MsgDataToMQ, 80)
				pushMsgList := make([]*pbMsg.MsgDataToMQ, 80)
				log.Debug(triggerID, "msg arrived channel", "channel id", channelID, msgList, msgChannelValue.userID, len(msgList))
				for _, v := range msgList {
					log.Debug(triggerID, "msg come to storage center", v.String())
					isHistory := utils.GetSwitchFromOptions(v.MsgData.Options, constant.IsHistory)
					isSenderSync := utils.GetSwitchFromOptions(v.MsgData.Options, constant.IsSenderSync)
					if isHistory {
						storageMsgList = append(storageMsgList, v)
					}
					if !(!isSenderSync && msgChannelValue.userID == v.MsgData.SendID) {
						pushMsgList = append(pushMsgList, v)
					}
				}

				//switch msgChannelValue.msg.MsgData.SessionType {
				//case constant.SingleChatType:
				//case constant.GroupChatType:
				//case constant.NotificationChatType:
				//default:
				//	log.NewError(msgFromMQ.OperationID, "SessionType error", msgFromMQ.String())
				//	return
				//}

				err := saveUserChatList(msgChannelValue.userID, storageMsgList, triggerID)
				if err != nil {
					singleMsgFailedCount += uint64(len(storageMsgList))
					log.NewError(triggerID, "single data insert to mongo err", err.Error(), storageMsgList)
				} else {
					singleMsgSuccessCountMutex.Lock()
					singleMsgSuccessCount += uint64(len(storageMsgList))
					singleMsgSuccessCountMutex.Unlock()
					for _, v := range pushMsgList {
						sendMessageToPush(v, msgChannelValue.userID)
					}

				}
			}
		}
	}
}
func (och *OfflineHistoryConsumerHandler) MessagesDistributionHandle() {
	UserAggregationMsgs := make(map[string][]*pbMsg.MsgDataToMQ, ChannelNum)
	for {
		select {
		case cmd := <-och.msgDistributionCh:
			switch cmd.Cmd {
			case ConsumerMsgs:
				triggerChannelValue := cmd.Value.(TriggerChannelValue)
				triggerID := triggerChannelValue.triggerID
				consumerMessages := triggerChannelValue.cmsgList
				//Aggregation map[userid]message list
				log.Debug(triggerID, "batch messages come to distribution center", len(consumerMessages))
				for i := 0; i < len(consumerMessages); i++ {
					msgFromMQ := pbMsg.MsgDataToMQ{}
					err := proto.Unmarshal(consumerMessages[i].Value, &msgFromMQ)
					if err != nil {
						log.Error(triggerID, "msg_transfer Unmarshal msg err", "msg", string(consumerMessages[i].Value), "err", err.Error())
						return
					}
					log.Debug(triggerID, "single msg come to distribution center", msgFromMQ.String())
					if oldM, ok := UserAggregationMsgs[string(consumerMessages[i].Key)]; ok {
						oldM = append(oldM, &msgFromMQ)
						UserAggregationMsgs[string(consumerMessages[i].Key)] = oldM
					} else {
						m := make([]*pbMsg.MsgDataToMQ, 0, 100)
						m = append(m, &msgFromMQ)
						UserAggregationMsgs[string(consumerMessages[i].Key)] = m
					}
				}
				log.Debug(triggerID, "generate map list users len", len(UserAggregationMsgs))
				for userID, v := range UserAggregationMsgs {
					if len(v) >= 0 {
						channelID := getHashCode(userID) % ChannelNum
						go func(cID uint32, userID string, messages []*pbMsg.MsgDataToMQ) {
							och.chArrays[cID] <- Cmd2Value{Cmd: UserMessages, Value: MsgChannelValue{userID: userID, msgList: messages, triggerID: triggerID}}
						}(channelID, userID, v)
					}
				}
			}
		}
	}

}
func (mc *OfflineHistoryConsumerHandler) handleChatWs2Mongo(cMsg *sarama.ConsumerMessage, msgKey string, sess sarama.ConsumerGroupSession) {
	msg := cMsg.Value
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
			singleMsgSuccessCountMutex.Lock()
			singleMsgSuccessCount++
			singleMsgSuccessCountMutex.Unlock()
			log.NewDebug(msgFromMQ.OperationID, "sendMessageToPush cost time ", time.Since(now))
		}
		if !isSenderSync && msgKey == msgFromMQ.MsgData.SendID {
		} else {
			go sendMessageToPush(&msgFromMQ, msgKey)
		}
		log.NewDebug(operationID, "saveSingleMsg cost time ", time.Since(now))
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
		log.NewDebug(operationID, "saveGroupMsg cost time ", time.Since(now))

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
	sess.MarkMessage(cMsg, "")
	log.NewDebug(msgFromMQ.OperationID, "msg_transfer handle topic data to database success...", msgFromMQ.String())
}
func (mc *OfflineHistoryConsumerHandler) handleChatWs2MongoLowReliability(cMsg *sarama.ConsumerMessage, msgKey string, sess sarama.ConsumerGroupSession) {
	msg := cMsg.Value
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
	isSenderSync := utils.GetSwitchFromOptions(msgFromMQ.MsgData.Options, constant.IsSenderSync)
	if isHistory {
		seq, err := db.DB.IncrUserSeq(msgKey)
		if err != nil {
			log.NewError(operationID, "data insert to redis err", err.Error(), string(msg))
			return
		}
		sess.MarkMessage(cMsg, "")
		msgFromMQ.MsgData.Seq = uint32(seq)
		log.Debug(operationID, "send ch msg is ", msgFromMQ.String())
		//mc.msgCh <- Cmd2Value{Cmd: Msg, Value: MsgChannelValue{msgKey, msgFromMQ}}
		//err := saveUserChat(msgKey, &msgFromMQ)
		//if err != nil {
		//	singleMsgFailedCount++
		//	log.NewError(operationID, "single data insert to mongo err", err.Error(), msgFromMQ.String())
		//	return
		//}
		//singleMsgSuccessCountMutex.Lock()
		//singleMsgSuccessCount++
		//singleMsgSuccessCountMutex.Unlock()
		//log.NewDebug(msgFromMQ.OperationID, "sendMessageToPush cost time ", time.Since(now))
	} else {
		if !(!isSenderSync && msgKey == msgFromMQ.MsgData.SendID) {
			go sendMessageToPush(&msgFromMQ, msgKey)
		}
	}
}

func (OfflineHistoryConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (OfflineHistoryConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

//func (mc *OfflineHistoryConsumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession,
//	claim sarama.ConsumerGroupClaim) error { // a instance in the consumer group
//	//log.NewDebug("", "offline new session msg come", claim.HighWaterMarkOffset(), claim.Topic(), claim.Partition())
//	//for msg := range claim.Messages() {
//	//	log.NewDebug("", "kafka get info to delay mongo", "msgTopic", msg.Topic, "msgPartition", msg.Partition, "offline")
//	//	//mc.msgHandle[msg.Topic](msg.Value, string(msg.Key))
//	//}
//	for msg := range claim.Messages() {
//		if GetOnlineTopicStatus() == OnlineTopicVacancy {
//			log.NewDebug("", "vacancy offline kafka get info to mongo", "msgTopic", msg.Topic, "msgPartition", msg.Partition, "msg", string(msg.Value))
//			mc.msgHandle[msg.Topic](msg, string(msg.Key), sess)
//		} else {
//			select {
//			case <-mc.cmdCh:
//				log.NewDebug("", "cmd offline kafka get info to mongo", "msgTopic", msg.Topic, "msgPartition", msg.Partition, "msg", string(msg.Value))
//			case <-time.After(time.Millisecond * time.Duration(100)):
//				log.NewDebug("", "timeout offline kafka get info to mongo", "msgTopic", msg.Topic, "msgPartition", msg.Partition, "msg", string(msg.Value))
//			}
//			mc.msgHandle[msg.Topic](msg, string(msg.Key), sess)
//		}
//	}
//
//	return nil
//}
func (och *OfflineHistoryConsumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim) error { // a instance in the consumer group
	log.NewDebug("", "online new session msg come", claim.HighWaterMarkOffset(), claim.Topic(), claim.Partition())
	cMsg := make([]*sarama.ConsumerMessage, 0, 500)
	t := time.NewTicker(time.Duration(500) * time.Millisecond)
	var triggerID string
	for msg := range claim.Messages() {
		//och.TriggerCmd(OnlineTopicBusy)
		cMsg = append(cMsg, msg)
		select {
		case <-t.C:
			if len(cMsg) >= 0 {
				triggerID = utils.OperationIDGenerator()
				log.Debug(triggerID, "timer trigger msg consumer start", len(cMsg))
				och.msgDistributionCh <- Cmd2Value{Cmd: ConsumerMsgs, Value: TriggerChannelValue{
					triggerID: triggerID, cmsgList: cMsg}}
				sess.MarkMessage(msg, "")
				cMsg = cMsg[0:0]
				log.Debug(triggerID, "timer trigger msg consumer end", len(cMsg))
			}
		default:
			if len(cMsg) >= 500 {
				triggerID = utils.OperationIDGenerator()
				log.Debug(triggerID, "length trigger msg consumer start", len(cMsg))
				och.msgDistributionCh <- Cmd2Value{Cmd: ConsumerMsgs, Value: TriggerChannelValue{
					triggerID: triggerID, cmsgList: cMsg}}
				sess.MarkMessage(msg, "")
				cMsg = cMsg[0:0]
				log.Debug(triggerID, "length trigger msg consumer end", len(cMsg))
			}

		}
		log.NewDebug("", "online kafka get info to mongo", "msgTopic", msg.Topic, "msgPartition", msg.Partition, "online", msg.Offset, claim.HighWaterMarkOffset())

	}
	return nil
}
