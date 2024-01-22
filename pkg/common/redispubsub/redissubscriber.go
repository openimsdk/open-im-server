package redispubsub

import (
	"context"
	"github.com/redis/go-redis/v9"
)

var ctx = context.Background()

type Subscriber struct {
	client  redis.UniversalClient
	channel string
}

func NewSubscriber(client redis.UniversalClient, channel string) *Subscriber {
	return &Subscriber{client: client, channel: channel}
}

func (s *Subscriber) OnMessage(ctx context.Context, callback func(string)) error {
	messageChannel := s.client.Subscribe(ctx, s.channel).Channel()

	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case msg := <-messageChannel:
				callback(msg.Payload)
			}
		}
	}()

	return nil
}
