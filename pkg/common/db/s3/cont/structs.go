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

import "github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/s3"

type InitiateUploadResult struct {
	UploadID string             `json:"uploadID"` // 上传ID
	PartSize int64              `json:"partSize"` // 分片大小
	Sign     *s3.AuthSignResult `json:"sign"`     // 分片信息
}

type UploadResult struct {
	Hash string `json:"hash"`
	Size int64  `json:"size"`
	Key  string `json:"key"`
}
