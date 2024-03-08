package rpccache

func newListMap[V comparable](values []V, err error) (*listMap[V], error) {
	if err != nil {
		return nil, err
	}
	lm := &listMap[V]{
		List: values,
		Map:  make(map[V]struct{}, len(values)),
	}
	for _, value := range values {
		lm.Map[value] = struct{}{}
	}
	return lm, nil
}

type listMap[V comparable] struct {
	List []V
	Map  map[V]struct{}
}
