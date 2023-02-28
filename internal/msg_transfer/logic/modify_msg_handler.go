package logic

import (
	"Open_IM/pkg/base_info"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	kfk "Open_IM/pkg/common/kafka"
	"Open_IM/pkg/common/log"
	pbMsg "Open_IM/pkg/proto/msg"
	server_api_params "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"encoding/json"
	"github.com/Shopify/sarama"

	"github.com/golang/protobuf/proto"
)

type ModifyMsgConsumerHandler struct {
	msgHandle              map[string]fcb
	modifyMsgConsumerGroup *kfk.MConsumerGroup
}

func (mmc *ModifyMsgConsumerHandler) Init() {
	mmc.msgHandle = make(map[string]fcb)
	mmc.msgHandle[config.Config.Kafka.MsgToModify.Topic] = mmc.ModifyMsg
	mmc.modifyMsgConsumerGroup = kfk.NewMConsumerGroup(&kfk.MConsumerGroupConfig{KafkaVersion: sarama.V2_0_0_0,
		OffsetsInitial: sarama.OffsetNewest, IsReturnErr: false}, []string{config.Config.Kafka.MsgToModify.Topic},
		config.Config.Kafka.MsgToModify.Addr, config.Config.Kafka.ConsumerGroupID.MsgToModify)
}

