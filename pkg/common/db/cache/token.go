package cache

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/pkg/common/cachekey"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/utils/stringutil"
	"github.com/redis/go-redis/v9"
)

func NewTokenCacheModel(rdb redis.UniversalClient) TokenModel {
	return &tokenCache{
		rdb: rdb,
	}
}

type TokenModel interface {
	AddTokenFlag(ctx context.Context, userID string, platformID int, token string, flag int) error
	GetTokensWithoutError(ctx context.Context, userID string, platformID int) (map[string]int, error)
	SetTokenMapByUidPid(ctx context.Context, userID string, platformID int, m map[string]int) error
	DeleteTokenByUidPid(ctx context.Context, userID string, platformID int, fields []string) error
}

type tokenCache struct {
	rdb redis.UniversalClient
}

func (c *tokenCache) AddTokenFlag(ctx context.Context, userID string, platformID int, token string, flag int) error {
	return errs.Wrap(c.rdb.HSet(ctx, cachekey.GetTokenKey(userID, platformID), token, flag).Err())
}

func (c *tokenCache) GetTokensWithoutError(ctx context.Context, userID string, platformID int) (map[string]int, error) {
	m, err := c.rdb.HGetAll(ctx, cachekey.GetTokenKey(userID, platformID)).Result()
	if err != nil {
		return nil, errs.Wrap(err)
	}
	mm := make(map[string]int)
	for k, v := range m {
		mm[k] = stringutil.StringToInt(v)
	}

	return mm, nil
}

func (c *tokenCache) SetTokenMapByUidPid(ctx context.Context, userID string, platformID int, m map[string]int) error {
	mm := make(map[string]any)
	for k, v := range m {
		mm[k] = v
	}
	return errs.Wrap(c.rdb.HSet(ctx, cachekey.GetTokenKey(userID, platformID), mm).Err())
}

func (c *tokenCache) DeleteTokenByUidPid(ctx context.Context, userID string, platformID int, fields []string) error {
	return errs.Wrap(c.rdb.HDel(ctx, cachekey.GetTokenKey(userID, platformID), fields...).Err())
}
