package api

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/a2r"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/conversation"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func NewConversation(discov discoveryregistry.SvcDiscoveryRegistry) *Conversation {
	conn, err := discov.GetConn(context.Background(), config.Config.RpcRegisterName.OpenImConversationName)
	if err != nil {
		// panic(err)
	}
	return &Conversation{conn: conn, discov: discov}
}

type Conversation struct {
	conn   *grpc.ClientConn
	discov discoveryregistry.SvcDiscoveryRegistry
}

func (o *Conversation) client(ctx context.Context) (conversation.ConversationClient, error) {
	c, err := o.discov.GetConn(ctx, config.Config.RpcRegisterName.OpenImConversationName)
	if err != nil {
		return nil, err
	}
	return conversation.NewConversationClient(c), nil
}

func (o *Conversation) GetAllConversations(c *gin.Context) {
	a2r.Call(conversation.ConversationClient.GetAllConversations, o.client, c)
}

func (o *Conversation) GetConversation(c *gin.Context) {
	a2r.Call(conversation.ConversationClient.GetConversation, o.client, c)
}

func (o *Conversation) GetConversations(c *gin.Context) {
	a2r.Call(conversation.ConversationClient.GetConversations, o.client, c)
}

// deprecated
func (o *Conversation) SetConversation(c *gin.Context) {
	a2r.Call(conversation.ConversationClient.SetConversation, o.client, c)
}

// deprecated
func (o *Conversation) BatchSetConversations(c *gin.Context) {
	a2r.Call(conversation.ConversationClient.BatchSetConversations, o.client, c)
}

func (o *Conversation) SetRecvMsgOpt(c *gin.Context) {
	a2r.Call(conversation.ConversationClient.SetRecvMsgOpt, o.client, c)
}

func (o *Conversation) ModifyConversationField(c *gin.Context) {
	a2r.Call(conversation.ConversationClient.ModifyConversationField, o.client, c)
}

func (o *Conversation) GetConversationsHasReadAndMaxSeq(c *gin.Context) {
	a2r.Call(conversation.ConversationClient.GetConversationsHasReadAndMaxSeq, o.client, c)
}

func (o *Conversation) SetConversations(c *gin.Context) {
	a2r.Call(conversation.ConversationClient.SetConversations, o.client, c)
}
