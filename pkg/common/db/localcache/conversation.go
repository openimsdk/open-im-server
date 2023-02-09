package localcache

import (
	discoveryRegistry "Open_IM/pkg/discoveryregistry"
	"context"
	"sync"
)

type ConversationLocalCacheInterface interface {
	GetRecvMsgNotNotifyUserIDs(ctx context.Context, groupID string) ([]string, error)
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

func (g *ConversationLocalCache) GetRecvMsgNotNotifyUserIDs(ctx context.Context, groupID string) ([]string, error) {
	g.client.GetConn()
	return []string{}
}
