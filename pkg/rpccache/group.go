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
	return &GroupLocalCache{
		client: client,
		local: localcache.New[any](
			localcache.WithLocalSlotNum(config.Config.LocalCache.Group.SlotNum),
			localcache.WithLocalSlotSize(config.Config.LocalCache.Group.SlotSize),
		),
	}
}

type GroupLocalCache struct {
	client rpcclient.GroupRpcClient
	local  localcache.Cache[any]
}

func (g *GroupLocalCache) GetGroupMemberIDs(ctx context.Context, groupID string) ([]string, error) {
	return localcache.AnyValue[[]string](g.local.Get(ctx, cachekey.GetGroupMemberIDsKey(groupID), func(ctx context.Context) (any, error) {
		return g.client.GetGroupMemberIDs(ctx, groupID)
	}))
}
