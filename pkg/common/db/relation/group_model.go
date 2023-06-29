package relation

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/ormutil"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"gorm.io/gorm"
)

var _ relation.GroupModelInterface = (*GroupGorm)(nil)

type GroupGorm struct {
	*MetaDB
}

func NewGroupDB(db *gorm.DB) relation.GroupModelInterface {
	return &GroupGorm{NewMetaDB(db, &relation.GroupModel{})}
}

func (g *GroupGorm) NewTx(tx any) relation.GroupModelInterface {
	return &GroupGorm{NewMetaDB(tx.(*gorm.DB), &relation.GroupModel{})}
}

func (g *GroupGorm) Create(ctx context.Context, groups []*relation.GroupModel) (err error) {
	return utils.Wrap(g.DB.Create(&groups).Error, "")
}

func (g *GroupGorm) UpdateMap(ctx context.Context, groupID string, args map[string]interface{}) (err error) {
	return utils.Wrap(g.DB.Where("group_id = ?", groupID).Model(&relation.GroupModel{}).Updates(args).Error, "")
}

func (g *GroupGorm) UpdateStatus(ctx context.Context, groupID string, status int32) (err error) {
	return utils.Wrap(g.DB.Where("group_id = ?", groupID).Model(&relation.GroupModel{}).Updates(map[string]any{"status": status}).Error, "")
}

func (g *GroupGorm) Find(ctx context.Context, groupIDs []string) (groups []*relation.GroupModel, err error) {
	return groups, utils.Wrap(g.DB.Where("group_id in (?)", groupIDs).Find(&groups).Error, "")
}

func (g *GroupGorm) Take(ctx context.Context, groupID string) (group *relation.GroupModel, err error) {
	group = &relation.GroupModel{}
	return group, utils.Wrap(g.DB.Where("group_id = ?", groupID).Take(group).Error, "")
}

func (g *GroupGorm) Search(ctx context.Context, keyword string, pageNumber, showNumber int32) (total uint32, groups []*relation.GroupModel, err error) {
	return ormutil.GormSearch[relation.GroupModel](g.DB, []string{"name"}, keyword, pageNumber, showNumber)
}

func (g *GroupGorm) GetGroupIDsByGroupType(ctx context.Context, groupType int) (groupIDs []string, err error) {
	return groupIDs, utils.Wrap(g.DB.Model(&relation.GroupModel{}).Where("group_type = ? ", groupType).Pluck("group_id", &groupIDs).Error, "")
}
