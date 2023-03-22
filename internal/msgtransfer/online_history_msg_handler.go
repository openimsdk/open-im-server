package msgtransfer

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/controller"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/kafka"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/mcontext"
	pbMsg "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/msg"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	"sync"
	"time"
)

const ConsumerMsgs = 3
const AggregationMessages = 4
const MongoMessages = 5
const ChannelNum = 100

type MsgChannelValue struct {
	aggregationID string //maybe userID or super groupID
	ctx           context.Context
	ctxMsgList    []*ContextMsg
	lastSeq       uint64
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
	message *pbMsg.MsgDataToMQ
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

	//producerToPush   *kafka.Producer
	//producerToModify *kafka.Producer
	//producerToMongo  *kafka.Producer

	msgDatabase controller.MsgDatabase
}

func NewOnlineHistoryRedisConsumerHandler(database controller.MsgDatabase) *OnlineHistoryRedisConsumerHandler {
	var och OnlineHistoryRedisConsumerHandler
	och.msgDatabase = database
	och.msgDistributionCh = make(chan Cmd2Value) //no buffer channel
	go och.MessagesDistributionHandle()
	for i := 0; i < ChannelNum; i++ {
		och.chArrays[i] = make(chan Cmd2Value, 50)
		go och.Run(i)
	}
	//och.producerToPush = kafka.NewKafkaProducer(config.Config.Kafka.Ms2pschat.Addr, config.Config.Kafka.Ms2pschat.Topic)
	//och.producerToModify = kafka.NewKafkaProducer(config.Config.Kafka.MsgToModify.Addr, config.Config.Kafka.MsgToModify.Topic)
	//och.producerToMongo = kafka.NewKafkaProducer(config.Config.Kafka.MsgToMongo.Addr, config.Config.Kafka.MsgToMongo.Topic)
	och.historyConsumerGroup = kafka.NewMConsumerGroup(&kafka.MConsumerGroupConfig{KafkaVersion: sarama.V2_0_0_0,
		OffsetsInitial: sarama.OffsetNewest, IsReturnErr: false}, []string{config.Config.Kafka.Ws2mschat.Topic},
		config.Config.Kafka.Ws2mschat.Addr, config.Config.Kafka.ConsumerGroupID.MsgToRedis)
	//statistics.NewStatistics(&och.singleMsgSuccessCount, config.Config.ModuleName.MsgTransferName, fmt.Sprintf("%d second singleMsgCount insert to mongo", constant.StatisticsTimeInterval), constant.StatisticsTimeInterval)
	return &och
}

