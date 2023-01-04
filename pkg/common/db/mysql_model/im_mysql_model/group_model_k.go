package im_mysql_model

import (
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/trace_log"
	"Open_IM/pkg/utils"
	"context"
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
}

func (*Group) Create(ctx context.Context, groupList []*Group) (err error) {
	defer trace_log.SetContextInfo(ctx, utils.GetSelfFuncName(), err, "groupList", groupList)
	err = utils.Wrap(db.DB.MysqlDB.DefaultGormDB().Create(&groupList).Error, "")
	return err
}

func (*Group) Delete(ctx context.Context, groupIDList []string) (err error) {
	defer trace_log.SetContextInfo(ctx, utils.GetSelfFuncName(), err, "groupIDList", groupIDList)
	err = utils.Wrap(db.DB.MysqlDB.DefaultGormDB().Where("group_id in (?)", groupIDList).Delete(&Group{}).Error, "")
	return err
}

func (*Group) UpdateByMap(ctx context.Context, groupID string, args map[string]interface{}) (err error) {
	defer trace_log.SetContextInfo(ctx, utils.GetSelfFuncName(), err, "groupID", groupID, "args", args)
	err = utils.Wrap(db.DB.MysqlDB.DefaultGormDB().Where("group_id = ?", groupID).Updates(args).Error, "")
	return err
}

func (*Group) Update(ctx context.Context, groups []*Group) (err error) {
	defer trace_log.SetContextInfo(ctx, utils.GetSelfFuncName(), err, "groups", groups)
	return utils.Wrap(db.DB.MysqlDB.DefaultGormDB().Updates(&groups).Error, "")
}

func (*Group) Find(ctx context.Context, groupIDList []string) (groupList []*Group, err error) {
	defer trace_log.SetContextInfo(ctx, utils.GetSelfFuncName(), err, "groupIDList", groupIDList)
	err = utils.Wrap(db.DB.MysqlDB.DefaultGormDB().Where("group_id in (?)", groupIDList).Find(&groupList).Error, "")
	return groupList, err
}

func (*Group) Take(ctx context.Context, groupID string) (group *Group, err error) {
	defer trace_log.SetContextInfo(ctx, utils.GetSelfFuncName(), err, "groupID", groupID)
	group = &Group{}
	err = utils.Wrap(db.DB.MysqlDB.DefaultGormDB().Where("group_id = ?", groupID).Take(group).Error, "")
	return group, err
}
