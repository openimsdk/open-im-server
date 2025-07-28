package redis

import (
	"context"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/cachekey"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/redis/go-redis/v9"
)

func NewClientConfigCache(rdb redis.UniversalClient, mgo database.ClientConfig) cache.ClientConfigCache {
	rc := newRocksCacheClient(rdb)
	return &ClientConfigCache{
		mgo:      mgo,
		rcClient: rc,
		delete:   rc.GetBatchDeleter(),
	}
}

type ClientConfigCache struct {
	mgo      database.ClientConfig
	rcClient *rocksCacheClient
	delete   cache.BatchDeleter
}

func (x *ClientConfigCache) getExpireTime(userID string) time.Duration {
	if userID == "" {
		return time.Hour * 24
	} else {
		return time.Hour
	}
}

func (x *ClientConfigCache) getClientConfigKey(userID string) string {
	return cachekey.GetClientConfigKey(userID)
}

func (x *ClientConfigCache) GetConfig(ctx context.Context, userID string) (map[string]string, error) {
	return getCache(ctx, x.rcClient, x.getClientConfigKey(userID), x.getExpireTime(userID), func(ctx context.Context) (map[string]string, error) {
		return x.mgo.Get(ctx, userID)
	})
}

func (x *ClientConfigCache) DeleteUserCache(ctx context.Context, userIDs []string) error {
	keys := make([]string, 0, len(userIDs))
	for _, userID := range userIDs {
		keys = append(keys, x.getClientConfigKey(userID))
	}
	return x.delete.ExecDelWithKeys(ctx, keys)
}

func (x *ClientConfigCache) GetUserConfig(ctx context.Context, userID string) (map[string]string, error) {
	config, err := x.GetConfig(ctx, "")
	if err != nil {
		return nil, err
	}
	if userID != "" {
		userConfig, err := x.GetConfig(ctx, userID)
		if err != nil {
			return nil, err
		}
		for k, v := range userConfig {
			config[k] = v
		}
	}
	return config, nil
}
