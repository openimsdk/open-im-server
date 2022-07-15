package msg

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	//sdk "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	//"github.com/golang/protobuf/jsonpb"
	//"github.com/golang/protobuf/proto"
)

func SuperGroupNotification(operationID, sendID, recvID string) {

	//var tips sdk.TipsComm
	//var err error
	//marshaler := jsonpb.Marshaler{
	//	OrigName:     true,
	//	EnumsAsInts:  false,
	//	EmitDefaults: false,
	//}
	//tips.JsonDetail, _ = marshaler.MarshalToString(m)
	n := &NotificationMsg{
		SendID:      sendID,
		RecvID:      recvID,
		MsgFrom:     constant.SysMsgType,
		ContentType: constant.SuperGroupUpdateNotification,
		SessionType: constant.SingleChatType,
		OperationID: operationID,
	}
	//n.Content, err = proto.Marshal(&tips)
	//if err != nil {
	//	log.NewError(operationID, utils.GetSelfFuncName(), "proto.Marshal failed")
	//	return
	//}
	log.NewInfo(operationID, utils.GetSelfFuncName(), string(n.Content))
	Notification(n)
}
