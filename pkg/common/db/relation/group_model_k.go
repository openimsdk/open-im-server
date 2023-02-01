package relation

import (
	"Open_IM/pkg/common/db/table"
	"Open_IM/pkg/common/tracelog"
	"Open_IM/pkg/utils"
	"context"
	"gorm.io/gorm"
)

type GroupGorm struct {
	DB *gorm.DB
}

func NewGroupDB(db *gorm.DB) *GroupGorm {
	return &GroupGorm{DB: db}
}

func (g *GroupGorm) Create(ctx context.Context, groups []*table.GroupModel, tx ...*gorm.DB) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groups", groups)
	}()
	return utils.Wrap(getDBConn(g.DB, tx).Create(&groups).Error, "")
}

func (g *GroupGorm) Delete(ctx context.Context, groupIDs []string, tx ...*gorm.DB) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupIDs", groupIDs)
	}()
	return utils.Wrap(getDBConn(g.DB, tx).Where("group_id in (?)", groupIDs).Delete(&table.GroupModel{}).Error, "")
}

func (g *GroupGorm) UpdateByMap(ctx context.Context, groupID string, args map[string]interface{}, tx ...*gorm.DB) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupID", groupID, "args", args)
	}()
	return utils.Wrap(getDBConn(g.DB, tx).Where("group_id = ?", groupID).Model(g).Updates(args).Error, "")
}

func (g *GroupGorm) Update(ctx context.Context, groups []*table.GroupModel, tx ...*gorm.DB) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groups", groups)
	}()
	return utils.Wrap(getDBConn(g.DB, tx).Updates(&groups).Error, "")
}

func (g *GroupGorm) Find(ctx context.Context, groupIDs []string, tx ...*gorm.DB) (groups []*table.GroupModel, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupIDs", groupIDs, "groups", groups)
	}()
	return groups, utils.Wrap(getDBConn(g.DB, tx).Where("group_id in (?)", groupIDs).Find(&groups).Error, "")
}

func (g *GroupGorm) Take(ctx context.Context, groupID string, tx ...*gorm.DB) (group *table.GroupModel, err error) {
	group = &table.GroupModel{}
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupID", groupID, "group", *group)
	}()
	return group, utils.Wrap(getDBConn(g.DB, tx).Where("group_id = ?", groupID).Take(group).Error, "")
}
