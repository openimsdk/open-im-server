package rpcli

import (
	"context"
	"github.com/openimsdk/protocol/conversation"
	"google.golang.org/grpc"
)

func NewConversationClient(cc grpc.ClientConnInterface) *ConversationClient {
	return &ConversationClient{conversation.NewConversationClient(cc)}
}

type ConversationClient struct {
	conversation.ConversationClient
}

func (x *ConversationClient) SetConversationMaxSeq(ctx context.Context, conversationID string, ownerUserIDs []string, maxSeq int64) error {
	if len(ownerUserIDs) == 0 {
		return nil
	}
	req := &conversation.SetConversationMaxSeqReq{ConversationID: conversationID, OwnerUserID: ownerUserIDs, MaxSeq: maxSeq}
	return ignoreResp(x.ConversationClient.SetConversationMaxSeq(ctx, req))
}

func (x *ConversationClient) SetConversations(ctx context.Context, ownerUserIDs []string, info *conversation.ConversationReq) error {
	if len(ownerUserIDs) == 0 {
		return nil
	}
	req := &conversation.SetConversationsReq{UserIDs: ownerUserIDs, Conversation: info}
	return ignoreResp(x.ConversationClient.SetConversations(ctx, req))
}

func (x *ConversationClient) GetConversationsByConversationIDs(ctx context.Context, conversationIDs []string) ([]*conversation.Conversation, error) {
	if len(conversationIDs) == 0 {
		return nil, nil
	}
	req := &conversation.GetConversationsByConversationIDReq{ConversationIDs: conversationIDs}
	return extractField(ctx, x.ConversationClient.GetConversationsByConversationID, req, (*conversation.GetConversationsByConversationIDResp).GetConversations)
}

func (x *ConversationClient) GetConversationsByConversationID(ctx context.Context, conversationID string) (*conversation.Conversation, error) {
	return firstValue(x.GetConversationsByConversationIDs(ctx, []string{conversationID}))
}

func (x *ConversationClient) SetConversationMinSeq(ctx context.Context, conversationID string, ownerUserIDs []string, minSeq int64) error {
	if len(ownerUserIDs) == 0 {
		return nil
	}
	req := &conversation.SetConversationMinSeqReq{ConversationID: conversationID, OwnerUserID: ownerUserIDs, MinSeq: minSeq}
	return ignoreResp(x.ConversationClient.SetConversationMinSeq(ctx, req))
}

func (x *ConversationClient) GetConversation(ctx context.Context, conversationID string, ownerUserID string) (*conversation.Conversation, error) {
	req := &conversation.GetConversationReq{ConversationID: conversationID, OwnerUserID: ownerUserID}
	return extractField(ctx, x.ConversationClient.GetConversation, req, (*conversation.GetConversationResp).GetConversation)
}

func (x *ConversationClient) GetConversations(ctx context.Context, conversationIDs []string, ownerUserID string) ([]*conversation.Conversation, error) {
	if len(conversationIDs) == 0 {
		return nil, nil
	}
	req := &conversation.GetConversationsReq{ConversationIDs: conversationIDs, OwnerUserID: ownerUserID}
	return extractField(ctx, x.ConversationClient.GetConversations, req, (*conversation.GetConversationsResp).GetConversations)
}

func (x *ConversationClient) GetConversationIDs(ctx context.Context, ownerUserID string) ([]string, error) {
	req := &conversation.GetConversationIDsReq{UserID: ownerUserID}
	return extractField(ctx, x.ConversationClient.GetConversationIDs, req, (*conversation.GetConversationIDsResp).GetConversationIDs)
}

func (x *ConversationClient) GetPinnedConversationIDs(ctx context.Context, ownerUserID string) ([]string, error) {
	req := &conversation.GetPinnedConversationIDsReq{UserID: ownerUserID}
	return extractField(ctx, x.ConversationClient.GetPinnedConversationIDs, req, (*conversation.GetPinnedConversationIDsResp).GetConversationIDs)
}

func (x *ConversationClient) CreateGroupChatConversations(ctx context.Context, groupID string, userIDs []string) error {
	if len(userIDs) == 0 {
		return nil
	}
	req := &conversation.CreateGroupChatConversationsReq{GroupID: groupID, UserIDs: userIDs}
	return ignoreResp(x.ConversationClient.CreateGroupChatConversations(ctx, req))
}

func (x *ConversationClient) CreateSingleChatConversations(ctx context.Context, req *conversation.CreateSingleChatConversationsReq) error {
	return ignoreResp(x.ConversationClient.CreateSingleChatConversations(ctx, req))
}

func (x *ConversationClient) GetConversationOfflinePushUserIDs(ctx context.Context, conversationID string, userIDs []string) ([]string, error) {
	if len(userIDs) == 0 {
		return nil, nil
	}
	req := &conversation.GetConversationOfflinePushUserIDsReq{ConversationID: conversationID, UserIDs: userIDs}
	return extractField(ctx, x.ConversationClient.GetConversationOfflinePushUserIDs, req, (*conversation.GetConversationOfflinePushUserIDsResp).GetUserIDs)
}
