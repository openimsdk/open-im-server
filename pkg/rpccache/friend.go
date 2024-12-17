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

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/cachekey"
	"github.com/openimsdk/protocol/relation"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/localcache"
	"github.com/openimsdk/tools/log"
	"github.com/redis/go-redis/v9"
)

func NewFriendLocalCache(localCache *config.LocalCache, cli redis.UniversalClient) *FriendLocalCache {
	lc := localCache.Friend
	log.ZDebug(context.Background(), "FriendLocalCache", "topic", lc.Topic, "slotNum", lc.SlotNum, "slotSize", lc.SlotSize, "enable", lc.Enable())
	x := &FriendLocalCache{
		local: localcache.New[[]byte](
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
	local localcache.Cache[[]byte]
}

func (f *FriendLocalCache) IsFriend(ctx context.Context, possibleFriendUserID, userID string) (val bool, err error) {
	res, err := f.isFriend(ctx, possibleFriendUserID, userID)
	if err != nil {
		return false, err
	}
	return res.InUser1Friends, nil
}

func (f *FriendLocalCache) isFriend(ctx context.Context, possibleFriendUserID, userID string) (val *relation.IsFriendResp, err error) {
	log.ZDebug(ctx, "FriendLocalCache isFriend req", "possibleFriendUserID", possibleFriendUserID, "userID", userID)
	defer func() {
		if err == nil {
			log.ZDebug(ctx, "FriendLocalCache isFriend return", "possibleFriendUserID", possibleFriendUserID, "userID", userID, "value", val)
		} else {
			log.ZError(ctx, "FriendLocalCache isFriend return", err, "possibleFriendUserID", possibleFriendUserID, "userID", userID)
		}
	}()
	var cache cacheProto[relation.IsFriendResp]
	return cache.Unmarshal(f.local.GetLink(ctx, cachekey.GetIsFriendKey(possibleFriendUserID, userID), func(ctx context.Context) ([]byte, error) {
		log.ZDebug(ctx, "FriendLocalCache isFriend rpc", "possibleFriendUserID", possibleFriendUserID, "userID", userID)
		return cache.Marshal(relation.IsFriendCaller.Invoke(ctx, &relation.IsFriendReq{UserID1: userID, UserID2: possibleFriendUserID}))
	}, cachekey.GetFriendIDsKey(possibleFriendUserID)))
}

// IsBlack possibleBlackUserID selfUserID.
func (f *FriendLocalCache) IsBlack(ctx context.Context, possibleBlackUserID, userID string) (val bool, err error) {
	res, err := f.isBlack(ctx, possibleBlackUserID, userID)
	if err != nil {
		return false, err
	}
	return res.InUser2Blacks, nil
}

// IsBlack possibleBlackUserID selfUserID.
func (f *FriendLocalCache) isBlack(ctx context.Context, possibleBlackUserID, userID string) (val *relation.IsBlackResp, err error) {
	log.ZDebug(ctx, "FriendLocalCache isBlack req", "possibleBlackUserID", possibleBlackUserID, "userID", userID)
	defer func() {
		if err == nil {
			log.ZDebug(ctx, "FriendLocalCache isBlack return", "possibleBlackUserID", possibleBlackUserID, "userID", userID, "value", val)
		} else {
			log.ZError(ctx, "FriendLocalCache isBlack return", err, "possibleBlackUserID", possibleBlackUserID, "userID", userID)
		}
	}()
	var cache cacheProto[relation.IsBlackResp]
	return cache.Unmarshal(f.local.GetLink(ctx, cachekey.GetIsBlackIDsKey(possibleBlackUserID, userID), func(ctx context.Context) ([]byte, error) {
		log.ZDebug(ctx, "FriendLocalCache IsBlack rpc", "possibleBlackUserID", possibleBlackUserID, "userID", userID)
		return cache.Marshal(relation.IsBlackCaller.Invoke(ctx, &relation.IsBlackReq{UserID1: possibleBlackUserID, UserID2: userID}))
	}, cachekey.GetBlackIDsKey(userID)))
}
