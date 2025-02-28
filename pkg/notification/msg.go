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

package rpcclient

import (
	"context"
	"encoding/json"
	"time"

	"google.golang.org/protobuf/proto"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mq/memamq"
	"github.com/openimsdk/tools/utils/idutil"
	"github.com/openimsdk/tools/utils/jsonutil"
	"github.com/openimsdk/tools/utils/timeutil"
)

func newContentTypeConf(conf *config.Notification) map[int32]config.NotificationConfig {
	return map[int32]config.NotificationConfig{
		// group
		constant.GroupCreatedNotification:                 conf.GroupCreated,
		constant.GroupInfoSetNotification:                 conf.GroupInfoSet,
		constant.JoinGroupApplicationNotification:         conf.JoinGroupApplication,
		constant.MemberQuitNotification:                   conf.MemberQuit,
		constant.GroupApplicationAcceptedNotification:     conf.GroupApplicationAccepted,
		constant.GroupApplicationRejectedNotification:     conf.GroupApplicationRejected,
		constant.GroupOwnerTransferredNotification:        conf.GroupOwnerTransferred,
		constant.MemberKickedNotification:                 conf.MemberKicked,
		constant.MemberInvitedNotification:                conf.MemberInvited,
		constant.MemberEnterNotification:                  conf.MemberEnter,
		constant.GroupDismissedNotification:               conf.GroupDismissed,
		constant.GroupMutedNotification:                   conf.GroupMuted,
		constant.GroupCancelMutedNotification:             conf.GroupCancelMuted,
		constant.GroupMemberMutedNotification:             conf.GroupMemberMuted,
		constant.GroupMemberCancelMutedNotification:       conf.GroupMemberCancelMuted,
		constant.GroupMemberInfoSetNotification:           conf.GroupMemberInfoSet,
		constant.GroupMemberSetToAdminNotification:        conf.GroupMemberSetToAdmin,
		constant.GroupMemberSetToOrdinaryUserNotification: conf.GroupMemberSetToOrdinary,
		constant.GroupInfoSetAnnouncementNotification:     conf.GroupInfoSetAnnouncement,
		constant.GroupInfoSetNameNotification:             conf.GroupInfoSetName,
		// user
		constant.UserInfoUpdatedNotification:  conf.UserInfoUpdated,
		constant.UserStatusChangeNotification: conf.UserStatusChanged,
		// friend
		constant.FriendApplicationNotification:         conf.FriendApplicationAdded,
		constant.FriendApplicationApprovedNotification: conf.FriendApplicationApproved,
		constant.FriendApplicationRejectedNotification: conf.FriendApplicationRejected,
		constant.FriendAddedNotification:               conf.FriendAdded,
		constant.FriendDeletedNotification:             conf.FriendDeleted,
		constant.FriendRemarkSetNotification:           conf.FriendRemarkSet,
		constant.BlackAddedNotification:                conf.BlackAdded,
		constant.BlackDeletedNotification:              conf.BlackDeleted,
		constant.FriendInfoUpdatedNotification:         conf.FriendInfoUpdated,
		constant.FriendsInfoUpdateNotification:         conf.FriendInfoUpdated, // use the same FriendInfoUpdated
		// conversation
		constant.ConversationChangeNotification:      conf.ConversationChanged,
		constant.ConversationUnreadNotification:      conf.ConversationChanged,
		constant.ConversationPrivateChatNotification: conf.ConversationSetPrivate,
		// msg
		constant.MsgRevokeNotification:  {IsSendMsg: false, ReliabilityLevel: constant.ReliableNotificationNoMsg},
		constant.HasReadReceipt:         {IsSendMsg: false, ReliabilityLevel: constant.ReliableNotificationNoMsg},
		constant.DeleteMsgsNotification: {IsSendMsg: false, ReliabilityLevel: constant.ReliableNotificationNoMsg},
	}
}

