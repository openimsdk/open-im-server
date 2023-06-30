package msgtransfer

import (
	"context"
	"encoding/json"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/controller"
	unRelationTb "github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/unrelation"
	kfk "github.com/OpenIMSDK/Open-IM-Server/pkg/common/kafka"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/mcontext"
	pbMsg "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/msg"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"github.com/Shopify/sarama"

	"google.golang.org/protobuf/proto"
)

type ModifyMsgConsumerHandler struct {
	modifyMsgConsumerGroup *kfk.MConsumerGroup

	extendMsgDatabase controller.ExtendMsgDatabase
	extendSetMsgModel unRelationTb.ExtendMsgSetModel
}

func NewModifyMsgConsumerHandler(database controller.ExtendMsgDatabase) *ModifyMsgConsumerHandler {
	return &ModifyMsgConsumerHandler{
		modifyMsgConsumerGroup: kfk.NewMConsumerGroup(&kfk.MConsumerGroupConfig{KafkaVersion: sarama.V2_0_0_0,
			OffsetsInitial: sarama.OffsetNewest, IsReturnErr: false}, []string{config.Config.Kafka.MsgToModify.Topic},
			config.Config.Kafka.Addr, config.Config.Kafka.ConsumerGroupID.MsgToModify),
		extendMsgDatabase: database,
	}
}

func (ModifyMsgConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (ModifyMsgConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (mmc *ModifyMsgConsumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		ctx := mmc.modifyMsgConsumerGroup.GetContextFromMsg(msg)
		log.ZDebug(ctx, "kafka get info to mysql", "ModifyMsgConsumerHandler", msg.Topic, "msgPartition", msg.Partition, "msg", string(msg.Value), "key", string(msg.Key))
		if len(msg.Value) != 0 {
			mmc.ModifyMsg(ctx, msg, string(msg.Key), sess)
		} else {
			log.ZError(ctx, "msg get from kafka but is nil", nil, "key", msg.Key)
		}
		sess.MarkMessage(msg, "")
	}
	return nil
}

func (mmc *ModifyMsgConsumerHandler) ModifyMsg(ctx context.Context, cMsg *sarama.ConsumerMessage, msgKey string, _ sarama.ConsumerGroupSession) {
	msgFromMQ := pbMsg.MsgDataToModifyByMQ{}
	operationID := mcontext.GetOperationID(ctx)
	err := proto.Unmarshal(cMsg.Value, &msgFromMQ)
	if err != nil {
		log.ZError(ctx, "msg_transfer Unmarshal msg err", err, "msg", string(cMsg.Value))
		return
	}
	log.ZDebug(ctx, "proto.Unmarshal MsgDataToMQ", "msgs", msgFromMQ.String())
	for _, msg := range msgFromMQ.Messages {
		isReactionFromCache := utils.GetSwitchFromOptions(msg.Options, constant.IsReactionFromCache)
		if !isReactionFromCache {
			continue
		}
		ctx = mcontext.SetOperationID(ctx, operationID)
		if msg.ContentType == constant.ReactionMessageModifier {
			notification := &sdkws.ReactionMessageModifierNotification{}
			if err := json.Unmarshal(msg.Content, notification); err != nil {
				continue
			}
			if notification.IsExternalExtensions {
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
				for _, v := range notification.SuccessReactionExtensions {
					reactionExtensionList[v.TypeKey] = unRelationTb.KeyValueModel{
						TypeKey:          v.TypeKey,
						Value:            v.Value,
						LatestUpdateTime: v.LatestUpdateTime,
					}
				}

				if err := mmc.extendMsgDatabase.InsertExtendMsg(ctx, notification.ConversationID, notification.SessionType, &extendMsg); err != nil {
					// log.ZError(ctx, "MsgFirstModify InsertExtendMsg failed",  notification.ConversationID, notification.SessionType, extendMsg, err.Error())
					continue
				}
			} else {
				if err := mmc.extendMsgDatabase.InsertOrUpdateReactionExtendMsgSet(ctx, notification.ConversationID, notification.SessionType, notification.ClientMsgID, notification.MsgFirstModifyTime, mmc.extendSetMsgModel.Pb2Model(notification.SuccessReactionExtensions)); err != nil {
					// log.NewError(operationID, "InsertOrUpdateReactionExtendMsgSet failed")
				}
			}
		} else if msg.ContentType == constant.ReactionMessageDeleter {
			notification := &sdkws.ReactionMessageDeleteNotification{}
			if err := json.Unmarshal(msg.Content, notification); err != nil {
				continue
			}
			if err := mmc.extendMsgDatabase.DeleteReactionExtendMsgSet(ctx, notification.ConversationID, notification.SessionType, notification.ClientMsgID, notification.MsgFirstModifyTime, mmc.extendSetMsgModel.Pb2Model(notification.SuccessReactionExtensions)); err != nil {
				// log.NewError(operationID, "InsertOrUpdateReactionExtendMsgSet failed")
			}
		}
	}
}
