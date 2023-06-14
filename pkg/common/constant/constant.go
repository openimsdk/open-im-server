package constant

const (

	///ContentType
	//UserRelated
	Text     = 101
	Picture  = 102
	Voice    = 103
	Video    = 104
	File     = 105
	AtText   = 106
	Merger   = 107
	Card     = 108
	Location = 109
	Custom   = 110
	Revoke   = 111
	Typing   = 113
	Quote    = 114

	AdvancedText = 117

	CustomNotTriggerConversation = 119
	CustomOnlineOnly             = 120
	ReactionMessageModifier      = 121
	ReactionMessageDeleter       = 122

	Common             = 200
	GroupMsg           = 201
	SignalMsg          = 202
	CustomNotification = 203

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
	FriendInfoUpdatedNotification         = 1209

	ConversationChangeNotification = 1300 // change conversation opt

	UserNotificationBegin       = 1301
	UserInfoUpdatedNotification = 1303 //SetSelfInfoTip              = 204
	UserNotificationEnd         = 1399
	OANotification              = 1400

	GroupNotificationBegin = 1500

	GroupCreatedNotification                 = 1501
	GroupInfoSetNotification                 = 1502
	JoinGroupApplicationNotification         = 1503
	MemberQuitNotification                   = 1504
	GroupApplicationAcceptedNotification     = 1505
	GroupApplicationRejectedNotification     = 1506
	GroupOwnerTransferredNotification        = 1507
	MemberKickedNotification                 = 1508
	MemberInvitedNotification                = 1509
	MemberEnterNotification                  = 1510
	GroupDismissedNotification               = 1511
	GroupMemberMutedNotification             = 1512
	GroupMemberCancelMutedNotification       = 1513
	GroupMutedNotification                   = 1514
	GroupCancelMutedNotification             = 1515
	GroupMemberInfoSetNotification           = 1516
	GroupMemberSetToAdminNotification        = 1517
	GroupMemberSetToOrdinaryUserNotification = 1518
	GroupInfoSetAnnouncementNotification     = 1519
	GroupInfoSetNameNotification             = 1520

	SignalingNotificationBegin = 1600
	SignalingNotification      = 1601
	SignalingNotificationEnd   = 1649

	SuperGroupNotificationBegin  = 1650
	SuperGroupUpdateNotification = 1651
	MsgDeleteNotification        = 1652
	SuperGroupNotificationEnd    = 1699

	ConversationPrivateChatNotification = 1701
	ConversationUnreadNotification      = 1702

	MsgRevokeNotification = 2101

	BusinessNotificationBegin = 2000
	BusinessNotification      = 2001
	BusinessNotificationEnd   = 2099

	ClearConversationNotification = 2101
	DeleteMsgsNotification        = 2102

	HasReadReceipt = 2200

	NotificationEnd = 5000

	//status
	MsgNormal  = 1
	MsgDeleted = 4

	//MsgFrom
	UserMsgType = 100
	SysMsgType  = 200

	//SessionType
	SingleChatType       = 1
	GroupChatType        = 2
	SuperGroupChatType   = 3
	NotificationChatType = 4
	//token
	NormalToken  = 0
	InValidToken = 1
	KickedToken  = 2
	ExpiredToken = 3

	//MultiTerminalLogin
	DefalutNotKick = 0
	//Full-end login, but the same end is mutually exclusive
	AllLoginButSameTermKick = 1
	//Only one of the endpoints can log in
	SingleTerminalLogin = 2
	//The web side can be online at the same time, and the other side can only log in at one end
	WebAndOther = 3
	//The PC side is mutually exclusive, and the mobile side is mutually exclusive, but the web side can be online at the same time
	PcMobileAndWeb = 4
	//The PC terminal can be online at the same time,but other terminal only one of the endpoints can login
	PCAndOther = 5

	OnlineStatus  = "online"
	OfflineStatus = "offline"
	Registered    = "registered"
	UnRegistered  = "unregistered"

	//MsgReceiveOpt
	ReceiveMessage          = 0
	NotReceiveMessage       = 1
	ReceiveNotNotifyMessage = 2

	//OptionsKey
	IsHistory                  = "history"
	IsPersistent               = "persistent"
	IsOfflinePush              = "offlinePush"
	IsUnreadCount              = "unreadCount"
	IsConversationUpdate       = "conversationUpdate"
	IsSenderSync               = "senderSync"
	IsNotPrivate               = "notPrivate"
	IsSenderConversationUpdate = "senderConversationUpdate"
	IsSenderNotificationPush   = "senderNotificationPush"
	IsReactionFromCache        = "reactionFromCache"
	IsNotNotification          = "isNotNotification"
	IsSendMsg                  = "isSendMsg"

	//GroupStatus
	GroupOk              = 0
	GroupBanChat         = 1
	GroupStatusDismissed = 2
	GroupStatusMuted     = 3

	//GroupType
	NormalGroup  = 0
	SuperGroup   = 1
	WorkingGroup = 2

	GroupBaned          = 3
	GroupBanPrivateChat = 4

	//UserJoinGroupSource
	JoinByAdmin = 1

	JoinByInvitation = 2
	JoinBySearch     = 3
	JoinByQRCode     = 4

	//Minio
	MinioDurationTimes = 3600
	//Aws
	AwsDurationTimes = 3600

	//callbackCommand
	CallbackBeforeSendSingleMsgCommand                   = "callbackBeforeSendSingleMsgCommand"
	CallbackAfterSendSingleMsgCommand                    = "callbackAfterSendSingleMsgCommand"
	CallbackBeforeSendGroupMsgCommand                    = "callbackBeforeSendGroupMsgCommand"
	CallbackAfterSendGroupMsgCommand                     = "callbackAfterSendGroupMsgCommand"
	CallbackMsgModifyCommand                             = "callbackMsgModifyCommand"
	CallbackUserOnlineCommand                            = "callbackUserOnlineCommand"
	CallbackUserOfflineCommand                           = "callbackUserOfflineCommand"
	CallbackUserKickOffCommand                           = "callbackUserKickOffCommand"
	CallbackOfflinePushCommand                           = "callbackOfflinePushCommand"
	CallbackOnlinePushCommand                            = "callbackOnlinePushCommand"
	CallbackSuperGroupOnlinePushCommand                  = "callbackSuperGroupOnlinePushCommand"
	CallbackBeforeAddFriendCommand                       = "callbackBeforeAddFriendCommand"
	CallbackBeforeCreateGroupCommand                     = "callbackBeforeCreateGroupCommand"
	CallbackBeforeMemberJoinGroupCommand                 = "callbackBeforeMemberJoinGroupCommand"
	CallbackBeforeSetGroupMemberInfoCommand              = "CallbackBeforeSetGroupMemberInfoCommand"
	CallbackBeforeSetMessageReactionExtensionCommand     = "callbackBeforeSetMessageReactionExtensionCommand"
	CallbackBeforeDeleteMessageReactionExtensionsCommand = "callbackBeforeDeleteMessageReactionExtensionsCommand"
	CallbackGetMessageListReactionExtensionsCommand      = "callbackGetMessageListReactionExtensionsCommand"
	CallbackAddMessageListReactionExtensionsCommand      = "callbackAddMessageListReactionExtensionsCommand"

	//callback actionCode
	ActionAllow     = 0
	ActionForbidden = 1
	//callback callbackHandleCode
	CallbackHandleSuccess = 0
	CallbackHandleFailed  = 1

	// minioUpload
	OtherType = 1
	VideoType = 2
	ImageType = 3

	// sendMsgStaus
	MsgStatusNotExist = 0
	MsgIsSending      = 1
	MsgSendSuccessed  = 2
	MsgSendFailed     = 3
)

