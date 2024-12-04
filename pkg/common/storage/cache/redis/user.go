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
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
)

const (
	userExpireTime            = time.Second * 60 * 60 * 12
	userOlineStatusExpireTime = time.Second * 60 * 60 * 24
	statusMod                 = 501
)

type UserCacheRedis struct {
	cache.BatchDeleter
	rdb        redis.UniversalClient
	userDB     database.User
	expireTime time.Duration
	rcClient   *rockscache.Client
}

func NewUserCacheRedis(rdb redis.UniversalClient, localCache *config.LocalCache, userDB database.User, options *rockscache.Options) cache.UserCache {
	batchHandler := NewBatchDeleterRedis(rdb, options, []string{localCache.User.Topic})
	u := localCache.User
	log.ZDebug(context.Background(), "user local cache init", "Topic", u.Topic, "SlotNum", u.SlotNum, "SlotSize", u.SlotSize, "enable", u.Enable())
	return &UserCacheRedis{
		BatchDeleter: batchHandler,
		rdb:          rdb,
		userDB:       userDB,
		expireTime:   userExpireTime,
		rcClient:     rockscache.NewClient(rdb, *options),
	}
}

func (u *UserCacheRedis) getUserID(user *model.User) string {
	return user.UserID
}

func (u *UserCacheRedis) CloneUserCache() cache.UserCache {
	return &UserCacheRedis{
		BatchDeleter: u.BatchDeleter.Clone(),
		rdb:          u.rdb,
		userDB:       u.userDB,
		expireTime:   u.expireTime,
		rcClient:     u.rcClient,
	}
}

func (u *UserCacheRedis) getUserInfoKey(userID string) string {
	return cachekey.GetUserInfoKey(userID)
}

func (u *UserCacheRedis) getUserGlobalRecvMsgOptKey(userID string) string {
	return cachekey.GetUserGlobalRecvMsgOptKey(userID)
}

func (u *UserCacheRedis) GetUserInfo(ctx context.Context, userID string) (userInfo *model.User, err error) {
	return getCache(ctx, u.rcClient, u.getUserInfoKey(userID), u.expireTime, func(ctx context.Context) (*model.User, error) {
		return u.userDB.Take(ctx, userID)
	})
}

func (u *UserCacheRedis) GetUsersInfo(ctx context.Context, userIDs []string) ([]*model.User, error) {
	return batchGetCache2(ctx, u.rcClient, u.expireTime, userIDs, u.getUserInfoKey, u.getUserID, u.userDB.Find)
}

func (u *UserCacheRedis) DelUsersInfo(userIDs ...string) cache.UserCache {
	keys := make([]string, 0, len(userIDs))
	for _, userID := range userIDs {
		keys = append(keys, u.getUserInfoKey(userID))
	}
	cache := u.CloneUserCache()
	cache.AddKeys(keys...)

	return cache
}

func (u *UserCacheRedis) GetUserGlobalRecvMsgOpt(ctx context.Context, userID string) (opt int, err error) {
	return getCache(
		ctx,
		u.rcClient,
		u.getUserGlobalRecvMsgOptKey(userID),
		u.expireTime,
		func(ctx context.Context) (int, error) {
			return u.userDB.GetUserGlobalRecvMsgOpt(ctx, userID)
		},
	)
}

func (u *UserCacheRedis) DelUsersGlobalRecvMsgOpt(userIDs ...string) cache.UserCache {
	keys := make([]string, 0, len(userIDs))
	for _, userID := range userIDs {
		keys = append(keys, u.getUserGlobalRecvMsgOptKey(userID))
	}
	cache := u.CloneUserCache()
	cache.AddKeys(keys...)

	return cache
}
