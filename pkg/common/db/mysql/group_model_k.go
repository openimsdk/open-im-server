package mysql

import (
	"Open_IM/pkg/common/trace_log"
	"Open_IM/pkg/utils"
	"context"
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
	DB                     *gorm.DB
}

func NewGroupDB() *Group {
	var group Group
	db := ConnectToDB()
	db = InitModel(db, &group)
	return &group
}

func (*Group) Create(ctx context.Context, groups []*Group) (err error) {
	defer func() {
		trace_log.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groups", groups)
	}()
	err = utils.Wrap(GroupDB.Create(&groups).Error, "")
	return err
}

func (g *Group) Delete(ctx context.Context, groupIDs []string, tx ...*gorm.DB) (err error) {
	defer func() {
		trace_log.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupIDs", groupIDs)
	}()
	return utils.Wrap(getDBConn(g.DB, tx...).Where("group_id in (?)", groupIDs).Delete(&Group{}).Error, "")
}

func (*Group) UpdateByMap(ctx context.Context, groupID string, args map[string]interface{}) (err error) {
	defer func() {
		trace_log.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupID", groupID, "args", args)
	}()
	return utils.Wrap(GroupDB.Where("group_id = ?", groupID).Updates(args).Error, "")
}

func (*Group) Update(ctx context.Context, groups []*Group) (err error) {
	defer func() {
		trace_log.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groups", groups)
	}()
	return utils.Wrap(GroupDB.Updates(&groups).Error, "")
}

func (*Group) Find(ctx context.Context, groupIDs []string) (groups []*Group, err error) {
	defer func() {
		trace_log.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupIDs", groupIDs, "groups", groups)
	}()
	err = utils.Wrap(GroupDB.Where("group_id in (?)", groupIDs).Find(&groups).Error, "")
	return groups, err
}

func (*Group) Take(ctx context.Context, groupID string) (group *Group, err error) {
	group = &Group{}
	defer func() {
		trace_log.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupID", groupID, "group", *group)
	}()
	err = utils.Wrap(GroupDB.Where("group_id = ?", groupID).Take(group).Error, "")
	return group, err
}