func newSessionTypeConf() map[int32]int32 {
	return map[int32]int32{
		// group
		constant.GroupCreatedNotification:                 constant.ReadGroupChatType,
		constant.GroupInfoSetNotification:                 constant.ReadGroupChatType,
		constant.JoinGroupApplicationNotification:         constant.SingleChatType,
		constant.MemberQuitNotification:                   constant.ReadGroupChatType,
		constant.GroupApplicationAcceptedNotification:     constant.SingleChatType,
		constant.GroupApplicationRejectedNotification:     constant.SingleChatType,
		constant.GroupOwnerTransferredNotification:        constant.ReadGroupChatType,
		constant.MemberKickedNotification:                 constant.ReadGroupChatType,
		constant.MemberInvitedNotification:                constant.ReadGroupChatType,
		constant.MemberEnterNotification:                  constant.ReadGroupChatType,
		constant.GroupDismissedNotification:               constant.ReadGroupChatType,
		constant.GroupMutedNotification:                   constant.ReadGroupChatType,
		constant.GroupCancelMutedNotification:             constant.ReadGroupChatType,
		constant.GroupMemberMutedNotification:             constant.ReadGroupChatType,
		constant.GroupMemberCancelMutedNotification:       constant.ReadGroupChatType,
		constant.GroupMemberInfoSetNotification:           constant.ReadGroupChatType,
		constant.GroupMemberSetToAdminNotification:        constant.ReadGroupChatType,
		constant.GroupMemberSetToOrdinaryUserNotification: constant.ReadGroupChatType,
		constant.GroupInfoSetAnnouncementNotification:     constant.ReadGroupChatType,
		constant.GroupInfoSetNameNotification:             constant.ReadGroupChatType,
		// user
		constant.UserInfoUpdatedNotification:  constant.SingleChatType,
		constant.UserStatusChangeNotification: constant.SingleChatType,
		// friend
		constant.FriendApplicationNotification:         constant.SingleChatType,
		constant.FriendApplicationApprovedNotification: constant.SingleChatType,
		constant.FriendApplicationRejectedNotification: constant.SingleChatType,
		constant.FriendAddedNotification:               constant.SingleChatType,
		constant.FriendDeletedNotification:             constant.SingleChatType,
		constant.FriendRemarkSetNotification:           constant.SingleChatType,
		constant.BlackAddedNotification:                constant.SingleChatType,
		constant.BlackDeletedNotification:              constant.SingleChatType,
		constant.FriendInfoUpdatedNotification:         constant.SingleChatType,
		constant.FriendsInfoUpdateNotification:         constant.SingleChatType,
		// conversation
		constant.ConversationChangeNotification:      constant.SingleChatType,
		constant.ConversationUnreadNotification:      constant.SingleChatType,
		constant.ConversationPrivateChatNotification: constant.SingleChatType,
		// delete
		constant.DeleteMsgsNotification: constant.SingleChatType,
	}
}

type NotificationSender struct {
	contentTypeConf map[int32]config.NotificationConfig
	sessionTypeConf map[int32]int32
	sendMsg         func(ctx context.Context, req *msg.SendMsgReq) (*msg.SendMsgResp, error)
	getUserInfo     func(ctx context.Context, userID string) (*sdkws.UserInfo, error)
	queue           *memamq.MemoryQueue
}

func WithQueue(queue *memamq.MemoryQueue) NotificationSenderOptions {
	return func(s *NotificationSender) {
		s.queue = queue
	}
}

type NotificationSenderOptions func(*NotificationSender)

func WithLocalSendMsg(sendMsg func(ctx context.Context, req *msg.SendMsgReq) (*msg.SendMsgResp, error)) NotificationSenderOptions {
	return func(s *NotificationSender) {
		s.sendMsg = sendMsg
	}
}

func WithRpcClient(sendMsg func(ctx context.Context, req *msg.SendMsgReq) (*msg.SendMsgResp, error)) NotificationSenderOptions {
	return func(s *NotificationSender) {
		s.sendMsg = func(ctx context.Context, req *msg.SendMsgReq) (*msg.SendMsgResp, error) {
			return sendMsg(ctx, req)
		}
	}
}

func WithUserRpcClient(getUserInfo func(ctx context.Context, userID string) (*sdkws.UserInfo, error)) NotificationSenderOptions {
	return func(s *NotificationSender) {
		s.getUserInfo = getUserInfo
	}
}

const (
	notificationWorkerCount = 16
	notificationBufferSize  = 1024 * 1024 * 2
)

func NewNotificationSender(conf *config.Notification, opts ...NotificationSenderOptions) *NotificationSender {
	notificationSender := &NotificationSender{contentTypeConf: newContentTypeConf(conf), sessionTypeConf: newSessionTypeConf()}
	for _, opt := range opts {
		opt(notificationSender)
	}
	if notificationSender.queue == nil {
		notificationSender.queue = memamq.NewMemoryQueue(notificationWorkerCount, notificationBufferSize)
	}
	return notificationSender
}

