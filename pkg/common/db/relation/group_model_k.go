package relation

import (
	"Open_IM/pkg/common/tracelog"
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

func NewGroupDB(db *gorm.DB) *Group {
	var group Group
	group.DB = db
	return &group
}

func (g *Group) Create(ctx context.Context, groups []*Group, tx ...*gorm.DB) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groups", groups)
	}()
	return utils.Wrap(getDBConn(g.DB, tx).Create(&groups).Error, "")
}

func (g *Group) Delete(ctx context.Context, groupIDs []string, tx ...*gorm.DB) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupIDs", groupIDs)
	}()
	return utils.Wrap(getDBConn(g.DB, tx).Where("group_id in (?)", groupIDs).Delete(&Group{}).Error, "")
}

func (g *Group) UpdateByMap(ctx context.Context, groupID string, args map[string]interface{}, tx ...*gorm.DB) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupID", groupID, "args", args)
	}()
	return utils.Wrap(getDBConn(g.DB, tx).Where("group_id = ?", groupID).Model(g).Updates(args).Error, "")
}

func (g *Group) Update(ctx context.Context, groups []*Group, tx ...*gorm.DB) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groups", groups)
	}()
	return utils.Wrap(getDBConn(g.DB, tx).Updates(&groups).Error, "")
}

func (g *Group) Find(ctx context.Context, groupIDs []string, tx ...*gorm.DB) (groups []*Group, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupIDs", groupIDs, "groups", groups)
	}()
	return groups, utils.Wrap(getDBConn(g.DB, tx).Where("group_id in (?)", groupIDs).Find(&groups).Error, "")
}

func (g *Group) Take(ctx context.Context, groupID string, tx ...*gorm.DB) (group *Group, err error) {
	group = &Group{}
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupID", groupID, "group", *group)
	}()
	return group, utils.Wrap(getDBConn(g.DB, tx).Where("group_id = ?", groupID).Take(group).Error, "")
}
