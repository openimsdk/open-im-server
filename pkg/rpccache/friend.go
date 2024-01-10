package rpccache

import (
	"context"
	"github.com/OpenIMSDK/tools/log"
	"github.com/openimsdk/open-im-server/v3/pkg/common/cachekey"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/localcache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/localcache/option"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient"
	"github.com/redis/go-redis/v9"
)

func NewFriendLocalCache(client rpcclient.FriendRpcClient, cli redis.UniversalClient) *FriendLocalCache {
	return &FriendLocalCache{
		local:  localcache.New[any](localcache.WithRedisDeleteSubscribe(config.Config.LocalCache.Friend.Topic, cli)),
		client: client,
	}
}

type FriendLocalCache struct {
	local  localcache.Cache[any]
	client rpcclient.FriendRpcClient
}

func (f *FriendLocalCache) GetFriendIDs(ctx context.Context, ownerUserID string) (val []string, err error) {
	log.ZDebug(ctx, "FriendLocalCache GetFriendIDs req", "ownerUserID", ownerUserID)
	defer func() {
		if err == nil {
			log.ZDebug(ctx, "FriendLocalCache GetFriendIDs return", "value", val)
		} else {
			log.ZError(ctx, "FriendLocalCache GetFriendIDs return", err)
		}
	}()
	return localcache.AnyValue[[]string](f.local.Get(ctx, cachekey.GetFriendIDsKey(ownerUserID), func(ctx context.Context) (any, error) {
		log.ZDebug(ctx, "FriendLocalCache GetFriendIDs call rpc", "ownerUserID", ownerUserID)
		return f.client.GetFriendIDs(ctx, ownerUserID)
	}))
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
	return localcache.AnyValue[bool](f.local.Get(ctx, cachekey.GetIsFriendKey(possibleFriendUserID, userID), func(ctx context.Context) (any, error) {
		log.ZDebug(ctx, "FriendLocalCache IsFriend rpc", "possibleFriendUserID", possibleFriendUserID, "userID", userID)
		return f.client.IsFriend(ctx, possibleFriendUserID, userID)
	}, option.NewOption().WithLink(cachekey.GetFriendIDsKey(possibleFriendUserID), cachekey.GetFriendIDsKey(userID))))
}

func (f *FriendLocalCache) IsBlack(ctx context.Context, possibleBlackUserID, userID string) (val bool, err error) {
	log.ZDebug(ctx, "FriendLocalCache IsBlack req", "possibleBlackUserID", possibleBlackUserID, "userID", userID)
	defer func() {
		if err == nil {
			log.ZDebug(ctx, "FriendLocalCache IsBlack return", "value", val)
		} else {
			log.ZError(ctx, "FriendLocalCache IsBlack return", err)
		}
	}()
	return localcache.AnyValue[bool](f.local.Get(ctx, cachekey.GetIsBlackIDsKey(possibleBlackUserID, userID), func(ctx context.Context) (any, error) {
		log.ZDebug(ctx, "FriendLocalCache IsBlack rpc", "possibleBlackUserID", possibleBlackUserID, "userID", userID)
		return f.client.IsBlack(ctx, possibleBlackUserID, userID)
	}, option.NewOption().WithLink(cachekey.GetBlackIDsKey(possibleBlackUserID), cachekey.GetBlackIDsKey(userID))))
}
