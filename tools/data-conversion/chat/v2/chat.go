package v2

import (
	"time"
)

// Register 注册信息表
type Register struct {
	UserID      string    `gorm:"column:user_id;primary_key;type:char(64)" json:"userID"`
	DeviceID    string    `gorm:"column:device_id;type:varchar(255)" json:"deviceID"`
	IP          string    `gorm:"column:ip;type:varchar(32)" json:"ip"`
	Platform    string    `gorm:"column:platform;type:varchar(32)" json:"platform"`
	AccountType string    `gorm:"column:account_type;type:varchar(32)" json:"accountType"` //email phone account
	Mode        string    `gorm:"column:mode;type:varchar(32)"`                            //user admin
	CreateTime  time.Time `gorm:"column:create_time" json:"createTime"`
}

// Account 账号密码表
type Account struct {
	UserID         string    `gorm:"column:user_id;primary_key;type:char(64)" json:"userID"`
	Password       string    `gorm:"column:password;type:varchar(255)" json:"password"`
	CreateTime     time.Time `gorm:"column:create_time" json:"createTime"`
	ChangeTime     time.Time `gorm:"column:change_time" json:"changeTime"`
	OperatorUserID string    `gorm:"column:operator_user_id;type:varchar(64)" json:"operatorUserID"`
}

// Attribute 用户属性表
type Attribute struct {
	UserID         string    `gorm:"column:user_id;primary_key;type:char(64)" json:"userID"`
	Account        string    `gorm:"column:account;type:char(64)" json:"account"`
	PhoneNumber    string    `gorm:"column:phone_number;type:varchar(32)" json:"phoneNumber"`
	AreaCode       string    `gorm:"column:area_code;type:varchar(8)" json:"areaCode"`
	Email          string    `gorm:"column:email;type:varchar(64)"  json:"email"`
	Nickname       string    `gorm:"column:nickname;type:varchar(64)"  json:"nickname"`
	FaceURL        string    `gorm:"column:face_url;type:varchar(255)"  json:"faceURL"`
	Gender         int32     `gorm:"column:gender" json:"gender"`
	Birth          uint32    `gorm:"column:birth" json:"birth"`
	CreateTime     time.Time `gorm:"column:create_time" json:"createTime"`
	ChangeTime     time.Time `gorm:"column:change_time" json:"changeTime"`
	BirthTime      time.Time `gorm:"column:birth_time" json:"birthTime"`
	Level          int32     `gorm:"column:level;default:1"   json:"level"`
	AllowVibration int32     `gorm:"column:allow_vibration;default:1" json:"allowVibration"`
	AllowBeep      int32     `gorm:"column:allow_beep;default:1" json:"allowBeep"`
	AllowAddFriend int32     `gorm:"column:allow_add_friend;default:1" json:"allowAddFriend"`
}

// 封号表
type ForbiddenAccount struct {
	UserID         string    `gorm:"column:user_id;index:userID;primary_key;type:char(64)" json:"userID"`
	CreateTime     time.Time `gorm:"column:create_time"  json:"createTime"`
	Reason         string    `gorm:"column:reason;type:varchar(255)"  json:"reason"`
	OperatorUserID string    `gorm:"column:operator_user_id;type:varchar(255)" json:"operatorUserID"`
}

// 用户登录信息表
type UserLoginRecord struct {
	UserID    string    `gorm:"column:user_id;size:64" json:"userID"`
	LoginTime time.Time `gorm:"column:login_time" json:"loginTime"`
	IP        string    `gorm:"column:ip;type:varchar(32)" json:"ip"`
	DeviceID  string    `gorm:"column:device_id;type:varchar(255)" json:"deviceID"`
	Platform  string    `gorm:"column:platform;type:varchar(32)" json:"platform"`
}

// 禁止ip登录 注册
type IPForbidden struct {
	IP            string    `gorm:"column:ip;primary_key;type:char(32)" json:"ip"`
	LimitRegister int32     `gorm:"column:limit_register" json:"limitRegister"`
	LimitLogin    int32     `gorm:"column:limit_login" json:"limitLogin"`
	CreateTime    time.Time `gorm:"column:create_time" json:"createTime"`
}

// 限制userID只能在某些ip登录
type LimitUserLoginIP struct {
	UserID     string    `gorm:"column:user_id;primary_key;type:char(64)" json:"userID"`
	IP         string    `gorm:"column:ip;primary_key;type:char(32)" json:"ip"`
	CreateTime time.Time `gorm:"column:create_time"  json:"createTime"`
}

// 邀请码被注册使用
type InvitationRegister struct {
	InvitationCode string    `gorm:"column:invitation_code;primary_key;type:char(32)" json:"invitationCode"`
	CreateTime     time.Time `gorm:"column:create_time" json:"createTime"`
	UsedByUserID   string    `gorm:"column:user_id;index:userID;type:char(64)" json:"usedByUserID"`
}

type SignalRecord struct {
	FileName    string    `gorm:"column:file_name;primary_key;type:char(128)" json:"fileName"`
	MediaType   string    `gorm:"column:media_type;type:char(64);index:media_type_index" json:"mediaType"`
	RoomType    string    `gorm:"column:room_type;type:char(20)" json:"roomType"`
	SenderID    string    `gorm:"column:sender_id;type:char(64);index:sender_id_index" json:"senderID"`
	RecvID      string    `gorm:"column:recv_id;type:char(64);index:recv_id_index" json:"recvID"`
	GroupID     string    `gorm:"column:group_id;type:char(64)" json:"groupID"`
	DownloadURL string    `gorm:"column:download_url;type:text" json:"downloadURL"`
	CreateTime  time.Time `gorm:"create_time;index:create_time_index" json:"createTime"`
}
