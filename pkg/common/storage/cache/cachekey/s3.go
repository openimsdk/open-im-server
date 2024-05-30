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

package cachekey

import "strconv"

const (
	object         = "OBJECT:"
	s3             = "S3:"
	minioImageInfo = "MINIO:IMAGE:"
	minioThumbnail = "MINIO:THUMBNAIL:"
)

func GetObjectKey(engine string, name string) string {
	return object + engine + ":" + name
}

func GetS3Key(engine string, name string) string {
	return s3 + engine + ":" + name
}

func GetObjectImageInfoKey(key string) string {
	return minioImageInfo + key
}

func GetMinioImageThumbnailKey(key string, format string, width int, height int) string {
	return minioThumbnail + format + ":w" + strconv.Itoa(width) + ":h" + strconv.Itoa(height) + ":" + key
}
