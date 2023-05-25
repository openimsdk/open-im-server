package msg

import (
	"context"
	"encoding/json"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/tokenverify"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/msg"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
)

func (m *msgServer) RevokeMsg(ctx context.Context, req *msg.RevokeMsgReq) (*msg.RevokeMsgResp, error) {
	if req.UserID == "" {
		return nil, errs.ErrArgs.Wrap("user_id is empty")
	}
	if req.RecvID == "" && req.GroupID == "" {
		return nil, errs.ErrArgs.Wrap("recv_id and group_id are empty")
	}
	if req.RecvID != "" && req.GroupID != "" {
		return nil, errs.ErrArgs.Wrap("recv_id and group_id cannot exist at the same time")
	}
	if err := tokenverify.CheckAccessV3(ctx, req.RecvID); err != nil {
		return nil, err
	}
	var sessionType int32
	var conversationID string
	if req.GroupID == "" {
		sessionType = constant.SingleChatType
		conversationID = utils.GenConversationUniqueKeyForSingle(req.UserID, req.RecvID)
	} else {
		sessionType = constant.SuperGroupChatType
		conversationID = utils.GenConversationUniqueKeyForGroup(req.GroupID)
	}
	tips := sdkws.RevokeMsgTips{
		RevokerUserID:  req.UserID,
		ClientMsgID:    "",
		RevokeTime:     utils.GetCurrentTimestampByMill(),
		Seq:            req.Seq,
		SesstionType:   sessionType,
		ConversationID: conversationID,
	}
	detail, err := json.Marshal(&tips)
	if err != nil {
		return nil, err
	}
	notificationElem := sdkws.NotificationElem{Detail: string(detail)}
	content, err := json.Marshal(&notificationElem)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	msgData := sdkws.MsgData{
		SendID:      req.UserID,
		RecvID:      req.RecvID,
		GroupID:     req.GroupID,
		Content:     content,
		MsgFrom:     constant.SysMsgType,
		ContentType: constant.MsgRevokeNotification,
		SessionType: sessionType,
		CreateTime:  utils.GetCurrentTimestampByMill(),
		ClientMsgID: utils.GetMsgID(req.UserID),
		Options: config.GetOptionsByNotification(config.NotificationConf{
			IsSendMsg:        true,
			ReliabilityLevel: 2,
		}),
		OfflinePushInfo: nil,
	}
	if msgData.SessionType == constant.SuperGroupChatType {
		msgData.GroupID = msgData.RecvID
	}
	_, err = m.SendMsg(ctx, &msg.SendMsgReq{MsgData: &msgData})
	if err != nil {
		return nil, err
	}
	return &msg.RevokeMsgResp{}, nil
}
