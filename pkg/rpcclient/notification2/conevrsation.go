package notification2

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/rpcclient"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

type ConversationNotificationSender struct {
	*rpcclient.MsgClient
}

func NewConversationNotificationSender(client discoveryregistry.SvcDiscoveryRegistry) *ConversationNotificationSender {
	return &ConversationNotificationSender{rpcclient.NewMsgClient(client)}
}

func (c *ConversationNotificationSender) SetConversationNotification(ctx context.Context, sendID, recvID string, contentType int, m proto.Message, tips *sdkws.TipsComm) {
	var err error
	tips.Detail, err = proto.Marshal(m)
	if err != nil {
		return
	}
	marshaler := jsonpb.Marshaler{
		OrigName:     true,
		EnumsAsInts:  false,
		EmitDefaults: false,
	}
	tips.JsonDetail, _ = marshaler.MarshalToString(m)
	var n rpcclient.NotificationMsg
	n.SendID = sendID
	n.RecvID = recvID
	n.ContentType = int32(contentType)
	n.SessionType = constant.SingleChatType
	n.MsgFrom = constant.SysMsgType
	n.Content, err = proto.Marshal(tips)
	if err != nil {
		return
	}
	c.Notification(ctx, &n)
}

// SetPrivate调用
func (c *ConversationNotificationSender) ConversationSetPrivateNotification(ctx context.Context, sendID, recvID string, isPrivateChat bool) {
	conversationSetPrivateTips := &sdkws.ConversationSetPrivateTips{
		RecvID:    recvID,
		SendID:    sendID,
		IsPrivate: isPrivateChat,
	}
	var tips sdkws.TipsComm
	var tipsMsg string
	if isPrivateChat == true {
		tipsMsg = config.Config.Notification.ConversationSetPrivate.DefaultTips.OpenTips
	} else {
		tipsMsg = config.Config.Notification.ConversationSetPrivate.DefaultTips.CloseTips
	}
	tips.DefaultTips = tipsMsg
	c.SetConversationNotification(ctx, sendID, recvID, constant.ConversationPrivateChatNotification, conversationSetPrivateTips, &tips)
}

// 会话改变
func (c *ConversationNotificationSender) ConversationChangeNotification(ctx context.Context, userID string) {
	ConversationChangedTips := &sdkws.ConversationUpdateTips{
		UserID: userID,
	}
	var tips sdkws.TipsComm
	tips.DefaultTips = config.Config.Notification.ConversationOptUpdate.DefaultTips.Tips
	c.SetConversationNotification(ctx, userID, userID, constant.ConversationOptChangeNotification, ConversationChangedTips, &tips)
}

// 会话未读数同步
func (c *ConversationNotificationSender) ConversationUnreadChangeNotification(ctx context.Context, userID, conversationID string, updateUnreadCountTime int64) {
	ConversationChangedTips := &sdkws.ConversationUpdateTips{
		UserID:                userID,
		ConversationIDList:    []string{conversationID},
		UpdateUnreadCountTime: updateUnreadCountTime,
	}
	var tips sdkws.TipsComm
	tips.DefaultTips = config.Config.Notification.ConversationOptUpdate.DefaultTips.Tips
	c.SetConversationNotification(ctx, userID, userID, constant.ConversationUnreadNotification, ConversationChangedTips, &tips)
}
