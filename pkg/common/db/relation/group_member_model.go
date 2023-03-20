package relation

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/ormutil"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"gorm.io/gorm"
)

var _ relation.GroupMemberModelInterface = (*GroupMemberGorm)(nil)

type GroupMemberGorm struct {
	*MetaDB
}

func NewGroupMemberDB(db *gorm.DB) relation.GroupMemberModelInterface {
	return &GroupMemberGorm{NewMetaDB(db, &relation.GroupMemberModel{})}
}

func (g *GroupMemberGorm) NewTx(tx any) relation.GroupMemberModelInterface {
	return &GroupMemberGorm{NewMetaDB(tx.(*gorm.DB), &relation.GroupMemberModel{})}
}

func (g *GroupMemberGorm) Create(ctx context.Context, groupMemberList []*relation.GroupMemberModel) (err error) {
	return utils.Wrap(g.DB.Create(&groupMemberList).Error, "")
}

func (g *GroupMemberGorm) Delete(ctx context.Context, groupID string, userIDs []string) (err error) {
	return utils.Wrap(g.DB.Where("group_id = ? and user_id in (?)", groupID, userIDs).Delete(&relation.GroupMemberModel{}).Error, "")
}

func (g *GroupMemberGorm) DeleteGroup(ctx context.Context, groupIDs []string) (err error) {
	return utils.Wrap(g.DB.Where("group_id in (?)", groupIDs).Delete(&relation.GroupMemberModel{}).Error, "")
}

func (g *GroupMemberGorm) Update(ctx context.Context, groupID string, userID string, data map[string]any) (err error) {
	return utils.Wrap(g.DB.Model(&relation.GroupMemberModel{}).Where("group_id = ? and user_id = ?", groupID, userID).Updates(data).Error, "")
}

func (g *GroupMemberGorm) UpdateRoleLevel(ctx context.Context, groupID string, userID string, roleLevel int32) (rowsAffected int64, err error) {
	db := g.DB.Model(&relation.GroupMemberModel{}).Where("group_id = ? and user_id = ?", groupID, userID).Updates(map[string]any{
		"role_level": roleLevel,
	})
	return db.RowsAffected, utils.Wrap(db.Error, "")
}

func (g *GroupMemberGorm) Find(ctx context.Context, groupIDs []string, userIDs []string, roleLevels []int32) (groupList []*relation.GroupMemberModel, err error) {
	db := g.DB
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

func (g *GroupMemberGorm) Take(ctx context.Context, groupID string, userID string) (groupMember *relation.GroupMemberModel, err error) {
	groupMember = &relation.GroupMemberModel{}
	return groupMember, utils.Wrap(g.DB.Where("group_id = ? and user_id = ?", groupID, userID).Take(groupMember).Error, "")
}

func (g *GroupMemberGorm) TakeOwner(ctx context.Context, groupID string) (groupMember *relation.GroupMemberModel, err error) {
	groupMember = &relation.GroupMemberModel{}
	return groupMember, utils.Wrap(g.DB.Where("group_id = ? and role_level = ?", groupID, constant.GroupOwner).Take(groupMember).Error, "")
}

func (g *GroupMemberGorm) SearchMember(ctx context.Context, keyword string, groupIDs []string, userIDs []string, roleLevels []int32, pageNumber, showNumber int32) (total uint32, groupList []*relation.GroupMemberModel, err error) {
	db := g.DB
	ormutil.GormIn(&db, "group_id", groupIDs)
	ormutil.GormIn(&db, "user_id", userIDs)
	ormutil.GormIn(&db, "role_level", roleLevels)
	return ormutil.GormSearch[relation.GroupMemberModel](db, []string{"nickname"}, keyword, pageNumber, showNumber)
}

func (g *GroupMemberGorm) MapGroupMemberNum(ctx context.Context, groupIDs []string) (count map[string]uint32, err error) {
	return ormutil.MapCount(g.DB.Where("group_id in (?)", groupIDs), "group_id")
}

func (g *GroupMemberGorm) FindJoinUserID(ctx context.Context, groupIDs []string) (groupUsers map[string][]string, err error) {
	var items []struct {
		GroupID string `gorm:"group_id"`
		UserID  string `gorm:"user_id"`
	}
	if err := g.DB.Model(&relation.GroupMemberModel{}).Where("group_id in (?)", groupIDs).Find(&items).Error; err != nil {
		return nil, utils.Wrap(err, "")
	}
	groupUsers = make(map[string][]string)
	for _, item := range items {
		groupUsers[item.GroupID] = append(groupUsers[item.GroupID], item.UserID)
	}
	return groupUsers, nil
}

func (g *GroupMemberGorm) FindMemberUserID(ctx context.Context, groupID string) (userIDs []string, err error) {
	return userIDs, utils.Wrap(g.DB.Model(&relation.GroupMemberModel{}).Where("group_id = ?", groupID).Pluck("user_id", &userIDs).Error, "")
}
