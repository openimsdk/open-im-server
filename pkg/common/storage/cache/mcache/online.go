package mcache

import (
	"context"
	"sync"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache"
)

var (
	globalOnlineCache cache.OnlineCache
	globalOnlineOnce  sync.Once
)

func NewOnlineCache() cache.OnlineCache {
	globalOnlineOnce.Do(func() {
		globalOnlineCache = &onlineCache{
			user: make(map[string]map[int32]struct{}),
		}
	})
	return globalOnlineCache
}

type onlineCache struct {
	lock sync.RWMutex
	user map[string]map[int32]struct{}
}

func (x *onlineCache) GetOnline(ctx context.Context, userID string) ([]int32, error) {
	x.lock.RLock()
	defer x.lock.RUnlock()
	pSet, ok := x.user[userID]
	if !ok {
		return nil, nil
	}
	res := make([]int32, 0, len(pSet))
	for k := range pSet {
		res = append(res, k)
	}
	return res, nil
}

func (x *onlineCache) SetUserOnline(ctx context.Context, userID string, online, offline []int32) error {
	x.lock.Lock()
	defer x.lock.Unlock()
	pSet, ok := x.user[userID]
	if ok {
		for _, p := range offline {
			delete(pSet, p)
		}
	}
	if len(online) > 0 {
		if !ok {
			pSet = make(map[int32]struct{})
			x.user[userID] = pSet
		}
		for _, p := range online {
			pSet[p] = struct{}{}
		}
	}
	if len(pSet) == 0 {
		delete(x.user, userID)
	}
	return nil
}

func (x *onlineCache) GetAllOnlineUsers(ctx context.Context, cursor uint64) (map[string][]int32, uint64, error) {
	if cursor != 0 {
		return nil, 0, nil
	}
	x.lock.RLock()
	defer x.lock.RUnlock()
	res := make(map[string][]int32)
	for k, v := range x.user {
		pSet := make([]int32, 0, len(v))
		for p := range v {
			pSet = append(pSet, p)
		}
		res[k] = pSet
	}
	return res, 0, nil
}
