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

func (*Group) Create(ctx context.Context, groups []*Group) (err error) {
	defer func() {
		trace_log.SetContextInfo(ctx, utils.GetFuncName(1), err, "groups", groups)
	}()
	err = utils.Wrap(db.DB.MysqlDB.DefaultGormDB().Create(&groups).Error, "")
	return err
}

func (*Group) Delete(ctx context.Context, groupIDs []string) (err error) {
	defer func() {
		trace_log.SetContextInfo(ctx, utils.GetFuncName(1), err, "groupIDs", groupIDs)
	}()
	err = utils.Wrap(db.DB.MysqlDB.DefaultGormDB().Where("group_id in (?)", groupIDs).Delete(&Group{}).Error, "")
	return err
}

func (*Group) UpdateByMap(ctx context.Context, groupID string, args map[string]interface{}) (err error) {
	defer func() {
		trace_log.SetContextInfo(ctx, utils.GetFuncName(1), err, "groupID", groupID, "args", args)
	}()
	return utils.Wrap(db.DB.MysqlDB.DefaultGormDB().Where("group_id = ?", groupID).Updates(args).Error, "")
}

func (*Group) Update(ctx context.Context, groups []*Group) (err error) {
	defer func() {
		trace_log.SetContextInfo(ctx, utils.GetFuncName(1), err, "groups", groups)
	}()
	return utils.Wrap(db.DB.MysqlDB.DefaultGormDB().Updates(&groups).Error, "")
}

func (*Group) Find(ctx context.Context, groupIDs []string) (groupList []*Group, err error) {
	defer func() {
		trace_log.SetContextInfo(ctx, utils.GetFuncName(1), err, "groupIDList", groupIDs, "groupList", groupList)
	}()
	err = utils.Wrap(db.DB.MysqlDB.DefaultGormDB().Where("group_id in (?)", groupIDs).Find(&groupList).Error, "")
	return groupList, err
}

func (*Group) Take(ctx context.Context, groupID string) (group *Group, err error) {
	group = &Group{}
	defer func() {
		trace_log.SetContextInfo(ctx, utils.GetFuncName(1), err, "groupID", groupID, "group", *group)
	}()
	err = utils.Wrap(db.DB.MysqlDB.DefaultGormDB().Where("group_id = ?", groupID).Take(group).Error, "")
	return group, err
}
