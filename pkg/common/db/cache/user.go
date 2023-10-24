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
	"errors"
	"hash/crc32"
	"strconv"
	"time"

	"github.com/OpenIMSDK/tools/log"

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
	SetUserStatus(ctx context.Context, userID string, status, platformID int32) error
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
	//var keys []string
	//for _, userID := range userIDs {
	//	keys = append(keys, u.getUserInfoKey(userID))
	//}
	//return batchGetCache(
	//	ctx,
	//	u.rcClient,
	//	keys,
	//	u.expireTime,
	//	func(user *relationtb.UserModel, keys []string) (int, error) {
	//		for i, key := range keys {
	//			if key == u.getUserInfoKey(user.UserID) {
	//				return i, nil
	//			}
	//		}
	//		return 0, errIndex
	//	},
	//	func(ctx context.Context) ([]*relationtb.UserModel, error) {
	//		return u.userDB.Find(ctx, userIDs)
	//	},
	//)
	return batchGetCache2(ctx, u.rcClient, u.expireTime, userIDs, func(userID string) string {
		return u.getUserInfoKey(userID)
	}, func(ctx context.Context, userID string) (*relationtb.UserModel, error) {
		return u.userDB.Take(ctx, userID)
	})
}

func (u *UserCacheRedis) DelUsersInfo(userIDs ...string) UserCache {
	keys := make([]string, 0, len(userIDs))
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
	keys := make([]string, 0, len(userIDs))
	for _, userID := range userIDs {
		keys = append(keys, u.getUserGlobalRecvMsgOptKey(userID))
	}
	cache := u.NewCache()
	cache.AddKeys(keys...)

	return cache
}

// GetUserStatus get user status.
func (u *UserCacheRedis) GetUserStatus(ctx context.Context, userIDs []string) ([]*user.OnlineStatus, error) {
	userStatus := make([]*user.OnlineStatus, 0, len(userIDs))
	for _, userID := range userIDs {
		UserIDNum := crc32.ChecksumIEEE([]byte(userID))
		modKey := strconv.Itoa(int(UserIDNum % statusMod))
		var onlineStatus user.OnlineStatus
		key := olineStatusKey + modKey
		result, err := u.rdb.HGet(ctx, key, userID).Result()
		if err != nil {
			if errors.Is(err, redis.Nil) {
				// key or field does not exist
				userStatus = append(userStatus, &user.OnlineStatus{
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
		onlineStatus.Status = constant.Online
		userStatus = append(userStatus, &onlineStatus)
	}

	return userStatus, nil
}

// SetUserStatus Set the user status and save it in redis.
func (u *UserCacheRedis) SetUserStatus(ctx context.Context, userID string, status, platformID int32) error {
	UserIDNum := crc32.ChecksumIEEE([]byte(userID))
	modKey := strconv.Itoa(int(UserIDNum % statusMod))
	key := olineStatusKey + modKey
	log.ZDebug(ctx, "SetUserStatus args", "userID", userID, "status", status,
		"platformID", platformID, "modKey", modKey, "key", key)
	isNewKey, err := u.rdb.Exists(ctx, key).Result()
	if err != nil {
		return errs.Wrap(err)
	}
	if isNewKey == 0 {
		if status == constant.Online {
			onlineStatus := user.OnlineStatus{
				UserID:      userID,
				Status:      constant.Online,
				PlatformIDs: []int32{platformID},
			}
			jsonData, err2 := json.Marshal(&onlineStatus)
			if err2 != nil {
				return errs.Wrap(err2)
			}
			_, err2 = u.rdb.HSet(ctx, key, userID, string(jsonData)).Result()
			if err2 != nil {
				return errs.Wrap(err2)
			}
			u.rdb.Expire(ctx, key, userOlineStatusExpireTime)

			return nil
		}
	}

	isNil := false
	result, err := u.rdb.HGet(ctx, key, userID).Result()
	if err != nil {
		if errors.Is(err, redis.Nil) {
			isNil = true
		} else {
			return errs.Wrap(err)
		}
	}

	if status == constant.Offline {
		err = u.refreshStatusOffline(ctx, userID, status, platformID, isNil, err, result, key)
		if err != nil {
			return err
		}
	} else {
		err = u.refreshStatusOnline(ctx, userID, platformID, isNil, err, result, key)
		if err != nil {
			return errs.Wrap(err)
		}
	}

	return nil
}

func (u *UserCacheRedis) refreshStatusOffline(ctx context.Context, userID string, status, platformID int32, isNil bool, err error, result, key string) error {
	if isNil {
		log.ZWarn(ctx, "this user not online,maybe trigger order not right",
			err, "userStatus", status)

		return nil
	}
	var onlineStatus user.OnlineStatus
	err = json.Unmarshal([]byte(result), &onlineStatus)
	if err != nil {
		return errs.Wrap(err)
	}
	var newPlatformIDs []int32
	for _, val := range onlineStatus.PlatformIDs {
		if val != platformID {
			newPlatformIDs = append(newPlatformIDs, val)
		}
	}
	if newPlatformIDs == nil {
		_, err = u.rdb.HDel(ctx, key, userID).Result()
		if err != nil {
			return errs.Wrap(err)
		}
	} else {
		onlineStatus.PlatformIDs = newPlatformIDs
		newjsonData, err := json.Marshal(&onlineStatus)
		if err != nil {
			return errs.Wrap(err)
		}
		_, err = u.rdb.HSet(ctx, key, userID, string(newjsonData)).Result()
		if err != nil {
			return errs.Wrap(err)
		}
	}

	return nil
}

func (u *UserCacheRedis) refreshStatusOnline(ctx context.Context, userID string, platformID int32, isNil bool, err error, result, key string) error {
	var onlineStatus user.OnlineStatus
	if !isNil {
		err2 := json.Unmarshal([]byte(result), &onlineStatus)
		if err != nil {
			return errs.Wrap(err2)
		}
		onlineStatus.PlatformIDs = RemoveRepeatedElementsInList(append(onlineStatus.PlatformIDs, platformID))
	} else {
		onlineStatus.PlatformIDs = append(onlineStatus.PlatformIDs, platformID)
	}
	onlineStatus.Status = constant.Online
	onlineStatus.UserID = userID
	newjsonData, err := json.Marshal(&onlineStatus)
	if err != nil {
		return errs.Wrap(err)
	}
	_, err = u.rdb.HSet(ctx, key, userID, string(newjsonData)).Result()
	if err != nil {
		return errs.Wrap(err)
	}

	return nil
}

type Comparable interface {
	~int | ~string | ~float64 | ~int32
}

func RemoveRepeatedElementsInList[T Comparable](slc []T) []T {
	var result []T
	tempMap := map[T]struct{}{}
	for _, e := range slc {
		if _, found := tempMap[e]; !found {
			tempMap[e] = struct{}{}
			result = append(result, e)
		}
	}

	return result
}
