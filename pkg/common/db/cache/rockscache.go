package cache

import (
	"context"
	"encoding/json"
	"errors"

	"time"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"github.com/dtm-labs/rockscache"
)

const scanCount = 3000

var errIndex = errors.New("err index")

type metaCache interface {
	ExecDel(ctx context.Context) error
	// delete key rapid
	DelKey(ctx context.Context, key string) error
	AddKeys(keys ...string)
	ClearKeys()
	GetPreDelKeys() []string
}

func NewMetaCacheRedis(rcClient *rockscache.Client, keys ...string) metaCache {
	return &metaCacheRedis{rcClient: rcClient, keys: keys}
}

type metaCacheRedis struct {
	rcClient *rockscache.Client
	keys     []string
}

func (m *metaCacheRedis) ExecDel(ctx context.Context) error {
	if len(m.keys) > 0 {
		log.ZDebug(ctx, "DelKey", "keys", m.keys)
		return m.rcClient.TagAsDeletedBatch2(ctx, m.keys)
	}
	return nil
}

func (m *metaCacheRedis) DelKey(ctx context.Context, key string) error {
	return m.rcClient.TagAsDeleted2(ctx, key)
}

func (m *metaCacheRedis) AddKeys(keys ...string) {
	m.keys = append(m.keys, keys...)
}

func (m *metaCacheRedis) ClearKeys() {
	m.keys = []string{}
}

func (m *metaCacheRedis) GetPreDelKeys() []string {
	return m.keys
}

func GetDefaultOpt() rockscache.Options {
	opts := rockscache.NewDefaultOptions()
	opts.StrongConsistency = true
	opts.RandomExpireAdjustment = 0.2
	return opts
}

func getCache[T any](ctx context.Context, rcClient *rockscache.Client, key string, expire time.Duration, fn func(ctx context.Context) (T, error)) (T, error) {
	var t T
	var write bool
	v, err := rcClient.Fetch2(ctx, key, expire, func() (s string, err error) {
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

func batchGetCache[T any](ctx context.Context, rcClient *rockscache.Client, keys []string, expire time.Duration, keyIndexFn func(t T, keys []string) (int, error), fn func(ctx context.Context) ([]T, error)) ([]T, error) {
	batchMap, err := rcClient.FetchBatch2(ctx, keys, expire, func(idxs []int) (m map[int]string, err error) {
		values := make(map[int]string)
		tArrays, err := fn(ctx)
		if err != nil {
			return nil, err
		}
		for _, v := range tArrays {
			index, err := keyIndexFn(v, keys)
			if err != nil {
				continue
			}
			bs, err := json.Marshal(v)
			if err != nil {
				return nil, utils.Wrap(err, "marshal failed")
			}
			values[index] = string(bs)
		}
		return values, nil
	})
	if err != nil {
		return nil, err
	}
	var tArrays []T
	for _, v := range batchMap {
		if v != "" {
			var t T
			err = json.Unmarshal([]byte(v), &t)
			if err != nil {
				return nil, utils.Wrap(err, "unmarshal failed")
			}
			tArrays = append(tArrays, t)
		}
	}
	return tArrays, nil
}

func batchGetCacheMap[T any](ctx context.Context, rcClient *rockscache.Client, keys []string, originKeys []string, expire time.Duration, keyIndexFn func(t T, keys []string) (int, error), fn func(ctx context.Context) (map[string]T, error)) (map[string]T, error) {
	batchMap, err := rcClient.FetchBatch2(ctx, keys, expire, func(idxs []int) (m map[int]string, err error) {
		values := make(map[int]string)
		tArrays, err := fn(ctx)
		if err != nil {
			return nil, err
		}
		for _, v := range tArrays {
			index, err := keyIndexFn(v, keys)
			if err != nil {
				continue
			}
			bs, err := json.Marshal(v)
			if err != nil {
				return nil, utils.Wrap(err, "marshal failed")
			}
			values[index] = string(bs)
		}
		return values, nil
	})
	if err != nil {
		return nil, err
	}
	tMap := make(map[string]T)
	for i, v := range batchMap {
		if v != "" {
			var t T
			err = json.Unmarshal([]byte(v), &t)
			if err != nil {
				return nil, utils.Wrap(err, "unmarshal failed")
			}
			tMap[keys[i]] = t
		}
	}
	return tMap, nil
}
