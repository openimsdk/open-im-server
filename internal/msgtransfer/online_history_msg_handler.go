package msgtransfer

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db/cache"
	"Open_IM/pkg/common/db/controller"
	"Open_IM/pkg/common/kafka"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/tracelog"
	pbMsg "Open_IM/pkg/proto/msg"
	pbPush "Open_IM/pkg/proto/push"
	"Open_IM/pkg/statistics"
	"Open_IM/pkg/utils"
	"context"
	"fmt"
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
	triggerID     string
	msgList       []*pbMsg.MsgDataToMQ
	lastSeq       uint64
}

type TriggerChannelValue struct {
	triggerID string
	cMsgList  []*sarama.ConsumerMessage
}

type Cmd2Value struct {
	Cmd   int
	Value interface{}
}

type OnlineHistoryRedisConsumerHandler struct {
	historyConsumerGroup *kafka.MConsumerGroup
	chArrays             [ChannelNum]chan Cmd2Value
	msgDistributionCh    chan Cmd2Value

	singleMsgSuccessCount      uint64
	singleMsgFailedCount       uint64
	singleMsgSuccessCountMutex sync.Mutex
	singleMsgFailedCountMutex  sync.Mutex

	producerToPush   *kafka.Producer
	producerToModify *kafka.Producer
	producerToMongo  *kafka.Producer

	msgInterface controller.MsgInterface
	cache        cache.Cache
}

func (och *OnlineHistoryRedisConsumerHandler) Init() {
	och.msgDistributionCh = make(chan Cmd2Value) //no buffer channel
	go och.MessagesDistributionHandle()
	for i := 0; i < ChannelNum; i++ {
		och.chArrays[i] = make(chan Cmd2Value, 50)
		go och.Run(i)
	}
	och.producerToPush = kafka.NewKafkaProducer(config.Config.Kafka.Ms2pschat.Addr, config.Config.Kafka.Ms2pschat.Topic)
	och.producerToModify = kafka.NewKafkaProducer(config.Config.Kafka.MsgToModify.Addr, config.Config.Kafka.MsgToModify.Topic)
	och.producerToMongo = kafka.NewKafkaProducer(config.Config.Kafka.MsgToMongo.Addr, config.Config.Kafka.MsgToMongo.Topic)
	och.historyConsumerGroup = kafka.NewMConsumerGroup(&kafka.MConsumerGroupConfig{KafkaVersion: sarama.V2_0_0_0,
		OffsetsInitial: sarama.OffsetNewest, IsReturnErr: false}, []string{config.Config.Kafka.Ws2mschat.Topic},
		config.Config.Kafka.Ws2mschat.Addr, config.Config.Kafka.ConsumerGroupID.MsgToRedis)

	statistics.NewStatistics(&och.singleMsgSuccessCount, config.Config.ModuleName.MsgTransferName, fmt.Sprintf("%d second singleMsgCount insert to mongo", constant.StatisticsTimeInterval), constant.StatisticsTimeInterval)
}

