package controller

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db/cache"
	"Open_IM/pkg/common/db/relation"
	relationTb "Open_IM/pkg/common/db/table/relation"
	unRelationTb "Open_IM/pkg/common/db/table/unrelation"
	"Open_IM/pkg/common/db/unrelation"
	"Open_IM/pkg/utils"
	"context"
	"fmt"
	"github.com/dtm-labs/rockscache"
	_ "github.com/dtm-labs/rockscache"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

//type GroupInterface GroupDataBaseInterface

type GroupInterface interface {
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
	// SuperGroup
	FindSuperGroup(ctx context.Context, groupIDs []string) ([]*unRelationTb.SuperGroupModel, error)
	FindJoinSuperGroup(ctx context.Context, userID string) (superGroup *unRelationTb.UserToSuperGroupModel, err error)
	CreateSuperGroup(ctx context.Context, groupID string, initMemberIDList []string) error
	DeleteSuperGroup(ctx context.Context, groupID string) error
	DeleteSuperGroupMember(ctx context.Context, groupID string, userIDs []string) error
	CreateSuperGroupMember(ctx context.Context, groupID string, userIDs []string) error
}

var _ GroupInterface = (*GroupController)(nil)

func NewGroupInterface(database GroupDataBaseInterface) GroupInterface {
	return &GroupController{database: database}
}

type GroupController struct {
	database GroupDataBaseInterface
}

func (g *GroupController) FindGroupMemberUserID(ctx context.Context, groupID string) ([]string, error) {
	return g.database.FindGroupMemberUserID(ctx, groupID)
}

func (g *GroupController) CreateGroup(ctx context.Context, groups []*relationTb.GroupModel, groupMembers []*relationTb.GroupMemberModel) error {
	return g.database.CreateGroup(ctx, groups, groupMembers)
}

func (g *GroupController) TakeGroup(ctx context.Context, groupID string) (group *relationTb.GroupModel, err error) {
	return g.database.TakeGroup(ctx, groupID)
}

func (g *GroupController) FindGroup(ctx context.Context, groupIDs []string) (groups []*relationTb.GroupModel, err error) {
	return g.database.FindGroup(ctx, groupIDs)
}

func (g *GroupController) SearchGroup(ctx context.Context, keyword string, pageNumber, showNumber int32) (uint32, []*relationTb.GroupModel, error) {
	return g.database.SearchGroup(ctx, keyword, pageNumber, showNumber)
}

func (g *GroupController) UpdateGroup(ctx context.Context, groupID string, data map[string]any) error {
	return g.database.UpdateGroup(ctx, groupID, data)
}

func (g *GroupController) DismissGroup(ctx context.Context, groupID string) error {
	return g.database.DismissGroup(ctx, groupID)
}

func (g *GroupController) GetGroupIDsByGroupType(ctx context.Context, groupType int) (groupIDs []string, err error) {
	return g.database.GetGroupIDsByGroupType(ctx, groupType)
}

func (g *GroupController) TakeGroupMember(ctx context.Context, groupID string, userID string) (groupMember *relationTb.GroupMemberModel, err error) {
	return g.database.TakeGroupMember(ctx, groupID, userID)
}

func (g *GroupController) TakeGroupOwner(ctx context.Context, groupID string) (*relationTb.GroupMemberModel, error) {
	return g.database.TakeGroupOwner(ctx, groupID)
}

func (g *GroupController) FindGroupMember(ctx context.Context, groupIDs []string, userIDs []string, roleLevels []int32) ([]*relationTb.GroupMemberModel, error) {
	return g.database.FindGroupMember(ctx, groupIDs, userIDs, roleLevels)
}

func (g *GroupController) PageGroupMember(ctx context.Context, groupIDs []string, userIDs []string, roleLevels []int32, pageNumber, showNumber int32) (uint32, []*relationTb.GroupMemberModel, error) {
	return g.database.PageGroupMember(ctx, groupIDs, userIDs, roleLevels, pageNumber, showNumber)
}

