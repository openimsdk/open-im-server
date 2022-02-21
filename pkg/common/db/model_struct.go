package db

import "time"

type Register struct {
	Account  string `gorm:"column:account;primary_key;type:char(255)" json:"account"`
	Password string `gorm:"column:password;type:varchar(255)" json:"password"`
	Ex       string `gorm:"column:ex;size:1024" json:"ex"`
}

//
//message FriendInfo{
//string OwnerUserID = 1;
//string Remark = 2;
//int64 CreateTime = 3;
//UserInfo FriendUser = 4;
//int32 AddSource = 5;
//string OperatorUserID = 6;
//string Ex = 7;
//}
//open_im_sdk.FriendInfo(FriendUser) != imdb.Friend(FriendUserID)
type Friend struct {
	OwnerUserID    string    `gorm:"column:owner_user_id;primary_key;size:64"`
	FriendUserID   string    `gorm:"column:friend_user_id;primary_key;size:64"`
	Remark         string    `gorm:"column:remark;size:255"`
	CreateTime     time.Time `gorm:"column:create_time"`
	AddSource      int32     `gorm:"column:add_source"`
	OperatorUserID string    `gorm:"column:operator_user_id;size:64"`
	Ex             string    `gorm:"column:ex;size:1024"`
}

//message FriendRequest{
//string  FromUserID = 1;
//string ToUserID = 2;
//int32 HandleResult = 3;
//string ReqMsg = 4;
//int64 CreateTime = 5;
//string HandlerUserID = 6;
//string HandleMsg = 7;
//int64 HandleTime = 8;
//string Ex = 9;
//}
//open_im_sdk.FriendRequest(nickname, farce url ...) != imdb.FriendRequest
type FriendRequest struct {
	FromUserID    string    `gorm:"column:from_user_id;primary_key;size:64"`
	ToUserID      string    `gorm:"column:to_user_id;primary_key;size:64"`
	HandleResult  int32     `gorm:"column:handle_result"`
	ReqMsg        string    `gorm:"column:req_msg;size:255"`
	CreateTime    time.Time `gorm:"column:create_time"`
	HandlerUserID string    `gorm:"column:handler_user_id;size:64"`
	HandleMsg     string    `gorm:"column:handle_msg;size:255"`
	HandleTime    time.Time `gorm:"column:handle_time"`
	Ex            string    `gorm:"column:ex;size:1024"`
}

func (FriendRequest) TableName() string {
	return "friend_requests"
}

//message GroupInfo{
//  string GroupID = 1;
//  string GroupName = 2;
//  string Notification = 3;
//  string Introduction = 4;
//  string FaceUrl = 5;
//  string OwnerUserID = 6;
//  uint32 MemberCount = 8;
//  int64 CreateTime = 7;
//  string Ex = 9;
//  int32 Status = 10;
//  string CreatorUserID = 11;
//  int32 GroupType = 12;
//}
//  open_im_sdk.GroupInfo (OwnerUserID ,  MemberCount )> imdb.Group
type Group struct {
	//`json:"operationID" binding:"required"`
	//`protobuf:"bytes,1,opt,name=GroupID" json:"GroupID,omitempty"` `json:"operationID" binding:"required"`
	GroupID       string    `gorm:"column:group_id;primary_key;size:64" json:"groupID" binding:"required"`
	GroupName     string    `gorm:"column:name;size:255" json:"groupName"`
	Notification  string    `gorm:"column:notification;size:255" json:"notification"`
	Introduction  string    `gorm:"column:introduction;size:255" json:"introduction"`
	FaceURL       string    `gorm:"column:face_url;size:255" json:"faceURL"`
	CreateTime    time.Time `gorm:"column:create_time"`
	Ex            string    `gorm:"column:ex" json:"ex;size:1024" json:"ex"`
	Status        int32     `gorm:"column:status"`
	CreatorUserID string    `gorm:"column:creator_user_id;size:64"`
	GroupType     int32     `gorm:"column:group_type"`
}

//message GroupMemberFullInfo {
//string GroupID = 1 ;
//string UserID = 2 ;
//int32 roleLevel = 3;
//int64 JoinTime = 4;
//string NickName = 5;
//string FaceUrl = 6;
//int32 JoinSource = 8;
//string OperatorUserID = 9;
//string Ex = 10;
//int32 AppMangerLevel = 7; //if >0
//}  open_im_sdk.GroupMemberFullInfo(AppMangerLevel) > imdb.GroupMember
type GroupMember struct {
	GroupID        string    `gorm:"column:group_id;primary_key;size:64"`
	UserID         string    `gorm:"column:user_id;primary_key;size:64"`
	Nickname       string    `gorm:"column:nickname;size:255"`
	FaceURL        string    `gorm:"column:user_group_face_url;size:255"`
	RoleLevel      int32     `gorm:"column:role_level"`
	JoinTime       time.Time `gorm:"column:join_time"`
	JoinSource     int32     `gorm:"column:join_source"`
	OperatorUserID string    `gorm:"column:operator_user_id;size:64"`
	Ex             string    `gorm:"column:ex;size:1024"`
}

