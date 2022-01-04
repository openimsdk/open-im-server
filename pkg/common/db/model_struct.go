package db

import "time"

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
	OwnerUserID    string    `gorm:"column:owner_user_id;primary_key;"`
	Remark         string    `gorm:"column:remark"`
	CreateTime     time.Time `gorm:"column:create_time"`
	FriendUserID   string    `gorm:"column:friend_user_id;primary_key;"`
	AddSource      int32     `gorm:"column:add_source"`
	OperatorUserID string    `gorm:"column:operator_user_id"`
	Ex             string    `gorm:"column:ex"`
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
//open_im_sdk.FriendRequest == imdb.FriendRequest
type FriendRequest struct {
	FromUserID    string    `gorm:"column:from_user_id;primaryKey;"`
	ToUserID      string    `gorm:"column:to_user_id;primaryKey;"`
	HandleResult  int32     `gorm:"column:handle_result"`
	ReqMsg        string    `gorm:"column:req_msg"`
	CreateTime    time.Time `gorm:"column:create_time"`
	HandlerUserID string    `gorm:"column:handler_user_id"`
	HandleMsg     string    `gorm:"column:handle_msg"`
	HandleTime    time.Time `gorm:"column:handle_time"`
	Ex            string    `gorm:"column:ex"`
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
	GroupID      string `gorm:"column:group_id;primaryKey;" json:"groupID" binding:"required"`
	GroupName    string `gorm:"column:name" json:"groupName"`
	Notification string `gorm:"column:notification" json:"notification"`
	Introduction string `gorm:"column:introduction" json:"introduction"`
	FaceUrl      string `gorm:"column:face_url" json:"faceUrl"`

	CreateTime    time.Time `gorm:"column:create_time"`
	Status        int32     `gorm:"column:status"`
	CreatorUserID string    `gorm:"column:creator_user_id"`
	GroupType     int32     `gorm:"column:group_type"`
	Ex            string    `gorm:"column:ex" json:"ex"`
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
	GroupID        string    `gorm:"column:group_id;primaryKey;"`
	UserID         string    `gorm:"column:user_id;primaryKey;"`
	RoleLevel      int32     `gorm:"column:role_level"`
	JoinTime       time.Time `gorm:"column:join_time"`
	Nickname       string    `gorm:"column:nickname"`
	FaceUrl        string    `gorm:"column:user_group_face_url"`
	JoinSource     int32     `gorm:"column:join_source"`
	OperatorUserID string    `gorm:"column:operator_user_id"`
	Ex             string    `gorm:"column:ex"`
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
	UserID       string    `gorm:"column:user_id;primaryKey;"`
	GroupID      string    `gorm:"column:group_id;primaryKey;"`
	HandleResult int32     `gorm:"column:handle_result"`
	ReqMsg       string    `gorm:"column:req_msg"`
	HandledMsg   string    `gorm:"column:handle_msg"`
	ReqTime      time.Time `gorm:"column:req_time"`
	HandleUserID string    `gorm:"column:handle_user_id"`
	HandledTime  time.Time `gorm:"column:handle_time"`
	Ex           string    `gorm:"column:ex"`
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
	UserID         string    `gorm:"column:user_id;primaryKey;"`
	Nickname       string    `gorm:"column:name"`
	FaceUrl        string    `gorm:"column:face_url"`
	Gender         int32     `gorm:"column:gender"`
	PhoneNumber    string    `gorm:"column:phone_number"`
	Birth          time.Time `gorm:"column:birth"`
	Email          string    `gorm:"column:email"`
	Ex             string    `gorm:"column:ex"`
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
	CreateTime     time.Time `gorm:"column:create_time"`
	BlockUserID    string    `gorm:"column:block_user_id;primary_key;size:64"`
	AddSource      int32     `gorm:"column:add_source"`
	OperatorUserID string    `gorm:"column:operator_user_id;size:64"`
	Ex             string    `gorm:"column:ex;size:1024"`
}