func (g *GroupController) SearchGroupMember(ctx context.Context, keyword string, groupIDs []string, userIDs []string, roleLevels []int32, pageNumber, showNumber int32) (uint32, []*relationTb.GroupMemberModel, error) {
	return g.database.SearchGroupMember(ctx, keyword, groupIDs, userIDs, roleLevels, pageNumber, showNumber)
}

func (g *GroupController) HandlerGroupRequest(ctx context.Context, groupID string, userID string, handledMsg string, handleResult int32, member *relationTb.GroupMemberModel) error {
	return g.database.HandlerGroupRequest(ctx, groupID, userID, handledMsg, handleResult, member)
}

func (g *GroupController) DeleteGroupMember(ctx context.Context, groupID string, userIDs []string) error {
	return g.database.DeleteGroupMember(ctx, groupID, userIDs)
}

func (g *GroupController) MapGroupMemberUserID(ctx context.Context, groupIDs []string) (map[string]*relationTb.GroupSimpleUserID, error) {
	return g.database.MapGroupMemberUserID(ctx, groupIDs)
}

func (g *GroupController) MapGroupMemberNum(ctx context.Context, groupIDs []string) (map[string]uint32, error) {
	return g.database.MapGroupMemberNum(ctx, groupIDs)
}

func (g *GroupController) TransferGroupOwner(ctx context.Context, groupID string, oldOwnerUserID, newOwnerUserID string, roleLevel int32) error {
	return g.database.TransferGroupOwner(ctx, groupID, oldOwnerUserID, newOwnerUserID, roleLevel)
}

func (g *GroupController) UpdateGroupMembers(ctx context.Context, data []*relationTb.BatchUpdateGroupMember) error {
	return g.database.UpdateGroupMembers(ctx, data)
}

func (g *GroupController) UpdateGroupMember(ctx context.Context, groupID string, userID string, data map[string]any) error {
	return g.database.UpdateGroupMember(ctx, groupID, userID, data)
}

func (g *GroupController) CreateGroupRequest(ctx context.Context, requests []*relationTb.GroupRequestModel) error {
	return g.database.CreateGroupRequest(ctx, requests)
}

func (g *GroupController) TakeGroupRequest(ctx context.Context, groupID string, userID string) (*relationTb.GroupRequestModel, error) {
	return g.database.TakeGroupRequest(ctx, groupID, userID)
}

func (g *GroupController) PageGroupRequestUser(ctx context.Context, userID string, pageNumber, showNumber int32) (uint32, []*relationTb.GroupRequestModel, error) {
	return g.database.PageGroupRequestUser(ctx, userID, pageNumber, showNumber)
}

func (g *GroupController) FindSuperGroup(ctx context.Context, groupIDs []string) ([]*unRelationTb.SuperGroupModel, error) {
	return g.database.FindSuperGroup(ctx, groupIDs)
}

func (g *GroupController) FindJoinSuperGroup(ctx context.Context, userID string) (*unRelationTb.UserToSuperGroupModel, error) {
	return g.database.FindJoinSuperGroup(ctx, userID)
}

func (g *GroupController) CreateSuperGroup(ctx context.Context, groupID string, initMemberIDList []string) error {
	return g.database.CreateSuperGroup(ctx, groupID, initMemberIDList)
}

func (g *GroupController) DeleteSuperGroup(ctx context.Context, groupID string) error {
	return g.database.DeleteSuperGroup(ctx, groupID)
}

func (g *GroupController) DeleteSuperGroupMember(ctx context.Context, groupID string, userIDs []string) error {
	return g.database.DeleteSuperGroupMember(ctx, groupID, userIDs)
}

func (g *GroupController) CreateSuperGroupMember(ctx context.Context, groupID string, userIDs []string) error {
	return g.database.CreateSuperGroupMember(ctx, groupID, userIDs)
}

