package redis

import (
	"context"
	"github.com/dtm-labs/rockscache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/cachekey"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/tools/utils/datautil"
	"github.com/redis/go-redis/v9"
	"time"
)

func NewApplicationRedisCache(db database.Application, rdb redis.UniversalClient) *ApplicationRedisCache {
	return &ApplicationRedisCache{
		db:         db,
		rcClient:   rockscache.NewClient(rdb, *GetRocksCacheOptions()),
		deleter:    NewBatchDeleterRedis(rdb, GetRocksCacheOptions(), nil),
		expireTime: time.Hour * 24 * 7,
	}
}

type ApplicationRedisCache struct {
	db         database.Application
	rcClient   *rockscache.Client
	deleter    *BatchDeleterRedis
	expireTime time.Duration
}

func (a *ApplicationRedisCache) LatestVersion(ctx context.Context, platform string) (*model.Application, error) {
	return getCache(ctx, a.rcClient, cachekey.GetApplicationLatestVersionKey(platform), a.expireTime, func(ctx context.Context) (*model.Application, error) {
		return a.db.LatestVersion(ctx, platform)
	})
}

func (a *ApplicationRedisCache) DeleteCache(ctx context.Context, platforms []string) error {
	if len(platforms) == 0 {
		return nil
	}
	return a.deleter.ExecDelWithKeys(ctx, datautil.Slice(platforms, func(platform string) string {
		return cachekey.GetApplicationLatestVersionKey(platform)
	}))
}
