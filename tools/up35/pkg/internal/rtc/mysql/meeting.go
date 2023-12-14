package relation

import (
	"time"
)

type MeetingInfo struct {
	RoomID      string    `gorm:"column:room_id;primary_key;size:128;index:room_id;index:status,priority:1"`
	MeetingName string    `gorm:"column:meeting_name;size:64"`
	HostUserID  string    `gorm:"column:host_user_id;size:64;index:host_user_id"`
	Status      int64     `gorm:"column:status;index:status,priority:2"`
	StartTime   int64     `gorm:"column:start_time"`
	EndTime     int64     `gorm:"column:end_time"`
	CreateTime  time.Time `gorm:"column:create_time"`
	Ex          string    `gorm:"column:ex;size:1024"`
}

func (MeetingInfo) TableName() string {
	return "meeting"
}

type MeetingInvitationInfo struct {
	RoomID     string    `gorm:"column:room_id;primary_key;size:128"`
	UserID     string    `gorm:"column:user_id;primary_key;size:64;index:user_id"`
	CreateTime time.Time `gorm:"column:create_time"`
}

func (MeetingInvitationInfo) TableName() string {
	return "meeting_invitation"
}

type MeetingVideoRecord struct {
	RoomID     string    `gorm:"column:room_id;size:128"`
	FileURL    string    `gorm:"column:file_url"`
	CreateTime time.Time `gorm:"column:create_time"`
}

func (MeetingVideoRecord) TableName() string {
	return "meeting_video_record"
}
