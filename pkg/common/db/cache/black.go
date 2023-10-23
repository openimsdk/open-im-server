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

package cache

import (
	"context"
	"time"

	"github.com/dtm-labs/rockscache"
	"github.com/redis/go-redis/v9"

	relationtb "github.com/openimsdk/open-im-server/v3/pkg/common/db/table/relation"
)

const (
	blackIDsKey     = "BLACK_IDS:"
	blackExpireTime = time.Second * 60 * 60 * 12
)

// args fn will exec when no data in msgCache.
type BlackCache interface {
	// get blackIDs from msgCache
	metaCache
	NewCache() BlackCache
	GetBlackIDs(ctx context.Context, userID string) (blackIDs []string, err error)
	// del user's blackIDs msgCache, exec when a user's black list changed
	DelBlackIDs(ctx context.Context, userID string) BlackCache
}

type BlackCacheRedis struct {
	metaCache
	expireTime time.Duration
	rcClient   *rockscache.Client
	blackDB    relationtb.BlackModelInterface
}

func NewBlackCacheRedis(
	rdb redis.UniversalClient,
	blackDB relationtb.BlackModelInterface,
	options rockscache.Options,
) BlackCache {
	rcClient := rockscache.NewClient(rdb, options)

	return &BlackCacheRedis{
		expireTime: blackExpireTime,
		rcClient:   rcClient,
		metaCache:  NewMetaCacheRedis(rcClient),
		blackDB:    blackDB,
	}
}

func (b *BlackCacheRedis) NewCache() BlackCache {
	return &BlackCacheRedis{expireTime: b.expireTime, rcClient: b.rcClient, blackDB: b.blackDB, metaCache: NewMetaCacheRedis(b.rcClient, b.metaCache.GetPreDelKeys()...)}
}

func (b *BlackCacheRedis) getBlackIDsKey(ownerUserID string) string {
	return blackIDsKey + ownerUserID
}

func (b *BlackCacheRedis) GetBlackIDs(ctx context.Context, userID string) (blackIDs []string, err error) {
	return getCache(ctx, b.rcClient, b.getBlackIDsKey(userID), b.expireTime, func(ctx context.Context) ([]string, error) {
		return b.blackDB.FindBlackUserIDs(ctx, userID)
	})
}

func (b *BlackCacheRedis) DelBlackIDs(ctx context.Context, userID string) BlackCache {
	cache := b.NewCache()
	cache.AddKeys(b.getBlackIDsKey(userID))

	return cache
}
