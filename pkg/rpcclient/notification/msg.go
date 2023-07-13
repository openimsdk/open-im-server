package notification

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/rpcclient"
)

type MsgNotificationSender struct {
	*rpcclient.NotificationSender
}

func NewMsgNotificationSender(opts ...rpcclient.NotificationSenderOptions) *MsgNotificationSender {
	return &MsgNotificationSender{rpcclient.NewNotificationSender(opts...)}
}

func (m *MsgNotificationSender) UserDeleteMsgsNotification(ctx context.Context, userID, conversationID string, seqs []int64) error {
	tips := sdkws.DeleteMsgsTips{
		UserID:         userID,
		ConversationID: conversationID,
		Seqs:           seqs,
	}
	return m.Notification(ctx, userID, userID, constant.MsgDeleteNotification, &tips)
}

func (m *MsgNotificationSender) MarkAsReadNotification(ctx context.Context, conversationID string, sesstionType int32, sendID, recvID string, seqs []int64, hasReadSeq int64) error {
	tips := &sdkws.MarkAsReadTips{
		MarkAsReadUserID: sendID,
		ConversationID:   conversationID,
		Seqs:             seqs,
		HasReadSeq:       hasReadSeq,
	}
	return m.NotificationWithSesstionType(ctx, sendID, recvID, constant.HasReadReceipt, sesstionType, tips)
}
