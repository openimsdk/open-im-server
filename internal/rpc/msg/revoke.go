package msg

import (
	"context"
	"encoding/json"
	"time"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	unRelationTb "github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/unrelation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/tokenverify"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/msg"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
)

func (m *msgServer) RevokeMsg(ctx context.Context, req *msg.RevokeMsgReq) (*msg.RevokeMsgResp, error) {
	defer log.ZInfo(ctx, "RevokeMsg return line")
	//if req.UserID == "" {
	//	return nil, errs.ErrArgs.Wrap("user_id is empty")
	//}
	//if req.RecvID == "" && req.GroupID == "" {
	//	return nil, errs.ErrArgs.Wrap("recv_id and group_id are empty")
	//}
	//if req.RecvID != "" && req.GroupID != "" {
	//	return nil, errs.ErrArgs.Wrap("recv_id and group_id cannot exist at the same time")
	//}
	if req.ConversationID == "" {
		return nil, errs.ErrArgs.Wrap("conversation_id is empty")
	}
	if req.Seq < 0 {
		return nil, errs.ErrArgs.Wrap("seq is invalid")
	}
	if err := tokenverify.CheckAccessV3(ctx, req.UserID); err != nil {
		return nil, err
	}
	user, err := m.User.GetUserInfo(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	msgs, err := m.MsgDatabase.GetMsgBySeqs(ctx, req.UserID, req.ConversationID, []int64{req.Seq})
	if err != nil {
		return nil, err
	}
	if len(msgs) == 0 {
		return nil, errs.ErrRecordNotFound.Wrap("msg not found")
	}
	data, _ := json.Marshal(msgs[0])
	log.ZInfo(ctx, "GetMsgBySeqs", "conversationID", req.ConversationID, "seq", req.Seq, "msg", string(data))
	//sendID := msgs[0]
	//if !tokenverify.IsAppManagerUid(ctx) {
	//	if req.GroupID == "" {
	//		if req.UserID != sendID {
	//			return nil, errs.ErrNoPermission.Wrap("no permission")
	//		}
	//	} else {
	//		members, err := m.Group.GetGroupMemberInfoMap(ctx, req.GroupID, utils.Distinct([]string{req.UserID, sendID}), true)
	//		if err != nil {
	//			return nil, err
	//		}
	//		if req.UserID != sendID {
	//			roleLevel := members[req.UserID].RoleLevel
	//			switch members[req.UserID].RoleLevel {
	//			case constant.GroupOwner:
	//			case constant.GroupAdmin:
	//				if roleLevel != constant.GroupOrdinaryUsers {
	//					return nil, errs.ErrNoPermission.Wrap("no permission")
	//				}
	//			default:
	//				return nil, errs.ErrNoPermission.Wrap("no permission")
	//			}
	//		}
	//	}
	//}
	err = m.MsgDatabase.RevokeMsg(ctx, req.ConversationID, req.Seq, &unRelationTb.RevokeModel{
		UserID:   req.UserID,
		Nickname: user.Nickname,
		Time:     time.Now().UnixMilli(),
	})
	if err != nil {
		return nil, err
	}
	tips := sdkws.RevokeMsgTips{
		RevokerUserID:  req.UserID,
		ClientMsgID:    "",
		RevokeTime:     utils.GetCurrentTimestampByMill(),
		Seq:            req.Seq,
		SesstionType:   msgs[0].SessionType,
		ConversationID: req.ConversationID,
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
		RecvID:      msgs[0].RecvID,
		GroupID:     msgs[0].GroupID,
		Content:     content,
		MsgFrom:     constant.SysMsgType,
		ContentType: constant.MsgRevokeNotification,
		SessionType: msgs[0].SessionType,
		CreateTime:  utils.GetCurrentTimestampByMill(),
		ClientMsgID: utils.GetMsgID(req.UserID),
		Options: config.GetOptionsByNotification(config.NotificationConf{
			IsSendMsg:        false,
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
