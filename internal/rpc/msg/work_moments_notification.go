package msg

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	pbOffice "Open_IM/pkg/proto/office"
	"Open_IM/pkg/utils"
	"github.com/golang/protobuf/proto"
)

func WorkMomentSendNotification(operationID, sendID, recvID string, notificationMsg *pbOffice.WorkMomentNotificationMsg) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), recvID)
	bytes, err := proto.Marshal(notificationMsg)
	if err != nil {
		log.NewError(operationID, utils.GetSelfFuncName(), "proto marshal failed", err.Error())
	}
	WorkMomentNotification(operationID, sendID, recvID, bytes)
}

func WorkMomentNotification(operationID, sendID, recvID string, content []byte) {
	n := &NotificationMsg{
		SendID:      sendID,
		RecvID:      recvID,
		Content:     content,
		MsgFrom:     constant.UserMsgType,
		ContentType: constant.WorkMomentNotification,
		SessionType: constant.UserMsgType,
		OperationID: operationID,
	}
	Notification(n)
}
