package api

import (
	"github.com/OpenIMSDK/Open-IM-Server/pkg/a2r"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/conversation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/rpcclient"
	"github.com/gin-gonic/gin"
)

type ConversationApi rpcclient.Conversation

func NewConversationApi(discov discoveryregistry.SvcDiscoveryRegistry) ConversationApi {
	return ConversationApi(*rpcclient.NewConversation(discov))
}

func (o *ConversationApi) GetAllConversations(c *gin.Context) {
	a2r.Call(conversation.ConversationClient.GetAllConversations, o.Client, c)
}

func (o *ConversationApi) GetConversation(c *gin.Context) {
	a2r.Call(conversation.ConversationClient.GetConversation, o.Client, c)
}

func (o *ConversationApi) GetConversations(c *gin.Context) {
	a2r.Call(conversation.ConversationClient.GetConversations, o.Client, c)
}

// deprecated
func (o *ConversationApi) SetConversation(c *gin.Context) {
	a2r.Call(conversation.ConversationClient.SetConversation, o.Client, c)
}

// deprecated
func (o *ConversationApi) BatchSetConversations(c *gin.Context) {
	a2r.Call(conversation.ConversationClient.BatchSetConversations, o.Client, c)
}

func (o *ConversationApi) SetRecvMsgOpt(c *gin.Context) {
	a2r.Call(conversation.ConversationClient.SetRecvMsgOpt, o.Client, c)
}

func (o *ConversationApi) ModifyConversationField(c *gin.Context) {
	a2r.Call(conversation.ConversationClient.ModifyConversationField, o.Client, c)
}

func (o *ConversationApi) SetConversations(c *gin.Context) {
	a2r.Call(conversation.ConversationClient.SetConversations, o.Client, c)
}