type notificationOpt struct {
	RpcGetUsername bool
	SendMessage    *bool
}

type NotificationOptions func(*notificationOpt)

func WithRpcGetUserName() NotificationOptions {
	return func(opt *notificationOpt) {
		opt.RpcGetUsername = true
	}
}
func WithSendMessage(sendMessage *bool) NotificationOptions {
	return func(opt *notificationOpt) {
		opt.SendMessage = sendMessage
	}
}

func (s *NotificationSender) send(ctx context.Context, sendID, recvID string, contentType, sessionType int32, m proto.Message, opts ...NotificationOptions) {
	//ctx = mcontext.WithMustInfoCtx([]string{mcontext.GetOperationID(ctx), mcontext.GetOpUserID(ctx), mcontext.GetOpUserPlatform(ctx), mcontext.GetConnID(ctx)})
	ctx = context.WithoutCancel(ctx)
	ctx, cancel := context.WithTimeout(ctx, time.Second*time.Duration(5))
	defer cancel()
	n := sdkws.NotificationElem{Detail: jsonutil.StructToJsonString(m)}
	content, err := json.Marshal(&n)
	if err != nil {
		log.ZWarn(ctx, "json.Marshal failed", err, "sendID", sendID, "recvID", recvID, "contentType", contentType, "msg", jsonutil.StructToJsonString(m))
		return
	}
	notificationOpt := &notificationOpt{}
	for _, opt := range opts {
		opt(notificationOpt)
	}
	var req msg.SendMsgReq
	var msg sdkws.MsgData
	var userInfo *sdkws.UserInfo
	if notificationOpt.WithRpcGetUsername && s.getUserInfo != nil {
		userInfo, err = s.getUserInfo(ctx, sendID)
		if err != nil {
			log.ZWarn(ctx, "getUserInfo failed", err, "sendID", sendID)
			return
		}
		msg.SenderNickname = userInfo.Nickname
		msg.SenderFaceURL = userInfo.FaceURL
	}
	var offlineInfo sdkws.OfflinePushInfo
	msg.SendID = sendID
	msg.RecvID = recvID
	msg.Content = content
	msg.MsgFrom = constant.SysMsgType
	msg.ContentType = contentType
	msg.SessionType = sessionType
	if msg.SessionType == constant.ReadGroupChatType {
		msg.GroupID = recvID
	}
	msg.CreateTime = timeutil.GetCurrentTimestampByMill()
	msg.ClientMsgID = idutil.GetMsgIDByMD5(sendID)
	optionsConfig := s.contentTypeConf[contentType]
	if sendID == recvID && contentType == constant.HasReadReceipt {
		optionsConfig.ReliabilityLevel = constant.UnreliableNotification
	}
	options := config.GetOptionsByNotification(optionsConfig, notificationOpt.SendMessage)
	s.SetOptionsByContentType(ctx, options, contentType)
	msg.Options = options
	// fill Notification OfflinePush by config
	offlineInfo.Title = optionsConfig.OfflinePush.Title
	offlineInfo.Desc = optionsConfig.OfflinePush.Desc
	offlineInfo.Ex = optionsConfig.OfflinePush.Ext
	msg.OfflinePushInfo = &offlineInfo
	req.MsgData = &msg
	_, err = s.sendMsg(ctx, &req)
	if err != nil {
		log.ZWarn(ctx, "SendMsg failed", err, "req", req.String())
	}
}

func (s *NotificationSender) NotificationWithSessionType(ctx context.Context, sendID, recvID string, contentType, sessionType int32, m proto.Message, opts ...NotificationOptions) {
	if err := s.queue.Push(func() { s.send(ctx, sendID, recvID, contentType, sessionType, m, opts...) }); err != nil {
		log.ZWarn(ctx, "Push to queue failed", err, "sendID", sendID, "recvID", recvID, "msg", jsonutil.StructToJsonString(m))
	}
}

func (s *NotificationSender) Notification(ctx context.Context, sendID, recvID string, contentType int32, m proto.Message, opts ...NotificationOptions) {
	s.NotificationWithSessionType(ctx, sendID, recvID, contentType, s.sessionTypeConf[contentType], m, opts...)
}

func (s *NotificationSender) SetOptionsByContentType(_ context.Context, options map[string]bool, contentType int32) {
	switch contentType {
	case constant.UserStatusChangeNotification:
		options[constant.IsSenderSync] = false
	default:
	}
}
