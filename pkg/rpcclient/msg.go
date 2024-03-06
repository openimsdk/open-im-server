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
	"fmt"

	"github.com/OpenIMSDK/protocol/constant"
	"github.com/OpenIMSDK/protocol/msg"
	"github.com/OpenIMSDK/protocol/sdkws"
	"github.com/OpenIMSDK/tools/discoveryregistry"
	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/log"
	"github.com/OpenIMSDK/tools/utils"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
)

func newContentTypeConf(conf *config.GlobalConfig) map[int32]config.NotificationConf {
	return map[int32]config.NotificationConf{
		// group
		constant.GroupCreatedNotification:                 conf.Notification.GroupCreated,
		constant.GroupInfoSetNotification:                 conf.Notification.GroupInfoSet,
		constant.JoinGroupApplicationNotification:         conf.Notification.JoinGroupApplication,
		constant.MemberQuitNotification:                   conf.Notification.MemberQuit,
		constant.GroupApplicationAcceptedNotification:     conf.Notification.GroupApplicationAccepted,
		constant.GroupApplicationRejectedNotification:     conf.Notification.GroupApplicationRejected,
		constant.GroupOwnerTransferredNotification:        conf.Notification.GroupOwnerTransferred,
		constant.MemberKickedNotification:                 conf.Notification.MemberKicked,
		constant.MemberInvitedNotification:                conf.Notification.MemberInvited,
		constant.MemberEnterNotification:                  conf.Notification.MemberEnter,
		constant.GroupDismissedNotification:               conf.Notification.GroupDismissed,
		constant.GroupMutedNotification:                   conf.Notification.GroupMuted,
		constant.GroupCancelMutedNotification:             conf.Notification.GroupCancelMuted,
		constant.GroupMemberMutedNotification:             conf.Notification.GroupMemberMuted,
		constant.GroupMemberCancelMutedNotification:       conf.Notification.GroupMemberCancelMuted,
		constant.GroupMemberInfoSetNotification:           conf.Notification.GroupMemberInfoSet,
		constant.GroupMemberSetToAdminNotification:        conf.Notification.GroupMemberSetToAdmin,
		constant.GroupMemberSetToOrdinaryUserNotification: conf.Notification.GroupMemberSetToOrdinary,
		constant.GroupInfoSetAnnouncementNotification:     conf.Notification.GroupInfoSetAnnouncement,
		constant.GroupInfoSetNameNotification:             conf.Notification.GroupInfoSetName,
		// user
		constant.UserInfoUpdatedNotification:  conf.Notification.UserInfoUpdated,
		constant.UserStatusChangeNotification: conf.Notification.UserStatusChanged,
		// friend
		constant.FriendApplicationNotification:         conf.Notification.FriendApplicationAdded,
		constant.FriendApplicationApprovedNotification: conf.Notification.FriendApplicationApproved,
		constant.FriendApplicationRejectedNotification: conf.Notification.FriendApplicationRejected,
		constant.FriendAddedNotification:               conf.Notification.FriendAdded,
		constant.FriendDeletedNotification:             conf.Notification.FriendDeleted,
		constant.FriendRemarkSetNotification:           conf.Notification.FriendRemarkSet,
		constant.BlackAddedNotification:                conf.Notification.BlackAdded,
		constant.BlackDeletedNotification:              conf.Notification.BlackDeleted,
		constant.FriendInfoUpdatedNotification:         conf.Notification.FriendInfoUpdated,
		constant.FriendsInfoUpdateNotification:         conf.Notification.FriendInfoUpdated, //use the same FriendInfoUpdated
		// conversation
		constant.ConversationChangeNotification:      conf.Notification.ConversationChanged,
		constant.ConversationUnreadNotification:      conf.Notification.ConversationChanged,
		constant.ConversationPrivateChatNotification: conf.Notification.ConversationSetPrivate,
		// msg
		constant.MsgRevokeNotification:  {IsSendMsg: false, ReliabilityLevel: constant.ReliableNotificationNoMsg},
		constant.HasReadReceipt:         {IsSendMsg: false, ReliabilityLevel: constant.ReliableNotificationNoMsg},
		constant.DeleteMsgsNotification: {IsSendMsg: false, ReliabilityLevel: constant.ReliableNotificationNoMsg},
	}
}