type Group interface {
	CreateGroup(ctx context.Context, groups []*relationTb.GroupModel, groupMembers []*relationTb.GroupMemberModel) error
	TakeGroup(ctx context.Context, groupID string) (group *relationTb.GroupModel, err error)
	FindGroup(ctx context.Context, groupIDs []string) (groups []*relationTb.GroupModel, err error)
	SearchGroup(ctx context.Context, keyword string, pageNumber, showNumber int32) (uint32, []*relationTb.GroupModel, error)
	UpdateGroup(ctx context.Context, groupID string, data map[string]any) error
	DismissGroup(ctx context.Context, groupID string) error // 解散群，并删除群成员
	GetGroupIDsByGroupType(ctx context.Context, groupType int) (groupIDs []string, err error)
}

type GroupMember interface {
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
}

type GroupRequest interface {
	CreateGroupRequest(ctx context.Context, requests []*relationTb.GroupRequestModel) error
	TakeGroupRequest(ctx context.Context, groupID string, userID string) (*relationTb.GroupRequestModel, error)
	PageGroupRequestUser(ctx context.Context, userID string, pageNumber, showNumber int32) (uint32, []*relationTb.GroupRequestModel, error)
}

type SuperGroup interface {
	FindSuperGroup(ctx context.Context, groupIDs []string) ([]*unrelationTb.SuperGroupModel, error)
	FindJoinSuperGroup(ctx context.Context, userID string) (*unrelationTb.UserToSuperGroupModel, error)
	CreateSuperGroup(ctx context.Context, groupID string, initMemberIDList []string) error
	DeleteSuperGroup(ctx context.Context, groupID string) error
	DeleteSuperGroupMember(ctx context.Context, groupID string, userIDs []string) error
	CreateSuperGroupMember(ctx context.Context, groupID string, userIDs []string) error
}

type GroupDataBase1 interface {
	Group
	GroupMember
	GroupRequest
	SuperGroup
}

type GroupDataBaseInterface interface {
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
	// SuperGroup
	FindSuperGroup(ctx context.Context, groupIDs []string) ([]*unRelationTb.SuperGroupModel, error)
	FindJoinSuperGroup(ctx context.Context, userID string) (*unRelationTb.UserToSuperGroupModel, error)
	CreateSuperGroup(ctx context.Context, groupID string, initMemberIDList []string) error
	DeleteSuperGroup(ctx context.Context, groupID string) error
	DeleteSuperGroupMember(ctx context.Context, groupID string, userIDs []string) error
	CreateSuperGroupMember(ctx context.Context, groupID string, userIDs []string) error
}

func NewGroupDatabase(db *gorm.DB, rdb redis.UniversalClient, mgoClient *mongo.Client) GroupDataBaseInterface {
	groupDB := relation.NewGroupDB(db)
	groupMemberDB := relation.NewGroupMemberDB(db)
	groupRequestDB := relation.NewGroupRequest(db)
	newDB := *db
	SuperGroupMongoDriver := unrelation.NewSuperGroupMongoDriver(mgoClient)
	database := &GroupDataBase{
		groupDB:        groupDB,
		groupMemberDB:  groupMemberDB,
		groupRequestDB: groupRequestDB,
		db:             &newDB,
		cache: cache.NewGroupCacheRedis(rdb, groupDB, groupMemberDB, groupRequestDB, SuperGroupMongoDriver, rockscache.Options{
			RandomExpireAdjustment: 0.2,
			DisableCacheRead:       false,
			DisableCacheDelete:     false,
			StrongConsistency:      true,
		}),
		mongoDB: SuperGroupMongoDriver,
	}
	return database
}

var _ GroupDataBaseInterface = (*GroupDataBase)(nil)

