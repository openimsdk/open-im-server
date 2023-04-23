package rpcclient

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/msg"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
)

type MsgClient struct {
	*MetaClient
}

func NewMsgClient(zk discoveryregistry.SvcDiscoveryRegistry) *MsgClient {
	return &MsgClient{NewMetaClient(zk, config.Config.RpcRegisterName.OpenImMsgName)}
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

func (c *MsgClient) Notification(ctx context.Context, notificationMsg *NotificationMsg) error {
	var err error
	var req msg.SendMsgReq
	var msg sdkws.MsgData
	var offlineInfo sdkws.OfflinePushInfo
	var title, desc, ex string
	var pushEnable, unReadCount bool
	msg.SendID = notificationMsg.SendID
	msg.RecvID = notificationMsg.RecvID
	msg.Content = notificationMsg.Content
	msg.MsgFrom = notificationMsg.MsgFrom
	msg.ContentType = notificationMsg.ContentType
	msg.SessionType = notificationMsg.SessionType
	msg.CreateTime = utils.GetCurrentTimestampByMill()
	msg.ClientMsgID = utils.GetMsgID(notificationMsg.SendID)
	msg.Options = make(map[string]bool, 7)
	msg.SenderNickname = notificationMsg.SenderNickname
	msg.SenderFaceURL = notificationMsg.SenderFaceURL
	switch notificationMsg.SessionType {
	case constant.GroupChatType, constant.SuperGroupChatType:
		msg.RecvID = ""
		msg.GroupID = notificationMsg.RecvID
	}
	offlineInfo.IOSBadgeCount = config.Config.IOSPush.BadgeCount
	offlineInfo.IOSPushSound = config.Config.IOSPush.PushSound
	switch msg.ContentType {
	case constant.GroupCreatedNotification:
		title = config.Config.Notification.GroupCreated.OfflinePush.Title
		desc = config.Config.Notification.GroupCreated.OfflinePush.Desc
		ex = config.Config.Notification.GroupCreated.OfflinePush.Ext
	case constant.GroupInfoSetNotification:
		title = config.Config.Notification.GroupInfoSet.OfflinePush.Title
		desc = config.Config.Notification.GroupInfoSet.OfflinePush.Desc
		ex = config.Config.Notification.GroupInfoSet.OfflinePush.Ext
	case constant.JoinGroupApplicationNotification:
		title = config.Config.Notification.JoinGroupApplication.OfflinePush.Title
		desc = config.Config.Notification.JoinGroupApplication.OfflinePush.Desc
		ex = config.Config.Notification.JoinGroupApplication.OfflinePush.Ext
	case constant.MemberQuitNotification:
		title = config.Config.Notification.MemberQuit.OfflinePush.Title
		desc = config.Config.Notification.MemberQuit.OfflinePush.Desc
		ex = config.Config.Notification.MemberQuit.OfflinePush.Ext
	case constant.GroupApplicationAcceptedNotification:
		title = config.Config.Notification.GroupApplicationAccepted.OfflinePush.Title
		desc = config.Config.Notification.GroupApplicationAccepted.OfflinePush.Desc
		ex = config.Config.Notification.GroupApplicationAccepted.OfflinePush.Ext
	case constant.GroupApplicationRejectedNotification:
		title = config.Config.Notification.GroupApplicationRejected.OfflinePush.Title
		desc = config.Config.Notification.GroupApplicationRejected.OfflinePush.Desc
		ex = config.Config.Notification.GroupApplicationRejected.OfflinePush.Ext
	case constant.GroupOwnerTransferredNotification:
		title = config.Config.Notification.GroupOwnerTransferred.OfflinePush.Title
		desc = config.Config.Notification.GroupOwnerTransferred.OfflinePush.Desc
		ex = config.Config.Notification.GroupOwnerTransferred.OfflinePush.Ext
	case constant.MemberKickedNotification:
		title = config.Config.Notification.MemberKicked.OfflinePush.Title
		desc = config.Config.Notification.MemberKicked.OfflinePush.Desc
		ex = config.Config.Notification.MemberKicked.OfflinePush.Ext
	case constant.MemberInvitedNotification:
		title = config.Config.Notification.MemberInvited.OfflinePush.Title
		desc = config.Config.Notification.MemberInvited.OfflinePush.Desc
		ex = config.Config.Notification.MemberInvited.OfflinePush.Ext
	case constant.MemberEnterNotification:
		title = config.Config.Notification.MemberEnter.OfflinePush.Title
		desc = config.Config.Notification.MemberEnter.OfflinePush.Desc
		ex = config.Config.Notification.MemberEnter.OfflinePush.Ext
	case constant.UserInfoUpdatedNotification:
		title = config.Config.Notification.UserInfoUpdated.OfflinePush.Title
		desc = config.Config.Notification.UserInfoUpdated.OfflinePush.Desc
		ex = config.Config.Notification.UserInfoUpdated.OfflinePush.Ext
	case constant.FriendApplicationNotification:
		title = config.Config.Notification.FriendApplication.OfflinePush.Title
		desc = config.Config.Notification.FriendApplication.OfflinePush.Desc
		ex = config.Config.Notification.FriendApplication.OfflinePush.Ext
	case constant.FriendApplicationApprovedNotification:
		title = config.Config.Notification.FriendApplicationApproved.OfflinePush.Title
		desc = config.Config.Notification.FriendApplicationApproved.OfflinePush.Desc
		ex = config.Config.Notification.FriendApplicationApproved.OfflinePush.Ext
	case constant.FriendApplicationRejectedNotification:
		title = config.Config.Notification.FriendApplicationRejected.OfflinePush.Title
		desc = config.Config.Notification.FriendApplicationRejected.OfflinePush.Desc
		ex = config.Config.Notification.FriendApplicationRejected.OfflinePush.Ext
	case constant.FriendAddedNotification:
		title = config.Config.Notification.FriendAdded.OfflinePush.Title
		desc = config.Config.Notification.FriendAdded.OfflinePush.Desc
		ex = config.Config.Notification.FriendAdded.OfflinePush.Ext
	case constant.FriendDeletedNotification:
		title = config.Config.Notification.FriendDeleted.OfflinePush.Title
		desc = config.Config.Notification.FriendDeleted.OfflinePush.Desc
		ex = config.Config.Notification.FriendDeleted.OfflinePush.Ext
	case constant.FriendRemarkSetNotification:
		title = config.Config.Notification.FriendRemarkSet.OfflinePush.Title
		desc = config.Config.Notification.FriendRemarkSet.OfflinePush.Desc
		ex = config.Config.Notification.FriendRemarkSet.OfflinePush.Ext
	case constant.BlackAddedNotification:
		title = config.Config.Notification.BlackAdded.OfflinePush.Title
		desc = config.Config.Notification.BlackAdded.OfflinePush.Desc
		ex = config.Config.Notification.BlackAdded.OfflinePush.Ext
	case constant.BlackDeletedNotification:
		title = config.Config.Notification.BlackDeleted.OfflinePush.Title
		desc = config.Config.Notification.BlackDeleted.OfflinePush.Desc
		ex = config.Config.Notification.BlackDeleted.OfflinePush.Ext
	case constant.ConversationOptChangeNotification:
		title = config.Config.Notification.ConversationOptUpdate.OfflinePush.Title
		desc = config.Config.Notification.ConversationOptUpdate.OfflinePush.Desc
		ex = config.Config.Notification.ConversationOptUpdate.OfflinePush.Ext

	case constant.GroupDismissedNotification:
		title = config.Config.Notification.GroupDismissed.OfflinePush.Title
		desc = config.Config.Notification.GroupDismissed.OfflinePush.Desc
		ex = config.Config.Notification.GroupDismissed.OfflinePush.Ext

	case constant.GroupMutedNotification:
		title = config.Config.Notification.GroupMuted.OfflinePush.Title
		desc = config.Config.Notification.GroupMuted.OfflinePush.Desc
		ex = config.Config.Notification.GroupMuted.OfflinePush.Ext
	case constant.GroupCancelMutedNotification:
		title = config.Config.Notification.GroupCancelMuted.OfflinePush.Title
		desc = config.Config.Notification.GroupCancelMuted.OfflinePush.Desc
		ex = config.Config.Notification.GroupCancelMuted.OfflinePush.Ext
	case constant.GroupMemberMutedNotification:
		title = config.Config.Notification.GroupMemberMuted.OfflinePush.Title
		desc = config.Config.Notification.GroupMemberMuted.OfflinePush.Desc
		ex = config.Config.Notification.GroupMemberMuted.OfflinePush.Ext
	case constant.GroupMemberCancelMutedNotification:
		title = config.Config.Notification.GroupMemberCancelMuted.OfflinePush.Title
		desc = config.Config.Notification.GroupMemberCancelMuted.OfflinePush.Desc
		ex = config.Config.Notification.GroupMemberCancelMuted.OfflinePush.Ext
	case constant.GroupMemberInfoSetNotification:
		title = config.Config.Notification.GroupMemberInfoSet.OfflinePush.Title
		desc = config.Config.Notification.GroupMemberInfoSet.OfflinePush.Desc
		ex = config.Config.Notification.GroupMemberInfoSet.OfflinePush.Ext
	case constant.ConversationPrivateChatNotification:
		title = config.Config.Notification.ConversationSetPrivate.OfflinePush.Title
		desc = config.Config.Notification.ConversationSetPrivate.OfflinePush.Desc
		ex = config.Config.Notification.ConversationSetPrivate.OfflinePush.Ext
	case constant.FriendInfoUpdatedNotification:
		title = config.Config.Notification.FriendInfoUpdated.OfflinePush.Title
		desc = config.Config.Notification.FriendInfoUpdated.OfflinePush.Desc
		ex = config.Config.Notification.FriendInfoUpdated.OfflinePush.Ext
	case constant.DeleteMessageNotification:
	case constant.ConversationUnreadNotification, constant.SuperGroupUpdateNotification:
	}
	utils.SetSwitchFromOptions(msg.Options, constant.IsUnreadCount, unReadCount)
	utils.SetSwitchFromOptions(msg.Options, constant.IsOfflinePush, pushEnable)
	offlineInfo.Title = title
	offlineInfo.Desc = desc
	offlineInfo.Ex = ex
	msg.OfflinePushInfo = &offlineInfo
	req.MsgData = &msg
	_, err = c.SendMsg(ctx, &req)
	return err
}
