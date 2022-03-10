/*
** description("").
** copyright('open-im,www.open-im.io').
** author("fg,Gordon@open-im.io").
** time(2021/3/22 15:33).
 */
package logic

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/kafka"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/statistics"
	"fmt"
)

var (
	rpcServer    RPCServer
	pushCh       PushConsumerHandler
	pushTerminal []int32
	producer     *kafka.Producer
	count        uint64
)

func Init(rpcPort int) {
	log.NewPrivateLog(config.Config.ModuleName.PushName)
	rpcServer.Init(rpcPort)
	pushCh.Init()
	pushTerminal = []int32{constant.IOSPlatformID, constant.AndroidPlatformID}
}
func init() {
	producer = kafka.NewKafkaProducer(config.Config.Kafka.Ws2mschat.Addr, config.Config.Kafka.Ws2mschat.Topic)
	statistics.NewStatistics(&count, config.Config.ModuleName.PushName, fmt.Sprintf("%d second push to msg_gateway count", 300), 300)
}

func Run() {
	go rpcServer.run()
	go pushCh.pushConsumerGroup.RegisterHandleAndConsumer(&pushCh)
}
