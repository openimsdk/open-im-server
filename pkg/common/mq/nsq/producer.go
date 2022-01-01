package nsq

import (
	"Open_IM/pkg/common/mq"
	"github.com/golang/protobuf/proto"
	"github.com/nsqio/go-nsq"
)

type nsqProducer struct {
	topic string

	producer *nsq.Producer
}

var _ mq.Producer = (*nsqProducer)(nil)

func NewNsqProducer(addr string, topic string) (*nsqProducer, error) {
	config := nsq.NewConfig()
	producer, err := nsq.NewProducer(addr, config)
	if err != nil {
		return nil, err
	}

	return &nsqProducer{
		topic:    topic,
		producer: producer,
	}, nil

}

func (p *nsqProducer) SendMessage(message proto.Message, key ...string) (int32, int64, error) {
	bytes, err := proto.Marshal(message)
	if err != nil {
		return 0, 0, err
	}
	return 0, 0, p.producer.Publish(p.topic, bytes)
}
