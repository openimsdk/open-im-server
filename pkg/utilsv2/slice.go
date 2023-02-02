package utilsv2

import (
	"sort"
)

// DistinctAny 切片去重
func DistinctAny[T any, V comparable](ts []T, fn func(t T) V) []T {
	v := make([]T, 0, len(ts))
	tmp := map[V]struct{}{}
	for i := 0; i < len(ts); i++ {
		t := ts[i]
		k := fn(t)
		if _, ok := tmp[k]; !ok {
			tmp[k] = struct{}{}
			v = append(v, t)
		}
	}
	return v
}

// Distinct 切片去重
func Distinct[T comparable](ts []T) []T {
	return DistinctAny(ts, func(t T) T {
		return t
	})
}

// DeleteAt 删除切片元素, 支持负数删除倒数第几个
func DeleteAt[T any](ts []T, index ...int) []T {
	switch len(index) {
	case 0:
		return ts
	case 1:
		i := index[0]
		if i < 0 {
			i = len(ts) + i
		}
		if len(ts) <= i {
			return ts
		}
		return append(ts[:i], ts[i+1:]...)
	default:
		tmp := make(map[int]struct{})
		for _, i := range index {
			if i < 0 {
				i = len(ts) + i
			}
			tmp[i] = struct{}{}
		}
		v := make([]T, 0, len(ts))
		for i := 0; i < len(ts); i++ {
			if _, ok := tmp[i]; !ok {
				v = append(v, ts[i])
			}
		}
		return v
	}
}

// IndexAny 获取元素所在的下标
func IndexAny[T any, V comparable](ts []T, t T, fn func(t T) V) int {
	k := fn(t)
	for i := 0; i < len(ts); i++ {
		if fn(ts[i]) == k {
			return i
		}
	}
	return -1
}

// IndexOf 可比较的元素index下标
func IndexOf[T comparable](ts []T, t T) int {
	return IndexAny(ts, t, func(t T) T {
		return t
	})
}

// IsContain 是否包含元素
func IsContain[T comparable](ts []T, t T) bool {
	return IndexOf(ts, t) >= 0
}

// SliceToMap 切片转map
func SliceToMap[T any, K comparable](ts []T, fn func(t T) K) map[K]T {
	kv := make(map[K]T)
	for i := 0; i < len(ts); i++ {
		t := ts[i]
		k := fn(t)
		kv[k] = t
	}
	return kv
}

// MapKey map获取所有key
func MapKey[K comparable, V any](kv map[K]V) []K {
	ks := make([]K, 0, len(kv))
	for k := range kv {
		ks = append(ks, k)
	}
	return ks
}

// MapValue map获取所有key
func MapValue[K comparable, V any](kv map[K]V) []V {
	vs := make([]V, 0, len(kv))
	for k := range kv {
		vs = append(vs, kv[k])
	}
	return vs
}

// Sort 排序
func Sort[T Ordered](ts []T, asc bool) []T {
	SortAny(ts, func(a, b T) bool {
		if asc {
			return a < b
		} else {
			return a > b
		}
	})
	return ts
}

// SortAny 排序
func SortAny[T any](ts []T, fn func(a, b T) bool) {
	sort.Sort(&sortSlice[T]{
		ts: ts,
		fn: fn,
	})
}

type sortSlice[T any] struct {
	ts []T
	fn func(a, b T) bool
}

func (o *sortSlice[T]) Len() int {
	return len(o.ts)
}

func (o *sortSlice[T]) Less(i, j int) bool {
	return o.fn(o.ts[i], o.ts[j])
}

func (o *sortSlice[T]) Swap(i, j int) {
	o.ts[i], o.ts[j] = o.ts[j], o.ts[i]
}

type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr | ~float32 | ~float64 | ~string
}
