// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package msg

import (
	"context"
	"encoding/json"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/table/relation"
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
	if err := authverify.CheckAccessV3(ctx, req.UserID, m.config.Share.IMAdminUserID); err != nil {
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
	if !authverify.IsAppManagerUid(ctx, m.config.Share.IMAdminUserID) {
		switch msgs[0].SessionType {
		case constant.SingleChatType:
			if err := authverify.CheckAccessV3(ctx, msgs[0].SendID, m.config.Share.IMAdminUserID); err != nil {
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
					if members[msgs[0].SendID].RoleLevel != constant.GroupOrdinaryUsers {
						return nil, errs.ErrNoPermission.WrapMsg("no permission")
					}
				default:
					return nil, errs.ErrNoPermission.WrapMsg("no permission")
				}
			}
			if member := members[req.UserID]; member != nil {
				role = member.RoleLevel
			}
		default:
			return nil, errs.ErrInternalServer.WrapMsg("msg sessionType not supported")
		}
	}
	now := time.Now().UnixMilli()
	err = m.MsgDatabase.RevokeMsg(ctx, req.ConversationID, req.Seq, &relation.RevokeModel{
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

	if len(m.config.Share.IMAdminUserID) > 0 {
		flag = datautil.Contain(revokerUserID, m.config.Share.IMAdminUserID...)
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
