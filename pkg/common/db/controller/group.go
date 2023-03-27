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
	_ "github.com/dtm-labs/rockscache"
	"github.com/go-redis/redis/v8"
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
	DismissGroup(ctx context.Context, groupID string) error // 解散群，并删除群成员
	GetGroupIDsByGroupType(ctx context.Context, groupType int) (groupIDs []string, err error)
	// GroupMember
	TakeGroupMember(ctx context.Context, groupID string, userID string) (groupMember *relationTb.GroupMemberModel, err error)
	TakeGroupOwner(ctx context.Context, groupID string) (*relationTb.GroupMemberModel, error)
	FindGroupMember(ctx context.Context, groupIDs []string, userIDs []string, roleLevels []int32) ([]*relationTb.GroupMemberModel, error)
	FindGroupMemberUserID(ctx context.Context, groupID string) ([]string, error)
	PageGroupMember(ctx context.Context, groupIDs []string, userIDs []string, roleLevels []int32, pageNumber, showNumber int32) (uint32, []*relationTb.GroupMemberModel, error)
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
	FindJoinSuperGroup(ctx context.Context, userID string) (*unRelationTb.UserToSuperGroupModel, error)
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

func (g *groupDatabase) CreateGroup(ctx context.Context, groups []*relationTb.GroupModel, groupMembers []*relationTb.GroupMemberModel) error {
	return g.tx.Transaction(func(tx any) error {
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
		m := make(map[string]struct{})
		var cache = g.cache.NewCache()
		for _, groupMember := range groupMembers {
			if _, ok := m[groupMember.GroupID]; !ok {
				m[groupMember.GroupID] = struct{}{}
				cache = cache.DelGroupMemberIDs(groupMember.GroupID).DelGroupMembersHash(groupMember.GroupID).DelJoinedGroupID(groupMember.UserID).DelGroupsMemberNum(groupMember.GroupID)
			}
		}
		return g.cache.ExecDel(ctx)
	})
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
	return g.tx.Transaction(func(tx any) error {
		if err := g.groupDB.NewTx(tx).UpdateMap(ctx, groupID, data); err != nil {
			return err
		}
		return g.cache.DelGroupsInfo(groupID).ExecDel(ctx)
	})
}

func (g *groupDatabase) DismissGroup(ctx context.Context, groupID string) error {
	return g.tx.Transaction(func(tx any) error {
		if err := g.groupDB.NewTx(tx).UpdateStatus(ctx, groupID, constant.GroupStatusDismissed); err != nil {
			return err
		}
		if err := g.groupMemberDB.NewTx(tx).DeleteGroup(ctx, []string{groupID}); err != nil {
			return err
		}
		userIDs, err := g.cache.GetGroupMemberIDs(ctx, groupID)
		if err != nil {
			return err
		}
		return g.cache.DelJoinedGroupID(userIDs...).DelGroupsInfo(groupID).DelGroupMemberIDs(groupID).DelGroupsMemberNum(groupID).DelGroupMembersHash(groupID).ExecDel(ctx)
	})
}

func (g *groupDatabase) TakeGroupMember(ctx context.Context, groupID string, userID string) (groupMember *relationTb.GroupMemberModel, err error) {
	return g.cache.GetGroupMemberInfo(ctx, groupID, userID)
}

func (g *groupDatabase) TakeGroupOwner(ctx context.Context, groupID string) (*relationTb.GroupMemberModel, error) {
	return g.groupMemberDB.TakeOwner(ctx, groupID) // todo cache group owner
}

func (g *groupDatabase) FindGroupMember(ctx context.Context, groupIDs []string, userIDs []string, roleLevels []int32) ([]*relationTb.GroupMemberModel, error) {
	return g.cache.GetGroupMembersInfo(ctx, groupIDs[0], userIDs, roleLevels) // todo cache group find
}

func (g *groupDatabase) PageGroupMember(ctx context.Context, groupIDs []string, userIDs []string, roleLevels []int32, pageNumber, showNumber int32) (uint32, []*relationTb.GroupMemberModel, error) {
	return g.groupMemberDB.SearchMember(ctx, "", groupIDs, userIDs, roleLevels, pageNumber, showNumber)
}

func (g *groupDatabase) SearchGroupMember(ctx context.Context, keyword string, groupIDs []string, userIDs []string, roleLevels []int32, pageNumber, showNumber int32) (uint32, []*relationTb.GroupMemberModel, error) {
	return g.groupMemberDB.SearchMember(ctx, keyword, groupIDs, userIDs, roleLevels, pageNumber, showNumber)
}

func (g *groupDatabase) HandlerGroupRequest(ctx context.Context, groupID string, userID string, handledMsg string, handleResult int32, member *relationTb.GroupMemberModel) error {
	return g.tx.Transaction(func(tx any) error {
		if err := g.groupRequestDB.NewTx(tx).UpdateHandler(ctx, groupID, userID, handledMsg, handleResult); err != nil {
			return err
		}
		if member != nil {
			if err := g.groupMemberDB.NewTx(tx).Create(ctx, []*relationTb.GroupMemberModel{member}); err != nil {
				return err
			}
			return g.cache.DelGroupMembersHash(groupID).DelGroupMemberIDs(groupID).DelGroupsMemberNum(groupID).DelJoinedGroupID(member.UserID).ExecDel(ctx)
		}
		return nil
	})
}

