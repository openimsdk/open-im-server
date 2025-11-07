package msg

import (
	"context"
	"encoding/json"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"

	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"github.com/openimsdk/open-im-server/v3/pkg/common/servererrs"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mcontext"
	"github.com/openimsdk/tools/utils/datautil"
)

func (m *msgServer) RevokeMsg(ctx context.Context, req *msg.RevokeMsgReq) (*msg.RevokeMsgResp, error) {
	if req.UserID == "" {
		return nil, errs.ErrArgs.WrapMsg("user_id is empty")
	}
	if req.ConversationID == "" {
		return nil, errs.ErrArgs.WrapMsg("conversation_id is empty")
	}
	if req.Seq < 0 {
		return nil, errs.ErrArgs.WrapMsg("seq is invalid")
	}
	if err := authverify.CheckAccess(ctx, req.UserID); err != nil {
		return nil, err
	}
	user, err := m.UserLocalCache.GetUserInfo(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	_, _, msgs, err := m.MsgDatabase.GetMsgBySeqs(ctx, req.UserID, req.ConversationID, []int64{req.Seq})
	if err != nil {
		return nil, err
	}
	if len(msgs) == 0 || msgs[0] == nil {
		return nil, errs.ErrRecordNotFound.WrapMsg("msg not found")
	}
	if msgs[0].ContentType == constant.MsgRevokeNotification {
		return nil, servererrs.ErrMsgAlreadyRevoke.WrapMsg("msg already revoke")
	}

	data, _ := json.Marshal(msgs[0])
	log.ZDebug(ctx, "GetMsgBySeqs", "conversationID", req.ConversationID, "seq", req.Seq, "msg", string(data))
	var role int32
	if !authverify.IsAdmin(ctx) {
		sessionType := msgs[0].SessionType
		switch sessionType {
		case constant.SingleChatType:
			if err := authverify.CheckAccess(ctx, msgs[0].SendID); err != nil {
				return nil, err
			}
			role = user.AppMangerLevel
		case constant.ReadGroupChatType:
			members, err := m.GroupLocalCache.GetGroupMemberInfoMap(ctx, msgs[0].GroupID, datautil.Distinct([]string{req.UserID, msgs[0].SendID}))
			if err != nil {
				return nil, err
			}
			if req.UserID != msgs[0].SendID {
				switch members[req.UserID].RoleLevel {
				case constant.GroupOwner:
				case constant.GroupAdmin:
					if sendMember, ok := members[msgs[0].SendID]; ok {
						if sendMember.RoleLevel != constant.GroupOrdinaryUsers {
							return nil, errs.ErrNoPermission.WrapMsg("no permission")
						}
					}
				default:
					return nil, errs.ErrNoPermission.WrapMsg("no permission")
				}
			}
			if member := members[req.UserID]; member != nil {
				role = member.RoleLevel
			}
		default:
			return nil, errs.ErrInternalServer.WrapMsg("msg sessionType not supported", "sessionType", sessionType)
		}
	}
	now := time.Now().UnixMilli()
	err = m.MsgDatabase.RevokeMsg(ctx, req.ConversationID, req.Seq, &model.RevokeModel{
		Role:     role,
		UserID:   req.UserID,
		Nickname: user.Nickname,
		Time:     now,
	})
	if err != nil {
		return nil, err
	}
	revokerUserID := mcontext.GetOpUserID(ctx)
	var flag bool

	if len(m.config.Share.IMAdminUser.UserIDs) > 0 {
		flag = datautil.Contain(revokerUserID, m.adminUserIDs...)
	}
	tips := sdkws.RevokeMsgTips{
		RevokerUserID:  revokerUserID,
		ClientMsgID:    msgs[0].ClientMsgID,
		RevokeTime:     now,
		Seq:            req.Seq,
		SesstionType:   msgs[0].SessionType,
		ConversationID: req.ConversationID,
		IsAdminRevoke:  flag,
	}
	var recvID string
	if msgs[0].SessionType == constant.ReadGroupChatType {
		recvID = msgs[0].GroupID
	} else {
		recvID = msgs[0].RecvID
	}
	m.notificationSender.NotificationWithSessionType(ctx, req.UserID, recvID, constant.MsgRevokeNotification, msgs[0].SessionType, &tips)
	m.webhookAfterRevokeMsg(ctx, &m.config.WebhooksConfig.AfterRevokeMsg, req)
	return &msg.RevokeMsgResp{}, nil
}
