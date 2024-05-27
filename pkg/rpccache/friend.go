// Copyright © 2024 OpenIM. All rights reserved.
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

package rpccache

import (
	"context"
	cachekey2 "github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/cachekey"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/localcache"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient"
	"github.com/openimsdk/tools/log"
	"github.com/redis/go-redis/v9"
)

func NewFriendLocalCache(client rpcclient.FriendRpcClient, localCache *config.LocalCache, cli redis.UniversalClient) *FriendLocalCache {
	lc := localCache.Friend
	log.ZDebug(context.Background(), "FriendLocalCache", "topic", lc.Topic, "slotNum", lc.SlotNum, "slotSize", lc.SlotSize, "enable", lc.Enable())
	x := &FriendLocalCache{
		client: client,
		local: localcache.New[any](
			localcache.WithLocalSlotNum(lc.SlotNum),
			localcache.WithLocalSlotSize(lc.SlotSize),
			localcache.WithLinkSlotNum(lc.SlotNum),
			localcache.WithLocalSuccessTTL(lc.Success()),
			localcache.WithLocalFailedTTL(lc.Failed()),
		),
	}
	if lc.Enable() {
		go subscriberRedisDeleteCache(context.Background(), cli, lc.Topic, x.local.DelLocal)
	}
	return x
}

type FriendLocalCache struct {
	client rpcclient.FriendRpcClient
	local  localcache.Cache[any]
}

func (f *FriendLocalCache) IsFriend(ctx context.Context, possibleFriendUserID, userID string) (val bool, err error) {
	log.ZDebug(ctx, "FriendLocalCache IsFriend req", "possibleFriendUserID", possibleFriendUserID, "userID", userID)
	defer func() {
		if err == nil {
			log.ZDebug(ctx, "FriendLocalCache IsFriend return", "value", val)
		} else {
			log.ZError(ctx, "FriendLocalCache IsFriend return", err)
		}
	}()
	return localcache.AnyValue[bool](f.local.GetLink(ctx, cachekey2.GetIsFriendKey(possibleFriendUserID, userID), func(ctx context.Context) (any, error) {
		log.ZDebug(ctx, "FriendLocalCache IsFriend rpc", "possibleFriendUserID", possibleFriendUserID, "userID", userID)
		return f.client.IsFriend(ctx, possibleFriendUserID, userID)
	}, cachekey2.GetFriendIDsKey(possibleFriendUserID)))
}

// IsBlack possibleBlackUserID selfUserID.
func (f *FriendLocalCache) IsBlack(ctx context.Context, possibleBlackUserID, userID string) (val bool, err error) {
	log.ZDebug(ctx, "FriendLocalCache IsBlack req", "possibleBlackUserID", possibleBlackUserID, "userID", userID)
	defer func() {
		if err == nil {
			log.ZDebug(ctx, "FriendLocalCache IsBlack return", "value", val)
		} else {
			log.ZError(ctx, "FriendLocalCache IsBlack return", err)
		}
	}()
	return localcache.AnyValue[bool](f.local.GetLink(ctx, cachekey2.GetIsBlackIDsKey(possibleBlackUserID, userID), func(ctx context.Context) (any, error) {
		log.ZDebug(ctx, "FriendLocalCache IsBlack rpc", "possibleBlackUserID", possibleBlackUserID, "userID", userID)
		return f.client.IsBlack(ctx, possibleBlackUserID, userID)
	}, cachekey2.GetBlackIDsKey(userID)))
}
