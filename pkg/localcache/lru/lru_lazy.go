// Copyright © 2024 OpenIM. All rights reserved.
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

package lru

import (
	"sync"
	"time"

	"github.com/hashicorp/golang-lru/v2/simplelru"
)

type layLruItem[V any] struct {
	lock    sync.Mutex
	expires int64
	err     error
	value   V
}

func NewLayLRU[K comparable, V any](size int, successTTL, failedTTL time.Duration, target Target, onEvict EvictCallback[K, V]) *LayLRU[K, V] {
	var cb simplelru.EvictCallback[K, *layLruItem[V]]
	if onEvict != nil {
		cb = func(key K, value *layLruItem[V]) {
			onEvict(key, value.value)
		}
	}
	core, err := simplelru.NewLRU[K, *layLruItem[V]](size, cb)
	if err != nil {
		panic(err)
	}
	return &LayLRU[K, V]{
		core:       core,
		successTTL: successTTL,
		failedTTL:  failedTTL,
		target:     target,
	}
}

type LayLRU[K comparable, V any] struct {
	lock       sync.Mutex
	core       *simplelru.LRU[K, *layLruItem[V]]
	successTTL time.Duration
	failedTTL  time.Duration
	target     Target
}

func (x *LayLRU[K, V]) Get(key K, fetch func() (V, error)) (V, error) {
	x.lock.Lock()
	v, ok := x.core.Get(key)
	if ok {
		x.lock.Unlock()
		v.lock.Lock()
		expires, value, err := v.expires, v.value, v.err
		if expires != 0 && expires > time.Now().UnixMilli() {
			v.lock.Unlock()
			x.target.IncrGetHit()
			return value, err
		}
	} else {
		v = &layLruItem[V]{}
		x.core.Add(key, v)
		v.lock.Lock()
		x.lock.Unlock()
	}
	defer v.lock.Unlock()
	if v.expires > time.Now().UnixMilli() {
		return v.value, v.err
	}
	v.value, v.err = fetch()
	if v.err == nil {
		v.expires = time.Now().Add(x.successTTL).UnixMilli()
		x.target.IncrGetSuccess()
	} else {
		v.expires = time.Now().Add(x.failedTTL).UnixMilli()
		x.target.IncrGetFailed()
	}
	return v.value, v.err
}

func (x *LayLRU[K, V]) GetBatch(keys []K, fetch func(keys []K) (map[K]V, error)) (map[K]V, error) {
	var (
		err  error
		once sync.Once
	)

	res := make(map[K]V)
	queries := make([]K, 0)
	setVs := make(map[K]*layLruItem[V])
	for _, key := range keys {
		x.lock.Lock()
		v, ok := x.core.Get(key)
		x.lock.Unlock()
		if ok {
			v.lock.Lock()
			expires, value, err1 := v.expires, v.value, v.err
			v.lock.Unlock()
			if expires != 0 && expires > time.Now().UnixMilli() {
				x.target.IncrGetHit()
				res[key] = value
				if err1 != nil {
					once.Do(func() {
						err = err1
					})
				}
				continue
			}
		}
		queries = append(queries, key)
	}
	values, err1 := fetch(queries)
	if err1 != nil {
		once.Do(func() {
			err = err1
		})
	}
	for key, val := range values {
		v := &layLruItem[V]{}
		v.value = val

		if err == nil {
			v.expires = time.Now().Add(x.successTTL).UnixMilli()
			x.target.IncrGetSuccess()
		} else {
			v.expires = time.Now().Add(x.failedTTL).UnixMilli()
			x.target.IncrGetFailed()
		}
		setVs[key] = v
		x.lock.Lock()
		x.core.Add(key, v)
		x.lock.Unlock()
		res[key] = val
	}

	return res, err
}

//func (x *LayLRU[K, V]) Has(key K) bool {
//	x.lock.Lock()
//	defer x.lock.Unlock()
//	return x.core.Contains(key)
//}

func (x *LayLRU[K, V]) Set(key K, value V) {
	x.lock.Lock()
	defer x.lock.Unlock()
	x.core.Add(key, &layLruItem[V]{value: value, expires: time.Now().Add(x.successTTL).UnixMilli()})
}

func (x *LayLRU[K, V]) SetHas(key K, value V) bool {
	x.lock.Lock()
	defer x.lock.Unlock()
	if x.core.Contains(key) {
		x.core.Add(key, &layLruItem[V]{value: value, expires: time.Now().Add(x.successTTL).UnixMilli()})
		return true
	}
	return false
}

func (x *LayLRU[K, V]) Del(key K) bool {
	x.lock.Lock()
	ok := x.core.Remove(key)
	x.lock.Unlock()
	if ok {
		x.target.IncrDelHit()
	} else {
		x.target.IncrDelNotFound()
	}
	return ok
}

func (x *LayLRU[K, V]) Stop() {

}
