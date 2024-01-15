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
		local:  localcache.New[any](localcache.WithRedisDeleteSubscribe(config.Config.LocalCache.Conversation.Topic, cli)),
		client: client,
	}
}

type ConversationLocalCache struct {
	local  localcache.Cache[any]
	client rpcclient.ConversationRpcClient
}

func (c *ConversationLocalCache) GetConversationIDs(ctx context.Context, ownerUserID string) ([]string, error) {
	return localcache.AnyValue[[]string](c.local.Get(ctx, cachekey.GetConversationIDsKey(ownerUserID), func(ctx context.Context) (any, error) {
		return c.client.GetConversationIDs(ctx, ownerUserID)
	}))
}
