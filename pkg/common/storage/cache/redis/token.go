// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package redis

import (
	"context"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/cachekey"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/utils/stringutil"
	"github.com/redis/go-redis/v9"
	"time"
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

func (c *tokenCache) getExpireTime(t int64) time.Duration {
	return time.Hour * 24 * time.Duration(t)
}
