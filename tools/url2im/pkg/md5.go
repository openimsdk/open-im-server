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

package pkg

import (
	"crypto/md5"
	"encoding/hex"
	"hash"
	"io"
)

func NewMd5Reader(r io.Reader) *Md5Reader {
	return &Md5Reader{h: md5.New(), r: r}
}

type Md5Reader struct {
	h hash.Hash
	r io.Reader
}

func (r *Md5Reader) Read(p []byte) (n int, err error) {
	n, err = r.r.Read(p)
	if err == nil && n > 0 {
		r.h.Write(p[:n])
	}
	return
}

func (r *Md5Reader) Md5() string {
	return hex.EncodeToString(r.h.Sum(nil))
}
