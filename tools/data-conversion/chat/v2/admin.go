package v2

import (
	"time"
)

// AppVersion pc端版本管理
type AppVersion struct {
	Version     string `gorm:"column:version;size:64" json:"version"`
	Type        int    `gorm:"column:type;primary_key" json:"type"`
	UpdateTime  int    `gorm:"column:update_time" json:"update_time"`
	ForceUpdate bool   `gorm:"column:force_update" json:"force_update"`
	FileName    string `gorm:"column:file_name" json:"file_name"`
	YamlName    string `gorm:"column:yaml_name" json:"yaml_name"`
	UpdateLog   string `gorm:"column:update_log" json:"update_log"`
}

// Admin 后台管理员
type Admin struct {
	Account    string    `gorm:"column:account;primary_key;type:char(64)" json:"account"`
	Password   string    `gorm:"column:Password;type:char(64)" json:"password"`
	FaceURL    string    `gorm:"column:FaceURL;type:char(64)" json:"faceURL"`
	Nickname   string    `gorm:"column:Nickname;type:char(64)" json:"nickname"`
	UserID     string    `gorm:"column:UserID;type:char(64)" json:"userID"` //openIM userID
	Level      int32     `gorm:"column:level;default:1"   json:"level"`
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"`
}

// RegisterAddFriend 注册时默认好友
type RegisterAddFriend struct {
	UserID     string    `gorm:"column:user_id;primary_key;type:char(64)" json:"userID"`
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"`
}

// RegisterAddGroup 注册时默认群组
type RegisterAddGroup struct {
	GroupID    string    `gorm:"column:group_id;primary_key;type:char(64)" json:"userID"`
	CreateTime time.Time `gorm:"column:create_time" json:"createTime"`
}

// ClientInitConfig 系统相关配置项
type ClientInitConfig struct {
	DiscoverPageURL            string `gorm:"column:discover_page_url;size:128" json:"discoverPageURL"`
	OrdinaryUserAddFriend      int32  `gorm:"column:ordinary_user_add_friend; default:1"   json:"ordinaryUserAddFriend"`
	BossUserID                 string `gorm:"column:boss_user_id;type:char(64)" json:"bossUserID"`
	AdminURL                   string `gorm:"column:admin_url;type:char(128)" json:"adminURL"`
	AllowSendMsgNotFriend      int32  `gorm:"column:allow_send_msg_not_friend;default:1" json:"allowSendMsgNotFriend"`
	NeedInvitationCodeRegister int32  `gorm:"column:need_invitation_code_register;default:0" json:"needInvitationCodeRegister"`
}

type Applet struct {
	ID         string    `gorm:"column:id;primary_key;size:64"`
	Name       string    `gorm:"column:name;uniqueIndex;size:64"`
	AppID      string    `gorm:"column:app_id;uniqueIndex;size:255"`
	Icon       string    `gorm:"column:icon;size:255"`
	URL        string    `gorm:"column:url;size:255"`
	MD5        string    `gorm:"column:md5;size:255"`
	Size       int64     `gorm:"column:size"`
	Version    string    `gorm:"column:version;size:64"`
	Priority   uint32    `gorm:"column:priority;size:64"`
	Status     uint8     `gorm:"column:status"`
	CreateTime time.Time `gorm:"column:create_time;autoCreateTime;size:64"`
}
