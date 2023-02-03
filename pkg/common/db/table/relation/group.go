package relation

import "time"

const (
	GroupModelTableName = "groups"
)

type GroupModel struct {
	GroupID                string    `gorm:"column:group_id;primary_key;size:64" json:"groupID" binding:"required"`
	GroupName              string    `gorm:"column:name;size:255" json:"groupName"`
	Notification           string    `gorm:"column:notification;size:255" json:"notification"`
	Introduction           string    `gorm:"column:introduction;size:255" json:"introduction"`
	FaceURL                string    `gorm:"column:face_url;size:255" json:"faceURL"`
	CreateTime             time.Time `gorm:"column:create_time;index:create_time;autoCreateTime"`
	Ex                     string    `gorm:"column:ex" json:"ex;size:1024" json:"ex"`
	Status                 int32     `gorm:"column:status"`
	CreatorUserID          string    `gorm:"column:creator_user_id;size:64"`
	GroupType              int32     `gorm:"column:group_type"`
	NeedVerification       int32     `gorm:"column:need_verification"`
	LookMemberInfo         int32     `gorm:"column:look_member_info" json:"lookMemberInfo"`
	ApplyMemberFriend      int32     `gorm:"column:apply_member_friend" json:"applyMemberFriend"`
	NotificationUpdateTime time.Time `gorm:"column:notification_update_time"`
	NotificationUserID     string    `gorm:"column:notification_user_id;size:64"`
}

func (GroupModel) TableName() string {
	return GroupModelTableName
}

type GroupModelInterface interface {
}