const (
	WriteDiffusion = 0
	ReadDiffusion  = 1
)

const (
	UnreliableNotification    = 1
	ReliableNotificationNoMsg = 2
	ReliableNotificationMsg   = 3
)

const (
	AtAllString       = "AtAllTag"
	AtNormal          = 0
	AtMe              = 1
	AtAll             = 2
	AtAllAtMe         = 3
	GroupNotification = 4
)

var ContentType2PushContent = map[int64]string{
	Picture:   "[PICTURE]",
	Voice:     "[VOICE]",
	Video:     "[VIDEO]",
	File:      "[File]",
	Text:      "[TEXT]",
	AtText:    "[@TEXT]",
	GroupMsg:  "[GROUPMSG]]",
	Common:    "[NEWMSG]",
	SignalMsg: "[SIGNALINVITE]",
}

const (
	FieldRecvMsgOpt    = 1
	FieldIsPinned      = 2
	FieldAttachedInfo  = 3
	FieldIsPrivateChat = 4
	FieldGroupAtType   = 5
	FieldEx            = 7
	FieldUnread        = 8
	FieldBurnDuration  = 9
	FieldHasReadSeq    = 10
)

const (
	AppOrdinaryUsers = 1
	AppAdmin         = 2

	GroupOwner         = 100
	GroupAdmin         = 60
	GroupOrdinaryUsers = 20

	GroupResponseAgree  = 1
	GroupResponseRefuse = -1

	FriendResponseAgree  = 1
	FriendResponseRefuse = -1

	Male   = 1
	Female = 2
)

const OperationID = "operationID"
const OpUserID = "opUserID"
const ConnID = "connID"
const OpUserPlatform = "platform"
const Token = "token"
const RpcCustomHeader = "customHeader" // rpc中间件自定义ctx参数
const CheckKey = "CheckKey"
const TriggerID = "triggerID"
const RemoteAddr = "remoteAddr"

const (
	BecomeFriendByImport = 1 //管理员导入
	BecomeFriendByApply  = 2 //申请添加
)

const (
	ApplyNeedVerificationInviteDirectly = 0 // 申请需要同意 邀请直接进
	AllNeedVerification                 = 1 //所有人进群需要验证，除了群主管理员邀请进群
	Directly                            = 2 //直接进群
)

const (
	GroupRPCRecvSize = 30
	GroupRPCSendSize = 30
)

const FriendAcceptTip = "You have successfully become friends, so start chatting"

func GroupIsBanChat(status int32) bool {
	if status != GroupStatusMuted {
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

const LogFileName = "OpenIM.log"

const LocalHost = "0.0.0.0"

// flag parse
const (
	FlagPort   = "port"
	FlagWsPort = "ws_port"

	FlagPrometheusPort = "prometheus_port"
	FlagConf           = "config_folder_path"
)

const OpenIMCommonConfigKey = "OpenIMServerConfig"

const CallbackCommand = "command"
