package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dtm-labs/rockscache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/localcache"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/datautil"
	"github.com/redis/go-redis/v9"
)

const (
	rocksCacheTimeout = 11 * time.Second
)

// BatchDeleterRedis is a concrete implementation of the BatchDeleter interface based on Redis and RocksCache.
type BatchDeleterRedis struct {
	redisClient    redis.UniversalClient
	keys           []string
	rocksClient    *rockscache.Client
	redisPubTopics []string
}

// NewBatchDeleterRedis creates a new BatchDeleterRedis instance.
func NewBatchDeleterRedis(rcClient *rocksCacheClient, redisPubTopics []string) *BatchDeleterRedis {
	return &BatchDeleterRedis{
		redisClient:    rcClient.GetRedis(),
		rocksClient:    rcClient.GetClient(),
		redisPubTopics: redisPubTopics,
	}
}

// ExecDelWithKeys directly takes keys for batch deletion and publishes deletion information.
func (c *BatchDeleterRedis) ExecDelWithKeys(ctx context.Context, keys []string) error {
	distinctKeys := datautil.Distinct(keys)
	return c.execDel(ctx, distinctKeys)
}

// ChainExecDel is used for chain calls for batch deletion. It must call Clone to prevent memory pollution.
func (c *BatchDeleterRedis) ChainExecDel(ctx context.Context) error {
	distinctKeys := datautil.Distinct(c.keys)
	return c.execDel(ctx, distinctKeys)
}

// execDel performs batch deletion and publishes the keys that have been deleted to update the local cache information of other nodes.
func (c *BatchDeleterRedis) execDel(ctx context.Context, keys []string) error {
	if len(keys) > 0 {
		log.ZDebug(ctx, "delete cache", "topic", c.redisPubTopics, "keys", keys)
		// Batch delete keys
		err := ProcessKeysBySlot(ctx, c.redisClient, keys, func(ctx context.Context, slot int64, keys []string) error {
			return c.rocksClient.TagAsDeletedBatch2(ctx, keys)
		})
		if err != nil {
			return err
		}
		// Publish the keys that have been deleted to Redis to update the local cache information of other nodes
		if len(c.redisPubTopics) > 0 && len(keys) > 0 {
			keysByTopic := localcache.GetPublishKeysByTopic(c.redisPubTopics, keys)
			for topic, keys := range keysByTopic {
				if len(keys) > 0 {
					data, err := json.Marshal(keys)
					if err != nil {
						log.ZWarn(ctx, "keys json marshal failed", err, "topic", topic, "keys", keys)
					} else {
						if err := c.redisClient.Publish(ctx, topic, string(data)).Err(); err != nil {
							log.ZWarn(ctx, "redis publish cache delete error", err, "topic", topic, "keys", keys)
						}
					}
				}
			}
		}
	}
	return nil
}

// Clone creates a copy of BatchDeleterRedis for chain calls to prevent memory pollution.
func (c *BatchDeleterRedis) Clone() cache.BatchDeleter {
	return &BatchDeleterRedis{
		redisClient:    c.redisClient,
		keys:           c.keys,
		rocksClient:    c.rocksClient,
		redisPubTopics: c.redisPubTopics,
	}
}

// AddKeys adds keys to be deleted.
func (c *BatchDeleterRedis) AddKeys(keys ...string) {
	c.keys = append(c.keys, keys...)
}

type disableBatchDeleter struct{}

func (x disableBatchDeleter) ChainExecDel(ctx context.Context) error {
	return nil
}

func (x disableBatchDeleter) ExecDelWithKeys(ctx context.Context, keys []string) error {
	return nil
}

func (x disableBatchDeleter) Clone() cache.BatchDeleter {
	return x
}

func (x disableBatchDeleter) AddKeys(keys ...string) {}

func getCache[T any](ctx context.Context, rcClient *rocksCacheClient, key string, expire time.Duration, fn func(ctx context.Context) (T, error)) (T, error) {
	if rcClient.Disable() {
		return fn(ctx)
	}
	var t T
	var write bool
	v, err := rcClient.GetClient().Fetch2(ctx, key, expire, func() (s string, err error) {
		t, err = fn(ctx)
		if err != nil {
			//log.ZError(ctx, "getCache query database failed", err, "key", key)
			return "", err
		}
		bs, err := json.Marshal(t)
		if err != nil {
			return "", errs.WrapMsg(err, "marshal failed")
		}
		write = true

		return string(bs), nil
	})
	if err != nil {
		return t, errs.Wrap(err)
	}
	if write {
		return t, nil
	}
	if v == "" {
		return t, errs.ErrRecordNotFound.WrapMsg("cache is not found")
	}
	err = json.Unmarshal([]byte(v), &t)
	if err != nil {
		errInfo := fmt.Sprintf("cache json.Unmarshal failed, key:%s, value:%s, expire:%s", key, v, expire)
		return t, errs.WrapMsg(err, errInfo)
	}

	return t, nil
}