func newSessionTypeConf() map[int32]int32 {
	return map[int32]int32{
		// group
		constant.GroupCreatedNotification:                 constant.SuperGroupChatType,
		constant.GroupInfoSetNotification:                 constant.SuperGroupChatType,
		constant.JoinGroupApplicationNotification:         constant.SingleChatType,
		constant.MemberQuitNotification:                   constant.SuperGroupChatType,
		constant.GroupApplicationAcceptedNotification:     constant.SingleChatType,
		constant.GroupApplicationRejectedNotification:     constant.SingleChatType,
		constant.GroupOwnerTransferredNotification:        constant.SuperGroupChatType,
		constant.MemberKickedNotification:                 constant.SuperGroupChatType,
		constant.MemberInvitedNotification:                constant.SuperGroupChatType,
		constant.MemberEnterNotification:                  constant.SuperGroupChatType,
		constant.GroupDismissedNotification:               constant.SuperGroupChatType,
		constant.GroupMutedNotification:                   constant.SuperGroupChatType,
		constant.GroupCancelMutedNotification:             constant.SuperGroupChatType,
		constant.GroupMemberMutedNotification:             constant.SuperGroupChatType,
		constant.GroupMemberCancelMutedNotification:       constant.SuperGroupChatType,
		constant.GroupMemberInfoSetNotification:           constant.SuperGroupChatType,
		constant.GroupMemberSetToAdminNotification:        constant.SuperGroupChatType,
		constant.GroupMemberSetToOrdinaryUserNotification: constant.SuperGroupChatType,
		constant.GroupInfoSetAnnouncementNotification:     constant.SuperGroupChatType,
		constant.GroupInfoSetNameNotification:             constant.SuperGroupChatType,
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
	discov discoveryregistry.SvcDiscoveryRegistry
	Config *config.GlobalConfig
}

func NewMessage(discov discoveryregistry.SvcDiscoveryRegistry, config *config.GlobalConfig) *Message {
	conn, err := discov.GetConn(context.Background(), config.RpcRegisterName.OpenImMsgName)
	if err != nil {
		panic(err)
	}
	client := msg.NewMsgClient(conn)
	return &Message{discov: discov, conn: conn, Client: client, Config: config}
}

type MessageRpcClient Message

func NewMessageRpcClient(discov discoveryregistry.SvcDiscoveryRegistry, config *config.GlobalConfig) MessageRpcClient {
	return MessageRpcClient(*NewMessage(discov, config))
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
	return resp.MaxSeqs, err
}

func (m *MessageRpcClient) GetHasReadSeqs(ctx context.Context, userID string, conversationIDs []string) (map[string]int64, error) {
	resp, err := m.Client.GetHasReadSeqs(ctx, &msg.GetHasReadSeqsReq{
		UserID:          userID,
		ConversationIDs: conversationIDs,
	})
	return resp.MaxSeqs, err
}

func (m *MessageRpcClient) GetMsgByConversationIDs(ctx context.Context, docIDs []string, seqs map[string]int64) (map[string]*sdkws.MsgData, error) {
	resp, err := m.Client.GetMsgByConversationIDs(ctx, &msg.GetMsgByConversationIDsReq{
		ConversationIDs: docIDs,
		MaxSeqs:         seqs,
	})
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

func (m *MessageRpcClient) GetConversationMaxSeq(ctx context.Context, conversationID string) (int64, error) {
	resp, err := m.Client.GetConversationMaxSeq(ctx, &msg.GetConversationMaxSeqReq{ConversationID: conversationID})
	if err != nil {
		return 0, err
	}
	return resp.MaxSeq, nil
}

type NotificationSender struct {
	contentTypeConf map[int32]config.NotificationConf
	sessionTypeConf map[int32]int32
	sendMsg         func(ctx context.Context, req *msg.SendMsgReq) (*msg.SendMsgResp, error)
	getUserInfo     func(ctx context.Context, userID string) (*sdkws.UserInfo, error)
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

func NewNotificationSender(config *config.GlobalConfig, opts ...NotificationSenderOptions) *NotificationSender {
	notificationSender := &NotificationSender{contentTypeConf: newContentTypeConf(config), sessionTypeConf: newSessionTypeConf()}
	for _, opt := range opts {
		opt(notificationSender)
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

func (s *NotificationSender) NotificationWithSesstionType(ctx context.Context, sendID, recvID string, contentType, sesstionType int32, m proto.Message, opts ...NotificationOptions) (err error) {
	n := sdkws.NotificationElem{Detail: utils.StructToJsonString(m)}
	content, err := json.Marshal(&n)
	if err != nil {
		errInfo := fmt.Sprintf("MsgClient Notification json.Marshal failed, sendID:%s, recvID:%s, contentType:%d, msg:%s", sendID, recvID, contentType, m)
		return errs.Wrap(err, errInfo)
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
			errInfo := fmt.Sprintf("getUserInfo failed, sendID:%s", sendID)
			return errs.Wrap(err, errInfo)
		} else {
			msg.SenderNickname = userInfo.Nickname
			msg.SenderFaceURL = userInfo.FaceURL
		}
	}
	var offlineInfo sdkws.OfflinePushInfo
	var title, desc, ex string
	msg.SendID = sendID
	msg.RecvID = recvID
	msg.Content = content
	msg.MsgFrom = constant.SysMsgType
	msg.ContentType = contentType
	msg.SessionType = sesstionType
	if msg.SessionType == constant.SuperGroupChatType {
		msg.GroupID = recvID
	}
	msg.CreateTime = utils.GetCurrentTimestampByMill()
	msg.ClientMsgID = utils.GetMsgID(sendID)
	optionsConfig := s.contentTypeConf[contentType]
	if sendID == recvID && contentType == constant.HasReadReceipt {
		optionsConfig.ReliabilityLevel = constant.UnreliableNotification
	}
	options := config.GetOptionsByNotification(optionsConfig)
	s.SetOptionsByContentType(ctx, options, contentType)
	msg.Options = options
	offlineInfo.Title = title
	offlineInfo.Desc = desc
	offlineInfo.Ex = ex
	msg.OfflinePushInfo = &offlineInfo
	req.MsgData = &msg
	_, err = s.sendMsg(ctx, &req)
	if err != nil {
		errInfo := fmt.Sprintf("MsgClient Notification SendMsg failed, req:%s", &req)
		return errs.Wrap(err, errInfo)
	}
	return err
}

func (s *NotificationSender) Notification(ctx context.Context, sendID, recvID string, contentType int32, m proto.Message, opts ...NotificationOptions) error {
	return s.NotificationWithSesstionType(ctx, sendID, recvID, contentType, s.sessionTypeConf[contentType], m, opts...)
}

func (s *NotificationSender) SetOptionsByContentType(_ context.Context, options map[string]bool, contentType int32) {
	switch contentType {
	case constant.UserStatusChangeNotification:
		options[constant.IsSenderSync] = false
	default:
	}
}
