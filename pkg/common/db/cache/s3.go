package cache

import (
	"context"
	"strconv"
	"time"

	"github.com/dtm-labs/rockscache"
	"github.com/redis/go-redis/v9"

	"github.com/openimsdk/open-im-server/v3/pkg/common/db/s3"
	relationtb "github.com/openimsdk/open-im-server/v3/pkg/common/db/table/relation"
)

type ObjectCache interface {
	metaCache
	GetName(ctx context.Context, name string) (*relationtb.ObjectModel, error)
	DelObjectName(names ...string) ObjectCache
}

func NewObjectCacheRedis(rdb redis.UniversalClient, objDB relationtb.ObjectInfoModelInterface) ObjectCache {
	rcClient := rockscache.NewClient(rdb, rockscache.NewDefaultOptions())
	return &objectCacheRedis{
		rcClient:   rcClient,
		expireTime: time.Hour * 12,
		objDB:      objDB,
		metaCache:  NewMetaCacheRedis(rcClient),
	}
}

type objectCacheRedis struct {
	metaCache
	objDB      relationtb.ObjectInfoModelInterface
	rcClient   *rockscache.Client
	expireTime time.Duration
}

func (g *objectCacheRedis) NewCache() ObjectCache {
	return &objectCacheRedis{
		rcClient:   g.rcClient,
		expireTime: g.expireTime,
		objDB:      g.objDB,
		metaCache:  NewMetaCacheRedis(g.rcClient, g.metaCache.GetPreDelKeys()...),
	}
}

func (g *objectCacheRedis) DelObjectName(names ...string) ObjectCache {
	objectCache := g.NewCache()
	keys := make([]string, 0, len(names))
	for _, name := range names {
		keys = append(keys, g.getObjectKey(name))
	}
	objectCache.AddKeys(keys...)
	return objectCache
}

func (g *objectCacheRedis) getObjectKey(name string) string {
	return "OBJECT:" + name
}

func (g *objectCacheRedis) GetName(ctx context.Context, name string) (*relationtb.ObjectModel, error) {
	return getCache(ctx, g.rcClient, g.getObjectKey(name), g.expireTime, func(ctx context.Context) (*relationtb.ObjectModel, error) {
		return g.objDB.Take(ctx, name)
	})
}

type S3Cache interface {
	metaCache
	GetKey(ctx context.Context, engine string, key string) (*s3.ObjectInfo, error)
	DelS3Key(engine string, keys ...string) S3Cache
}

func NewS3Cache(rdb redis.UniversalClient, s3 s3.Interface) S3Cache {
	rcClient := rockscache.NewClient(rdb, rockscache.NewDefaultOptions())
	return &s3CacheRedis{
		rcClient:   rcClient,
		expireTime: time.Hour * 12,
		s3:         s3,
		metaCache:  NewMetaCacheRedis(rcClient),
	}
}

type s3CacheRedis struct {
	metaCache
	s3         s3.Interface
	rcClient   *rockscache.Client
	expireTime time.Duration
}

func (g *s3CacheRedis) NewCache() S3Cache {
	return &s3CacheRedis{
		rcClient:   g.rcClient,
		expireTime: g.expireTime,
		s3:         g.s3,
		metaCache:  NewMetaCacheRedis(g.rcClient, g.metaCache.GetPreDelKeys()...),
	}
}

func (g *s3CacheRedis) DelS3Key(engine string, keys ...string) S3Cache {
	s3cache := g.NewCache()
	ks := make([]string, 0, len(keys))
	for _, key := range keys {
		ks = append(ks, g.getS3Key(engine, key))
	}
	s3cache.AddKeys(ks...)
	return s3cache
}

func (g *s3CacheRedis) getS3Key(engine string, name string) string {
	return "S3:" + engine + ":" + name
}

func (g *s3CacheRedis) GetKey(ctx context.Context, engine string, name string) (*s3.ObjectInfo, error) {
	return getCache(ctx, g.rcClient, g.getS3Key(engine, name), g.expireTime, func(ctx context.Context) (*s3.ObjectInfo, error) {
		return g.s3.StatObject(ctx, name)
	})
}

type MinioCache interface {
	metaCache
	GetImageObjectKeyInfo(ctx context.Context, key string, fn func(ctx context.Context) (*MinioImageInfo, error)) (*MinioImageInfo, error)
	GetThumbnailKey(ctx context.Context, key string, format string, width int, height int, minioCache func(ctx context.Context) (string, error)) (string, error)
	DelObjectImageInfoKey(keys ...string) MinioCache
	DelImageThumbnailKey(key string, format string, width int, height int) MinioCache
}

func NewMinioCache(rdb redis.UniversalClient) MinioCache {
	rcClient := rockscache.NewClient(rdb, rockscache.NewDefaultOptions())
	return &minioCacheRedis{
		rcClient:   rcClient,
		expireTime: time.Hour * 24 * 7,
		metaCache:  NewMetaCacheRedis(rcClient),
	}
}

type minioCacheRedis struct {
	metaCache
	rcClient   *rockscache.Client
	expireTime time.Duration
}

func (g *minioCacheRedis) NewCache() MinioCache {
	return &minioCacheRedis{
		rcClient:   g.rcClient,
		expireTime: g.expireTime,
		metaCache:  NewMetaCacheRedis(g.rcClient, g.metaCache.GetPreDelKeys()...),
	}
}

func (g *minioCacheRedis) DelObjectImageInfoKey(keys ...string) MinioCache {
	s3cache := g.NewCache()
	ks := make([]string, 0, len(keys))
	for _, key := range keys {
		ks = append(ks, g.getObjectImageInfoKey(key))
	}
	s3cache.AddKeys(ks...)
	return s3cache
}

func (g *minioCacheRedis) DelImageThumbnailKey(key string, format string, width int, height int) MinioCache {
	s3cache := g.NewCache()
	s3cache.AddKeys(g.getMinioImageThumbnailKey(key, format, width, height))
	return s3cache
}

func (g *minioCacheRedis) getObjectImageInfoKey(key string) string {
	return "MINIO:IMAGE:" + key
}

func (g *minioCacheRedis) getMinioImageThumbnailKey(key string, format string, width int, height int) string {
	return "MINIO:THUMBNAIL:" + format + ":w" + strconv.Itoa(width) + ":h" + strconv.Itoa(height) + ":" + key
}

func (g *minioCacheRedis) GetImageObjectKeyInfo(ctx context.Context, key string, fn func(ctx context.Context) (*MinioImageInfo, error)) (*MinioImageInfo, error) {
	info, err := getCache(ctx, g.rcClient, g.getObjectImageInfoKey(key), g.expireTime, fn)
	if err != nil {
		return nil, err
	}
	return info, nil
}

func (g *minioCacheRedis) GetThumbnailKey(ctx context.Context, key string, format string, width int, height int, minioCache func(ctx context.Context) (string, error)) (string, error) {
	return getCache(ctx, g.rcClient, g.getMinioImageThumbnailKey(key, format, width, height), g.expireTime, minioCache)
}

type MinioImageInfo struct {
	IsImg  bool   `json:"isImg"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
	Format string `json:"format"`
	Etag   string `json:"etag"`
}
