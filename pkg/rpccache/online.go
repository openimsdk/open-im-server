package rpccache

import (
	"context"
	"fmt"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcli"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/user"
	"math/rand"
	"strconv"
	"sync"
	"sync/atomic"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/cachekey"
	"github.com/openimsdk/open-im-server/v3/pkg/localcache"
	"github.com/openimsdk/open-im-server/v3/pkg/localcache/lru"
	"github.com/openimsdk/open-im-server/v3/pkg/util/useronline"
	"github.com/openimsdk/tools/db/cacheutil"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mcontext"
	"github.com/redis/go-redis/v9"
)

func NewOnlineCache(client *rpcli.UserClient, group *GroupLocalCache, rdb redis.UniversalClient, fullUserCache bool, fn func(ctx context.Context, userID string, platformIDs []int32)) (*OnlineCache, error) {
	l := &sync.Mutex{}
	x := &OnlineCache{
		client:        client,
		group:         group,
		fullUserCache: fullUserCache,
		Lock:          l,
		Cond:          sync.NewCond(l),
	}

	ctx := mcontext.SetOperationID(context.TODO(), strconv.FormatInt(time.Now().UnixNano()+int64(rand.Uint32()), 10))

	switch x.fullUserCache {
	case true:
		log.ZDebug(ctx, "fullUserCache is true")
		x.mapCache = cacheutil.NewCache[string, []int32]()
		go func() {
			if err := x.initUsersOnlineStatus(ctx); err != nil {
				log.ZError(ctx, "initUsersOnlineStatus failed", err)
			}
		}()
	case false:
		log.ZDebug(ctx, "fullUserCache is false")
		x.lruCache = lru.NewSlotLRU(1024, localcache.LRUStringHash, func() lru.LRU[string, []int32] {
			return lru.NewLayLRU[string, []int32](2048, cachekey.OnlineExpire/2, time.Second*3, localcache.EmptyTarget{}, func(key string, value []int32) {})
		})
		x.CurrentPhase.Store(DoSubscribeOver)
		x.Cond.Broadcast()
	}

	go func() {
		x.doSubscribe(ctx, rdb, fn)
	}()
	return x, nil
}

const (
	Begin uint32 = iota
	DoOnlineStatusOver
	DoSubscribeOver
)

type OnlineCache struct {
	client *rpcli.UserClient
	group  *GroupLocalCache

	// fullUserCache if enabled, caches the online status of all users using mapCache;
	// otherwise, only a portion of users' online statuses (regardless of whether they are online) will be cached using lruCache.
	fullUserCache bool

	lruCache lru.LRU[string, []int32]
	mapCache *cacheutil.Cache[string, []int32]

	Lock         *sync.Mutex
	Cond         *sync.Cond
	CurrentPhase atomic.Uint32
}

func (o *OnlineCache) initUsersOnlineStatus(ctx context.Context) (err error) {
	log.ZDebug(ctx, "init users online status begin")

	var (
		totalSet      atomic.Int64
		maxTries      = 5
		retryInterval = time.Second * 5

		resp *user.GetAllOnlineUsersResp
	)

	defer func(t time.Time) {
		log.ZInfo(ctx, "init users online status end", "cost", time.Since(t), "totalSet", totalSet.Load())
		o.CurrentPhase.Store(DoOnlineStatusOver)
		o.Cond.Broadcast()
	}(time.Now())

	retryOperation := func(operation func() error, operationName string) error {
		for i := 0; i < maxTries; i++ {
			if err = operation(); err != nil {
				log.ZWarn(ctx, fmt.Sprintf("initUsersOnlineStatus: %s failed", operationName), err)
				time.Sleep(retryInterval)
			} else {
				return nil
			}
		}
		return err
	}

	cursor := uint64(0)
	for resp == nil || resp.NextCursor != 0 {
		if err = retryOperation(func() error {
			resp, err = o.client.GetAllOnlineUsers(ctx, cursor)
			if err != nil {
				return err
			}

			for _, u := range resp.StatusList {
				if u.Status == constant.Online {
					o.setUserOnline(u.UserID, u.PlatformIDs)
				}
				totalSet.Add(1)
			}
			cursor = resp.NextCursor
			return nil
		}, "getAllOnlineUsers"); err != nil {
			return err
		}
	}

	return nil
}

