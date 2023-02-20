package relation

import (
	"Open_IM/pkg/common/db/table/relation"
	"Open_IM/pkg/common/tracelog"
	"Open_IM/pkg/utils"
	"context"
	"gorm.io/gorm"
)

type GroupRequestGorm struct {
	DB *gorm.DB
}

func (g *GroupRequestGorm) NewTx(tx any) relation.GroupRequestModelInterface {
	return &GroupRequestGorm{
		DB: tx.(*gorm.DB),
	}
}

func NewGroupRequest(db *gorm.DB) relation.GroupRequestModelInterface {
	return &GroupRequestGorm{
		DB: db,
	}
}

func (g *GroupRequestGorm) Create(ctx context.Context, groupRequests []*relation.GroupRequestModel) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupRequests", groupRequests)
	}()
	return utils.Wrap(g.DB.Create(&groupRequests).Error, utils.GetSelfFuncName())
}

func (g *GroupRequestGorm) UpdateHandler(ctx context.Context, groupID string, userID string, handledMsg string, handleResult int32) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupID", groupID, "userID", userID, "handledMsg", handledMsg, "handleResult", handleResult)
	}()
	return utils.Wrap(g.DB.Model(&relation.GroupRequestModel{}).Where("group_id = ? and user_id = ? ", groupID, userID).Updates(map[string]any{
		"handle_msg":    handledMsg,
		"handle_result": handleResult,
	}).Error, utils.GetSelfFuncName())
}

func (g *GroupRequestGorm) Take(ctx context.Context, groupID string, userID string) (groupRequest *relation.GroupRequestModel, err error) {
	groupRequest = &relation.GroupRequestModel{}
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupID", groupID, "userID", userID, "groupRequest", *groupRequest)
	}()
	return groupRequest, utils.Wrap(g.DB.Where("group_id = ? and user_id = ? ", groupID, userID).Take(groupRequest).Error, utils.GetSelfFuncName())
}

func (g *GroupRequestGorm) Page(ctx context.Context, userID string, pageNumber, showNumber int32) (total uint32, groups []*relation.GroupRequestModel, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "pageNumber", pageNumber, "showNumber", showNumber, "total", total, "groups", groups)
	}()
	return gormSearch[relation.GroupRequestModel](g.DB.Where("user_id = ?", userID), nil, "", pageNumber, showNumber)
}
