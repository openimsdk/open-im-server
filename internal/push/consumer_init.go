package push

import (
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/prome"
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
	//statistics.NewStatistics(&c.successCount, config.Config.ModuleName.PushName, fmt.Sprintf("%d second push to msg_gateway count", constant.StatisticsTimeInterval), constant.StatisticsTimeInterval)
	go c.pushCh.pushConsumerGroup.RegisterHandleAndConsumer(&c.pushCh)
}
