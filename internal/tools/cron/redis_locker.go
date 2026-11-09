package cron

import (
	"context"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"

	"github.com/openimsdk/tools/log"
)

func NewRedisLocker(client redis.UniversalClient) *RedisLocker {
	return &RedisLocker{
		client: client,
		script: redis.NewScript(strings.TrimSpace(`
if redis.call("get", KEYS[1]) == ARGV[1] then
    return redis.call("del", KEYS[1])
else
    return 0
end
`)),
	}
}

type RedisLocker struct {
	client redis.UniversalClient
	script *redis.Script
}

func (e *RedisLocker) getKey(name string) string {
	return "CRON_LOCKED:" + name
}

func (e *RedisLocker) lock(ctx context.Context, name string, owner string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	return e.client.SetNX(ctx, e.getKey(name), owner, lockLeaseTTL).Result()
}

func (e *RedisLocker) unlock(ctx context.Context, name string, owner string) error {
	ctx, cancel := context.WithTimeout(ctx, time.Second)
	defer cancel()
	return e.script.Run(ctx, e.client, []string{e.getKey(name)}, owner).Err()
}

func (e *RedisLocker) ExecuteWithLock(ctx context.Context, taskName string, task func()) {
	owner := uuid.New().String()
	ok, err := e.lock(ctx, taskName, owner)
	if err != nil {
		log.ZWarn(ctx, "cron lock get lock", err, "taskName", taskName)
		return
	}
	log.ZDebug(ctx, "cron lock get lock", "taskName", taskName, "ok", ok, "owner", owner)
	if !ok {
		return
	}
	defer func() {
		err := e.unlock(ctx, taskName, owner)
		if err == nil {
			log.ZDebug(ctx, "cron lock unlock", "taskName", taskName, "owner", owner)
		} else {
			log.ZWarn(ctx, "cron lock unlock", err, "taskName", taskName, "owner", owner)
		}
	}()
	task()
}
