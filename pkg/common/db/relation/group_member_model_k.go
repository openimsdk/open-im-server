package relation

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db/table/relation"
	"Open_IM/pkg/common/tracelog"
	"Open_IM/pkg/utils"
	"context"
	"gorm.io/gorm"
)

type GroupMemberGorm struct {
	DB *gorm.DB
}

func NewGroupMemberDB(db *gorm.DB) *GroupMemberGorm {
	return &GroupMemberGorm{DB: db}
}

func (g *GroupMemberGorm) Create(ctx context.Context, groupMemberList []*relation.GroupMemberModel, tx ...*gorm.DB) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupMemberList", groupMemberList)
	}()
	return utils.Wrap(getDBConn(g.DB, tx).Create(&groupMemberList).Error, "")
}

func (g *GroupMemberGorm) Delete(ctx context.Context, groupMembers []*relation.GroupMemberModel, tx ...*gorm.DB) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupMembers", groupMembers)
	}()
	return utils.Wrap(getDBConn(g.DB, tx).Delete(groupMembers).Error, "")
}

func (g *GroupMemberGorm) DeleteGroup(ctx context.Context, groupIDs []string, tx ...*gorm.DB) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupIDs", groupIDs)
	}()
	return utils.Wrap(getDBConn(g.DB, tx).Where("group_id in (?)", groupIDs).Delete(&relation.GroupMemberModel{}).Error, "")
}

func (g *GroupMemberGorm) UpdateByMap(ctx context.Context, groupID string, userID string, args map[string]interface{}, tx ...*gorm.DB) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupID", groupID, "userID", userID, "args", args)
	}()
	return utils.Wrap(getDBConn(g.DB, tx).Model(&relation.GroupMemberModel{}).Where("group_id = ? and user_id = ?", groupID, userID).Updates(args).Error, "")
}

func (g *GroupMemberGorm) Update(ctx context.Context, groupMembers []*relation.GroupMemberModel, tx ...*gorm.DB) (err error) {
	defer func() { tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupMembers", groupMembers) }()
	return utils.Wrap(getDBConn(g.DB, tx).Updates(&groupMembers).Error, "")
}

func (g *GroupMemberGorm) Find(ctx context.Context, groupIDs []string, userIDs []string, roleLevels []int32, tx ...*gorm.DB) (groupList []*relation.GroupMemberModel, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupIDs", groupIDs, "userIDs", userIDs, "groupList", groupList)
	}()
	db := getDBConn(g.DB, tx)
	if len(groupIDs) > 0 {
		db = db.Where("group_id in (?)", groupIDs)
	}
	if len(userIDs) > 0 {
		db = db.Where("user_id in (?)", userIDs)
	}
	if len(roleLevels) > 0 {
		db = db.Where("role_level in (?)", roleLevels)
	}
	return groupList, utils.Wrap(db.Find(&groupList).Error, "")
}

//func (g *GroupMemberGorm) Find(ctx context.Context, groupMembers []*relation.GroupMemberModel, tx ...*gorm.DB) (groupList []*relation.GroupMemberModel, err error) {
//	defer func() {
//		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupMembers", groupMembers, "groupList", groupList)
//	}()
//	var where [][]interface{}
//	for _, groupMember := range groupMembers {
//		where = append(where, []interface{}{groupMember.GroupID, groupMember.UserID})
//	}
//	return groupList, utils.Wrap(getDBConn(g.DB, tx).Where("(group_id, user_id) in ?", where).Find(&groupList).Error, "")
//}

func (g *GroupMemberGorm) FindGroupUser(ctx context.Context, groupIDs []string, userIDs []string, roleLevels []int32, tx ...*gorm.DB) (groupList []*relation.GroupMemberModel, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupIDs", groupIDs, "userIDs", userIDs, "groupList", groupList)
	}()
	db := getDBConn(g.DB, tx)
	if len(groupList) > 0 {
		db = db.Where("group_id in (?)", groupIDs)
	}
	if len(userIDs) > 0 {
		db = db.Where("user_id in (?)", userIDs)
	}
	if len(roleLevels) > 0 {
		db = db.Where("role_level in (?)", roleLevels)
	}
	return groupList, utils.Wrap(db.Find(&groupList).Error, "")
}

func (g *GroupMemberGorm) Take(ctx context.Context, groupID string, userID string, tx ...*gorm.DB) (groupMember *relation.GroupMemberModel, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupID", groupID, "userID", userID, "groupMember", *groupMember)
	}()
	groupMember = &relation.GroupMemberModel{}
	return groupMember, utils.Wrap(getDBConn(g.DB, tx).Where("group_id = ? and user_id = ?", groupID, userID).Take(groupMember).Error, "")
}

func (g *GroupMemberGorm) TakeOwner(ctx context.Context, groupID string, tx ...*gorm.DB) (groupMember *relation.GroupMemberModel, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupID", groupID, "groupMember", *groupMember)
	}()
	groupMember = &relation.GroupMemberModel{}
	return groupMember, utils.Wrap(getDBConn(g.DB, tx).Where("group_id = ? and role_level = ?", groupID, constant.GroupOwner).Take(groupMember).Error, "")
}

func (g *GroupMemberGorm) SearchMember(ctx context.Context, keyword string, groupIDs []string, userIDs []string, roleLevels []int32, pageNumber, showNumber int32, tx ...*gorm.DB) (total int32, groupList []*relation.GroupMemberModel, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "keyword", keyword, "groupIDs", groupIDs, "userIDs", userIDs, "roleLevels", roleLevels, "pageNumber", pageNumber, "showNumber", showNumber, "total", total, "groupList", groupList)
	}()
	db := getDBConn(g.DB, tx)
	gormIn(&db, "group_id", groupIDs)
	gormIn(&db, "user_id", userIDs)
	gormIn(&db, "role_level", roleLevels)
	return gormSearch[relation.GroupMemberModel](db, []string{"nickname"}, keyword, pageNumber, showNumber)
}
