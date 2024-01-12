package localcache

import (
	"context"
	"github.com/openimsdk/open-im-server/v3/pkg/common/localcache/link"
	opt "github.com/openimsdk/open-im-server/v3/pkg/common/localcache/option"
	"hash/fnv"
	"time"
	"unsafe"
)

const TimingWheelSize = 500

type Cache[V any] interface {
	Get(ctx context.Context, key string, fetch func(ctx context.Context) (V, error), opts ...*opt.Option) (V, error)
	Del(ctx context.Context, key ...string)
}

func New[V any](opts ...Option) Cache[V] {
	opt := defaultOption()
	for _, o := range opts {
		o(opt)
	}
	c := &cache[V]{
		opt:  opt,
		link: link.New(opt.localSlotNum),
		n:    uint64(opt.localSlotNum),
	}
	c.timingWheel = NewTimeWheel[string, V](TimingWheelSize, time.Second, c.exec)
	for i := 0; i < opt.localSlotNum; i++ {
		c.slots[i] = NewLRU[string, V](opt.localSlotSize, opt.localSuccessTTL, opt.localFailedTTL, opt.target, c.onEvict)
	}
	go func() {
		c.opt.delCh(c.del)
	}()
	return c
}

type cache[V any] struct {
	n           uint64
	slots       []*LRU[string, V]
	opt         *option
	link        link.Link
	timingWheel *TimeWheel[string, V]
}

func (c *cache[V]) index(key string) uint64 {
	h := fnv.New64a()
	_, _ = h.Write(*(*[]byte)(unsafe.Pointer(&key)))
	return h.Sum64() % c.n
}

func (c *cache[V]) onEvict(key string, value V) {
	lks := c.link.Del(key)
	for k := range lks {
		if key != k { // prevent deadlock
			c.slots[c.index(k)].Del(k)
		}
	}
}

func (c *cache[V]) del(key ...string) {
	for _, k := range key {
		lks := c.link.Del(k)
		c.slots[c.index(k)].Del(k)
		for k := range lks {
			c.slots[c.index(k)].Del(k)
		}
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
		return c.slots[c.index(key)].Get(key, func() (V, error) {
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
func (c *cache[V]) exec(key string, value V) {

}
