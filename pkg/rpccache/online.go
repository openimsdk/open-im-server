package rpccache

import (
	"context"
	"math/rand"
	"strconv"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/cachekey"
	"github.com/openimsdk/open-im-server/v3/pkg/localcache"
	"github.com/openimsdk/open-im-server/v3/pkg/localcache/lru"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient"
	"github.com/openimsdk/open-im-server/v3/pkg/util/useronline"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/tools/db/cacheutil"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mcontext"
	"github.com/redis/go-redis/v9"
)

func NewOnlineCache(user rpcclient.UserRpcClient, group *GroupLocalCache, rdb redis.UniversalClient, fullUserCache bool, fn func(ctx context.Context, userID string, platformIDs []int32)) (*OnlineCache, error) {
	x := &OnlineCache{
		user:          user,
		group:         group,
		fullUserCache: fullUserCache,
	}

	switch x.fullUserCache {
	case true:
		x.mapCache = cacheutil.NewCache[string, []int32]()
		if err := x.initUsersOnlineStatus(mcontext.SetOperationID(context.TODO(), strconv.FormatInt(time.Now().UnixNano()+int64(rand.Uint32()), 10))); err != nil {
			return nil, err
		}
	case false:
		x.lruCache = lru.NewSlotLRU(1024, localcache.LRUStringHash, func() lru.LRU[string, []int32] {
			return lru.NewLayLRU[string, []int32](2048, cachekey.OnlineExpire/2, time.Second*3, localcache.EmptyTarget{}, func(key string, value []int32) {})
		})
	}

	go func() {
		ctx := mcontext.SetOperationID(context.Background(), cachekey.OnlineChannel+strconv.FormatUint(rand.Uint64(), 10))
		for message := range rdb.Subscribe(ctx, cachekey.OnlineChannel).Channel() {
			userID, platformIDs, err := useronline.ParseUserOnlineStatus(message.Payload)
			if err != nil {
				log.ZError(ctx, "OnlineCache setHasUserOnline redis subscribe parseUserOnlineStatus", err, "payload", message.Payload, "channel", message.Channel)
				continue
			}

			switch x.fullUserCache {
			case true:
				if len(platformIDs) == 0 {
					// offline
					x.mapCache.Delete(userID)
				} else {
					x.mapCache.Store(userID, platformIDs)
				}
			case false:
				storageCache := x.setHasUserOnline(userID, platformIDs)
				log.ZDebug(ctx, "OnlineCache setHasUserOnline", "userID", userID, "platformIDs", platformIDs, "payload", message.Payload, "storageCache", storageCache)
				if fn != nil {
					fn(ctx, userID, platformIDs)
				}
			}

		}
	}()
	return x, nil
}

type OnlineCache struct {
	user  rpcclient.UserRpcClient
	group *GroupLocalCache

	// fullUserCache if enabled, caches the online status of all users using mapCache;
	// otherwise, only a portion of users' online statuses (regardless of whether they are online) will be cached using lruCache.
	fullUserCache bool

	lruCache lru.LRU[string, []int32]
	mapCache *cacheutil.Cache[string, []int32]
}

func (o *OnlineCache) initUsersOnlineStatus(ctx context.Context) error {
	log.ZDebug(ctx, "init users online status begin")

	var (
		totalSet int
	)

	time.Sleep(time.Second * 10)

	defer func(t time.Time) {
		log.ZWarn(ctx, "init users online status end", nil, "cost", time.Since(t), "totalSet", totalSet)
	}(time.Now())

	for page := int32(1); ; page++ {
		resp, err := o.user.GetAllUserID(ctx, page, constant.ParamMaxLength)
		if err != nil {
			return err
		}

		usersStatus, err := o.user.GetUsersOnlinePlatform(ctx, resp.UserIDs)
		if err != nil {
			return err
		}

		for _, user := range usersStatus {
			if user.Status == constant.Online {
				o.setUserOnline(user.UserID, user.PlatformIDs)
			}
			totalSet++
		}

		if len(resp.UserIDs) < constant.ParamMaxLength {
			break
		}
	}
	return nil
}

