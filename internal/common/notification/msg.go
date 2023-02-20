package notification

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/proto/sdkws"
	"context"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

func (c *Check) DeleteMessageNotification(ctx context.Context, userID string, seqs []int64, operationID string) {
	DeleteMessageTips := sdkws.DeleteMessageTips{UserID: userID, Seqs: seqs}
	c.MessageNotification(ctx, userID, userID, constant.DeleteMessageNotification, &DeleteMessageTips)
}

func (c *Check) MessageNotification(ctx context.Context, sendID, recvID string, contentType int32, m proto.Message) {
	var err error
	var tips sdkws.TipsComm
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
	var n NotificationMsg
	n.SendID = sendID
	n.RecvID = recvID
	n.ContentType = contentType
	n.SessionType = constant.SingleChatType
	n.MsgFrom = constant.SysMsgType
	n.Content, err = proto.Marshal(&tips)
	if err != nil {
		return
	}
	c.Notification(ctx, &n)
}
