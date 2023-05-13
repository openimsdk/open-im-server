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

package logic

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/kafka"
	promePkg "Open_IM/pkg/common/prometheus"
	"Open_IM/pkg/statistics"
	"fmt"
	"sync"
)

const OnlineTopicBusy = 1
const OnlineTopicVacancy = 0
const Msg = 2
const ConsumerMsgs = 3
const AggregationMessages = 4
const MongoMessages = 5
const ChannelNum = 100

var (
	persistentCH          PersistentConsumerHandler
	historyCH             OnlineHistoryRedisConsumerHandler
	historyMongoCH        OnlineHistoryMongoConsumerHandler
	modifyCH              ModifyMsgConsumerHandler
	producer              *kafka.Producer
	producerToModify      *kafka.Producer
	producerToMongo       *kafka.Producer
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
	if config.Config.Prometheus.Enable {
		initPrometheus()
	}
	persistentCH.Init()   // ws2mschat save mysql
	historyCH.Init(cmdCh) //
	historyMongoCH.Init()
	modifyCH.Init()
	onlineTopicStatus = OnlineTopicVacancy
	//offlineHistoryCH.Init(cmdCh)
	statistics.NewStatistics(&singleMsgSuccessCount, config.Config.ModuleName.MsgTransferName, fmt.Sprintf("%d second singleMsgCount insert to mongo", constant.StatisticsTimeInterval), constant.StatisticsTimeInterval)
	statistics.NewStatistics(&groupMsgCount, config.Config.ModuleName.MsgTransferName, fmt.Sprintf("%d second groupMsgCount insert to mongo", constant.StatisticsTimeInterval), constant.StatisticsTimeInterval)
	producer = kafka.NewKafkaProducer(config.Config.Kafka.Ms2pschat.Addr, config.Config.Kafka.Ms2pschat.Topic)
	producerToModify = kafka.NewKafkaProducer(config.Config.Kafka.MsgToModify.Addr, config.Config.Kafka.MsgToModify.Topic)
	producerToMongo = kafka.NewKafkaProducer(config.Config.Kafka.MsgToMongo.Addr, config.Config.Kafka.MsgToMongo.Topic)
}
func Run(promethuesPort int) {
	//register mysqlConsumerHandler to
	if config.Config.ChatPersistenceMysql {
		go persistentCH.persistentConsumerGroup.RegisterHandleAndConsumer(&persistentCH)
	} else {
		fmt.Println("not start mysql consumer")
	}
	go historyCH.historyConsumerGroup.RegisterHandleAndConsumer(&historyCH)
	go historyMongoCH.historyConsumerGroup.RegisterHandleAndConsumer(&historyMongoCH)
	go modifyCH.modifyMsgConsumerGroup.RegisterHandleAndConsumer(&modifyCH)
	//go offlineHistoryCH.historyConsumerGroup.RegisterHandleAndConsumer(&offlineHistoryCH)
	go func() {
		err := promePkg.StartPromeSrv(promethuesPort)
		if err != nil {
			panic(err)
		}
	}()
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
