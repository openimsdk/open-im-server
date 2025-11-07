package localcache

func AnyValue[V any](v any, err error) (V, error) {
	if err != nil {
		var zero V
		return zero, err
	}
	return v.(V), nil
}
