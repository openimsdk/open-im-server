package local

import (
	"github.com/hashicorp/golang-lru/v2/simplelru"
	"sync"
	"time"
)

type waitItem[V any] struct {
	lock    sync.Mutex
	expires int64
	active  bool
	err     error
	value   V
}

func NewLRU[K comparable, V any](size int, successTTL, failedTTL time.Duration, target Target) *LRU[K, V] {
	core, err := simplelru.NewLRU[K, *waitItem[V]](size, nil)
	if err != nil {
		panic(err)
	}
	return &LRU[K, V]{
		core:       core,
		successTTL: successTTL,
		failedTTL:  failedTTL,
		target:     target,
	}
}

type LRU[K comparable, V any] struct {
	lock       sync.Mutex
	core       *simplelru.LRU[K, *waitItem[V]]
	successTTL time.Duration
	failedTTL  time.Duration
	target     Target
}

func (x *LRU[K, V]) Get(key K, fetch func() (V, error)) (V, error) {
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
		v = &waitItem[V]{}
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

func (x *LRU[K, V]) Del(key K) bool {
	x.lock.Lock()
	ok := x.core.Remove(key)
	x.lock.Unlock()
	return ok
}
