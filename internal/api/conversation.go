package api

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/internal/api/a2r"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/conversation"
	"github.com/gin-gonic/gin"
)

var _ context.Context // 解决goland编辑器bug

func NewConversation(c discoveryregistry.SvcDiscoveryRegistry) *Conversation {
	return &Conversation{c: c}
}

type Conversation struct {
	c discoveryregistry.SvcDiscoveryRegistry
}

func (o *Conversation) client() (conversation.ConversationClient, error) {
	conn, err := o.c.GetConn(config.Config.RpcRegisterName.OpenImConversationName)
	if err != nil {
		return nil, err
	}
	return conversation.NewConversationClient(conn), nil
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

func (o *Conversation) SetConversation(c *gin.Context) {
	a2r.Call(conversation.ConversationClient.SetConversation, o.client, c)
}

func (o *Conversation) BatchSetConversations(c *gin.Context) {
	a2r.Call(conversation.ConversationClient.BatchSetConversations, o.client, c)
}

func (o *Conversation) SetRecvMsgOpt(c *gin.Context) {
	a2r.Call(conversation.ConversationClient.SetRecvMsgOpt, o.client, c)
}

func (o *Conversation) ModifyConversationField(c *gin.Context) {
	a2r.Call(conversation.ConversationClient.ModifyConversationField, o.client, c)
}
