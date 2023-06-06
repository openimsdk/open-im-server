package msg

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/msg"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
)

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
	err = m.Conversation.SetHasReadSeq(ctx, req.UserID, req.ConversationID, conversations[0].ConversationType, req.Seqs[len(req.Seqs)-1])
	if err != nil {
		return
	}
	if err = m.sendMarkAsReadNotification(ctx, req.ConversationID, conversations[0].ConversationType, req.UserID, m.conversationAndGetRecvID(conversations[0], req.UserID), req.Seqs); err != nil {
		return
	}
	return &msg.MarkMsgsAsReadResp{}, nil
}

func (m *msgServer) sendMarkAsReadNotification(ctx context.Context, conversationID string, sesstionType int32, sendID, recvID string, seqs []int64) error {
	tips := &sdkws.MarkAsReadTips{
		MarkAsReadUserID: sendID,
		ConversationID:   conversationID,
		Seqs:             seqs,
	}
	m.notificationSender.NotificationWithSesstionType(ctx, sendID, recvID, constant.HasReadReceipt, sesstionType, tips)
	return nil
}
