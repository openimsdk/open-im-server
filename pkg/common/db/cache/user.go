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

	relationTb "github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	"github.com/dtm-labs/rockscache"
	"github.com/redis/go-redis/v9"
)

const (
	userExpireTime          = time.Second * 60 * 60 * 12
	userInfoKey             = "USER_INFO:"
	userGlobalRecvMsgOptKey = "USER_GLOBAL_RECV_MSG_OPT_KEY:"
)

type UserCache interface {
	metaCache
	NewCache() UserCache
	GetUserInfo(ctx context.Context, userID string) (userInfo *relationTb.UserModel, err error)
	GetUsersInfo(ctx context.Context, userIDs []string) ([]*relationTb.UserModel, error)
	DelUsersInfo(userIDs ...string) UserCache
	GetUserGlobalRecvMsgOpt(ctx context.Context, userID string) (opt int, err error)
	DelUsersGlobalRecvMsgOpt(userIDs ...string) UserCache
}

type UserCacheRedis struct {
	metaCache
	userDB     relationTb.UserModelInterface
	expireTime time.Duration
	rcClient   *rockscache.Client
}

func NewUserCacheRedis(rdb redis.UniversalClient, userDB relationTb.UserModelInterface, options rockscache.Options) UserCache {
	rcClient := rockscache.NewClient(rdb, options)
	return &UserCacheRedis{
		metaCache:  NewMetaCacheRedis(rcClient),
		userDB:     userDB,
		expireTime: userExpireTime,
		rcClient:   rcClient,
	}
}

func (u *UserCacheRedis) NewCache() UserCache {
	return &UserCacheRedis{
		metaCache:  NewMetaCacheRedis(u.rcClient, u.metaCache.GetPreDelKeys()...),
		userDB:     u.userDB,
		expireTime: u.expireTime,
		rcClient:   u.rcClient,
	}
}

func (u *UserCacheRedis) getUserInfoKey(userID string) string {
	return userInfoKey + userID
}

func (u *UserCacheRedis) getUserGlobalRecvMsgOptKey(userID string) string {
	return userGlobalRecvMsgOptKey + userID
}

func (u *UserCacheRedis) GetUserInfo(ctx context.Context, userID string) (userInfo *relationTb.UserModel, err error) {
	return getCache(ctx, u.rcClient, u.getUserInfoKey(userID), u.expireTime, func(ctx context.Context) (*relationTb.UserModel, error) {
		return u.userDB.Take(ctx, userID)
	})
}

func (u *UserCacheRedis) GetUsersInfo(ctx context.Context, userIDs []string) ([]*relationTb.UserModel, error) {
	var keys []string
	for _, userID := range userIDs {
		keys = append(keys, u.getUserInfoKey(userID))
	}
	return batchGetCache(ctx, u.rcClient, keys, u.expireTime, func(user *relationTb.UserModel, keys []string) (int, error) {
		for i, key := range keys {
			if key == u.getUserInfoKey(user.UserID) {
				return i, nil
			}
		}
		return 0, errIndex
	}, func(ctx context.Context) ([]*relationTb.UserModel, error) {
		return u.userDB.Find(ctx, userIDs)
	})
}

func (u *UserCacheRedis) DelUsersInfo(userIDs ...string) UserCache {
	var keys []string
	for _, userID := range userIDs {
		keys = append(keys, u.getUserInfoKey(userID))
	}
	cache := u.NewCache()
	cache.AddKeys(keys...)
	return cache
}

func (u *UserCacheRedis) GetUserGlobalRecvMsgOpt(ctx context.Context, userID string) (opt int, err error) {
	return getCache(ctx, u.rcClient, u.getUserGlobalRecvMsgOptKey(userID), u.expireTime, func(ctx context.Context) (int, error) {
		return u.userDB.GetUserGlobalRecvMsgOpt(ctx, userID)
	})
}

func (u *UserCacheRedis) DelUsersGlobalRecvMsgOpt(userIDs ...string) UserCache {
	var keys []string
	for _, userID := range userIDs {
		keys = append(keys, u.getUserGlobalRecvMsgOptKey(userID))
	}
	cache := u.NewCache()
	cache.AddKeys(keys...)
	return cache
}
