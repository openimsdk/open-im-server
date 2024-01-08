package localcache

import (
	"context"
	"github.com/openimsdk/open-im-server/v3/pkg/common/localcache/link"
	"github.com/openimsdk/open-im-server/v3/pkg/common/localcache/local"
	opt "github.com/openimsdk/open-im-server/v3/pkg/common/localcache/option"
)

type Cache[V any] interface {
	Get(ctx context.Context, key string, fetch func(ctx context.Context) (V, error), opts ...*opt.Option) (V, error)
	Del(ctx context.Context, key ...string)
}

func New[V any](opts ...Option) Cache[V] {
	opt := defaultOption()
	for _, o := range opts {
		o(opt)
	}
	c := &cache[V]{opt: opt, link: link.New(opt.localSlotNum)}
	c.local = local.NewCache[V](opt.localSlotNum, opt.localSlotSize, opt.localSuccessTTL, opt.localFailedTTL, opt.target, c.onEvict)
	go func() {
		c.opt.delCh(c.del)
	}()
	return c
}

type cache[V any] struct {
	opt   *option
	link  link.Link
	local local.Cache[V]
}

func (c *cache[V]) onEvict(key string, value V) {
	for k := range c.link.Del(key) {
		if key != k {
			c.local.Del(k)
		}
	}
}

func (c *cache[V]) del(key ...string) {
	for _, k := range key {
		c.local.Del(k)
	}
}

func (c *cache[V]) Get(ctx context.Context, key string, fetch func(ctx context.Context) (V, error), opts ...*opt.Option) (V, error) {
	enable := c.opt.enable
	if len(opts) > 0 && opts[0].Enable != nil {
		enable = *opts[0].Enable
	}
	if enable {
		if len(opts) > 0 && len(opts[0].Link) > 0 {
			c.link.Link(key, opts[0].Link...)
		}
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
