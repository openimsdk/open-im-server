package rpccache

import (
	"context"
	"github.com/OpenIMSDK/tools/log"
	"github.com/openimsdk/localcache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/cachekey"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient"
	"github.com/redis/go-redis/v9"
)

func NewFriendLocalCache(client rpcclient.FriendRpcClient, cli redis.UniversalClient) *FriendLocalCache {
	lc := config.Config.LocalCache.Friend
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
	return localcache.AnyValue[bool](f.local.GetLink(ctx, cachekey.GetIsFriendKey(possibleFriendUserID, userID), func(ctx context.Context) (any, error) {
		log.ZDebug(ctx, "FriendLocalCache IsFriend rpc", "possibleFriendUserID", possibleFriendUserID, "userID", userID)
		return f.client.IsFriend(ctx, possibleFriendUserID, userID)
	}, cachekey.GetFriendIDsKey(possibleFriendUserID)))
}

// IsBlack possibleBlackUserID selfUserID
func (f *FriendLocalCache) IsBlack(ctx context.Context, possibleBlackUserID, userID string) (val bool, err error) {
	log.ZDebug(ctx, "FriendLocalCache IsBlack req", "possibleBlackUserID", possibleBlackUserID, "userID", userID)
	defer func() {
		if err == nil {
			log.ZDebug(ctx, "FriendLocalCache IsBlack return", "value", val)
		} else {
			log.ZError(ctx, "FriendLocalCache IsBlack return", err)
		}
	}()
	return localcache.AnyValue[bool](f.local.GetLink(ctx, cachekey.GetIsBlackIDsKey(possibleBlackUserID, userID), func(ctx context.Context) (any, error) {
		log.ZDebug(ctx, "FriendLocalCache IsBlack rpc", "possibleBlackUserID", possibleBlackUserID, "userID", userID)
		return f.client.IsBlack(ctx, possibleBlackUserID, userID)
	}, cachekey.GetBlackIDsKey(userID)))
}
