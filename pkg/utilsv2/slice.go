package utilsv2

func DuplicateRemovalAny[T any, V comparable](ts []T, fn func(t T) V) []T {
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

func DuplicateRemoval[T comparable](ts []T) []T {
	return DuplicateRemovalAny(ts, func(t T) T {
		return t
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
		v := make([]T, 0, len(ts))
		for i := 0; i < len(ts); i++ {
			if _, ok := tmp[i]; !ok {
				v = append(v, ts[i])
			}
		}
		return v
	}
}
