package notification

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/tracelog"
	sdkws "Open_IM/pkg/proto/sdkws"
	"Open_IM/pkg/utils"
	"context"
	"github.com/golang/protobuf/jsonpb"
	"github.com/golang/protobuf/proto"
)

func SetConversationNotification(operationID, sendID, recvID string, contentType int, m proto.Message, tips sdkws.TipsComm) {
	log.NewInfo(operationID, "args: ", sendID, recvID, contentType, m.String(), tips.String())
	var err error
	tips.Detail, err = proto.Marshal(m)
	if err != nil {
		log.NewError(operationID, "Marshal failed ", err.Error(), m.String())
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
	n.ContentType = int32(contentType)
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

// SetPrivate调用
func ConversationSetPrivateNotification(ctx context.Context, sendID, recvID string, isPrivateChat bool) {
	conversationSetPrivateTips := &sdkws.ConversationSetPrivateTips{
		RecvID:    recvID,
		SendID:    sendID,
		IsPrivate: isPrivateChat,
	}
	var tips sdkws.TipsComm
	var tipsMsg string
	if isPrivateChat == true {
		tipsMsg = config.Config.Notification.ConversationSetPrivate.DefaultTips.OpenTips
	} else {
		tipsMsg = config.Config.Notification.ConversationSetPrivate.DefaultTips.CloseTips
	}
	tips.DefaultTips = tipsMsg
	SetConversationNotification(tracelog.GetOperationID(ctx), sendID, recvID, constant.ConversationPrivateChatNotification, conversationSetPrivateTips, tips)
}

// 会话改变
func ConversationChangeNotification(ctx context.Context, userID string) {

	ConversationChangedTips := &sdkws.ConversationUpdateTips{
		UserID: userID,
	}
	var tips sdkws.TipsComm
	tips.DefaultTips = config.Config.Notification.ConversationOptUpdate.DefaultTips.Tips
	SetConversationNotification(tracelog.GetOperationID(ctx), userID, userID, constant.ConversationOptChangeNotification, ConversationChangedTips, tips)
}

// 会话未读数同步
func ConversationUnreadChangeNotification(ctx context.Context, userID, conversationID string, updateUnreadCountTime int64) {

	ConversationChangedTips := &sdkws.ConversationUpdateTips{
		UserID:                userID,
		ConversationIDList:    []string{conversationID},
		UpdateUnreadCountTime: updateUnreadCountTime,
	}
	var tips sdkws.TipsComm
	tips.DefaultTips = config.Config.Notification.ConversationOptUpdate.DefaultTips.Tips
	SetConversationNotification(tracelog.GetOperationID(ctx), userID, userID, constant.ConversationUnreadNotification, ConversationChangedTips, tips)
}
