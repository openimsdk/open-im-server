package msgtransfer

import (
	"OpenIM/pkg/apistruct"
	"OpenIM/pkg/common/config"
	"OpenIM/pkg/common/constant"
	"OpenIM/pkg/common/db/controller"
	unRelationTb "OpenIM/pkg/common/db/table/unrelation"
	kfk "OpenIM/pkg/common/kafka"
	"OpenIM/pkg/common/log"
	"OpenIM/pkg/common/tracelog"
	pbMsg "OpenIM/pkg/proto/msg"
	sdkws "OpenIM/pkg/proto/sdkws"
	"OpenIM/pkg/utils"
	"context"
	"encoding/json"
	"github.com/Shopify/sarama"

	"github.com/golang/protobuf/proto"
)

type ModifyMsgConsumerHandler struct {
	modifyMsgConsumerGroup *kfk.MConsumerGroup

	extendMsgInterface controller.ExtendMsgInterface
}

func (mmc *ModifyMsgConsumerHandler) Init() {
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
			mmc.ModifyMsg(msg, string(msg.Key), sess)
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
		ctx := context.Background()
		tracelog.SetOperationID(ctx, msgDataToMQ.OperationID)
		if msgDataToMQ.MsgData.ContentType == constant.ReactionMessageModifier {
			notification := &apistruct.ReactionMessageModifierNotification{}
			if err := json.Unmarshal(msgDataToMQ.MsgData.Content, notification); err != nil {
				continue
			}
			if notification.IsExternalExtensions {
				log.NewInfo(msgDataToMQ.OperationID, "msg:", notification, "this is external extensions")
				continue
			}
			if !notification.IsReact {
				// first time to modify
				var reactionExtensionList = make(map[string]unRelationTb.KeyValueModel)
				extendMsg := unRelationTb.ExtendMsgModel{
					ReactionExtensionList: reactionExtensionList,
					ClientMsgID:           notification.ClientMsgID,
					MsgFirstModifyTime:    notification.MsgFirstModifyTime,
				}
				for _, v := range notification.SuccessReactionExtensionList {
					reactionExtensionList[v.TypeKey] = unRelationTb.KeyValueModel{
						TypeKey:          v.TypeKey,
						Value:            v.Value,
						LatestUpdateTime: v.LatestUpdateTime,
					}
				}

				if err := mmc.extendMsgInterface.InsertExtendMsg(ctx, notification.SourceID, notification.SessionType, &extendMsg); err != nil {
					log.NewError(msgDataToMQ.OperationID, "MsgFirstModify InsertExtendMsg failed", notification.SourceID, notification.SessionType, extendMsg, err.Error())
					continue
				}
			} else {
				var reactionExtensionList = make(map[string]*sdkws.KeyValue)
				for _, v := range notification.SuccessReactionExtensionList {
					reactionExtensionList[v.TypeKey] = &sdkws.KeyValue{
						TypeKey:          v.TypeKey,
						Value:            v.Value,
						LatestUpdateTime: v.LatestUpdateTime,
					}
				}
				// is already modify
				if err := mmc.extendMsgInterface.InsertOrUpdateReactionExtendMsgSet(ctx, notification.SourceID, notification.SessionType, notification.ClientMsgID, notification.MsgFirstModifyTime, reactionExtensionList); err != nil {
					log.NewError(msgDataToMQ.OperationID, "InsertOrUpdateReactionExtendMsgSet failed")
				}
			}
		} else if msgDataToMQ.MsgData.ContentType == constant.ReactionMessageDeleter {
			notification := &apistruct.ReactionMessageDeleteNotification{}
			if err := json.Unmarshal(msgDataToMQ.MsgData.Content, notification); err != nil {
				continue
			}
			if err := mmc.extendMsgInterface.DeleteReactionExtendMsgSet(ctx, notification.SourceID, notification.SessionType, notification.ClientMsgID, notification.MsgFirstModifyTime, notification.SuccessReactionExtensionList); err != nil {
				log.NewError(msgDataToMQ.OperationID, "InsertOrUpdateReactionExtendMsgSet failed")
			}
		}
	}

}
