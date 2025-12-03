package cache

import (
	"context"
	relationtb "github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/tools/s3"
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
