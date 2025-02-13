package mcache

import (
	"context"
	"strconv"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/cachekey"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/tools/errs"
	"github.com/redis/go-redis/v9"
)

func NewThirdCache(cache database.Cache) cache.ThirdCache {
	return &thirdCache{
		cache: cache,
	}
}

type thirdCache struct {
	cache database.Cache
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

func (c *thirdCache) get(ctx context.Context, key string) (string, error) {
	res, err := c.cache.Get(ctx, []string{key})
	if err != nil {
		return "", err
	}
	if val, ok := res[key]; ok {
		return val, nil
	}
	return "", errs.Wrap(redis.Nil)
}

func (c *thirdCache) SetFcmToken(ctx context.Context, account string, platformID int, fcmToken string, expireTime int64) (err error) {
	return errs.Wrap(c.cache.Set(ctx, c.getFcmAccountTokenKey(account, platformID), fcmToken, time.Duration(expireTime)*time.Second))
}

func (c *thirdCache) GetFcmToken(ctx context.Context, account string, platformID int) (string, error) {
	return c.get(ctx, c.getFcmAccountTokenKey(account, platformID))
}

func (c *thirdCache) DelFcmToken(ctx context.Context, account string, platformID int) error {
	return c.cache.Del(ctx, []string{c.getFcmAccountTokenKey(account, platformID)})
}

func (c *thirdCache) IncrUserBadgeUnreadCountSum(ctx context.Context, userID string) (int, error) {
	return c.cache.Incr(ctx, c.getUserBadgeUnreadCountSumKey(userID), 1)
}

func (c *thirdCache) SetUserBadgeUnreadCountSum(ctx context.Context, userID string, value int) error {
	return c.cache.Set(ctx, c.getUserBadgeUnreadCountSumKey(userID), strconv.Itoa(value), 0)
}

func (c *thirdCache) GetUserBadgeUnreadCountSum(ctx context.Context, userID string) (int, error) {
	str, err := c.get(ctx, c.getUserBadgeUnreadCountSumKey(userID))
	if err != nil {
		return 0, err
	}
	val, err := strconv.Atoi(str)
	if err != nil {
		return 0, errs.WrapMsg(err, "strconv.Atoi", "str", str)
	}
	return val, nil
}

func (c *thirdCache) SetGetuiToken(ctx context.Context, token string, expireTime int64) error {
	return c.cache.Set(ctx, c.getGetuiTokenKey(), token, time.Duration(expireTime)*time.Second)
}

func (c *thirdCache) GetGetuiToken(ctx context.Context) (string, error) {
	return c.get(ctx, c.getGetuiTokenKey())
}

func (c *thirdCache) SetGetuiTaskID(ctx context.Context, taskID string, expireTime int64) error {
	return c.cache.Set(ctx, c.getGetuiTaskIDKey(), taskID, time.Duration(expireTime)*time.Second)
}

func (c *thirdCache) GetGetuiTaskID(ctx context.Context) (string, error) {
	return c.get(ctx, c.getGetuiTaskIDKey())
}
