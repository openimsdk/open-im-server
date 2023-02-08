package controller

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db/cache"
	"Open_IM/pkg/common/db/relation"
	relation2 "Open_IM/pkg/common/db/table/relation"
	unrelation2 "Open_IM/pkg/common/db/table/unrelation"
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
	CreateGroup(ctx context.Context, groups []*relation2.GroupModel, groupMembers []*relation2.GroupMemberModel) error
	TakeGroup(ctx context.Context, groupID string) (group *relation2.GroupModel, err error)
	FindGroup(ctx context.Context, groupIDs []string) (groups []*relation2.GroupModel, err error)
	SearchGroup(ctx context.Context, keyword string, pageNumber, showNumber int32) (int32, []*relation2.GroupModel, error)
	UpdateGroup(ctx context.Context, groupID string, data map[string]any) error
	DismissGroup(ctx context.Context, groupID string) error // 解散群，并删除群成员
	// GroupMember
	TakeGroupMember(ctx context.Context, groupID string, userID string) (groupMember *relation2.GroupMemberModel, err error)
	TakeGroupOwner(ctx context.Context, groupID string) (*relation2.GroupMemberModel, error)
	FindGroupMember(ctx context.Context, groupIDs []string, userIDs []string, roleLevels []int32) ([]*relation2.GroupMemberModel, error)
	PageGroupMember(ctx context.Context, groupIDs []string, userIDs []string, roleLevels []int32, pageNumber, showNumber int32) (int32, []*relation2.GroupMemberModel, error)
	SearchGroupMember(ctx context.Context, keyword string, groupIDs []string, userIDs []string, roleLevels []int32, pageNumber, showNumber int32) (int32, []*relation2.GroupMemberModel, error)
	HandlerGroupRequest(ctx context.Context, groupID string, userID string, handledMsg string, handleResult int32, member *relation2.GroupMemberModel) error
	DeleteGroupMember(ctx context.Context, groupID string, userIDs []string) error
	MapGroupMemberUserID(ctx context.Context, groupIDs []string) (map[string][]string, error)
	MapGroupMemberNum(ctx context.Context, groupIDs []string) (map[string]uint32, error)
	TransferGroupOwner(ctx context.Context, groupID string, oldOwnerUserID, newOwnerUserID string, roleLevel int32) error // 转让群
	UpdateGroupMember(ctx context.Context, groupID, userID string, data map[string]any) error
	// GroupRequest
	CreateGroupRequest(ctx context.Context, requests []*relation2.GroupRequestModel) error
	TakeGroupRequest(ctx context.Context, groupID string, userID string) (*relation2.GroupRequestModel, error)
	PageGroupRequestUser(ctx context.Context, userID string, pageNumber, showNumber int32) (int32, []*relation2.GroupRequestModel, error)
	// SuperGroup
	FindSuperGroup(ctx context.Context, groupIDs []string) ([]*unrelation2.SuperGroupModel, error)
	FindJoinSuperGroup(ctx context.Context, userID string) (superGroup *unrelation2.UserToSuperGroupModel, err error)
	CreateSuperGroup(ctx context.Context, groupID string, initMemberIDList []string) error
	DeleteSuperGroup(ctx context.Context, groupID string) error
	DeleteSuperGroupMember(ctx context.Context, groupID string, userIDs []string) error
	CreateSuperGroupMember(ctx context.Context, groupID string, userIDs []string) error
}

var _ GroupInterface = (*GroupController)(nil)

type GroupController struct {
	database GroupDataBaseInterface
}

func (g *GroupController) CreateGroup(ctx context.Context, groups []*relation2.GroupModel, groupMembers []*relation2.GroupMemberModel) error {
	return g.database.CreateGroup(ctx, groups, groupMembers)
}

func (g *GroupController) TakeGroup(ctx context.Context, groupID string) (group *relation2.GroupModel, err error) {
	return g.TakeGroup(ctx, groupID)
}

