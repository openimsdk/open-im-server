package cache

import (
	"context"
	"github.com/openimsdk/tools/errs"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
)

type ThirdCache interface {
	SetFcmToken(ctx context.Context, account string, platformID int, fcmToken string, expireTime int64) (err error)
	GetFcmToken(ctx context.Context, account string, platformID int) (string, error)
	DelFcmToken(ctx context.Context, account string, platformID int) error
	IncrUserBadgeUnreadCountSum(ctx context.Context, userID string) (int, error)
	SetUserBadgeUnreadCountSum(ctx context.Context, userID string, value int) error
	GetUserBadgeUnreadCountSum(ctx context.Context, userID string) (int, error)
	SetGetuiToken(ctx context.Context, token string, expireTime int64) error
	GetGetuiToken(ctx context.Context) (string, error)
	SetGetuiTaskID(ctx context.Context, taskID string, expireTime int64) error
	GetGetuiTaskID(ctx context.Context) (string, error)
}

func NewThirdCache(rdb redis.UniversalClient) ThirdCache {
	return &thirdCache{rdb: rdb}
}

type thirdCache struct {
	rdb redis.UniversalClient
}

func (c *thirdCache) SetFcmToken(ctx context.Context, account string, platformID int, fcmToken string, expireTime int64) (err error) {
	return errs.Wrap(c.rdb.Set(ctx, FCM_TOKEN+account+":"+strconv.Itoa(platformID), fcmToken, time.Duration(expireTime)*time.Second).Err())
}

func (c *thirdCache) GetFcmToken(ctx context.Context, account string, platformID int) (string, error) {
	val, err := c.rdb.Get(ctx, FCM_TOKEN+account+":"+strconv.Itoa(platformID)).Result()
	if err != nil {
		return "", errs.Wrap(err)
	}
	return val, nil
}

func (c *thirdCache) DelFcmToken(ctx context.Context, account string, platformID int) error {
	return errs.Wrap(c.rdb.Del(ctx, FCM_TOKEN+account+":"+strconv.Itoa(platformID)).Err())
}

func (c *thirdCache) IncrUserBadgeUnreadCountSum(ctx context.Context, userID string) (int, error) {
	seq, err := c.rdb.Incr(ctx, userBadgeUnreadCountSum+userID).Result()

	return int(seq), errs.Wrap(err)
}

func (c *thirdCache) SetUserBadgeUnreadCountSum(ctx context.Context, userID string, value int) error {
	return errs.Wrap(c.rdb.Set(ctx, userBadgeUnreadCountSum+userID, value, 0).Err())
}

func (c *thirdCache) GetUserBadgeUnreadCountSum(ctx context.Context, userID string) (int, error) {
	val, err := c.rdb.Get(ctx, userBadgeUnreadCountSum+userID).Int()
	return val, errs.Wrap(err)
}

func (c *thirdCache) SetGetuiToken(ctx context.Context, token string, expireTime int64) error {
	return errs.Wrap(c.rdb.Set(ctx, getuiToken, token, time.Duration(expireTime)*time.Second).Err())
}

func (c *thirdCache) GetGetuiToken(ctx context.Context) (string, error) {
	val, err := c.rdb.Get(ctx, getuiToken).Result()
	if err != nil {
		return "", errs.Wrap(err)
	}
	return val, nil
}

func (c *thirdCache) SetGetuiTaskID(ctx context.Context, taskID string, expireTime int64) error {
	return errs.Wrap(c.rdb.Set(ctx, getuiTaskID, taskID, time.Duration(expireTime)*time.Second).Err())
}

func (c *thirdCache) GetGetuiTaskID(ctx context.Context) (string, error) {
	val, err := c.rdb.Get(ctx, getuiTaskID).Result()
	if err != nil {
		return "", errs.Wrap(err)
	}
	return val, nil
}