func (och *OnlineHistoryRedisConsumerHandler) Run(channelID int) {
	for {
		select {
		case cmd := <-och.chArrays[channelID]:
			switch cmd.Cmd {
			case AggregationMessages:
				msgChannelValue := cmd.Value.(MsgChannelValue)
				ctxMsgList := msgChannelValue.ctxMsgList
				ctx := msgChannelValue.ctx
				storageMsgList := make([]*pbMsg.MsgDataToMQ, 0, 80)
				storagePushMsgList := make([]*ContextMsg, 0, 80)
				notStoragePushMsgList := make([]*ContextMsg, 0, 80)
				log.ZDebug(ctx, "msg arrived channel", "channel id", channelID, "msgList length", len(ctxMsgList), "aggregationID", msgChannelValue.aggregationID)
				var modifyMsgList []*pbMsg.MsgDataToMQ
				//ctx := mcontext.NewCtx("redis consumer")
				//mcontext.SetOperationID(ctx, triggerID)
				for _, v := range ctxMsgList {
					log.ZDebug(ctx, "msg come to storage center", v.message.String())
					isHistory := utils.GetSwitchFromOptions(v.message.MsgData.Options, constant.IsHistory)
					isSenderSync := utils.GetSwitchFromOptions(v.message.MsgData.Options, constant.IsSenderSync)
					if isHistory {
						storageMsgList = append(storageMsgList, v.message)
						storagePushMsgList = append(storagePushMsgList, v)
					} else {
						if !(!isSenderSync && msgChannelValue.aggregationID == v.message.MsgData.SendID) {
							notStoragePushMsgList = append(notStoragePushMsgList, v)
						}
					}
					if v.message.MsgData.ContentType == constant.ReactionMessageModifier || v.message.MsgData.ContentType == constant.ReactionMessageDeleter {
						modifyMsgList = append(modifyMsgList, v.message)
					}
				}
				if len(modifyMsgList) > 0 {
					och.msgDatabase.MsgToModifyMQ(ctx, msgChannelValue.aggregationID, "", modifyMsgList)
				}
				log.ZDebug(ctx, "msg storage length", len(storageMsgList), "push length", len(notStoragePushMsgList))
				if len(storageMsgList) > 0 {
					lastSeq, err := och.msgDatabase.BatchInsertChat2Cache(ctx, msgChannelValue.aggregationID, storageMsgList)
					if err != nil {
						log.ZError(ctx, "batch data insert to redis err", err, "storageMsgList", storageMsgList)
						och.singleMsgFailedCountMutex.Lock()
						och.singleMsgFailedCount += uint64(len(storageMsgList))
						och.singleMsgFailedCountMutex.Unlock()
					} else {
						och.singleMsgSuccessCountMutex.Lock()
						och.singleMsgSuccessCount += uint64(len(storageMsgList))
						och.singleMsgSuccessCountMutex.Unlock()
						och.msgDatabase.MsgToMongoMQ(ctx, msgChannelValue.aggregationID, "", storageMsgList, lastSeq)
						for _, v := range storagePushMsgList {
							och.msgDatabase.MsgToPushMQ(v.ctx, msgChannelValue.aggregationID, v.message)
						}
						for _, v := range notStoragePushMsgList {
							och.msgDatabase.MsgToPushMQ(v.ctx, msgChannelValue.aggregationID, v.message)
						}
					}
				} else {
					for _, v := range notStoragePushMsgList {
						p, o, err := och.msgDatabase.MsgToPushMQ(v.ctx, msgChannelValue.aggregationID, v.message)
						if err != nil {
							log.ZError(v.ctx, "kafka send failed", err, "msg", v.message.String(), "pid", p, "offset", o)
						}
					}
				}
			}
		}
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
				log.ZDebug(ctx, "batch messages come to distribution center", len(consumerMessages))
				for i := 0; i < len(consumerMessages); i++ {
					ctxMsg := &ContextMsg{}
					msgFromMQ := pbMsg.MsgDataToMQ{}
					err := proto.Unmarshal(consumerMessages[i].Value, &msgFromMQ)
					if err != nil {
						log.ZError(ctx, "msg_transfer Unmarshal msg err", err, string(consumerMessages[i].Value))
						return
					}
					ctxMsg.ctx = kafka.GetContextWithMQHeader(consumerMessages[i].Headers)
					ctxMsg.message = &msgFromMQ
					log.ZDebug(ctx, "single msg come to distribution center", msgFromMQ.String(), string(consumerMessages[i].Key))
					if oldM, ok := aggregationMsgs[string(consumerMessages[i].Key)]; ok {
						oldM = append(oldM, ctxMsg)
						aggregationMsgs[string(consumerMessages[i].Key)] = oldM
					} else {
						m := make([]*ContextMsg, 0, 100)
						m = append(m, ctxMsg)
						aggregationMsgs[string(consumerMessages[i].Key)] = m
					}
				}
				log.ZDebug(ctx, "generate map list users len", len(aggregationMsgs))
				for aggregationID, v := range aggregationMsgs {
					if len(v) >= 0 {
						hashCode := utils.GetHashCode(aggregationID)
						channelID := hashCode % ChannelNum
						log.ZDebug(ctx, "generate channelID", "hashCode", hashCode, "channelID", channelID, "aggregationID", aggregationID)
						och.chArrays[channelID] <- Cmd2Value{Cmd: AggregationMessages, Value: MsgChannelValue{aggregationID: aggregationID, ctxMsgList: v, ctx: ctx}}
					}
				}
			}
		}
	}
}

func (OnlineHistoryRedisConsumerHandler) Setup(_ sarama.ConsumerGroupSession) error   { return nil }
func (OnlineHistoryRedisConsumerHandler) Cleanup(_ sarama.ConsumerGroupSession) error { return nil }

func (och *OnlineHistoryRedisConsumerHandler) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error { // a instance in the consumer group
	for {
		if sess == nil {
			log.NewWarn("", " sess == nil, waiting ")
			time.Sleep(100 * time.Millisecond)
		} else {
			break
		}
	}
	rwLock := new(sync.RWMutex)
	log.ZDebug(context.Background(), "online new session msg come", claim.HighWaterMarkOffset(), claim.Topic(), claim.Partition())
	cMsg := make([]*sarama.ConsumerMessage, 0, 1000)
	t := time.NewTicker(time.Duration(100) * time.Millisecond)
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
					log.ZDebug(ctx, "timer trigger msg consumer start", len(ccMsg))
					for i := 0; i < len(ccMsg)/split; i++ {
						//log.Debug()
						och.msgDistributionCh <- Cmd2Value{Cmd: ConsumerMsgs, Value: TriggerChannelValue{
							ctx: ctx, cMsgList: ccMsg[i*split : (i+1)*split]}}
					}
					if (len(ccMsg) % split) > 0 {
						och.msgDistributionCh <- Cmd2Value{Cmd: ConsumerMsgs, Value: TriggerChannelValue{
							ctx: ctx, cMsgList: ccMsg[split*(len(ccMsg)/split):]}}
					}
					log.ZDebug(ctx, "timer trigger msg consumer end", len(ccMsg))
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
