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

package msggateway

import (
	"encoding/json"

	"github.com/openimsdk/tools/errs"
)

type Encoder interface {
	Encode(data any) ([]byte, error)
	Decode(encodeData []byte, decodeData any) error
}

type GobEncoder struct{}

func NewGobEncoder() *GobEncoder {
	return &GobEncoder{}
}

func (g *GobEncoder) Encode(data any) ([]byte, error) {
	b, err := json.Marshal(data)
	if err != nil {
		return nil, errs.New("Encoder.Encode failed", "action", "encode")
	}
	return b, nil
}

func (g *GobEncoder) Decode(encodeData []byte, decodeData any) error {
	err := json.Unmarshal(encodeData, decodeData)
	if err != nil {
		return errs.New("Encoder.Decode failed", "action", "decode")
	}
	return nil
}
