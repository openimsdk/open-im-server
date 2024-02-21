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
	"bytes"
	"encoding/gob"

	"github.com/OpenIMSDK/tools/errs"
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
	buff := bytes.Buffer{}
	enc := gob.NewEncoder(&buff)
	err := enc.Encode(data)
	if err != nil {
		return nil, err
	}
	return buff.Bytes(), nil
}

func (g *GobEncoder) Decode(encodeData []byte, decodeData any) error {
	buff := bytes.NewBuffer(encodeData)
	dec := gob.NewDecoder(buff)
	err := dec.Decode(decodeData)
	if err != nil {
		return errs.Wrap(err)
	}
	return nil
}
