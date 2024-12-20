package rpcli

import (
	"context"
	"github.com/openimsdk/protocol/conversation"
)

func NewConversationClient(cli conversation.ConversationClient) *ConversationClient {
	return &ConversationClient{cli}
}

type ConversationClient struct {
	conversation.ConversationClient
}

func (x *ConversationClient) SetConversationMaxSeq(ctx context.Context, conversationID string, ownerUserIDs []string, maxSeq int64) error {
	req := &conversation.SetConversationMaxSeqReq{ConversationID: conversationID, OwnerUserID: ownerUserIDs, MaxSeq: maxSeq}
	return ignoreResp(x.ConversationClient.SetConversationMaxSeq(ctx, req))
}

func (x *ConversationClient) SetConversations(ctx context.Context, userIDs []string, info *conversation.ConversationReq) error {
	req := &conversation.SetConversationsReq{UserIDs: userIDs, Conversation: info}
	return ignoreResp(x.ConversationClient.SetConversations(ctx, req))
}

func (x *ConversationClient) GetConversationsByConversationIDs(ctx context.Context, conversationIDs []string) ([]*conversation.Conversation, error) {
	req := &conversation.GetConversationsByConversationIDReq{ConversationIDs: conversationIDs}
	return extractField(ctx, x.ConversationClient.GetConversationsByConversationID, req, (*conversation.GetConversationsByConversationIDResp).GetConversations)
}

func (x *ConversationClient) GetConversationsByConversationID(ctx context.Context, conversationID string) (*conversation.Conversation, error) {
	return firstValue(x.GetConversationsByConversationIDs(ctx, []string{conversationID}))
}

func (x *ConversationClient) SetConversationMinSeq(ctx context.Context, conversationID string, ownerUserIDs []string, minSeq int64) error {
	req := &conversation.SetConversationMinSeqReq{ConversationID: conversationID, OwnerUserID: ownerUserIDs, MinSeq: minSeq}
	return ignoreResp(x.ConversationClient.SetConversationMinSeq(ctx, req))
}

func (x *ConversationClient) GetConversation(ctx context.Context, conversationID string, ownerUserID string) (*conversation.Conversation, error) {
	req := &conversation.GetConversationReq{ConversationID: conversationID, OwnerUserID: ownerUserID}
	return extractField(ctx, x.ConversationClient.GetConversation, req, (*conversation.GetConversationResp).GetConversation)
}

func (x *ConversationClient) GetConversations(ctx context.Context, conversationIDs []string, ownerUserID string) ([]*conversation.Conversation, error) {
	req := &conversation.GetConversationsReq{ConversationIDs: conversationIDs, OwnerUserID: ownerUserID}
	return extractField(ctx, x.ConversationClient.GetConversations, req, (*conversation.GetConversationsResp).GetConversations)
}
