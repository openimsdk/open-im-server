package redis

import (
	"context"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/cachekey"
	"github.com/openimsdk/tools/s3/minio"
	"github.com/redis/go-redis/v9"
)

func NewMinioCache(rdb redis.UniversalClient) minio.Cache {
	rc := newRocksCacheClient(rdb)
	return &minioCacheRedis{
		BatchDeleter: rc.GetBatchDeleter(),
		rcClient:     rc,
		expireTime:   time.Hour * 24 * 7,
	}
}

type minioCacheRedis struct {
	cache.BatchDeleter
	rcClient   *rocksCacheClient
	expireTime time.Duration
}

func (g *minioCacheRedis) getObjectImageInfoKey(key string) string {
	return cachekey.GetObjectImageInfoKey(key)
}

func (g *minioCacheRedis) getMinioImageThumbnailKey(key string, format string, width int, height int) string {
	return cachekey.GetMinioImageThumbnailKey(key, format, width, height)
}

func (g *minioCacheRedis) DelObjectImageInfoKey(ctx context.Context, keys ...string) error {
	ks := make([]string, 0, len(keys))
	for _, key := range keys {
		ks = append(ks, g.getObjectImageInfoKey(key))
	}
	return g.BatchDeleter.ExecDelWithKeys(ctx, ks)
}

func (g *minioCacheRedis) DelImageThumbnailKey(ctx context.Context, key string, format string, width int, height int) error {
	return g.BatchDeleter.ExecDelWithKeys(ctx, []string{g.getMinioImageThumbnailKey(key, format, width, height)})

}

func (g *minioCacheRedis) GetImageObjectKeyInfo(ctx context.Context, key string, fn func(ctx context.Context) (*minio.ImageInfo, error)) (*minio.ImageInfo, error) {
	info, err := getCache(ctx, g.rcClient, g.getObjectImageInfoKey(key), g.expireTime, fn)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func (g *minioCacheRedis) GetThumbnailKey(ctx context.Context, key string, format string, width int, height int, minioCache func(ctx context.Context) (string, error)) (string, error) {
	return getCache(ctx, g.rcClient, g.getMinioImageThumbnailKey(key, format, width, height), g.expireTime, minioCache)
}
