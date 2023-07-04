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

package splitter

type SplitResult struct {
	Item []string
}
type Splitter struct {
	splitCount int
	data       []string
}

func NewSplitter(splitCount int, data []string) *Splitter {
	return &Splitter{splitCount: splitCount, data: data}
}
func (s *Splitter) GetSplitResult() (result []*SplitResult) {
	remain := len(s.data) % s.splitCount
	integer := len(s.data) / s.splitCount
	for i := 0; i < integer; i++ {
		r := new(SplitResult)
		r.Item = s.data[i*s.splitCount : (i+1)*s.splitCount]
		result = append(result, r)
	}
	if remain > 0 {
		r := new(SplitResult)
		r.Item = s.data[integer*s.splitCount:]
		result = append(result, r)
	}
	return result
}
