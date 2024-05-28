// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package msgtransfer

import (
	"context"
	"github.com/IBM/sarama"
	"github.com/go-redis/redis"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/controller"
	"github.com/openimsdk/open-im-server/v3/pkg/msgprocessor"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient"
	"github.com/openimsdk/open-im-server/v3/pkg/tools/batcher"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mcontext"
	"github.com/openimsdk/tools/mq/kafka"
	"github.com/openimsdk/tools/utils/stringutil"
	"google.golang.org/protobuf/proto"
	"strconv"
	"strings"
	"time"
)

const (
	size           = 500
	mainDataBuffer = 500
	subChanBuffer  = 50
	worker         = 50
	interval       = 100 * time.Millisecond
)

type ContextMsg struct {
	message *sdkws.MsgData
	ctx     context.Context
}

type OnlineHistoryRedisConsumerHandler struct {
	historyConsumerGroup *kafka.MConsumerGroup

	redisMessageBatches *batcher.Batcher[sarama.ConsumerMessage]

	msgDatabase           controller.CommonMsgDatabase
	conversationRpcClient *rpcclient.ConversationRpcClient
	groupRpcClient        *rpcclient.GroupRpcClient
}

func NewOnlineHistoryRedisConsumerHandler(kafkaConf *config.Kafka, database controller.CommonMsgDatabase,
	conversationRpcClient *rpcclient.ConversationRpcClient, groupRpcClient *rpcclient.GroupRpcClient) (*OnlineHistoryRedisConsumerHandler, error) {
	historyConsumerGroup, err := kafka.NewMConsumerGroup(kafkaConf.Build(), kafkaConf.ToRedisGroupID, []string{kafkaConf.ToRedisTopic}, false)
	if err != nil {
		return nil, err
	}
	var och OnlineHistoryRedisConsumerHandler
	och.msgDatabase = database

	b := batcher.New[sarama.ConsumerMessage](
		batcher.WithSize(size),
		batcher.WithWorker(worker),
		batcher.WithInterval(interval),
		batcher.WithDataBuffer(mainDataBuffer),
		batcher.WithSyncWait(true),
		batcher.WithBuffer(subChanBuffer),
	)
	b.Sharding = func(key string) int {
		hashCode := stringutil.GetHashCode(key)
		return int(hashCode) % och.redisMessageBatches.Worker()
	}
	b.Key = func(consumerMessage *sarama.ConsumerMessage) string {
		return string(consumerMessage.Key)
	}
	b.Do = och.do
	och.redisMessageBatches = b
	och.conversationRpcClient = conversationRpcClient
	och.groupRpcClient = groupRpcClient
	och.historyConsumerGroup = historyConsumerGroup
	return &och, err
}
func (och *OnlineHistoryRedisConsumerHandler) do(ctx context.Context, channelID int, val *batcher.Msg[sarama.ConsumerMessage]) {
	ctx = mcontext.WithTriggerIDContext(ctx, val.TriggerID())
	ctxMessages := och.parseConsumerMessages(ctx, val.Val())
	ctx = withAggregationCtx(ctx, ctxMessages)
	log.ZInfo(ctx, "msg arrived channel", "channel id", channelID, "msgList length", len(ctxMessages),
		"key", val.Key())

	storageMsgList, notStorageMsgList, storageNotificationList, notStorageNotificationList :=
		och.categorizeMessageLists(ctxMessages)
	log.ZDebug(ctx, "number of categorized messages", "storageMsgList", len(storageMsgList), "notStorageMsgList",
		len(notStorageMsgList), "storageNotificationList", len(storageNotificationList), "notStorageNotificationList",
		len(notStorageNotificationList))

	conversationIDMsg := msgprocessor.GetChatConversationIDByMsg(ctxMessages[0].message)
	conversationIDNotification := msgprocessor.GetNotificationConversationIDByMsg(ctxMessages[0].message)
	och.handleMsg(ctx, val.Key(), conversationIDMsg, storageMsgList, notStorageMsgList)
	och.handleNotification(ctx, val.Key(), conversationIDNotification, storageNotificationList, notStorageNotificationList)
}

