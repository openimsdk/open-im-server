/*
** description("").
** copyright('open-im,www.open-im.io').
** author("fg,Gordon@open-im.io").
** time(2021/3/22 15:33).
 */
package push

import (
	"OpenIM/pkg/common/config"
	"OpenIM/pkg/common/constant"
	"OpenIM/pkg/common/prome"
	"OpenIM/pkg/statistics"
	"fmt"
)

type Consumer struct {
	pushCh       ConsumerHandler
	successCount uint64
}

func NewConsumer(pusher *Pusher) *Consumer {
	return &Consumer{
		pushCh: *NewConsumerHandler(pusher),
	}
}

func (c *Consumer) initPrometheus() {
	prome.NewMsgOfflinePushSuccessCounter()
	prome.NewMsgOfflinePushFailedCounter()
}

func (c *Consumer) Start() {
	statistics.NewStatistics(&c.successCount, config.Config.ModuleName.PushName, fmt.Sprintf("%d second push to msg_gateway count", constant.StatisticsTimeInterval), constant.StatisticsTimeInterval)
	go c.pushCh.pushConsumerGroup.RegisterHandleAndConsumer(&c.pushCh)
}
