/*
** description("").
** copyright('tuoyun,www.tuoyun.net').
** author("fg,Gordon@tuoyun.net").
** time(2021/5/11 15:37).
 */
package msgtransfer

import (
	"OpenIM/pkg/common/config"
	"OpenIM/pkg/common/constant"
	"OpenIM/pkg/common/db/controller"
	kfk "OpenIM/pkg/common/kafka"
	"OpenIM/pkg/common/log"
	pbMsg "OpenIM/pkg/proto/msg"
	"OpenIM/pkg/utils"

	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
)

type PersistentConsumerHandler struct {
	persistentConsumerGroup *kfk.MConsumerGroup
	chatLogInterface        controller.ChatLogInterface
}

func (pc *PersistentConsumerHandler) Init() {
	pc.persistentConsumerGroup = kfk.NewMConsumerGroup(&kfk.MConsumerGroupConfig{KafkaVersion: sarama.V2_0_0_0,
		OffsetsInitial: sarama.OffsetNewest, IsReturnErr: false}, []string{config.Config.Kafka.Ws2mschat.Topic},
		config.Config.Kafka.Ws2mschat.Addr, config.Config.Kafka.ConsumerGroupID.MsgToMySql)
}

func (pc *PersistentConsumerHandler) handleChatWs2Mysql(cMsg *sarama.ConsumerMessage, msgKey string, _ sarama.ConsumerGroupSession) {
	msg := cMsg.Value
	log.NewInfo("msg come here mysql!!!", "", "msg", string(msg), msgKey)
	var tag bool
	msgFromMQ := pbMsg.MsgDataToMQ{}
	err := proto.Unmarshal(msg, &msgFromMQ)
	if err != nil {
		log.NewError(msgFromMQ.OperationID, "msg_transfer Unmarshal msg err", "msg", string(msg), "err", err.Error())
		return
	}
	log.Debug(msgFromMQ.OperationID, "proto.Unmarshal MsgDataToMQ", msgFromMQ.String())
	//Control whether to store history messages (mysql)
	isPersist := utils.GetSwitchFromOptions(msgFromMQ.MsgData.Options, constant.IsPersistent)
	//Only process receiver data
	if isPersist {
		switch msgFromMQ.MsgData.SessionType {
		case constant.SingleChatType, constant.NotificationChatType:
			if msgKey == msgFromMQ.MsgData.RecvID {
				tag = true
			}
		case constant.GroupChatType:
			if msgKey == msgFromMQ.MsgData.SendID {
				tag = true
			}
		case constant.SuperGroupChatType:
			tag = true
		}
		if tag {
			log.NewInfo(msgFromMQ.OperationID, "msg_transfer msg persisting", string(msg))
			if err = pc.chatLogInterface.CreateChatLog(msgFromMQ); err != nil {
				log.NewError(msgFromMQ.OperationID, "Message insert failed", "err", err.Error(), "msg", msgFromMQ.String())
				return
			}
		}

	}
}
func (PersistentConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (PersistentConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (pc *PersistentConsumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.NewDebug("", "kafka get info to mysql", "msgTopic", msg.Topic, "msgPartition", msg.Partition, "msg", string(msg.Value), "key", string(msg.Key))
		if len(msg.Value) != 0 {
			pc.handleChatWs2Mysql(msg, string(msg.Key), sess)
		} else {
			log.Error("", "msg get from kafka but is nil", msg.Key)
		}
		sess.MarkMessage(msg, "")
	}
	return nil
}
