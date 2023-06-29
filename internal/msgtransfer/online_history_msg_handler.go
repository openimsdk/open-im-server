package msgtransfer

import (
	"context"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/controller"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/kafka"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/mcontext"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/rpcclient"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"github.com/Shopify/sarama"
	"github.com/go-redis/redis"
	"google.golang.org/protobuf/proto"
)

const ConsumerMsgs = 3
const SourceMessages = 4
const MongoMessages = 5
const ChannelNum = 100

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
	Value interface{}
}
type ContextMsg struct {
	message *sdkws.MsgData
	ctx     context.Context
}

type OnlineHistoryRedisConsumerHandler struct {
	historyConsumerGroup *kafka.MConsumerGroup
	chArrays             [ChannelNum]chan Cmd2Value
	msgDistributionCh    chan Cmd2Value

	singleMsgSuccessCount      uint64
	singleMsgFailedCount       uint64
	singleMsgSuccessCountMutex sync.Mutex
	singleMsgFailedCountMutex  sync.Mutex

	msgDatabase           controller.CommonMsgDatabase
	conversationRpcClient *rpcclient.ConversationRpcClient
	groupRpcClient        *rpcclient.GroupRpcClient
}

func NewOnlineHistoryRedisConsumerHandler(database controller.CommonMsgDatabase, conversationRpcClient *rpcclient.ConversationRpcClient, groupRpcClient *rpcclient.GroupRpcClient) *OnlineHistoryRedisConsumerHandler {
	var och OnlineHistoryRedisConsumerHandler
	och.msgDatabase = database
	och.msgDistributionCh = make(chan Cmd2Value) //no buffer channel
	go och.MessagesDistributionHandle()
	for i := 0; i < ChannelNum; i++ {
		och.chArrays[i] = make(chan Cmd2Value, 50)
		go och.Run(i)
	}
	och.conversationRpcClient = conversationRpcClient
	och.groupRpcClient = groupRpcClient
	och.historyConsumerGroup = kafka.NewMConsumerGroup(&kafka.MConsumerGroupConfig{KafkaVersion: sarama.V2_0_0_0,
		OffsetsInitial: sarama.OffsetNewest, IsReturnErr: false}, []string{config.Config.Kafka.LatestMsgToRedis.Topic},
		config.Config.Kafka.Addr, config.Config.Kafka.ConsumerGroupID.MsgToRedis)
	//statistics.NewStatistics(&och.singleMsgSuccessCount, config.Config.ModuleName.MsgTransferName, fmt.Sprintf("%d second singleMsgCount insert to mongo", constant.StatisticsTimeInterval), constant.StatisticsTimeInterval)
	return &och
}

func (och *OnlineHistoryRedisConsumerHandler) Run(channelID int) {
	for {
		select {
		case cmd := <-och.chArrays[channelID]:
			switch cmd.Cmd {
			case SourceMessages:
				msgChannelValue := cmd.Value.(MsgChannelValue)
				ctxMsgList := msgChannelValue.ctxMsgList
				ctx := msgChannelValue.ctx
				log.ZDebug(ctx, "msg arrived channel", "channel id", channelID, "msgList length", len(ctxMsgList), "uniqueKey", msgChannelValue.uniqueKey)
				storageMsgList, notStorageMsgList, storageNotificationList, notStorageNotificationList, modifyMsgList := och.getPushStorageMsgList(ctxMsgList)
				log.ZDebug(ctx, "msg lens", "storageMsgList", len(storageMsgList), "notStorageMsgList", len(notStorageMsgList),
					"storageNotificationList", len(storageNotificationList), "notStorageNotificationList", len(notStorageNotificationList), "modifyMsgList", len(modifyMsgList))
				conversationIDMsg := utils.GetChatConversationIDByMsg(ctxMsgList[0].message)
				conversationIDNotification := utils.GetNotificationConversationID(ctxMsgList[0].message)
				och.handleMsg(ctx, msgChannelValue.uniqueKey, conversationIDMsg, storageMsgList, notStorageMsgList)
				och.handleNotification(ctx, msgChannelValue.uniqueKey, conversationIDNotification, storageNotificationList, notStorageNotificationList)
				if err := och.msgDatabase.MsgToModifyMQ(ctx, msgChannelValue.uniqueKey, conversationIDNotification, modifyMsgList); err != nil {
					log.ZError(ctx, "msg to modify mq error", err, "uniqueKey", msgChannelValue.uniqueKey, "modifyMsgList", modifyMsgList)
				}
			}
		}
	}
}

