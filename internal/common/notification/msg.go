package notification

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/proto/sdkws"
	"Open_IM/pkg/utils"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

func (c *Check) DeleteMessageNotification(opUserID, userID string, seqList []uint32, operationID string) {
	DeleteMessageTips := sdkws.DeleteMessageTips{OpUserID: opUserID, UserID: userID, SeqList: seqList}
	c.MessageNotification(operationID, userID, userID, constant.DeleteMessageNotification, &DeleteMessageTips)
}

func (c *Check) MessageNotification(operationID, sendID, recvID string, contentType int32, m proto.Message) {
	log.Debug(operationID, utils.GetSelfFuncName(), "args: ", m.String(), contentType)
	var err error
	var tips sdkws.TipsComm
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
	c.Notification(&n)
}
