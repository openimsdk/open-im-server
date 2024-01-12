package localcache

import (
	"fmt"
	"sync/atomic"
)

type Target interface {
	IncrGetHit()
	IncrGetSuccess()
	IncrGetFailed()

	IncrDelHit()
	IncrDelNotFound()
}

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

type emptyTarget struct{}

func (e emptyTarget) IncrGetHit() {}

func (e emptyTarget) IncrGetSuccess() {}

func (e emptyTarget) IncrGetFailed() {}

func (e emptyTarget) IncrDelHit() {}

func (e emptyTarget) IncrDelNotFound() {}
