/*
** description("").
** copyright('open-im,www.open-im.io').
** author("fg,Gordon@open-im.io").
** time(2021/3/22 15:33).
 */
package push

import (
	fcm "Open_IM/internal/push/fcm"
	"Open_IM/internal/push/getui"
	jpush "Open_IM/internal/push/jpush"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db/cache"
	"Open_IM/pkg/common/prome"
	"Open_IM/pkg/statistics"
	"fmt"
)

type Push struct {
	rpcServer     RPCServer
	pushCh        ConsumerHandler
	offlinePusher OfflinePusher
	successCount  uint64
}

func (p *Push) Init(rpcPort int) error {
	redisClient, err := cache.NewRedis()
	if err != nil {
		return err
	}
	var cacheInterface cache.Cache = redisClient
	p.rpcServer.Init(rpcPort, cacheInterface)
	p.pushCh.Init()
	statistics.NewStatistics(&p.successCount, config.Config.ModuleName.PushName, fmt.Sprintf("%d second push to msg_gateway count", constant.StatisticsTimeInterval), constant.StatisticsTimeInterval)
	if *config.Config.Push.Getui.Enable {
		p.offlinePusher = getui.NewClient(cacheInterface)
	}
	if config.Config.Push.Jpns.Enable {
		p.offlinePusher = jpush.NewClient()
	}
	if config.Config.Push.Fcm.Enable {
		p.offlinePusher = fcm.NewClient(cacheInterface)
	}
	return nil
}

func (p *Push) initPrometheus() {
	prome.NewMsgOfflinePushSuccessCounter()
	prome.NewMsgOfflinePushFailedCounter()
}

func (p *Push) Run(prometheusPort int) {
	go p.rpcServer.run()
	go p.pushCh.pushConsumerGroup.RegisterHandleAndConsumer(&p.pushCh)
	go func() {
		err := prome.StartPromeSrv(prometheusPort)
		if err != nil {
			panic(err)
		}
	}()
}
