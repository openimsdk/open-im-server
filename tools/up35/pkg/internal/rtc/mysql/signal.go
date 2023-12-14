package relation

import (
	"time"
)

type SignalModel struct {
	SID           string    `gorm:"column:sid;type:char(128);primary_key"`
	InviterUserID string    `gorm:"column:inviter_user_id;type:char(64);index:inviter_user_id_index"`
	CustomData    string    `gorm:"column:custom_data;type:text"`
	GroupID       string    `gorm:"column:group_id;type:char(64)"`
	RoomID        string    `gorm:"column:room_id;primary_key;type:char(128)"`
	Timeout       int32     `gorm:"column:timeout"`
	MediaType     string    `gorm:"column:media_type;type:char(64)"`
	PlatformID    int32     `gorm:"column:platform_id"`
	SessionType   int32     `gorm:"column:sesstion_type"`
	InitiateTime  time.Time `gorm:"column:initiate_time"`
	EndTime       time.Time `gorm:"column:end_time"`
	FileURL       string    `gorm:"column:file_url"                                                  json:"-"`

	Title         string `gorm:"column:title;size:128"`
	Desc          string `gorm:"column:desc;size:1024"`
	Ex            string `gorm:"column:ex;size:1024"`
	IOSPushSound  string `gorm:"column:ios_push_sound"`
	IOSBadgeCount bool   `gorm:"column:ios_badge_count"`
	SignalInfo    string `gorm:"column:signal_info;size:1024"`
}

func (SignalModel) TableName() string {
	return "signal"
}

type SignalInvitationModel struct {
	UserID       string    `gorm:"column:user_id;primary_key"`
	SID          string    `gorm:"column:sid;type:char(128);primary_key"`
	Status       int32     `gorm:"column:status"`
	InitiateTime time.Time `gorm:"column:initiate_time;primary_key"`
	HandleTime   time.Time `gorm:"column:handle_time"`
}

func (SignalInvitationModel) TableName() string {
	return "signal_invitation"
}
