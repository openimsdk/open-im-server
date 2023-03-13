package cache

import (
	"OpenIM/pkg/utils"
	"context"
	"encoding/json"
	"github.com/dtm-labs/rockscache"
	"time"
)

const scanCount = 3000


func GetDefaultOpt() rockscache.Options {
	opts := rockscache.NewDefaultOptions()
	opts.StrongConsistency = true
	opts.RandomExpireAdjustment = 0.2
	return opts
}

func GetCache[T any](ctx context.Context, rcClient *rockscache.Client, key string, expire time.Duration, fn func(ctx context.Context) (T, error)) (T, error) {
	var t T
	var write bool
	v, err := rcClient.Fetch(key, expire, func() (s string, err error) {
		t, err = fn(ctx)
		if err != nil {
			return "", err
		}
		bs, err := json.Marshal(t)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		write = true
		return string(bs), nil
	})
	if err != nil {
		return t, err
	}
	if write {
		return t, nil
	}
	err = json.Unmarshal([]byte(v), &t)
	if err != nil {
		return t, utils.Wrap(err, "")
	}
	return t, nil
}

func GetCacheFor[E any, T any](ctx context.Context, list []E, fn func(ctx context.Context, item E) (T, error)) ([]T, error) {
	rs := make([]T, 0, len(list))
	for _, e := range list {
		r, err := fn(ctx, e)
		if err != nil {
			return nil, err
		}
		rs = append(rs, r)
	}
	return rs, nil
}
