package redispubsub

import "github.com/redis/go-redis/v9"

type Publisher struct {
	client  redis.UniversalClient
	channel string
}

func NewPublisher(client redis.UniversalClient, channel string) *Publisher {
	return &Publisher{client: client, channel: channel}
}

func (p *Publisher) Publish(message string) error {
	return p.client.Publish(ctx, p.channel, message).Err()
}