type GroupDataBase struct {
	groupDB        relationTb.GroupModelInterface
	groupMemberDB  relationTb.GroupMemberModelInterface
	groupRequestDB relationTb.GroupRequestModelInterface
	db             *gorm.DB

	//cache   cache.GroupCache
	cache   *cache.GroupCacheRedis
	mongoDB *unrelation.SuperGroupMongoDriver
}

func (g *GroupDataBase) delGroupMemberCache(ctx context.Context, groupID string, userIDs []string) error {
	for _, userID := range userIDs {
		if err := g.cache.DelJoinedGroupID(ctx, userID); err != nil {
			return err
		}
		if err := g.cache.DelJoinedSuperGroupIDs(ctx, userID); err != nil {
			return err
		}
	}
	if err := g.cache.DelGroupMemberIDs(ctx, groupID); err != nil {
		return err
	}
	if err := g.cache.DelGroupMemberNum(ctx, groupID); err != nil {
		return err
	}
	if err := g.cache.DelGroupMembersHash(ctx, groupID); err != nil {
		return err
	}
	return nil
}

func (g *GroupDataBase) FindGroupMemberUserID(ctx context.Context, groupID string) ([]string, error) {
	return g.cache.GetGroupMemberIDs(ctx, groupID)
}

func (g *GroupDataBase) CreateGroup(ctx context.Context, groups []*relationTb.GroupModel, groupMembers []*relationTb.GroupMemberModel) error {
	return g.db.Transaction(func(tx *gorm.DB) error {
		if len(groups) > 0 {
			if err := g.groupDB.Create(ctx, groups, tx); err != nil {
				return err
			}
		}
		if len(groupMembers) > 0 {
			if err := g.groupMemberDB.Create(ctx, groupMembers, tx); err != nil {
				return err
			}
			//if err := g.cache.DelJoinedGroupIDs(ctx, utils.Slice(groupMembers, func(e *relationTb.GroupMemberModel) string {
			//	return e.UserID
			//})); err != nil {
			//	return err
			//}
		}
		return nil
	})
}

func (g *GroupDataBase) TakeGroup(ctx context.Context, groupID string) (group *relationTb.GroupModel, err error) {
	//return g.cache.GetGroupInfo(ctx, groupID)
	return cache.GetCache(ctx, g.rcClient, g.getGroupInfoKey(groupID), g.expireTime, func(ctx context.Context) (*relationTb.GroupModel, error) {
		return g.group.Take(ctx, groupID)
	})
}

func (g *GroupDataBase) FindGroup(ctx context.Context, groupIDs []string) (groups []*relationTb.GroupModel, err error) {
	return g.cache.GetGroupsInfo(ctx, groupIDs)
}

func (g *GroupDataBase) SearchGroup(ctx context.Context, keyword string, pageNumber, showNumber int32) (uint32, []*relationTb.GroupModel, error) {
	return g.groupDB.Search(ctx, keyword, pageNumber, showNumber)
}

func (g *GroupDataBase) UpdateGroup(ctx context.Context, groupID string, data map[string]any) error {
	return g.db.Transaction(func(tx *gorm.DB) error {
		if err := g.groupDB.UpdateMap(ctx, groupID, data, tx); err != nil {
			return err
		}
		if err := g.cache.DelGroupInfo(ctx, groupID); err != nil {
			return err
		}
		return nil
	})
}

func (g *GroupDataBase) DismissGroup(ctx context.Context, groupID string) error {
	return g.db.Transaction(func(tx *gorm.DB) error {
		if err := g.groupDB.UpdateStatus(ctx, groupID, constant.GroupStatusDismissed, tx); err != nil {
			return err
		}
		if err := g.groupMemberDB.DeleteGroup(ctx, []string{groupID}, tx); err != nil {
			return err
		}
		userIDs, err := g.cache.GetGroupMemberIDs(ctx, groupID)
		if err != nil {
			return err
		}
		if err := g.delGroupMemberCache(ctx, groupID, userIDs); err != nil {
			return err
		}
		return nil
	})
}

