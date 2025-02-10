package redis

import (
	"context"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/cachekey"
	"github.com/openimsdk/tools/errs"
	"github.com/redis/go-redis/v9"
)

func NewThirdCache(rdb redis.UniversalClient) cache.ThirdCache {
	return &thirdCache{rdb: rdb}
}

type thirdCache struct {
	rdb redis.UniversalClient
}

func (c *thirdCache) getGetuiTokenKey() string {
	return cachekey.GetGetuiTokenKey()
}

func (c *thirdCache) getGetuiTaskIDKey() string {
	return cachekey.GetGetuiTaskIDKey()
}

func (c *thirdCache) getUserBadgeUnreadCountSumKey(userID string) string {
	return cachekey.GetUserBadgeUnreadCountSumKey(userID)
}

func (c *thirdCache) getFcmAccountTokenKey(account string, platformID int) string {
	return cachekey.GetFcmAccountTokenKey(account, platformID)
}

func (c *thirdCache) SetFcmToken(ctx context.Context, account string, platformID int, fcmToken string, expireTime int64) (err error) {
	return errs.Wrap(c.rdb.Set(ctx, c.getFcmAccountTokenKey(account, platformID), fcmToken, time.Duration(expireTime)*time.Second).Err())
}

func (c *thirdCache) GetFcmToken(ctx context.Context, account string, platformID int) (string, error) {
	val, err := c.rdb.Get(ctx, c.getFcmAccountTokenKey(account, platformID)).Result()
	if err != nil {
		return "", errs.Wrap(err)
	}
	return val, nil
}

func (c *thirdCache) DelFcmToken(ctx context.Context, account string, platformID int) error {
	return errs.Wrap(c.rdb.Del(ctx, c.getFcmAccountTokenKey(account, platformID)).Err())
}

func (c *thirdCache) IncrUserBadgeUnreadCountSum(ctx context.Context, userID string) (int, error) {
	seq, err := c.rdb.Incr(ctx, c.getUserBadgeUnreadCountSumKey(userID)).Result()

	return int(seq), errs.Wrap(err)
}

func (c *thirdCache) SetUserBadgeUnreadCountSum(ctx context.Context, userID string, value int) error {
	return errs.Wrap(c.rdb.Set(ctx, c.getUserBadgeUnreadCountSumKey(userID), value, 0).Err())
}

func (c *thirdCache) GetUserBadgeUnreadCountSum(ctx context.Context, userID string) (int, error) {
	val, err := c.rdb.Get(ctx, c.getUserBadgeUnreadCountSumKey(userID)).Int()
	return val, errs.Wrap(err)
}

func (c *thirdCache) SetGetuiToken(ctx context.Context, token string, expireTime int64) error {
	return errs.Wrap(c.rdb.Set(ctx, c.getGetuiTokenKey(), token, time.Duration(expireTime)*time.Second).Err())
}

func (c *thirdCache) GetGetuiToken(ctx context.Context) (string, error) {
	val, err := c.rdb.Get(ctx, c.getGetuiTokenKey()).Result()
	if err != nil {
		return "", errs.Wrap(err)
	}
	return val, nil
}

func (c *thirdCache) SetGetuiTaskID(ctx context.Context, taskID string, expireTime int64) error {
	return errs.Wrap(c.rdb.Set(ctx, c.getGetuiTaskIDKey(), taskID, time.Duration(expireTime)*time.Second).Err())
}

func (c *thirdCache) GetGetuiTaskID(ctx context.Context) (string, error) {
	val, err := c.rdb.Get(ctx, c.getGetuiTaskIDKey()).Result()
	if err != nil {
		return "", errs.Wrap(err)
	}
	return val, nil
}
