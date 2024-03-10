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

package link

import (
	"hash/fnv"
	"sync"
	"unsafe"
)

type Link interface {
	Link(key string, link ...string)
	Del(key string) map[string]struct{}
}

func newLinkKey() *linkKey {
	return &linkKey{
		data: make(map[string]map[string]struct{}),
	}
}

type linkKey struct {
	lock sync.Mutex
	data map[string]map[string]struct{}
}

func (x *linkKey) link(key string, link ...string) {
	x.lock.Lock()
	defer x.lock.Unlock()
	v, ok := x.data[key]
	if !ok {
		v = make(map[string]struct{})
		x.data[key] = v
	}
	for _, k := range link {
		v[k] = struct{}{}
	}
}

func (x *linkKey) del(key string) map[string]struct{} {
	x.lock.Lock()
	defer x.lock.Unlock()
	ks, ok := x.data[key]
	if !ok {
		return nil
	}
	delete(x.data, key)
	return ks
}

func New(n int) Link {
	if n <= 0 {
		panic("must be greater than 0")
	}
	slots := make([]*linkKey, n)
	for i := 0; i < len(slots); i++ {
		slots[i] = newLinkKey()
	}
	return &slot{
		n:     uint64(n),
		slots: slots,
	}
}

type slot struct {
	n     uint64
	slots []*linkKey
}

func (x *slot) index(s string) uint64 {
	h := fnv.New64a()
	_, _ = h.Write(*(*[]byte)(unsafe.Pointer(&s)))
	return h.Sum64() % x.n
}

func (x *slot) Link(key string, link ...string) {
	if len(link) == 0 {
		return
	}
	mk := key
	lks := make([]string, len(link))
	for i, k := range link {
		lks[i] = k
	}
	x.slots[x.index(mk)].link(mk, lks...)
	for _, lk := range lks {
		x.slots[x.index(lk)].link(lk, mk)
	}
}

func (x *slot) Del(key string) map[string]struct{} {
	return x.delKey(key)
}

func (x *slot) delKey(k string) map[string]struct{} {
	del := make(map[string]struct{})
	stack := []string{k}
	for len(stack) > 0 {
		curr := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		if _, ok := del[curr]; ok {
			continue
		}
		del[curr] = struct{}{}
		childKeys := x.slots[x.index(curr)].del(curr)
		for ck := range childKeys {
			stack = append(stack, ck)
		}
	}
	return del
}
