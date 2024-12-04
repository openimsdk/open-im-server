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

package cache

import (
	"context"

	"github.com/openimsdk/tools/s3"

	relationtb "github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
)

type ObjectCache interface {
	BatchDeleter
	CloneObjectCache() ObjectCache
	GetName(ctx context.Context, engine string, name string) (*relationtb.Object, error)
	DelObjectName(engine string, names ...string) ObjectCache
}

type S3Cache interface {
	BatchDeleter
	GetKey(ctx context.Context, engine string, key string) (*s3.ObjectInfo, error)
	DelS3Key(engine string, keys ...string) S3Cache
}

// TODO integrating minio.Cache and MinioCache interfaces.
type MinioCache interface {
	BatchDeleter
	GetImageObjectKeyInfo(ctx context.Context, key string, fn func(ctx context.Context) (*MinioImageInfo, error)) (*MinioImageInfo, error)
	GetThumbnailKey(ctx context.Context, key string, format string, width int, height int, minioCache func(ctx context.Context) (string, error)) (string, error)
	DelObjectImageInfoKey(keys ...string) MinioCache
	DelImageThumbnailKey(key string, format string, width int, height int) MinioCache
}

type MinioImageInfo struct {
	IsImg  bool   `json:"isImg"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Format string `json:"format"`
	Etag   string `json:"etag"`
}
