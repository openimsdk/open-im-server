package msg

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/msg"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
)

func (m *msgServer) GetConversationsHasReadAndMaxSeq(ctx context.Context, req *msg.GetConversationsHasReadAndMaxSeqReq) (*msg.GetConversationsHasReadAndMaxSeqResp, error) {
	conversationIDs, err := m.ConversationLocalCache.GetConversationIDs(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	hasReadSeqs, err := m.MsgDatabase.GetHasReadSeqs(ctx, req.UserID, conversationIDs)
	if err != nil {
		return nil, err
	}
	maxSeqs, err := m.MsgDatabase.GetMaxSeqs(ctx, conversationIDs)
	if err != nil {
		return nil, err
	}
	resp := &msg.GetConversationsHasReadAndMaxSeqResp{Seqs: make(map[string]*msg.Seqs)}
	for conversarionID, maxSeq := range maxSeqs {
		resp.Seqs[conversarionID] = &msg.Seqs{
			HasReadSeq: hasReadSeqs[conversarionID],
			MaxSeq:     maxSeq,
		}
	}
	return resp, nil
}

func (m *msgServer) MarkMsgsAsRead(ctx context.Context, req *msg.MarkMsgsAsReadReq) (resp *msg.MarkMsgsAsReadResp, err error) {
	if len(req.Seqs) < 1 {
		return nil, errs.ErrArgs.Wrap("seqs must not be empty")
	}
	conversations, err := m.Conversation.GetConversationsByConversationID(ctx, []string{req.ConversationID})
	if err != nil {
		return
	}
	err = m.MsgDatabase.MarkSingleChatMsgsAsRead(ctx, req.ConversationID, req.UserID, req.Seqs)
	if err != nil {
		return
	}
	hasReadSeq := req.Seqs[len(req.Seqs)-1]
	err = m.MsgDatabase.SetHasReadSeq(ctx, req.UserID, req.ConversationID, hasReadSeq)
	if err != nil {
		return
	}
	if err = m.sendMarkAsReadNotification(ctx, req.ConversationID, conversations[0].ConversationType, req.UserID, m.conversationAndGetRecvID(conversations[0], req.UserID), req.Seqs, hasReadSeq); err != nil {
		return
	}
	return &msg.MarkMsgsAsReadResp{}, nil
}

func (m *msgServer) sendMarkAsReadNotification(ctx context.Context, conversationID string, sesstionType int32, sendID, recvID string, seqs []int64, hasReadSeq int64) error {
	tips := &sdkws.MarkAsReadTips{
		MarkAsReadUserID: sendID,
		ConversationID:   conversationID,
		Seqs:             seqs,
		HasReadSeq:       hasReadSeq,
	}
	m.notificationSender.NotificationWithSesstionType(ctx, sendID, recvID, constant.HasReadReceipt, sesstionType, tips)
	return nil
}
