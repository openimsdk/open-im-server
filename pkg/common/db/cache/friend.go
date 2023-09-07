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

	"github.com/OpenIMSDK/tools/utils"

	relationtb "github.com/openimsdk/open-im-server/v3/pkg/common/db/table/relation"
)

const (
	friendExpireTime    = time.Second * 60 * 60 * 12
	friendIDsKey        = "FRIEND_IDS:"
	TwoWayFriendsIDsKey = "COMMON_FRIENDS_IDS:"
	friendKey           = "FRIEND_INFO:"
)

// args fn will exec when no data in msgCache.
type FriendCache interface {
	metaCache
	NewCache() FriendCache
	GetFriendIDs(ctx context.Context, ownerUserID string) (friendIDs []string, err error)
	// call when friendID List changed
	DelFriendIDs(ownerUserID ...string) FriendCache
	// get single friendInfo from msgCache
	GetFriend(ctx context.Context, ownerUserID, friendUserID string) (friend *relationtb.FriendModel, err error)
	// del friend when friend info changed
	DelFriend(ownerUserID, friendUserID string) FriendCache
}

type FriendCacheRedis struct {
	metaCache
	friendDB   relationtb.FriendModelInterface
	expireTime time.Duration
	rcClient   *rockscache.Client
}

func NewFriendCacheRedis(
	rdb redis.UniversalClient,
	friendDB relationtb.FriendModelInterface,
	options rockscache.Options,
) FriendCache {
	rcClient := rockscache.NewClient(rdb, options)
	return &FriendCacheRedis{
		metaCache:  NewMetaCacheRedis(rcClient),
		friendDB:   friendDB,
		expireTime: friendExpireTime,
		rcClient:   rcClient,
	}
}

func (c *FriendCacheRedis) NewCache() FriendCache {
	return &FriendCacheRedis{
		rcClient:   c.rcClient,
		metaCache:  NewMetaCacheRedis(c.rcClient, c.metaCache.GetPreDelKeys()...),
		friendDB:   c.friendDB,
		expireTime: c.expireTime,
	}
}

func (f *FriendCacheRedis) getFriendIDsKey(ownerUserID string) string {
	return friendIDsKey + ownerUserID
}

func (f *FriendCacheRedis) getTwoWayFriendsIDsKey(ownerUserID string) string {
	return TwoWayFriendsIDsKey + ownerUserID
}

func (f *FriendCacheRedis) getFriendKey(ownerUserID, friendUserID string) string {
	return friendKey + ownerUserID + "-" + friendUserID
}

func (f *FriendCacheRedis) GetFriendIDs(ctx context.Context, ownerUserID string) (friendIDs []string, err error) {
	return getCache(
		ctx,
		f.rcClient,
		f.getFriendIDsKey(ownerUserID),
		f.expireTime,
		func(ctx context.Context) ([]string, error) {
			return f.friendDB.FindFriendUserIDs(ctx, ownerUserID)
		},
	)
}

func (f *FriendCacheRedis) DelFriendIDs(ownerUserID ...string) FriendCache {
	new := f.NewCache()
	var keys []string
	for _, userID := range ownerUserID {
		keys = append(keys, f.getFriendIDsKey(userID))
	}
	new.AddKeys(keys...)
	return new
}

// todo.
func (f *FriendCacheRedis) GetTwoWayFriendIDs(
	ctx context.Context,
	ownerUserID string,
) (twoWayFriendIDs []string, err error) {
	friendIDs, err := f.GetFriendIDs(ctx, ownerUserID)
	if err != nil {
		return nil, err
	}
	for _, friendID := range friendIDs {
		friendFriendID, err := f.GetFriendIDs(ctx, friendID)
		if err != nil {
			return nil, err
		}
		if utils.IsContain(ownerUserID, friendFriendID) {
			twoWayFriendIDs = append(twoWayFriendIDs, ownerUserID)
		}
	}
	return twoWayFriendIDs, nil
}

func (f *FriendCacheRedis) DelTwoWayFriendIDs(ctx context.Context, ownerUserID string) FriendCache {
	new := f.NewCache()
	new.AddKeys(f.getTwoWayFriendsIDsKey(ownerUserID))
	return new
}

func (f *FriendCacheRedis) GetFriend(
	ctx context.Context,
	ownerUserID, friendUserID string,
) (friend *relationtb.FriendModel, err error) {
	return getCache(
		ctx,
		f.rcClient,
		f.getFriendKey(ownerUserID, friendUserID),
		f.expireTime,
		func(ctx context.Context) (*relationtb.FriendModel, error) {
			return f.friendDB.Take(ctx, ownerUserID, friendUserID)
		},
	)
}

func (f *FriendCacheRedis) DelFriend(ownerUserID, friendUserID string) FriendCache {
	new := f.NewCache()
	new.AddKeys(f.getFriendKey(ownerUserID, friendUserID))
	return new
}
