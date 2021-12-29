/*
** description("").
** copyright('tuoyun,www.tuoyun.net').
** author("fg,Gordon@tuoyun.net").
** time(2021/5/11 15:37).
 */
package logic

import (
	"Open_IM/pkg/common/mq"
	"strings"

	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db/mysql_model/im_mysql_msg_model"
	"Open_IM/pkg/common/log"
	kfk "Open_IM/pkg/common/mq/kafka"
	pbMsg "Open_IM/pkg/proto/chat"
	"Open_IM/pkg/utils"

	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
)

type PersistentConsumerHandler struct {
	persistentConsumerGroup mq.Consumer
}

func (pc *PersistentConsumerHandler) Init() {
	pc.persistentConsumerGroup = kfk.NewMConsumerGroup(&kfk.MConsumerGroupConfig{KafkaVersion: sarama.V0_10_2_0,
		OffsetsInitial: sarama.OffsetNewest, IsReturnErr: false},
		config.Config.Kafka.Ws2mschat.Addr, config.Config.Kafka.ConsumerGroupID.MsgToMySql)

	pc.persistentConsumerGroup.RegisterMessageHandler(config.Config.Kafka.Ws2mschat.Topic, mq.MessageHandleFunc(pc.handleChatWs2Mysql))
}

func (pc *PersistentConsumerHandler) handleChatWs2Mysql(message *mq.Message) error {
	msg, msgKey := message.Value, string(message.Key)
	log.InfoByKv("chat come here mysql!!!", "", "chat", string(msg))
	pbData := pbMsg.WSToMsgSvrChatMsg{}
	err := proto.Unmarshal(msg, &pbData)
	if err != nil {
		log.ErrorByKv("msg_transfer Unmarshal chat err", "", "chat", string(msg), "err", err.Error())
		return nil // not retry
	}
	Options := utils.JsonStringToMap(pbData.Options)
	//Control whether to store history messages (mysql)
	isPersist := utils.GetSwitchFromOptions(Options, "persistent")
	//Only process receiver data
	if isPersist {
		if msgKey == pbData.RecvID && pbData.SessionType == constant.SingleChatType {
			log.InfoByKv("msg_transfer chat persisting", pbData.OperationID)
			if err = im_mysql_msg_model.InsertMessageToChatLog(pbData); err != nil {
				log.ErrorByKv("Message insert failed", pbData.OperationID, "err", err.Error(), "chat", pbData.String())
				return err
			}
		} else if pbData.SessionType == constant.GroupChatType && msgKey == "0" {
			pbData.RecvID = strings.Split(pbData.RecvID, " ")[1]
			log.InfoByKv("msg_transfer chat persisting", pbData.OperationID)
			if err = im_mysql_msg_model.InsertMessageToChatLog(pbData); err != nil {
				log.ErrorByKv("Message insert failed", pbData.OperationID, "err", err.Error(), "chat", pbData.String())
				return err
			}
		}

	}

	return nil
}
