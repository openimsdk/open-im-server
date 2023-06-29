package cache

import (
	"context"
	"time"

	relationTb "github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	"github.com/dtm-labs/rockscache"
	"github.com/redis/go-redis/v9"
)

const (
	blackIDsKey     = "BLACK_IDS:"
	blackExpireTime = time.Second * 60 * 60 * 12
)

// args fn will exec when no data in msgCache
type BlackCache interface {
	//get blackIDs from msgCache
	metaCache
	NewCache() BlackCache
	GetBlackIDs(ctx context.Context, userID string) (blackIDs []string, err error)
	//del user's blackIDs msgCache, exec when a user's black list changed
	DelBlackIDs(ctx context.Context, userID string) BlackCache
}

type BlackCacheRedis struct {
	metaCache
	expireTime time.Duration
	rcClient   *rockscache.Client
	blackDB    relationTb.BlackModelInterface
}

func NewBlackCacheRedis(rdb redis.UniversalClient, blackDB relationTb.BlackModelInterface, options rockscache.Options) BlackCache {
	rcClient := rockscache.NewClient(rdb, options)
	return &BlackCacheRedis{
		expireTime: blackExpireTime,
		rcClient:   rcClient,
		metaCache:  NewMetaCacheRedis(rcClient),
		blackDB:    blackDB,
	}
}

func (b *BlackCacheRedis) NewCache() BlackCache {
	return &BlackCacheRedis{
		expireTime: b.expireTime,
		rcClient:   b.rcClient,
		blackDB:    b.blackDB,
		metaCache:  NewMetaCacheRedis(b.rcClient, b.metaCache.GetPreDelKeys()...),
	}
}

func (b *BlackCacheRedis) getBlackIDsKey(ownerUserID string) string {
	return blackIDsKey + ownerUserID
}

func (b *BlackCacheRedis) GetBlackIDs(ctx context.Context, userID string) (blackIDs []string, err error) {
	return getCache(ctx, b.rcClient, b.getBlackIDsKey(userID), b.expireTime, func(ctx context.Context) ([]string, error) {
		return b.blackDB.FindBlackUserIDs(ctx, userID)
	})
}

func (b *BlackCacheRedis) DelBlackIDs(ctx context.Context, userID string) BlackCache {
	cache := b.NewCache()
	cache.AddKeys(b.getBlackIDsKey(userID))
	return cache
}
