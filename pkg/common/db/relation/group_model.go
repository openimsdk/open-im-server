package relation

import (
	"OpenIM/pkg/common/db/table/relation"
	"OpenIM/pkg/common/tracelog"
	"OpenIM/pkg/utils"
	"context"
	"gorm.io/gorm"
)

var _ relation.GroupModelInterface = (*GroupGorm)(nil)

type GroupGorm struct {
	DB *gorm.DB
}

func NewGroupDB(db *gorm.DB) relation.GroupModelInterface {
	return &GroupGorm{DB: db}
}

func (g *GroupGorm) NewTx(tx any) relation.GroupModelInterface {
	return &GroupGorm{DB: tx.(*gorm.DB)}
}

func (g *GroupGorm) Create(ctx context.Context, groups []*relation.GroupModel) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groups", groups)
	}()
	return utils.Wrap(g.DB.Create(&groups).Error, "")
}

func (g *GroupGorm) UpdateMap(ctx context.Context, groupID string, args map[string]interface{}) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupID", groupID, "args", args)
	}()
	return utils.Wrap(g.DB.Where("group_id = ?", groupID).Model(&relation.GroupModel{}).Updates(args).Error, "")
}

func (g *GroupGorm) UpdateStatus(ctx context.Context, groupID string, status int32) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupID", groupID, "status", status)
	}()
	return utils.Wrap(g.DB.Where("group_id = ?", groupID).Model(&relation.GroupModel{}).Updates(map[string]any{"status": status}).Error, "")
}

func (g *GroupGorm) Find(ctx context.Context, groupIDs []string) (groups []*relation.GroupModel, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupIDs", groupIDs, "groups", groups)
	}()
	return groups, utils.Wrap(g.DB.Where("group_id in (?)", groupIDs).Find(&groups).Error, "")
}

func (g *GroupGorm) Take(ctx context.Context, groupID string) (group *relation.GroupModel, err error) {
	group = &relation.GroupModel{}
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupID", groupID, "group", group)
	}()
	return group, utils.Wrap(g.DB.Where("group_id = ?", groupID).Take(group).Error, "")
}

func (g *GroupGorm) Search(ctx context.Context, keyword string, pageNumber, showNumber int32) (total uint32, groups []*relation.GroupModel, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "keyword", keyword, "pageNumber", pageNumber, "showNumber", showNumber, "total", total, "groups", groups)
	}()
	return gormSearch[relation.GroupModel](g.DB, []string{"name"}, keyword, pageNumber, showNumber)
}

func (g *GroupGorm) GetGroupIDsByGroupType(ctx context.Context, groupType int) (groupIDs []string, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupType", groupType, "groupIDs", groupIDs)
	}()
	return groupIDs, utils.Wrap(g.DB.Model(&relation.GroupModel{}).Where("group_type = ? ", groupType).Pluck("group_id", &groupIDs).Error, "")
}
