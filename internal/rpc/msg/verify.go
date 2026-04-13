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
	"math/rand"
	"strconv"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/servererrs"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/datautil"
	"github.com/openimsdk/tools/utils/encrypt"
	"github.com/openimsdk/tools/utils/timeutil"
)

var ExcludeContentType = []int{constant.HasReadReceipt}

type Validator interface {
	validate(pb *msg.SendMsgReq) (bool, int32, string)
}

type MessageRevoked struct {
	RevokerID                   string `json:"revokerID"`
	RevokerRole                 int32  `json:"revokerRole"`
	ClientMsgID                 string `json:"clientMsgID"`
	RevokerNickname             string `json:"revokerNickname"`
	RevokeTime                  int64  `json:"revokeTime"`
	SourceMessageSendTime       int64  `json:"sourceMessageSendTime"`
	SourceMessageSendID         string `json:"sourceMessageSendID"`
	SourceMessageSenderNickname string `json:"sourceMessageSenderNickname"`
	SessionType                 int32  `json:"sessionType"`
	Seq                         uint32 `json:"seq"`
}

func (m *msgServer) messageVerification(ctx context.Context, data *msg.SendMsgReq) error {
	switch data.MsgData.SessionType {
	case constant.SingleChatType:
		// 单聊发送权限校验已迁移至 modifyMessageByUserMessageReceiveOpt
		return nil
	case constant.ReadGroupChatType:
		groupInfo, err := m.GroupLocalCache.GetGroupInfo(ctx, data.MsgData.GroupID)
		if err != nil {
			log.ZError(ctx, "messageVerification group: GetGroupInfo failed", err,
				"groupID", data.MsgData.GroupID, "sendID", data.MsgData.SendID,
				"contentType", data.MsgData.ContentType, "clientMsgID", data.MsgData.ClientMsgID)
			return err
		}
		if groupInfo.Status == constant.GroupStatusDismissed &&
			data.MsgData.ContentType != constant.GroupDismissedNotification {
			return servererrs.ErrDismissedAlready.Wrap()
		}
		if groupInfo.GroupType == constant.SuperGroup {
			return nil
		}

		if datautil.Contain(data.MsgData.SendID, m.config.Share.IMAdminUserID...) {
			return nil
		}
		if data.MsgData.ContentType <= constant.NotificationEnd &&
			data.MsgData.ContentType >= constant.NotificationBegin {
			return nil
		}
		memberIDs, err := m.GroupLocalCache.GetGroupMemberIDMap(ctx, data.MsgData.GroupID)
		if err != nil {
			log.ZError(ctx, "messageVerification group: GetGroupMemberIDMap failed", err,
				"groupID", data.MsgData.GroupID, "sendID", data.MsgData.SendID,
				"contentType", data.MsgData.ContentType, "clientMsgID", data.MsgData.ClientMsgID)
			return err
		}
		if _, ok := memberIDs[data.MsgData.SendID]; !ok {
			return servererrs.ErrNotInGroupYet.Wrap()
		}

		groupMemberInfo, err := m.GroupLocalCache.GetGroupMember(ctx, data.MsgData.GroupID, data.MsgData.SendID)
		if err != nil {
			if errs.ErrRecordNotFound.Is(err) {
				return servererrs.ErrNotInGroupYet.WrapMsg(err.Error())
			}
			log.ZError(ctx, "messageVerification group: GetGroupMember failed", err,
				"groupID", data.MsgData.GroupID, "sendID", data.MsgData.SendID,
				"contentType", data.MsgData.ContentType, "clientMsgID", data.MsgData.ClientMsgID)
			return err
		}
		if groupMemberInfo.RoleLevel == constant.GroupOwner {
			return nil
		} else {
			if groupMemberInfo.MuteEndTime >= time.Now().UnixMilli() {
				return servererrs.ErrMutedInGroup.Wrap()
			}
			if groupInfo.Status == constant.GroupStatusMuted && groupMemberInfo.RoleLevel != constant.GroupAdmin {
				return servererrs.ErrMutedGroup.Wrap()
			}
			// AllowSendMsg == 1 时仅群主/管理员可发消息
			if groupInfo.AllowSendMsg == 1 && groupMemberInfo.RoleLevel != constant.GroupAdmin {
				return servererrs.ErrNoPermission.WrapMsg("only owner or admin can send messages in this group")
			}
		}
		return nil
	default:
		return nil
	}
}

