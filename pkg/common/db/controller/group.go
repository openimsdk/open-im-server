package controller

import (
	"context"
	"fmt"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/cache"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/relation"
	relationTb "github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	unRelationTb "github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/unrelation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/tx"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/unrelation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"github.com/dtm-labs/rockscache"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type GroupDatabase interface {
	// Group
	CreateGroup(ctx context.Context, groups []*relationTb.GroupModel, groupMembers []*relationTb.GroupMemberModel) error
	TakeGroup(ctx context.Context, groupID string) (group *relationTb.GroupModel, err error)
	FindGroup(ctx context.Context, groupIDs []string) (groups []*relationTb.GroupModel, err error)
	SearchGroup(ctx context.Context, keyword string, pageNumber, showNumber int32) (uint32, []*relationTb.GroupModel, error)
	UpdateGroup(ctx context.Context, groupID string, data map[string]any) error
	DismissGroup(ctx context.Context, groupID string, deleteMember bool) error // 解散群，并删除群成员
	GetGroupIDsByGroupType(ctx context.Context, groupType int) (groupIDs []string, err error)
	// GroupMember
	TakeGroupMember(ctx context.Context, groupID string, userID string) (groupMember *relationTb.GroupMemberModel, err error)
	TakeGroupOwner(ctx context.Context, groupID string) (*relationTb.GroupMemberModel, error)
	FindGroupMember(ctx context.Context, groupIDs []string, userIDs []string, roleLevels []int32) ([]*relationTb.GroupMemberModel, error)
	FindGroupMemberUserID(ctx context.Context, groupID string) ([]string, error)
	FindGroupMemberNum(ctx context.Context, groupID string) (uint32, error)
	FindUserManagedGroupID(ctx context.Context, userID string) (groupIDs []string, err error)
	PageGroupRequest(ctx context.Context, groupIDs []string, pageNumber, showNumber int32) (uint32, []*relationTb.GroupRequestModel, error)
	//PageGroupMember(ctx context.Context, groupIDs []string, userIDs []string, roleLevels []int32, pageNumber, showNumber int32) (uint32, []*relationTb.GroupMemberModel, error)
	PageGetJoinGroup(ctx context.Context, userID string, pageNumber, showNumber int32) (total uint32, totalGroupMembers []*relationTb.GroupMemberModel, err error)
	PageGetGroupMember(ctx context.Context, groupID string, pageNumber, showNumber int32) (total uint32, totalGroupMembers []*relationTb.GroupMemberModel, err error)
	SearchGroupMember(ctx context.Context, keyword string, groupIDs []string, userIDs []string, roleLevels []int32, pageNumber, showNumber int32) (uint32, []*relationTb.GroupMemberModel, error)
	HandlerGroupRequest(ctx context.Context, groupID string, userID string, handledMsg string, handleResult int32, member *relationTb.GroupMemberModel) error
	DeleteGroupMember(ctx context.Context, groupID string, userIDs []string) error
	MapGroupMemberUserID(ctx context.Context, groupIDs []string) (map[string]*relationTb.GroupSimpleUserID, error)
	MapGroupMemberNum(ctx context.Context, groupIDs []string) (map[string]uint32, error)
	TransferGroupOwner(ctx context.Context, groupID string, oldOwnerUserID, newOwnerUserID string, roleLevel int32) error // 转让群
	UpdateGroupMember(ctx context.Context, groupID string, userID string, data map[string]any) error
	UpdateGroupMembers(ctx context.Context, data []*relationTb.BatchUpdateGroupMember) error
	// GroupRequest
	CreateGroupRequest(ctx context.Context, requests []*relationTb.GroupRequestModel) error
	TakeGroupRequest(ctx context.Context, groupID string, userID string) (*relationTb.GroupRequestModel, error)
	PageGroupRequestUser(ctx context.Context, userID string, pageNumber, showNumber int32) (uint32, []*relationTb.GroupRequestModel, error)
	// SuperGroupModelInterface
	FindSuperGroup(ctx context.Context, groupIDs []string) ([]*unRelationTb.SuperGroupModel, error)
	FindJoinSuperGroup(ctx context.Context, userID string) ([]string, error)
	CreateSuperGroup(ctx context.Context, groupID string, initMemberIDList []string) error
	DeleteSuperGroup(ctx context.Context, groupID string) error
	DeleteSuperGroupMember(ctx context.Context, groupID string, userIDs []string) error
	CreateSuperGroupMember(ctx context.Context, groupID string, userIDs []string) error
}

