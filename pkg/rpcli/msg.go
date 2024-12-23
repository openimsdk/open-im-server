package rpcli

import (
	"context"
	"github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/protocol/sdkws"
	"google.golang.org/grpc"
)

func NewMsgClient(cc grpc.ClientConnInterface) *MsgClient {
	return &MsgClient{msg.NewMsgClient(cc)}
}

type MsgClient struct {
	msg.MsgClient
}

func (x *MsgClient) GetMaxSeqs(ctx context.Context, conversationIDs []string) (map[string]int64, error) {
	if len(conversationIDs) == 0 {
		return nil, nil
	}
	req := &msg.GetMaxSeqsReq{ConversationIDs: conversationIDs}
	return extractField(ctx, x.MsgClient.GetMaxSeqs, req, (*msg.SeqsInfoResp).GetMaxSeqs)
}

func (x *MsgClient) GetMsgByConversationIDs(ctx context.Context, conversationIDs []string, maxSeqs map[string]int64) (map[string]*sdkws.MsgData, error) {
	if len(conversationIDs) == 0 || len(maxSeqs) == 0 {
		return nil, nil
	}
	req := &msg.GetMsgByConversationIDsReq{ConversationIDs: conversationIDs, MaxSeqs: maxSeqs}
	return extractField(ctx, x.MsgClient.GetMsgByConversationIDs, req, (*msg.GetMsgByConversationIDsResp).GetMsgDatas)
}

func (x *MsgClient) GetHasReadSeqs(ctx context.Context, conversationIDs []string, userID string) (map[string]int64, error) {
	if len(conversationIDs) == 0 {
		return nil, nil
	}
	req := &msg.GetHasReadSeqsReq{ConversationIDs: conversationIDs, UserID: userID}
	return extractField(ctx, x.MsgClient.GetHasReadSeqs, req, (*msg.SeqsInfoResp).GetMaxSeqs)
}

func (x *MsgClient) SetUserConversationMaxSeq(ctx context.Context, conversationID string, ownerUserIDs []string, maxSeq int64) error {
	if len(ownerUserIDs) == 0 {
		return nil
	}
	req := &msg.SetUserConversationMaxSeqReq{ConversationID: conversationID, OwnerUserID: ownerUserIDs, MaxSeq: maxSeq}
	return ignoreResp(x.MsgClient.SetUserConversationMaxSeq(ctx, req))
}

func (x *MsgClient) SetUserConversationMin(ctx context.Context, conversationID string, ownerUserIDs []string, minSeq int64) error {
	if len(ownerUserIDs) == 0 {
		return nil
	}
	req := &msg.SetUserConversationsMinSeqReq{ConversationID: conversationID, UserIDs: ownerUserIDs, Seq: minSeq}
	return ignoreResp(x.MsgClient.SetUserConversationsMinSeq(ctx, req))
}

func (x *MsgClient) GetLastMessageSeqByTime(ctx context.Context, conversationID string, lastTime int64) (int64, error) {
	req := &msg.GetLastMessageSeqByTimeReq{ConversationID: conversationID, Time: lastTime}
	return extractField(ctx, x.MsgClient.GetLastMessageSeqByTime, req, (*msg.GetLastMessageSeqByTimeResp).GetSeq)
}

func (x *MsgClient) GetConversationMaxSeq(ctx context.Context, conversationID string) (int64, error) {
	req := &msg.GetConversationMaxSeqReq{ConversationID: conversationID}
	return extractField(ctx, x.MsgClient.GetConversationMaxSeq, req, (*msg.GetConversationMaxSeqResp).GetMaxSeq)
}

func (x *MsgClient) GetActiveConversation(ctx context.Context, conversationIDs []string) ([]*msg.ActiveConversation, error) {
	if len(conversationIDs) == 0 {
		return nil, nil
	}
	req := &msg.GetActiveConversationReq{ConversationIDs: conversationIDs}
	return extractField(ctx, x.MsgClient.GetActiveConversation, req, (*msg.GetActiveConversationResp).GetConversations)
}

func (x *MsgClient) GetSeqMessage(ctx context.Context, userID string, conversations []*msg.ConversationSeqs) (map[string]*sdkws.PullMsgs, error) {
	if len(conversations) == 0 {
		return nil, nil
	}
	req := &msg.GetSeqMessageReq{UserID: userID, Conversations: conversations}
	return extractField(ctx, x.MsgClient.GetSeqMessage, req, (*msg.GetSeqMessageResp).GetMsgs)
}

func (x *MsgClient) SetUserConversationsMinSeq(ctx context.Context, conversationID string, userIDs []string, seq int64) error {
	if len(userIDs) == 0 {
		return nil
	}
	req := &msg.SetUserConversationsMinSeqReq{ConversationID: conversationID, UserIDs: userIDs, Seq: seq}
	return ignoreResp(x.MsgClient.SetUserConversationsMinSeq(ctx, req))
}