func (m *msgServer) encapsulateMsgData(msg *sdkws.MsgData) {
	msg.ServerMsgID = GetMsgID(msg.SendID)
	if msg.SendTime == 0 {
		msg.SendTime = timeutil.GetCurrentTimestampByMill()
	}
	switch msg.ContentType {
	case constant.Text:
		fallthrough
	case constant.Picture:
		fallthrough
	case constant.Voice:
		fallthrough
	case constant.Video:
		fallthrough
	case constant.File:
		fallthrough
	case constant.AtText:
		fallthrough
	case constant.Merger:
		fallthrough
	case constant.Card:
		fallthrough
	case constant.Location:
		fallthrough
	case constant.Custom:
		fallthrough
	case constant.Quote:
	case constant.Revoke:
		datautil.SetSwitchFromOptions(msg.Options, constant.IsUnreadCount, false)
		datautil.SetSwitchFromOptions(msg.Options, constant.IsOfflinePush, false)
	case constant.HasReadReceipt:
		datautil.SetSwitchFromOptions(msg.Options, constant.IsConversationUpdate, false)
		datautil.SetSwitchFromOptions(msg.Options, constant.IsSenderConversationUpdate, false)
		datautil.SetSwitchFromOptions(msg.Options, constant.IsUnreadCount, false)
		datautil.SetSwitchFromOptions(msg.Options, constant.IsOfflinePush, false)
	case constant.Typing:
		datautil.SetSwitchFromOptions(msg.Options, constant.IsHistory, false)
		datautil.SetSwitchFromOptions(msg.Options, constant.IsPersistent, false)
		datautil.SetSwitchFromOptions(msg.Options, constant.IsSenderSync, false)
		datautil.SetSwitchFromOptions(msg.Options, constant.IsConversationUpdate, false)
		datautil.SetSwitchFromOptions(msg.Options, constant.IsSenderConversationUpdate, false)
		datautil.SetSwitchFromOptions(msg.Options, constant.IsUnreadCount, false)
		datautil.SetSwitchFromOptions(msg.Options, constant.IsOfflinePush, false)
	}
}

func GetMsgID(sendID string) string {
	t := timeutil.GetCurrentTimeFormatted()
	return encrypt.Md5(t + "-" + sendID + "-" + strconv.Itoa(rand.Int()))
}