func (g *GroupController) FindGroup(ctx context.Context, groupIDs []string) (groups []*relation2.GroupModel, err error) {
	return g.database.FindGroup(ctx, groupIDs)
}

func (g *GroupController) SearchGroup(ctx context.Context, keyword string, pageNumber, showNumber int32) (int32, []*relation2.GroupModel, error) {
	return g.database.SearchGroup(ctx, keyword, pageNumber, showNumber)
}

func (g *GroupController) UpdateGroup(ctx context.Context, groupID string, data map[string]any) error {
	return g.database.UpdateGroup(ctx, groupID, data)
}

func (g *GroupController) DismissGroup(ctx context.Context, groupID string) error {
	return g.database.DismissGroup(ctx, groupID)
}

func (g *GroupController) TakeGroupMember(ctx context.Context, groupID string, userID string) (groupMember *relation2.GroupMemberModel, err error) {
	return g.database.TakeGroupMember(ctx, groupID, userID)
}

func (g *GroupController) TakeGroupOwner(ctx context.Context, groupID string) (*relation2.GroupMemberModel, error) {
	return g.database.TakeGroupOwner(ctx, groupID)
}

func (g *GroupController) FindGroupMember(ctx context.Context, groupIDs []string, userIDs []string, roleLevels []int32) ([]*relation2.GroupMemberModel, error) {
	return g.database.FindGroupMember(ctx, groupIDs, userIDs, roleLevels)
}

func (g *GroupController) PageGroupMember(ctx context.Context, groupIDs []string, userIDs []string, roleLevels []int32, pageNumber, showNumber int32) (int32, []*relation2.GroupMemberModel, error) {
	return g.database.PageGroupMember(ctx, groupIDs, userIDs, roleLevels, pageNumber, showNumber)
}

func (g *GroupController) SearchGroupMember(ctx context.Context, keyword string, groupIDs []string, userIDs []string, roleLevels []int32, pageNumber, showNumber int32) (int32, []*relation2.GroupMemberModel, error) {
	return g.database.SearchGroupMember(ctx, keyword, groupIDs, userIDs, roleLevels, pageNumber, showNumber)
}

func (g *GroupController) HandlerGroupRequest(ctx context.Context, groupID string, userID string, handledMsg string, handleResult int32, member *relation2.GroupMemberModel) error {
	return g.database.HandlerGroupRequest(ctx, groupID, userID, handledMsg, handleResult, member)
}

func (g *GroupController) DeleteGroupMember(ctx context.Context, groupID string, userIDs []string) error {
	return g.database.DeleteGroupMember(ctx, groupID, userIDs)
}

func (g *GroupController) MapGroupMemberUserID(ctx context.Context, groupIDs []string) (map[string][]string, error) {
	return g.database.MapGroupMemberUserID(ctx, groupIDs)
}

func (g *GroupController) MapGroupMemberNum(ctx context.Context, groupIDs []string) (map[string]uint32, error) {
	return g.database.MapGroupMemberNum(ctx, groupIDs)
}

func (g *GroupController) TransferGroupOwner(ctx context.Context, groupID string, oldOwnerUserID, newOwnerUserID string, roleLevel int32) error {
	return g.database.TransferGroupOwner(ctx, groupID, oldOwnerUserID, newOwnerUserID, roleLevel)
}

func (g *GroupController) UpdateGroupMember(ctx context.Context, groupID, userID string, data map[string]any) error {
	return g.database.UpdateGroupMember(ctx, groupID, userID, data)
}

func (g *GroupController) CreateGroupRequest(ctx context.Context, requests []*relation2.GroupRequestModel) error {
	return g.database.CreateGroupRequest(ctx, requests)
}

func (g *GroupController) TakeGroupRequest(ctx context.Context, groupID string, userID string) (*relation2.GroupRequestModel, error) {
	return g.database.TakeGroupRequest(ctx, groupID, userID)
}

