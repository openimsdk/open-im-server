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
	
	// 获取消息并添加容错处理
	_, _, msgs, err := m.MsgDatabase.GetMsgBySeqs(ctx, req.UserID, req.ConversationID, []int64{req.Seq})
	if err != nil {
		// 记录错误但继续执行
		log.ZWarn(ctx, "GetMsgBySeqs error when revoking message", err, 
				"userID", req.UserID, "conversationID", req.ConversationID, "seq", req.Seq)
	} else if len(msgs) == 0 || msgs[0] == nil {
		// 检查seq是否在当前会话的有效范围内
		maxSeq, err := m.MsgDatabase.GetMaxSeq(ctx, req.ConversationID)
		if err != nil {
			log.ZWarn(ctx, "GetMaxSeq error when revoking message", err, 
					"conversationID", req.ConversationID)
		} else if req.Seq > maxSeq {
			return nil, errs.ErrArgs.WrapMsg("seq exceeds maxSeq")
		}
		
		// 记录警告但继续执行
		log.ZWarn(ctx, "Message not found when revoking, but will proceed", nil,
				"userID", req.UserID, "conversationID", req.ConversationID, "seq", req.Seq)
	}

	// 如果消息不存在，创建一个最小化的替代对象用于撤回通知
	var msgToRevoke *sdkws.MsgData
	if len(msgs) == 0 || msgs[0] == nil {
		// 创建一个最小的消息对象，包含必要的字段
		msgToRevoke = &sdkws.MsgData{
			SendID:         req.UserID,  // 使用撤回者作为发送者
			Seq:            req.Seq,
			SessionType:    getSessionTypeFromConversationID(req.ConversationID), // 辅助函数获取会话类型
			ClientMsgID:    "missing_" + strconv.FormatInt(req.Seq, 10), // 生成一个临时ID
			SendTime:       time.Now().UnixMilli() - 1000, // 设置为稍早的时间
		}
		
		// 根据会话类型设置GroupID或RecvID
		if msgToRevoke.SessionType == constant.ReadGroupChatType {
			// 从会话ID提取群组ID
			if strings.HasPrefix(req.ConversationID, "sg_") {
				msgToRevoke.GroupID = req.ConversationID[3:] // 移除"sg_"前缀
			}
		} else {
			// 对于单聊，需要从会话ID解析出接收者ID
			if strings.HasPrefix(req.ConversationID, "si_") || strings.HasPrefix(req.ConversationID, "sp_") {
				parts := strings.Split(req.ConversationID[3:], "_")
				if len(parts) == 2 {
					// 会话ID一般格式为: si_发送者ID_接收者ID
					// 如果当前用户是发送者，则接收者是另一方
					if parts[0] == req.UserID {
						msgToRevoke.RecvID = parts[1]
					} else {
						msgToRevoke.RecvID = parts[0]
					}
				} else {
					// 无法解析时设置为空
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
	// 通常会话ID格式为: "单聊前缀{userID}" 或 "群聊前缀{groupID}"
	if strings.HasPrefix(conversationID, "sp_") || strings.HasPrefix(conversationID, "si_") {
		return constant.SingleChatType
	} else if strings.HasPrefix(conversationID, "sg_") {
		return constant.ReadGroupChatType
	}
	// 默认返回单聊类型
	return constant.SingleChatType
}
