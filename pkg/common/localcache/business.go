package localcache

import (
	"context"
	"encoding/json"
	"github.com/OpenIMSDK/tools/log"
	"github.com/dtm-labs/rockscache"
	"github.com/redis/go-redis/v9"
)

func WithRedisDeleteSubscribe(topic string, cli redis.UniversalClient) Option {
	return WithDeleteLocal(func(fn func(key ...string)) {
		if fn == nil {
			log.ZDebug(context.Background(), "WithRedisDeleteSubscribe fn is nil", "topic", topic)
			return
		}
		msg := cli.Subscribe(context.Background(), topic).Channel()
		for m := range msg {
			log.ZDebug(context.Background(), "WithRedisDeleteSubscribe delete", "topic", m.Channel, "payload", m.Payload)
			var key []string
			if err := json.Unmarshal([]byte(m.Payload), &key); err != nil {
				log.ZError(context.Background(), "WithRedisDeleteSubscribe json unmarshal error", err, "topic", topic, "payload", m.Payload)
				continue
			}
			if len(key) == 0 {
				continue
			}
			fn(key...)
		}
	})
}

func WithRedisDeletePublish(topic string, cli redis.UniversalClient) Option {
	return WithDeleteKeyBefore(func(ctx context.Context, key ...string) {
		data, err := json.Marshal(key)
		if err != nil {
			log.ZError(ctx, "json marshal error", err, "topic", topic, "key", key)
			return
		}
		if err := cli.Publish(ctx, topic, data).Err(); err != nil {
			log.ZError(ctx, "redis publish error", err, "topic", topic, "key", key)
		} else {
			log.ZDebug(ctx, "redis publish success", "topic", topic, "key", key)
		}
	})
}

func WithRedisDelete(cli redis.UniversalClient) Option {
	return WithDeleteKeyBefore(func(ctx context.Context, key ...string) {
		for _, s := range key {
			if err := cli.Del(ctx, s).Err(); err != nil {
				log.ZError(ctx, "redis delete error", err, "key", s)
			} else {
				log.ZDebug(ctx, "redis delete success", "key", s)
			}
		}
	})
}

func WithRocksCacheDelete(cli *rockscache.Client) Option {
	return WithDeleteKeyBefore(func(ctx context.Context, key ...string) {
		for _, k := range key {
			if err := cli.TagAsDeleted2(ctx, k); err != nil {
				log.ZError(ctx, "rocksdb delete error", err, "key", k)
			} else {
				log.ZDebug(ctx, "rocksdb delete success", "key", k)
			}
		}
	})
}
