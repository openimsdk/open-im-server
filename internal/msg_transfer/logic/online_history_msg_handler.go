package logic

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	kfk "Open_IM/pkg/common/kafka"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbMsg "Open_IM/pkg/proto/chat"
	pbPush "Open_IM/pkg/proto/push"
	"Open_IM/pkg/utils"
	"context"
	"errors"
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	"strings"
	"time"
)

type MsgChannelValue struct {
	userID string
	msg    pbMsg.MsgDataToMQ
}
type fcb func(cMsg *sarama.ConsumerMessage, msgKey string, sess sarama.ConsumerGroupSession)
type Cmd2Value struct {
	Cmd   int
	Value interface{}
}
type OnlineHistoryConsumerHandler struct {
	msgHandle            map[string]fcb
	historyConsumerGroup *kfk.MConsumerGroup
	cmdCh                chan Cmd2Value
	msgCh                chan Cmd2Value
}

func (och *OnlineHistoryConsumerHandler) Init(cmdCh chan Cmd2Value) {
	och.msgHandle = make(map[string]fcb)
	och.cmdCh = cmdCh
	och.msgCh = make(chan Cmd2Value, 1000)
	if config.Config.ReliableStorage {
		och.msgHandle[config.Config.Kafka.Ws2mschat.Topic] = och.handleChatWs2Mongo
	} else {
		och.msgHandle[config.Config.Kafka.Ws2mschat.Topic] = och.handleChatWs2MongoLowReliability
		for i := 0; i < 10; i++ {
			go och.Run()

		}
	}
	och.historyConsumerGroup = kfk.NewMConsumerGroup(&kfk.MConsumerGroupConfig{KafkaVersion: sarama.V0_10_2_0,
		OffsetsInitial: sarama.OffsetNewest, IsReturnErr: false}, []string{config.Config.Kafka.Ws2mschat.Topic},
		config.Config.Kafka.Ws2mschat.Addr, config.Config.Kafka.ConsumerGroupID.MsgToMongo)

}
func (och *OnlineHistoryConsumerHandler) TriggerCmd(status int) {
	operationID := utils.OperationIDGenerator()
	err := sendCmd(och.cmdCh, Cmd2Value{Cmd: status, Value: ""}, 1)
	if err != nil {
		log.Error(operationID, "TriggerCmd failed ", err.Error(), status)
		return
	}
	log.Debug(operationID, "TriggerCmd success", status)

}
func sendCmd(ch chan Cmd2Value, value Cmd2Value, timeout int64) error {
	var flag = 0
	select {
	case ch <- value:
		flag = 1
	case <-time.After(time.Second * time.Duration(timeout)):
		flag = 2
	}
	if flag == 1 {
		return nil
	} else {
		return errors.New("send cmd timeout")
	}
}
func (och *OnlineHistoryConsumerHandler) Run() {
	for {
		select {
		case cmd := <-och.msgCh:
			switch cmd.Cmd {
			case Msg:
				msgChannelValue := cmd.Value.(MsgChannelValue)
				msg := msgChannelValue.msg
				log.Debug(msg.OperationID, "msg arrived channel", msg.String())
				isSenderSync := utils.GetSwitchFromOptions(msg.MsgData.Options, constant.IsSenderSync)
				//switch msgChannelValue.msg.MsgData.SessionType {
				//case constant.SingleChatType:
				//case constant.GroupChatType:
				//case constant.NotificationChatType:
				//default:
				//	log.NewError(msgFromMQ.OperationID, "SessionType error", msgFromMQ.String())
				//	return
				//}

				err := saveUserChat(msgChannelValue.userID, &msg)
				if err != nil {
					singleMsgFailedCount++
					log.NewError(msg.OperationID, "single data insert to mongo err", err.Error(), msg.String())
				} else {
					singleMsgSuccessCountMutex.Lock()
					singleMsgSuccessCount++
					singleMsgSuccessCountMutex.Unlock()
					if !(!isSenderSync && msgChannelValue.userID == msg.MsgData.SendID) {
						go sendMessageToPush(&msg, msgChannelValue.userID)
					}
				}
			}
		}
	}
}
func (och *OnlineHistoryConsumerHandler) handleChatWs2Mongo(cMsg *sarama.ConsumerMessage, msgKey string, sess sarama.ConsumerGroupSession) {
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
	sess.MarkMessage(cMsg, "")
	log.NewDebug(msgFromMQ.OperationID, "msg_transfer handle topic data to database success...", msgFromMQ.String())
}
func (och *OnlineHistoryConsumerHandler) handleChatWs2MongoLowReliability(cMsg *sarama.ConsumerMessage, msgKey string, sess sarama.ConsumerGroupSession) {
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
		och.msgCh <- Cmd2Value{Cmd: Msg, Value: MsgChannelValue{msgKey, msgFromMQ}}
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

func (OnlineHistoryConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (OnlineHistoryConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (och *OnlineHistoryConsumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim) error { // a instance in the consumer group
	log.NewDebug("", "online new session msg come", claim.HighWaterMarkOffset(), claim.Topic(), claim.Partition())
	for msg := range claim.Messages() {
		SetOnlineTopicStatus(OnlineTopicBusy)
		//och.TriggerCmd(OnlineTopicBusy)
		log.NewDebug("", "online kafka get info to mongo", "msgTopic", msg.Topic, "msgPartition", msg.Partition, "online", msg.Offset, claim.HighWaterMarkOffset())
		och.msgHandle[msg.Topic](msg, string(msg.Key), sess)
		if claim.HighWaterMarkOffset()-msg.Offset <= 1 {
			log.Debug("", "online msg consume end", claim.HighWaterMarkOffset(), msg.Offset)
			SetOnlineTopicStatus(OnlineTopicVacancy)
			och.TriggerCmd(OnlineTopicVacancy)
		}
	}
	return nil
}
func sendMessageToPush(message *pbMsg.MsgDataToMQ, pushToUserID string) {
	log.Info(message.OperationID, "msg_transfer send message to push", "message", message.String())
	rpcPushMsg := pbPush.PushMsgReq{OperationID: message.OperationID, MsgData: message.MsgData, PushToUserID: pushToUserID}
	mqPushMsg := pbMsg.PushMsgDataToMQ{OperationID: message.OperationID, MsgData: message.MsgData, PushToUserID: pushToUserID}
	grpcConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImPushName)
	if grpcConn == nil {
		log.Error(rpcPushMsg.OperationID, "rpc dial failed", "push data", rpcPushMsg.String())
		pid, offset, err := producer.SendMessage(&mqPushMsg)
		if err != nil {
			log.Error(mqPushMsg.OperationID, "kafka send failed", "send data", message.String(), "pid", pid, "offset", offset, "err", err.Error())
		}
		return
	}
	msgClient := pbPush.NewPushMsgServiceClient(grpcConn)
	_, err := msgClient.PushMsg(context.Background(), &rpcPushMsg)
	if err != nil {
		log.Error(rpcPushMsg.OperationID, "rpc send failed", rpcPushMsg.OperationID, "push data", rpcPushMsg.String(), "err", err.Error())
		pid, offset, err := producer.SendMessage(&mqPushMsg)
		if err != nil {
			log.Error(message.OperationID, "kafka send failed", mqPushMsg.OperationID, "send data", mqPushMsg.String(), "pid", pid, "offset", offset, "err", err.Error())
		}
	} else {
		log.Info(message.OperationID, "rpc send success", rpcPushMsg.OperationID, "push data", rpcPushMsg.String())

	}
}
