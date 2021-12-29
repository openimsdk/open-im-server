package logic

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/mq"
	"Open_IM/pkg/common/mq/kafka"
)

var (
	persistentCH PersistentConsumerHandler
	historyCH    HistoryConsumerHandler
	producer     mq.Producer
)

func Init() {
	log.NewPrivateLog(config.Config.ModuleName.MsgTransferName)
	persistentCH.Init()
	historyCH.Init()
	producer = kafka.NewKafkaProducer(config.Config.Kafka.Ms2pschat.Addr, config.Config.Kafka.Ms2pschat.Topic)
}

func Run() {
	//register mysqlConsumerHandler to
	go persistentCH.persistentConsumerGroup.Start()
	go historyCH.historyConsumerGroup.Start()
}
