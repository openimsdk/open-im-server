package relation

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/ormutil"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"gorm.io/gorm"
)

type GroupRequestGorm struct {
	*MetaDB
}

func NewGroupRequest(db *gorm.DB) relation.GroupRequestModelInterface {
	return &GroupRequestGorm{
		NewMetaDB(db, &relation.GroupRequestModel{}),
	}
}

func (g *GroupRequestGorm) NewTx(tx any) relation.GroupRequestModelInterface {
	return &GroupRequestGorm{NewMetaDB(tx.(*gorm.DB), &relation.GroupRequestModel{})}
}

func (g *GroupRequestGorm) Create(ctx context.Context, groupRequests []*relation.GroupRequestModel) (err error) {
	return utils.Wrap(g.DB.WithContext(ctx).Create(&groupRequests).Error, utils.GetSelfFuncName())
}

func (g *GroupRequestGorm) Delete(ctx context.Context, groupID string, userID string) (err error) {
	return utils.Wrap(g.DB.WithContext(ctx).Where("group_id = ? and user_id = ? ", groupID, userID).Delete(&relation.GroupRequestModel{}).Error, utils.GetSelfFuncName())
}

func (g *GroupRequestGorm) UpdateHandler(ctx context.Context, groupID string, userID string, handledMsg string, handleResult int32) (err error) {
	return utils.Wrap(g.DB.WithContext(ctx).Model(&relation.GroupRequestModel{}).Where("group_id = ? and user_id = ? ", groupID, userID).Updates(map[string]any{
		"handle_msg":    handledMsg,
		"handle_result": handleResult,
	}).Error, utils.GetSelfFuncName())
}

func (g *GroupRequestGorm) Take(ctx context.Context, groupID string, userID string) (groupRequest *relation.GroupRequestModel, err error) {
	groupRequest = &relation.GroupRequestModel{}
	return groupRequest, utils.Wrap(g.DB.WithContext(ctx).Where("group_id = ? and user_id = ? ", groupID, userID).Take(groupRequest).Error, utils.GetSelfFuncName())
}

func (g *GroupRequestGorm) Page(ctx context.Context, userID string, pageNumber, showNumber int32) (total uint32, groups []*relation.GroupRequestModel, err error) {
	return ormutil.GormSearch[relation.GroupRequestModel](g.DB.WithContext(ctx).Where("user_id = ?", userID), nil, "", pageNumber, showNumber)
}

func (g *GroupRequestGorm) PageGroup(ctx context.Context, groupIDs []string, pageNumber, showNumber int32) (total uint32, groups []*relation.GroupRequestModel, err error) {
	return ormutil.GormPage[relation.GroupRequestModel](g.DB.WithContext(ctx).Where("group_id in ?", groupIDs), pageNumber, showNumber)
}
