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
	"strconv"
	"strings"
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
	if err := authverify.CheckAccessV3(ctx, req.UserID, m.config.Share.IMAdminUserID); err != nil {
		return nil, err
	}
	user, err := m.UserLocalCache.GetUserInfo(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	
	// Get message and add fault tolerance handling
	_, _, msgs, err := m.MsgDatabase.GetMsgBySeqs(ctx, req.UserID, req.ConversationID, []int64{req.Seq})
	if err != nil {
		// Log the error but continue execution
		log.ZWarn(ctx, "GetMsgBySeqs error when revoking message", err, 
				"userID", req.UserID, "conversationID", req.ConversationID, "seq", req.Seq)
	} else if len(msgs) == 0 || msgs[0] == nil {
		// Check if seq is within valid range for the current conversation
		maxSeq, err := m.MsgDatabase.GetMaxSeq(ctx, req.ConversationID)
		if err != nil {
			log.ZWarn(ctx, "GetMaxSeq error when revoking message", err, 
					"conversationID", req.ConversationID)
		} else if req.Seq > maxSeq {
			return nil, errs.ErrArgs.WrapMsg("seq exceeds maxSeq")
		}
		
		// Log warning but continue execution
		log.ZWarn(ctx, "Message not found when revoking, but will proceed", nil,
				"userID", req.UserID, "conversationID", req.ConversationID, "seq", req.Seq)
	}

	// If message doesn't exist, create a minimal substitute for revocation notification
	var msgToRevoke *sdkws.MsgData
	if len(msgs) == 0 || msgs[0] == nil {
		// Create a minimal message object with necessary fields
		msgToRevoke = &sdkws.MsgData{
			SendID:         req.UserID,  // Use revoker as sender
			Seq:            req.Seq,
			SessionType:    getSessionTypeFromConversationID(req.ConversationID), // Helper function to get session type
			ClientMsgID:    "missing_" + strconv.FormatInt(req.Seq, 10), // Generate a temporary ID
			SendTime:       time.Now().UnixMilli() - 1000, // Set to a slightly earlier time
		}
		
		// Set GroupID or RecvID based on session type
		if msgToRevoke.SessionType == constant.ReadGroupChatType {
			// Extract group ID from conversation ID
			if strings.HasPrefix(req.ConversationID, "sg_") {
				msgToRevoke.GroupID = req.ConversationID[3:] // Remove "sg_" prefix
			}
		} else {
			// For single chat, parse receiver ID from conversation ID
			if strings.HasPrefix(req.ConversationID, "si_") || strings.HasPrefix(req.ConversationID, "sp_") {
				parts := strings.Split(req.ConversationID[3:], "_")
				if len(parts) == 2 {
					// Conversation ID format is typically: si_senderID_receiverID
					// If current user is the sender, then receiver is the other party
					if parts[0] == req.UserID {
						msgToRevoke.RecvID = parts[1]
					} else {
						msgToRevoke.RecvID = parts[0]
					}
				} else {
					// Set empty if parsing fails
					msgToRevoke.RecvID = ""
				}
			} else {
				msgToRevoke.RecvID = ""
			}
		}
	} else {
		msgToRevoke = msgs[0]
		if msgToRevoke.ContentType == constant.MsgRevokeNotification {
			return nil, servererrs.ErrMsgAlreadyRevoke.WrapMsg("msg already revoke")
		}
	}

	data, _ := json.Marshal(msgToRevoke)
	log.ZDebug(ctx, "GetMsgBySeqs", "conversationID", req.ConversationID, "seq", req.Seq, "msg", string(data))
	var role int32
	if !authverify.IsAppManagerUid(ctx, m.config.Share.IMAdminUserID) {
		sessionType := msgToRevoke.SessionType
		switch sessionType {
		case constant.SingleChatType:
			if err := authverify.CheckAccessV3(ctx, msgToRevoke.SendID, m.config.Share.IMAdminUserID); err != nil {
				return nil, err
			}
			role = user.AppMangerLevel
		case constant.ReadGroupChatType:
			members, err := m.GroupLocalCache.GetGroupMemberInfoMap(ctx, msgToRevoke.GroupID, datautil.Distinct([]string{req.UserID, msgToRevoke.SendID}))
			if err != nil {
				return nil, err
			}
			if req.UserID != msgToRevoke.SendID {
				switch members[req.UserID].RoleLevel {
				case constant.GroupOwner:
				case constant.GroupAdmin:
					if sendMember, ok := members[msgToRevoke.SendID]; ok {
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

	if len(m.config.Share.IMAdminUserID) > 0 {
		flag = datautil.Contain(revokerUserID, m.config.Share.IMAdminUserID...)
	}
	tips := sdkws.RevokeMsgTips{
		RevokerUserID:  revokerUserID,
		ClientMsgID:    msgToRevoke.ClientMsgID,
		RevokeTime:     now,
		Seq:            req.Seq,
		SesstionType:   msgToRevoke.SessionType,
		ConversationID: req.ConversationID,
		IsAdminRevoke:  flag,
	}
	var recvID string
	if msgToRevoke.SessionType == constant.ReadGroupChatType {
		recvID = msgToRevoke.GroupID
	} else {
		recvID = msgToRevoke.RecvID
	}
	m.notificationSender.NotificationWithSessionType(ctx, req.UserID, recvID, constant.MsgRevokeNotification, msgToRevoke.SessionType, &tips)
	m.webhookAfterRevokeMsg(ctx, &m.config.WebhooksConfig.AfterRevokeMsg, req)
	return &msg.RevokeMsgResp{}, nil
}

func getSessionTypeFromConversationID(conversationID string) int32 {
	// Conversation ID format is typically: "single chat prefix{userID}" or "group chat prefix{groupID}"
	if strings.HasPrefix(conversationID, "sp_") || strings.HasPrefix(conversationID, "si_") {
		return constant.SingleChatType
	} else if strings.HasPrefix(conversationID, "sg_") {
		return constant.ReadGroupChatType
	}
	// Default to single chat type
	return constant.SingleChatType
}
