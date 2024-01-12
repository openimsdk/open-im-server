package localcache

import (
	"context"
	"github.com/openimsdk/localcache/link"
	"github.com/openimsdk/localcache/local"
	lopt "github.com/openimsdk/localcache/option"
)

type Cache[V any] interface {
	Get(ctx context.Context, key string, fetch func(ctx context.Context) (V, error), opts ...*lopt.Option) (V, error)
	Del(ctx context.Context, key ...string)
}

func New[V any](opts ...Option) Cache[V] {
	opt := defaultOption()
	for _, o := range opts {
		o(opt)
	}
	c := cache[V]{opt: opt}
	if opt.localSlotNum > 0 && opt.localSlotSize > 0 {
		c.local = local.NewCache[V](opt.localSlotNum, opt.localSlotSize, opt.localSuccessTTL, opt.localFailedTTL, opt.target, c.onEvict)
		go func() {
			c.opt.delCh(c.del)
		}()
		if opt.linkSlotNum > 0 {
			c.link = link.New(opt.linkSlotNum)
		}
	}
	return &c
}

type cache[V any] struct {
	opt   *option
	link  link.Link
	local local.Cache[V]
}

func (c *cache[V]) onEvict(key string, value V) {
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
	for _, k := range key {
		lks := c.link.Del(k)
		c.local.Del(k)
		for k := range lks {
			c.local.Del(k)
		}
	}
}

func (c *cache[V]) Get(ctx context.Context, key string, fetch func(ctx context.Context) (V, error), opts ...*lopt.Option) (V, error) {
	if c.local != nil {
		return c.local.Get(key, func() (V, error) {
			if c.link != nil {
				for _, o := range opts {
					c.link.Link(key, o.Link...)
				}
			}
			return fetch(ctx)
		})
	} else {
		return fetch(ctx)
	}
}

func (c *cache[V]) Del(ctx context.Context, key ...string) {
	if len(key) == 0 {
		return
	}
	for _, fn := range c.opt.delFn {
		fn(ctx, key...)
	}
	if c.local != nil {
		c.del(key...)
	}
}
