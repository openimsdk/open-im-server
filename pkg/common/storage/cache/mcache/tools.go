package mcache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/tools/log"
)

func getCache[V any](ctx context.Context, cache database.Cache, key string, expireTime time.Duration, fn func(ctx context.Context) (V, error)) (V, error) {
	getDB := func() (V, bool, error) {
		res, err := cache.Get(ctx, []string{key})
		if err != nil {
			var val V
			return val, false, err
		}
		var val V
		if str, ok := res[key]; ok {
			if json.Unmarshal([]byte(str), &val) != nil {
				return val, false, err
			}
			return val, true, nil
		}
		return val, false, nil
	}
	dbVal, ok, err := getDB()
	if err != nil {
		return dbVal, err
	}
	if ok {
		return dbVal, nil
	}
	lockValue, err := cache.Lock(ctx, key, time.Minute)
	if err != nil {
		return dbVal, err
	}
	defer func() {
		if err := cache.Unlock(ctx, key, lockValue); err != nil {
			log.ZError(ctx, "unlock cache key", err, "key", key, "value", lockValue)
		}
	}()
	dbVal, ok, err = getDB()
	if err != nil {
		return dbVal, err
	}
	if ok {
		return dbVal, nil
	}
	val, err := fn(ctx)
	if err != nil {
		return val, err
	}
	data, err := json.Marshal(val)
	if err != nil {
		return val, err
	}
	if err := cache.Set(ctx, key, string(data), expireTime); err != nil {
		return val, err
	}
	return val, nil
}