func (ModifyMsgConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (ModifyMsgConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (mmc *ModifyMsgConsumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		log.NewDebug("", "kafka get info to mysql", "ModifyMsgConsumerHandler", msg.Topic, "msgPartition", msg.Partition, "msg", string(msg.Value), "key", string(msg.Key))
		if len(msg.Value) != 0 {
			mmc.msgHandle[msg.Topic](msg, string(msg.Key), sess)
		} else {
			log.Error("", "msg get from kafka but is nil", msg.Key)
		}
		sess.MarkMessage(msg, "")
	}
	return nil
}

func (mmc *ModifyMsgConsumerHandler) ModifyMsg(cMsg *sarama.ConsumerMessage, msgKey string, _ sarama.ConsumerGroupSession) {
	log.NewInfo("msg come here ModifyMsg!!!", "", "msg", string(cMsg.Value), msgKey)
	msgFromMQ := pbMsg.MsgDataToModifyByMQ{}
	err := proto.Unmarshal(cMsg.Value, &msgFromMQ)
	if err != nil {
		log.NewError(msgFromMQ.TriggerID, "msg_transfer Unmarshal msg err", "msg", string(cMsg.Value), "err", err.Error())
		return
	}
	log.Debug(msgFromMQ.TriggerID, "proto.Unmarshal MsgDataToMQ", msgFromMQ.String())
	for _, msgDataToMQ := range msgFromMQ.MessageList {
		isReactionFromCache := utils.GetSwitchFromOptions(msgDataToMQ.MsgData.Options, constant.IsReactionFromCache)
		if !isReactionFromCache {
			continue
		}
		if msgDataToMQ.MsgData.ContentType == constant.ReactionMessageModifier {
			notification := &base_info.ReactionMessageModifierNotification{}
			if err := json.Unmarshal(msgDataToMQ.MsgData.Content, notification); err != nil {
				continue
			}
			if notification.IsExternalExtensions {
				log.NewInfo(msgDataToMQ.OperationID, "msg:", notification, "this is external extensions")
				continue
			}

			if msgDataToMQ.MsgData.SessionType == constant.SuperGroupChatType && utils.GetSwitchFromOptions(msgDataToMQ.MsgData.Options, constant.IsHistory) {
				if msgDataToMQ.MsgData.Seq == 0 {
					log.NewError(msgDataToMQ.OperationID, "seq==0, error msg", msgDataToMQ.MsgData)
					continue
				}
				msg, err := db.DB.GetMsgBySeqIndex(notification.SourceID, notification.Seq, msgDataToMQ.OperationID)
				if (msg != nil && msg.Seq != notification.Seq) || err != nil {
					if err != nil {
						log.NewError(msgDataToMQ.OperationID, "GetMsgBySeqIndex failed", notification, err.Error())
					}
					msgs, indexes, err := db.DB.GetSuperGroupMsgBySeqListMongo(notification.SourceID, []uint32{notification.Seq}, msgDataToMQ.OperationID)
					if err != nil {
						log.NewError(msgDataToMQ.OperationID, "GetSuperGroupMsgBySeqListMongo failed", notification, err.Error())
						continue
					}
					var index int
					if len(msgs) < 1 || len(indexes) < 1 {
						log.NewError(msgDataToMQ.OperationID, "GetSuperGroupMsgBySeqListMongo failed", notification, "len<1", msgs, indexes)
						continue
					} else {
						msg = msgs[0]
						index = indexes[msg.Seq]
					}
					msg.IsReact = true
					if err := db.DB.ReplaceMsgByIndex(notification.SourceID, msg, index); err != nil {
						log.NewError(msgDataToMQ.OperationID, "ReplaceMsgByIndex failed", notification.SourceID, *msg)
					}

				} else {
					msg.IsReact = true
					if err = db.DB.ReplaceMsgBySeq(notification.SourceID, msg, msgDataToMQ.OperationID); err != nil {
						log.NewError(msgDataToMQ.OperationID, "ReplaceMsgBySeq failed", notification.SourceID, *msg)
					}
				}
			}

			if !notification.IsReact {
				// first time to modify
				var reactionExtensionList = make(map[string]db.KeyValue)
				extendMsg := db.ExtendMsg{
					ReactionExtensionList: reactionExtensionList,
					ClientMsgID:           notification.ClientMsgID,
					MsgFirstModifyTime:    notification.MsgFirstModifyTime,
				}
				for _, v := range notification.SuccessReactionExtensionList {
					reactionExtensionList[v.TypeKey] = db.KeyValue{
						TypeKey:          v.TypeKey,
						Value:            v.Value,
						LatestUpdateTime: v.LatestUpdateTime,
					}
				}
				// modify old msg
				if err := db.DB.InsertExtendMsg(notification.SourceID, notification.SessionType, &extendMsg); err != nil {
					log.NewError(msgDataToMQ.OperationID, "MsgFirstModify InsertExtendMsg failed", notification.SourceID, notification.SessionType, extendMsg, err.Error())
					continue
				}

			} else {
				var reactionExtensionList = make(map[string]*server_api_params.KeyValue)
				for _, v := range notification.SuccessReactionExtensionList {
					reactionExtensionList[v.TypeKey] = &server_api_params.KeyValue{
						TypeKey:          v.TypeKey,
						Value:            v.Value,
						LatestUpdateTime: v.LatestUpdateTime,
					}
				}
				// is already modify
				if err := db.DB.InsertOrUpdateReactionExtendMsgSet(notification.SourceID, notification.SessionType, notification.ClientMsgID, notification.MsgFirstModifyTime, reactionExtensionList); err != nil {
					log.NewError(msgDataToMQ.OperationID, "InsertOrUpdateReactionExtendMsgSet failed")
				}
			}
		} else if msgDataToMQ.MsgData.ContentType == constant.ReactionMessageDeleter {
			notification := &base_info.ReactionMessageDeleteNotification{}
			if err := json.Unmarshal(msgDataToMQ.MsgData.Content, notification); err != nil {
				continue
			}
			if err := db.DB.DeleteReactionExtendMsgSet(notification.SourceID, notification.SessionType, notification.ClientMsgID, notification.MsgFirstModifyTime, notification.SuccessReactionExtensionList); err != nil {
				log.NewError(msgDataToMQ.OperationID, "InsertOrUpdateReactionExtendMsgSet failed")
			}
		}
	}

}

func UnMarshallSetReactionMsgContent(content []byte) (notification *base_info.ReactionMessageModifierNotification, err error) {

	return notification, nil
}
