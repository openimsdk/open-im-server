package constant

const (

	//group admin
	GroupAdmin = 1
	//feiend related
	BlackListFlag = 1
	NotFriendFlag = 0
	FriendFlag    = 1

	//Websocket Protocol
	WSGetNewestSeq = 1001
	WSPullMsg      = 1002
	WSSendMsg      = 1003
	WSPushMsg      = 2001

	///ContentType
	//UserRelated
	Text    = 101
	Picture = 102
	Voice   = 103
	Video   = 104
	File    = 105

	SyncSenderMsg = 108
	//SysRelated
	AddFriendTip      = 201
	AgreeAddFriendTip = 202
	KickOnlineTip     = 203

	//MsgFrom
	UserMsgType = 100
	SysMsgType  = 200

	//SessionType
	SingleChatType = 1
	GroupChatType  = 2
)

var ContentType2PushContent = map[int64]string{
	Picture: "[picture]",
	Voice:   "[voice]",
	Video:   "[video]",
	File:    "[file]",
}
