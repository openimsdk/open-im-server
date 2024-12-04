package redis

import (
	"context"
	"encoding/json"
	"time"
	"unsafe"

	"github.com/dtm-labs/rockscache"
	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/singleflight"

	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
)

func getRocksCacheRedisClient(cli *rockscache.Client) redis.UniversalClient {
	type Client struct {
		rdb redis.UniversalClient
		_   rockscache.Options
		_   singleflight.Group
	}
	return (*Client)(unsafe.Pointer(cli)).rdb
}

func batchGetCache2[K comparable, V any](ctx context.Context, rcClient *rockscache.Client, expire time.Duration, ids []K, idKey func(id K) string, vId func(v *V) K, fn func(ctx context.Context, ids []K) ([]*V, error)) ([]*V, error) {
	if len(ids) == 0 {
		return nil, nil
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
	slotKeys, err := groupKeysBySlot(ctx, getRocksCacheRedisClient(rcClient), findKeys)
	if err != nil {
		return nil, err
	}
	result := make([]*V, 0, len(findKeys))
	for _, keys := range slotKeys {
		indexCache, err := rcClient.FetchBatch2(ctx, keys, expire, func(idx []int) (map[int]string, error) {
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