func (m *msgServer) modifyMessageByUserMessageReceiveOpt(ctx context.Context, userID, conversationID string, sessionType int, pb *msg.SendMsgReq) (bool, error) {
	// 第一优先级：接收方全局接收设置
	// NotReceiveMessage 直接丢弃，无需执行后续任何权限或偏好查询
	opt, err := m.UserLocalCache.GetUserGlobalMsgRecvOpt(ctx, userID)
	if err != nil {
		return false, err
	}
	if opt == constant.NotReceiveMessage {
		return false, nil
	}
	if opt == constant.ReceiveNotNotifyMessage {
		if pb.MsgData.Options == nil {
			pb.MsgData.Options = make(map[string]bool, 10)
		}
		datautil.SetSwitchFromOptions(pb.MsgData.Options, constant.IsOfflinePush, false)
		// 全局静音：仅关闭离线推送，仍需继续执行发送权限校验 + 会话级偏好校验
	}

	// 第二优先级：单聊发送权限校验（从 messageVerification 迁移）
	// 仅对非通知类消息生效（调用方已通过 !isNotification 做过前置过滤）
	if sessionType == constant.SingleChatType {
		// 管理员跳过发送权限拦截，直接进入接收偏好校验
		if !datautil.Contain(pb.MsgData.SendID, m.config.Share.IMAdminUserID...) {
			// 黑名单拦截
			black, err := m.FriendLocalCache.IsBlack(ctx, pb.MsgData.SendID, pb.MsgData.RecvID)
			if err != nil {
				log.ZError(ctx, "modifyMessageByUserMessageReceiveOpt: IsBlack failed", err,
					"sendID", pb.MsgData.SendID, "recvID", pb.MsgData.RecvID,
					"contentType", pb.MsgData.ContentType, "clientMsgID", pb.MsgData.ClientMsgID)
				return false, err
			}
			if black {
				return false, servererrs.ErrBlockedByPeer.Wrap()
			}

			// 接收方消息接收权限（MsgReceiveSetting）
			// 0=所有人可发送，1=仅好友可发送，2=所有人不可发送
			recvUserInfo, err := m.UserLocalCache.GetUserInfo(ctx, pb.MsgData.RecvID)
			if err != nil {
				log.ZError(ctx, "modifyMessageByUserMessageReceiveOpt: GetUserInfo(recv) failed", err,
					"sendID", pb.MsgData.SendID, "recvID", pb.MsgData.RecvID,
					"contentType", pb.MsgData.ContentType, "clientMsgID", pb.MsgData.ClientMsgID)
				return false, err
			}

			// skipFriendVerify: MsgReceiveSetting=1 已确认好友关系，无需再做 FriendVerify 重复查询
			skipFriendVerify := false
			switch recvUserInfo.MsgReceiveSetting {
			case 2: // MsgReceiveSettingNobody
				return false, servererrs.ErrMsgReceiveNotAllowed.Wrap()
			case 1: // MsgReceiveSettingFriends
				isFriend, err := m.FriendLocalCache.IsFriend(ctx, pb.MsgData.RecvID, pb.MsgData.SendID)
				if err != nil {
					log.ZError(ctx, "modifyMessageByUserMessageReceiveOpt: IsFriend failed (MsgReceiveSetting)", err,
						"sendID", pb.MsgData.SendID, "recvID", pb.MsgData.RecvID,
						"contentType", pb.MsgData.ContentType, "clientMsgID", pb.MsgData.ClientMsgID)
					return false, err
				}
				if !isFriend {
					return false, servererrs.ErrMsgReceiveNotAllowed.Wrap()
				}
				// 已确认好友关系，触发 webhook 后跳过 FriendVerify，直接进入接收偏好校验
				if err := m.webhookBeforeSendSingleMsg(ctx, &m.config.WebhooksConfig.BeforeSendSingleMsg, pb); err != nil {
					log.ZError(ctx, "modifyMessageByUserMessageReceiveOpt: webhookBeforeSendSingleMsg failed (friends-only)", err,
						"sendID", pb.MsgData.SendID, "recvID", pb.MsgData.RecvID,
						"contentType", pb.MsgData.ContentType, "clientMsgID", pb.MsgData.ClientMsgID)
					return false, err
				}
				skipFriendVerify = true
			}

			if !skipFriendVerify {
				// MsgReceiveSetting==0（所有人可发），触发 webhook，再按全局 FriendVerify 兜底
				if err := m.webhookBeforeSendSingleMsg(ctx, &m.config.WebhooksConfig.BeforeSendSingleMsg, pb); err != nil {
					log.ZError(ctx, "modifyMessageByUserMessageReceiveOpt: webhookBeforeSendSingleMsg failed", err,
						"sendID", pb.MsgData.SendID, "recvID", pb.MsgData.RecvID,
						"contentType", pb.MsgData.ContentType, "clientMsgID", pb.MsgData.ClientMsgID)
					return false, err
				}
				if m.config.RpcConfig.FriendVerify {
					friend, err := m.FriendLocalCache.IsFriend(ctx, pb.MsgData.SendID, pb.MsgData.RecvID)
					if err != nil {
						log.ZError(ctx, "modifyMessageByUserMessageReceiveOpt: IsFriend failed (FriendVerify)", err,
							"sendID", pb.MsgData.SendID, "recvID", pb.MsgData.RecvID,
							"contentType", pb.MsgData.ContentType, "clientMsgID", pb.MsgData.ClientMsgID)
						return false, err
					}
					if !friend {
						return false, servererrs.ErrNotPeersFriend.Wrap()
					}
				}
			}
		}
	}

	// 第三优先级：会话级接收偏好
	singleOpt, err := m.ConversationLocalCache.GetSingleConversationRecvMsgOpt(ctx, userID, conversationID)
	if errs.ErrRecordNotFound.Is(err) {
		return true, nil
	} else if err != nil {
		return false, err
	}
	switch singleOpt {
	case constant.ReceiveMessage:
		return true, nil
	case constant.NotReceiveMessage:
		if datautil.Contain(int(pb.MsgData.ContentType), ExcludeContentType...) {
			return true, nil
		}
		return false, nil
	case constant.ReceiveNotNotifyMessage:
		if pb.MsgData.Options == nil {
			pb.MsgData.Options = make(map[string]bool, 10)
		}
		datautil.SetSwitchFromOptions(pb.MsgData.Options, constant.IsOfflinePush, false)
		return true, nil
	}
	return true, nil
}
