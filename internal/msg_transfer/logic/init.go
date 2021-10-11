package logic

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/kafka"
	"Open_IM/pkg/common/log"
)

var (
	persistentCH PersistentConsumerHandler
	historyCH    HistoryConsumerHandler
	producer     *kafka.Producer
)

func Init() {
	log.NewPrivateLog(config.Config.ModuleName.MsgTransferName)
	persistentCH.Init()
	historyCH.Init()
	producer = kafka.NewKafkaProducer(config.Config.Kafka.Ms2pschat.Addr, config.Config.Kafka.Ms2pschat.Topic)
}
func Run() {
	//register mysqlConsumerHandler to
	go persistentCH.persistentConsumerGroup.RegisterHandleAndConsumer(&persistentCH)
	go historyCH.historyConsumerGroup.RegisterHandleAndConsumer(&historyCH)
}
