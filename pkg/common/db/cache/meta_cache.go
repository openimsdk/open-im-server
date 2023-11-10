// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cache

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/OpenIMSDK/tools/mw/specialerror"

	"github.com/dtm-labs/rockscache"

	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/log"
	"github.com/OpenIMSDK/tools/utils"
)

const (
	scanCount     = 3000
	maxRetryTimes = 5
	retryInterval = time.Millisecond * 100
)

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
	return &metaCacheRedis{rcClient: rcClient, keys: keys, maxRetryTimes: maxRetryTimes, retryInterval: retryInterval}
}

type metaCacheRedis struct {
	rcClient      *rockscache.Client
	keys          []string
	maxRetryTimes int
	retryInterval time.Duration
}

func (m *metaCacheRedis) ExecDel(ctx context.Context) error {
	if len(m.keys) > 0 {
		log.ZDebug(ctx, "delete cache", "keys", m.keys)
		for _, key := range m.keys {
			for i := 0; i < m.maxRetryTimes; i++ {
				if err := m.rcClient.TagAsDeleted(key); err != nil {
					log.ZError(ctx, "delete cache failed", err, "key", key)
					time.Sleep(m.retryInterval)
					continue
				}
				break
			}

			//retryTimes := 0
			//for {
			//	m.rcClient.TagAsDeleted()
			//	if err := m.rcClient.TagAsDeletedBatch2(ctx, []string{key}); err != nil {
			//		if retryTimes >= m.maxRetryTimes {
			//			err = errs.ErrInternalServer.Wrap(
			//				fmt.Sprintf(
			//					"delete cache error: %v, keys: %v, retry times %d, please check redis server",
			//					err,
			//					key,
			//					retryTimes,
			//				),
			//			)
			//			log.ZWarn(ctx, "delete cache failed, please handle keys", err, "keys", key)
			//			return err
			//		}
			//		retryTimes++
			//	} else {
			//		break
			//	}
			//}
		}
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
	if v == "" {
		return t, errs.ErrRecordNotFound.Wrap("cache is not found")
	}
	err = json.Unmarshal([]byte(v), &t)
	if err != nil {
		log.ZError(ctx, "cache json.Unmarshal failed", err, "key", key, "value", v, "expire", expire)

		return t, utils.Wrap(err, "")
	}

	return t, nil
}

// func batchGetCache[T any](ctx context.Context, rcClient *rockscache.Client, keys []string, expire time.Duration, keyIndexFn func(t T, keys []string) (int, error), fn func(ctx context.Context) ([]T,
// error)) ([]T, error) {
//	batchMap, err := rcClient.FetchBatch2(ctx, keys, expire, func(idxs []int) (m map[int]string, err error) {
//		values := make(map[int]string)
//		tArrays, err := fn(ctx)
//		if err != nil {
//			return nil, err
//		}
//		for _, v := range tArrays {
//			index, err := keyIndexFn(v, keys)
//			if err != nil {
//				continue
//			}
//			bs, err := json.Marshal(v)
//			if err != nil {
//				return nil, utils.Wrap(err, "marshal failed")
//			}
//			values[index] = string(bs)
//		}
//		return values, nil
//	})
//	if err != nil {
//		return nil, err
//	}
//	var tArrays []T
//	for _, v := range batchMap {
//		if v != "" {
//			var t T
//			err = json.Unmarshal([]byte(v), &t)
//			if err != nil {
//				return nil, utils.Wrap(err, "unmarshal failed")
//			}
//			tArrays = append(tArrays, t)
//		}
//	}
//	return tArrays, nil
//}

func batchGetCache2[T any, K comparable](
	ctx context.Context,
	rcClient *rockscache.Client,
	expire time.Duration,
	keys []K,
	keyFn func(key K) string,
	fns func(ctx context.Context, key K) (T, error),
) ([]T, error) {
	if len(keys) == 0 {
		return nil, nil
	}
	res := make([]T, 0, len(keys))
	for _, key := range keys {
		val, err := getCache(ctx, rcClient, keyFn(key), expire, func(ctx context.Context) (T, error) {
			return fns(ctx, key)
		})
		if err != nil {
			if errs.ErrRecordNotFound.Is(specialerror.ErrCode(errs.Unwrap(err))) {
				continue
			}
			return nil, err
		}
		res = append(res, val)
	}

	return res, nil
}

//func batchGetCacheMap[T any](
//	ctx context.Context,
//	rcClient *rockscache.Client,
//	keys, originKeys []string,
//	expire time.Duration,
//	keyIndexFn func(s string, keys []string) (int, error),
//	fn func(ctx context.Context) (map[string]T, error),
//) (map[string]T, error) {
//	batchMap, err := rcClient.FetchBatch2(ctx, keys, expire, func(idxs []int) (m map[int]string, err error) {
//		tArrays, err := fn(ctx)
//		if err != nil {
//			return nil, err
//		}
//		values := make(map[int]string)
//		for k, v := range tArrays {
//			index, err := keyIndexFn(k, originKeys)
//			if err != nil {
//				continue
//			}
//			bs, err := json.Marshal(v)
//			if err != nil {
//				return nil, utils.Wrap(err, "marshal failed")
//			}
//			values[index] = string(bs)
//		}
//		return values, nil
//	})
//	if err != nil {
//		return nil, err
//	}
//	tMap := make(map[string]T)
//	for i, v := range batchMap {
//		if v != "" {
//			var t T
//			err = json.Unmarshal([]byte(v), &t)
//			if err != nil {
//				return nil, utils.Wrap(err, "unmarshal failed")
//			}
//			tMap[originKeys[i]] = t
//		}
//	}
//	return tMap, nil
//}