func (och *OnlineHistoryRedisConsumerHandler) Run(channelID int) {
	for {
		select {
		case cmd := <-och.chArrays[channelID]:
			switch cmd.Cmd {
			case AggregationMessages:
				msgChannelValue := cmd.Value.(MsgChannelValue)
				msgList := msgChannelValue.msgList
				triggerID := msgChannelValue.triggerID
				storageMsgList := make([]*pbMsg.MsgDataToMQ, 0, 80)
				notStoragePushMsgList := make([]*pbMsg.MsgDataToMQ, 0, 80)
				log.Debug(triggerID, "msg arrived channel", "channel id", channelID, msgList, msgChannelValue.aggregationID, len(msgList))
				var modifyMsgList []*pbMsg.MsgDataToMQ
				ctx := context.Background()
				tracelog.SetOperationID(ctx, triggerID)
				for _, v := range msgList {
					log.Debug(triggerID, "msg come to storage center", v.String())
					isHistory := utils.GetSwitchFromOptions(v.MsgData.Options, constant.IsHistory)
					isSenderSync := utils.GetSwitchFromOptions(v.MsgData.Options, constant.IsSenderSync)
					if isHistory {
						storageMsgList = append(storageMsgList, v)
						//log.NewWarn(triggerID, "storageMsgList to mongodb  client msgID: ", v.MsgData.ClientMsgID)
					} else {
						if !(!isSenderSync && msgChannelValue.aggregationID == v.MsgData.SendID) {
							notStoragePushMsgList = append(notStoragePushMsgList, v)
						}
					}
					if v.MsgData.ContentType == constant.ReactionMessageModifier || v.MsgData.ContentType == constant.ReactionMessageDeleter {
						modifyMsgList = append(modifyMsgList, v)
					}
				}
				if len(modifyMsgList) > 0 {
					och.sendMessageToModifyMQ(ctx, msgChannelValue.aggregationID, triggerID, modifyMsgList)
				}
				log.Debug(triggerID, "msg storage length", len(storageMsgList), "push length", len(notStoragePushMsgList))
				if len(storageMsgList) > 0 {
					lastSeq, err := och.msgInterface.BatchInsertChat2Cache(ctx, msgChannelValue.aggregationID, storageMsgList)
					if err != nil {
						och.singleMsgFailedCountMutex.Lock()
						och.singleMsgFailedCount += uint64(len(storageMsgList))
						och.singleMsgFailedCountMutex.Unlock()
						log.NewError(triggerID, "single data insert to redis err", err.Error(), storageMsgList)
					} else {
						och.singleMsgSuccessCountMutex.Lock()
						och.singleMsgSuccessCount += uint64(len(storageMsgList))
						och.singleMsgSuccessCountMutex.Unlock()
						och.SendMessageToMongoCH(ctx, msgChannelValue.aggregationID, triggerID, storageMsgList, lastSeq)
						for _, v := range storageMsgList {
							och.sendMessageToPushMQ(ctx, v, msgChannelValue.aggregationID)
						}
						for _, x := range notStoragePushMsgList {
							och.sendMessageToPushMQ(ctx, x, msgChannelValue.aggregationID)
						}
					}
				} else {
					for _, v := range notStoragePushMsgList {
						och.sendMessageToPushMQ(ctx, v, msgChannelValue.aggregationID)
					}
				}
			}
		}
	}
}

