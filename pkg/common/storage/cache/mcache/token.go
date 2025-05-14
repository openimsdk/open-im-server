package mcache

import (
	"context"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/cachekey"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
)

func NewTokenCacheModel(cache database.Cache, accessExpire int64) cache.TokenModel {
	c := &tokenCache{cache: cache}
	c.accessExpire = c.getExpireTime(accessExpire)
	return c
}

type tokenCache struct {
	cache        database.Cache
	accessExpire time.Duration
}

func (x *tokenCache) getTokenKey(userID string, platformID int, token string) string {
	return cachekey.GetTokenKey(userID, platformID) + ":" + token
}

func (x *tokenCache) SetTokenFlag(ctx context.Context, userID string, platformID int, token string, flag int) error {
	return x.cache.Set(ctx, x.getTokenKey(userID, platformID, token), strconv.Itoa(flag), x.accessExpire)
}

// SetTokenFlagEx set token and flag with expire time
func (x *tokenCache) SetTokenFlagEx(ctx context.Context, userID string, platformID int, token string, flag int) error {
	return x.SetTokenFlag(ctx, userID, platformID, token, flag)
}

func (x *tokenCache) GetTokensWithoutError(ctx context.Context, userID string, platformID int) (map[string]int, error) {
	prefix := x.getTokenKey(userID, platformID, "")
	m, err := x.cache.Prefix(ctx, prefix)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	mm := make(map[string]int)
	for k, v := range m {
		state, err := strconv.Atoi(v)
		if err != nil {
			log.ZError(ctx, "token value is not int", err, "value", v, "userID", userID, "platformID", platformID)
			continue
		}
		mm[strings.TrimPrefix(k, prefix)] = state
	}
	return mm, nil
}

func (x *tokenCache) HasTemporaryToken(ctx context.Context, userID string, platformID int, token string) error {
	key := cachekey.GetTemporaryTokenKey(userID, platformID, token)
	if _, err := x.cache.Get(ctx, []string{key}); err != nil {
		return err
	}
	return nil
}

func (x *tokenCache) GetAllTokensWithoutError(ctx context.Context, userID string) (map[int]map[string]int, error) {
	prefix := cachekey.UidPidToken + userID + ":"
	tokens, err := x.cache.Prefix(ctx, prefix)
	if err != nil {
		return nil, err
	}
	res := make(map[int]map[string]int)
	for key, flagStr := range tokens {
		flag, err := strconv.Atoi(flagStr)
		if err != nil {
			log.ZError(ctx, "token value is not int", err, "key", key, "value", flagStr, "userID", userID)
			continue
		}
		arr := strings.SplitN(strings.TrimPrefix(key, prefix), ":", 2)
		if len(arr) != 2 {
			log.ZError(ctx, "token value is not int", err, "key", key, "value", flagStr, "userID", userID)
			continue
		}
		platformID, err := strconv.Atoi(arr[0])
		if err != nil {
			log.ZError(ctx, "token value is not int", err, "key", key, "value", flagStr, "userID", userID)
			continue
		}
		token := arr[1]
		if token == "" {
			log.ZError(ctx, "token value is not int", err, "key", key, "value", flagStr, "userID", userID)
			continue
		}
		tk, ok := res[platformID]
		if !ok {
			tk = make(map[string]int)
			res[platformID] = tk
		}
		tk[token] = flag
	}
	return res, nil
}

func (x *tokenCache) SetTokenMapByUidPid(ctx context.Context, userID string, platformID int, m map[string]int) error {
	for token, flag := range m {
		err := x.SetTokenFlag(ctx, userID, platformID, token, flag)
		if err != nil {
			return err
		}
	}
	return nil
}

func (x *tokenCache) BatchSetTokenMapByUidPid(ctx context.Context, tokens map[string]map[string]any) error {
	for prefix, tokenFlag := range tokens {
		for token, flag := range tokenFlag {
			flagStr := fmt.Sprintf("%v", flag)
			if err := x.cache.Set(ctx, prefix+":"+token, flagStr, x.accessExpire); err != nil {
				return err
			}
		}
	}
	return nil
}

func (x *tokenCache) DeleteTokenByUidPid(ctx context.Context, userID string, platformID int, fields []string) error {
	keys := make([]string, 0, len(fields))
	for _, token := range fields {
		keys = append(keys, x.getTokenKey(userID, platformID, token))
	}
	return x.cache.Del(ctx, keys)
}

func (x *tokenCache) getExpireTime(t int64) time.Duration {
	return time.Hour * 24 * time.Duration(t)
}

func (x *tokenCache) DeleteTokenByTokenMap(ctx context.Context, userID string, tokens map[int][]string) error {
	keys := make([]string, 0, len(tokens))
	for platformID, ts := range tokens {
		for _, t := range ts {
			keys = append(keys, x.getTokenKey(userID, platformID, t))
		}
	}
	return x.cache.Del(ctx, keys)
}

func (x *tokenCache) DeleteAndSetTemporary(ctx context.Context, userID string, platformID int, fields []string) error {
	keys := make([]string, 0, len(fields))
	for _, f := range fields {
		keys = append(keys, x.getTokenKey(userID, platformID, f))
	}
	if err := x.cache.Del(ctx, keys); err != nil {
		return err
	}

	for _, f := range fields {
		k := cachekey.GetTemporaryTokenKey(userID, platformID, f)
		if err := x.cache.Set(ctx, k, "", time.Minute*5); err != nil {
			return errs.Wrap(err)
		}
	}

	return nil
}
