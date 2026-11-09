package mqbuild

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/tools/mq"
	"github.com/openimsdk/tools/mq/kafka"
	"github.com/openimsdk/tools/mq/redismq"
	"github.com/openimsdk/tools/mq/simmq"
)

type Builder interface {
	GetTopicProducer(ctx context.Context, topic string) (mq.Producer, error)
	GetTopicConsumer(ctx context.Context, topic string) (mq.Consumer, error)
}

const (
	TopicToRedis       = "toRedis"
	TopicToMongo       = "toMongo"
	TopicToPush        = "toPush"
	TopicToOfflinePush = "toOfflinePush"
)

const (
	GroupRedis       = "redis"
	GroupMongo       = "mongo"
	GroupPush        = "push"
	GroupOfflinePush = "offlinePush"
)

func NewBuilder(queue config.EngineSelector, kafka *config.Kafka, redis redis.UniversalClient) (Builder, error) {
	engine, err := config.ValidateQueueEngine(queue.Engine, config.Standalone())
	if err != nil {
		return nil, err
	}
	switch engine {
	case config.QueueEngineKafka:
		if kafka == nil {
			return nil, fmt.Errorf("nil kafka config")
		}
		return newKafkaBuilder(kafka), nil
	case config.QueueEngineRedis:
		return newRedisBuilder(redis), nil
	case config.QueueEngineMemory:
		return standaloneBuilder{}, nil
	default:
		return nil, fmt.Errorf("unsupported queue engine %s", queue.Engine)
	}
}

func newKafkaBuilder(kafka *config.Kafka) Builder {
	topics := MergeTopics(KafkaTopics(kafka))
	return &kafkaBuilder{
		addr:         kafka.Address,
		config:       kafka.Build(),
		logicalTopic: LogicalTopicNames(topics),
		topicGroupID: TopicGroupID(topics),
	}
}

func newRedisBuilder(redis redis.UniversalClient) Builder {
	return redismq.NewBuilder(redis, TopicGroupID(DefaultTopics()), redismq.Config{StreamPrefix: "mq:"})
}

type QueueTopicsConfig struct {
	ToRedis       QueueTopicConfig
	ToMongo       QueueTopicConfig
	ToPush        QueueTopicConfig
	ToOfflinePush QueueTopicConfig
}

type QueueTopicConfig struct {
	Topic   string
	GroupID string
}

func DefaultTopics() QueueTopicsConfig {
	return QueueTopicsConfig{
		ToRedis:       QueueTopicConfig{Topic: TopicToRedis, GroupID: GroupRedis},
		ToMongo:       QueueTopicConfig{Topic: TopicToMongo, GroupID: GroupMongo},
		ToPush:        QueueTopicConfig{Topic: TopicToPush, GroupID: GroupPush},
		ToOfflinePush: QueueTopicConfig{Topic: TopicToOfflinePush, GroupID: GroupOfflinePush},
	}
}

func KafkaTopics(kafka *config.Kafka) QueueTopicsConfig {
	if kafka == nil {
		return QueueTopicsConfig{}
	}
	return QueueTopicsConfig{
		ToRedis:       QueueTopicConfig{Topic: kafka.ToRedisTopic, GroupID: kafka.ToRedisGroupID},
		ToMongo:       QueueTopicConfig{Topic: kafka.ToMongoTopic, GroupID: kafka.ToMongoGroupID},
		ToPush:        QueueTopicConfig{Topic: kafka.ToPushTopic, GroupID: kafka.ToPushGroupID},
		ToOfflinePush: QueueTopicConfig{Topic: kafka.ToOfflinePushTopic, GroupID: kafka.ToOfflineGroupID},
	}
}

func MergeTopics(topics QueueTopicsConfig) QueueTopicsConfig {
	defaults := DefaultTopics()
	fillTopic(&topics.ToRedis, defaults.ToRedis)
	fillTopic(&topics.ToMongo, defaults.ToMongo)
	fillTopic(&topics.ToPush, defaults.ToPush)
	fillTopic(&topics.ToOfflinePush, defaults.ToOfflinePush)
	return topics
}

func fillTopic(topic *QueueTopicConfig, defaultTopic QueueTopicConfig) {
	if topic.Topic == "" {
		topic.Topic = defaultTopic.Topic
	}
	if topic.GroupID == "" {
		topic.GroupID = defaultTopic.GroupID
	}
}

func LogicalTopicNames(topics QueueTopicsConfig) map[string]string {
	return map[string]string{
		TopicToRedis:       topics.ToRedis.Topic,
		TopicToMongo:       topics.ToMongo.Topic,
		TopicToPush:        topics.ToPush.Topic,
		TopicToOfflinePush: topics.ToOfflinePush.Topic,
	}
}

func TopicGroupID(topics QueueTopicsConfig) map[string]string {
	return map[string]string{
		topics.ToRedis.Topic:       topics.ToRedis.GroupID,
		topics.ToMongo.Topic:       topics.ToMongo.GroupID,
		topics.ToPush.Topic:        topics.ToPush.GroupID,
		topics.ToOfflinePush.Topic: topics.ToOfflinePush.GroupID,
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
	logicalTopic map[string]string
	topicGroupID map[string]string
}

func (x *kafkaBuilder) GetTopicProducer(ctx context.Context, topic string) (mq.Producer, error) {
	realTopic, ok := x.logicalTopic[topic]
	if !ok {
		return nil, fmt.Errorf("topic %s not found", topic)
	}
	return kafka.NewKafkaProducerV2(x.config, x.addr, realTopic)
}

func (x *kafkaBuilder) GetTopicConsumer(ctx context.Context, topic string) (mq.Consumer, error) {
	realTopic, ok := x.logicalTopic[topic]
	if !ok {
		return nil, fmt.Errorf("topic %s not found", topic)
	}
	groupID, ok := x.topicGroupID[realTopic]
	if !ok {
		return nil, fmt.Errorf("topic %s groupID not found", realTopic)
	}
	return kafka.NewMConsumerGroupV2(ctx, x.config, groupID, []string{realTopic}, true)
}
