package rpccache

import (
	"context"
	"errors"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/cachekey"
	"github.com/openimsdk/open-im-server/v3/pkg/localcache"
	"github.com/openimsdk/open-im-server/v3/pkg/localcache/lru"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mcontext"
	"github.com/redis/go-redis/v9"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

func NewOnlineCache(user rpcclient.UserRpcClient, group *GroupLocalCache, rdb redis.UniversalClient) *OnlineCache {
	x := &OnlineCache{
		user:  user,
		group: group,
		local: lru.NewSlotLRU(1024, localcache.LRUStringHash, func() lru.LRU[string, []int32] {
			return lru.NewLayLRU[string, []int32](2048, cachekey.OnlineExpire, time.Second*3, localcache.EmptyTarget{}, func(key string, value []int32) {})
		}),
	}
	go func() {
		parseUserOnlineStatus := func(payload string) (string, []int32, error) {
			arr := strings.Split(payload, ":")
			if len(arr) == 0 {
				return "", nil, errors.New("invalid data")
			}
			userID := arr[len(arr)-1]
			if userID == "" {
				return "", nil, errors.New("userID is empty")
			}
			platformIDs := make([]int32, len(arr)-1)
			for i := range platformIDs {
				platformID, err := strconv.Atoi(arr[i])
				if err != nil {
					return "", nil, err
				}
				platformIDs[i] = int32(platformID)
			}
			return userID, platformIDs, nil
		}
		ctx := mcontext.SetOperationID(context.Background(), cachekey.OnlineChannel+strconv.FormatUint(rand.Uint64(), 10))
		for message := range rdb.Subscribe(ctx, cachekey.OnlineChannel).Channel() {
			userID, platformIDs, err := parseUserOnlineStatus(message.Payload)
			if err != nil {
				log.ZError(ctx, "OnlineCache redis subscribe parseUserOnlineStatus", err, "payload", message.Payload, "channel", message.Channel)
				continue
			}
			log.ZDebug(ctx, "OnlineCache setUserOnline", "userID", userID, "platformIDs", platformIDs, "payload", message.Payload)
			x.setUserOnline(userID, platformIDs)
			//if err := x.setUserOnline(ctx, userID, platformIDs); err != nil {
			//	log.ZError(ctx, "redis subscribe setUserOnline", err, "payload", message.Payload, "channel", message.Channel)
			//}
		}
	}()
	return x
}

type OnlineCache struct {
	user  rpcclient.UserRpcClient
	group *GroupLocalCache
	local lru.LRU[string, []int32]
}

func (o *OnlineCache) getUserOnlineKey(userID string) string {
	return "<u>" + userID
}

func (o *OnlineCache) GetUserOnlinePlatform(ctx context.Context, userID string) ([]int32, error) {
	return o.local.Get(userID, func() ([]int32, error) {
		return o.user.GetUserOnlinePlatform(ctx, userID)
	})
}

func (o *OnlineCache) GetUserOnline(ctx context.Context, userID string) (bool, error) {
	platformIDs, err := o.GetUserOnlinePlatform(ctx, userID)
	if err != nil {
		return false, err
	}
	return len(platformIDs) > 0, nil
}

func (o *OnlineCache) GetUsersOnline(ctx context.Context, userIDs []string) ([]string, error) {
	onlineUserIDs := make([]string, 0, len(userIDs))
	for _, userID := range userIDs {
		online, err := o.GetUserOnline(ctx, userID)
		if err != nil {
			return nil, err
		}
		if online {
			onlineUserIDs = append(onlineUserIDs, userID)
		}
	}
	log.ZDebug(ctx, "OnlineCache GetUsersOnline", "userIDs", userIDs, "onlineUserIDs", onlineUserIDs)
	return onlineUserIDs, nil
}

func (o *OnlineCache) GetGroupOnline(ctx context.Context, groupID string) ([]string, error) {
	userIDs, err := o.group.GetGroupMemberIDs(ctx, groupID)
	if err != nil {
		return nil, err
	}
	var onlineUserIDs []string
	for _, userID := range userIDs {
		online, err := o.GetUserOnline(ctx, userID)
		if err != nil {
			return nil, err
		}
		if online {
			onlineUserIDs = append(onlineUserIDs, userID)
		}
	}
	log.ZDebug(ctx, "OnlineCache GetGroupOnline", "groupID", groupID, "onlineUserIDs", onlineUserIDs)
	return onlineUserIDs, nil
}

func (o *OnlineCache) setUserOnline(userID string, platformIDs []int32) {
	o.local.SetHas(o.getUserOnlineKey(userID), platformIDs)
}
