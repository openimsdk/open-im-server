package localcache

import (
	"context"
	"github.com/OpenIMSDK/openKeeper"
	"sync"
)

type ConversationLocalCacheInterface interface {
	GetRecvMsgNotNotifyUserIDs(ctx context.Context, groupID string) []string
}

type ConversationLocalCache struct {
	lock                              sync.Mutex
	SuperGroupRecvMsgNotNotifyUserIDs map[string][]string
	zkClient                          *openKeeper.ZkClient
}

func NewConversationLocalCache(zkClient *openKeeper.ZkClient) ConversationLocalCache {
	return ConversationLocalCache{
		SuperGroupRecvMsgNotNotifyUserIDs: make(map[string][]string, 0),
		zkClient:                          zkClient,
	}
}

func (g *ConversationLocalCache) GetRecvMsgNotNotifyUserIDs(ctx context.Context, groupID string) []string {
	return []string{}
}