// 获取消息/通知 存储的消息列表， 不存储并且推送的消息列表，
func (och *OnlineHistoryRedisConsumerHandler) getPushStorageMsgList(totalMsgs []*ContextMsg) (storageMsgList, notStorageMsgList, storageNotificatoinList, notStorageNotificationList, modifyMsgList []*sdkws.MsgData) {
	isStorage := func(msg *sdkws.MsgData) bool {
		options2 := utils.Options(msg.Options)
		if options2.IsHistory() {
			return true
		} else {
			// if !(!options2.IsSenderSync() && conversationID == msg.MsgData.SendID) {
			// 	return false
			// }
			return false
		}
	}
	for _, v := range totalMsgs {
		options := utils.Options(v.message.Options)
		if !options.IsNotNotification() {
			// clone msg from notificationMsg
			if options.IsSendMsg() {
				msg := proto.Clone(v.message).(*sdkws.MsgData)
				// 消息
				if v.message.Options != nil {
					msg.Options = utils.NewMsgOptions()
				}
				if options.IsOfflinePush() {
					v.message.Options = utils.WithOptions(utils.Options(v.message.Options), utils.WithOfflinePush(false))
					msg.Options = utils.WithOptions(utils.Options(msg.Options), utils.WithOfflinePush(true))
				}
				if options.IsUnreadCount() {
					v.message.Options = utils.WithOptions(utils.Options(v.message.Options), utils.WithUnreadCount(false))
					msg.Options = utils.WithOptions(utils.Options(msg.Options), utils.WithUnreadCount(true))
				}
				storageMsgList = append(storageMsgList, msg)
			}
			if isStorage(v.message) {
				storageNotificatoinList = append(storageNotificatoinList, v.message)
			} else {
				notStorageNotificationList = append(notStorageNotificationList, v.message)
			}
		} else {
			if isStorage(v.message) {
				storageMsgList = append(storageMsgList, v.message)
			} else {
				notStorageMsgList = append(notStorageMsgList, v.message)
			}
		}
		if v.message.ContentType == constant.ReactionMessageModifier || v.message.ContentType == constant.ReactionMessageDeleter {
			modifyMsgList = append(modifyMsgList, v.message)
		}
	}
	return
}

func (och *OnlineHistoryRedisConsumerHandler) handleNotification(ctx context.Context, key, conversationID string, storageList, notStorageList []*sdkws.MsgData) {
	och.toPushTopic(ctx, key, conversationID, notStorageList)
	if len(storageList) > 0 {
		lastSeq, _, err := och.msgDatabase.BatchInsertChat2Cache(ctx, conversationID, storageList)
		if err != nil {
			log.ZError(ctx, "notification batch insert to redis error", err, "conversationID", conversationID, "storageList", storageList)
			return
		}
		log.ZDebug(ctx, "success to next topic", "conversationID", conversationID)
		och.msgDatabase.MsgToMongoMQ(ctx, key, conversationID, storageList, lastSeq)
		och.toPushTopic(ctx, key, conversationID, storageList)
	}
}

func (och *OnlineHistoryRedisConsumerHandler) toPushTopic(ctx context.Context, key, conversationID string, msgs []*sdkws.MsgData) {
	for _, v := range msgs {
		och.msgDatabase.MsgToPushMQ(ctx, key, conversationID, v)
	}
}

func (och *OnlineHistoryRedisConsumerHandler) handleMsg(ctx context.Context, key, conversationID string, storageList, notStorageList []*sdkws.MsgData) {
	och.toPushTopic(ctx, key, conversationID, notStorageList)
	if len(storageList) > 0 {
		lastSeq, isNewConversation, err := och.msgDatabase.BatchInsertChat2Cache(ctx, conversationID, storageList)
		if err != nil && errs.Unwrap(err) != redis.Nil {
			log.ZError(ctx, "batch data insert to redis err", err, "storageMsgList", storageList)
			och.singleMsgFailedCountMutex.Lock()
			och.singleMsgFailedCount += uint64(len(storageList))
			och.singleMsgFailedCountMutex.Unlock()
			return
		}
		if isNewConversation {
			if storageList[0].SessionType == constant.SuperGroupChatType {
				log.ZInfo(ctx, "group chat first create conversation", "conversationID", conversationID)
				userIDs, err := och.groupRpcClient.GetGroupMemberIDs(ctx, storageList[0].GroupID)
				if err != nil {
					log.ZWarn(ctx, "get group member ids error", err, "conversationID", conversationID)
				} else {
					if err := och.conversationRpcClient.GroupChatFirstCreateConversation(ctx, storageList[0].GroupID, userIDs); err != nil {
						log.ZWarn(ctx, "single chat first create conversation error", err, "conversationID", conversationID)
					}
				}
			} else {
				if err := och.conversationRpcClient.SingleChatFirstCreateConversation(ctx, storageList[0].RecvID, storageList[0].SendID); err != nil {
					log.ZWarn(ctx, "single chat first create conversation error", err, "conversationID", conversationID)
				}
			}
		}

		log.ZDebug(ctx, "success incr to next topic")
		och.singleMsgSuccessCountMutex.Lock()
		och.singleMsgSuccessCount += uint64(len(storageList))
		och.singleMsgSuccessCountMutex.Unlock()
		och.msgDatabase.MsgToMongoMQ(ctx, key, conversationID, storageList, lastSeq)
		och.toPushTopic(ctx, key, conversationID, storageList)
	}
}

