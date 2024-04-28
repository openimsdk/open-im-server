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
	"github.com/openimsdk/open-im-server/v3/pkg/util/batcher"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/IBM/sarama"
	"github.com/go-redis/redis"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/controller"
	"github.com/openimsdk/open-im-server/v3/pkg/msgprocessor"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mcontext"
	"github.com/openimsdk/tools/mq/kafka"
	"github.com/openimsdk/tools/utils/idutil"
	"github.com/openimsdk/tools/utils/stringutil"
	"google.golang.org/protobuf/proto"
)

const (
	ConsumerMsgs   = 3
	SourceMessages = 4
	MongoMessages  = 5
	ChannelNum     = 100
)

type MsgChannelValue struct {
	uniqueKey  string
	ctx        context.Context
	ctxMsgList []*ContextMsg
}

type TriggerChannelValue struct {
	ctx      context.Context
	cMsgList []*sarama.ConsumerMessage
}

type Cmd2Value struct {
	Cmd   int
	Value any
}
type ContextMsg struct {
	message *sdkws.MsgData
	ctx     context.Context
}

type OnlineHistoryRedisConsumerHandler struct {
	historyConsumerGroup *kafka.MConsumerGroup
	chArrays             [ChannelNum]chan Cmd2Value
	msgDistributionCh    chan Cmd2Value

	redisMessageBatches *batcher.Batcher[sarama.ConsumerMessage]

	msgDatabase           controller.CommonMsgDatabase
	conversationRpcClient *rpcclient.ConversationRpcClient
	groupRpcClient        *rpcclient.GroupRpcClient
}

func NewOnlineHistoryRedisConsumerHandler(kafkaConf *config.Kafka, database controller.CommonMsgDatabase,
	conversationRpcClient *rpcclient.ConversationRpcClient, groupRpcClient *rpcclient.GroupRpcClient) (*OnlineHistoryRedisConsumerHandler, error) {
	historyConsumerGroup, err := kafka.NewMConsumerGroup(kafkaConf.Build(), kafkaConf.ToRedisGroupID, []string{kafkaConf.ToRedisTopic})
	if err != nil {
		return nil, err
	}
	var och OnlineHistoryRedisConsumerHandler
	och.msgDatabase = database

	b := batcher.New[sarama.ConsumerMessage]()
	b.Sharding = func(key string) int {
		hashCode := stringutil.GetHashCode(key)
		return int(hashCode) % och.redisMessageBatches.Worker()
	}
	b.Key = func(consumerMessage *sarama.ConsumerMessage) string {
		return string(consumerMessage.Key)
	}
	och.redisMessageBatches = b

	err = b.Start()
	if err != nil {
		return nil, err
	}
	//och.msgDistributionCh = make(chan Cmd2Value) // no buffer channel
	//go och.MessagesDistributionHandle()
	//for i := 0; i < ChannelNum; i++ {
	//	och.chArrays[i] = make(chan Cmd2Value, 50)
	//	go och.Run(i)
	//}
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
				if options.IsOfflinePush() {
					v.message.Options = msgprocessor.WithOptions(
						v.message.Options,
						msgprocessor.WithOfflinePush(false),
					)
					msg.Options = msgprocessor.WithOptions(msg.Options, msgprocessor.WithOfflinePush(true))
				}
				if options.IsUnreadCount() {
					v.message.Options = msgprocessor.WithOptions(
						v.message.Options,
						msgprocessor.WithUnreadCount(false),
					)
					msg.Options = msgprocessor.WithOptions(msg.Options, msgprocessor.WithUnreadCount(true))
				}
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
			log.ZError(ctx, "MsgToMongoMQ error", err)
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
			log.ZError(ctx, "MsgToMongoMQ error", err)
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

			session.MarkMessage(msg, "")

		case <-session.Context().Done():
			return nil
		}
	}

	var (
		split    = 1000
		rwLock   = new(sync.RWMutex)
		messages = make([]*sarama.ConsumerMessage, 0, 1000)
		ticker   = time.NewTicker(time.Millisecond * 100)

		wg      = sync.WaitGroup{}
		running = new(atomic.Bool)
	)
	running.Store(true)

	wg.Add(1)
	go func() {
		defer wg.Done()

		for {
			select {
			case <-ticker.C:
				// if the buffer is empty and running is false, return loop.
				if len(messages) == 0 {
					if !running.Load() {
						return
					}

					continue
				}

				rwLock.Lock()
				buffer := make([]*sarama.ConsumerMessage, 0, len(messages))
				buffer = append(buffer, messages...)

				// reuse slice, set cap to 0
				messages = messages[:0]
				rwLock.Unlock()

				start := time.Now()
				ctx := mcontext.WithTriggerIDContext(context.Background(), idutil.OperationIDGenerator())
				log.ZDebug(ctx, "timer trigger msg consumer start", "length", len(buffer))
				for i := 0; i < len(buffer)/split; i++ {
					och.msgDistributionCh <- Cmd2Value{Cmd: ConsumerMsgs, Value: TriggerChannelValue{
						ctx: ctx, cMsgList: buffer[i*split : (i+1)*split],
					}}
				}
				if (len(buffer) % split) > 0 {
					och.msgDistributionCh <- Cmd2Value{Cmd: ConsumerMsgs, Value: TriggerChannelValue{
						ctx: ctx, cMsgList: buffer[split*(len(buffer)/split):],
					}}
				}

				log.ZDebug(ctx, "timer trigger msg consumer end",
					"length", len(buffer), "time_cost", time.Since(start),
				)
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()

		for running.Load() {
			select {
			case msg, ok := <-claim.Messages():
				if !ok {
					running.Store(false)
					return
				}

				if len(msg.Value) == 0 {
					continue
				}

				rwLock.Lock()
				messages = append(messages, msg)
				rwLock.Unlock()

				session.MarkMessage(msg, "")

			case <-session.Context().Done():
				running.Store(false)
				return
			}
		}
	}()

	wg.Wait()
	return nil
}
