package logic

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	kfk "Open_IM/pkg/common/kafka"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbMsg "Open_IM/pkg/proto/chat"
	pbPush "Open_IM/pkg/proto/push"
	"Open_IM/pkg/utils"
	"context"
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	"strings"
)

type fcb func(msg []byte, msgKey string)

type HistoryConsumerHandler struct {
	msgHandle            map[string]fcb
	historyConsumerGroup *kfk.MConsumerGroup
}

func (mc *HistoryConsumerHandler) Init() {
	mc.msgHandle = make(map[string]fcb)
	mc.msgHandle[config.Config.Kafka.Ws2mschat.Topic] = mc.handleChatWs2Mongo
	mc.historyConsumerGroup = kfk.NewMConsumerGroup(&kfk.MConsumerGroupConfig{KafkaVersion: sarama.V0_10_2_0,
		OffsetsInitial: sarama.OffsetNewest, IsReturnErr: false}, []string{config.Config.Kafka.Ws2mschat.Topic},
		config.Config.Kafka.Ws2mschat.Addr, config.Config.Kafka.ConsumerGroupID.MsgToMongo)

}

func (mc *HistoryConsumerHandler) handleChatWs2Mongo(msg []byte, msgKey string) {
	log.NewInfo("msg come mongo!!!", "", "msg", string(msg))
	time := utils.GetCurrentTimestampByNano()
	msgFromMQ := pbMsg.MsgDataToMQ{}
	err := proto.Unmarshal(msg, &msgFromMQ)
	if err != nil {
		log.ErrorByKv("msg_transfer Unmarshal msg err", "", "msg", string(msg), "err", err.Error())
		return
	}
	operationID := msgFromMQ.OperationID
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
				log.NewError(operationID, "single data insert to mongo err", err.Error(), msgFromMQ.String())
				return
			}
			log.NewDebug(msgFromMQ.OperationID, "sendMessageToPush cost time ", utils.GetCurrentTimestampByNano()-time)
		}
		if !isSenderSync && msgKey == msgFromMQ.MsgData.SendID {
		} else {
			go sendMessageToPush(&msgFromMQ, msgKey)
		}
		log.NewDebug(operationID, "saveUserChat cost time ", utils.GetCurrentTimestampByNano()-time)
	case constant.GroupChatType:
		log.NewDebug(msgFromMQ.OperationID, "msg_transfer msg type = GroupChatType", isHistory, isPersist)
		if isHistory {
			err := saveUserChat(msgFromMQ.MsgData.RecvID, &msgFromMQ)
			if err != nil {
				log.NewError(operationID, "group data insert to mongo err", msgFromMQ.String(), msgFromMQ.MsgData.RecvID, err.Error())
				return
			}
		}
		go sendMessageToPush(&msgFromMQ, msgFromMQ.MsgData.RecvID)
	default:
		log.NewError(msgFromMQ.OperationID, "SessionType error", msgFromMQ.String())
		return
	}
	log.NewDebug(msgFromMQ.OperationID, "msg_transfer handle topic data to database success...", msgFromMQ.String())
}

func (HistoryConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (HistoryConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (mc *HistoryConsumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.InfoByKv("kafka get info to mongo", "", "msgTopic", msg.Topic, "msgPartition", msg.Partition, "msg", string(msg.Value))
		mc.msgHandle[msg.Topic](msg.Value, string(msg.Key))
		sess.MarkMessage(msg, "")
	}
	return nil
}
func sendMessageToPush(message *pbMsg.MsgDataToMQ, pushToUserID string) {
	log.InfoByKv("msg_transfer send message to push", message.OperationID, "message", message.String())
	rpcPushMsg := pbPush.PushMsgReq{OperationID: message.OperationID, MsgData: message.MsgData, PushToUserID: pushToUserID}
	mqPushMsg := pbMsg.PushMsgDataToMQ{OperationID: message.OperationID, MsgData: message.MsgData, PushToUserID: pushToUserID}
	grpcConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImPushName)
	if grpcConn == nil {
		log.ErrorByKv("rpc dial failed", rpcPushMsg.OperationID, "push data", rpcPushMsg.String())
		pid, offset, err := producer.SendMessage(&mqPushMsg)
		if err != nil {
			log.ErrorByKv("kafka send failed", mqPushMsg.OperationID, "send data", message.String(), "pid", pid, "offset", offset, "err", err.Error())
		}
		return
	}
	msgClient := pbPush.NewPushMsgServiceClient(grpcConn)
	_, err := msgClient.PushMsg(context.Background(), &rpcPushMsg)
	if err != nil {
		log.ErrorByKv("rpc send failed", rpcPushMsg.OperationID, "push data", rpcPushMsg.String(), "err", err.Error())
		pid, offset, err := producer.SendMessage(&mqPushMsg)
		if err != nil {
			log.ErrorByKv("kafka send failed", mqPushMsg.OperationID, "send data", mqPushMsg.String(), "pid", pid, "offset", offset, "err", err.Error())
		}
	} else {
		log.InfoByKv("rpc send success", rpcPushMsg.OperationID, "push data", rpcPushMsg.String())

	}
}
