/*
** description("").
** copyright('tuoyun,www.tuoyun.net').
** author("fg,Gordon@tuoyun.net").
** time(2021/5/11 15:37).
 */
package logic

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db/mysql_model/im_mysql_msg_model"
	kfk "Open_IM/pkg/common/kafka"
	"Open_IM/pkg/common/log"
	pbMsg "Open_IM/pkg/proto/chat"
	"Open_IM/pkg/utils"
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
)

type PersistentConsumerHandler struct {
	msgHandle               map[string]fcb
	persistentConsumerGroup *kfk.MConsumerGroup
}

func (pc *PersistentConsumerHandler) Init() {
	pc.msgHandle = make(map[string]fcb)
	pc.msgHandle[config.Config.Kafka.Ws2mschat.Topic] = pc.handleChatWs2Mysql
	pc.persistentConsumerGroup = kfk.NewMConsumerGroup(&kfk.MConsumerGroupConfig{KafkaVersion: sarama.V0_10_2_0,
		OffsetsInitial: sarama.OffsetNewest, IsReturnErr: false}, []string{config.Config.Kafka.Ws2mschat.Topic},
		config.Config.Kafka.Ws2mschat.Addr, config.Config.Kafka.ConsumerGroupID.MsgToMySql)

}

func (pc *PersistentConsumerHandler) handleChatWs2Mysql(msg []byte, msgKey string) {
	log.NewInfo("msg come here mysql!!!", "", "msg", string(msg))
	msgFromMQ := pbMsg.MsgDataToMQ{}
	err := proto.Unmarshal(msg, &msgFromMQ)
	if err != nil {
		log.ErrorByKv("msg_transfer Unmarshal msg err", "", "msg", string(msg), "err", err.Error())
		return
	}
	//Control whether to store history messages (mysql)
	isPersist := utils.GetSwitchFromOptions(msgFromMQ.MsgData.Options, constant.IsPersistent)
	//Only process receiver data
	if isPersist {
		if msgKey == msgFromMQ.MsgData.RecvID && msgFromMQ.MsgData.SessionType == constant.SingleChatType {
			log.InfoByKv("msg_transfer msg persisting", msgFromMQ.OperationID)
			if err = im_mysql_msg_model.InsertMessageToChatLog(msgFromMQ); err != nil {
				log.ErrorByKv("Message insert failed", msgFromMQ.OperationID, "err", err.Error(), "msg", msgFromMQ.String())
				return
			}
		} else if msgFromMQ.MsgData.SessionType == constant.GroupChatType && msgKey == msgFromMQ.MsgData.SendID {
			log.InfoByKv("msg_transfer msg persisting", msgFromMQ.OperationID)
			if err = im_mysql_msg_model.InsertMessageToChatLog(msgFromMQ); err != nil {
				log.ErrorByKv("Message insert failed", msgFromMQ.OperationID, "err", err.Error(), "msg", msgFromMQ.String())
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
		log.InfoByKv("kafka get info to mysql", "", "msgTopic", msg.Topic, "msgPartition", msg.Partition, "msg", string(msg.Value))
		pc.msgHandle[msg.Topic](msg.Value, string(msg.Key))
		sess.MarkMessage(msg, "")
	}
	return nil
}
