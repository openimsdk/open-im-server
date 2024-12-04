// Copyright © 2024 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package rpccache

import (
	"google.golang.org/protobuf/proto"

	"github.com/openimsdk/tools/errs"
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
