package lru

import (
	"fmt"
	"hash/fnv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
	"unsafe"
)

type cacheTarget struct {
	getHit      int64
	getSuccess  int64
	getFailed   int64
	delHit      int64
	delNotFound int64
}

func (r *cacheTarget) IncrGetHit() {
	atomic.AddInt64(&r.getHit, 1)
}

func (r *cacheTarget) IncrGetSuccess() {
	atomic.AddInt64(&r.getSuccess, 1)
}

func (r *cacheTarget) IncrGetFailed() {
	atomic.AddInt64(&r.getFailed, 1)
}

func (r *cacheTarget) IncrDelHit() {
	atomic.AddInt64(&r.delHit, 1)
}

func (r *cacheTarget) IncrDelNotFound() {
	atomic.AddInt64(&r.delNotFound, 1)
}

func (r *cacheTarget) String() string {
	return fmt.Sprintf("getHit: %d, getSuccess: %d, getFailed: %d, delHit: %d, delNotFound: %d", r.getHit, r.getSuccess, r.getFailed, r.delHit, r.delNotFound)
}

func TestName(t *testing.T) {
	target := &cacheTarget{}
	l := NewSlotLRU[string, string](100, func(k string) uint64 {
		h := fnv.New64a()
		h.Write(*(*[]byte)(unsafe.Pointer(&k)))
		return h.Sum64()
	}, func() LRU[string, string] {
		return NewExpirationLRU[string, string](100, time.Second*60, time.Second, target, nil)
	})

	fn := func(key string, n int, fetch func() (string, error)) {
		for i := 0; i < n; i++ {
			v, err := l.Get(key, fetch)
			func(v ...any) {}(v, err)
		}
	}

	tmp := make(map[string]struct{})

	var wg sync.WaitGroup
	for i := 0; i < 10000; i++ {
		wg.Add(1)
		key := fmt.Sprintf("key_%d", i%200)
		tmp[key] = struct{}{}
		go func() {
			defer wg.Done()
			fn(key, 10000, func() (string, error) {
				return "value_" + key, nil
			})
		}()
	}
	wg.Wait()
	t.Log(len(tmp))
	t.Log(target.String())

}
