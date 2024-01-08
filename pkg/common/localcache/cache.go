package localcache

import (
	"context"
	"github.com/openimsdk/open-im-server/v3/pkg/common/localcache/local"
)

type Cache[V any] interface {
	Get(ctx context.Context, key string, fetch func(ctx context.Context) (V, error)) (V, error)
	Del(ctx context.Context, key ...string)
}

func New[V any](opts ...Option) Cache[V] {
	opt := defaultOption()
	for _, o := range opts {
		o(opt)
	}
	if opt.enable {
		lc := local.NewCache[V](opt.localSlotNum, opt.localSlotSize, opt.localSuccessTTL, opt.localFailedTTL, opt.target)
		c := &cache[V]{
			opt:   opt,
			local: lc,
		}
		go func() {
			c.opt.delCh(c.del)
		}()
		return c
	} else {
		return &cache[V]{
			opt: opt,
		}
	}
}

type cache[V any] struct {
	opt   *option
	local local.Cache[V]
}

func (c *cache[V]) del(key ...string) {
	for _, k := range key {
		c.local.Del(k)
	}
}

func (c *cache[V]) Get(ctx context.Context, key string, fetch func(ctx context.Context) (V, error)) (V, error) {
	if c.opt.enable {
		return c.local.Get(key, func() (V, error) {
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
	if c.opt.enable {
		c.del(key...)
	}
}
