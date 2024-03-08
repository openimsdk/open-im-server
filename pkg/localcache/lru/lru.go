package lru

import "github.com/hashicorp/golang-lru/v2/simplelru"

type EvictCallback[K comparable, V any] simplelru.EvictCallback[K, V]

type LRU[K comparable, V any] interface {
	Get(key K, fetch func() (V, error)) (V, error)
	Del(key K) bool
	Stop()
}

type Target interface {
	IncrGetHit()
	IncrGetSuccess()
	IncrGetFailed()

	IncrDelHit()
	IncrDelNotFound()
}
