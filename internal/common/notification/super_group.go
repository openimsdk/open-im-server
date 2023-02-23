package notification

import (
	"OpenIM/pkg/common/constant"
	"context"
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