func (och *OnlineHistoryRedisConsumerHandler) parseConsumerMessages(ctx context.Context, consumerMessages []*sarama.ConsumerMessage) []*ContextMsg {
	var ctxMessages []*ContextMsg
	for i := 0; i < len(consumerMessages); i++ {
		ctxMsg := &ContextMsg{}
		msgFromMQ := &sdkws.MsgData{}
		err := proto.Unmarshal(consumerMessages[i].Value, msgFromMQ)
		if err != nil {
			log.ZWarn(ctx, "msg_transfer Unmarshal msg err", err, string(consumerMessages[i].Value))
			continue
		}
		var arr []string
		for i, header := range consumerMessages[i].Headers {
			arr = append(arr, strconv.Itoa(i), string(header.Key), string(header.Value))
		}
		log.ZDebug(ctx, "consumer.kafka.GetContextWithMQHeader", "len", len(consumerMessages[i].Headers),
			"header", strings.Join(arr, ", "))
		ctxMsg.ctx = kafka.GetContextWithMQHeader(consumerMessages[i].Headers)
		ctxMsg.message = msgFromMQ
		log.ZDebug(ctx, "message parse finish", "message", msgFromMQ, "key",
			string(consumerMessages[i].Key))
		ctxMessages = append(ctxMessages, ctxMsg)
	}
	return ctxMessages
}

// Get messages/notifications stored message list, not stored and pushed message list.
func (och *OnlineHistoryRedisConsumerHandler) categorizeMessageLists(totalMsgs []*ContextMsg) (storageMsgList,
	notStorageMsgList, storageNotificationList, notStorageNotificationList []*ContextMsg) {
	for _, v := range totalMsgs {
		options := msgprocessor.Options(v.message.Options)
		if !options.IsNotNotification() {
			// clone msg from notificationMsg
			if options.IsSendMsg() {
				msg := proto.Clone(v.message).(*sdkws.MsgData)
				// message
				if v.message.Options != nil {
					msg.Options = msgprocessor.NewMsgOptions()
				}
				msg.Options = msgprocessor.WithOptions(msg.Options,
					msgprocessor.WithOfflinePush(options.IsOfflinePush()),
					msgprocessor.WithUnreadCount(options.IsUnreadCount()),
				)
				v.message.Options = msgprocessor.WithOptions(
					v.message.Options,
					msgprocessor.WithOfflinePush(false),
					msgprocessor.WithUnreadCount(false),
				)
				ctxMsg := &ContextMsg{
					message: msg,
					ctx:     v.ctx,
				}
				storageMsgList = append(storageMsgList, ctxMsg)
			}
			if options.IsHistory() {
				storageNotificationList = append(storageNotificationList, v)
			} else {
				notStorageNotificationList = append(notStorageNotificationList, v)
			}
		} else {
			if options.IsHistory() {
				storageMsgList = append(storageMsgList, v)
			} else {
				notStorageMsgList = append(notStorageMsgList, v)
			}
		}
	}
	return
}

func (och *OnlineHistoryRedisConsumerHandler) handleMsg(ctx context.Context, key, conversationID string, storageList, notStorageList []*ContextMsg) {
	och.toPushTopic(ctx, key, conversationID, notStorageList)
	var storageMessageList []*sdkws.MsgData
	for _, msg := range storageList {
		storageMessageList = append(storageMessageList, msg.message)
	}
	if len(storageMessageList) > 0 {
		msg := storageMessageList[0]
		lastSeq, isNewConversation, err := och.msgDatabase.BatchInsertChat2Cache(ctx, conversationID, storageMessageList)
		if err != nil && errs.Unwrap(err) != redis.Nil {
			log.ZError(ctx, "batch data insert to redis err", err, "storageMsgList", storageMessageList)
			return
		}
		if isNewConversation {
			switch msg.SessionType {
			case constant.ReadGroupChatType:
				log.ZInfo(ctx, "group chat first create conversation", "conversationID",
					conversationID)
				userIDs, err := och.groupRpcClient.GetGroupMemberIDs(ctx, msg.GroupID)
				if err != nil {
					log.ZWarn(ctx, "get group member ids error", err, "conversationID",
						conversationID)
				} else {
					if err := och.conversationRpcClient.GroupChatFirstCreateConversation(ctx,
						msg.GroupID, userIDs); err != nil {
						log.ZWarn(ctx, "single chat first create conversation error", err,
							"conversationID", conversationID)
					}
				}
			case constant.SingleChatType, constant.NotificationChatType:
				if err := och.conversationRpcClient.SingleChatFirstCreateConversation(ctx, msg.RecvID,
					msg.SendID, conversationID, msg.SessionType); err != nil {
					log.ZWarn(ctx, "single chat or notification first create conversation error", err,
						"conversationID", conversationID, "sessionType", msg.SessionType)
				}
			default:
				log.ZWarn(ctx, "unknown session type", nil, "sessionType",
					msg.SessionType)
			}
		}

		log.ZDebug(ctx, "success incr to next topic")
		err = och.msgDatabase.MsgToMongoMQ(ctx, key, conversationID, storageMessageList, lastSeq)
		if err != nil {
			log.ZError(ctx, "Msg To MongoDB MQ error", err, "conversationID",
				conversationID, "storageList", storageMessageList, "lastSeq", lastSeq)
		}
		och.toPushTopic(ctx, key, conversationID, storageList)
	}
}

