/*
** description("").
** copyright('open-im,www.open-im.io').
** author("fg,Gordon@open-im.io").
** time(2021/3/22 15:33).
 */
package logic

import (
	"Open_IM/src/common/config"
	"Open_IM/src/common/kafka"
	"Open_IM/src/common/log"
	"Open_IM/src/utils"
)

var (
	rpcServer    RPCServer
	pushCh       PushConsumerHandler
	pushTerminal []int32
	producer     *kafka.Producer
)

func Init(rpcPort int) {
	log.NewPrivateLog(config.Config.ModuleName.PushName)
	rpcServer.Init(rpcPort)
	pushCh.Init()
	pushTerminal = []int32{utils.IOSPlatformID}
}
func init() {
	producer = kafka.NewKafkaProducer(config.Config.Kafka.Ws2mschat.Addr, config.Config.Kafka.Ws2mschat.Topic)
}

func Run() {
	go rpcServer.run()
	go pushCh.pushConsumerGroup.RegisterHandleAndConsumer(&pushCh)
}
