package redis

import (
	"context"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/cachekey"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/tools/s3"
	"github.com/openimsdk/tools/s3/cont"
	"github.com/redis/go-redis/v9"
)

func NewObjectCacheRedis(rdb redis.UniversalClient, objDB database.ObjectInfo) cache.ObjectCache {
	rc := newRocksCacheClient(rdb)
	return &objectCacheRedis{
		BatchDeleter: rc.GetBatchDeleter(),
		rcClient:     rc,
		expireTime:   time.Hour * 12,
		objDB:        objDB,
	}
}

type objectCacheRedis struct {
	cache.BatchDeleter
	objDB      database.ObjectInfo
	rcClient   *rocksCacheClient
	expireTime time.Duration
}

func (g *objectCacheRedis) getObjectKey(engine string, name string) string {
	return cachekey.GetObjectKey(engine, name)
}

func (g *objectCacheRedis) CloneObjectCache() cache.ObjectCache {
	return &objectCacheRedis{
		BatchDeleter: g.BatchDeleter.Clone(),
		rcClient:     g.rcClient,
		expireTime:   g.expireTime,
		objDB:        g.objDB,
	}
}

func (g *objectCacheRedis) DelObjectName(engine string, names ...string) cache.ObjectCache {
	objectCache := g.CloneObjectCache()
	keys := make([]string, 0, len(names))
	for _, name := range names {
		keys = append(keys, g.getObjectKey(name, engine))
	}
	objectCache.AddKeys(keys...)
	return objectCache
}

func (g *objectCacheRedis) GetName(ctx context.Context, engine string, name string) (*model.Object, error) {
	return getCache(ctx, g.rcClient, g.getObjectKey(name, engine), g.expireTime, func(ctx context.Context) (*model.Object, error) {
		return g.objDB.Take(ctx, engine, name)
	})
}

func NewS3Cache(rdb redis.UniversalClient, s3 s3.Interface) cont.S3Cache {
	rc := newRocksCacheClient(rdb)
	return &s3CacheRedis{
		BatchDeleter: rc.GetBatchDeleter(),
		rcClient:     rc,
		expireTime:   time.Hour * 12,
		s3:           s3,
	}
}

type s3CacheRedis struct {
	cache.BatchDeleter
	s3         s3.Interface
	rcClient   *rocksCacheClient
	expireTime time.Duration
}

func (g *s3CacheRedis) getS3Key(engine string, name string) string {
	return cachekey.GetS3Key(engine, name)
}

func (g *s3CacheRedis) DelS3Key(ctx context.Context, engine string, keys ...string) error {
	ks := make([]string, 0, len(keys))
	for _, key := range keys {
		ks = append(ks, g.getS3Key(engine, key))
	}
	return g.BatchDeleter.ExecDelWithKeys(ctx, ks)
}

func (g *s3CacheRedis) GetKey(ctx context.Context, engine string, name string) (*s3.ObjectInfo, error) {
	return getCache(ctx, g.rcClient, g.getS3Key(engine, name), g.expireTime, func(ctx context.Context) (*s3.ObjectInfo, error) {
		return g.s3.StatObject(ctx, name)
	})
}
