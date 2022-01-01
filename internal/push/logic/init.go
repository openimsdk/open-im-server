/*
** description("").
** copyright('open-im,www.open-im.io').
** author("fg,Gordon@open-im.io").
** time(2021/3/22 15:33).
 */
package logic

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/mq"
	"Open_IM/pkg/common/mq/kafka"
	"Open_IM/pkg/common/mq/nsq"
	"Open_IM/pkg/utils"
)

var (
	rpcServer    RPCServer
	pushCh       PushConsumerHandler
	pushTerminal []int32
	producer     mq.Producer
)

func Init(rpcPort int) {
	log.NewPrivateLog(config.Config.ModuleName.PushName)
	rpcServer.Init(rpcPort)
	pushCh.Init()
	pushTerminal = []int32{utils.IOSPlatformID, utils.AndroidPlatformID}
}
func init() {
	cfg := config.Config.MQ.Ws2mschat
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
	go rpcServer.run()
	go pushCh.pushConsumerGroup.Start()
}
