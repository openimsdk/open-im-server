package notification

import (
	"context"
	"github.com/apache/dubbo-go"
	pbconv "github.com/OpenIMSDK/protocol/conversation"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient"
)

type Conversation struct {
	Client pbconv.ConversationClient
	conn   dubbo-go.ClientConnInterface
	discov discoveryregistry.SvcDiscoveryRegistry
}

func (c *Conversation) SendMessage(ctx context.Context, req *pbconv.SendMessageReq) (*pbconv.SendMessageResp, error) {
	return c.Client.SendMessage(ctx, req)
}

func (c *Conversation) GetConversation(ctx context.Context, req *pbconv.GetConversationReq) (*pbconv.GetConversationResp, error) {
	return c.Client.GetConversation(ctx, req)
}

func (c *Conversation) DeleteConversation(ctx context.Context, req *pbconv.DeleteConversationReq) (*pbconv.DeleteConversationResp, error) {
	return c.Client.DeleteConversation(ctx, req)
}
