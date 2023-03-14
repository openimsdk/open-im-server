package cache

import (
	"OpenIM/pkg/common/db/relation"
	"OpenIM/pkg/common/tracelog"
	"OpenIM/pkg/utils"
	"context"
	"github.com/dtm-labs/rockscache"
	"github.com/go-redis/redis/v8"
	"time"
)

const (
	blackIDsKey     = "BLACK_IDS:"
	blackExpireTime = time.Second * 60 * 60 * 12
)

// args fn will exec when no data in cache
type BlackCache interface {
	//get blackIDs from cache
	GetBlackIDs(ctx context.Context, userID string, fn func(ctx context.Context, userID string) ([]string, error)) (blackIDs []string, err error)
	//del user's blackIDs cache, exec when a user's black list changed
	DelBlackIDs(ctx context.Context, userID string) (err error)
}

type BlackCacheRedis struct {
	expireTime time.Duration
	rcClient   *rockscache.Client
	black      *relation.BlackGorm
}

func NewBlackCacheRedis(rdb redis.UniversalClient, blackDB BlackCache, options rockscache.Options) *BlackCacheRedis {
	return &BlackCacheRedis{
		expireTime: blackExpireTime,
		rcClient:   rockscache.NewClient(rdb, options),
	}
}

func (b *BlackCacheRedis) getBlackIDsKey(ownerUserID string) string {
	return blackIDsKey + ownerUserID
}

func (b *BlackCacheRedis) GetBlackIDs(ctx context.Context, userID string) (blackIDs []string, err error) {
	return GetCache(ctx, b.rcClient, b.getBlackIDsKey(userID), b.expireTime, func(ctx context.Context) ([]string, error) {
		return b.black.FindBlackUserIDs(ctx, userID)
	})
}

func (b *BlackCacheRedis) DelBlackIDs(ctx context.Context, userID string) (err error) {
	return b.rcClient.TagAsDeleted(b.getBlackIDsKey(userID))
}
