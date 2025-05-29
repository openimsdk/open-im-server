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
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/rpcli"
	"github.com/openimsdk/tools/discovery"

	"github.com/openimsdk/open-im-server/v3/pkg/common/prommetrics"

	"github.com/IBM/sarama"
	"github.com/go-redis/redis"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/controller"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/kafka"
	"github.com/openimsdk/open-im-server/v3/pkg/msgprocessor"
	"github.com/openimsdk/open-im-server/v3/pkg/tools/batcher"
	"github.com/openimsdk/protocol/constant"
	pbconv "github.com/openimsdk/protocol/conversation"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mcontext"
	"github.com/openimsdk/tools/utils/stringutil"
	"google.golang.org/protobuf/proto"
)

const (
	size              = 500
	mainDataBuffer    = 500
	subChanBuffer     = 50
	worker            = 50
	interval          = 100 * time.Millisecond
	hasReadChanBuffer = 1000
)

type ContextMsg struct {
	message *sdkws.MsgData
	ctx     context.Context
}

// This structure is used for asynchronously writing the sender’s read sequence (seq) regarding a message into MongoDB.
// For example, if the sender sends a message with a seq of 10, then their own read seq for this conversation should be set to 10.
type userHasReadSeq struct {
	conversationID string
	userHasReadMap map[string]int64
}

type OnlineHistoryRedisConsumerHandler struct {
	historyConsumerGroup *kafka.MConsumerGroup

	redisMessageBatches *batcher.Batcher[sarama.ConsumerMessage]

	msgTransferDatabase         controller.MsgTransferDatabase
	conversationUserHasReadChan chan *userHasReadSeq
	wg                          sync.WaitGroup

	groupClient        *rpcli.GroupClient
	conversationClient *rpcli.ConversationClient
}

func NewOnlineHistoryRedisConsumerHandler(ctx context.Context, client discovery.SvcDiscoveryRegistry, config *Config, database controller.MsgTransferDatabase) (*OnlineHistoryRedisConsumerHandler, error) {
	kafkaConf := config.KafkaConfig
	historyConsumerGroup, err := kafka.NewMConsumerGroup(kafkaConf.Build(), kafkaConf.ToRedisGroupID, []string{kafkaConf.ToRedisTopic}, false)
	if err != nil {
		return nil, err
	}
	groupConn, err := client.GetConn(ctx, config.Share.RpcRegisterName.Group)
	if err != nil {
		return nil, err
	}
	conversationConn, err := client.GetConn(ctx, config.Share.RpcRegisterName.Conversation)
	if err != nil {
		return nil, err
	}
	var och OnlineHistoryRedisConsumerHandler
	och.msgTransferDatabase = database
	och.conversationUserHasReadChan = make(chan *userHasReadSeq, hasReadChanBuffer)
	och.groupClient = rpcli.NewGroupClient(groupConn)
	och.conversationClient = rpcli.NewConversationClient(conversationConn)
	och.wg.Add(1)

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
	och.historyConsumerGroup = historyConsumerGroup

	return &och, nil
}
func (och *OnlineHistoryRedisConsumerHandler) do(ctx context.Context, channelID int, val *batcher.Msg[sarama.ConsumerMessage]) {
	ctx = mcontext.WithTriggerIDContext(ctx, val.TriggerID())
	ctxMessages := och.parseConsumerMessages(ctx, val.Val())
	ctx = withAggregationCtx(ctx, ctxMessages)
	log.ZInfo(ctx, "msg arrived channel", "channel id", channelID, "msgList length", len(ctxMessages), "key", val.Key())
	och.doSetReadSeq(ctx, ctxMessages)

	storageMsgList, notStorageMsgList, storageNotificationList, notStorageNotificationList :=
		och.categorizeMessageLists(ctxMessages)
	log.ZDebug(ctx, "number of categorized messages", "storageMsgList", len(storageMsgList), "notStorageMsgList",
		len(notStorageMsgList), "storageNotificationList", len(storageNotificationList), "notStorageNotificationList", len(notStorageNotificationList))

	conversationIDMsg := msgprocessor.GetChatConversationIDByMsg(ctxMessages[0].message)
	conversationIDNotification := msgprocessor.GetNotificationConversationIDByMsg(ctxMessages[0].message)
	och.handleMsg(ctx, val.Key(), conversationIDMsg, storageMsgList, notStorageMsgList)
	och.handleNotification(ctx, val.Key(), conversationIDNotification, storageNotificationList, notStorageNotificationList)
}

