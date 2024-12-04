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
	"time"

	"github.com/dtm-labs/rockscache"
	"github.com/redis/go-redis/v9"

	"github.com/openimsdk/tools/log"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/cachekey"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
)

const (
	blackExpireTime = time.Second * 60 * 60 * 12
)

type BlackCacheRedis struct {
	cache.BatchDeleter
	expireTime time.Duration
	rcClient   *rockscache.Client
	blackDB    database.Black
}

func NewBlackCacheRedis(rdb redis.UniversalClient, localCache *config.LocalCache, blackDB database.Black, options *rockscache.Options) cache.BlackCache {
	batchHandler := NewBatchDeleterRedis(rdb, options, []string{localCache.Friend.Topic})
	b := localCache.Friend
	log.ZDebug(context.Background(), "black local cache init", "Topic", b.Topic, "SlotNum", b.SlotNum, "SlotSize", b.SlotSize, "enable", b.Enable())
	return &BlackCacheRedis{
		BatchDeleter: batchHandler,
		expireTime:   blackExpireTime,
		rcClient:     rockscache.NewClient(rdb, *options),
		blackDB:      blackDB,
	}
}

func (b *BlackCacheRedis) CloneBlackCache() cache.BlackCache {
	return &BlackCacheRedis{
		BatchDeleter: b.BatchDeleter.Clone(),
		expireTime:   b.expireTime,
		rcClient:     b.rcClient,
		blackDB:      b.blackDB,
	}
}

func (b *BlackCacheRedis) getBlackIDsKey(ownerUserID string) string {
	return cachekey.GetBlackIDsKey(ownerUserID)
}

func (b *BlackCacheRedis) GetBlackIDs(ctx context.Context, userID string) (blackIDs []string, err error) {
	return getCache(
		ctx,
		b.rcClient,
		b.getBlackIDsKey(userID),
		b.expireTime,
		func(ctx context.Context) ([]string, error) {
			return b.blackDB.FindBlackUserIDs(ctx, userID)
		},
	)
}

func (b *BlackCacheRedis) DelBlackIDs(_ context.Context, userID string) cache.BlackCache {
	cache := b.CloneBlackCache()
	cache.AddKeys(b.getBlackIDsKey(userID))

	return cache
}