func (o *OnlineCache) getUserOnlinePlatform(ctx context.Context, userID string) ([]int32, error) {
	platformIDs, err := o.lruCache.Get(userID, func() ([]int32, error) {
		return o.user.GetUserOnlinePlatform(ctx, userID)
	})
	if err != nil {
		log.ZError(ctx, "OnlineCache GetUserOnlinePlatform", err, "userID", userID)
		return nil, err
	}
	//log.ZDebug(ctx, "OnlineCache GetUserOnlinePlatform", "userID", userID, "platformIDs", platformIDs)
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

// func (o *OnlineCache) GetUserOnlinePlatformBatch(ctx context.Context, userIDs []string) (map[string]int32, error) {
// 	platformIDs, err := o.getUserOnlinePlatform(ctx, userIDs)
// 	if err != nil {
// 		return nil, err
// 	}
// 	tmp := make([]int32, len(platformIDs))
// 	copy(tmp, platformIDs)
// 	return platformIDs, nil
// }

func (o *OnlineCache) GetUserOnline(ctx context.Context, userID string) (bool, error) {
	platformIDs, err := o.getUserOnlinePlatform(ctx, userID)
	if err != nil {
		return false, err
	}
	return len(platformIDs) > 0, nil
}

//func (o *OnlineCache) getUserOnlinePlatformBatch(ctx context.Context, userIDs []string) (map[string][]int32, error) {
//	platformIDsMap, err := o.lruCache.GetBatch(userIDs, func(missingUsers []string) (map[string][]int32, error) {
//		platformIDsMap := make(map[string][]int32)
//
//		usersStatus, err := o.user.GetUsersOnlinePlatform(ctx, missingUsers)
//		if err != nil {
//			return nil, err
//		}
//
//		for _, user := range usersStatus {
//			platformIDsMap[user.UserID] = user.PlatformIDs
//		}
//
//		return platformIDsMap, nil
//	})
//	if err != nil {
//		log.ZError(ctx, "OnlineCache GetUserOnlinePlatform", err, "userID", userIDs)
//		return nil, err
//	}
//
//	//log.ZDebug(ctx, "OnlineCache GetUserOnlinePlatform", "userID", userID, "platformIDs", platformIDs)
//	return platformIDsMap, nil
//}

func (o *OnlineCache) GetUsersOnline(ctx context.Context, userIDs []string) ([]string, []string, error) {
	t := time.Now()

	var (
		onlineUserIDs  = make([]string, 0, len(userIDs))
		offlineUserIDs = make([]string, 0, len(userIDs))
	)

	//userOnlineMap, err := o.getUserOnlinePlatformBatch(ctx, userIDs)
	//if err != nil {
	//	return nil, nil, err
	//}
	//
	//for key, value := range userOnlineMap {
	//	if len(value) > 0 {
	//		onlineUserIDs = append(onlineUserIDs, key)
	//	} else {
	//		offlineUserIDs = append(offlineUserIDs, key)
	//	}
	//}

	switch o.fullUserCache {
	case true:
		for _, userID := range userIDs {
			if _, ok := o.mapCache.Load(userID); ok {
				onlineUserIDs = append(onlineUserIDs, userID)
			} else {
				offlineUserIDs = append(offlineUserIDs, userID)
			}
		}
	case false:
	}

	log.ZWarn(ctx, "get users online", nil, "online users length", len(userIDs), "offline users length", len(offlineUserIDs), "cost", time.Since(t))
	return userIDs, offlineUserIDs, nil
}

//func (o *OnlineCache) GetUsersOnline(ctx context.Context, userIDs []string) ([]string, error) {
//	onlineUserIDs := make([]string, 0, len(userIDs))
//	for _, userID := range userIDs {
// online, err := o.GetUserOnline(ctx, userID)
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

func (o *OnlineCache) setUserOnline(userID string, platformIDs []int32) {
	switch o.fullUserCache {
	case true:
		o.mapCache.Store(userID, platformIDs)
	case false:
		o.lruCache.Set(userID, platformIDs)
	}
}

func (o *OnlineCache) setHasUserOnline(userID string, platformIDs []int32) bool {
	return o.lruCache.SetHas(userID, platformIDs)
}
