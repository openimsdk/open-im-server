package logic

import (
	"Open_IM/src/common/config"
	"Open_IM/src/common/constant"
	kfk "Open_IM/src/common/kafka"
	"Open_IM/src/common/log"
	pbMsg "Open_IM/src/proto/chat"
	pbPush "Open_IM/src/proto/push"
	"Open_IM/src/push/content_struct"
	"Open_IM/src/utils"
	"context"
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	"github.com/skiffer-git/grpc-etcdv3/getcdv3"
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
	if pbData.SessionType == constant.SingleChatType {
		log.Info("", "", "msg_transfer chat type = SingleChatType", isHistory, isPersist)
		if isHistory {
			if msgKey == pbSaveData.RecvID {
				err := saveUserChat(pbData.RecvID, &pbSaveData)
				if err != nil {
					log.ErrorByKv("data insert to mongo err", pbSaveData.OperationID, "data", pbSaveData.String(), "err", err.Error())
				}
				pbSaveData.Options = pbData.Options
				pbSaveData.OfflineInfo = pbData.OfflineInfo
				sendMessageToPush(&pbSaveData)
			} else if msgKey == pbSaveData.SendID {
				err := saveUserChat(pbData.SendID, &pbSaveData)
				if err != nil {
					log.ErrorByKv("data insert to mongo err", pbSaveData.OperationID, "data", pbSaveData.String(), "err", err.Error())
				}
				//if isSenderSync {
				//	pbSaveData.ContentType = constant.SyncSenderMsg
				//	log.WarnByKv("SyncSenderMsg come here", pbData.OperationID, pbSaveData.String())
				//	sendMessageToPush(&pbSaveData)
				//}
			}

		}
		log.InfoByKv("msg_transfer handle topic success...", "", "")
	} else if pbData.SessionType == constant.GroupChatType {
		log.Info("", "", "msg_transfer chat type = GroupChatType")
		uidAndGroupID := strings.Split(pbData.RecvID, " ")
		if pbData.ContentType == constant.AtText {
			atContent := content_struct.AtTextContent{
				Text:       pbData.Content,
				AtUserList: pbData.ForceList,
			}
			if utils.IsContain(uidAndGroupID[0], pbData.ForceList) {
				atContent.IsAtSelf = true
			}
			pbSaveData.Content = utils.StructToJsonString(atContent)
		}

		saveUserChat(uidAndGroupID[0], &pbSaveData)
		pbSaveData.Options = pbData.Options
		pbSaveData.OfflineInfo = pbData.OfflineInfo
		sendMessageToPush(&pbSaveData)
		log.InfoByKv("msg_transfer handle topic success...", "", "")
	} else {
		log.Error("", "", "msg_transfer recv chat err, chat.MsgFrom = %d", pbData.SessionType)
	}

}

func (HistoryConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (HistoryConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (mc *HistoryConsumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.InfoByKv("kafka get info to mongo", "", "msgTopic", msg.Topic, "msgPartition", msg.Partition, "chat", string(msg.Value))
		mc.msgHandle[msg.Topic](msg.Value, string(msg.Key))
	}
	return nil
}
func sendMessageToPush(message *pbMsg.MsgSvrToPushSvrChatMsg) {
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
	defer grpcConn.Close()
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