func (g *groupDatabase) DeleteGroupMember(ctx context.Context, groupID string, userIDs []string) error {
	return g.tx.Transaction(func(tx any) error {
		if err := g.groupMemberDB.NewTx(tx).Delete(ctx, groupID, userIDs); err != nil {
			return err
		}
		return g.cache.DelGroupMembersHash(groupID).DelGroupMemberIDs(groupID).DelGroupsMemberNum(groupID).DelJoinedGroupID(userIDs...).ExecDel(ctx)
	})
}

func (g *groupDatabase) MapGroupMemberUserID(ctx context.Context, groupIDs []string) (map[string]*relationTb.GroupSimpleUserID, error) {
	return g.cache.GetGroupMemberHashMap(ctx, groupIDs)
}

func (g *groupDatabase) MapGroupMemberNum(ctx context.Context, groupIDs []string) (map[string]uint32, error) {
	return g.groupMemberDB.MapGroupMemberNum(ctx, groupIDs)
}

func (g *groupDatabase) TransferGroupOwner(ctx context.Context, groupID string, oldOwnerUserID, newOwnerUserID string, roleLevel int32) error {
	return g.tx.Transaction(func(tx any) error {
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
		return g.cache.DelGroupMembersInfo(groupID, oldOwnerUserID, newOwnerUserID).ExecDel(ctx)
	})
}

func (g *groupDatabase) UpdateGroupMember(ctx context.Context, groupID string, userID string, data map[string]any) error {
	return g.tx.Transaction(func(tx any) error {
		if err := g.groupMemberDB.NewTx(tx).Update(ctx, groupID, userID, data); err != nil {
			return err
		}
		return g.cache.DelGroupMembersInfo(groupID, userID).ExecDel(ctx)
	})
}

func (g *groupDatabase) UpdateGroupMembers(ctx context.Context, data []*relationTb.BatchUpdateGroupMember) error {
	return g.tx.Transaction(func(tx any) error {
		var cache = g.cache.NewCache()
		for _, item := range data {
			if err := g.groupMemberDB.NewTx(tx).Update(ctx, item.GroupID, item.UserID, item.Map); err != nil {
				return err
			}
			cache = cache.DelGroupMembersInfo(item.GroupID, item.UserID)
		}
		return cache.ExecDel(ctx)
	})
}

func (g *groupDatabase) CreateGroupRequest(ctx context.Context, requests []*relationTb.GroupRequestModel) error {
	return g.groupRequestDB.Create(ctx, requests)
}

func (g *groupDatabase) TakeGroupRequest(ctx context.Context, groupID string, userID string) (*relationTb.GroupRequestModel, error) {
	return g.groupRequestDB.Take(ctx, groupID, userID)
}

func (g *groupDatabase) PageGroupRequestUser(ctx context.Context, userID string, pageNumber, showNumber int32) (uint32, []*relationTb.GroupRequestModel, error) {
	return g.groupRequestDB.Page(ctx, userID, pageNumber, showNumber)
}

func (g *groupDatabase) FindSuperGroup(ctx context.Context, groupIDs []string) ([]*unRelationTb.SuperGroupModel, error) {
	return g.mongoDB.FindSuperGroup(ctx, groupIDs)
}

func (g *groupDatabase) FindJoinSuperGroup(ctx context.Context, userID string) (*unRelationTb.UserToSuperGroupModel, error) {
	return g.mongoDB.GetSuperGroupByUserID(ctx, userID)
}

func (g *groupDatabase) CreateSuperGroup(ctx context.Context, groupID string, initMemberIDList []string) error {
	return g.ctxTx.Transaction(ctx, func(ctx context.Context) error {
		return g.mongoDB.CreateSuperGroup(ctx, groupID, initMemberIDList)
	})
}

func (g *groupDatabase) DeleteSuperGroup(ctx context.Context, groupID string) error {
	return g.ctxTx.Transaction(ctx, func(ctx context.Context) error {
		return g.mongoDB.DeleteSuperGroup(ctx, groupID)
	})
}

func (g *groupDatabase) DeleteSuperGroupMember(ctx context.Context, groupID string, userIDs []string) error {
	return g.ctxTx.Transaction(ctx, func(ctx context.Context) error {
		return g.mongoDB.RemoverUserFromSuperGroup(ctx, groupID, userIDs)
	})
}

func (g *groupDatabase) CreateSuperGroupMember(ctx context.Context, groupID string, userIDs []string) error {
	return g.ctxTx.Transaction(ctx, func(ctx context.Context) error {
		return g.mongoDB.AddUserToSuperGroup(ctx, groupID, userIDs)
	})
}