func NewGroupDatabase(
	group relationTb.GroupModelInterface,
	member relationTb.GroupMemberModelInterface,
	request relationTb.GroupRequestModelInterface,
	tx tx.Tx,
	ctxTx tx.CtxTx,
	superGroup unRelationTb.SuperGroupModelInterface,
	cache cache.GroupCache,
) GroupDatabase {
	database := &groupDatabase{
		groupDB:        group,
		groupMemberDB:  member,
		groupRequestDB: request,
		tx:             tx,
		ctxTx:          ctxTx,
		cache:          cache,
		mongoDB:        superGroup,
	}
	return database
}

func InitGroupDatabase(db *gorm.DB, rdb redis.UniversalClient, database *mongo.Database) GroupDatabase {
	rcOptions := rockscache.NewDefaultOptions()
	rcOptions.StrongConsistency = true
	rcOptions.RandomExpireAdjustment = 0.2
	return NewGroupDatabase(
		relation.NewGroupDB(db),
		relation.NewGroupMemberDB(db),
		relation.NewGroupRequest(db),
		tx.NewGorm(db),
		tx.NewMongo(database.Client()),
		unrelation.NewSuperGroupMongoDriver(database),
		cache.NewGroupCacheRedis(rdb, relation.NewGroupDB(db), relation.NewGroupMemberDB(db), relation.NewGroupRequest(db), unrelation.NewSuperGroupMongoDriver(database), rcOptions),
	)
}

type groupDatabase struct {
	groupDB        relationTb.GroupModelInterface
	groupMemberDB  relationTb.GroupMemberModelInterface
	groupRequestDB relationTb.GroupRequestModelInterface
	tx             tx.Tx
	ctxTx          tx.CtxTx
	cache          cache.GroupCache
	mongoDB        unRelationTb.SuperGroupModelInterface
}

func (g *groupDatabase) GetGroupIDsByGroupType(ctx context.Context, groupType int) (groupIDs []string, err error) {
	return g.groupDB.GetGroupIDsByGroupType(ctx, groupType)
}

func (g *groupDatabase) FindGroupMemberUserID(ctx context.Context, groupID string) ([]string, error) {
	return g.cache.GetGroupMemberIDs(ctx, groupID)
}

func (g *groupDatabase) FindGroupMemberNum(ctx context.Context, groupID string) (uint32, error) {
	num, err := g.cache.GetGroupMemberNum(ctx, groupID)
	if err != nil {
		return 0, err
	}
	return uint32(num), nil
}

func (g *groupDatabase) CreateGroup(ctx context.Context, groups []*relationTb.GroupModel, groupMembers []*relationTb.GroupMemberModel) error {
	var cache = g.cache.NewCache()
	if err := g.tx.Transaction(func(tx any) error {
		if len(groups) > 0 {
			if err := g.groupDB.NewTx(tx).Create(ctx, groups); err != nil {
				return err
			}
		}
		if len(groupMembers) > 0 {
			if err := g.groupMemberDB.NewTx(tx).Create(ctx, groupMembers); err != nil {
				return err
			}
		}
		createGroupIDs := utils.DistinctAnyGetComparable(groups, func(group *relationTb.GroupModel) string {
			return group.GroupID
		})
		m := make(map[string]struct{})

		for _, groupMember := range groupMembers {
			if _, ok := m[groupMember.GroupID]; !ok {
				m[groupMember.GroupID] = struct{}{}
				cache = cache.DelGroupMemberIDs(groupMember.GroupID).DelGroupMembersHash(groupMember.GroupID).DelGroupsMemberNum(groupMember.GroupID)
			}
			cache = cache.DelJoinedGroupID(groupMember.UserID).DelGroupMembersInfo(groupMember.GroupID, groupMember.UserID)
		}
		cache = cache.DelGroupsInfo(createGroupIDs...)
		return nil
	}); err != nil {
		return err
	}
	return cache.ExecDel(ctx)
}

func (g *groupDatabase) TakeGroup(ctx context.Context, groupID string) (group *relationTb.GroupModel, err error) {
	return g.cache.GetGroupInfo(ctx, groupID)
}