func (g *GroupController) PageGroupRequestUser(ctx context.Context, userID string, pageNumber, showNumber int32) (int32, []*relation2.GroupRequestModel, error) {
	return g.database.PageGroupRequestUser(ctx, userID, pageNumber, showNumber)
}

//	func (g *GroupController) TakeSuperGroup(ctx context.Context, groupID string) (superGroup *unrelation2.SuperGroupModel, err error) {
//		return g.database.TakeSuperGroup(ctx, groupID)
//	}
func (g *GroupController) FindSuperGroup(ctx context.Context, groupIDs []string) ([]*unrelation2.SuperGroupModel, error) {
	return g.database.FindSuperGroup(ctx, groupIDs)
}

func (g *GroupController) FindJoinSuperGroup(ctx context.Context, userID string) (*unrelation2.UserToSuperGroupModel, error) {
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

type GroupDataBaseInterface interface {
	CreateGroup(ctx context.Context, groups []*relation2.GroupModel, groupMembers []*relation2.GroupMemberModel) error
	TakeGroup(ctx context.Context, groupID string) (group *relation2.GroupModel, err error)
	FindGroup(ctx context.Context, groupIDs []string) (groups []*relation2.GroupModel, err error)
	SearchGroup(ctx context.Context, keyword string, pageNumber, showNumber int32) (int32, []*relation2.GroupModel, error)
	UpdateGroup(ctx context.Context, groupID string, data map[string]any) error
	DismissGroup(ctx context.Context, groupID string) error // 解散群，并删除群成员
	// GroupMember
	TakeGroupMember(ctx context.Context, groupID string, userID string) (groupMember *relation2.GroupMemberModel, err error)
	TakeGroupOwner(ctx context.Context, groupID string) (*relation2.GroupMemberModel, error)
	FindGroupMember(ctx context.Context, groupIDs []string, userIDs []string, roleLevels []int32) ([]*relation2.GroupMemberModel, error)
	PageGroupMember(ctx context.Context, groupIDs []string, userIDs []string, roleLevels []int32, pageNumber, showNumber int32) (int32, []*relation2.GroupMemberModel, error)
	SearchGroupMember(ctx context.Context, keyword string, groupIDs []string, userIDs []string, roleLevels []int32, pageNumber, showNumber int32) (int32, []*relation2.GroupMemberModel, error)
	HandlerGroupRequest(ctx context.Context, groupID string, userID string, handledMsg string, handleResult int32, member *relation2.GroupMemberModel) error
	DeleteGroupMember(ctx context.Context, groupID string, userIDs []string) error
	MapGroupMemberUserID(ctx context.Context, groupIDs []string) (map[string][]string, error)
	MapGroupMemberNum(ctx context.Context, groupIDs []string) (map[string]uint32, error)
	TransferGroupOwner(ctx context.Context, groupID string, oldOwnerUserID, newOwnerUserID string, roleLevel int32) error // 转让群
	UpdateGroupMember(ctx context.Context, groupID, userID string, data map[string]any) error
	// GroupRequest
	CreateGroupRequest(ctx context.Context, requests []*relation2.GroupRequestModel) error
	TakeGroupRequest(ctx context.Context, groupID string, userID string) (*relation2.GroupRequestModel, error)
	PageGroupRequestUser(ctx context.Context, userID string, pageNumber, showNumber int32) (int32, []*relation2.GroupRequestModel, error)
	// SuperGroup
	FindSuperGroup(ctx context.Context, groupIDs []string) ([]*unrelation2.SuperGroupModel, error)
	FindJoinSuperGroup(ctx context.Context, userID string) (*unrelation2.UserToSuperGroupModel, error)
	CreateSuperGroup(ctx context.Context, groupID string, initMemberIDList []string) error
	DeleteSuperGroup(ctx context.Context, groupID string) error
	DeleteSuperGroupMember(ctx context.Context, groupID string, userIDs []string) error
	CreateSuperGroupMember(ctx context.Context, groupID string, userIDs []string) error
}

func newGroupDatabase(db *gorm.DB, rdb redis.UniversalClient, mgoClient *mongo.Client) GroupDataBaseInterface {
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
		cache: cache.NewGroupCache(rdb, groupDB, groupMemberDB, groupRequestDB, SuperGroupMongoDriver, rockscache.Options{
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
	groupDB        relation2.GroupModelInterface
	groupMemberDB  relation2.GroupMemberModelInterface
	groupRequestDB relation2.GroupRequestModelInterface
	db             *gorm.DB

	cache   *cache.GroupCache
	mongoDB *unrelation.SuperGroupMongoDriver
}

func (g *GroupDataBase) CreateGroup(ctx context.Context, groups []*relation2.GroupModel, groupMembers []*relation2.GroupMemberModel) error {
	if len(groups) > 0 && len(groupMembers) > 0 {
		return g.db.Transaction(func(tx *gorm.DB) error {
			if err := g.groupDB.Create(ctx, groups, tx); err != nil {
				return err
			}
			return g.groupMemberDB.Create(ctx, groupMembers, tx)
		})
	}
	if len(groups) > 0 {
		return g.groupDB.Create(ctx, groups)
	}
	if len(groupMembers) > 0 {
		return g.groupMemberDB.Create(ctx, groupMembers)
	}
	return nil
}

func (g *GroupDataBase) TakeGroup(ctx context.Context, groupID string) (group *relation2.GroupModel, err error) {
	return g.groupDB.Take(ctx, groupID)
}

func (g *GroupDataBase) FindGroup(ctx context.Context, groupIDs []string) (groups []*relation2.GroupModel, err error) {
	return g.groupDB.Find(ctx, groupIDs)
}

func (g *GroupDataBase) SearchGroup(ctx context.Context, keyword string, pageNumber, showNumber int32) (int32, []*relation2.GroupModel, error) {
	return g.groupDB.Search(ctx, keyword, pageNumber, showNumber)
}

func (g *GroupDataBase) UpdateGroup(ctx context.Context, groupID string, data map[string]any) error {
	return g.groupDB.UpdateMap(ctx, groupID, data)
}

func (g *GroupDataBase) DismissGroup(ctx context.Context, groupID string) error {
	return g.db.Transaction(func(tx *gorm.DB) error {
		if err := g.groupDB.UpdateStatus(ctx, groupID, constant.GroupStatusDismissed, tx); err != nil {
			return err
		}
		return g.groupMemberDB.DeleteGroup(ctx, []string{groupID}, tx)
	})
}

func (g *GroupDataBase) TakeGroupMember(ctx context.Context, groupID string, userID string) (groupMember *relation2.GroupMemberModel, err error) {
	return g.groupMemberDB.Take(ctx, groupID, userID)
}

func (g *GroupDataBase) TakeGroupOwner(ctx context.Context, groupID string) (*relation2.GroupMemberModel, error) {
	return g.groupMemberDB.TakeOwner(ctx, groupID)
}

func (g *GroupDataBase) FindGroupMember(ctx context.Context, groupIDs []string, userIDs []string, roleLevels []int32) ([]*relation2.GroupMemberModel, error) {
	return g.groupMemberDB.Find(ctx, groupIDs, userIDs, roleLevels)
}

func (g *GroupDataBase) PageGroupMember(ctx context.Context, groupIDs []string, userIDs []string, roleLevels []int32, pageNumber, showNumber int32) (int32, []*relation2.GroupMemberModel, error) {
	return g.groupMemberDB.SearchMember(ctx, "", groupIDs, userIDs, roleLevels, pageNumber, showNumber)
}

func (g *GroupDataBase) SearchGroupMember(ctx context.Context, keyword string, groupIDs []string, userIDs []string, roleLevels []int32, pageNumber, showNumber int32) (int32, []*relation2.GroupMemberModel, error) {
	return g.groupMemberDB.SearchMember(ctx, keyword, groupIDs, userIDs, roleLevels, pageNumber, showNumber)
}

func (g *GroupDataBase) HandlerGroupRequest(ctx context.Context, groupID string, userID string, handledMsg string, handleResult int32, member *relation2.GroupMemberModel) error {
	if member == nil {
		return g.groupRequestDB.UpdateHandler(ctx, groupID, userID, handledMsg, handleResult)
	}
	return g.db.Transaction(func(tx *gorm.DB) error {
		if err := g.groupRequestDB.UpdateHandler(ctx, groupID, userID, handledMsg, handleResult, tx); err != nil {
			return err
		}
		return g.groupMemberDB.Create(ctx, []*relation2.GroupMemberModel{member}, tx)
	})
}

func (g *GroupDataBase) DeleteGroupMember(ctx context.Context, groupID string, userIDs []string) error {
	return g.groupMemberDB.Delete(ctx, groupID, userIDs)
}

func (g *GroupDataBase) MapGroupMemberUserID(ctx context.Context, groupIDs []string) (map[string][]string, error) {
	return g.groupMemberDB.FindJoinUserID(ctx, groupIDs)
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
		return nil
	})
}

func (g *GroupDataBase) UpdateGroupMember(ctx context.Context, groupID, userID string, data map[string]any) error {
	return g.groupMemberDB.Update(ctx, groupID, userID, data)
}

func (g *GroupDataBase) CreateGroupRequest(ctx context.Context, requests []*relation2.GroupRequestModel) error {
	return g.groupRequestDB.Create(ctx, requests)
}

func (g *GroupDataBase) TakeGroupRequest(ctx context.Context, groupID string, userID string) (*relation2.GroupRequestModel, error) {
	return g.groupRequestDB.Take(ctx, groupID, userID)
}

func (g *GroupDataBase) PageGroupRequestUser(ctx context.Context, userID string, pageNumber, showNumber int32) (int32, []*relation2.GroupRequestModel, error) {
	return g.groupRequestDB.Page(ctx, userID, pageNumber, showNumber)
}

func (g *GroupDataBase) FindSuperGroup(ctx context.Context, groupIDs []string) ([]*unrelation2.SuperGroupModel, error) {
	return g.mongoDB.FindSuperGroup(ctx, groupIDs)
}

func (g *GroupDataBase) FindJoinSuperGroup(ctx context.Context, userID string) (*unrelation2.UserToSuperGroupModel, error) {
	return g.mongoDB.GetSuperGroupByUserID(ctx, userID)
}

func (g *GroupDataBase) CreateSuperGroup(ctx context.Context, groupID string, initMemberIDList []string) error {
	return g.mongoDB.Transaction(ctx, func(s unrelation2.SuperGroupModelInterface, tx any) error {
		return s.CreateSuperGroup(ctx, groupID, initMemberIDList, tx)
	})
}

func (g *GroupDataBase) DeleteSuperGroup(ctx context.Context, groupID string) error {
	return g.mongoDB.Transaction(ctx, func(s unrelation2.SuperGroupModelInterface, tx any) error {
		return s.DeleteSuperGroup(ctx, groupID, tx)
	})
}

func (g *GroupDataBase) DeleteSuperGroupMember(ctx context.Context, groupID string, userIDs []string) error {
	return g.mongoDB.Transaction(ctx, func(s unrelation2.SuperGroupModelInterface, tx any) error {
		return s.RemoverUserFromSuperGroup(ctx, groupID, userIDs, tx)
	})
}

func (g *GroupDataBase) CreateSuperGroupMember(ctx context.Context, groupID string, userIDs []string) error {
	return g.mongoDB.Transaction(ctx, func(s unrelation2.SuperGroupModelInterface, tx any) error {
		return s.AddUserToSuperGroup(ctx, groupID, userIDs, tx)
	})
}
