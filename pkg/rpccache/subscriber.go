package rpccache

import (
	"context"
	"encoding/json"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/redis/go-redis/v9"
)

func subscriberRedisDeleteCache(ctx context.Context, client redis.UniversalClient, channel string, del func(ctx context.Context, key ...string)) {
	defer func() {
		if r := recover(); r != nil {
			log.ZPanic(ctx, "subscriberRedisDeleteCache Panic", errs.ErrPanic(r))
		}
	}()
	for message := range client.Subscribe(ctx, channel).Channel() {
		log.ZDebug(ctx, "subscriberRedisDeleteCache", "channel", channel, "payload", message.Payload)
		var keys []string
		if err := json.Unmarshal([]byte(message.Payload), &keys); err != nil {
			log.ZError(ctx, "subscriberRedisDeleteCache json.Unmarshal error", err)
			continue
		}
		if len(keys) == 0 {
			continue
		}
		del(ctx, keys...)
	}
}