func (g *groupDatabase) FindGroup(ctx context.Context, groupIDs []string) (groups []*relationTb.GroupModel, err error) {
	return g.cache.GetGroupsInfo(ctx, groupIDs)
}

func (g *groupDatabase) SearchGroup(ctx context.Context, keyword string, pageNumber, showNumber int32) (uint32, []*relationTb.GroupModel, error) {
	return g.groupDB.Search(ctx, keyword, pageNumber, showNumber)
}

func (g *groupDatabase) UpdateGroup(ctx context.Context, groupID string, data map[string]any) error {
	if err := g.groupDB.UpdateMap(ctx, groupID, data); err != nil {
		return err
	}
	return g.cache.DelGroupsInfo(groupID).ExecDel(ctx)
}

func (g *groupDatabase) DismissGroup(ctx context.Context, groupID string, deleteMember bool) error {
	cache := g.cache.NewCache()
	if err := g.tx.Transaction(func(tx any) error {
		if err := g.groupDB.NewTx(tx).UpdateStatus(ctx, groupID, constant.GroupStatusDismissed); err != nil {
			return err
		}
		if deleteMember {
			if err := g.groupMemberDB.NewTx(tx).DeleteGroup(ctx, []string{groupID}); err != nil {
				return err
			}
			userIDs, err := g.cache.GetGroupMemberIDs(ctx, groupID)
			if err != nil {
				return err
			}
			cache = cache.DelJoinedGroupID(userIDs...).DelGroupMemberIDs(groupID).DelGroupsMemberNum(groupID).DelGroupMembersHash(groupID)
		}
		cache = cache.DelGroupsInfo(groupID)
		return nil
	}); err != nil {
		return err
	}
	return cache.ExecDel(ctx)
}

func (g *groupDatabase) TakeGroupMember(ctx context.Context, groupID string, userID string) (groupMember *relationTb.GroupMemberModel, err error) {
	return g.cache.GetGroupMemberInfo(ctx, groupID, userID)
}

func (g *groupDatabase) TakeGroupOwner(ctx context.Context, groupID string) (*relationTb.GroupMemberModel, error) {
	return g.groupMemberDB.TakeOwner(ctx, groupID) // todo cache group owner
}

func (g *groupDatabase) FindUserManagedGroupID(ctx context.Context, userID string) (groupIDs []string, err error) {
	return g.groupMemberDB.FindUserManagedGroupID(ctx, userID)
}

func (g *groupDatabase) PageGroupRequest(ctx context.Context, groupIDs []string, pageNumber, showNumber int32) (uint32, []*relationTb.GroupRequestModel, error) {
	return g.groupRequestDB.PageGroup(ctx, groupIDs, pageNumber, showNumber)
}

func (g *groupDatabase) FindGroupMember(ctx context.Context, groupIDs []string, userIDs []string, roleLevels []int32) (totalGroupMembers []*relationTb.GroupMemberModel, err error) {
	if roleLevels == nil {
		for _, groupID := range groupIDs {
			groupMembers, err := g.cache.GetGroupMembersInfo(ctx, groupID, userIDs)
			if err != nil {
				return nil, err
			}
			totalGroupMembers = append(totalGroupMembers, groupMembers...)
		}
		return totalGroupMembers, nil
	}
	return g.groupMemberDB.Find(ctx, groupIDs, userIDs, roleLevels)
}

func (g *groupDatabase) PageGetJoinGroup(ctx context.Context, userID string, pageNumber, showNumber int32) (total uint32, totalGroupMembers []*relationTb.GroupMemberModel, err error) {
	groupIDs, err := g.cache.GetJoinedGroupIDs(ctx, userID)
	if err != nil {
		return 0, nil, err
	}
	for _, groupID := range utils.Paginate(groupIDs, int(pageNumber), int(showNumber)) {
		groupMembers, err := g.cache.GetGroupMembersInfo(ctx, groupID, []string{userID})
		if err != nil {
			return 0, nil, err
		}
		totalGroupMembers = append(totalGroupMembers, groupMembers...)
	}
	return uint32(len(groupIDs)), totalGroupMembers, nil
}

