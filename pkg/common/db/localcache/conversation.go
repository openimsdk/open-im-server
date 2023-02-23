package localcache

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/discoveryregistry"
	"Open_IM/pkg/proto/conversation"
	"context"
	"sync"
)

type ConversationLocalCacheInterface interface {
	GetRecvMsgNotNotifyUserIDs(ctx context.Context, groupID string) ([]string, error)
}

type ConversationLocalCache struct {
	lock                              sync.Mutex
	SuperGroupRecvMsgNotNotifyUserIDs map[string][]string
	client                            discoveryregistry.SvcDiscoveryRegistry
}

func NewConversationLocalCache(client discoveryregistry.SvcDiscoveryRegistry) ConversationLocalCache {
	return ConversationLocalCache{
		SuperGroupRecvMsgNotNotifyUserIDs: make(map[string][]string, 0),
		client:                            client,
	}
}

func (g *ConversationLocalCache) GetRecvMsgNotNotifyUserIDs(ctx context.Context, groupID string) ([]string, error) {
	conn, err := g.client.GetConn(config.Config.RpcRegisterName.OpenImConversationName)
	if err != nil {
		return nil, err
	}
	client := conversation.NewConversationClient(conn)
	resp, err := client.GetRecvMsgNotNotifyUserIDs(ctx, &conversation.GetRecvMsgNotNotifyUserIDsReq{
		GroupID: groupID,
	})
	if err != nil {
		return nil, err
	}
	return resp.UserIDs, nil
}
