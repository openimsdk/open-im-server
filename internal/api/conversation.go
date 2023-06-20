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
		panic(err)
	}
	client := conversation.NewConversationClient(conn)
	return &Conversation{discov: discov, conn: conn, client: client}
}

type Conversation struct {
	client conversation.ConversationClient
	conn   *grpc.ClientConn
	discov discoveryregistry.SvcDiscoveryRegistry
}

func (o *Conversation) Client() conversation.ConversationClient {
	return o.client
}

func (o *Conversation) GetAllConversations(c *gin.Context) {
	a2r.Call(conversation.ConversationClient.GetAllConversations, o.Client, c)
}

func (o *Conversation) GetConversation(c *gin.Context) {
	a2r.Call(conversation.ConversationClient.GetConversation, o.Client, c)
}

func (o *Conversation) GetConversations(c *gin.Context) {
	a2r.Call(conversation.ConversationClient.GetConversations, o.Client, c)
}

// deprecated
func (o *Conversation) SetConversation(c *gin.Context) {
	a2r.Call(conversation.ConversationClient.SetConversation, o.Client, c)
}

// deprecated
func (o *Conversation) BatchSetConversations(c *gin.Context) {
	a2r.Call(conversation.ConversationClient.BatchSetConversations, o.Client, c)
}

func (o *Conversation) SetRecvMsgOpt(c *gin.Context) {
	a2r.Call(conversation.ConversationClient.SetRecvMsgOpt, o.Client, c)
}

func (o *Conversation) ModifyConversationField(c *gin.Context) {
	a2r.Call(conversation.ConversationClient.ModifyConversationField, o.Client, c)
}

func (o *Conversation) SetConversations(c *gin.Context) {
	a2r.Call(conversation.ConversationClient.SetConversations, o.Client, c)
}