func (g *groupDatabase) PageGetGroupMember(ctx context.Context, groupID string, pageNumber, showNumber int32) (total uint32, totalGroupMembers []*relationTb.GroupMemberModel, err error) {
	groupMemberIDs, err := g.cache.GetGroupMemberIDs(ctx, groupID)
	if err != nil {
		return 0, nil, err
	}
	pageIDs := utils.Paginate(groupMemberIDs, int(pageNumber), int(showNumber))
	if len(pageIDs) == 0 {
		return uint32(len(groupMemberIDs)), nil, nil
	}
	members, err := g.cache.GetGroupMembersInfo(ctx, groupID, pageIDs)
	if err != nil {
		return 0, nil, err
	}
	return uint32(len(groupMemberIDs)), members, nil
}

func (g *groupDatabase) SearchGroupMember(ctx context.Context, keyword string, groupIDs []string, userIDs []string, roleLevels []int32, pageNumber, showNumber int32) (uint32, []*relationTb.GroupMemberModel, error) {
	return g.groupMemberDB.SearchMember(ctx, keyword, groupIDs, userIDs, roleLevels, pageNumber, showNumber)
}

func (g *groupDatabase) HandlerGroupRequest(ctx context.Context, groupID string, userID string, handledMsg string, handleResult int32, member *relationTb.GroupMemberModel) error {
	cache := g.cache.NewCache()
	if err := g.tx.Transaction(func(tx any) error {
		if err := g.groupRequestDB.NewTx(tx).UpdateHandler(ctx, groupID, userID, handledMsg, handleResult); err != nil {
			return err
		}
		if member != nil {
			if err := g.groupMemberDB.NewTx(tx).Create(ctx, []*relationTb.GroupMemberModel{member}); err != nil {
				return err
			}
			cache = cache.DelGroupMembersHash(groupID).DelGroupMemberIDs(groupID).DelGroupsMemberNum(groupID).DelJoinedGroupID(member.UserID)
		}
		return nil
	}); err != nil {
		return err
	}
	return cache.ExecDel(ctx)
}

func (g *groupDatabase) DeleteGroupMember(ctx context.Context, groupID string, userIDs []string) error {
	if err := g.groupMemberDB.Delete(ctx, groupID, userIDs); err != nil {
		return err
	}
	return g.cache.DelGroupMembersHash(groupID).DelGroupMemberIDs(groupID).DelGroupsMemberNum(groupID).DelJoinedGroupID(userIDs...).DelGroupMembersInfo(groupID, userIDs...).ExecDel(ctx)
}

func (g *groupDatabase) MapGroupMemberUserID(ctx context.Context, groupIDs []string) (map[string]*relationTb.GroupSimpleUserID, error) {
	return g.cache.GetGroupMemberHashMap(ctx, groupIDs)
}

func (g *groupDatabase) MapGroupMemberNum(ctx context.Context, groupIDs []string) (m map[string]uint32, err error) {
	m = make(map[string]uint32)
	for _, groupID := range groupIDs {
		num, err := g.cache.GetGroupMemberNum(ctx, groupID)
		if err != nil {
			return nil, err
		}
		m[groupID] = uint32(num)
	}
	return m, nil
}

func (g *groupDatabase) TransferGroupOwner(ctx context.Context, groupID string, oldOwnerUserID, newOwnerUserID string, roleLevel int32) error {
	if err := g.tx.Transaction(func(tx any) error {
		rowsAffected, err := g.groupMemberDB.NewTx(tx).UpdateRoleLevel(ctx, groupID, oldOwnerUserID, roleLevel)
		if err != nil {
			return err
		}
		if rowsAffected != 1 {
			return utils.Wrap(fmt.Errorf("oldOwnerUserID %s rowsAffected = %d", oldOwnerUserID, rowsAffected), "")
		}
		rowsAffected, err = g.groupMemberDB.NewTx(tx).UpdateRoleLevel(ctx, groupID, newOwnerUserID, constant.GroupOwner)
		if err != nil {
			return err
		}
		if rowsAffected != 1 {
			return utils.Wrap(fmt.Errorf("newOwnerUserID %s rowsAffected = %d", newOwnerUserID, rowsAffected), "")
		}
		return nil
	}); err != nil {
		return err
	}
	return g.cache.DelGroupMembersInfo(groupID, oldOwnerUserID, newOwnerUserID).ExecDel(ctx)
}

