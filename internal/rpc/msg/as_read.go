package msg

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/msg"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
)

func (m *msgServer) MarkMsgsAsRead(ctx context.Context, req *msg.MarkMsgsAsReadReq) (resp *msg.MarkMsgsAsReadResp, err error) {
	recvID, err := m.getConversationAndGetRecvID(ctx, req.ConversationID, req.UserID)
	if err != nil {
		return
	}
	err = m.MsgDatabase.MarkSingleChatMsgsAsRead(ctx, req.ConversationID, req.UserID, req.Seqs)
	if err != nil {
		return
	}
	if err = m.sendMarkAsReadNotification(ctx, req.ConversationID, req.UserID, recvID, req.Seqs); err != nil {
		return
	}
	return &msg.MarkMsgsAsReadResp{}, nil
}

func (m *msgServer) sendMarkAsReadNotification(ctx context.Context, conversationID string, sendID, recvID string, seqs []int64) error {
	tips := &sdkws.MarkAsReadTips{
		MarkAsReadUserID: sendID,
		ConversationID:   conversationID,
		Seqs:             seqs,
	}
	m.notificationSender.NotificationWithSesstionType(ctx, sendID, recvID, constant.HasReadReceipt, constant.SingleChatType, tips)
	return nil
}
