package msg

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	pbOffice "Open_IM/pkg/proto/office"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

func WorkMomentSendNotification(operationID, sendID, recvID string, notificationMsg *pbOffice.WorkMomentNotificationMsg) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), sendID, recvID, notificationMsg)
	WorkMomentNotification(operationID, sendID, recvID, notificationMsg)
}

func WorkMomentNotification(operationID, sendID, recvID string, m proto.Message) {
	marshaler := jsonpb.Marshaler{
		OrigName:     true,
		EnumsAsInts:  false,
		EmitDefaults: false,
	}
	var tips open_im_sdk.TipsComm
	var err error
	tips.JsonDetail, _ = marshaler.MarshalToString(m)
	n := &NotificationMsg{
		SendID:      sendID,
		RecvID:      recvID,
		MsgFrom:     constant.UserMsgType,
		ContentType: constant.WorkMomentNotification,
		SessionType: constant.SingleChatType,
		OperationID: operationID,
	}
	n.Content, err = proto.Marshal(m)
	if err != nil {
		log.NewError(operationID, utils.GetSelfFuncName(), "proto marshal falied", err.Error())
		return
	}
	Notification(n)
}