func (g *GroupDataBase) TakeGroupMember(ctx context.Context, groupID string, userID string) (groupMember *relationTb.GroupMemberModel, err error) {
	return g.cache.GetGroupMemberInfo(ctx, groupID, userID)
}

func (g *GroupDataBase) TakeGroupOwner(ctx context.Context, groupID string) (*relationTb.GroupMemberModel, error) {
	return g.groupMemberDB.TakeOwner(ctx, groupID) // todo cache group owner
}

func (g *GroupDataBase) FindGroupMember(ctx context.Context, groupIDs []string, userIDs []string, roleLevels []int32) ([]*relationTb.GroupMemberModel, error) {
	return g.groupMemberDB.Find(ctx, groupIDs, userIDs, roleLevels) // todo cache group find
}

func (g *GroupDataBase) PageGroupMember(ctx context.Context, groupIDs []string, userIDs []string, roleLevels []int32, pageNumber, showNumber int32) (uint32, []*relationTb.GroupMemberModel, error) {
	return g.groupMemberDB.SearchMember(ctx, "", groupIDs, userIDs, roleLevels, pageNumber, showNumber)
}

func (g *GroupDataBase) SearchGroupMember(ctx context.Context, keyword string, groupIDs []string, userIDs []string, roleLevels []int32, pageNumber, showNumber int32) (uint32, []*relationTb.GroupMemberModel, error) {
	return g.groupMemberDB.SearchMember(ctx, keyword, groupIDs, userIDs, roleLevels, pageNumber, showNumber)
}

func (g *GroupDataBase) HandlerGroupRequest(ctx context.Context, groupID string, userID string, handledMsg string, handleResult int32, member *relationTb.GroupMemberModel) error {
	return g.db.Transaction(func(tx *gorm.DB) error {
		if err := g.groupRequestDB.UpdateHandler(ctx, groupID, userID, handledMsg, handleResult, tx); err != nil {
			return err
		}
		if member != nil {
			if err := g.groupMemberDB.Create(ctx, []*relationTb.GroupMemberModel{member}, tx); err != nil {
				return err
			}
			if err := g.delGroupMemberCache(ctx, groupID, []string{userID}); err != nil {
				return err
			}
		}
		return nil
	})
}

func (g *GroupDataBase) DeleteGroupMember(ctx context.Context, groupID string, userIDs []string) error {
	return g.db.Transaction(func(tx *gorm.DB) error {
		if err := g.groupMemberDB.Delete(ctx, groupID, userIDs, tx); err != nil {
			return err
		}
		if err := g.delGroupMemberCache(ctx, groupID, userIDs); err != nil {
			return err
		}
		return nil
	})
}

func (g *GroupDataBase) MapGroupMemberUserID(ctx context.Context, groupIDs []string) (map[string]*relationTb.GroupSimpleUserID, error) {
	return g.cache.GetGroupMemberHash1(ctx, groupIDs)
}

func (g *GroupDataBase) MapGroupMemberNum(ctx context.Context, groupIDs []string) (map[string]uint32, error) {
	return g.groupMemberDB.MapGroupMemberNum(ctx, groupIDs)
}

func (g *GroupDataBase) TransferGroupOwner(ctx context.Context, groupID string, oldOwnerUserID, newOwnerUserID string, roleLevel int32) error {
	return g.db.Transaction(func(tx *gorm.DB) error {
		rowsAffected, err := g.groupMemberDB.UpdateRoleLevel(ctx, groupID, oldOwnerUserID, roleLevel, tx)
		if err != nil {
			return err
		}
		if rowsAffected != 1 {
			return utils.Wrap(fmt.Errorf("oldOwnerUserID %s rowsAffected = %d", oldOwnerUserID, rowsAffected), "")
		}
		rowsAffected, err = g.groupMemberDB.UpdateRoleLevel(ctx, groupID, newOwnerUserID, constant.GroupOwner, tx)
		if err != nil {
			return err
		}
		if rowsAffected != 1 {
			return utils.Wrap(fmt.Errorf("newOwnerUserID %s rowsAffected = %d", newOwnerUserID, rowsAffected), "")
		}
		if err := g.delGroupMemberCache(ctx, groupID, []string{oldOwnerUserID, newOwnerUserID}); err != nil {
			return err
		}
		return nil
	})
}

