package rpcclient

import (
	"context"
	"encoding/json"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/msg"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"

	"google.golang.org/protobuf/proto"
	// "google.golang.org/protobuf/proto"
)

func newContentTypeConf() map[int32]config.NotificationConf {
	return map[int32]config.NotificationConf{
		// group
		constant.GroupCreatedNotification:                 config.Config.Notification.GroupCreated,
		constant.GroupInfoSetNotification:                 config.Config.Notification.GroupInfoSet,
		constant.JoinGroupApplicationNotification:         config.Config.Notification.JoinGroupApplication,
		constant.MemberQuitNotification:                   config.Config.Notification.MemberQuit,
		constant.GroupApplicationAcceptedNotification:     config.Config.Notification.GroupApplicationAccepted,
		constant.GroupApplicationRejectedNotification:     config.Config.Notification.GroupApplicationRejected,
		constant.GroupOwnerTransferredNotification:        config.Config.Notification.GroupOwnerTransferred,
		constant.MemberKickedNotification:                 config.Config.Notification.MemberKicked,
		constant.MemberInvitedNotification:                config.Config.Notification.MemberInvited,
		constant.MemberEnterNotification:                  config.Config.Notification.MemberEnter,
		constant.GroupDismissedNotification:               config.Config.Notification.GroupDismissed,
		constant.GroupMutedNotification:                   config.Config.Notification.GroupMuted,
		constant.GroupCancelMutedNotification:             config.Config.Notification.GroupCancelMuted,
		constant.GroupMemberMutedNotification:             config.Config.Notification.GroupMemberMuted,
		constant.GroupMemberCancelMutedNotification:       config.Config.Notification.GroupMemberCancelMuted,
		constant.GroupMemberInfoSetNotification:           config.Config.Notification.GroupMemberInfoSet,
		constant.GroupMemberSetToAdminNotification:        config.Config.Notification.GroupMemberSetToAdmin,
		constant.GroupMemberSetToOrdinaryUserNotification: config.Config.Notification.GroupMemberSetToOrdinary,
		// user
		constant.UserInfoUpdatedNotification: config.Config.Notification.UserInfoUpdated,
		// friend
		constant.FriendApplicationNotification:         config.Config.Notification.FriendApplicationAdded,
		constant.FriendApplicationApprovedNotification: config.Config.Notification.FriendApplicationApproved,
		constant.FriendApplicationRejectedNotification: config.Config.Notification.FriendApplicationRejected,
		constant.FriendAddedNotification:               config.Config.Notification.FriendAdded,
		constant.FriendDeletedNotification:             config.Config.Notification.FriendDeleted,
		constant.FriendRemarkSetNotification:           config.Config.Notification.FriendRemarkSet,
		constant.BlackAddedNotification:                config.Config.Notification.BlackAdded,
		constant.BlackDeletedNotification:              config.Config.Notification.BlackDeleted,
		constant.FriendInfoUpdatedNotification:         config.Config.Notification.FriendInfoUpdated,
		// conversation
		constant.ConversationChangeNotification:      config.Config.Notification.ConversationChanged,
		constant.ConversationUnreadNotification:      config.Config.Notification.ConversationChanged,
		constant.ConversationPrivateChatNotification: config.Config.Notification.ConversationSetPrivate,
	}
}

func newSessionTypeConf() map[int32]int32 {
	return map[int32]int32{
		// group
		constant.GroupCreatedNotification:                 constant.SuperGroupChatType,
		constant.GroupInfoSetNotification:                 constant.SuperGroupChatType,
		constant.JoinGroupApplicationNotification:         constant.SuperGroupChatType,
		constant.MemberQuitNotification:                   constant.SingleChatType,
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
		// user
		constant.UserInfoUpdatedNotification: constant.SingleChatType,
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
		// conversation
		constant.ConversationChangeNotification:      constant.SingleChatType,
		constant.ConversationUnreadNotification:      constant.SingleChatType,
		constant.ConversationPrivateChatNotification: constant.SingleChatType,
	}
}

type MsgClient struct {
	*MetaClient
	contentTypeConf map[int32]config.NotificationConf
	sessionTypeConf map[int32]int32
}

func NewMsgClient(zk discoveryregistry.SvcDiscoveryRegistry) *MsgClient {
	return &MsgClient{MetaClient: NewMetaClient(zk, config.Config.RpcRegisterName.OpenImMsgName), contentTypeConf: newContentTypeConf(), sessionTypeConf: newSessionTypeConf()}
}

func (m *MsgClient) SendMsg(ctx context.Context, req *msg.SendMsgReq) (*msg.SendMsgResp, error) {
	cc, err := m.getConn(ctx)
	if err != nil {
		return nil, err
	}
	resp, err := msg.NewMsgClient(cc).SendMsg(ctx, req)
	return resp, err
}

func (m *MsgClient) GetMaxSeq(ctx context.Context, req *sdkws.GetMaxSeqReq) (*sdkws.GetMaxSeqResp, error) {
	cc, err := m.getConn(ctx)
	if err != nil {
		return nil, err
	}
	resp, err := msg.NewMsgClient(cc).GetMaxSeq(ctx, req)
	return resp, err
}

func (m *MsgClient) PullMessageBySeqList(ctx context.Context, req *sdkws.PullMessageBySeqsReq) (*sdkws.PullMessageBySeqsResp, error) {
	cc, err := m.getConn(ctx)
	if err != nil {
		return nil, err
	}
	resp, err := msg.NewMsgClient(cc).PullMessageBySeqs(ctx, req)
	return resp, err
}

type NotificationElem struct {
	Detail interface{} `json:"detail,omitempty"`
}

func (c *MsgClient) Notification(ctx context.Context, sendID, recvID string, contentType int32, m proto.Message, opts ...utils.OptionsOpt) error {
	n := NotificationElem{Detail: m}
	content, err := json.Marshal(n)
	if err != nil {
		log.ZError(ctx, "MsgClient Notification json.Marshal failed", err, "sendID", sendID, "recvID", recvID, "contentType", contentType, "msg", m)
		return err
	}
	var req msg.SendMsgReq
	var msg sdkws.MsgData
	var offlineInfo sdkws.OfflinePushInfo
	var title, desc, ex string
	msg.SendID = sendID
	msg.RecvID = recvID
	msg.Content = content
	msg.MsgFrom = constant.SysMsgType
	msg.ContentType = contentType
	msg.SessionType = c.sessionTypeConf[contentType]
	if msg.SessionType == constant.SuperGroupChatType {
		msg.GroupID = recvID
	}
	msg.CreateTime = utils.GetCurrentTimestampByMill()
	msg.ClientMsgID = utils.GetMsgID(sendID)
	options := config.GetOptionsByNotification(c.contentTypeConf[contentType])
	options = utils.WithOptions(options, opts...)
	msg.Options = options
	offlineInfo.Title = title
	offlineInfo.Desc = desc
	offlineInfo.Ex = ex
	msg.OfflinePushInfo = &offlineInfo
	req.MsgData = &msg
	_, err = c.SendMsg(ctx, &req)
	if err == nil {
		log.ZDebug(ctx, "MsgClient Notification SendMsg success", "req", &req)
	} else {
		log.ZError(ctx, "MsgClient Notification SendMsg failed %s\n", err, "req", &req)
	}
	return err
}
