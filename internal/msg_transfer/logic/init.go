package logic

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/mq"
	"Open_IM/pkg/common/mq/kafka"
	"Open_IM/pkg/common/mq/nsq"
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

	cfg := config.Config.MQ.Ms2pschat
	switch cfg.Type {
	case "kafka":
		producer = kafka.NewKafkaProducer(cfg.Addr, cfg.Topic)
	case "nsq":
		p, err := nsq.NewNsqProducer(cfg.Addr[0], cfg.Topic)
		if err != nil {
			panic(err)
		}
		producer = p
	}
}

func Run() {
	//register mysqlConsumerHandler to
	go persistentCH.persistentConsumerGroup.Start()
	go historyCH.historyConsumerGroup.Start()
}
