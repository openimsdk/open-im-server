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
	//l := NewInertiaLRU[string, string](1000, time.Second*20, time.Second*5, target)

	fn := func(key string, n int, fetch func() (string, error)) {
		for i := 0; i < n; i++ {
			//v, err := l.Get(key, fetch)
			//if err == nil {
			//	t.Log("key", key, "value", v)
			//} else {
			//	t.Error("key", key, err)
			//}
			v, err := l.Get(key, fetch)
			//time.Sleep(time.Second / 100)
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
			//t.Log(key)
			fn(key, 10000, func() (string, error) {

				return "value_" + key, nil
			})
		}()

		//wg.Add(1)
		//go func() {
		//	defer wg.Done()
		//	for i := 0; i < 10; i++ {
		//		l.Del(key)
		//		time.Sleep(time.Second / 3)
		//	}
		//}()
	}
	wg.Wait()
	t.Log(len(tmp))
	t.Log(target.String())

}
