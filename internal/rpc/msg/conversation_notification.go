package msg

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	open_im_sdk "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

func conversationNotification(contentType int32, m proto.Message, operationID, userID string) {
	var err error
	var tips open_im_sdk.TipsComm
	tips.Detail, err = proto.Marshal(m)
	if err != nil {
		log.Error(operationID, utils.GetSelfFuncName(), "Marshal failed ", err.Error(), m.String())
		return
	}
	marshaler := jsonpb.Marshaler{
		OrigName:     true,
		EnumsAsInts:  false,
		EmitDefaults: false,
	}
	tips.JsonDetail, _ = marshaler.MarshalToString(m)
	cn := config.Config.Notification
	switch contentType {
	case constant.ConversationOptChangeNotification:
		tips.DefaultTips = cn.ConversationOptUpdate.DefaultTips.Tips
	}
	var n NotificationMsg
	n.SendID = userID
	n.RecvID = userID
	n.ContentType = contentType
	n.SessionType = constant.SingleChatType
	n.MsgFrom = constant.SysMsgType
	n.OperationID = operationID
	n.Content, err = proto.Marshal(&tips)
	if err != nil {
		log.Error(operationID, utils.GetSelfFuncName(), "Marshal failed ", err.Error(), tips.String())
		return
	}
	Notification(&n)
}

// 客户端调用设置opt接口后调用
func SetReceiveMessageOptNotification(operationID, opUserID, userID string) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), "operation user: ", opUserID, "operation id: ", userID)
	conversationUpdateTips := open_im_sdk.ConversationUpdateTips{
		UserID: userID,
	}
	conversationNotification(constant.ConversationOptChangeNotification, &conversationUpdateTips, operationID, userID)
}
