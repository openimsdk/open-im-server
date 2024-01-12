package localcache

import (
	"context"
	"time"
)

func defaultOption() *option {
	return &option{
		enable:          true,
		localSlotNum:    500,   //LRU的slot数量
		localSlotSize:   20000, //每个LRU的大小
		localSuccessTTL: time.Minute,
		localFailedTTL:  time.Second * 5,
		delFn:           make([]func(ctx context.Context, key ...string), 0, 2),
		target:          emptyTarget{},
	}
}

type option struct {
	enable          bool
	localSlotNum    int
	localSlotSize   int
	localSuccessTTL time.Duration
	localFailedTTL  time.Duration
	delFn           []func(ctx context.Context, key ...string)
	delCh           func(fn func(key ...string))
	target          Target
}

type Option func(o *option)

func WithDisable() Option {
	return func(o *option) {
		o.enable = false
	}
}

func WithLocalSlotNum(localSlotNum int) Option {
	if localSlotNum < 1 {
		panic("localSlotNum should be greater than 0")
	}
	return func(o *option) {
		o.localSlotNum = localSlotNum
	}
}

func WithLocalSlotSize(localSlotSize int) Option {
	if localSlotSize < 1 {
		panic("localSlotSize should be greater than 0")
	}
	return func(o *option) {
		o.localSlotSize = localSlotSize
	}
}

func WithLocalSuccessTTL(localSuccessTTL time.Duration) Option {
	if localSuccessTTL < 0 {
		panic("localSuccessTTL should be greater than 0")
	}
	return func(o *option) {
		o.localSuccessTTL = localSuccessTTL
	}
}

func WithLocalFailedTTL(localFailedTTL time.Duration) Option {
	if localFailedTTL < 0 {
		panic("localFailedTTL should be greater than 0")
	}
	return func(o *option) {
		o.localFailedTTL = localFailedTTL
	}
}

func WithTarget(target Target) Option {
	if target == nil {
		panic("target should not be nil")
	}
	return func(o *option) {
		o.target = target
	}
}

func WithDeleteKeyBefore(fn func(ctx context.Context, key ...string)) Option {
	if fn == nil {
		panic("fn should not be nil")
	}
	return func(o *option) {
		o.delFn = append(o.delFn, fn)
	}
}

func WithDeleteLocal(fn func(fn func(key ...string))) Option {
	if fn == nil {
		panic("fn should not be nil")
	}
	return func(o *option) {
		o.delCh = fn
	}
}
