package kafka

import (
	"github.com/Shopify/sarama"
	"sync"
)

type Consumer struct {
	addr          []string
	WG            sync.WaitGroup
	Topic         string
	PartitionList []int32
	Consumer      sarama.Consumer
}

func NewKafkaConsumer(addr []string, topic string) *Consumer {
	p := Consumer{}
	p.Topic = topic
	p.addr = addr

	consumer, err := sarama.NewConsumer(p.addr, nil)
	if err != nil {
		panic(err.Error())
		return nil
	}
	p.Consumer = consumer

	partitionList, err := consumer.Partitions(p.Topic)
	if err != nil {
		panic(err.Error())
		return nil
	}
	p.PartitionList = partitionList

	return &p
}
