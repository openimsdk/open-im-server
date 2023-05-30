package rpcclient

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	discoveryRegistry "github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	pbConversation "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/conversation"
	"google.golang.org/grpc"
)

type ConversationClient struct {
	conn *grpc.ClientConn
}

func NewConversationClient(discov discoveryRegistry.SvcDiscoveryRegistry) *ConversationClient {
	conn, err := discov.GetConn(context.Background(), config.Config.RpcRegisterName.OpenImConversationName)
	if err != nil {
		panic(err)
	}
	return &ConversationClient{conn: conn}
}

func (c *ConversationClient) ModifyConversationField(ctx context.Context, req *pbConversation.ModifyConversationFieldReq) error {
	_, err := pbConversation.NewConversationClient(c.conn).ModifyConversationField(ctx, req)
	return err
}

func (c *ConversationClient) GetSingleConversationRecvMsgOpt(ctx context.Context, userID, conversationID string) (int32, error) {
	var req pbConversation.GetConversationReq
	req.OwnerUserID = userID
	req.ConversationID = conversationID
	conversation, err := pbConversation.NewConversationClient(c.conn).GetConversation(ctx, &req)
	if err != nil {
		return 0, err
	}
	return conversation.GetConversation().RecvMsgOpt, err
}

func (c *ConversationClient) SingleChatFirstCreateConversation(ctx context.Context, recvID, sendID string) error {
	_, err := pbConversation.NewConversationClient(c.conn).CreateSingleChatConversations(ctx, &pbConversation.CreateSingleChatConversationsReq{RecvID: recvID, SendID: sendID})
	return err
}

func (c *ConversationClient) GroupChatFirstCreateConversation(ctx context.Context, groupID string, userIDs []string) error {
	_, err := pbConversation.NewConversationClient(c.conn).CreateGroupChatConversations(ctx, &pbConversation.CreateGroupChatConversationsReq{UserIDs: userIDs, GroupID: groupID})
	return err
}

func (c *ConversationClient) DelGroupChatConversations(ctx context.Context, ownerUserIDs []string, groupID string, maxSeq int64) error {
	_, err := pbConversation.NewConversationClient(c.conn).DelGroupChatConversations(ctx, &pbConversation.DelGroupChatConversationsReq{OwnerUserID: ownerUserIDs, GroupID: groupID, MaxSeq: maxSeq})
	return err
}

func (c *ConversationClient) GetConversationIDs(ctx context.Context, ownerUserID string) ([]string, error) {
	resp, err := pbConversation.NewConversationClient(c.conn).GetConversationIDs(ctx, &pbConversation.GetConversationIDsReq{UserID: ownerUserID})
	return resp.ConversationIDs, err
}

func (c *ConversationClient) GetConversation(ctx context.Context, ownerUserID, conversationID string) (*pbConversation.Conversation, error) {
	resp, err := pbConversation.NewConversationClient(c.conn).GetConversation(ctx, &pbConversation.GetConversationReq{OwnerUserID: ownerUserID, ConversationID: conversationID})
	return resp.Conversation, err
}

func (c *ConversationClient) GetConversationByConversationID(ctx context.Context, conversationID string) (*pbConversation.Conversation, error) {
	resp, err := pbConversation.NewConversationClient(c.conn).GetConversationByConversationID(ctx, &pbConversation.GetConversationByConversationIDReq{ConversationID: conversationID})
	return resp.Conversation, err
}
