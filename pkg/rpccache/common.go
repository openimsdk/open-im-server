package rpccache

import (
	"github.com/openimsdk/tools/errs"
	"google.golang.org/protobuf/proto"
)

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

func respProtoMarshal(resp proto.Message, err error) ([]byte, error) {
	if err != nil {
		return nil, err
	}
	return proto.Marshal(resp)
}

func cacheUnmarshal[V any](resp []byte, err error) (*V, error) {
	if err != nil {
		return nil, err
	}
	var val V
	if err := proto.Unmarshal(resp, any(&val).(proto.Message)); err != nil {
		return nil, errs.WrapMsg(err, "local cache proto.Unmarshal error")
	}
	return &val, nil
}

type cacheProto[V any] struct{}

func (cacheProto[V]) Marshal(resp *V, err error) ([]byte, error) {
	if err != nil {
		return nil, err
	}
	return proto.Marshal(any(resp).(proto.Message))
}

func (cacheProto[V]) Unmarshal(resp []byte, err error) (*V, error) {
	if err != nil {
		return nil, err
	}
	var val V
	if err := proto.Unmarshal(resp, any(&val).(proto.Message)); err != nil {
		return nil, errs.WrapMsg(err, "local cache proto.Unmarshal error")
	}
	return &val, nil
}
