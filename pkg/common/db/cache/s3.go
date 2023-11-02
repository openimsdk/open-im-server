package cache

import (
	"context"
	"github.com/OpenIMSDK/tools/errs"
	"github.com/dtm-labs/rockscache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/s3"
	relationtb "github.com/openimsdk/open-im-server/v3/pkg/common/db/table/relation"
	"github.com/redis/go-redis/v9"
	"image"
	"time"
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
	return "OBJECT_INFO:" + name
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
	GetImageObjectKeyInfo(ctx context.Context, key string, fn func(ctx context.Context, key string) (*MinioImageInfo, image.Image, error)) (*MinioImageInfo, image.Image, error)
	GetThumbnailKey(ctx context.Context, key string, format string, width int, height int, minioCache func(ctx context.Context) (string, error)) (string, error)
	//DelS3Key(engine string, keys ...string) S3Cache
}

func NewMinioCache(rdb redis.UniversalClient, s3 s3.Interface) MinioCache {
	rcClient := rockscache.NewClient(rdb, rockscache.NewDefaultOptions())
	return &minioCacheRedis{
		rcClient:   rcClient,
		expireTime: time.Hour * 12,
		s3:         s3,
		metaCache:  NewMetaCacheRedis(rcClient),
	}
}

type minioCacheRedis struct {
	metaCache
	s3         s3.Interface
	rcClient   *rockscache.Client
	expireTime time.Duration
}

func (g *minioCacheRedis) NewCache() MinioCache {
	return &minioCacheRedis{
		rcClient:   g.rcClient,
		expireTime: g.expireTime,
		s3:         g.s3,
		metaCache:  NewMetaCacheRedis(g.rcClient, g.metaCache.GetPreDelKeys()...),
	}
}

//func (g *minioCacheRedis) DelS3Key(engine string, keys ...string) MinioCache {
//	s3cache := g.NewCache()
//	ks := make([]string, 0, len(keys))
//	for _, key := range keys {
//		ks = append(ks, g.getS3Key(engine, key))
//	}
//	s3cache.AddKeys(ks...)
//	return s3cache
//}

func (g *minioCacheRedis) getMinioImageInfoKey(name string) string {
	return "MINIO:IMAGE:" + name
}

func (g *minioCacheRedis) getMinioImageThumbnailKey(name string) string {
	return "MINIO:THUMBNAIL:" + name
}

func (g *minioCacheRedis) GetImageObjectKeyInfo(ctx context.Context, key string, fn func(ctx context.Context, key string) (*MinioImageInfo, image.Image, error)) (*MinioImageInfo, image.Image, error) {
	var img image.Image
	info, err := getCache(ctx, g.rcClient, g.getMinioImageInfoKey(key), g.expireTime, func(ctx context.Context) (info *MinioImageInfo, err error) {
		info, img, err = fn(ctx, key)
		return
	})
	if err != nil {
		return nil, nil, err
	}
	if !info.IsImg {
		return nil, nil, errs.ErrData.Wrap("object not image")
	}
	return info, img, nil
}

func (g *minioCacheRedis) GetThumbnailKey(ctx context.Context, key string, format string, width int, height int, minioCache func(ctx context.Context, key string, format string, width int, height int) (string, error)) (string, error) {
	return getCache(ctx, g.rcClient, g.getMinioImageThumbnailKey(key), g.expireTime, func(ctx context.Context) (string, error) {
		info, img, err := g.GetImageObjectKeyInfo(ctx, key, getInfo)
		if err != nil {
			return "", err
		}
		return minioCache(ctx, key, format, width, height, info, img)
	})
}

type MinioImageInfo struct {
	IsImg  bool   `json:"i,omitempty"`
	Width  int    `json:"w,omitempty"`
	Height int    `json:"h,omitempty"`
	Format string `json:"f,omitempty"`
	Etag   string `json:"e,omitempty"`
}
