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

	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mq/memamq"
	"github.com/openimsdk/tools/system/program"
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

type Message struct {
	conn   grpc.ClientConnInterface
	Client msg.MsgClient
	discov discovery.SvcDiscoveryRegistry
}

func NewMessage(discov discovery.SvcDiscoveryRegistry, rpcRegisterName string) *Message {
	conn, err := discov.GetConn(context.Background(), rpcRegisterName)
	if err != nil {
		program.ExitWithError(err)
	}
	client := msg.NewMsgClient(conn)
	return &Message{discov: discov, conn: conn, Client: client}
}

type MessageRpcClient Message

func NewMessageRpcClient(discov discovery.SvcDiscoveryRegistry, rpcRegisterName string) MessageRpcClient {
	return MessageRpcClient(*NewMessage(discov, rpcRegisterName))
}

// SendMsg sends a message through the gRPC client and returns the response.
// It wraps any encountered error for better error handling and context understanding.
func (m *MessageRpcClient) SendMsg(ctx context.Context, req *msg.SendMsgReq) (*msg.SendMsgResp, error) {
	resp, err := m.Client.SendMsg(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// SetUserConversationsMinSeq set min seq
func (m *MessageRpcClient) SetUserConversationsMinSeq(ctx context.Context, req *msg.SetUserConversationsMinSeqReq) (*msg.SetUserConversationsMinSeqResp, error) {
	resp, err := m.Client.SetUserConversationsMinSeq(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// GetMaxSeq retrieves the maximum sequence number from the gRPC client.
// Errors during the gRPC call are wrapped to provide additional context.
func (m *MessageRpcClient) GetMaxSeq(ctx context.Context, req *sdkws.GetMaxSeqReq) (*sdkws.GetMaxSeqResp, error) {
	resp, err := m.Client.GetMaxSeq(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (m *MessageRpcClient) GetMaxSeqs(ctx context.Context, conversationIDs []string) (map[string]int64, error) {
	log.ZDebug(ctx, "GetMaxSeqs", "conversationIDs", conversationIDs)
	resp, err := m.Client.GetMaxSeqs(ctx, &msg.GetMaxSeqsReq{
		ConversationIDs: conversationIDs,
	})
	if err != nil {
		return nil, err
	}
	return resp.MaxSeqs, err
}

func (m *MessageRpcClient) GetHasReadSeqs(ctx context.Context, userID string, conversationIDs []string) (map[string]int64, error) {
	resp, err := m.Client.GetHasReadSeqs(ctx, &msg.GetHasReadSeqsReq{
		UserID:          userID,
		ConversationIDs: conversationIDs,
	})
	if err != nil {
		return nil, err
	}
	return resp.MaxSeqs, err
}

func (m *MessageRpcClient) GetMsgByConversationIDs(ctx context.Context, docIDs []string, seqs map[string]int64) (map[string]*sdkws.MsgData, error) {
	resp, err := m.Client.GetMsgByConversationIDs(ctx, &msg.GetMsgByConversationIDsReq{
		ConversationIDs: docIDs,
		MaxSeqs:         seqs,
	})
	if err != nil {
		return nil, err
	}
	return resp.MsgDatas, err
}

// PullMessageBySeqList retrieves messages by their sequence numbers using the gRPC client.
// It directly forwards the request to the gRPC client and returns the response along with any error encountered.
func (m *MessageRpcClient) PullMessageBySeqList(ctx context.Context, req *sdkws.PullMessageBySeqsReq) (*sdkws.PullMessageBySeqsResp, error) {
	resp, err := m.Client.PullMessageBySeqs(ctx, req)
	if err != nil {
		// Wrap the error to provide more context if the gRPC call fails.
		return nil, err
	}
	return resp, nil
}

func (m *MessageRpcClient) GetConversationsHasReadAndMaxSeq(ctx context.Context, req *msg.GetConversationsHasReadAndMaxSeqReq) (*msg.GetConversationsHasReadAndMaxSeqResp, error) {
	resp, err := m.Client.GetConversationsHasReadAndMaxSeq(ctx, req)
	if err != nil {
		// Wrap the error to provide more context if the gRPC call fails.
		return nil, err
	}
	return resp, nil
}

func (m *MessageRpcClient) GetSeqMessage(ctx context.Context, req *msg.GetSeqMessageReq) (*msg.GetSeqMessageResp, error) {
	return m.Client.GetSeqMessage(ctx, req)
}

func (m *MessageRpcClient) GetConversationMaxSeq(ctx context.Context, conversationID string) (int64, error) {
	resp, err := m.Client.GetConversationMaxSeq(ctx, &msg.GetConversationMaxSeqReq{ConversationID: conversationID})
	if err != nil {
		return 0, err
	}
	return resp.MaxSeq, nil
}

func (m *MessageRpcClient) DestructMsgs(ctx context.Context, ts int64) error {
	_, err := m.Client.DestructMsgs(ctx, &msg.DestructMsgsReq{Timestamp: ts})
	return err
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

func WithRpcClient(msgRpcClient *MessageRpcClient) NotificationSenderOptions {
	return func(s *NotificationSender) {
		s.sendMsg = msgRpcClient.SendMsg
	}
}

func WithUserRpcClient(userRpcClient *UserRpcClient) NotificationSenderOptions {
	return func(s *NotificationSender) {
		s.getUserInfo = userRpcClient.GetUserInfo
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
	WithRpcGetUsername bool
}

type NotificationOptions func(*notificationOpt)

func WithRpcGetUserName() NotificationOptions {
	return func(opt *notificationOpt) {
		opt.WithRpcGetUsername = true
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
	options := config.GetOptionsByNotification(optionsConfig)
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
