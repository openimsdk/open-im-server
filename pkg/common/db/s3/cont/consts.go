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

const (
	// HashPath defines the storage path for hash data within the 'openim' directory.
	hashPath = "openim/data/hash/"

	// TempPath specifies the directory for temporary files in the 'openim' structure.
	tempPath = "openim/temp/"

	// DirectPath indicates the directory for direct uploads or access within the 'openim' structure.
	DirectPath = "openim/direct"

	// UploadTypeMultipart represents the identifier for multipart uploads,
	// allowing large files to be uploaded in chunks.
	UploadTypeMultipart = 1

	// UploadTypePresigned signifies the use of presigned URLs for uploads,
	// facilitating secure, authorized file transfers without requiring direct access to the storage credentials.
	UploadTypePresigned = 2

	// PartSeparator is used as a delimiter in multipart upload processes,
	// separating individual file parts.
	partSeparator = ","
)
