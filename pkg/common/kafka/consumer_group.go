/*
** description("").
** copyright('tuoyun,www.tuoyun.net').
** author("fg,Gordon@tuoyun.net").
** time(2021/5/11 9:36).
 */
package kafka

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"

	"github.com/Shopify/sarama"
)

type MConsumerGroup struct {
	sarama.ConsumerGroup
	groupID string
	topics  []string
}

type MConsumerGroupConfig struct {
	KafkaVersion   sarama.KafkaVersion
	OffsetsInitial int64
	IsReturnErr    bool
}

func NewMConsumerGroup(consumerConfig *MConsumerGroupConfig, topics, addrs []string, groupID string) *MConsumerGroup {
	config := sarama.NewConfig()
	config.Version = consumerConfig.KafkaVersion
	config.Consumer.Offsets.Initial = consumerConfig.OffsetsInitial
	config.Consumer.Return.Errors = consumerConfig.IsReturnErr
	consumerGroup, err := sarama.NewConsumerGroup(addrs, groupID, config)
	if err != nil {
		panic(err.Error())
	}
	return &MConsumerGroup{
		consumerGroup,
		groupID,
		topics,
	}
}

func (mc *MConsumerGroup) GetContextFromMsg(cMsg *sarama.ConsumerMessage) context.Context {
	return GetContextWithMQHeader(cMsg.Headers)

}

func (mc *MConsumerGroup) RegisterHandleAndConsumer(handler sarama.ConsumerGroupHandler) {
	log.ZDebug(context.Background(), "register consumer group", "groupID", mc.groupID)
	ctx := context.Background()
	for {
		err := mc.ConsumerGroup.Consume(ctx, mc.topics, handler)
		if err != nil {
			panic(err.Error())
		}
	}
}
