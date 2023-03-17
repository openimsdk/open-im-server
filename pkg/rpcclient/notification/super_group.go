package notification

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	//"github.com/golang/protobuf/jsonpb"
	//"github.com/golang/protobuf/proto"
)

func (c *Check) SuperGroupNotification(ctx context.Context, sendID, recvID string) {
	n := &NotificationMsg{
		SendID:      sendID,
		RecvID:      recvID,
		MsgFrom:     constant.SysMsgType,
		ContentType: constant.SuperGroupUpdateNotification,
		SessionType: constant.SingleChatType,
	}
	c.Notification(ctx, n)
}
