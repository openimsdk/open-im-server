package redis

import (
	"context"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/cachekey"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/redis/go-redis/v9"
)

const (
	blackExpireTime = time.Second * 60 * 60 * 12
)

type BlackCacheRedis struct {
	cache.BatchDeleter
	expireTime time.Duration
	rcClient   *rocksCacheClient
	blackDB    database.Black
}

func NewBlackCacheRedis(rdb redis.UniversalClient, localCache *config.LocalCache, blackDB database.Black) cache.BlackCache {
	rc := newRocksCacheClient(rdb)
	return &BlackCacheRedis{
		BatchDeleter: rc.GetBatchDeleter(localCache.Friend.Topic),
		expireTime:   blackExpireTime,
		rcClient:     rc,
		blackDB:      blackDB,
	}
}

func (b *BlackCacheRedis) CloneBlackCache() cache.BlackCache {
	return &BlackCacheRedis{
		BatchDeleter: b.BatchDeleter.Clone(),
		expireTime:   b.expireTime,
		rcClient:     b.rcClient,
		blackDB:      b.blackDB,
	}
}

func (b *BlackCacheRedis) getBlackIDsKey(ownerUserID string) string {
	return cachekey.GetBlackIDsKey(ownerUserID)
}

func (b *BlackCacheRedis) GetBlackIDs(ctx context.Context, userID string) (blackIDs []string, err error) {
	return getCache(
		ctx,
		b.rcClient,
		b.getBlackIDsKey(userID),
		b.expireTime,
		func(ctx context.Context) ([]string, error) {
			return b.blackDB.FindBlackUserIDs(ctx, userID)
		},
	)
}

func (b *BlackCacheRedis) DelBlackIDs(_ context.Context, userID string) cache.BlackCache {
	cache := b.CloneBlackCache()
	cache.AddKeys(b.getBlackIDsKey(userID))

	return cache
}
