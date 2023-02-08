package relation

import (
	"Open_IM/pkg/common/db/table/relation"
	"Open_IM/pkg/common/tracelog"
	"Open_IM/pkg/utils"
	"context"
	"gorm.io/gorm"
)

var _ relation.GroupRequestModelInterface = (*GroupRequestGorm)(nil)

type GroupRequestGorm struct {
	DB *gorm.DB
}

func NewGroupRequest(db *gorm.DB) relation.GroupRequestModelInterface {
	return &GroupRequestGorm{
		DB: db,
	}
}

func (g *GroupRequestGorm) Create(ctx context.Context, groupRequests []*relation.GroupRequestModel, tx ...any) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupRequests", groupRequests)
	}()
	return utils.Wrap(getDBConn(g.DB, tx).Create(&groupRequests).Error, utils.GetSelfFuncName())
}

//func (g *GroupRequestGorm) Delete(ctx context.Context, groupRequests []*relation.GroupRequestModel, tx ...any) (err error) {
//	defer func() {
//		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupRequests", groupRequests)
//	}()
//	return utils.Wrap(getDBConn(g.DB, tx).Delete(&groupRequests).Error, utils.GetSelfFuncName())
//}

//func (g *GroupRequestGorm) UpdateMap(ctx context.Context, groupID string, userID string, args map[string]interface{}, tx ...any) (err error) {
//	defer func() {
//		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupID", groupID, "userID", userID, "args", args)
//	}()
//	return utils.Wrap(getDBConn(g.DB, tx).Model(&relation.GroupRequestModel{}).Where("group_id = ? and user_id = ? ", groupID, userID).Updates(args).Error, utils.GetSelfFuncName())
//}

func (g *GroupRequestGorm) UpdateHandler(ctx context.Context, groupID string, userID string, handledMsg string, handleResult int32, tx ...any) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupID", groupID, "userID", userID, "handledMsg", handledMsg, "handleResult", handleResult)
	}()
	return utils.Wrap(getDBConn(g.DB, tx).Model(&relation.GroupRequestModel{}).Where("group_id = ? and user_id = ? ", groupID, userID).Updates(map[string]any{
		"handle_msg":    handledMsg,
		"handle_result": handleResult,
	}).Error, utils.GetSelfFuncName())
}

//func (g *GroupRequestGorm) Update(ctx context.Context, groupRequests []*relation.GroupRequestModel, tx ...any) (err error) {
//	defer func() {
//		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupRequests", groupRequests)
//	}()
//	return utils.Wrap(getDBConn(g.DB, tx).Updates(&groupRequests).Error, utils.GetSelfFuncName())
//}

//func (g *GroupRequestGorm) Find(ctx context.Context, groupRequests []*relation.GroupRequestModel, tx ...any) (resultGroupRequests []*relation.GroupRequestModel, err error) {
//	defer func() {
//		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupRequests", groupRequests, "resultGroupRequests", resultGroupRequests)
//	}()
//	var where [][]interface{}
//	for _, groupMember := range groupRequests {
//		where = append(where, []interface{}{groupMember.GroupID, groupMember.UserID})
//	}
//	return resultGroupRequests, utils.Wrap(getDBConn(g.DB, tx).Where("(group_id, user_id) in ?", where).Find(&resultGroupRequests).Error, utils.GetSelfFuncName())
//}

func (g *GroupRequestGorm) Take(ctx context.Context, groupID string, userID string, tx ...any) (groupRequest *relation.GroupRequestModel, err error) {
	groupRequest = &relation.GroupRequestModel{}
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupID", groupID, "userID", userID, "groupRequest", *groupRequest)
	}()
	return groupRequest, utils.Wrap(getDBConn(g.DB, tx).Where("group_id = ? and user_id = ? ", groupID, userID).Take(groupRequest).Error, utils.GetSelfFuncName())
}

func (g *GroupRequestGorm) Page(ctx context.Context, userID string, pageNumber, showNumber int32, tx ...any) (total int32, groups []*relation.GroupRequestModel, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "pageNumber", pageNumber, "showNumber", showNumber, "total", total, "groups", groups)
	}()
	return gormSearch[relation.GroupRequestModel](getDBConn(g.DB, tx).Where("user_id = ?", userID), nil, "", pageNumber, showNumber)
}
