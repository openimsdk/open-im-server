package msg

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	//sdk "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	//"github.com/golang/protobuf/jsonpb"
	//"github.com/golang/protobuf/proto"
)

func SuperGroupNotification(operationID, sendID, recvID string, isKicked bool) {
	m := make(map[string]bool)
	m["kicked"] = isKicked
	n := &NotificationMsg{
		SendID:      sendID,
		RecvID:      recvID,
		MsgFrom:     constant.SysMsgType,
		ContentType: constant.SuperGroupUpdateNotification,
		SessionType: constant.SingleChatType,
		OperationID: operationID,
	}
	n.Content = utils.StructToJsonBytes(m)

	log.NewInfo(operationID, utils.GetSelfFuncName(), string(n.Content))
	Notification(n)
}
