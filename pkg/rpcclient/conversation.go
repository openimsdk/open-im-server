package rpcclient

import (
	"context"
	"fmt"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	discoveryRegistry "github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	pbConversation "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/conversation"
)

type ConversationClient struct {
	*MetaClient
}

func NewConversationClient(zk discoveryRegistry.SvcDiscoveryRegistry) *ConversationClient {
	return &ConversationClient{NewMetaClient(zk, config.Config.RpcRegisterName.OpenImConversationName)}
}

func (c *ConversationClient) ModifyConversationField(ctx context.Context, req *pbConversation.ModifyConversationFieldReq) error {
	cc, err := c.getConn(ctx)
	if err != nil {
		return err
	}
	_, err = pbConversation.NewConversationClient(cc).ModifyConversationField(ctx, req)
	return err
}

func (c *ConversationClient) GetSingleConversationRecvMsgOpt(ctx context.Context, userID, conversationID string) (int32, error) {
	cc, err := c.getConn(ctx)
	if err != nil {
		return 0, err
	}
	var req pbConversation.GetConversationReq
	req.OwnerUserID = userID
	req.ConversationID = conversationID
	conversation, err := pbConversation.NewConversationClient(cc).GetConversation(ctx, &req)
	if err != nil {
		return 0, err
	}
	return conversation.GetConversation().RecvMsgOpt, err
}

func (c *ConversationClient) SingleChatFirstCreateConversation(ctx context.Context, recvID, sendID string) error {
	cc, err := c.getConn(ctx)
	if err != nil {
		return err
	}
	_, err = pbConversation.NewConversationClient(cc).CreateSingleChatConversations(ctx, &pbConversation.CreateSingleChatConversationsReq{RecvID: recvID, SendID: sendID})
	return err
}

func (c *ConversationClient) GroupChatFirstCreateConversation(ctx context.Context, groupID string, userIDs []string) error {
	cc, err := c.getConn(ctx)
	if err != nil {
		return err
	}
	_, err = pbConversation.NewConversationClient(cc).CreateGroupChatConversations(ctx, &pbConversation.CreateGroupChatConversationsReq{UserIDs: userIDs, GroupID: groupID})
	return err
}

func (c *ConversationClient) DelGroupChatConversations(ctx context.Context, ownerUserIDs []string, groupID string, maxSeq int64) error {
	cc, err := c.getConn(ctx)
	if err != nil {
		return err
	}
	_, err = pbConversation.NewConversationClient(cc).DelGroupChatConversations(ctx, &pbConversation.DelGroupChatConversationsReq{OwnerUserID: ownerUserIDs, GroupID: groupID, MaxSeq: maxSeq})
	return err
}

func (c *ConversationClient) GetConversationIDs(ctx context.Context, ownerUserID string) ([]string, error) {
	cc, err := c.getConn(ctx)
	if err != nil {
		return nil, err
	}
	resp, err := pbConversation.NewConversationClient(cc).GetConversationIDs(ctx, &pbConversation.GetConversationIDsReq{UserID: ownerUserID})
	if err != nil {
		return nil, err
	}
	return resp.ConversationIDs, nil
}

func (c *ConversationClient) GetConversation(ctx context.Context, ownerUserID, conversationID string) (*pbConversation.Conversation, error) {
	cc, err := c.getConn(ctx)
	if err != nil {
		return nil, err
	}
	resp, err := pbConversation.NewConversationClient(cc).GetConversation(ctx, &pbConversation.GetConversationReq{OwnerUserID: ownerUserID, ConversationID: conversationID})
	if err != nil {
		return nil, err
	}
	return resp.Conversation, nil
}

func (c *ConversationClient) GetConversationsByConversationID(ctx context.Context, conversationIDs []string) ([]*pbConversation.Conversation, error) {
	cc, err := c.getConn(ctx)
	if err != nil {
		return nil, err
	}
	resp, err := pbConversation.NewConversationClient(cc).GetConversationsByConversationID(ctx, &pbConversation.GetConversationsByConversationIDReq{ConversationIDs: conversationIDs})
	if err != nil {
		return nil, err
	}
	if len(resp.Conversations) == 0 {
		return nil, errs.ErrRecordNotFound.Wrap(fmt.Sprintf("conversationIDs: %v not found", conversationIDs))
	}
	return resp.Conversations, nil
}
