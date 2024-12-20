package rpcli

import (
	"context"
	"github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/protocol/sdkws"
)

func NewMsgClient(cli msg.MsgClient) *MsgClient {
	return &MsgClient{cli}
}

type MsgClient struct {
	msg.MsgClient
}

func (x *MsgClient) cli() msg.MsgClient {
	return x.MsgClient
}

func (x *MsgClient) GetMaxSeqs(ctx context.Context, conversationIDs []string) (map[string]int64, error) {
	req := &msg.GetMaxSeqsReq{ConversationIDs: conversationIDs}
	return extractField(ctx, x.cli().GetMaxSeqs, req, (*msg.SeqsInfoResp).GetMaxSeqs)
}

func (x *MsgClient) GetMsgByConversationIDs(ctx context.Context, conversationIDs []string, maxSeqs map[string]int64) (map[string]*sdkws.MsgData, error) {
	req := &msg.GetMsgByConversationIDsReq{ConversationIDs: conversationIDs, MaxSeqs: maxSeqs}
	return extractField(ctx, x.cli().GetMsgByConversationIDs, req, (*msg.GetMsgByConversationIDsResp).GetMsgDatas)
}

func (x *MsgClient) GetHasReadSeqs(ctx context.Context, conversationIDs []string, userID string) (map[string]int64, error) {
	req := &msg.GetHasReadSeqsReq{ConversationIDs: conversationIDs, UserID: userID}
	return extractField(ctx, x.cli().GetHasReadSeqs, req, (*msg.SeqsInfoResp).GetMaxSeqs)
}

func (x *MsgClient) SetUserConversationMaxSeq(ctx context.Context, conversationID string, ownerUserIDs []string, maxSeq int64) error {
	req := &msg.SetUserConversationMaxSeqReq{ConversationID: conversationID, OwnerUserID: ownerUserIDs, MaxSeq: maxSeq}
	return ignoreResp(x.cli().SetUserConversationMaxSeq(ctx, req))
}

func (x *MsgClient) SetUserConversationMin(ctx context.Context, conversationID string, ownerUserIDs []string, minSeq int64) error {
	req := &msg.SetUserConversationsMinSeqReq{ConversationID: conversationID, UserIDs: ownerUserIDs, Seq: minSeq}
	return ignoreResp(x.cli().SetUserConversationsMinSeq(ctx, req))
}

func (x *MsgClient) GetLastMessageSeqByTime(ctx context.Context, conversationID string, lastTime int64) (int64, error) {
	req := &msg.GetLastMessageSeqByTimeReq{ConversationID: conversationID, Time: lastTime}
	return extractField(ctx, x.cli().GetLastMessageSeqByTime, req, (*msg.GetLastMessageSeqByTimeResp).GetSeq)
}

func (x *MsgClient) GetConversationMaxSeq(ctx context.Context, conversationID string) (int64, error) {
	req := &msg.GetConversationMaxSeqReq{ConversationID: conversationID}
	return extractField(ctx, x.cli().GetConversationMaxSeq, req, (*msg.GetConversationMaxSeqResp).GetMaxSeq)
}