func (och *OnlineHistoryRedisConsumerHandler) doSetReadSeq(ctx context.Context, msgs []*ContextMsg) {

	var conversationID string
	var userSeqMap map[string]int64
	for _, msg := range msgs {
		if msg.message.ContentType != constant.HasReadReceipt {
			continue
		}
		var elem sdkws.NotificationElem
		if err := json.Unmarshal(msg.message.Content, &elem); err != nil {
			log.ZWarn(ctx, "handlerConversationRead Unmarshal NotificationElem msg err", err, "msg", msg)
			continue
		}
		var tips sdkws.MarkAsReadTips
		if err := json.Unmarshal([]byte(elem.Detail), &tips); err != nil {
			log.ZWarn(ctx, "handlerConversationRead Unmarshal MarkAsReadTips msg err", err, "msg", msg)
			continue
		}
		//The conversation ID for each batch of messages processed by the batcher is the same.
		conversationID = tips.ConversationID
		if len(tips.Seqs) > 0 {
			for _, seq := range tips.Seqs {
				if tips.HasReadSeq < seq {
					tips.HasReadSeq = seq
				}
			}
			clear(tips.Seqs)
			tips.Seqs = nil
		}
		if tips.HasReadSeq < 0 {
			continue
		}
		if userSeqMap == nil {
			userSeqMap = make(map[string]int64)
		}

		if userSeqMap[tips.MarkAsReadUserID] > tips.HasReadSeq {
			continue
		}
		userSeqMap[tips.MarkAsReadUserID] = tips.HasReadSeq
	}
	if userSeqMap == nil {
		return
	}
	if len(conversationID) == 0 {
		log.ZWarn(ctx, "conversation err", nil, "conversationID", conversationID)
	}
	if err := och.msgTransferDatabase.SetHasReadSeqToDB(ctx, conversationID, userSeqMap); err != nil {
		log.ZWarn(ctx, "set read seq to db error", err, "conversationID", conversationID, "userSeqMap", userSeqMap)
	}

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
	log.ZInfo(ctx, "handle storage msg")
	for _, storageMsg := range storageList {
		log.ZDebug(ctx, "handle storage msg", "msg", storageMsg.message.String())
	}

	och.toPushTopic(ctx, key, conversationID, notStorageList)
	var storageMessageList []*sdkws.MsgData
	for _, msg := range storageList {
		storageMessageList = append(storageMessageList, msg.message)
	}
	if len(storageMessageList) > 0 {
		msg := storageMessageList[0]
		lastSeq, isNewConversation, userSeqMap, err := och.msgTransferDatabase.BatchInsertChat2Cache(ctx, conversationID, storageMessageList)
		if err != nil && !errors.Is(errs.Unwrap(err), redis.Nil) {
			log.ZWarn(ctx, "batch data insert to redis err", err, "storageMsgList", storageMessageList)
			return
		}
		log.ZInfo(ctx, "BatchInsertChat2Cache end")
		err = och.msgTransferDatabase.SetHasReadSeqs(ctx, conversationID, userSeqMap)
		if err != nil {
			log.ZWarn(ctx, "SetHasReadSeqs error", err, "userSeqMap", userSeqMap, "conversationID", conversationID)
			prommetrics.SeqSetFailedCounter.Inc()
		}
		och.conversationUserHasReadChan <- &userHasReadSeq{
			conversationID: conversationID,
			userHasReadMap: userSeqMap,
		}

		if isNewConversation {
			ctx := storageList[0].ctx
			switch msg.SessionType {
			case constant.ReadGroupChatType:
				log.ZDebug(ctx, "group chat first create conversation", "conversationID",
					conversationID)

				userIDs, err := och.groupClient.GetGroupMemberUserIDs(ctx, msg.GroupID)
				if err != nil {
					log.ZWarn(ctx, "get group member ids error", err, "conversationID",
						conversationID)
				} else {
					log.ZInfo(ctx, "GetGroupMemberIDs end")

					if err := och.conversationClient.CreateGroupChatConversations(ctx, msg.GroupID, userIDs); err != nil {
						log.ZWarn(ctx, "single chat first create conversation error", err,
							"conversationID", conversationID)
					}
				}
			case constant.SingleChatType, constant.NotificationChatType:
				req := &pbconv.CreateSingleChatConversationsReq{
					RecvID:           msg.RecvID,
					SendID:           msg.SendID,
					ConversationID:   conversationID,
					ConversationType: msg.SessionType,
				}
				if err := och.conversationClient.CreateSingleChatConversations(ctx, req); err != nil {
					log.ZWarn(ctx, "single chat or notification first create conversation error", err,
						"conversationID", conversationID, "sessionType", msg.SessionType)
				}
			default:
				log.ZWarn(ctx, "unknown session type", nil, "sessionType",
					msg.SessionType)
			}
		}

		log.ZInfo(ctx, "success incr to next topic")
		err = och.msgTransferDatabase.MsgToMongoMQ(ctx, key, conversationID, storageMessageList, lastSeq)
		if err != nil {
			log.ZError(ctx, "Msg To MongoDB MQ error", err, "conversationID",
				conversationID, "storageList", storageMessageList, "lastSeq", lastSeq)
		}
		log.ZInfo(ctx, "MsgToMongoMQ end")

		och.toPushTopic(ctx, key, conversationID, storageList)
		log.ZInfo(ctx, "toPushTopic end")
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
		lastSeq, _, _, err := och.msgTransferDatabase.BatchInsertChat2Cache(ctx, conversationID, storageMessageList)
		if err != nil {
			log.ZError(ctx, "notification batch insert to redis error", err, "conversationID", conversationID,
				"storageList", storageMessageList)
			return
		}
		log.ZDebug(ctx, "success to next topic", "conversationID", conversationID)
		err = och.msgTransferDatabase.MsgToMongoMQ(ctx, key, conversationID, storageMessageList, lastSeq)
		if err != nil {
			log.ZError(ctx, "Msg To MongoDB MQ error", err, "conversationID",
				conversationID, "storageList", storageMessageList, "lastSeq", lastSeq)
		}
		och.toPushTopic(ctx, key, conversationID, storageList)
	}
}
func (och *OnlineHistoryRedisConsumerHandler) HandleUserHasReadSeqMessages(ctx context.Context) {
	defer func() {
		if r := recover(); r != nil {
			log.ZPanic(ctx, "HandleUserHasReadSeqMessages Panic", errs.ErrPanic(r))
		}
	}()

	defer och.wg.Done()

	for msg := range och.conversationUserHasReadChan {
		if err := och.msgTransferDatabase.SetHasReadSeqToDB(ctx, msg.conversationID, msg.userHasReadMap); err != nil {
			log.ZWarn(ctx, "set read seq to db error", err, "conversationID", msg.conversationID, "userSeqMap", msg.userHasReadMap)
		}
	}

	log.ZInfo(ctx, "Channel closed, exiting handleUserHasReadSeqMessages")
}
func (och *OnlineHistoryRedisConsumerHandler) Close() {
	close(och.conversationUserHasReadChan)
	och.wg.Wait()
}

func (och *OnlineHistoryRedisConsumerHandler) toPushTopic(ctx context.Context, key, conversationID string, msgs []*ContextMsg) {
	for _, v := range msgs {
		log.ZDebug(ctx, "push msg to topic", "msg", v.message.String())
		_, _, _ = och.msgTransferDatabase.MsgToPushMQ(v.ctx, key, conversationID, v.message)
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
	log.ZDebug(context.Background(), "online new session msg come", "highWaterMarkOffset",
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
