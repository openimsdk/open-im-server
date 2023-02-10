package notification

import (
	"Open_IM/pkg/common/constant"
	pbChat "Open_IM/pkg/proto/msg"
	"strings"
)

type NotificationMsg struct {
	SendID         string
	RecvID         string
	Content        []byte //  sdkws.TipsComm
	MsgFrom        int32
	ContentType    int32
	SessionType    int32
	OperationID    string
	SenderNickname string
	SenderFaceURL  string
}

func Notification(n *NotificationMsg) {
	var req pbChat.SendMsgReq
	var msg sdkws.MsgData
	var offlineInfo sdkws.OfflinePushInfo
	var title, desc, ex string
	var pushSwitch, unReadCount bool
	var reliabilityLevel int
	req.OperationID = n.OperationID
	msg.SendID = n.SendID
	msg.RecvID = n.RecvID
	msg.Content = n.Content
	msg.MsgFrom = n.MsgFrom
	msg.ContentType = n.ContentType
	msg.SessionType = n.SessionType
	msg.CreateTime = utils.GetCurrentTimestampByMill()
	msg.ClientMsgID = utils.GetMsgID(n.SendID)
	msg.Options = make(map[string]bool, 7)
	msg.SenderNickname = n.SenderNickname
	msg.SenderFaceURL = n.SenderFaceURL
	switch n.SessionType {
	case constant.GroupChatType, constant.SuperGroupChatType:
		msg.RecvID = ""
		msg.GroupID = n.RecvID
	}
	offlineInfo.IOSBadgeCount = config.Config.IOSPush.BadgeCount
	offlineInfo.IOSPushSound = config.Config.IOSPush.PushSound
	switch msg.ContentType {
	case constant.GroupCreatedNotification:
		pushSwitch = config.Config.Notification.GroupCreated.OfflinePush.PushSwitch
		title = config.Config.Notification.GroupCreated.OfflinePush.Title
		desc = config.Config.Notification.GroupCreated.OfflinePush.Desc
		ex = config.Config.Notification.GroupCreated.OfflinePush.Ext
		reliabilityLevel = config.Config.Notification.GroupCreated.Conversation.ReliabilityLevel
		unReadCount = config.Config.Notification.GroupCreated.Conversation.UnreadCount
	case constant.GroupInfoSetNotification:
		pushSwitch = config.Config.Notification.GroupInfoSet.OfflinePush.PushSwitch
		title = config.Config.Notification.GroupInfoSet.OfflinePush.Title
		desc = config.Config.Notification.GroupInfoSet.OfflinePush.Desc
		ex = config.Config.Notification.GroupInfoSet.OfflinePush.Ext
		reliabilityLevel = config.Config.Notification.GroupInfoSet.Conversation.ReliabilityLevel
		unReadCount = config.Config.Notification.GroupInfoSet.Conversation.UnreadCount
	case constant.JoinGroupApplicationNotification:
		pushSwitch = config.Config.Notification.JoinGroupApplication.OfflinePush.PushSwitch
		title = config.Config.Notification.JoinGroupApplication.OfflinePush.Title
		desc = config.Config.Notification.JoinGroupApplication.OfflinePush.Desc
		ex = config.Config.Notification.JoinGroupApplication.OfflinePush.Ext
		reliabilityLevel = config.Config.Notification.JoinGroupApplication.Conversation.ReliabilityLevel
		unReadCount = config.Config.Notification.JoinGroupApplication.Conversation.UnreadCount
	case constant.MemberQuitNotification:
		pushSwitch = config.Config.Notification.MemberQuit.OfflinePush.PushSwitch
		title = config.Config.Notification.MemberQuit.OfflinePush.Title
		desc = config.Config.Notification.MemberQuit.OfflinePush.Desc
		ex = config.Config.Notification.MemberQuit.OfflinePush.Ext
		reliabilityLevel = config.Config.Notification.MemberQuit.Conversation.ReliabilityLevel
		unReadCount = config.Config.Notification.MemberQuit.Conversation.UnreadCount
	case constant.GroupApplicationAcceptedNotification:
		pushSwitch = config.Config.Notification.GroupApplicationAccepted.OfflinePush.PushSwitch
		title = config.Config.Notification.GroupApplicationAccepted.OfflinePush.Title
		desc = config.Config.Notification.GroupApplicationAccepted.OfflinePush.Desc
		ex = config.Config.Notification.GroupApplicationAccepted.OfflinePush.Ext
		reliabilityLevel = config.Config.Notification.GroupApplicationAccepted.Conversation.ReliabilityLevel
		unReadCount = config.Config.Notification.GroupApplicationAccepted.Conversation.UnreadCount
	case constant.GroupApplicationRejectedNotification:
		pushSwitch = config.Config.Notification.GroupApplicationRejected.OfflinePush.PushSwitch
		title = config.Config.Notification.GroupApplicationRejected.OfflinePush.Title
		desc = config.Config.Notification.GroupApplicationRejected.OfflinePush.Desc
		ex = config.Config.Notification.GroupApplicationRejected.OfflinePush.Ext
		reliabilityLevel = config.Config.Notification.GroupApplicationRejected.Conversation.ReliabilityLevel
		unReadCount = config.Config.Notification.GroupApplicationRejected.Conversation.UnreadCount
	case constant.GroupOwnerTransferredNotification:
		pushSwitch = config.Config.Notification.GroupOwnerTransferred.OfflinePush.PushSwitch
		title = config.Config.Notification.GroupOwnerTransferred.OfflinePush.Title
		desc = config.Config.Notification.GroupOwnerTransferred.OfflinePush.Desc
		ex = config.Config.Notification.GroupOwnerTransferred.OfflinePush.Ext
		reliabilityLevel = config.Config.Notification.GroupOwnerTransferred.Conversation.ReliabilityLevel
		unReadCount = config.Config.Notification.GroupOwnerTransferred.Conversation.UnreadCount
	case constant.MemberKickedNotification:
		pushSwitch = config.Config.Notification.MemberKicked.OfflinePush.PushSwitch
		title = config.Config.Notification.MemberKicked.OfflinePush.Title
		desc = config.Config.Notification.MemberKicked.OfflinePush.Desc
		ex = config.Config.Notification.MemberKicked.OfflinePush.Ext
		reliabilityLevel = config.Config.Notification.MemberKicked.Conversation.ReliabilityLevel
		unReadCount = config.Config.Notification.MemberKicked.Conversation.UnreadCount
	case constant.MemberInvitedNotification:
		pushSwitch = config.Config.Notification.MemberInvited.OfflinePush.PushSwitch
		title = config.Config.Notification.MemberInvited.OfflinePush.Title
		desc = config.Config.Notification.MemberInvited.OfflinePush.Desc
		ex = config.Config.Notification.MemberInvited.OfflinePush.Ext
		reliabilityLevel = config.Config.Notification.MemberInvited.Conversation.ReliabilityLevel
		unReadCount = config.Config.Notification.MemberInvited.Conversation.UnreadCount
	case constant.MemberEnterNotification:
		pushSwitch = config.Config.Notification.MemberEnter.OfflinePush.PushSwitch
		title = config.Config.Notification.MemberEnter.OfflinePush.Title
		desc = config.Config.Notification.MemberEnter.OfflinePush.Desc
		ex = config.Config.Notification.MemberEnter.OfflinePush.Ext
		reliabilityLevel = config.Config.Notification.MemberEnter.Conversation.ReliabilityLevel
		unReadCount = config.Config.Notification.MemberEnter.Conversation.UnreadCount
	case constant.UserInfoUpdatedNotification:
		pushSwitch = config.Config.Notification.UserInfoUpdated.OfflinePush.PushSwitch
		title = config.Config.Notification.UserInfoUpdated.OfflinePush.Title
		desc = config.Config.Notification.UserInfoUpdated.OfflinePush.Desc
		ex = config.Config.Notification.UserInfoUpdated.OfflinePush.Ext
		reliabilityLevel = config.Config.Notification.UserInfoUpdated.Conversation.ReliabilityLevel
		unReadCount = config.Config.Notification.UserInfoUpdated.Conversation.UnreadCount
	case constant.FriendApplicationNotification:
		pushSwitch = config.Config.Notification.FriendApplication.OfflinePush.PushSwitch
		title = config.Config.Notification.FriendApplication.OfflinePush.Title
		desc = config.Config.Notification.FriendApplication.OfflinePush.Desc
		ex = config.Config.Notification.FriendApplication.OfflinePush.Ext
		reliabilityLevel = config.Config.Notification.FriendApplication.Conversation.ReliabilityLevel
		unReadCount = config.Config.Notification.FriendApplication.Conversation.UnreadCount
	case constant.FriendApplicationApprovedNotification:
		pushSwitch = config.Config.Notification.FriendApplicationApproved.OfflinePush.PushSwitch
		title = config.Config.Notification.FriendApplicationApproved.OfflinePush.Title
		desc = config.Config.Notification.FriendApplicationApproved.OfflinePush.Desc
		ex = config.Config.Notification.FriendApplicationApproved.OfflinePush.Ext
		reliabilityLevel = config.Config.Notification.FriendApplicationApproved.Conversation.ReliabilityLevel
		unReadCount = config.Config.Notification.FriendApplicationApproved.Conversation.UnreadCount
	case constant.FriendApplicationRejectedNotification:
		pushSwitch = config.Config.Notification.FriendApplicationRejected.OfflinePush.PushSwitch
		title = config.Config.Notification.FriendApplicationRejected.OfflinePush.Title
		desc = config.Config.Notification.FriendApplicationRejected.OfflinePush.Desc
		ex = config.Config.Notification.FriendApplicationRejected.OfflinePush.Ext
		reliabilityLevel = config.Config.Notification.FriendApplicationRejected.Conversation.ReliabilityLevel
		unReadCount = config.Config.Notification.FriendApplicationRejected.Conversation.UnreadCount
	case constant.FriendAddedNotification:
		pushSwitch = config.Config.Notification.FriendAdded.OfflinePush.PushSwitch
		title = config.Config.Notification.FriendAdded.OfflinePush.Title
		desc = config.Config.Notification.FriendAdded.OfflinePush.Desc
		ex = config.Config.Notification.FriendAdded.OfflinePush.Ext
		reliabilityLevel = config.Config.Notification.FriendAdded.Conversation.ReliabilityLevel
		unReadCount = config.Config.Notification.FriendAdded.Conversation.UnreadCount
	case constant.FriendDeletedNotification:
		pushSwitch = config.Config.Notification.FriendDeleted.OfflinePush.PushSwitch
		title = config.Config.Notification.FriendDeleted.OfflinePush.Title
		desc = config.Config.Notification.FriendDeleted.OfflinePush.Desc
		ex = config.Config.Notification.FriendDeleted.OfflinePush.Ext
		reliabilityLevel = config.Config.Notification.FriendDeleted.Conversation.ReliabilityLevel
		unReadCount = config.Config.Notification.FriendDeleted.Conversation.UnreadCount
	case constant.FriendRemarkSetNotification:
		pushSwitch = config.Config.Notification.FriendRemarkSet.OfflinePush.PushSwitch
		title = config.Config.Notification.FriendRemarkSet.OfflinePush.Title
		desc = config.Config.Notification.FriendRemarkSet.OfflinePush.Desc
		ex = config.Config.Notification.FriendRemarkSet.OfflinePush.Ext
		reliabilityLevel = config.Config.Notification.FriendRemarkSet.Conversation.ReliabilityLevel
		unReadCount = config.Config.Notification.FriendRemarkSet.Conversation.UnreadCount
	case constant.BlackAddedNotification:
		pushSwitch = config.Config.Notification.BlackAdded.OfflinePush.PushSwitch
		title = config.Config.Notification.BlackAdded.OfflinePush.Title
		desc = config.Config.Notification.BlackAdded.OfflinePush.Desc
		ex = config.Config.Notification.BlackAdded.OfflinePush.Ext
		reliabilityLevel = config.Config.Notification.BlackAdded.Conversation.ReliabilityLevel
		unReadCount = config.Config.Notification.BlackAdded.Conversation.UnreadCount
	case constant.BlackDeletedNotification:
		pushSwitch = config.Config.Notification.BlackDeleted.OfflinePush.PushSwitch
		title = config.Config.Notification.BlackDeleted.OfflinePush.Title
		desc = config.Config.Notification.BlackDeleted.OfflinePush.Desc
		ex = config.Config.Notification.BlackDeleted.OfflinePush.Ext
		reliabilityLevel = config.Config.Notification.BlackDeleted.Conversation.ReliabilityLevel
		unReadCount = config.Config.Notification.BlackDeleted.Conversation.UnreadCount
	case constant.ConversationOptChangeNotification:
		pushSwitch = config.Config.Notification.ConversationOptUpdate.OfflinePush.PushSwitch
		title = config.Config.Notification.ConversationOptUpdate.OfflinePush.Title
		desc = config.Config.Notification.ConversationOptUpdate.OfflinePush.Desc
		ex = config.Config.Notification.ConversationOptUpdate.OfflinePush.Ext
		reliabilityLevel = config.Config.Notification.ConversationOptUpdate.Conversation.ReliabilityLevel
		unReadCount = config.Config.Notification.ConversationOptUpdate.Conversation.UnreadCount

	case constant.GroupDismissedNotification:
		pushSwitch = config.Config.Notification.GroupDismissed.OfflinePush.PushSwitch
		title = config.Config.Notification.GroupDismissed.OfflinePush.Title
		desc = config.Config.Notification.GroupDismissed.OfflinePush.Desc
		ex = config.Config.Notification.GroupDismissed.OfflinePush.Ext
		reliabilityLevel = config.Config.Notification.GroupDismissed.Conversation.ReliabilityLevel
		unReadCount = config.Config.Notification.GroupDismissed.Conversation.UnreadCount

	case constant.GroupMutedNotification:
		pushSwitch = config.Config.Notification.GroupMuted.OfflinePush.PushSwitch
		title = config.Config.Notification.GroupMuted.OfflinePush.Title
		desc = config.Config.Notification.GroupMuted.OfflinePush.Desc
		ex = config.Config.Notification.GroupMuted.OfflinePush.Ext
		reliabilityLevel = config.Config.Notification.GroupMuted.Conversation.ReliabilityLevel
		unReadCount = config.Config.Notification.GroupMuted.Conversation.UnreadCount

	case constant.GroupCancelMutedNotification:
		pushSwitch = config.Config.Notification.GroupCancelMuted.OfflinePush.PushSwitch
		title = config.Config.Notification.GroupCancelMuted.OfflinePush.Title
		desc = config.Config.Notification.GroupCancelMuted.OfflinePush.Desc
		ex = config.Config.Notification.GroupCancelMuted.OfflinePush.Ext
		reliabilityLevel = config.Config.Notification.GroupCancelMuted.Conversation.ReliabilityLevel
		unReadCount = config.Config.Notification.GroupCancelMuted.Conversation.UnreadCount

	case constant.GroupMemberMutedNotification:
		pushSwitch = config.Config.Notification.GroupMemberMuted.OfflinePush.PushSwitch
		title = config.Config.Notification.GroupMemberMuted.OfflinePush.Title
		desc = config.Config.Notification.GroupMemberMuted.OfflinePush.Desc
		ex = config.Config.Notification.GroupMemberMuted.OfflinePush.Ext
		reliabilityLevel = config.Config.Notification.GroupMemberMuted.Conversation.ReliabilityLevel
		unReadCount = config.Config.Notification.GroupMemberMuted.Conversation.UnreadCount

	case constant.GroupMemberCancelMutedNotification:
		pushSwitch = config.Config.Notification.GroupMemberCancelMuted.OfflinePush.PushSwitch
		title = config.Config.Notification.GroupMemberCancelMuted.OfflinePush.Title
		desc = config.Config.Notification.GroupMemberCancelMuted.OfflinePush.Desc
		ex = config.Config.Notification.GroupMemberCancelMuted.OfflinePush.Ext
		reliabilityLevel = config.Config.Notification.GroupMemberCancelMuted.Conversation.ReliabilityLevel
		unReadCount = config.Config.Notification.GroupMemberCancelMuted.Conversation.UnreadCount

	case constant.GroupMemberInfoSetNotification:
		pushSwitch = config.Config.Notification.GroupMemberInfoSet.OfflinePush.PushSwitch
		title = config.Config.Notification.GroupMemberInfoSet.OfflinePush.Title
		desc = config.Config.Notification.GroupMemberInfoSet.OfflinePush.Desc
		ex = config.Config.Notification.GroupMemberInfoSet.OfflinePush.Ext
		reliabilityLevel = config.Config.Notification.GroupMemberInfoSet.Conversation.ReliabilityLevel
		unReadCount = config.Config.Notification.GroupMemberInfoSet.Conversation.UnreadCount

	case constant.ConversationPrivateChatNotification:
		pushSwitch = config.Config.Notification.ConversationSetPrivate.OfflinePush.PushSwitch
		title = config.Config.Notification.ConversationSetPrivate.OfflinePush.Title
		desc = config.Config.Notification.ConversationSetPrivate.OfflinePush.Desc
		ex = config.Config.Notification.ConversationSetPrivate.OfflinePush.Ext
		reliabilityLevel = config.Config.Notification.ConversationSetPrivate.Conversation.ReliabilityLevel
		unReadCount = config.Config.Notification.ConversationSetPrivate.Conversation.UnreadCount
	case constant.FriendInfoUpdatedNotification:
		pushSwitch = config.Config.Notification.FriendInfoUpdated.OfflinePush.PushSwitch
		title = config.Config.Notification.FriendInfoUpdated.OfflinePush.Title
		desc = config.Config.Notification.FriendInfoUpdated.OfflinePush.Desc
		ex = config.Config.Notification.FriendInfoUpdated.OfflinePush.Ext
		reliabilityLevel = config.Config.Notification.FriendInfoUpdated.Conversation.ReliabilityLevel
		unReadCount = config.Config.Notification.FriendInfoUpdated.Conversation.UnreadCount
	case constant.DeleteMessageNotification:
		reliabilityLevel = constant.ReliableNotificationNoMsg
	case constant.ConversationUnreadNotification, constant.SuperGroupUpdateNotification:
		reliabilityLevel = constant.UnreliableNotification
	}
	switch reliabilityLevel {
	case constant.UnreliableNotification:
		utils.SetSwitchFromOptions(msg.Options, constant.IsHistory, false)
		utils.SetSwitchFromOptions(msg.Options, constant.IsPersistent, false)
		utils.SetSwitchFromOptions(msg.Options, constant.IsConversationUpdate, false)
		utils.SetSwitchFromOptions(msg.Options, constant.IsSenderConversationUpdate, false)
	case constant.ReliableNotificationNoMsg:
		utils.SetSwitchFromOptions(msg.Options, constant.IsConversationUpdate, false)
		utils.SetSwitchFromOptions(msg.Options, constant.IsSenderConversationUpdate, false)
	case constant.ReliableNotificationMsg:

	}
	utils.SetSwitchFromOptions(msg.Options, constant.IsUnreadCount, unReadCount)
	utils.SetSwitchFromOptions(msg.Options, constant.IsOfflinePush, pushSwitch)
	offlineInfo.Title = title
	offlineInfo.Desc = desc
	offlineInfo.Ex = ex
	msg.OfflinePushInfo = &offlineInfo
	req.MsgData = &msg
	etcdConn := rpc.GetDefaultConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImMsgName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetDefaultConn == nil"
		log.NewError(req.OperationID, errMsg)
		return
	}

	client := pbChat.NewMsgClient(etcdConn)
	reply, err := client.SendMsg(context.Background(), &req)
	if err != nil {
		log.NewError(req.OperationID, "SendMsg rpc failed, ", req.String(), err.Error())
	} else if reply.ErrCode != 0 {
		log.NewError(req.OperationID, "SendMsg rpc failed, ", req.String(), reply.ErrCode, reply.ErrMsg)
	}
}
