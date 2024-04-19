// Copyright © 2023 OpenIM. All rights reserved.
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

package servererrs

import "github.com/openimsdk/tools/errs"

var Relation = &relation{m: make(map[int]map[int]struct{})}

func init() {
	Relation.Add(errs.RecordNotFoundError, UserIDNotFoundError)
	Relation.Add(errs.RecordNotFoundError, GroupIDNotFoundError)
	Relation.Add(errs.DuplicateKeyError, GroupIDExisted)
}

type relation struct {
	m map[int]map[int]struct{}
}

func (r *relation) Add(codes ...int) {
	if len(codes) < 2 {
		panic("codes length must be greater than 2")
	}
	for i := 1; i < len(codes); i++ {
		parent := codes[i-1]
		s, ok := r.m[parent]
		if !ok {
			s = make(map[int]struct{})
			r.m[parent] = s
		}
		for _, code := range codes[i:] {
			s[code] = struct{}{}
		}
	}
}

func (r *relation) Is(parent, child int) bool {
	if parent == child {
		return true
	}
	s, ok := r.m[parent]
	if !ok {
		return false
	}
	_, ok = s[child]
	return ok
}