func (och *OnlineHistoryRedisConsumerHandler) handleNotification(ctx context.Context, key, conversationID string,
	storageList, notStorageList []*ContextMsg) {
	och.toPushTopic(ctx, key, conversationID, notStorageList)
	var storageMessageList []*sdkws.MsgData
	for _, msg := range storageList {
		storageMessageList = append(storageMessageList, msg.message)
	}
	if len(storageMessageList) > 0 {
		lastSeq, _, err := och.msgDatabase.BatchInsertChat2Cache(ctx, conversationID, storageMessageList)
		if err != nil {
			log.ZError(ctx, "notification batch insert to redis error", err, "conversationID", conversationID,
				"storageList", storageMessageList)
			return
		}
		log.ZDebug(ctx, "success to next topic", "conversationID", conversationID)
		err = och.msgDatabase.MsgToMongoMQ(ctx, key, conversationID, storageMessageList, lastSeq)
		if err != nil {
			log.ZError(ctx, "Msg To MongoDB MQ error", err, "conversationID",
				conversationID, "storageList", storageMessageList, "lastSeq", lastSeq)
		}
		och.toPushTopic(ctx, key, conversationID, storageList)
	}
}

func (och *OnlineHistoryRedisConsumerHandler) toPushTopic(_ context.Context, key, conversationID string, msgs []*ContextMsg) {
	for _, v := range msgs {
		och.msgDatabase.MsgToPushMQ(v.ctx, key, conversationID, v.message)
	}
}

func withAggregationCtx(ctx context.Context, values []*ContextMsg) context.Context {
	var allMessageOperationID string
	for i, v := range values {
		if opid := mcontext.GetOperationID(v.ctx); opid != "" {
			if i == 0 {
				allMessageOperationID += opid
			} else {
				allMessageOperationID += "$" + opid
			}
		}
	}
	return mcontext.SetOperationID(ctx, allMessageOperationID)
}

func (och *OnlineHistoryRedisConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error { return nil }
func (och *OnlineHistoryRedisConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error {
	return nil
}

func (och *OnlineHistoryRedisConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession,
	claim sarama.ConsumerGroupClaim) error { // a instance in the consumer group
	log.ZInfo(context.Background(), "online new session msg come", "highWaterMarkOffset",
		claim.HighWaterMarkOffset(), "topic", claim.Topic(), "partition", claim.Partition())
	och.redisMessageBatches.OnComplete = func(lastMessage *sarama.ConsumerMessage, totalCount int) {
		session.MarkMessage(lastMessage, "")
		session.Commit()
	}
	for {
		select {
		case msg, ok := <-claim.Messages():
			if !ok {
				return nil
			}

			if len(msg.Value) == 0 {
				continue
			}
			err := och.redisMessageBatches.Put(context.Background(), msg)
			if err != nil {
				log.ZWarn(context.Background(), "put msg to  error", err, "msg", msg)
			}
		case <-session.Context().Done():
			return nil
		}
	}
}
