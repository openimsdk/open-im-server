/*
** description("").
** copyright('tuoyun,www.tuoyun.net').
** author("fg,Gordon@tuoyun.net").
** time(2021/5/11 9:36).
 */
package kafka

import (
	"context"
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

func NewMConsumerGroup(consumerConfig *MConsumerGroupConfig, topics, addr []string, groupID string) *MConsumerGroup {
	config := sarama.NewConfig()
	config.Version = consumerConfig.KafkaVersion
	config.Consumer.Offsets.Initial = consumerConfig.OffsetsInitial
	config.Consumer.Return.Errors = consumerConfig.IsReturnErr
	client, err := sarama.NewClient(addr, config)
	if err != nil {
		panic(err)
	}
	consumerGroup, err := sarama.NewConsumerGroupFromClient(groupID, client)
	if err != nil {
		panic(err)
	}
	return &MConsumerGroup{
		consumerGroup,
		groupID,
		topics,
	}
}
func (mc *MConsumerGroup) RegisterHandleAndConsumer(handler sarama.ConsumerGroupHandler) {
	ctx := context.Background()
	for {
		err := mc.ConsumerGroup.Consume(ctx, mc.topics, handler)
		if err != nil {
			panic(err)
		}
	}
}