//message GroupRequest{
//string UserID = 1;
//string GroupID = 2;
//string HandleResult = 3;
//string ReqMsg = 4;
//string  HandleMsg = 5;
//int64 ReqTime = 6;
//string HandleUserID = 7;
//int64 HandleTime = 8;
//string Ex = 9;
//}open_im_sdk.GroupRequest == imdb.GroupRequest
type GroupRequest struct {
	UserID       string    `gorm:"column:user_id;primary_key;size:64"`
	GroupID      string    `gorm:"column:group_id;primary_key;size:64"`
	HandleResult int32     `gorm:"column:handle_result"`
	ReqMsg       string    `gorm:"column:req_msg;size:1024"`
	HandledMsg   string    `gorm:"column:handle_msg;size:1024"`
	ReqTime      time.Time `gorm:"column:req_time"`
	HandleUserID string    `gorm:"column:handle_user_id;size:64"`
	HandledTime  time.Time `gorm:"column:handle_time"`
	Ex           string    `gorm:"column:ex;size:1024"`
}

//string UserID = 1;
//string Nickname = 2;
//string FaceUrl = 3;
//int32 Gender = 4;
//string PhoneNumber = 5;
//string Birth = 6;
//string Email = 7;
//string Ex = 8;
//int64 CreateTime = 9;
//int32 AppMangerLevel = 10;
//open_im_sdk.User == imdb.User
type User struct {
	UserID         string    `gorm:"column:user_id;primary_key;size:64"`
	Nickname       string    `gorm:"column:name;size:255"`
	FaceURL        string    `gorm:"column:face_url;size:255"`
	Gender         int32     `gorm:"column:gender"`
	PhoneNumber    string    `gorm:"column:phone_number;size:32"`
	Birth          time.Time `gorm:"column:birth"`
	Email          string    `gorm:"column:email;size:64"`
	Ex             string    `gorm:"column:ex;size:1024"`
	CreateTime     time.Time `gorm:"column:create_time"`
	AppMangerLevel int32     `gorm:"column:app_manger_level"`
}

//message BlackInfo{
//string OwnerUserID = 1;
//int64 CreateTime = 2;
//PublicUserInfo BlackUserInfo = 4;
//int32 AddSource = 5;
//string OperatorUserID = 6;
//string Ex = 7;
//}
// open_im_sdk.BlackInfo(BlackUserInfo) != imdb.Black (BlockUserID)
type Black struct {
	OwnerUserID    string    `gorm:"column:owner_user_id;primary_key;size:64"`
	BlockUserID    string    `gorm:"column:block_user_id;primary_key;size:64"`
	CreateTime     time.Time `gorm:"column:create_time"`
	AddSource      int32     `gorm:"column:add_source"`
	OperatorUserID string    `gorm:"column:operator_user_id;size:64"`
	Ex             string    `gorm:"column:ex;size:1024"`
}

type ChatLog struct {
	ServerMsgID      string    `gorm:"column:server_msg_id;primary_key;type:char(64)" json:"serverMsgID"`
	ClientMsgID      string    `gorm:"column:client_msg_id;type:char(64)" json:"clientMsgID"`
	SendID           string    `gorm:"column:send_id;type:char(64)" json:"sendID"`
	RecvID           string    `gorm:"column:recv_id;type:char(64)" json:"recvID"`
	SenderPlatformID int32     `gorm:"column:sender_platform_id" json:"senderPlatformID"`
	SenderNickname   string    `gorm:"column:sender_nick_name;type:varchar(255)" json:"senderNickname"`
	SenderFaceURL    string    `gorm:"column:sender_face_url;type:varchar(255)" json:"senderFaceURL"`
	SessionType      int32     `gorm:"column:session_type" json:"sessionType"`
	MsgFrom          int32     `gorm:"column:msg_from" json:"msgFrom"`
	ContentType      int32     `gorm:"column:content_type" json:"contentType"`
	Content          string    `gorm:"column:content;type:varchar(3000)" json:"content"`
	Status           int32     `gorm:"column:status" json:"status"`
	SendTime         time.Time `gorm:"column:send_time" json:"sendTime"`
	CreateTime       time.Time `gorm:"column:create_time" json:"createTime"`
	Ex               string    `gorm:"column:ex;type:varchar(1024)" json:"ex"`
}

func (ChatLog) TableName() string {
	return "chat_logs"
}

type BlackList struct {
	UserId           string    `gorm:"column:uid"`
	BeginDisableTime time.Time `gorm:"column:begin_disable_time"`
	EndDisableTime   time.Time `gorm:"column:end_disable_time"`
}
type Conversation struct {
	OwnerUserID      string `gorm:"column:owner_user_id;primary_key;type:char(128)" json:"OwnerUserID"`
	ConversationID   string `gorm:"column:conversation_id;primary_key;type:char(128)" json:"conversationID"`
	ConversationType int32  `gorm:"column:conversation_type" json:"conversationType"`
	UserID           string `gorm:"column:user_id;type:char(64)" json:"userID"`
	GroupID          string `gorm:"column:group_id;type:char(128)" json:"groupID"`
	RecvMsgOpt       int32  `gorm:"column:recv_msg_opt" json:"recvMsgOpt"`
	UnreadCount      int32  `gorm:"column:unread_count" json:"unreadCount"`
	DraftTextTime    int64  `gorm:"column:draft_text_time" json:"draftTextTime"`
	IsPinned         bool   `gorm:"column:is_pinned" json:"isPinned"`
	AttachedInfo     string `gorm:"column:attached_info;type:varchar(1024)" json:"attachedInfo"`
	Ex               string `gorm:"column:ex;type:varchar(1024)" json:"ex"`
}
