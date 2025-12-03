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
	//TODO implement me
	panic("implement me")
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
