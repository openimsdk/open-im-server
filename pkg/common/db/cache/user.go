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
	"encoding/json"
	"hash/crc32"
	"strconv"
	"time"

	"github.com/OpenIMSDK/protocol/constant"

	"github.com/OpenIMSDK/protocol/user"
	"github.com/OpenIMSDK/tools/errs"

	"github.com/dtm-labs/rockscache"
	"github.com/redis/go-redis/v9"

	relationtb "github.com/openimsdk/open-im-server/v3/pkg/common/db/table/relation"
)

const (
	userExpireTime            = time.Second * 60 * 60 * 12
	userInfoKey               = "USER_INFO:"
	userGlobalRecvMsgOptKey   = "USER_GLOBAL_RECV_MSG_OPT_KEY:"
	olineStatusKey            = "ONLINE_STATUS:"
	userOlineStatusExpireTime = time.Second * 60 * 60 * 24
	statusMod                 = 501
	platformID                = "_PlatformIDSuffix"
)

type UserCache interface {
	metaCache
	NewCache() UserCache
	GetUserInfo(ctx context.Context, userID string) (userInfo *relationtb.UserModel, err error)
	GetUsersInfo(ctx context.Context, userIDs []string) ([]*relationtb.UserModel, error)
	DelUsersInfo(userIDs ...string) UserCache
	GetUserGlobalRecvMsgOpt(ctx context.Context, userID string) (opt int, err error)
	DelUsersGlobalRecvMsgOpt(userIDs ...string) UserCache
	GetUserStatus(ctx context.Context, userIDs []string) ([]*user.OnlineStatus, error)
	SetUserStatus(ctx context.Context, list []*user.OnlineStatus) error
}

type UserCacheRedis struct {
	metaCache
	rdb        redis.UniversalClient
	userDB     relationtb.UserModelInterface
	expireTime time.Duration
	rcClient   *rockscache.Client
}

func NewUserCacheRedis(
	rdb redis.UniversalClient,
	userDB relationtb.UserModelInterface,
	options rockscache.Options,
) UserCache {
	rcClient := rockscache.NewClient(rdb, options)
	return &UserCacheRedis{
		rdb:        rdb,
		metaCache:  NewMetaCacheRedis(rcClient),
		userDB:     userDB,
		expireTime: userExpireTime,
		rcClient:   rcClient,
	}
}

