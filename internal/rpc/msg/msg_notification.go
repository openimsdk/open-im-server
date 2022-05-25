package msg

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

func DeleteMessageNotification(opUserID, userID string, seqList []uint32, operationID string) {
	DeleteMessageTips := open_im_sdk.DeleteMessageTips{OpUserID: opUserID, UserID: userID, SeqList: seqList}
	MessageNotification(operationID, userID, userID, constant.DeleteMessageNotification, &DeleteMessageTips)
}

func MessageNotification(operationID, sendID, recvID string, contentType int32, m proto.Message) {
	log.Debug(operationID, utils.GetSelfFuncName(), "args: ", m.String(), contentType)
	var err error
	var tips open_im_sdk.TipsComm
	tips.Detail, err = proto.Marshal(m)
	if err != nil {
		log.Error(operationID, "Marshal failed ", err.Error(), m.String())
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
	n.OperationID = operationID
	n.Content, err = proto.Marshal(&tips)
	if err != nil {
		log.Error(operationID, "Marshal failed ", err.Error(), tips.String())
		return
	}
	Notification(&n)
}