func (g *groupDatabase) UpdateGroupMember(ctx context.Context, groupID string, userID string, data map[string]any) error {
	if err := g.groupMemberDB.Update(ctx, groupID, userID, data); err != nil {
		return err
	}
	return g.cache.DelGroupMembersInfo(groupID, userID).ExecDel(ctx)
}

func (g *groupDatabase) UpdateGroupMembers(ctx context.Context, data []*relationTb.BatchUpdateGroupMember) error {
	var cache = g.cache.NewCache()
	if err := g.tx.Transaction(func(tx any) error {
		for _, item := range data {
			if err := g.groupMemberDB.NewTx(tx).Update(ctx, item.GroupID, item.UserID, item.Map); err != nil {
				return err
			}
			cache = cache.DelGroupMembersInfo(item.GroupID, item.UserID)
		}
		return nil
	}); err != nil {
		return err
	}
	return cache.ExecDel(ctx)
}

func (g *groupDatabase) CreateGroupRequest(ctx context.Context, requests []*relationTb.GroupRequestModel) error {
	return g.tx.Transaction(func(tx any) error {
		db := g.groupRequestDB.NewTx(tx)
		for _, request := range requests {
			if err := db.Delete(ctx, request.GroupID, request.UserID); err != nil {
				return err
			}
		}
		return db.Create(ctx, requests)
	})
}

func (g *groupDatabase) TakeGroupRequest(ctx context.Context, groupID string, userID string) (*relationTb.GroupRequestModel, error) {
	return g.groupRequestDB.Take(ctx, groupID, userID)
}

func (g *groupDatabase) PageGroupRequestUser(ctx context.Context, userID string, pageNumber, showNumber int32) (uint32, []*relationTb.GroupRequestModel, error) {
	return g.groupRequestDB.Page(ctx, userID, pageNumber, showNumber)
}

func (g *groupDatabase) FindSuperGroup(ctx context.Context, groupIDs []string) (models []*unRelationTb.SuperGroupModel, err error) {
	return g.cache.GetSuperGroupMemberIDs(ctx, groupIDs...)
}

func (g *groupDatabase) FindJoinSuperGroup(ctx context.Context, userID string) ([]string, error) {
	return g.cache.GetJoinedSuperGroupIDs(ctx, userID)
}

func (g *groupDatabase) CreateSuperGroup(ctx context.Context, groupID string, initMemberIDs []string) error {
	if err := g.mongoDB.CreateSuperGroup(ctx, groupID, initMemberIDs); err != nil {
		return err
	}
	return g.cache.DelSuperGroupMemberIDs(groupID).DelJoinedSuperGroupIDs(initMemberIDs...).ExecDel(ctx)
}

func (g *groupDatabase) DeleteSuperGroup(ctx context.Context, groupID string) error {
	cache := g.cache.NewCache()
	if err := g.ctxTx.Transaction(ctx, func(ctx context.Context) error {
		if err := g.mongoDB.DeleteSuperGroup(ctx, groupID); err != nil {
			return err
		}
		models, err := g.cache.GetSuperGroupMemberIDs(ctx, groupID)
		if err != nil {
			return err
		}
		cache = cache.DelSuperGroupMemberIDs(groupID)
		if len(models) > 0 {
			cache = cache.DelJoinedSuperGroupIDs(models[0].MemberIDs...)
		}
		return nil
	}); err != nil {
		return err
	}
	return cache.ExecDel(ctx)
}

func (g *groupDatabase) DeleteSuperGroupMember(ctx context.Context, groupID string, userIDs []string) error {
	if err := g.mongoDB.RemoverUserFromSuperGroup(ctx, groupID, userIDs); err != nil {
		return err
	}
	return g.cache.DelSuperGroupMemberIDs(groupID).DelJoinedSuperGroupIDs(userIDs...).ExecDel(ctx)
}

func (g *groupDatabase) CreateSuperGroupMember(ctx context.Context, groupID string, userIDs []string) error {
	if err := g.mongoDB.AddUserToSuperGroup(ctx, groupID, userIDs); err != nil {
		return err
	}
	return g.cache.DelSuperGroupMemberIDs(groupID).DelJoinedSuperGroupIDs(userIDs...).ExecDel(ctx)
}
