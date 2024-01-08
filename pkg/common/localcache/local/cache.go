package local

import (
	"hash/fnv"
	"time"
	"unsafe"
)

type Cache[V any] interface {
	Get(key string, fetch func() (V, error)) (V, error)
	Del(key string) bool
}

func NewCache[V any](slotNum, slotSize int, successTTL, failedTTL time.Duration, target Target) Cache[V] {
	c := &cache[V]{
		n:      uint64(slotNum),
		slots:  make([]*LRU[string, V], slotNum),
		target: target,
	}
	for i := 0; i < slotNum; i++ {
		c.slots[i] = NewLRU[string, V](slotSize, successTTL, failedTTL, c.target)
	}
	return c
}

type cache[V any] struct {
	n      uint64
	slots  []*LRU[string, V]
	target Target
}

func (c *cache[V]) index(s string) uint64 {
	h := fnv.New64a()
	_, _ = h.Write(*(*[]byte)(unsafe.Pointer(&s)))
	//_, _ = h.Write([]byte(s))
	return h.Sum64() % c.n
}

func (c *cache[V]) Get(key string, fetch func() (V, error)) (V, error) {
	return c.slots[c.index(key)].Get(key, fetch)
}

func (c *cache[V]) Del(key string) bool {
	if c.slots[c.index(key)].Del(key) {
		c.target.IncrDelHit()
		return true
	} else {
		c.target.IncrDelNotFound()
		return false
	}
}
