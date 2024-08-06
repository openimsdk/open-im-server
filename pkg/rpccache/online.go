package rpccache

import (
	"context"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/cachekey"
	"github.com/openimsdk/open-im-server/v3/pkg/localcache"
	"github.com/openimsdk/open-im-server/v3/pkg/localcache/lru"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient"
	"github.com/openimsdk/open-im-server/v3/pkg/util/useronline"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mcontext"
	"github.com/redis/go-redis/v9"
	"math/rand"
	"strconv"
	"time"
)

func NewOnlineCache(user rpcclient.UserRpcClient, group *GroupLocalCache, rdb redis.UniversalClient, fn func(ctx context.Context, userID string, platformIDs []int32)) *OnlineCache {
	x := &OnlineCache{
		user:  user,
		group: group,
		local: lru.NewSlotLRU(1024, localcache.LRUStringHash, func() lru.LRU[string, []int32] {
			return lru.NewLayLRU[string, []int32](2048, cachekey.OnlineExpire/2, time.Second*3, localcache.EmptyTarget{}, func(key string, value []int32) {})
		}),
	}
	go func() {
		ctx := mcontext.SetOperationID(context.Background(), cachekey.OnlineChannel+strconv.FormatUint(rand.Uint64(), 10))
		for message := range rdb.Subscribe(ctx, cachekey.OnlineChannel).Channel() {
			userID, platformIDs, err := useronline.ParseUserOnlineStatus(message.Payload)
			if err != nil {
				log.ZError(ctx, "OnlineCache setUserOnline redis subscribe parseUserOnlineStatus", err, "payload", message.Payload, "channel", message.Channel)
				continue
			}
			storageCache := x.setUserOnline(userID, platformIDs)
			log.ZDebug(ctx, "OnlineCache setUserOnline", "userID", userID, "platformIDs", platformIDs, "payload", message.Payload, "storageCache", storageCache)
			if fn != nil {
				fn(ctx, userID, platformIDs)
			}
		}
	}()
	return x
}

type OnlineCache struct {
	user  rpcclient.UserRpcClient
	group *GroupLocalCache
	local lru.LRU[string, []int32]
}

func (o *OnlineCache) getUserOnlinePlatform(ctx context.Context, userID string) ([]int32, error) {
	platformIDs, err := o.local.Get(userID, func() ([]int32, error) {
		return o.user.GetUserOnlinePlatform(ctx, userID)
	})
	if err != nil {
		log.ZError(ctx, "OnlineCache GetUserOnlinePlatform", err, "userID", userID)
		return nil, err
	}
	log.ZDebug(ctx, "OnlineCache GetUserOnlinePlatform", "userID", userID, "platformIDs", platformIDs)
	return platformIDs, nil
}

func (o *OnlineCache) GetUserOnlinePlatform(ctx context.Context, userID string) ([]int32, error) {
	platformIDs, err := o.getUserOnlinePlatform(ctx, userID)
	if err != nil {
		return nil, err
	}
	tmp := make([]int32, len(platformIDs))
	copy(tmp, platformIDs)
	return platformIDs, nil
}

func (o *OnlineCache) GetUserOnline(ctx context.Context, userID string) (bool, error) {
	platformIDs, err := o.getUserOnlinePlatform(ctx, userID)
	if err != nil {
		return false, err
	}
	return len(platformIDs) > 0, nil
}

//func (o *OnlineCache) GetUsersOnline(ctx context.Context, userIDs []string) ([]string, error) {
//	onlineUserIDs := make([]string, 0, len(userIDs))
//	for _, userID := range userIDs {
//		online, err := o.GetUserOnline(ctx, userID)
//		if err != nil {
//			return nil, err
//		}
//		if online {
//			onlineUserIDs = append(onlineUserIDs, userID)
//		}
//	}
//	log.ZDebug(ctx, "OnlineCache GetUsersOnline", "userIDs", userIDs, "onlineUserIDs", onlineUserIDs)
//	return onlineUserIDs, nil
//}
//
//func (o *OnlineCache) GetGroupOnline(ctx context.Context, groupID string) ([]string, error) {
//	userIDs, err := o.group.GetGroupMemberIDs(ctx, groupID)
//	if err != nil {
//		return nil, err
//	}
//	var onlineUserIDs []string
//	for _, userID := range userIDs {
//		online, err := o.GetUserOnline(ctx, userID)
//		if err != nil {
//			return nil, err
//		}
//		if online {
//			onlineUserIDs = append(onlineUserIDs, userID)
//		}
//	}
//	log.ZDebug(ctx, "OnlineCache GetGroupOnline", "groupID", groupID, "onlineUserIDs", onlineUserIDs, "allUserID", userIDs)
//	return onlineUserIDs, nil
//}

func (o *OnlineCache) setUserOnline(userID string, platformIDs []int32) bool {
	return o.local.SetHas(userID, platformIDs)
}
