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

	"github.com/hashicorp/golang-lru/v2/expirable"
)

func NewExpirationLRU[K comparable, V any](size int, successTTL, failedTTL time.Duration, target Target, onEvict EvictCallback[K, V]) LRU[K, V] {
	var cb expirable.EvictCallback[K, *expirationLruItem[V]]
	if onEvict != nil {
		cb = func(key K, value *expirationLruItem[V]) {
			onEvict(key, value.value)
		}
	}
	core := expirable.NewLRU[K, *expirationLruItem[V]](size, cb, successTTL)
	return &ExpirationLRU[K, V]{
		core:       core,
		successTTL: successTTL,
		failedTTL:  failedTTL,
		target:     target,
	}
}

type expirationLruItem[V any] struct {
	lock  sync.RWMutex
	err   error
	value V
}

type ExpirationLRU[K comparable, V any] struct {
	lock       sync.Mutex
	core       *expirable.LRU[K, *expirationLruItem[V]]
	successTTL time.Duration
	failedTTL  time.Duration
	target     Target
}

func (x *ExpirationLRU[K, V]) GetBatch(keys []K, fetch func(keys []K) (map[K]V, error)) (map[K]V, error) {
	var (
		err     error
		results = make(map[K]V)
		misses  = make([]K, 0, len(keys))
	)

	for _, key := range keys {
		x.lock.Lock()
		v, ok := x.core.Get(key)
		x.lock.Unlock()
		if ok {
			x.target.IncrGetHit()
			v.lock.RLock()
			results[key] = v.value
			if v.err != nil && err == nil {
				err = v.err
			}
			v.lock.RUnlock()
			continue
		}
		misses = append(misses, key)
	}

	if len(misses) == 0 {
		return results, err
	}

	fetchValues, fetchErr := fetch(misses)
	if fetchErr != nil && err == nil {
		err = fetchErr
	}

	for key, val := range fetchValues {
		results[key] = val
		if fetchErr != nil {
			x.target.IncrGetFailed()
			continue
		}
		x.target.IncrGetSuccess()
		item := &expirationLruItem[V]{value: val}
		x.lock.Lock()
		x.core.Add(key, item)
		x.lock.Unlock()
	}

	// any keys not returned from fetch remain absent (no cache write)
	return results, err
}

func (x *ExpirationLRU[K, V]) Get(key K, fetch func() (V, error)) (V, error) {
	x.lock.Lock()
	v, ok := x.core.Get(key)
	if ok {
		x.lock.Unlock()
		x.target.IncrGetSuccess()
		v.lock.RLock()
		defer v.lock.RUnlock()
		return v.value, v.err
	} else {
		v = &expirationLruItem[V]{}
		x.core.Add(key, v)
		v.lock.Lock()
		x.lock.Unlock()
		defer v.lock.Unlock()
		v.value, v.err = fetch()
		if v.err == nil {
			x.target.IncrGetSuccess()
		} else {
			x.target.IncrGetFailed()
			x.core.Remove(key)
		}
		return v.value, v.err
	}
}

func (x *ExpirationLRU[K, V]) Del(key K) bool {
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

func (x *ExpirationLRU[K, V]) SetHas(key K, value V) bool {
	x.lock.Lock()
	defer x.lock.Unlock()
	if x.core.Contains(key) {
		x.core.Add(key, &expirationLruItem[V]{value: value})
		return true
	}
	return false
}

func (x *ExpirationLRU[K, V]) Set(key K, value V) {
	x.lock.Lock()
	defer x.lock.Unlock()
	x.core.Add(key, &expirationLruItem[V]{value: value})
}

func (x *ExpirationLRU[K, V]) Stop() {
}
