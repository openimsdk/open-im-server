package mqbuild

import (
	"context"
	"fmt"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/tools/mq"
	"github.com/openimsdk/tools/mq/kafka"
	"github.com/openimsdk/tools/mq/simmq"
)

type Builder interface {
	GetTopicProducer(ctx context.Context, topic string) (mq.Producer, error)
	GetTopicConsumer(ctx context.Context, topic string) (mq.Consumer, error)
}

func NewBuilder(kafka *config.Kafka) Builder {
	if config.Standalone() {
		return standaloneBuilder{}
	}
	return &kafkaBuilder{
		addr:   kafka.Address,
		config: kafka.Build(),
		topicGroupID: map[string]string{
			kafka.ToRedisTopic:       kafka.ToRedisGroupID,
			kafka.ToMongoTopic:       kafka.ToMongoGroupID,
			kafka.ToPushTopic:        kafka.ToPushGroupID,
			kafka.ToOfflinePushTopic: kafka.ToOfflineGroupID,
		},
	}
}

type standaloneBuilder struct{}

func (standaloneBuilder) GetTopicProducer(ctx context.Context, topic string) (mq.Producer, error) {
	return simmq.GetTopicProducer(topic), nil
}

func (standaloneBuilder) GetTopicConsumer(ctx context.Context, topic string) (mq.Consumer, error) {
	return simmq.GetTopicConsumer(topic), nil
}

type kafkaBuilder struct {
	addr         []string
	config       *kafka.Config
	topicGroupID map[string]string
}

func (x *kafkaBuilder) GetTopicProducer(ctx context.Context, topic string) (mq.Producer, error) {
	return kafka.NewKafkaProducerV2(x.config, x.addr, topic)
}

func (x *kafkaBuilder) GetTopicConsumer(ctx context.Context, topic string) (mq.Consumer, error) {
	groupID, ok := x.topicGroupID[topic]
	if !ok {
		return nil, fmt.Errorf("topic %s groupID not found", topic)
	}
	return kafka.NewMConsumerGroupV2(ctx, x.config, groupID, []string{topic}, true)
}
