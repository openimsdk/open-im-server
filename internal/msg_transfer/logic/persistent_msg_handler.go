/*
** description("").
** copyright('tuoyun,www.tuoyun.net').
** author("fg,Gordon@tuoyun.net").
** time(2021/5/11 15:37).
 */
package logic

import (
	"Open_IM/pkg/common/mq"
	"Open_IM/pkg/common/mq/nsq"
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
	cfg := config.Config.MQ.Ws2mschat
	switch cfg.Type {
	case "kafka":
		pc.persistentConsumerGroup = kfk.NewMConsumerGroup(&kfk.MConsumerGroupConfig{KafkaVersion: sarama.V0_10_2_0,
			OffsetsInitial: sarama.OffsetNewest, IsReturnErr: false},
			cfg.Addr, config.Config.MQ.ConsumerGroupID.MsgToMySql)
	case "nsq":
		nc, err := nsq.NewNsqConsumer(cfg.Addr, cfg.Topic, cfg.Channel)
		if err != nil {
			panic(err)
		}
		pc.persistentConsumerGroup = nc
	default:
		panic("unsupported mq type: " + cfg.Type)
	}

	pc.persistentConsumerGroup.RegisterMessageHandler(cfg.Topic, mq.MessageHandleFunc(pc.handleChatWs2Mysql))
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
