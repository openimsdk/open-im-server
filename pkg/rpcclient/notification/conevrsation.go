package notification

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/rpcclient"
)

type ConversationNotificationSender struct {
	*rpcclient.MsgClient
}

func NewConversationNotificationSender(client discoveryregistry.SvcDiscoveryRegistry) *ConversationNotificationSender {
	return &ConversationNotificationSender{rpcclient.NewMsgClient(client)}
}

// SetPrivate调用
func (c *ConversationNotificationSender) ConversationSetPrivateNotification(ctx context.Context, sendID, recvID string, isPrivateChat bool) error {
	tips := &sdkws.ConversationSetPrivateTips{
		RecvID:    recvID,
		SendID:    sendID,
		IsPrivate: isPrivateChat,
	}
	return c.Notification(ctx, sendID, recvID, constant.ConversationPrivateChatNotification, tips)
}

// 会话改变
func (c *ConversationNotificationSender) ConversationChangeNotification(ctx context.Context, userID string) error {
	tips := &sdkws.ConversationUpdateTips{
		UserID: userID,
	}
	return c.Notification(ctx, userID, userID, constant.ConversationChangeNotification, tips)
}

// 会话未读数同步
func (c *ConversationNotificationSender) ConversationUnreadChangeNotification(ctx context.Context, userID, conversationID string, updateUnreadCountTime int64) error {
	tips := &sdkws.ConversationUpdateTips{
		UserID:                userID,
		ConversationIDList:    []string{conversationID},
		UpdateUnreadCountTime: updateUnreadCountTime,
	}
	return c.Notification(ctx, userID, userID, constant.ConversationUnreadNotification, tips)
}
