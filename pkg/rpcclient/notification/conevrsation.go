package notification

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
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
func (c *ConversationNotificationSender) ConversationSetPrivateNotification(ctx context.Context, sendID, recvID string, isPrivateChat bool) {
	tips := &sdkws.ConversationSetPrivateTips{
		RecvID:    recvID,
		SendID:    sendID,
		IsPrivate: isPrivateChat,
	}
	c.Notification(ctx, sendID, recvID, constant.ConversationPrivateChatNotification, constant.SingleChatType, tips, config.Config.Notification.ConversationSetPrivate)
}

// 会话改变
func (c *ConversationNotificationSender) ConversationChangeNotification(ctx context.Context, userID string) {
	tips := &sdkws.ConversationUpdateTips{
		UserID: userID,
	}
	c.Notification(ctx, userID, userID, constant.ConversationChangeNotification, constant.SingleChatType, tips, config.Config.Notification.ConversationChanged)
}

// 会话未读数同步
func (c *ConversationNotificationSender) ConversationUnreadChangeNotification(ctx context.Context, userID, conversationID string, updateUnreadCountTime int64) {
	tips := &sdkws.ConversationUpdateTips{
		UserID:                userID,
		ConversationIDList:    []string{conversationID},
		UpdateUnreadCountTime: updateUnreadCountTime,
	}
	c.Notification(ctx, userID, userID, constant.ConversationUnreadNotification, constant.SingleChatType, tips, config.Config.Notification.ConversationChanged)
}