func (och *OnlineHistoryRedisConsumerHandler) MessagesDistributionHandle() {
	for {
		aggregationMsgs := make(map[string][]*ContextMsg, ChannelNum)
		select {
		case cmd := <-och.msgDistributionCh:
			switch cmd.Cmd {
			case ConsumerMsgs:
				triggerChannelValue := cmd.Value.(TriggerChannelValue)
				ctx := triggerChannelValue.ctx
				consumerMessages := triggerChannelValue.cMsgList
				//Aggregation map[userid]message list
				log.ZDebug(ctx, "batch messages come to distribution center", "length", len(consumerMessages))
				for i := 0; i < len(consumerMessages); i++ {
					ctxMsg := &ContextMsg{}
					msgFromMQ := &sdkws.MsgData{}
					err := proto.Unmarshal(consumerMessages[i].Value, msgFromMQ)
					if err != nil {
						log.ZError(ctx, "msg_transfer Unmarshal msg err", err, string(consumerMessages[i].Value))
						continue
					}
					var arr []string
					for i, header := range consumerMessages[i].Headers {
						arr = append(arr, strconv.Itoa(i), string(header.Key), string(header.Value))
					}
					log.ZInfo(ctx, "consumer.kafka.GetContextWithMQHeader", "len", len(consumerMessages[i].Headers), "header", strings.Join(arr, ", "))
					ctxMsg.ctx = kafka.GetContextWithMQHeader(consumerMessages[i].Headers)
					ctxMsg.message = msgFromMQ
					log.ZDebug(ctx, "single msg come to distribution center", "message", msgFromMQ, "key", string(consumerMessages[i].Key))
					//aggregationMsgs[string(consumerMessages[i].Key)] = append(aggregationMsgs[string(consumerMessages[i].Key)], ctxMsg)
					if oldM, ok := aggregationMsgs[string(consumerMessages[i].Key)]; ok {
						oldM = append(oldM, ctxMsg)
						aggregationMsgs[string(consumerMessages[i].Key)] = oldM
					} else {
						m := make([]*ContextMsg, 0, 100)
						m = append(m, ctxMsg)
						aggregationMsgs[string(consumerMessages[i].Key)] = m
					}
				}
				log.ZDebug(ctx, "generate map list users len", "length", len(aggregationMsgs))
				for uniqueKey, v := range aggregationMsgs {
					if len(v) >= 0 {
						hashCode := utils.GetHashCode(uniqueKey)
						channelID := hashCode % ChannelNum
						newCtx := withAggregationCtx(ctx, v)
						log.ZDebug(newCtx, "generate channelID", "hashCode", hashCode, "channelID", channelID, "uniqueKey", uniqueKey)
						och.chArrays[channelID] <- Cmd2Value{Cmd: SourceMessages, Value: MsgChannelValue{uniqueKey: uniqueKey, ctxMsgList: v, ctx: newCtx}}
					}
				}
			}
		}
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

func (och *OnlineHistoryRedisConsumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error { // a instance in the consumer group
	for {
		if sess == nil {
			log.ZWarn(context.Background(), "sess == nil, waiting", nil)
			time.Sleep(100 * time.Millisecond)
		} else {
			break
		}
	}
	rwLock := new(sync.RWMutex)
	log.ZDebug(context.Background(), "online new session msg come", "highWaterMarkOffset",
		claim.HighWaterMarkOffset(), "topic", claim.Topic(), "partition", claim.Partition())
	cMsg := make([]*sarama.ConsumerMessage, 0, 1000)
	t := time.NewTicker(time.Millisecond * 100)
	go func() {
		for {
			select {
			case <-t.C:
				if len(cMsg) > 0 {
					rwLock.Lock()
					ccMsg := make([]*sarama.ConsumerMessage, 0, 1000)
					for _, v := range cMsg {
						ccMsg = append(ccMsg, v)
					}
					cMsg = make([]*sarama.ConsumerMessage, 0, 1000)
					rwLock.Unlock()
					split := 1000
					ctx := mcontext.WithTriggerIDContext(context.Background(), utils.OperationIDGenerator())
					log.ZDebug(ctx, "timer trigger msg consumer start", "length", len(ccMsg))
					for i := 0; i < len(ccMsg)/split; i++ {
						//log.Debug()
						och.msgDistributionCh <- Cmd2Value{Cmd: ConsumerMsgs, Value: TriggerChannelValue{
							ctx: ctx, cMsgList: ccMsg[i*split : (i+1)*split]}}
					}
					if (len(ccMsg) % split) > 0 {
						och.msgDistributionCh <- Cmd2Value{Cmd: ConsumerMsgs, Value: TriggerChannelValue{
							ctx: ctx, cMsgList: ccMsg[split*(len(ccMsg)/split):]}}
					}
					log.ZDebug(ctx, "timer trigger msg consumer end", "length", len(ccMsg))
				}
			}
		}
	}()
	for msg := range claim.Messages() {
		rwLock.Lock()
		if len(msg.Value) != 0 {
			cMsg = append(cMsg, msg)
		}
		rwLock.Unlock()
		sess.MarkMessage(msg, "")
	}
	return nil
}
