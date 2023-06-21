package localcache

import (
	"context"
	"sync"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/conversation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/rpcclient"
)

type ConversationLocalCache struct {
	lock                              sync.Mutex
	superGroupRecvMsgNotNotifyUserIDs map[string]Hash
	conversationIDs                   map[string]Hash
	client                            *rpcclient.Conversation
}

type Hash struct {
	hash uint64
	ids  []string
}

func NewConversationLocalCache(discov discoveryregistry.SvcDiscoveryRegistry) *ConversationLocalCache {
	return &ConversationLocalCache{
		superGroupRecvMsgNotNotifyUserIDs: make(map[string]Hash),
		conversationIDs:                   make(map[string]Hash),
		client:                            rpcclient.NewConversation(discov),
	}
}

func (g *ConversationLocalCache) GetRecvMsgNotNotifyUserIDs(ctx context.Context, groupID string) ([]string, error) {
	resp, err := g.client.Client.GetRecvMsgNotNotifyUserIDs(ctx, &conversation.GetRecvMsgNotNotifyUserIDsReq{
		GroupID: groupID,
	})
	if err != nil {
		return nil, err
	}
	return resp.UserIDs, nil
}

func (g *ConversationLocalCache) GetConversationIDs(ctx context.Context, userID string) ([]string, error) {
	resp, err := g.client.Client.GetUserConversationIDsHash(ctx, &conversation.GetUserConversationIDsHashReq{
		OwnerUserID: userID,
	})
	if err != nil {
		return nil, err
	}
	g.lock.Lock()
	defer g.lock.Unlock()
	hash, ok := g.conversationIDs[userID]
	if !ok || hash.hash != resp.Hash {
		conversationIDsResp, err := g.client.Client.GetConversationIDs(ctx, &conversation.GetConversationIDsReq{
			UserID: userID,
		})
		if err != nil {
			return nil, err
		}
		g.conversationIDs[userID] = Hash{
			hash: resp.Hash,
			ids:  conversationIDsResp.ConversationIDs,
		}
		return conversationIDsResp.ConversationIDs, nil
	}
	return hash.ids, nil

}