func (och *OnlineHistoryRedisConsumerHandler) MessagesDistributionHandle() {
	for {
		aggregationMsgs := make(map[string][]*pbMsg.MsgDataToMQ, ChannelNum)
		select {
		case cmd := <-och.msgDistributionCh:
			switch cmd.Cmd {
			case ConsumerMsgs:
				triggerChannelValue := cmd.Value.(TriggerChannelValue)
				triggerID := triggerChannelValue.triggerID
				consumerMessages := triggerChannelValue.cMsgList
				//Aggregation map[userid]message list
				log.Debug(triggerID, "batch messages come to distribution center", len(consumerMessages))
				for i := 0; i < len(consumerMessages); i++ {
					msgFromMQ := pbMsg.MsgDataToMQ{}
					err := proto.Unmarshal(consumerMessages[i].Value, &msgFromMQ)
					if err != nil {
						log.Error(triggerID, "msg_transfer Unmarshal msg err", "msg", string(consumerMessages[i].Value), "err", err.Error())
						return
					}
					log.Debug(triggerID, "single msg come to distribution center", msgFromMQ.String(), string(consumerMessages[i].Key))
					if oldM, ok := aggregationMsgs[string(consumerMessages[i].Key)]; ok {
						oldM = append(oldM, &msgFromMQ)
						aggregationMsgs[string(consumerMessages[i].Key)] = oldM
					} else {
						m := make([]*pbMsg.MsgDataToMQ, 0, 100)
						m = append(m, &msgFromMQ)
						aggregationMsgs[string(consumerMessages[i].Key)] = m
					}
				}
				log.Debug(triggerID, "generate map list users len", len(aggregationMsgs))
				for aggregationID, v := range aggregationMsgs {
					if len(v) >= 0 {
						hashCode := utils.GetHashCode(aggregationID)
						channelID := hashCode % ChannelNum
						log.Debug(triggerID, "generate channelID", hashCode, channelID, aggregationID)
						och.chArrays[channelID] <- Cmd2Value{Cmd: AggregationMessages, Value: MsgChannelValue{aggregationID: aggregationID, msgList: v, triggerID: triggerID}}
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
	log.NewDebug("", "online new session msg come", claim.HighWaterMarkOffset(), claim.Topic(), claim.Partition())
	cMsg := make([]*sarama.ConsumerMessage, 0, 1000)
	t := time.NewTicker(time.Duration(100) * time.Millisecond)
	var triggerID string
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
					triggerID = utils.OperationIDGenerator()
					log.Debug(triggerID, "timer trigger msg consumer start", len(ccMsg))
					for i := 0; i < len(ccMsg)/split; i++ {
						//log.Debug()
						och.msgDistributionCh <- Cmd2Value{Cmd: ConsumerMsgs, Value: TriggerChannelValue{
							triggerID: triggerID, cMsgList: ccMsg[i*split : (i+1)*split]}}
					}
					if (len(ccMsg) % split) > 0 {
						och.msgDistributionCh <- Cmd2Value{Cmd: ConsumerMsgs, Value: TriggerChannelValue{
							triggerID: triggerID, cMsgList: ccMsg[split*(len(ccMsg)/split):]}}
					}
					log.Debug(triggerID, "timer trigger msg consumer end", len(cMsg))
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

func (och *OnlineHistoryRedisConsumerHandler) sendMessageToPushMQ(ctx context.Context, message *pbMsg.MsgDataToMQ, pushToUserID string) {
	log.Info(message.OperationID, utils.GetSelfFuncName(), "msg ", message.String(), pushToUserID)
	rpcPushMsg := pbPush.PushMsgReq{OperationID: message.OperationID, MsgData: message.MsgData, PushToUserID: pushToUserID}
	mqPushMsg := pbMsg.PushMsgDataToMQ{OperationID: message.OperationID, MsgData: message.MsgData, PushToUserID: pushToUserID}
	pid, offset, err := och.producerToPush.SendMessage(&mqPushMsg, mqPushMsg.PushToUserID, rpcPushMsg.OperationID)
	if err != nil {
		log.Error(mqPushMsg.OperationID, "kafka send failed", "send data", message.String(), "pid", pid, "offset", offset, "err", err.Error())
	}
	return
}

func (och *OnlineHistoryRedisConsumerHandler) sendMessageToModifyMQ(ctx context.Context, aggregationID string, triggerID string, messages []*pbMsg.MsgDataToMQ) {
	if len(messages) > 0 {
		pid, offset, err := och.producerToModify.SendMessage(&pbMsg.MsgDataToModifyByMQ{AggregationID: aggregationID, MessageList: messages, TriggerID: triggerID}, aggregationID, triggerID)
		if err != nil {
			log.Error(triggerID, "kafka send failed", "send data", len(messages), "pid", pid, "offset", offset, "err", err.Error(), "key", aggregationID)
		}
	}
}

func (och *OnlineHistoryRedisConsumerHandler) SendMessageToMongoCH(ctx context.Context, aggregationID string, triggerID string, messages []*pbMsg.MsgDataToMQ, lastSeq uint64) {
	if len(messages) > 0 {
		pid, offset, err := och.producerToMongo.SendMessage(&pbMsg.MsgDataToMongoByMQ{LastSeq: lastSeq, AggregationID: aggregationID, MessageList: messages, TriggerID: triggerID}, aggregationID, triggerID)
		if err != nil {
			log.Error(triggerID, "kafka send failed", "send data", len(messages), "pid", pid, "offset", offset, "err", err.Error(), "key", aggregationID)
		}
	}
}