func (u *UserCacheRedis) NewCache() UserCache {
	return &UserCacheRedis{
		rdb:        u.rdb,
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

func (u *UserCacheRedis) getUserStatusHashKey(userID string, Id int32) string {
	return userID + "_" + string(Id) + platformID
}

func (u *UserCacheRedis) GetUserInfo(ctx context.Context, userID string) (userInfo *relationtb.UserModel, err error) {
	return getCache(
		ctx,
		u.rcClient,
		u.getUserInfoKey(userID),
		u.expireTime,
		func(ctx context.Context) (*relationtb.UserModel, error) {
			return u.userDB.Take(ctx, userID)
		},
	)
}

func (u *UserCacheRedis) GetUsersInfo(ctx context.Context, userIDs []string) ([]*relationtb.UserModel, error) {
	var keys []string
	for _, userID := range userIDs {
		keys = append(keys, u.getUserInfoKey(userID))
	}
	return batchGetCache(
		ctx,
		u.rcClient,
		keys,
		u.expireTime,
		func(user *relationtb.UserModel, keys []string) (int, error) {
			for i, key := range keys {
				if key == u.getUserInfoKey(user.UserID) {
					return i, nil
				}
			}
			return 0, errIndex
		},
		func(ctx context.Context) ([]*relationtb.UserModel, error) {
			return u.userDB.Find(ctx, userIDs)
		},
	)
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

func (u *UserCacheRedis) DelUsersGlobalRecvMsgOpt(userIDs ...string) UserCache {
	var keys []string
	for _, userID := range userIDs {
		keys = append(keys, u.getUserGlobalRecvMsgOptKey(userID))
	}
	cache := u.NewCache()
	cache.AddKeys(keys...)
	return cache
}

func (u *UserCacheRedis) getOnlineStatusKey(userID string) string {
	return olineStatusKey + userID
}

// GetUserStatus get user status.
func (u *UserCacheRedis) GetUserStatus(ctx context.Context, userIDs []string) ([]*user.OnlineStatus, error) {
	var res []*user.OnlineStatus
	for _, userID := range userIDs {
		UserIDNum := crc32.ChecksumIEEE([]byte(userID))
		modKey := strconv.Itoa(int(UserIDNum % statusMod))
		var onlineStatus user.OnlineStatus
		key := olineStatusKey + modKey
		result, err := u.rdb.HGet(ctx, key, userID).Result()
		if err != nil {
			if err == redis.Nil {
				// key or field does not exist
				res = append(res, &user.OnlineStatus{
					UserID:      userID,
					Status:      constant.Offline,
					PlatformIDs: nil,
				})
				continue
			} else {
				return nil, errs.Wrap(err)
			}
		}
		err = json.Unmarshal([]byte(result), &onlineStatus)
		if err != nil {
			return nil, errs.Wrap(err)
		}
		onlineStatus.UserID = userID
		res = append(res, &onlineStatus)
	}
	return res, nil
}

// SetUserStatus Set the user status and save it in redis.
func (u *UserCacheRedis) SetUserStatus(ctx context.Context, list []*user.OnlineStatus) error {
	for _, status := range list {
		var isNewKey int64
		UserIDNum := crc32.ChecksumIEEE([]byte(status.UserID))
		modKey := strconv.Itoa(int(UserIDNum % statusMod))
		key := olineStatusKey + modKey
		jsonData, err := json.Marshal(status)
		if err != nil {
			return errs.Wrap(err)
		}
		isNewKey, err = u.rdb.Exists(ctx, key).Result()
		if err != nil {
			return errs.Wrap(err)
		}
		if isNewKey == 0 {
			_, err = u.rdb.HSet(ctx, key, status.UserID, string(jsonData)).Result()
			if err != nil {
				return errs.Wrap(err)
			}
			u.rdb.Expire(ctx, key, userOlineStatusExpireTime)
		} else {
			result, err := u.rdb.HGet(ctx, key, status.UserID).Result()
			if err != nil {
				return errs.Wrap(err)
			}
			var onlineStatus user.OnlineStatus
			err = json.Unmarshal([]byte(result), &onlineStatus)
			if err != nil {
				return errs.Wrap(err)
			}
			onlineStatus.UserID = status.UserID
			if status.Status == constant.Offline {
				var newPlatformIDs []int32
				for _, val := range onlineStatus.PlatformIDs {
					if val != status.PlatformIDs[0] {
						newPlatformIDs = append(newPlatformIDs, val)
					}
				}
				if newPlatformIDs == nil {
					onlineStatus.Status = constant.Offline
					onlineStatus.PlatformIDs = []int32{}
					newjsonData, err := json.Marshal(&onlineStatus)
					if err != nil {
						return errs.Wrap(err)
					}
					_, err = u.rdb.HSet(ctx, key, status.UserID, string(newjsonData)).Result()
					if err != nil {
						return errs.Wrap(err)
					}
				} else {
					onlineStatus.PlatformIDs = newPlatformIDs
					newjsonData, err := json.Marshal(&onlineStatus)
					if err != nil {
						return errs.Wrap(err)
					}
					_, err = u.rdb.HSet(ctx, key, status.UserID, string(newjsonData)).Result()
					if err != nil {
						return errs.Wrap(err)
					}
				}
			} else {
				onlineStatus.Status = constant.Online
				// Judging whether to be kicked out.
				flag := false
				for _, val := range onlineStatus.PlatformIDs {
					if val == status.PlatformIDs[0] {
						flag = true
						break
					}
				}
				if !flag {
					onlineStatus.PlatformIDs = append(onlineStatus.PlatformIDs, status.PlatformIDs[0])
				}
				newjsonData, err := json.Marshal(&onlineStatus)
				if err != nil {
					return errs.Wrap(err)
				}
				_, err = u.rdb.HSet(ctx, key, status.UserID, string(newjsonData)).Result()
				if err != nil {
					return errs.Wrap(err)
				}
			}
		}
	}
	return nil
}
