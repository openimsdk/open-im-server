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
	log.InfoByKv("chat come mongo!!!", "", "chat", string(msg))
	time := utils.GetCurrentTimestampByNano()
	pbData := pbMsg.WSToMsgSvrChatMsg{}
	err := proto.Unmarshal(msg, &pbData)
	if err != nil {
		log.ErrorByKv("msg_transfer Unmarshal chat err", "", "chat", string(msg), "err", err.Error())
		return
	}
	pbSaveData := pbMsg.MsgSvrToPushSvrChatMsg{}
	pbSaveData.SendID = pbData.SendID
	pbSaveData.SenderNickName = pbData.SenderNickName
	pbSaveData.SenderFaceURL = pbData.SenderFaceURL
	pbSaveData.ClientMsgID = pbData.ClientMsgID
	pbSaveData.SendTime = pbData.SendTime
	pbSaveData.Content = pbData.Content
	pbSaveData.MsgFrom = pbData.MsgFrom
	pbSaveData.ContentType = pbData.ContentType
	pbSaveData.SessionType = pbData.SessionType
	pbSaveData.MsgID = pbData.MsgID
	pbSaveData.OperationID = pbData.OperationID
	pbSaveData.RecvID = pbData.RecvID
	pbSaveData.PlatformID = pbData.PlatformID
	Options := utils.JsonStringToMap(pbData.Options)
	//Control whether to store offline messages (mongo)
	isHistory := utils.GetSwitchFromOptions(Options, "history")
	//Control whether to store history messages (mysql)
	isPersist := utils.GetSwitchFromOptions(Options, "persistent")
	switch pbData.SessionType {
	case constant.SingleChatType:
		log.NewDebug(pbSaveData.OperationID, "msg_transfer chat type = SingleChatType", isHistory, isPersist)
		if isHistory {
			if msgKey == pbSaveData.RecvID {
				err := saveUserChat(pbData.RecvID, &pbSaveData)
				if err != nil {
					log.NewError(pbSaveData.OperationID, "single data insert to mongo err", err.Error(), pbSaveData.String())
					return
				}

			} else if msgKey == pbSaveData.SendID {
				err := saveUserChat(pbData.SendID, &pbSaveData)
				if err != nil {
					log.NewError(pbSaveData.OperationID, "single data insert to mongo err", err.Error(), pbSaveData.String())
					return
				}

			}

			log.NewDebug(pbSaveData.OperationID, "saveUserChat cost time ", utils.GetCurrentTimestampByNano()-time)
		}
		if msgKey == pbSaveData.RecvID {
			pbSaveData.Options = pbData.Options
			pbSaveData.OfflineInfo = pbData.OfflineInfo
			go sendMessageToPush(&pbSaveData)
			log.NewDebug(pbSaveData.OperationID, "sendMessageToPush cost time ", utils.GetCurrentTimestampByNano()-time)
		}

	case constant.GroupChatType:
		log.NewDebug(pbSaveData.OperationID, "msg_transfer chat type = GroupChatType", isHistory, isPersist)
		if isHistory {
			uidAndGroupID := strings.Split(pbData.RecvID, " ")
			err := saveUserChat(uidAndGroupID[0], &pbSaveData)
			if err != nil {
				log.NewError(pbSaveData.OperationID, "group data insert to mongo err", pbSaveData.String(), uidAndGroupID[0], err.Error())
				return
			}
		}
		pbSaveData.Options = pbData.Options
		pbSaveData.OfflineInfo = pbData.OfflineInfo
		go sendMessageToPush(&pbSaveData)
	default:
		log.NewError(pbSaveData.OperationID, "SessionType error", pbSaveData.String())
		return
	}
	log.NewDebug(pbSaveData.OperationID, "msg_transfer handle topic data to database success...", pbSaveData.String())
}

func (HistoryConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (HistoryConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (mc *HistoryConsumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.InfoByKv("kafka get info to mongo", "", "msgTopic", msg.Topic, "msgPartition", msg.Partition, "chat", string(msg.Value))
		mc.msgHandle[msg.Topic](msg.Value, string(msg.Key))
		sess.MarkMessage(msg, "")
	}
	return nil
}
func sendMessageToPush(message *pbMsg.MsgSvrToPushSvrChatMsg) {
	log.InfoByKv("msg_transfer send message to push", message.OperationID, "message", message.String())
	msg := pbPush.PushMsgReq{}
	msg.OperationID = message.OperationID
	msg.PlatformID = message.PlatformID
	msg.Content = message.Content
	msg.ContentType = message.ContentType
	msg.SessionType = message.SessionType
	msg.RecvID = message.RecvID
	msg.SendID = message.SendID
	msg.SenderNickName = message.SenderNickName
	msg.SenderFaceURL = message.SenderFaceURL
	msg.ClientMsgID = message.ClientMsgID
	msg.MsgFrom = message.MsgFrom
	msg.Options = message.Options
	msg.RecvSeq = message.RecvSeq
	msg.SendTime = message.SendTime
	msg.MsgID = message.MsgID
	msg.OfflineInfo = message.OfflineInfo
	grpcConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImPushName)
	if grpcConn == nil {
		log.ErrorByKv("rpc dial failed", msg.OperationID, "push data", msg.String())
		pid, offset, err := producer.SendMessage(message)
		if err != nil {
			log.ErrorByKv("kafka send failed", msg.OperationID, "send data", message.String(), "pid", pid, "offset", offset, "err", err.Error())
		}
		return
	}
	msgClient := pbPush.NewPushMsgServiceClient(grpcConn)
	_, err := msgClient.PushMsg(context.Background(), &msg)
	if err != nil {
		log.ErrorByKv("rpc send failed", msg.OperationID, "push data", msg.String(), "err", err.Error())
		pid, offset, err := producer.SendMessage(message)
		if err != nil {
			log.ErrorByKv("kafka send failed", msg.OperationID, "send data", message.String(), "pid", pid, "offset", offset, "err", err.Error())
		}
	} else {
		log.InfoByKv("rpc send success", msg.OperationID, "push data", msg.String())

	}
}
