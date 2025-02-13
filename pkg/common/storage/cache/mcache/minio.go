package mcache

import (
	"context"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/cachekey"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/tools/s3/minio"
)

func NewMinioCache(cache database.Cache) minio.Cache {
	return &minioCache{
		cache:      cache,
		expireTime: time.Hour * 24 * 7,
	}
}

type minioCache struct {
	cache      database.Cache
	expireTime time.Duration
}

func (g *minioCache) getObjectImageInfoKey(key string) string {
	return cachekey.GetObjectImageInfoKey(key)
}

func (g *minioCache) getMinioImageThumbnailKey(key string, format string, width int, height int) string {
	return cachekey.GetMinioImageThumbnailKey(key, format, width, height)
}

func (g *minioCache) DelObjectImageInfoKey(ctx context.Context, keys ...string) error {
	ks := make([]string, 0, len(keys))
	for _, key := range keys {
		ks = append(ks, g.getObjectImageInfoKey(key))
	}
	return g.cache.Del(ctx, ks)
}

func (g *minioCache) DelImageThumbnailKey(ctx context.Context, key string, format string, width int, height int) error {
	return g.cache.Del(ctx, []string{g.getMinioImageThumbnailKey(key, format, width, height)})
}

func (g *minioCache) GetImageObjectKeyInfo(ctx context.Context, key string, fn func(ctx context.Context) (*minio.ImageInfo, error)) (*minio.ImageInfo, error) {
	return getCache[*minio.ImageInfo](ctx, g.cache, g.getObjectImageInfoKey(key), g.expireTime, fn)
}

func (g *minioCache) GetThumbnailKey(ctx context.Context, key string, format string, width int, height int, minioCache func(ctx context.Context) (string, error)) (string, error) {
	return getCache[string](ctx, g.cache, g.getMinioImageThumbnailKey(key, format, width, height), g.expireTime, minioCache)
}
