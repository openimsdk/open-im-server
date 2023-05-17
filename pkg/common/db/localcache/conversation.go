package localcache

import (
	"context"
	"sync"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/conversation"
)

type ConversationLocalCache struct {
	lock                              sync.Mutex
	SuperGroupRecvMsgNotNotifyUserIDs map[string]Hash
	ConversationIDs                   map[string]Hash
	client                            discoveryregistry.SvcDiscoveryRegistry
}

type Hash struct {
	hash uint64
	ids  []string
}

func NewConversationLocalCache(client discoveryregistry.SvcDiscoveryRegistry) *ConversationLocalCache {
	return &ConversationLocalCache{
		SuperGroupRecvMsgNotNotifyUserIDs: make(map[string]Hash),
		ConversationIDs:                   make(map[string]Hash),
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
	resp, err := client.GetUserConversationIDsHash(ctx, &conversation.GetUserConversationIDsHashReq{
		OwnerUserID: userID,
	})
	if err != nil {
		return nil, err
	}
	g.lock.Lock()
	defer g.lock.Unlock()
	hash, ok := g.ConversationIDs[userID]
	if !ok || hash.hash != resp.Hash {
		conversationIDsResp, err := client.GetConversationIDs(ctx, &conversation.GetConversationIDsReq{
			UserID: userID,
		})
		if err != nil {
			return nil, err
		}
		g.ConversationIDs[userID] = Hash{
			hash: resp.Hash,
			ids:  conversationIDsResp.ConversationIDs,
		}
		return conversationIDsResp.ConversationIDs, nil
	}
	return hash.ids, nil

}
