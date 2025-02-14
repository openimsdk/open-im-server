package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/dtm-labs/rockscache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/redis/go-redis/v9"
)

// GetRocksCacheOptions returns the default configuration options for RocksCache.
func GetRocksCacheOptions() *rockscache.Options {
	opts := rockscache.NewDefaultOptions()
	opts.LockExpire = rocksCacheTimeout
	opts.WaitReplicasTimeout = rocksCacheTimeout
	opts.StrongConsistency = true
	opts.RandomExpireAdjustment = 0.2

	return &opts
}

func newRocksCacheClient(rdb redis.UniversalClient) *rocksCacheClient {
	if rdb == nil {
		return &rocksCacheClient{}
	}
	rc := &rocksCacheClient{
		rdb:    rdb,
		client: rockscache.NewClient(rdb, *GetRocksCacheOptions()),
	}
	return rc
}

type rocksCacheClient struct {
	rdb    redis.UniversalClient
	client *rockscache.Client
}

func (x *rocksCacheClient) GetClient() *rockscache.Client {
	return x.client
}

func (x *rocksCacheClient) Disable() bool {
	return x.client == nil
}

func (x *rocksCacheClient) GetRedis() redis.UniversalClient {
	return x.rdb
}

func (x *rocksCacheClient) GetBatchDeleter(topics ...string) cache.BatchDeleter {
	return NewBatchDeleterRedis(x, topics)
}

func batchGetCache2[K comparable, V any](ctx context.Context, rcClient *rocksCacheClient, expire time.Duration, ids []K, idKey func(id K) string, vId func(v *V) K, fn func(ctx context.Context, ids []K) ([]*V, error)) ([]*V, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	if rcClient.Disable() {
		return fn(ctx, ids)
	}
	findKeys := make([]string, 0, len(ids))
	keyId := make(map[string]K)
	for _, id := range ids {
		key := idKey(id)
		if _, ok := keyId[key]; ok {
			continue
		}
		keyId[key] = id
		findKeys = append(findKeys, key)
	}
	slotKeys, err := groupKeysBySlot(ctx, rcClient.GetRedis(), findKeys)
	if err != nil {
		return nil, err
	}
	result := make([]*V, 0, len(findKeys))
	for _, keys := range slotKeys {
		indexCache, err := rcClient.GetClient().FetchBatch2(ctx, keys, expire, func(idx []int) (map[int]string, error) {
			queryIds := make([]K, 0, len(idx))
			idIndex := make(map[K]int)
			for _, index := range idx {
				id := keyId[keys[index]]
				idIndex[id] = index
				queryIds = append(queryIds, id)
			}
			values, err := fn(ctx, queryIds)
			if err != nil {
				log.ZError(ctx, "batchGetCache query database failed", err, "keys", keys, "queryIds", queryIds)
				return nil, err
			}
			if len(values) == 0 {
				return map[int]string{}, nil
			}
			cacheIndex := make(map[int]string)
			for _, value := range values {
				id := vId(value)
				index, ok := idIndex[id]
				if !ok {
					continue
				}
				bs, err := json.Marshal(value)
				if err != nil {
					log.ZError(ctx, "marshal failed", err)
					return nil, err
				}
				cacheIndex[index] = string(bs)
			}
			return cacheIndex, nil
		})
		if err != nil {
			return nil, errs.WrapMsg(err, "FetchBatch2 failed")
		}
		for index, data := range indexCache {
			if data == "" {
				continue
			}
			var value V
			if err := json.Unmarshal([]byte(data), &value); err != nil {
				return nil, errs.WrapMsg(err, "Unmarshal failed")
			}
			if cb, ok := any(&value).(BatchCacheCallback[K]); ok {
				cb.BatchCache(keyId[keys[index]])
			}
			result = append(result, &value)
		}
	}
	return result, nil
}

type BatchCacheCallback[K comparable] interface {
	BatchCache(id K)
}
