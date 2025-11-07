package msg

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/pkg/notification"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/sdkws"
)

type MsgNotificationSender struct {
	*notification.NotificationSender
}

func NewMsgNotificationSender(config *Config, opts ...notification.NotificationSenderOptions) *MsgNotificationSender {
	return &MsgNotificationSender{notification.NewNotificationSender(&config.NotificationConfig, opts...)}
}

func (m *MsgNotificationSender) UserDeleteMsgsNotification(ctx context.Context, userID, conversationID string, seqs []int64) {
	tips := sdkws.DeleteMsgsTips{
		UserID:         userID,
		ConversationID: conversationID,
		Seqs:           seqs,
	}
	m.Notification(ctx, userID, userID, constant.DeleteMsgsNotification, &tips)
}

func (m *MsgNotificationSender) MarkAsReadNotification(ctx context.Context, conversationID string, sessionType int32, sendID, recvID string, seqs []int64, hasReadSeq int64) {
	tips := &sdkws.MarkAsReadTips{
		MarkAsReadUserID: sendID,
		ConversationID:   conversationID,
		Seqs:             seqs,
		HasReadSeq:       hasReadSeq,
	}
	m.NotificationWithSessionType(ctx, sendID, recvID, constant.HasReadReceipt, sessionType, tips)
}
