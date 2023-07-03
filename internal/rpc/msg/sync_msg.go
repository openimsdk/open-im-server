package msg

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/tokenverify"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
)

func (m *msgServer) PullMessageBySeqs(
	ctx context.Context,
	req *sdkws.PullMessageBySeqsReq,
) (*sdkws.PullMessageBySeqsResp, error) {
	resp := &sdkws.PullMessageBySeqsResp{}
	resp.Msgs = make(map[string]*sdkws.PullMsgs)
	resp.NotificationMsgs = make(map[string]*sdkws.PullMsgs)
	for _, seq := range req.SeqRanges {
		if !utils.IsNotification(seq.ConversationID) {
			conversation, err := m.Conversation.GetConversation(ctx, req.UserID, seq.ConversationID)
			if err != nil {
				log.ZError(ctx, "GetConversation error", err, "conversationID", seq.ConversationID)
				continue
			}
			minSeq, maxSeq, msgs, err := m.MsgDatabase.GetMsgBySeqsRange(
				ctx,
				req.UserID,
				seq.ConversationID,
				seq.Begin,
				seq.End,
				seq.Num,
				conversation.MaxSeq,
			)
			if err != nil {
				log.ZWarn(ctx, "GetMsgBySeqsRange error", err, "conversationID", seq.ConversationID, "seq", seq)
				continue
			}
			var isEnd bool
			switch req.Order {
			case sdkws.PullOrder_PullOrderAsc:
				isEnd = maxSeq <= seq.End
			case sdkws.PullOrder_PullOrderDesc:
				isEnd = seq.Begin <= minSeq
			}
			resp.Msgs[seq.ConversationID] = &sdkws.PullMsgs{Msgs: msgs, IsEnd: isEnd}
		} else {
			var seqs []int64
			for i := seq.Begin; i <= seq.End; i++ {
				seqs = append(seqs, i)
			}
			minSeq, maxSeq, notificationMsgs, err := m.MsgDatabase.GetMsgBySeqs(ctx, req.UserID, seq.ConversationID, seqs)
			if err != nil {
				log.ZWarn(ctx, "GetMsgBySeqs error", err, "conversationID", seq.ConversationID, "seq", seq)
				continue
			}
			var isEnd bool
			switch req.Order {
			case sdkws.PullOrder_PullOrderAsc:
				isEnd = maxSeq <= seq.End
			case sdkws.PullOrder_PullOrderDesc:
				isEnd = seq.Begin <= minSeq
			}
			resp.NotificationMsgs[seq.ConversationID] = &sdkws.PullMsgs{Msgs: notificationMsgs, IsEnd: isEnd}
		}
	}
	return resp, nil
}

func (m *msgServer) GetMaxSeq(ctx context.Context, req *sdkws.GetMaxSeqReq) (*sdkws.GetMaxSeqResp, error) {
	if err := tokenverify.CheckAccessV3(ctx, req.UserID); err != nil {
		return nil, err
	}
	conversationIDs, err := m.ConversationLocalCache.GetConversationIDs(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	for _, conversationID := range conversationIDs {
		conversationIDs = append(conversationIDs, utils.GetNotificationConversationIDByConversationID(conversationID))
	}
	log.ZDebug(ctx, "GetMaxSeq", "conversationIDs", conversationIDs)
	maxSeqs, err := m.MsgDatabase.GetMaxSeqs(ctx, conversationIDs)
	if err != nil {
		log.ZWarn(ctx, "GetMaxSeqs error", err, "conversationIDs", conversationIDs, "maxSeqs", maxSeqs)
		return nil, err
	}
	resp := new(sdkws.GetMaxSeqResp)
	resp.MaxSeqs = maxSeqs
	return resp, nil
}
