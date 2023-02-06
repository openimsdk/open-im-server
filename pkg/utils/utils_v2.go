package utils

import (
	"encoding/json"
	"sort"
)

// DistinctAny 去重
func DistinctAny[E any, K comparable](es []E, fn func(e E) K) []E {
	v := make([]E, 0, len(es))
	tmp := map[K]struct{}{}
	for i := 0; i < len(es); i++ {
		t := es[i]
		k := fn(t)
		if _, ok := tmp[k]; !ok {
			tmp[k] = struct{}{}
			v = append(v, t)
		}
	}
	return v
}

// Distinct 去重
func Distinct[T comparable](ts []T) []T {
	return DistinctAny(ts, func(t T) T {
		return t
	})
}

// Delete 删除切片元素, 支持负数删除倒数第几个
func Delete[E any](es []E, index ...int) []E {
	switch len(index) {
	case 0:
		return es
	case 1:
		i := index[0]
		if i < 0 {
			i = len(es) + i
		}
		if len(es) <= i {
			return es
		}
		return append(es[:i], es[i+1:]...)
	default:
		tmp := make(map[int]struct{})
		for _, i := range index {
			if i < 0 {
				i = len(es) + i
			}
			tmp[i] = struct{}{}
		}
		v := make([]E, 0, len(es))
		for i := 0; i < len(es); i++ {
			if _, ok := tmp[i]; !ok {
				v = append(v, es[i])
			}
		}
		return v
	}
}

// DeleteAt 删除切片元素, 支持负数删除倒数第几个
func DeleteAt[E any](es *[]E, index ...int) []E {
	v := Delete(*es, index...)
	*es = v
	return v
}

// IndexAny get the index of the element
func IndexAny[E any, K comparable](es []E, e E, fn func(e E) K) int {
	k := fn(e)
	for i := 0; i < len(es); i++ {
		if fn(es[i]) == k {
			return i
		}
	}
	return -1
}

// IndexOf get the index of the element
func IndexOf[E comparable](es []E, e E) int {
	return IndexAny(es, e, func(t E) E {
		return t
	})
}

// Contain 是否包含
func Contain[E comparable](es []E, e E) bool {
	return IndexOf(es, e) >= 0
}

// DuplicateAny 是否有重复的
func DuplicateAny[E any, K comparable](es []E, fn func(e E) K) bool {
	t := make(map[K]struct{})
	for _, e := range es {
		k := fn(e)
		if _, ok := t[k]; ok {
			return true
		}
		t[k] = struct{}{}
	}
	return false
}

// Duplicate 是否有重复的
func Duplicate[E comparable](es []E) bool {
	return DuplicateAny(es, func(e E) E {
		return e
	})
}

// SliceToMapOkAny slice to map (自定义类型, 筛选)
func SliceToMapOkAny[E any, K comparable, V any](es []E, fn func(e E) (K, V, bool)) map[K]V {
	kv := make(map[K]V)
	for i := 0; i < len(es); i++ {
		t := es[i]
		if k, v, ok := fn(t); ok {
			kv[k] = v
		}
	}
	return kv
}

// SliceToMapAny slice to map (自定义类型)
func SliceToMapAny[E any, K comparable, V any](es []E, fn func(e E) (K, V)) map[K]V {
	return SliceToMapOkAny(es, func(e E) (K, V, bool) {
		k, v := fn(e)
		return k, v, true
	})
}

// SliceToMap slice to map
func SliceToMap[E any, K comparable](es []E, fn func(e E) K) map[K]E {
	return SliceToMapOkAny[E, K, E](es, func(e E) (K, E, bool) {
		k := fn(e)
		return k, e, true
	})
}

// SliceSetAny slice to map[K]struct{}
func SliceSetAny[E any, K comparable](es []E, fn func(e E) K) map[K]struct{} {
	return SliceToMapAny(es, func(e E) (K, struct{}) {
		return fn(e), struct{}{}
	})
}

// Slice 批量转换切片类型
func Slice[E any, T any](es []E, fn func(e E) T) []T {
	v := make([]T, len(es))
	for i := 0; i < len(es); i++ {
		v = append(v, fn(es[i]))
	}
	return v
}

// SliceSet slice to map[E]struct{}
func SliceSet[E comparable](es []E) map[E]struct{} {
	return SliceSetAny(es, func(e E) E {
		return e
	})
}

// HasKey get whether the map contains key
func HasKey[K comparable, V any](m map[K]V, k K) bool {
	if m == nil {
		return false
	}
	_, ok := m[k]
	return ok
}

