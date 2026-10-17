package redis

import (
	"context"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/redis/go-redis/v9"
)

const lockPrefix = "LOCK:"

func NewLock(rdb redis.UniversalClient) cache.Lock {
	return &redisLock{rdb: rdb}
}

type redisLock struct {
	rdb redis.UniversalClient
}

func (x *redisLock) Lock(ctx context.Context, key string, timeout time.Duration) (string, error) {
	uid, err := uuid.NewUUID()
	if err != nil {
		return "", err
	}
	if timeout < time.Second {
		timeout = time.Minute * 2
	}
	value := hex.EncodeToString(uid[:])
	key = lockPrefix + key
	for {
		ok, err := x.rdb.SetNX(ctx, key, value, timeout).Result()
		if err != nil {
			return "", errs.WrapMsg(err, "get redis lock", "key", key)
		}
		if ok {
			return value, nil
		}
		timer := time.NewTimer(50 * time.Millisecond)
		select {
		case <-ctx.Done():
			timer.Stop()
			return "", context.Cause(ctx)
		case <-timer.C:
		}
	}
}

func (x *redisLock) Unlock(ctx context.Context, key, value string) {
	ctx, cancel := context.WithTimeout(context.WithoutCancel(ctx), 10*time.Second)
	defer cancel()
	script := "\nlocal value = redis.call(\"GET\", KEYS[1])\nif value == ARGV[1] then\n  return redis.call(\"DEL\", KEYS[1])\nend\nreturn 0"
	if err := x.rdb.Eval(ctx, script, []string{lockPrefix + key}, value).Err(); err != nil {
		log.ZWarn(ctx, "unlock redis lock", err, "key", key)
	}
}
