package localcache

import "github.com/hashicorp/golang-lru/v2/simplelru"

type EvictCallback[K comparable, V any] simplelru.EvictCallback[K, V]
