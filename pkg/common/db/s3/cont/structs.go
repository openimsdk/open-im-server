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

import "github.com/openimsdk/open-im-server/v3/pkg/common/db/s3"

type InitiateUploadResult struct {
	// UploadID uniquely identifies the upload session for tracking and management purposes.
	UploadID string `json:"uploadID"`

	// PartSize specifies the size of each part in a multipart upload. This is relevant for breaking down large uploads into manageable pieces.
	PartSize int64 `json:"partSize"`

	// Sign contains the authentication and signature information necessary for securely uploading each part. This could include signed URLs or tokens.
	Sign *s3.AuthSignResult `json:"sign"`
}

type UploadResult struct {
	Hash string `json:"hash"`
	Size int64  `json:"size"`
	Key  string `json:"key"`
}
