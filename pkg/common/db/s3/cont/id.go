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

package cont

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
)

type multipartUploadID struct {
	Type int    `json:"a,omitempty"`
	ID   string `json:"b,omitempty"`
	Key  string `json:"c,omitempty"`
	Size int64  `json:"d,omitempty"`
	Hash string `json:"e,omitempty"`
}

func newMultipartUploadID(id multipartUploadID) string {
	data, err := json.Marshal(id)
	if err != nil {
		panic(err)
	}
	return base64.StdEncoding.EncodeToString(data)
}

func parseMultipartUploadID(id string) (*multipartUploadID, error) {
	data, err := base64.StdEncoding.DecodeString(id)
	if err != nil {
		return nil, fmt.Errorf("invalid multipart upload id: %w", err)
	}
	var upload multipartUploadID
	if err := json.Unmarshal(data, &upload); err != nil {
		return nil, fmt.Errorf("invalid multipart upload id: %w", err)
	}
	return &upload, nil
}
