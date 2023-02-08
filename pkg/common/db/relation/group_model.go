package relation

import (
	"Open_IM/pkg/common/db/table/relation"
	"Open_IM/pkg/common/tracelog"
	"Open_IM/pkg/utils"
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

func (g *GroupGorm) Create(ctx context.Context, groups []*relation.GroupModel, tx ...any) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groups", groups)
	}()
	return utils.Wrap(getDBConn(g.DB, tx).Create(&groups).Error, "")
}

//func (g *GroupGorm) Delete(ctx context.Context, groupIDs []string, tx ...any) (err error) {
//	defer func() {
//		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupIDs", groupIDs)
//	}()
//	return utils.Wrap(getDBConn(g.DB, tx).Where("group_id in (?)", groupIDs).Delete(&relation.GroupModel{}).Error, "")
//}

func (g *GroupGorm) UpdateMap(ctx context.Context, groupID string, args map[string]interface{}, tx ...any) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupID", groupID, "args", args)
	}()
	return utils.Wrap(getDBConn(g.DB, tx).Where("group_id = ?", groupID).Model(&relation.GroupModel{}).Updates(args).Error, "")
}

func (g *GroupGorm) UpdateStatus(ctx context.Context, groupID string, status int32, tx ...any) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupID", groupID, "status", status)
	}()
	return utils.Wrap(getDBConn(g.DB, tx).Where("group_id = ?", groupID).Model(&relation.GroupModel{}).Updates(map[string]any{"status": status}).Error, "")
}

func (g *GroupGorm) Find(ctx context.Context, groupIDs []string, tx ...any) (groups []*relation.GroupModel, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupIDs", groupIDs, "groups", groups)
	}()
	return groups, utils.Wrap(getDBConn(g.DB, tx).Where("group_id in (?)", groupIDs).Find(&groups).Error, "")
}

func (g *GroupGorm) Take(ctx context.Context, groupID string, tx ...any) (group *relation.GroupModel, err error) {
	group = &relation.GroupModel{}
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupID", groupID, "group", group)
	}()
	return group, utils.Wrap(getDBConn(g.DB, tx).Where("group_id = ?", groupID).Take(group).Error, "")
}

func (g *GroupGorm) Search(ctx context.Context, keyword string, pageNumber, showNumber int32, tx ...any) (total int32, groups []*relation.GroupModel, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "keyword", keyword, "pageNumber", pageNumber, "showNumber", showNumber, "total", total, "groups", groups)
	}()
	return gormSearch[relation.GroupModel](getDBConn(g.DB, tx), []string{"name"}, keyword, pageNumber, showNumber)
}
