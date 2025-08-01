package redis

import (
	"context"
	"strconv"
	"sync"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/cachekey"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/utils/datautil"
	"github.com/redis/go-redis/v9"
)

type tokenCache struct {
	rdb          redis.UniversalClient
	accessExpire time.Duration
}

func NewTokenCacheModel(rdb redis.UniversalClient, accessExpire int64) cache.TokenModel {
	c := &tokenCache{rdb: rdb}
	c.accessExpire = c.getExpireTime(accessExpire)
	return c
}

func (c *tokenCache) SetTokenFlag(ctx context.Context, userID string, platformID int, token string, flag int) error {
	return errs.Wrap(c.rdb.HSet(ctx, cachekey.GetTokenKey(userID, platformID), token, flag).Err())
}

// SetTokenFlagEx set token and flag with expire time
func (c *tokenCache) SetTokenFlagEx(ctx context.Context, userID string, platformID int, token string, flag int) error {
	key := cachekey.GetTokenKey(userID, platformID)
	if err := c.rdb.HSet(ctx, key, token, flag).Err(); err != nil {
		return errs.Wrap(err)
	}
	if err := c.rdb.Expire(ctx, key, c.accessExpire).Err(); err != nil {
		return errs.Wrap(err)
	}
	return nil
}

func (c *tokenCache) GetTokensWithoutError(ctx context.Context, userID string, platformID int) (map[string]int, error) {
	m, err := c.rdb.HGetAll(ctx, cachekey.GetTokenKey(userID, platformID)).Result()
	if err != nil {
		return nil, errs.Wrap(err)
	}
	mm := make(map[string]int)
	for k, v := range m {
		state, err := strconv.Atoi(v)
		if err != nil {
			return nil, errs.WrapMsg(err, "redis token value is not int", "value", v, "userID", userID, "platformID", platformID)
		}
		mm[k] = state
	}
	return mm, nil
}

func (c *tokenCache) HasTemporaryToken(ctx context.Context, userID string, platformID int, token string) error {
	err := c.rdb.Get(ctx, cachekey.GetTemporaryTokenKey(userID, platformID, token)).Err()
	if err != nil {
		return errs.Wrap(err)
	}
	return nil
}

func (c *tokenCache) GetAllTokensWithoutError(ctx context.Context, userID string) (map[int]map[string]int, error) {
	var (
		res     = make(map[int]map[string]int)
		resLock = sync.Mutex{}
	)

	keys := cachekey.GetAllPlatformTokenKey(userID)
	if err := ProcessKeysBySlot(ctx, c.rdb, keys, func(ctx context.Context, slot int64, keys []string) error {
		pipe := c.rdb.Pipeline()
		mapRes := make([]*redis.MapStringStringCmd, len(keys))
		for i, key := range keys {
			mapRes[i] = pipe.HGetAll(ctx, key)
		}
		_, err := pipe.Exec(ctx)
		if err != nil {
			return err
		}
		for i, m := range mapRes {
			mm := make(map[string]int)
			for k, v := range m.Val() {
				state, err := strconv.Atoi(v)
				if err != nil {
					return errs.WrapMsg(err, "redis token value is not int", "value", v, "userID", userID)
				}
				mm[k] = state
			}
			resLock.Lock()
			res[cachekey.GetPlatformIDByTokenKey(keys[i])] = mm
			resLock.Unlock()
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *tokenCache) SetTokenMapByUidPid(ctx context.Context, userID string, platformID int, m map[string]int) error {
	mm := make(map[string]any)
	for k, v := range m {
		mm[k] = v
	}
	return errs.Wrap(c.rdb.HSet(ctx, cachekey.GetTokenKey(userID, platformID), mm).Err())
}

func (c *tokenCache) BatchSetTokenMapByUidPid(ctx context.Context, tokens map[string]map[string]any) error {
	keys := datautil.Keys(tokens)
	if err := ProcessKeysBySlot(ctx, c.rdb, keys, func(ctx context.Context, slot int64, keys []string) error {
		pipe := c.rdb.Pipeline()
		for k, v := range tokens {
			pipe.HSet(ctx, k, v)
		}
		_, err := pipe.Exec(ctx)
		if err != nil {
			return errs.Wrap(err)
		}
		return nil
	}); err != nil {
		return err
	}
	return nil
}

func (c *tokenCache) DeleteTokenByUidPid(ctx context.Context, userID string, platformID int, fields []string) error {
	return errs.Wrap(c.rdb.HDel(ctx, cachekey.GetTokenKey(userID, platformID), fields...).Err())
}

func (c *tokenCache) getExpireTime(t int64) time.Duration {
	return time.Hour * 24 * time.Duration(t)
}

// DeleteTokenByTokenMap tokens key is platformID, value is token slice
func (c *tokenCache) DeleteTokenByTokenMap(ctx context.Context, userID string, tokens map[int][]string) error {
	var (
		keys   = make([]string, 0, len(tokens))
		keyMap = make(map[string][]string)
	)
	for k, v := range tokens {
		k1 := cachekey.GetTokenKey(userID, k)
		keys = append(keys, k1)
		keyMap[k1] = v
	}

	if err := ProcessKeysBySlot(ctx, c.rdb, keys, func(ctx context.Context, slot int64, keys []string) error {
		pipe := c.rdb.Pipeline()
		for k, v := range tokens {
			pipe.HDel(ctx, cachekey.GetTokenKey(userID, k), v...)
		}
		_, err := pipe.Exec(ctx)
		if err != nil {
			return errs.Wrap(err)
		}
		return nil
	}); err != nil {
		return err
	}

	return nil
}

func (c *tokenCache) DeleteAndSetTemporary(ctx context.Context, userID string, platformID int, fields []string) error {
	for _, f := range fields {
		k := cachekey.GetTemporaryTokenKey(userID, platformID, f)
		if err := c.rdb.Set(ctx, k, "", time.Minute*5).Err(); err != nil {
			return errs.Wrap(err)
		}
	}
	key := cachekey.GetTokenKey(userID, platformID)
	if err := c.rdb.HDel(ctx, key, fields...).Err(); err != nil {
		return errs.Wrap(err)
	}
	return nil
}
