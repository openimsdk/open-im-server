package msgtransfer

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/controller"
	kfk "github.com/OpenIMSDK/Open-IM-Server/pkg/common/kafka"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/mcontext"
	pbMsg "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/msg"
	"github.com/Shopify/sarama"
	"google.golang.org/protobuf/proto"
)

type OnlineHistoryMongoConsumerHandler struct {
	historyConsumerGroup *kfk.MConsumerGroup
	msgDatabase          controller.CommonMsgDatabase
	notificationDatabase controller.NotificationDatabase
}

func NewOnlineHistoryMongoConsumerHandler(database controller.CommonMsgDatabase, notificationDatabase controller.NotificationDatabase) *OnlineHistoryMongoConsumerHandler {
	mc := &OnlineHistoryMongoConsumerHandler{
		historyConsumerGroup: kfk.NewMConsumerGroup(&kfk.MConsumerGroupConfig{KafkaVersion: sarama.V2_0_0_0,
			OffsetsInitial: sarama.OffsetNewest, IsReturnErr: false}, []string{config.Config.Kafka.MsgToMongo.Topic},
			config.Config.Kafka.MsgToMongo.Addr, config.Config.Kafka.ConsumerGroupID.MsgToMongo),
		msgDatabase:          database,
		notificationDatabase: notificationDatabase,
	}
	return mc
}

func (mc *OnlineHistoryMongoConsumerHandler) handleChatWs2Mongo(ctx context.Context, cMsg *sarama.ConsumerMessage, conversationID string, session sarama.ConsumerGroupSession) {
	msg := cMsg.Value
	msgFromMQ := pbMsg.MsgDataToMongoByMQ{}
	operationID := mcontext.GetOperationID(ctx)
	err := proto.Unmarshal(msg, &msgFromMQ)
	if err != nil {
		log.ZError(ctx, "unmarshall failed", err, "conversationID", conversationID, "len", len(msg))
		return
	}
	if len(msgFromMQ.MsgData) == 0 {
		log.ZError(ctx, "msgFromMQ.MsgData is empty", nil, "cMsg", cMsg)
		return
	}
	log.ZInfo(ctx, "mongo consumer recv msg", "msgs", msgFromMQ.MsgData)
	isNotification := msgFromMQ.MsgData[0].Options[constant.IsNotification]
	if isNotification {
		err = mc.notificationDatabase.BatchInsertChat2DB(ctx, msgFromMQ.ConversationID, msgFromMQ.MsgData, msgFromMQ.LastSeq)
		if err != nil {
			log.NewError(operationID, "single data insert to mongo err", err.Error(), msgFromMQ.MsgData, msgFromMQ.ConversationID, msgFromMQ.TriggerID)
		}
		err = mc.notificationDatabase.DeleteMessageFromCache(ctx, msgFromMQ.ConversationID, msgFromMQ.MsgData)
		if err != nil {
			log.NewError(operationID, "remove cache msg from redis err", err.Error(), msgFromMQ.MsgData, msgFromMQ.ConversationID, msgFromMQ.TriggerID)
		}
		for _, v := range msgFromMQ.MsgData {
			if v.ContentType == constant.DeleteMessageNotification {
				deleteMessageTips := sdkws.DeleteMessageTips{}
				err := proto.Unmarshal(v.Content, &deleteMessageTips)
				if err != nil {
					log.NewError(operationID, "tips unmarshal err:", err.Error(), v.String())
					continue
				}
				if totalUnExistSeqs, err := mc.notificationDatabase.DelMsgBySeqs(ctx, deleteMessageTips.UserID, deleteMessageTips.Seqs); err != nil {
					log.NewError(operationID, utils.GetSelfFuncName(), "DelMsgBySeqs args: ", deleteMessageTips.UserID, deleteMessageTips.Seqs, "error:", err.Error(), "totalUnExistSeqs: ", totalUnExistSeqs)
				}
			}
		}
	} else {
		err = mc.msgDatabase.BatchInsertChat2DB(ctx, msgFromMQ.ConversationID, msgFromMQ.MsgData, msgFromMQ.LastSeq)
		if err != nil {
			log.NewError(operationID, "single data insert to mongo err", err.Error(), msgFromMQ.MsgData, msgFromMQ.ConversationID, msgFromMQ.TriggerID)
		}
		err = mc.msgDatabase.DeleteMessageFromCache(ctx, msgFromMQ.ConversationID, msgFromMQ.MsgData)
		if err != nil {
			log.NewError(operationID, "remove cache msg from redis err", err.Error(), msgFromMQ.MsgData, msgFromMQ.ConversationID, msgFromMQ.TriggerID)
		}
		for _, v := range msgFromMQ.MsgData {
			if v.ContentType == constant.DeleteMessageNotification {
				deleteMessageTips := sdkws.DeleteMessageTips{}
				err := proto.Unmarshal(v.Content, &deleteMessageTips)
				if err != nil {
					log.NewError(operationID, "tips unmarshal err:", err.Error(), v.String())
					continue
				}
				if totalUnExistSeqs, err := mc.msgDatabase.DelMsgBySeqs(ctx, deleteMessageTips.UserID, deleteMessageTips.Seqs); err != nil {
					log.NewError(operationID, utils.GetSelfFuncName(), "DelMsgBySeqs args: ", deleteMessageTips.UserID, deleteMessageTips.Seqs, "error:", err.Error(), "totalUnExistSeqs: ", totalUnExistSeqs)
				}
			}
		}
	}
}

func (OnlineHistoryMongoConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (OnlineHistoryMongoConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }
func (mc *OnlineHistoryMongoConsumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error { // a instance in the consumer group
	log.ZDebug(context.Background(), "online new session msg come", "highWaterMarkOffset",
		claim.HighWaterMarkOffset(), "topic", claim.Topic(), "partition", claim.Partition())
	for msg := range claim.Messages() {
		ctx := mc.historyConsumerGroup.GetContextFromMsg(msg)
		if len(msg.Value) != 0 {
			log.ZDebug(ctx, "mongo consumer recv new msg", "conversationID", msg.Key, "offset", msg.Offset)
			mc.handleChatWs2Mongo(ctx, msg, string(msg.Key), sess)
		} else {
			log.ZError(ctx, "mongo msg get from kafka but is nil", nil, "conversationID", msg.Key)
		}
		sess.MarkMessage(msg, "")
	}
	return nil
}
