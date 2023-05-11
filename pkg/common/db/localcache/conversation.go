package localcache

import (
	"context"
	"sync"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/conversation"
)

type ConversationLocalCacheInterface interface {
	GetRecvMsgNotNotifyUserIDs(ctx context.Context, groupID string) ([]string, error)
	GetConversationIDs(ctx context.Context, userID string) ([]string, error)
}

type ConversationLocalCache struct {
	lock                              sync.Mutex
	SuperGroupRecvMsgNotNotifyUserIDs map[string][]string
	client                            discoveryregistry.SvcDiscoveryRegistry
}

func NewConversationLocalCache(client discoveryregistry.SvcDiscoveryRegistry) *ConversationLocalCache {
	return &ConversationLocalCache{
		SuperGroupRecvMsgNotNotifyUserIDs: make(map[string][]string, 0),
		client:                            client,
	}
}

func (g *ConversationLocalCache) GetRecvMsgNotNotifyUserIDs(ctx context.Context, groupID string) ([]string, error) {
	conn, err := g.client.GetConn(ctx, config.Config.RpcRegisterName.OpenImConversationName)
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

func (g *ConversationLocalCache) GetConversationIDs(ctx context.Context, userID string) ([]string, error) {
	conn, err := g.client.GetConn(ctx, config.Config.RpcRegisterName.OpenImConversationName)
	if err != nil {
		return nil, err
	}
	client := conversation.NewConversationClient(conn)
	resp, err := client.GetConversationIDs(ctx, &conversation.GetConversationIDsReq{
		UserID: userID,
	})
	if err != nil {
		return nil, err
	}
	return resp.ConversationIDs, nil
}
