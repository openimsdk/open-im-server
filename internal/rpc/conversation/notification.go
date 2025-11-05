package conversation

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/pkg/rpcli"
	"github.com/openimsdk/protocol/msg"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/notification"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/sdkws"
)

type ConversationNotificationSender struct {
	*notification.NotificationSender
}

func NewConversationNotificationSender(conf *config.Notification, msgClient *rpcli.MsgClient) *ConversationNotificationSender {
	return &ConversationNotificationSender{notification.NewNotificationSender(conf, notification.WithRpcClient(func(ctx context.Context, req *msg.SendMsgReq) (*msg.SendMsgResp, error) {
		return msgClient.SendMsg(ctx, req)
	}))}
}

// SetPrivate invote.
func (c *ConversationNotificationSender) ConversationSetPrivateNotification(ctx context.Context, sendID, recvID string,
	isPrivateChat bool, conversationID string,
) {
	tips := &sdkws.ConversationSetPrivateTips{
		RecvID:         recvID,
		SendID:         sendID,
		IsPrivate:      isPrivateChat,
		ConversationID: conversationID,
	}

	c.Notification(ctx, sendID, recvID, constant.ConversationPrivateChatNotification, tips)
}

func (c *ConversationNotificationSender) ConversationChangeNotification(ctx context.Context, userID string, conversationIDs []string) {
	tips := &sdkws.ConversationUpdateTips{
		UserID:             userID,
		ConversationIDList: conversationIDs,
	}

	c.Notification(ctx, userID, userID, constant.ConversationChangeNotification, tips)
}

func (c *ConversationNotificationSender) ConversationUnreadChangeNotification(
	ctx context.Context,
	userID, conversationID string,
	unreadCountTime, hasReadSeq int64,
) {
	tips := &sdkws.ConversationHasReadTips{
		UserID:          userID,
		ConversationID:  conversationID,
		HasReadSeq:      hasReadSeq,
		UnreadCountTime: unreadCountTime,
	}

	c.Notification(ctx, userID, userID, constant.ConversationUnreadNotification, tips)
}
