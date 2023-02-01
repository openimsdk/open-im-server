package relation

import (
	"Open_IM/pkg/common/db/table"
	"Open_IM/pkg/common/tracelog"
	"Open_IM/pkg/utils"
	"context"
	"gorm.io/gorm"
)

type GroupRequestGorm struct {
	DB *gorm.DB
}

func NewGroupRequest(db *gorm.DB) *GroupRequestGorm {
	return &GroupRequestGorm{
		DB: db,
	}
}

func (g *GroupRequestGorm) Create(ctx context.Context, groupRequests []*table.GroupRequestModel, tx ...*gorm.DB) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupRequests", groupRequests)
	}()
	return utils.Wrap(getDBConn(g.DB, tx).Create(&groupRequests).Error, utils.GetSelfFuncName())
}

func (g *GroupRequestGorm) Delete(ctx context.Context, groupRequests []*table.GroupRequestModel, tx ...*gorm.DB) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupRequests", groupRequests)
	}()
	return utils.Wrap(getDBConn(g.DB, tx).Delete(&groupRequests).Error, utils.GetSelfFuncName())
}

func (g *GroupRequestGorm) UpdateByMap(ctx context.Context, groupID string, userID string, args map[string]interface{}, tx ...*gorm.DB) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupID", groupID, "userID", userID, "args", args)
	}()
	return utils.Wrap(getDBConn(g.DB, tx).Where("group_id = ? and user_id = ? ", groupID, userID).Updates(args).Error, utils.GetSelfFuncName())
}

func (g *GroupRequestGorm) Update(ctx context.Context, groupRequests []*table.GroupRequestModel, tx ...*gorm.DB) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupRequests", groupRequests)
	}()
	return utils.Wrap(getDBConn(g.DB, tx).Updates(&groupRequests).Error, utils.GetSelfFuncName())
}

func (g *GroupRequestGorm) Find(ctx context.Context, groupRequests []*table.GroupRequestModel, tx ...*gorm.DB) (resultGroupRequests []*table.GroupRequestModel, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupRequests", groupRequests, "resultGroupRequests", resultGroupRequests)
	}()
	var where [][]interface{}
	for _, groupMember := range groupRequests {
		where = append(where, []interface{}{groupMember.GroupID, groupMember.UserID})
	}
	return resultGroupRequests, utils.Wrap(getDBConn(g.DB, tx).Where("(group_id, user_id) in ?", where).Find(&resultGroupRequests).Error, utils.GetSelfFuncName())
}

func (g *GroupRequestGorm) Take(ctx context.Context, groupID string, userID string, tx ...*gorm.DB) (groupRequest *table.GroupRequestModel, err error) {
	groupRequest = &table.GroupRequestModel{}
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupID", groupID, "userID", userID, "groupRequest", *groupRequest)
	}()
	return groupRequest, utils.Wrap(getDBConn(g.DB, tx).Where("group_id = ? and user_id = ? ", groupID, userID).Take(groupRequest).Error, utils.GetSelfFuncName())
}
