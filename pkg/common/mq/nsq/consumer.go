package nsq

import (
	"Open_IM/pkg/common/mq"
	"github.com/nsqio/go-nsq"
)

type nsqConsumer struct {
	lookupAddrs []string
	topic       string

	handlers []mq.MessageHandler
	consumer *nsq.Consumer
}

var _ mq.Consumer = (*nsqConsumer)(nil)

func NewNsqConsumer(lookupAddrs []string, topic, channel string) (*nsqConsumer, error) {
	config := nsq.NewConfig()
	consumer, err := nsq.NewConsumer(topic, channel, config)
	if err != nil {
		return nil, err
	}

	nc := &nsqConsumer{
		lookupAddrs: lookupAddrs,
		topic:       topic,
		handlers:    make([]mq.MessageHandler, 0),
		consumer:    consumer,
	}

	consumer.AddHandler(nsq.HandlerFunc(nc.consume))

	return nc, nil
}

func (c *nsqConsumer) RegisterMessageHandler(topic string, handler mq.MessageHandler) {
	if topic != c.topic {
		return
	}
	c.handlers = append(c.handlers, handler)
}

func (c *nsqConsumer) consume(msg *nsq.Message) error {
	for _, handler := range c.handlers {
		if err := handler.HandleMessage(&mq.Message{
			Value: msg.Body,
		}); err != nil {
			return err
		}
	}
	msg.Finish()

	return nil
}

func (c *nsqConsumer) Start() error {

	if err := c.consumer.ConnectToNSQLookupds(c.lookupAddrs); err != nil {
		return err
	}

	<-c.consumer.StopChan

	return nil
}