func (o *OnlineCache) doSubscribe(ctx context.Context, rdb redis.UniversalClient, fn func(ctx context.Context, userID string, platformIDs []int32)) {
	o.Lock.Lock()
	ch := rdb.Subscribe(ctx, cachekey.OnlineChannel).Channel()
	for o.CurrentPhase.Load() < DoOnlineStatusOver {
		o.Cond.Wait()
	}
	o.Lock.Unlock()
	log.ZInfo(ctx, "begin doSubscribe")

	doMessage := func(message *redis.Message) {
		userID, platformIDs, err := useronline.ParseUserOnlineStatus(message.Payload)
		if err != nil {
			log.ZError(ctx, "OnlineCache setHasUserOnline redis subscribe parseUserOnlineStatus", err, "payload", message.Payload, "channel", message.Channel)
			return
		}
		log.ZDebug(ctx, fmt.Sprintf("get subscribe %s message", cachekey.OnlineChannel), "useID", userID, "platformIDs", platformIDs)
		switch o.fullUserCache {
		case true:
			if len(platformIDs) == 0 {
				// offline
				o.mapCache.Delete(userID)
			} else {
				o.mapCache.Store(userID, platformIDs)
			}
		case false:
			storageCache := o.setHasUserOnline(userID, platformIDs)
			log.ZDebug(ctx, "OnlineCache setHasUserOnline", "userID", userID, "platformIDs", platformIDs, "payload", message.Payload, "storageCache", storageCache)
			if fn != nil {
				fn(ctx, userID, platformIDs)
			}
		}
	}

	if o.CurrentPhase.Load() == DoOnlineStatusOver {
		for done := false; !done; {
			select {
			case message := <-ch:
				doMessage(message)
			default:
				o.CurrentPhase.Store(DoSubscribeOver)
				o.Cond.Broadcast()
				done = true
			}
		}
	}

	for message := range ch {
		doMessage(message)
	}
}

func (o *OnlineCache) getUserOnlinePlatform(ctx context.Context, userID string) ([]int32, error) {
	platformIDs, err := o.lruCache.Get(userID, func() ([]int32, error) {
		return o.client.GetUserOnlinePlatform(ctx, userID)
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

func (o *OnlineCache) getUserOnlinePlatformBatch(ctx context.Context, userIDs []string) (map[string][]int32, error) {
	platformIDsMap, err := o.lruCache.GetBatch(userIDs, func(missingUsers []string) (map[string][]int32, error) {
		platformIDsMap := make(map[string][]int32)
		usersStatus, err := o.client.GetUsersOnlinePlatform(ctx, missingUsers)
		if err != nil {
			return nil, err
		}

		for _, u := range usersStatus {
			platformIDsMap[u.UserID] = u.PlatformIDs
		}

		return platformIDsMap, nil
	})
	if err != nil {
		log.ZError(ctx, "OnlineCache GetUserOnlinePlatform", err, "userID", userIDs)
		return nil, err
	}
	return platformIDsMap, nil
}

func (o *OnlineCache) GetUsersOnline(ctx context.Context, userIDs []string) ([]string, []string, error) {
	t := time.Now()

	var (
		onlineUserIDs  = make([]string, 0, len(userIDs))
		offlineUserIDs = make([]string, 0, len(userIDs))
	)

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
		userOnlineMap, err := o.getUserOnlinePlatformBatch(ctx, userIDs)
		if err != nil {
			return nil, nil, err
		}

		for key, value := range userOnlineMap {
			if len(value) > 0 {
				onlineUserIDs = append(onlineUserIDs, key)
			} else {
				offlineUserIDs = append(offlineUserIDs, key)
			}
		}
	}

	log.ZInfo(ctx, "get users online", "online users length", len(userIDs), "offline users length", len(offlineUserIDs), "cost", time.Since(t))
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
