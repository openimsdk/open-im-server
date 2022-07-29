/*
** description("").
** copyright('open-im,www.open-im.io').
** author("fg,Gordon@open-im.io").
** time(2021/3/22 15:33).
 */
package logic

import (
	pusher "Open_IM/internal/push"
	fcm "Open_IM/internal/push/fcm"
	"Open_IM/internal/push/getui"
	jpush "Open_IM/internal/push/jpush"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/kafka"
	"Open_IM/pkg/statistics"
	"fmt"
)

var (
	rpcServer     RPCServer
	pushCh        PushConsumerHandler
	pushTerminal  []int32
	producer      *kafka.Producer
	offlinePusher pusher.OfflinePusher
	successCount  uint64
)

func Init(rpcPort int) {

	rpcServer.Init(rpcPort)
	pushCh.Init()
	pushTerminal = []int32{constant.IOSPlatformID, constant.AndroidPlatformID}
}
func init() {
	producer = kafka.NewKafkaProducer(config.Config.Kafka.Ws2mschat.Addr, config.Config.Kafka.Ws2mschat.Topic)
	statistics.NewStatistics(&successCount, config.Config.ModuleName.PushName, fmt.Sprintf("%d second push to msg_gateway count", constant.StatisticsTimeInterval), constant.StatisticsTimeInterval)
	if config.Config.Push.Getui.Enable {
		offlinePusher = getui.GetuiClient
	}
	if config.Config.Push.Jpns.Enable {
		offlinePusher = jpush.JPushClient
	}

	if config.Config.Push.Fcm.Enable {
		offlinePusher = fcm.NewFcm()
	}
}

func Run() {
	go rpcServer.run()
	go pushCh.pushConsumerGroup.RegisterHandleAndConsumer(&pushCh)
}
