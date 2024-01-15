package rpccache

import (
	"context"
	"github.com/openimsdk/localcache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/cachekey"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient"
	"github.com/redis/go-redis/v9"
)

func NewConversationLocalCache(client rpcclient.ConversationRpcClient, cli redis.UniversalClient) *ConversationLocalCache {
	return &ConversationLocalCache{
		client: client,
		local: localcache.New[any](
			localcache.WithLocalSlotNum(config.Config.LocalCache.Conversation.SlotNum),
			localcache.WithLocalSlotSize(config.Config.LocalCache.Conversation.SlotSize),
		),
	}
}

type ConversationLocalCache struct {
	client rpcclient.ConversationRpcClient
	local  localcache.Cache[any]
}

func (c *ConversationLocalCache) GetConversationIDs(ctx context.Context, ownerUserID string) ([]string, error) {
	return localcache.AnyValue[[]string](c.local.Get(ctx, cachekey.GetConversationIDsKey(ownerUserID), func(ctx context.Context) (any, error) {
		return c.client.GetConversationIDs(ctx, ownerUserID)
	}))
}
