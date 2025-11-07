package localcache

import (
	"context"
	"hash/fnv"
	"unsafe"

	"github.com/openimsdk/open-im-server/v3/pkg/localcache/link"
	"github.com/openimsdk/open-im-server/v3/pkg/localcache/lru"
)

type Cache[V any] interface {
	Get(ctx context.Context, key string, fetch func(ctx context.Context) (V, error)) (V, error)
	GetLink(ctx context.Context, key string, fetch func(ctx context.Context) (V, error), link ...string) (V, error)
	Del(ctx context.Context, key ...string)
	DelLocal(ctx context.Context, key ...string)
	Stop()
}

func LRUStringHash(key string) uint64 {
	h := fnv.New64a()
	h.Write(*(*[]byte)(unsafe.Pointer(&key)))
	return h.Sum64()
}

func New[V any](opts ...Option) Cache[V] {
	opt := defaultOption()
	for _, o := range opts {
		o(opt)
	}

	c := cache[V]{opt: opt}
	if opt.localSlotNum > 0 && opt.localSlotSize > 0 {
		createSimpleLRU := func() lru.LRU[string, V] {
			if opt.expirationEvict {
				return lru.NewExpirationLRU(opt.localSlotSize, opt.localSuccessTTL, opt.localFailedTTL, opt.target, c.onEvict)
			} else {
				return lru.NewLazyLRU(opt.localSlotSize, opt.localSuccessTTL, opt.localFailedTTL, opt.target, c.onEvict)
			}
		}
		if opt.localSlotNum == 1 {
			c.local = createSimpleLRU()
		} else {
			c.local = lru.NewSlotLRU(opt.localSlotNum, LRUStringHash, createSimpleLRU)
		}
		if opt.linkSlotNum > 0 {
			c.link = link.New(opt.linkSlotNum)
		}
	}
	return &c
}

type cache[V any] struct {
	opt   *option
	link  link.Link
	local lru.LRU[string, V]
}

func (c *cache[V]) onEvict(key string, value V) {
	_ = value

	if c.link != nil {
		lks := c.link.Del(key)
		for k := range lks {
			if key != k { // prevent deadlock
				c.local.Del(k)
			}
		}
	}
}

func (c *cache[V]) del(key ...string) {
	if c.local == nil {
		return
	}
	for _, k := range key {
		c.local.Del(k)
		if c.link != nil {
			lks := c.link.Del(k)
			for k := range lks {
				c.local.Del(k)
			}
		}
	}
}

func (c *cache[V]) Get(ctx context.Context, key string, fetch func(ctx context.Context) (V, error)) (V, error) {
	return c.GetLink(ctx, key, fetch)
}

func (c *cache[V]) GetLink(ctx context.Context, key string, fetch func(ctx context.Context) (V, error), link ...string) (V, error) {
	if c.local != nil {
		return c.local.Get(key, func() (V, error) {
			if len(link) > 0 {
				c.link.Link(key, link...)
			}
			return fetch(ctx)
		})
	} else {
		return fetch(ctx)
	}
}

func (c *cache[V]) Del(ctx context.Context, key ...string) {
	for _, fn := range c.opt.delFn {
		fn(ctx, key...)
	}
	c.del(key...)
}

func (c *cache[V]) DelLocal(ctx context.Context, key ...string) {
	c.del(key...)
}

func (c *cache[V]) Stop() {
	c.local.Stop()
}
