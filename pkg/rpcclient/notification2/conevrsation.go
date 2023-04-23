package notification2

import (
	"context"
	"encoding/json"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/rpcclient"
	"github.com/golang/protobuf/proto"
)

type ConversationNotificationSender struct {
	*rpcclient.MsgClient
}

func NewConversationNotificationSender(client discoveryregistry.SvcDiscoveryRegistry) *ConversationNotificationSender {
	return &ConversationNotificationSender{rpcclient.NewMsgClient(client)}
}

func (c *ConversationNotificationSender) SetConversationNotification(ctx context.Context, sendID, recvID string, contentType int, m proto.Message) {
	var err error
	var n rpcclient.NotificationMsg
	n.SendID = sendID
	n.RecvID = recvID
	n.ContentType = int32(contentType)
	n.SessionType = constant.SingleChatType
	n.MsgFrom = constant.SysMsgType
	n.Content, err = json.Marshal(m)
	if err != nil {
		return
	}
	c.Notification(ctx, &n)
}

// SetPrivate调用
func (c *ConversationNotificationSender) ConversationSetPrivateNotification(ctx context.Context, sendID, recvID string, isPrivateChat bool) {
	tips := &sdkws.ConversationSetPrivateTips{
		RecvID:    recvID,
		SendID:    sendID,
		IsPrivate: isPrivateChat,
	}
	c.SetConversationNotification(ctx, sendID, recvID, constant.ConversationPrivateChatNotification, tips)
}

// 会话改变
func (c *ConversationNotificationSender) ConversationChangeNotification(ctx context.Context, userID string) {
	tips := &sdkws.ConversationUpdateTips{
		UserID: userID,
	}
	c.SetConversationNotification(ctx, userID, userID, constant.ConversationOptChangeNotification, tips)
}

// 会话未读数同步
func (c *ConversationNotificationSender) ConversationUnreadChangeNotification(ctx context.Context, userID, conversationID string, updateUnreadCountTime int64) {
	tips := &sdkws.ConversationUpdateTips{
		UserID:                userID,
		ConversationIDList:    []string{conversationID},
		UpdateUnreadCountTime: updateUnreadCountTime,
	}
	c.SetConversationNotification(ctx, userID, userID, constant.ConversationUnreadNotification, tips)
}
