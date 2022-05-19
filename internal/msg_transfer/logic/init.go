package logic

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/kafka"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/statistics"
	"fmt"
	"sync"
)

const OnlineTopicBusy = 1
const OnlineTopicVacancy = 0
const Msg = 2

var (
	persistentCH          PersistentConsumerHandler
	historyCH             OnlineHistoryConsumerHandler
	offlineHistoryCH      OfflineHistoryConsumerHandler
	producer              *kafka.Producer
	cmdCh                 chan Cmd2Value
	onlineTopicStatus     int
	w                     *sync.Mutex
	singleMsgSuccessCount uint64
	groupMsgCount         uint64
	singleMsgFailedCount  uint64

	singleMsgSuccessCountMutex sync.Mutex
)

func Init() {
	cmdCh = make(chan Cmd2Value, 10000)
	w = new(sync.Mutex)
	persistentCH.Init()
	historyCH.Init(cmdCh)
	onlineTopicStatus = OnlineTopicVacancy
	log.Debug("come msg transfer ts", config.Config.Kafka.ConsumerGroupID.MsgToMongoOffline, config.Config.Kafka.Ws2mschatOffline.Topic)
	offlineHistoryCH.Init(cmdCh)
	statistics.NewStatistics(&singleMsgSuccessCount, config.Config.ModuleName.MsgTransferName, fmt.Sprintf("%d second singleMsgCount insert to mongo", constant.StatisticsTimeInterval), constant.StatisticsTimeInterval)
	statistics.NewStatistics(&groupMsgCount, config.Config.ModuleName.MsgTransferName, fmt.Sprintf("%d second groupMsgCount insert to mongo", constant.StatisticsTimeInterval), constant.StatisticsTimeInterval)
	producer = kafka.NewKafkaProducer(config.Config.Kafka.Ms2pschat.Addr, config.Config.Kafka.Ms2pschat.Topic)
}
func Run() {
	//register mysqlConsumerHandler to
	if config.Config.ChatPersistenceMysql {
		go persistentCH.persistentConsumerGroup.RegisterHandleAndConsumer(&persistentCH)
	} else {
		fmt.Println("not start mysql consumer")
	}
	go historyCH.historyConsumerGroup.RegisterHandleAndConsumer(&historyCH)
	go offlineHistoryCH.historyConsumerGroup.RegisterHandleAndConsumer(&offlineHistoryCH)
}
func SetOnlineTopicStatus(status int) {
	w.Lock()
	defer w.Unlock()
	onlineTopicStatus = status
}
func GetOnlineTopicStatus() int {
	w.Lock()
	defer w.Unlock()
	return onlineTopicStatus
}
