package utilsv2

import "Open_IM/pkg/common/db/table"

//func DuplicateRemoval[T comparable](ts []T) []T {
//	v := make([]T, 0, len(ts))
//	tmp := map[T]struct{}{}
//	for _, t := range ts {
//		if _, ok := tmp[t]; !ok {
//			tmp[t] = struct{}{}
//			v = append(v, t)
//		}
//	}
//	return v
//}

func DuplicateRemovalAny[T any, V comparable](ts []T, fn func(t T) V) []T {
	v := make([]T, 0, len(ts))
	tmp := map[V]struct{}{}
	for _, t := range ts {
		k := fn(t)
		if _, ok := tmp[k]; !ok {
			tmp[k] = struct{}{}
			v = append(v, t)
		}
	}
	return v
}

func DuplicateRemoval[T comparable](ts []T) []T {
	return DuplicateRemovalAny(ts, func(t T) T {
		return t
	})
}

func demo() {

	groups := []*table.GroupModel{}

	groups = DuplicateRemovalAny(groups, func(t *table.GroupModel) string {
		return t.GroupID
	})

}

func DeleteAt[T any](ts []T, index ...int) []T {
	switch len(index) {
	case 0:
		return ts
	case 1:
		return append(ts[:index[0]], ts[index[0]+1:]...)
	default:
		tmp := make(map[int]struct{})
		for _, v := range index {
			tmp[v] = struct{}{}
		}
		for i, t := range ts {

		}

	}
	return nil
}
