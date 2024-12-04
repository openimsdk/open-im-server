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

package redis

import (
	"context"
	"time"

	"github.com/dtm-labs/rockscache"
	"github.com/redis/go-redis/v9"

	"github.com/openimsdk/tools/s3"
	"github.com/openimsdk/tools/s3/cont"
	"github.com/openimsdk/tools/s3/minio"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/cachekey"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
)

func NewObjectCacheRedis(rdb redis.UniversalClient, objDB database.ObjectInfo) cache.ObjectCache {
	opts := rockscache.NewDefaultOptions()
	batchHandler := NewBatchDeleterRedis(rdb, &opts, nil)
	return &objectCacheRedis{
		BatchDeleter: batchHandler,
		rcClient:     rockscache.NewClient(rdb, opts),
		expireTime:   time.Hour * 12,
		objDB:        objDB,
	}
}

type objectCacheRedis struct {
	cache.BatchDeleter
	objDB      database.ObjectInfo
	rcClient   *rockscache.Client
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
	opts := rockscache.NewDefaultOptions()
	batchHandler := NewBatchDeleterRedis(rdb, &opts, nil)
	return &s3CacheRedis{
		BatchDeleter: batchHandler,
		rcClient:     rockscache.NewClient(rdb, opts),
		expireTime:   time.Hour * 12,
		s3:           s3,
	}
}

type s3CacheRedis struct {
	cache.BatchDeleter
	s3         s3.Interface
	rcClient   *rockscache.Client
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

func NewMinioCache(rdb redis.UniversalClient) minio.Cache {
	opts := rockscache.NewDefaultOptions()
	batchHandler := NewBatchDeleterRedis(rdb, &opts, nil)
	return &minioCacheRedis{
		BatchDeleter: batchHandler,
		rcClient:     rockscache.NewClient(rdb, opts),
		expireTime:   time.Hour * 24 * 7,
	}
}

type minioCacheRedis struct {
	cache.BatchDeleter
	rcClient   *rockscache.Client
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
