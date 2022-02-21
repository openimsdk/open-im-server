package constant

const (

	//group admin
	//	OrdinaryMember = 0
	//	GroupOwner     = 1
	//	Administrator  = 2
	//group application
	//	Application      = 0
	//	AgreeApplication = 1

	//friend related
	BlackListFlag         = 1
	ApplicationFriendFlag = 0
	FriendFlag            = 1
	RefuseFriendFlag      = -1

	//Websocket Protocol
	WSGetNewestSeq     = 1001
	WSPullMsgBySeqList = 1002
	WSSendMsg          = 1003
	WSPushMsg          = 2001
	WSKickOnlineMsg    = 2002
	WSDataError        = 3001

	///ContentType
	//UserRelated
	Text           = 101
	Picture        = 102
	Voice          = 103
	Video          = 104
	File           = 105
	AtText         = 106
	Merger         = 107
	Card           = 108
	Location       = 109
	Custom         = 110
	Revoke         = 111
	HasReadReceipt = 112
	Typing         = 113
	Quote          = 114
	Common         = 200
	GroupMsg       = 201

	//SysRelated
	NotificationBegin = 1000

	FriendApplicationApprovedNotification = 1201 //add_friend_response
	FriendApplicationRejectedNotification = 1202 //add_friend_response
	FriendApplicationNotification         = 1203 //add_friend
	FriendAddedNotification               = 1204
	FriendDeletedNotification             = 1205 //delete_friend
	FriendRemarkSetNotification           = 1206 //set_friend_remark?
	BlackAddedNotification                = 1207 //add_black
	BlackDeletedNotification              = 1208 //remove_black

	UserInfoUpdatedNotification = 1303 //SetSelfInfoTip             = 204

	ConversationOptChangeNotification    = 1300 // change conversation opt

	GroupCreatedNotification             = 1501
	GroupInfoSetNotification             = 1502
	JoinGroupApplicationNotification     = 1503
	MemberQuitNotification               = 1504
	GroupApplicationAcceptedNotification = 1505
	GroupApplicationRejectedNotification = 1506
	GroupOwnerTransferredNotification    = 1507
	MemberKickedNotification             = 1508
	MemberInvitedNotification            = 1509
	MemberEnterNotification              = 1510
	NotificationEnd                      = 2000

	//MsgFrom
	UserMsgType = 100
	SysMsgType  = 200

	//SessionType
	SingleChatType = 1
	GroupChatType  = 2
	//token
	NormalToken  = 0
	InValidToken = 1
	KickedToken  = 2
	ExpiredToken = 3

	//MultiTerminalLogin
	//Full-end login, but the same end is mutually exclusive
	AllLoginButSameTermKick = 1
	//Only one of the endpoints can log in
	SingleTerminalLogin = 2
	//The web side can be online at the same time, and the other side can only log in at one end
	WebAndOther = 3
	//The PC side is mutually exclusive, and the mobile side is mutually exclusive, but the web side can be online at the same time
	PcMobileAndWeb = 4

	OnlineStatus  = "online"
	OfflineStatus = "offline"
	Registered    = "registered"
	UnRegistered  = "unregistered"

	//MsgReceiveOpt
	ReceiveMessage          = 0
	NotReceiveMessage       = 1
	ReceiveNotNotifyMessage = 2

	//OptionsKey
	IsHistory            = "history"
	IsPersistent         = "persistent"
	IsOfflinePush        = "offlinePush"
	IsUnreadCount        = "unreadCount"
	IsConversationUpdate = "conversationUpdate"
	IsSenderSync         = "senderSync"

	//GroupStatus
	GroupOk             = 0
	GroupBanChat        = 1
	GroupDisband        = 2
	GroupBaned          = 3
	GroupBanPrivateChat = 4

	//UserJoinGroupSource
	JoinByAdmin = 1

	//Minio
	MinioDurationTimes = 3600

	// verificationCode used for
	VerificationCodeForRegister = 1
	VerificationCodeForReset = 2
	VerificationCodeForRegisterSuffix = "_forRegister"
	VerificationCodeForResetSuffix = "_forReset"
)

var ContentType2PushContent = map[int64]string{
	Picture:  "[图片]",
	Voice:    "[语音]",
	Video:    "[视频]",
	File:     "[文件]",
	Text:     "你收到了一条文本消息",
	AtText:   "[有人@你]",
	GroupMsg: "你收到一条群聊消息",
	Common:   "你收到一条新消息",
}

const (
	AppOrdinaryUsers = 1
	AppAdmin         = 2

	GroupOrdinaryUsers = 1
	GroupOwner         = 2
	GroupAdmin         = 3

	GroupResponseAgree  = 1
	GroupResponseRefuse = -1

	FriendResponseAgree  = 1
	FriendResponseRefuse = -1

	Male   = 1
	Female = 2
)

const (
	UnreliableNotification    = 1
	ReliableNotificationNoMsg = 2
	ReliableNotificationMsg   = 3
)

const FriendAcceptTip = "You have successfully become friends, so start chatting"

func GroupIsBanChat(status int32) bool {
	if status != GroupBanChat {
		return false
	}
	return true
}

func GroupIsBanPrivateChat(status int32) bool {
	if status != GroupBanPrivateChat {
		return false
	}
	return true
}

 const BigVersion = "v3"