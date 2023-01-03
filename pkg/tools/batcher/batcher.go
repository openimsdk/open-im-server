package batcher

import (
	"Open_IM/pkg/common/log"
	"context"
	"errors"
	"hash/crc32"
	"sync"
	"time"
)

var (
	ErrorNotSetFunction = errors.New("not set do function")
)

var (
	DefaultSize     = 100
	DefaultBuffer   = 100
	DefaultWorker   = 5
	DefaultInterval = time.Second
)

type DoFuntion func(ctx context.Context, val map[string][]interface{})
type Option func(c *Config)
type Config struct {
	size     int           //Number of message aggregations
	buffer   int           //The number of caches running in a single coroutine
	worker   int           //Number of coroutines processed in parallel
	interval time.Duration //Time of message aggregations
}

func newDefaultConfig() *Config {
	return &Config{
		size:     DefaultSize,
		buffer:   DefaultBuffer,
		worker:   DefaultWorker,
		interval: DefaultInterval,
	}
}

type Batcher struct {
	config   Config
	Do       func(ctx context.Context, val map[string][]interface{})
	Sharding func(key string) int
	chans    []chan *msg
	wait     sync.WaitGroup
}
type msg struct {
	key string
	val interface{}
}

func NewBatcher(fn DoFuntion, opts ...Option) *Batcher {
	b := &Batcher{}
	b.Do = fn
	config := newDefaultConfig()
	for _, o := range opts {
		o(config)
	}
	b.chans = make([]chan *msg, b.config.worker)
	for i := 0; i < b.config.worker; i++ {
		b.chans[i] = make(chan *msg, b.config.buffer)
	}
	return b
}
func WithSize(s int) Option {
	return func(c *Config) {
		c.size = s
	}
}

func WithBuffer(b int) Option {
	return func(c *Config) {
		c.buffer = b
	}
}

func WithWorker(w int) Option {
	return func(c *Config) {
		c.worker = w
	}
}

func WithInterval(i time.Duration) Option {
	return func(c *Config) {
		c.interval = i
	}
}
func (b *Batcher) Start() error {
	if b.Do == nil {
		return ErrorNotSetFunction
	}
	if b.Sharding == nil {
		b.Sharding = func(key string) int {
			hasCode := int(crc32.ChecksumIEEE([]byte(key)))
			return hasCode % b.config.worker
		}
	}
	b.wait.Add(len(b.chans))
	for i, ch := range b.chans {
		go b.merge(i, ch)
	}
	return nil
}

func (b *Batcher) Add(key string, val interface{}) error {
	ch, msg := b.add(key, val)
	select {
	case ch <- msg:
	default:
		return ErrFull
	}
	return nil
}

func (b *Batcher) add(key string, val interface{}) (chan *msg, *msg) {
	sharding := b.Sharding(key) % b.opts.worker
	ch := b.chans[sharding]
	msg := &msg{key: key, val: val}
	return ch, msg
}

func (b *Batcher) merge(idx int, ch <-chan *msg) {
	defer b.wait.Done()

	var (
		msg        *msg
		count      int
		closed     bool
		lastTicker = true
		interval   = b.opts.interval
		vals       = make(map[string][]interface{}, b.opts.size)
	)
	if idx > 0 {
		interval = time.Duration(int64(idx) * (int64(b.opts.interval) / int64(b.opts.worker)))
	}
	ticker := time.NewTicker(interval)
	for {
		select {
		case msg = <-ch:
			if msg == nil {
				closed = true
				break
			}
			count++
			vals[msg.key] = append(vals[msg.key], msg.val)
			if count >= b.opts.size {
				break
			}
			continue
		case <-ticker.C:
			if lastTicker {
				ticker.Stop()
				ticker = time.NewTicker(b.opts.interval)
				lastTicker = false
			}
		}
		if len(vals) > 0 {
			ctx := context.Background()
			b.Do(ctx, vals)
			vals = make(map[string][]interface{}, b.opts.size)
			count = 0
		}
		if closed {
			ticker.Stop()
			return
		}
	}
}

func (b *Batcher) Close() {
	for _, ch := range b.chans {
		ch <- nil
	}
	b.wait.Wait()
}
