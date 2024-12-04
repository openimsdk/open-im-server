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
	"github.com/openimsdk/tools/utils/datautil"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/cachekey"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
)

const (
	friendExpireTime = time.Second * 60 * 60 * 12
)

// FriendCacheRedis is an implementation of the FriendCache interface using Redis.
type FriendCacheRedis struct {
	cache.BatchDeleter
	friendDB   database.Friend
	expireTime time.Duration
	rcClient   *rockscache.Client
	syncCount  int
}

// NewFriendCacheRedis creates a new instance of FriendCacheRedis.
func NewFriendCacheRedis(rdb redis.UniversalClient, localCache *config.LocalCache, friendDB database.Friend,
	options *rockscache.Options) cache.FriendCache {
	batchHandler := NewBatchDeleterRedis(rdb, options, []string{localCache.Friend.Topic})
	f := localCache.Friend
	log.ZDebug(context.Background(), "friend local cache init", "Topic", f.Topic, "SlotNum", f.SlotNum, "SlotSize", f.SlotSize, "enable", f.Enable())
	return &FriendCacheRedis{
		BatchDeleter: batchHandler,
		friendDB:     friendDB,
		expireTime:   friendExpireTime,
		rcClient:     rockscache.NewClient(rdb, *options),
	}
}

func (f *FriendCacheRedis) CloneFriendCache() cache.FriendCache {
	return &FriendCacheRedis{
		BatchDeleter: f.BatchDeleter.Clone(),
		friendDB:     f.friendDB,
		expireTime:   f.expireTime,
		rcClient:     f.rcClient,
	}
}

// getFriendIDsKey returns the key for storing friend IDs in the cache.
func (f *FriendCacheRedis) getFriendIDsKey(ownerUserID string) string {
	return cachekey.GetFriendIDsKey(ownerUserID)
}

func (f *FriendCacheRedis) getFriendMaxVersionKey(ownerUserID string) string {
	return cachekey.GetFriendMaxVersionKey(ownerUserID)
}

// getTwoWayFriendsIDsKey returns the key for storing two-way friend IDs in the cache.
func (f *FriendCacheRedis) getTwoWayFriendsIDsKey(ownerUserID string) string {
	return cachekey.GetTwoWayFriendsIDsKey(ownerUserID)
}

// getFriendKey returns the key for storing friend info in the cache.
func (f *FriendCacheRedis) getFriendKey(ownerUserID, friendUserID string) string {
	return cachekey.GetFriendKey(ownerUserID, friendUserID)
}

// GetFriendIDs retrieves friend IDs from the cache or the database if not found.
func (f *FriendCacheRedis) GetFriendIDs(ctx context.Context, ownerUserID string) (friendIDs []string, err error) {
	return getCache(ctx, f.rcClient, f.getFriendIDsKey(ownerUserID), f.expireTime, func(ctx context.Context) ([]string, error) {
		return f.friendDB.FindFriendUserIDs(ctx, ownerUserID)
	})
}

// DelFriendIDs deletes friend IDs from the cache.
func (f *FriendCacheRedis) DelFriendIDs(ownerUserIDs ...string) cache.FriendCache {
	newFriendCache := f.CloneFriendCache()
	keys := make([]string, 0, len(ownerUserIDs))
	for _, userID := range ownerUserIDs {
		keys = append(keys, f.getFriendIDsKey(userID))
	}
	newFriendCache.AddKeys(keys...)

	return newFriendCache
}

// GetTwoWayFriendIDs retrieves two-way friend IDs from the cache.
func (f *FriendCacheRedis) GetTwoWayFriendIDs(ctx context.Context, ownerUserID string) (twoWayFriendIDs []string, err error) {
	friendIDs, err := f.GetFriendIDs(ctx, ownerUserID)
	if err != nil {
		return nil, err
	}
	for _, friendID := range friendIDs {
		friendFriendID, err := f.GetFriendIDs(ctx, friendID)
		if err != nil {
			return nil, err
		}
		if datautil.Contain(ownerUserID, friendFriendID...) {
			twoWayFriendIDs = append(twoWayFriendIDs, ownerUserID)
		}
	}

	return twoWayFriendIDs, nil
}

// DelTwoWayFriendIDs deletes two-way friend IDs from the cache.
func (f *FriendCacheRedis) DelTwoWayFriendIDs(ctx context.Context, ownerUserID string) cache.FriendCache {
	newFriendCache := f.CloneFriendCache()
	newFriendCache.AddKeys(f.getTwoWayFriendsIDsKey(ownerUserID))

	return newFriendCache
}

// GetFriend retrieves friend info from the cache or the database if not found.
func (f *FriendCacheRedis) GetFriend(ctx context.Context, ownerUserID, friendUserID string) (friend *model.Friend, err error) {
	return getCache(ctx, f.rcClient, f.getFriendKey(ownerUserID,
		friendUserID), f.expireTime, func(ctx context.Context) (*model.Friend, error) {
		return f.friendDB.Take(ctx, ownerUserID, friendUserID)
	})
}

// DelFriend deletes friend info from the cache.
func (f *FriendCacheRedis) DelFriend(ownerUserID, friendUserID string) cache.FriendCache {
	newFriendCache := f.CloneFriendCache()
	newFriendCache.AddKeys(f.getFriendKey(ownerUserID, friendUserID))

	return newFriendCache
}

// DelFriends deletes multiple friend infos from the cache.
func (f *FriendCacheRedis) DelFriends(ownerUserID string, friendUserIDs []string) cache.FriendCache {
	newFriendCache := f.CloneFriendCache()

	for _, friendUserID := range friendUserIDs {
		key := f.getFriendKey(ownerUserID, friendUserID)
		newFriendCache.AddKeys(key) // Assuming AddKeys marks the keys for deletion
	}

	return newFriendCache
}

func (f *FriendCacheRedis) DelOwner(friendUserID string, ownerUserIDs []string) cache.FriendCache {
	newFriendCache := f.CloneFriendCache()

	for _, ownerUserID := range ownerUserIDs {
		key := f.getFriendKey(ownerUserID, friendUserID)
		newFriendCache.AddKeys(key) // Assuming AddKeys marks the keys for deletion
	}

	return newFriendCache
}

func (f *FriendCacheRedis) DelMaxFriendVersion(ownerUserIDs ...string) cache.FriendCache {
	newFriendCache := f.CloneFriendCache()
	for _, ownerUserID := range ownerUserIDs {
		key := f.getFriendMaxVersionKey(ownerUserID)
		newFriendCache.AddKeys(key) // Assuming AddKeys marks the keys for deletion
	}

	return newFriendCache
}

func (f *FriendCacheRedis) FindMaxFriendVersion(ctx context.Context, ownerUserID string) (*model.VersionLog, error) {
	return getCache(ctx, f.rcClient, f.getFriendMaxVersionKey(ownerUserID), f.expireTime, func(ctx context.Context) (*model.VersionLog, error) {
		return f.friendDB.FindIncrVersion(ctx, ownerUserID, 0, 0)
	})
}