func (g *GroupDataBase) UpdateGroupMember(ctx context.Context, groupID string, userID string, data map[string]any) error {
	return g.db.Transaction(func(tx *gorm.DB) error {
		if err := g.groupMemberDB.Update(ctx, groupID, userID, data, tx); err != nil {
			return err
		}
		if err := g.cache.DelGroupMemberInfo(ctx, groupID, userID); err != nil {
			return err
		}
		return nil
	})
}

func (g *GroupDataBase) UpdateGroupMembers(ctx context.Context, data []*relationTb.BatchUpdateGroupMember) error {
	return g.db.Transaction(func(tx *gorm.DB) error {
		for _, item := range data {
			if err := g.groupMemberDB.Update(ctx, item.GroupID, item.UserID, item.Map, tx); err != nil {
				return err
			}
			if err := g.cache.DelGroupMemberInfo(ctx, item.GroupID, item.UserID); err != nil {
				return err
			}
		}
		return nil
	})
}

func (g *GroupDataBase) CreateGroupRequest(ctx context.Context, requests []*relationTb.GroupRequestModel) error {
	return g.groupRequestDB.Create(ctx, requests)
}

func (g *GroupDataBase) TakeGroupRequest(ctx context.Context, groupID string, userID string) (*relationTb.GroupRequestModel, error) {
	return g.groupRequestDB.Take(ctx, groupID, userID)
}

func (g *GroupDataBase) PageGroupRequestUser(ctx context.Context, userID string, pageNumber, showNumber int32) (uint32, []*relationTb.GroupRequestModel, error) {
	return g.groupRequestDB.Page(ctx, userID, pageNumber, showNumber)
}

func (g *GroupDataBase) FindSuperGroup(ctx context.Context, groupIDs []string) ([]*table.SuperGroupModel, error) {
	return g.mongoDB.FindSuperGroup(ctx, groupIDs)
}

func (g *GroupDataBase) FindJoinSuperGroup(ctx context.Context, userID string) (*table.UserToSuperGroupModel, error) {
	return g.mongoDB.GetSuperGroupByUserID(ctx, userID)
}

func (g *GroupDataBase) CreateSuperGroup(ctx context.Context, groupID string, initMemberIDList []string) error {
	return unrelation.MongoTransaction(ctx, g.mongoDB.MgoClient, func(tx mongo.SessionContext) error {
		return g.mongoDB.CreateSuperGroup(ctx, groupID, initMemberIDList, tx)
	})
}

func (g *GroupDataBase) DeleteSuperGroup(ctx context.Context, groupID string) error {
	return unrelation.MongoTransaction(ctx, g.mongoDB.MgoClient, func(tx mongo.SessionContext) error {
		return g.mongoDB.DeleteSuperGroup(ctx, groupID, tx)
	})
}

func (g *GroupDataBase) DeleteSuperGroupMember(ctx context.Context, groupID string, userIDs []string) error {
	return unrelation.MongoTransaction(ctx, g.mongoDB.MgoClient, func(tx mongo.SessionContext) error {
		return g.mongoDB.RemoverUserFromSuperGroup(ctx, groupID, userIDs, tx)
	})
}

func (g *GroupDataBase) CreateSuperGroupMember(ctx context.Context, groupID string, userIDs []string) error {
	return unrelation.MongoTransaction(ctx, g.mongoDB.MgoClient, func(tx mongo.SessionContext) error {
		return g.mongoDB.AddUserToSuperGroup(ctx, groupID, userIDs, tx)
	})
}
