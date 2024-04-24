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
	"github.com/openimsdk/open-im-server/v3/pkg/common/cachekey"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	relationtb "github.com/openimsdk/open-im-server/v3/pkg/common/db/table/relation"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/datautil"
	"github.com/redis/go-redis/v9"
)

const (
	friendExpireTime = time.Second * 60 * 60 * 12
	// FriendIDsKey        = "FRIEND_IDS:"
	// TwoWayFriendsIDsKey = "COMMON_FRIENDS_IDS:"
	// friendKey           = "FRIEND_INFO:".
)

// FriendCache is an interface for caching friend-related data.
type FriendCache interface {
	metaCache
	NewCache() FriendCache
	GetFriendIDs(ctx context.Context, ownerUserID string) (friendIDs []string, err error)
	// Called when friendID list changed
	DelFriendIDs(ownerUserID ...string) FriendCache
	// Get single friendInfo from the cache
	GetFriend(ctx context.Context, ownerUserID, friendUserID string) (friend *relationtb.FriendModel, err error)
	// Delete friend when friend info changed
	DelFriend(ownerUserID, friendUserID string) FriendCache
	// Delete friends when friends' info changed
	DelFriends(ownerUserID string, friendUserIDs []string) FriendCache
}

// FriendCacheRedis is an implementation of the FriendCache interface using Redis.
type FriendCacheRedis struct {
	metaCache
	friendDB   relationtb.FriendModelInterface
	expireTime time.Duration
	rcClient   *rockscache.Client
}

// NewFriendCacheRedis creates a new instance of FriendCacheRedis.
func NewFriendCacheRedis(rdb redis.UniversalClient, localCache *config.LocalCache, friendDB relationtb.FriendModelInterface,
	options rockscache.Options) FriendCache {
	rcClient := rockscache.NewClient(rdb, options)
	mc := NewMetaCacheRedis(rcClient)
	f := localCache.Friend
	log.ZDebug(context.Background(), "friend local cache init", "Topic", f.Topic, "SlotNum", f.SlotNum, "SlotSize", f.SlotSize, "enable", f.Enable())
	mc.SetTopic(f.Topic)
	mc.SetRawRedisClient(rdb)
	return &FriendCacheRedis{
		metaCache:  mc,
		friendDB:   friendDB,
		expireTime: friendExpireTime,
		rcClient:   rcClient,
	}
}

// NewCache creates a new instance of FriendCacheRedis with the same configuration.
func (f *FriendCacheRedis) NewCache() FriendCache {
	return &FriendCacheRedis{
		rcClient:   f.rcClient,
		metaCache:  f.Copy(),
		friendDB:   f.friendDB,
		expireTime: f.expireTime,
	}
}

// getFriendIDsKey returns the key for storing friend IDs in the cache.
func (f *FriendCacheRedis) getFriendIDsKey(ownerUserID string) string {
	return cachekey.GetFriendIDsKey(ownerUserID)
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
func (f *FriendCacheRedis) DelFriendIDs(ownerUserIDs ...string) FriendCache {
	newGroupCache := f.NewCache()
	keys := make([]string, 0, len(ownerUserIDs))
	for _, userID := range ownerUserIDs {
		keys = append(keys, f.getFriendIDsKey(userID))
	}
	newGroupCache.AddKeys(keys...)

	return newGroupCache
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
func (f *FriendCacheRedis) DelTwoWayFriendIDs(ctx context.Context, ownerUserID string) FriendCache {
	newFriendCache := f.NewCache()
	newFriendCache.AddKeys(f.getTwoWayFriendsIDsKey(ownerUserID))

	return newFriendCache
}

// GetFriend retrieves friend info from the cache or the database if not found.
func (f *FriendCacheRedis) GetFriend(ctx context.Context, ownerUserID, friendUserID string) (friend *relationtb.FriendModel, err error) {
	return getCache(ctx, f.rcClient, f.getFriendKey(ownerUserID,
		friendUserID), f.expireTime, func(ctx context.Context) (*relationtb.FriendModel, error) {
		return f.friendDB.Take(ctx, ownerUserID, friendUserID)
	})
}

// DelFriend deletes friend info from the cache.
func (f *FriendCacheRedis) DelFriend(ownerUserID, friendUserID string) FriendCache {
	newFriendCache := f.NewCache()
	newFriendCache.AddKeys(f.getFriendKey(ownerUserID, friendUserID))

	return newFriendCache
}

// DelFriends deletes multiple friend infos from the cache.
func (f *FriendCacheRedis) DelFriends(ownerUserID string, friendUserIDs []string) FriendCache {
	newFriendCache := f.NewCache()

	for _, friendUserID := range friendUserIDs {
		key := f.getFriendKey(ownerUserID, friendUserID)
		newFriendCache.AddKeys(key) // Assuming AddKeys marks the keys for deletion
	}

	return newFriendCache
}
