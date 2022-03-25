package msg

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/utils"
)

func SetTagNotification(operationID, sendID, recvID, content string, contentType int32) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "args: ", sendID, recvID, content, contentType)
	var n NotificationMsg
	n.SendID = sendID
	n.RecvID = recvID
	n.ContentType = contentType
	n.SessionType = constant.SingleChatType
	n.MsgFrom = constant.UserMsgType
	n.OperationID = operationID
	n.Content = []byte(content)
	Notification(&n)
}
