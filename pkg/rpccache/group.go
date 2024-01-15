package rpccache

import (
	"context"
	"github.com/openimsdk/localcache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/cachekey"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient"
	"github.com/redis/go-redis/v9"
)

func NewGroupLocalCache(client rpcclient.GroupRpcClient, cli redis.UniversalClient) *GroupLocalCache {
	lc := config.Config.LocalCache.Group
	x := &GroupLocalCache{
		client: client,
		local: localcache.New[any](
			localcache.WithLocalSlotNum(lc.SlotNum),
			localcache.WithLocalSlotSize(lc.SlotSize),
		),
	}
	go subscriberRedisDeleteCache(context.Background(), cli, lc.Topic, x.local.DelLocal)
	return x
}

type GroupLocalCache struct {
	client rpcclient.GroupRpcClient
	local  localcache.Cache[any]
}

type listMap[V comparable] struct {
	List []V
	Map  map[V]struct{}
}

func newListMap[V comparable](values []V, err error) (*listMap[V], error) {
	if err != nil {
		return nil, err
	}
	lm := &listMap[V]{
		List: values,
		Map:  make(map[V]struct{}, len(values)),
	}
	for _, value := range values {
		lm.Map[value] = struct{}{}
	}
	return lm, nil
}

func (g *GroupLocalCache) getGroupMemberIDs(ctx context.Context, groupID string) (*listMap[string], error) {
	return localcache.AnyValue[*listMap[string]](g.local.Get(ctx, cachekey.GetGroupMemberIDsKey(groupID), func(ctx context.Context) (any, error) {
		return newListMap(g.client.GetGroupMemberIDs(ctx, groupID))
	}))
}

func (g *GroupLocalCache) GetGroupMemberIDs(ctx context.Context, groupID string) ([]string, error) {
	res, err := g.getGroupMemberIDs(ctx, groupID)
	if err != nil {
		return nil, err
	}
	return res.List, nil
}

func (g *GroupLocalCache) GetGroupMemberIDMap(ctx context.Context, groupID string) (map[string]struct{}, error) {
	res, err := g.getGroupMemberIDs(ctx, groupID)
	if err != nil {
		return nil, err
	}
	return res.Map, nil
}
