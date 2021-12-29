/*
** description("").
** copyright('tuoyun,www.tuoyun.net').
** author("fg,Gordon@tuoyun.net").
** time(2021/5/11 9:36).
 */
package kafka

import (
	"context"
	"fmt"
	"sync"

	"Open_IM/pkg/common/mq"

	"github.com/Shopify/sarama"
)

type kafkaConsumerGroup struct {
	sarama.ConsumerGroup
	groupID string

	mu       *sync.RWMutex
	handlers map[string][]mq.MessageHandler
}

var _ mq.Consumer = (*kafkaConsumerGroup)(nil)

type MConsumerGroupConfig struct {
	KafkaVersion   sarama.KafkaVersion
	OffsetsInitial int64
	IsReturnErr    bool
}

func NewMConsumerGroup(consumerConfig *MConsumerGroupConfig, addr []string, groupID string) *kafkaConsumerGroup {
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
	return &kafkaConsumerGroup{
		ConsumerGroup: consumerGroup,
		groupID:       groupID,
		handlers:      make(map[string][]mq.MessageHandler),
	}
}

func (mc *kafkaConsumerGroup) RegisterMessageHandler(topic string, handler mq.MessageHandler) {
	mc.mu.Lock()
	defer mc.mu.Unlock()

	handlers := mc.handlers[topic]
	handlers = append(handlers, handler)
	mc.handlers[topic] = handlers
}

func (mc *kafkaConsumerGroup) Start() error {
	topics := make([]string, 0, len(mc.handlers))
	for topic := range mc.handlers {
		topics = append(topics, topic)
	}

	ctx := context.Background()
	for {
		err := mc.ConsumerGroup.Consume(ctx, topics, mc)
		if err != nil {
			panic(err)
		}
	}
}

func (mc *kafkaConsumerGroup) Setup(sarama.ConsumerGroupSession) error   { return nil }
func (mc *kafkaConsumerGroup) Cleanup(sarama.ConsumerGroupSession) error { return nil }
func (mc *kafkaConsumerGroup) ConsumeClaim(sess sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {

		mc.mu.RLock()
		handlers, ok := mc.handlers[msg.Topic]
		mc.mu.RUnlock()
		if !ok {
			panic(fmt.Sprintf("no handlers for topic: %s", msg.Topic))
		}

		message := &mq.Message{
			Key:       msg.Key,
			Value:     msg.Value,
			Topic:     msg.Topic,
			Partition: msg.Partition,
			Offset:    msg.Offset,
			Timestamp: msg.Timestamp,
		}
		for _, handler := range handlers {
			for {
				if err := handler.HandleMessage(message); err == nil { // error is nil, auto commit
					sess.MarkMessage(msg, "")
					break
				}
			}
		}
	}

	return nil
}
