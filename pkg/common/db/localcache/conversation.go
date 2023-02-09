package localcache

import (
	discoveryRegistry "Open_IM/pkg/discovery_registry"
	"context"
	"sync"
)

type ConversationLocalCacheInterface interface {
	GetRecvMsgNotNotifyUserIDs(ctx context.Context, groupID string) []string
}

type ConversationLocalCache struct {
	lock                              sync.Mutex
	SuperGroupRecvMsgNotNotifyUserIDs map[string][]string
	client                            discoveryRegistry.SvcDiscoveryRegistry
}

func NewConversationLocalCache(client discoveryRegistry.SvcDiscoveryRegistry) ConversationLocalCache {
	return ConversationLocalCache{
		SuperGroupRecvMsgNotNotifyUserIDs: make(map[string][]string, 0),
		client:                            client,
	}
}

func (g *ConversationLocalCache) GetRecvMsgNotNotifyUserIDs(ctx context.Context, groupID string) []string {
	g.client.GetConn()
	return []string{}
}
