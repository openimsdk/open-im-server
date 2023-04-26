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

	"github.com/golang/protobuf/proto"
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
		constant.FriendApplicationNotification:         config.Config.Notification.FriendApplication,
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

type MsgClient struct {
	*MetaClient
	contentTypeConf map[int32]config.NotificationConf
}

func NewMsgClient(zk discoveryregistry.SvcDiscoveryRegistry) *MsgClient {
	return &MsgClient{MetaClient: NewMetaClient(zk, config.Config.RpcRegisterName.OpenImMsgName), contentTypeConf: newContentTypeConf()}
}

func (m *MsgClient) SendMsg(ctx context.Context, req *msg.SendMsgReq) (*msg.SendMsgResp, error) {
	cc, err := m.getConn()
	if err != nil {
		return nil, err
	}
	resp, err := msg.NewMsgClient(cc).SendMsg(ctx, req)
	return resp, err
}

func (m *MsgClient) GetMaxAndMinSeq(ctx context.Context, req *sdkws.GetMaxAndMinSeqReq) (*sdkws.GetMaxAndMinSeqResp, error) {
	cc, err := m.getConn()
	if err != nil {
		return nil, err
	}
	resp, err := msg.NewMsgClient(cc).GetMaxAndMinSeq(ctx, req)
	return resp, err
}

func (m *MsgClient) PullMessageBySeqList(ctx context.Context, req *sdkws.PullMessageBySeqsReq) (*sdkws.PullMessageBySeqsResp, error) {
	cc, err := m.getConn()
	if err != nil {
		return nil, err
	}
	resp, err := msg.NewMsgClient(cc).PullMessageBySeqs(ctx, req)
	return resp, err
}

func (c *MsgClient) Notification(ctx context.Context, sendID, recvID string, contentType, sessionType int32, m proto.Message, opts ...utils.OptionsOpt) error {
	content, err := json.Marshal(m)
	if err != nil {
		log.ZError(ctx, "MsgClient Notification json.Marshal failed", err)
		return err
	}
	var req msg.SendMsgReq
	var msg sdkws.MsgData
	var offlineInfo sdkws.OfflinePushInfo
	var title, desc, ex string
	msg.SendID = sendID
	msg.RecvID = recvID
	if sessionType == constant.SuperGroupChatType {
		msg.GroupID = recvID
	}
	msg.Content = content
	msg.MsgFrom = constant.SysMsgType
	msg.ContentType = contentType
	msg.SessionType = sessionType
	msg.CreateTime = utils.GetCurrentTimestampByMill()
	msg.ClientMsgID = utils.GetMsgID(sendID)
	// msg.Options = make(map[string]bool, 7)
	// todo notification get sender name and face url
	// msg.SenderNickname, msg.SenderFaceURL, err = c.getFaceURLAndName(sendID)
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