// Min get minimum value
func Min[E Ordered](e ...E) E {
	v := e[0]
	for _, t := range e[1:] {
		if v > t {
			v = t
		}
	}
	return v
}

// Max get maximum value
func Max[E Ordered](e ...E) E {
	v := e[0]
	for _, t := range e[1:] {
		if v < t {
			v = t
		}
	}
	return v
}

// BothExistAny 获取切片中共同存在的元素(交集)
func BothExistAny[E any, K comparable](es [][]E, fn func(e E) K) []E {
	if len(es) == 0 {
		return []E{}
	}
	var idx int
	ei := make([]map[K]E, len(es))
	for i := 0; i < len(ei); i++ {
		e := es[i]
		if len(e) == 0 {
			return []E{}
		}
		kv := make(map[K]E)
		for j := 0; j < len(e); j++ {
			t := e[j]
			k := fn(t)
			kv[k] = t
		}
		ei[i] = kv
		if len(kv) < len(ei[idx]) {
			idx = i
		}
	}
	v := make([]E, 0, len(ei[idx]))
	for k := range ei[idx] {
		all := true
		for i := 0; i < len(ei); i++ {
			if i == idx {
				continue
			}
			if _, ok := ei[i][k]; !ok {
				all = false
				break
			}
		}
		if !all {
			continue
		}
		v = append(v, ei[idx][k])
	}
	return v
}

// BothExist 获取切片中共同存在的元素(交集)
func BothExist[E comparable](es ...[]E) []E {
	return BothExistAny(es, func(e E) E {
		return e
	})
}

//// CompleteAny a中存在b的所有元素, 同时b中的所有元素a
//func CompleteAny[K comparable, E any](ks []K, es []E, fn func(e E) K) bool {
//	if len(ks) == 0 && len(es) == 0 {
//		return true
//	}
//	kn := make(map[K]uint8)
//	for _, e := range Distinct(ks) {
//		kn[e]++
//	}
//	for k := range SliceSetAny(es, fn) {
//		kn[k]++
//	}
//	for _, n := range kn {
//		if n != 2 {
//			return false
//		}
//	}
//	return true
//}

// Complete a和b去重后是否相等(忽略顺序)
func Complete[E comparable](a []E, b []E) bool {
	return len(Single(a, b)) == 0
}

// Keys get map keys
func Keys[K comparable, V any](kv map[K]V) []K {
	ks := make([]K, 0, len(kv))
	for k := range kv {
		ks = append(ks, k)
	}
	return ks
}

// Values get map values
func Values[K comparable, V any](kv map[K]V) []V {
	vs := make([]V, 0, len(kv))
	for k := range kv {
		vs = append(vs, kv[k])
	}
	return vs
}

// Sort basic type sorting
func Sort[E Ordered](es []E, asc bool) []E {
	SortAny(es, func(a, b E) bool {
		if asc {
			return a < b
		} else {
			return a > b
		}
	})
	return es
}

// SortAny custom sort method
func SortAny[E any](es []E, fn func(a, b E) bool) {
	sort.Sort(&sortSlice[E]{
		ts: es,
		fn: fn,
	})
}

// If true -> a, false -> b
func If[T any](isa bool, a, b T) T {
	if isa {
		return a
	}
	return b
}

// Equal 比较切片是否相对(包括元素顺序)
func Equal[E comparable](a []E, b []E) bool {
	if len(a) != len(b) {
		return false
	}
	for i := 0; i < len(a); i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// Single a中存在,b中不存在 或 b中存在,a中不存在
func Single[E comparable](a, b []E) []E {
	kn := make(map[E]uint8)
	for _, e := range Distinct(a) {
		kn[e]++
	}
	for _, e := range Distinct(b) {
		kn[e]++
	}
	v := make([]E, 0, len(kn))
	for k, n := range kn {
		if n == 1 {
			v = append(v, k)
		}
	}
	return v
}

func UniqueJoin(s ...string) string {
	data, _ := json.Marshal(s)
	return string(data)
}

type sortSlice[E any] struct {
	ts []E
	fn func(a, b E) bool
}

func (o *sortSlice[E]) Len() int {
	return len(o.ts)
}

func (o *sortSlice[E]) Less(i, j int) bool {
	return o.fn(o.ts[i], o.ts[j])
}

func (o *sortSlice[E]) Swap(i, j int) {
	o.ts[i], o.ts[j] = o.ts[j], o.ts[i]
}

// Ordered types that can be sorted
type Ordered interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr | ~float32 | ~float64 | ~string
}
