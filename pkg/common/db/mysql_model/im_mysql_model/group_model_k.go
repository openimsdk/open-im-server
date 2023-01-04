package im_mysql_model

import (
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/utils"
	"gorm.io/gorm"
	"time"
)

type Group struct {
	GroupID                string    `gorm:"column:group_id;primary_key;size:64" json:"groupID" binding:"required"`
	GroupName              string    `gorm:"column:name;size:255" json:"groupName"`
	Notification           string    `gorm:"column:notification;size:255" json:"notification"`
	Introduction           string    `gorm:"column:introduction;size:255" json:"introduction"`
	FaceURL                string    `gorm:"column:face_url;size:255" json:"faceURL"`
	CreateTime             time.Time `gorm:"column:create_time;index:create_time"`
	Ex                     string    `gorm:"column:ex" json:"ex;size:1024" json:"ex"`
	Status                 int32     `gorm:"column:status"`
	CreatorUserID          string    `gorm:"column:creator_user_id;size:64"`
	GroupType              int32     `gorm:"column:group_type"`
	NeedVerification       int32     `gorm:"column:need_verification"`
	LookMemberInfo         int32     `gorm:"column:look_member_info" json:"lookMemberInfo"`
	ApplyMemberFriend      int32     `gorm:"column:apply_member_friend" json:"applyMemberFriend"`
	NotificationUpdateTime time.Time `gorm:"column:notification_update_time"`
	NotificationUserID     string    `gorm:"column:notification_user_id;size:64"`
	DB                     *gorm.DB  `gorm:"-" json:"-"`
}

func (tb *Group) Create(groupList []*Group) error {
	return utils.Wrap(db.DB.MysqlDB.DefaultGormDB().Create(&groupList).Error, "")
}

func (tb *Group) Delete(groupIDList []string) error {
	return utils.Wrap(db.DB.MysqlDB.DefaultGormDB().Where("group_id in (?)", groupIDList).Delete(&Group{}).Error, "")
}

func (tb *Group) Get(groupIDs []string) ([]*Group, error) {
	var ms []*Group
	return ms, utils.Wrap(db.DB.MysqlDB.DefaultGormDB().Where("group_id in (?)", groupIDs).Find(&ms).Error, "")
}

func (tb *Group) Update(groups []*Group) error {
	return utils.Wrap(utils.Wrap(db.DB.MysqlDB.DefaultGormDB().Updates(groups).Error, ""), "")
}

func (*Group) Find(groupIDList []string) ([]*Group, error) {
	return nil, nil
}

func (*Group) Take(groupID string) (*Group, error) {
	return nil, nil
}
