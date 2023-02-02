package localcache

import (
	"Open_IM/pkg/proto/conversation"
	"context"
	"google.golang.org/grpc"
	"sync"
)

type ConversationLocalCache struct {
	lock                              sync.Mutex
	SuperGroupRecvMsgNotNotifyUserIDs map[string][]string
	rpc                               *grpc.ClientConn
	conversation                      conversation.ConversationClient
}

func NewConversationLocalCache(rpc *grpc.ClientConn) ConversationLocalCache {
	return ConversationLocalCache{
		SuperGroupRecvMsgNotNotifyUserIDs: make(map[string][]string, 0),
		rpc:                               rpc,
		conversation:                      conversation.NewConversationClient(rpc),
	}
}

func (g *ConversationLocalCache) GetRecvMsgNotNotifyUserIDs(ctx context.Context, groupID string) []string {
	return []string{}
}
