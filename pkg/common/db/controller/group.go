package controller

import (
	"Open_IM/pkg/common/db/cache"
	"Open_IM/pkg/common/db/relation"
	relation2 "Open_IM/pkg/common/db/table/relation"
	unrelation2 "Open_IM/pkg/common/db/table/unrelation"
	"Open_IM/pkg/common/db/unrelation"
	"context"
	"github.com/dtm-labs/rockscache"
	_ "github.com/dtm-labs/rockscache"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type GroupInterface interface {
	// group
	FindGroup(ctx context.Context, groupIDs []string) (groups []*relation2.GroupModel, err error)
	SearchGroup(ctx context.Context, name string, pageNumber, showNumber int32) (int32, []*relation2.GroupModel, error)
	TakeGroup(ctx context.Context, groupID string) (group *relation2.GroupModel, err error)
	FindJoinedGroup(ctx context.Context, userID string, pageNumber, showNumber int32) ([]*relation2.GroupModel, error)
	UpdateGroup(ctx context.Context, groupID string, data map[string]any) error
	DismissGroup(ctx context.Context, groupID string) error // 解散群，并删除群成员
	// groupMember
	CreateGroup(ctx context.Context, groups []*relation2.GroupModel, groupMember []*relation2.GroupMemberModel) error
	TakeGroupMember(ctx context.Context, groupID string, userID string) (groupMember *relation2.GroupMemberModel, err error)
	FindGroupMember(ctx context.Context, groupID string, userIDs []string) ([]*relation2.GroupMemberModel, error)
	FindGroupMemberAll(ctx context.Context, groupID string) ([]*relation2.GroupMemberModel, error)
	FindGroupMemberFilterList(ctx context.Context, groupID string, filter int32, begin int32, maxNumber int32) ([]*relation2.GroupMemberModel, error) // relation.GetGroupMemberByGroupID(req.GroupID, req.Filter, req.NextSeq, 30)
	SearchGroupMember(ctx context.Context, groupID, name string, pageNumber, showNumber int32) (int32, []*relation2.GroupMemberModel, error)
	FindGroupMemberByUserID(ctx context.Context, groupID string, userIDs []string) ([]*relation2.GroupMemberModel, error)
	TakeGroupOwner(ctx context.Context, groupID string) (*relation2.GroupMemberModel, error)
	FindGroupOwnerUser(ctx context.Context, groupID []string) ([]*relation2.GroupMemberModel, error)
	CreateGroupMember(ctx context.Context, groupMember []*relation2.GroupMemberModel) error
	HandlerGroupRequest(ctx context.Context, groupID string, userID string, handledMsg string, handleResult int32, member *relation2.GroupMemberModel) error
	DeleteGroupMember(ctx context.Context, groupID string, userIDs []string) error
	MapGroupHash(ctx context.Context, groupIDs []string) (map[string]uint64, error)
	MapGroupMemberNum(ctx context.Context, groupIDs []string) (map[string]int, error)
	MapGroupOwnerUserID(ctx context.Context, groupIDs []string) (map[string]string, error)
	TransferGroupOwner(ctx context.Context, groupID string, oldOwnerUserID, newOwnerUserID string) error // 转让群
	UpdateGroupMember(ctx context.Context, groupID, userID string, data map[string]any) error

	CreateGroupRequest(ctx context.Context, requests []*relation2.GroupRequestModel) error
	GetGroupRecvApplicationList(ctx context.Context, userID string) ([]*relation2.GroupRequestModel, error) // ?
	TakeGroupRequest(ctx context.Context, groupID string, userID string) (*relation2.GroupRequestModel, error)
	FindUserGroupRequest(ctx context.Context, userID string, pageNumber, showNumber int32) (int32, []*relation2.GroupRequestModel, error)
	// superGroup
	TakeSuperGroup(ctx context.Context, groupID string) (superGroup *unrelation2.SuperGroupModel, err error)
	CreateSuperGroup(ctx context.Context, groupID string, initMemberIDList []string) error
	DeleteSuperGroup(ctx context.Context, groupID string) error
	DeleteSuperGroupMember(ctx context.Context, groupID string, userIDs []string) error
	AddUserToSuperGroup(ctx context.Context, groupID string, userIDs []string) error
	FindJoinSuperGroup(ctx context.Context, userID string, pageNumber, showNumber int32) (total int32, groupIDs []string, err error)
	MapSuperGroupMemberNum(ctx context.Context, groupIDs []string) (map[string]uint32, error)
}

var _ GroupInterface = (*GroupController)(nil)

type GroupController struct {
	database GroupDataBaseInterface
}

func (g *GroupController) FindGroup(ctx context.Context, groupIDs []string) (groups []*relation2.GroupModel, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupController) SearchGroup(ctx context.Context, name string, pageNumber, showNumber int32) (int32, []*relation2.GroupModel, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupController) TakeGroup(ctx context.Context, groupID string) (group *relation2.GroupModel, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupController) FindJoinedGroup(ctx context.Context, userID string, pageNumber, showNumber int32) ([]*relation2.GroupModel, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupController) UpdateGroup(ctx context.Context, groupID string, data map[string]any) error {
	//TODO implement me
	panic("implement me")
}

func (g *GroupController) DismissGroup(ctx context.Context, groupID string) error {
	//TODO implement me
	panic("implement me")
}

func (g *GroupController) CreateGroup(ctx context.Context, groups []*relation2.GroupModel, groupMember []*relation2.GroupMemberModel) error {
	//TODO implement me
	panic("implement me")
}

func (g *GroupController) TakeGroupMember(ctx context.Context, groupID string, userID string) (groupMember *relation2.GroupMemberModel, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupController) FindGroupMember(ctx context.Context, groupID string, userIDs []string) ([]*relation2.GroupMemberModel, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupController) FindGroupMemberAll(ctx context.Context, groupID string) ([]*relation2.GroupMemberModel, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupController) FindGroupMemberFilterList(ctx context.Context, groupID string, filter int32, begin int32, maxNumber int32) ([]*relation2.GroupMemberModel, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupController) FindGroupMembersByID(ctx context.Context, groupID string, userIDs []string) (groups []*relation2.GroupMemberModel, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupController) SearchGroupMember(ctx context.Context, groupID, name string, pageNumber, showNumber int32) (int32, []*relation2.GroupMemberModel, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupController) FindGroupMemberByUserID(ctx context.Context, groupID string, userIDs []string) ([]*relation2.GroupMemberModel, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupController) TakeGroupOwner(ctx context.Context, groupID string) (*relation2.GroupMemberModel, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupController) FindGroupOwnerUser(ctx context.Context, groupID []string) ([]*relation2.GroupMemberModel, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupController) CreateGroupMember(ctx context.Context, groupMember []*relation2.GroupMemberModel) error {
	//TODO implement me
	panic("implement me")
}

func (g *GroupController) HandlerGroupRequest(ctx context.Context, groupID string, userID string, handledMsg string, handleResult int32, member *relation2.GroupMemberModel) error {
	//TODO implement me
	panic("implement me")
}

func (g *GroupController) DeleteGroupMember(ctx context.Context, groupID string, userIDs []string) error {
	//TODO implement me
	panic("implement me")
}

func (g *GroupController) MapGroupHash(ctx context.Context, groupIDs []string) (map[string]uint64, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupController) MapGroupMemberNum(ctx context.Context, groupIDs []string) (map[string]int, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupController) MapGroupOwnerUserID(ctx context.Context, groupIDs []string) (map[string]string, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupController) TransferGroupOwner(ctx context.Context, groupID string, oldOwnerUserID, newOwnerUserID string) error {
	//TODO implement me
	panic("implement me")
}

func (g *GroupController) UpdateGroupMember(ctx context.Context, groupID, userID string, data map[string]any) error {
	//TODO implement me
	panic("implement me")
}

func (g *GroupController) CreateGroupRequest(ctx context.Context, requests []*relation2.GroupRequestModel) error {
	//TODO implement me
	panic("implement me")
}

func (g *GroupController) GetGroupRecvApplicationList(ctx context.Context, userID string) ([]*relation2.GroupRequestModel, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupController) TakeGroupRequest(ctx context.Context, groupID string, userID string) (*relation2.GroupRequestModel, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupController) FindUserGroupRequest(ctx context.Context, userID string, pageNumber, showNumber int32) (int32, []*relation2.GroupRequestModel, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupController) TakeSuperGroup(ctx context.Context, groupID string) (superGroup *unrelation2.SuperGroupModel, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupController) CreateSuperGroup(ctx context.Context, groupID string, initMemberIDList []string) error {
	//TODO implement me
	panic("implement me")
}

func (g *GroupController) DeleteSuperGroup(ctx context.Context, groupID string) error {
	//TODO implement me
	panic("implement me")
}

func (g *GroupController) DeleteSuperGroupMember(ctx context.Context, groupID string, userIDs []string) error {
	//TODO implement me
	panic("implement me")
}

func (g *GroupController) AddUserToSuperGroup(ctx context.Context, groupID string, userIDs []string) error {
	//TODO implement me
	panic("implement me")
}

func (g *GroupController) FindJoinSuperGroup(ctx context.Context, userID string, pageNumber, showNumber int32) (total int32, groupIDs []string, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupController) MapSuperGroupMemberNum(ctx context.Context, groupIDs []string) (map[string]uint32, error) {
	//TODO implement me
	panic("implement me")
}

type GroupDataBaseInterface interface {
	// group
	FindGroup(ctx context.Context, groupIDs []string) (groups []*relation2.GroupModel, err error)
	SearchGroup(ctx context.Context, name string, pageNumber, showNumber int32) (int32, []*relation2.GroupModel, error)
	TakeGroup(ctx context.Context, groupID string) (group *relation2.GroupModel, err error)
	FindJoinedGroup(ctx context.Context, userID string, pageNumber, showNumber int32) ([]*relation2.GroupModel, error)
	UpdateGroup(ctx context.Context, groupID string, data map[string]any) error
	DismissGroup(ctx context.Context, groupID string) error // 解散群，并删除群成员
	// groupMember
	CreateGroup(ctx context.Context, groups []*relation2.GroupModel, groupMember []*relation2.GroupMemberModel) error
	TakeGroupMember(ctx context.Context, groupID string, userID string) (groupMember *relation2.GroupMemberModel, err error)
	FindGroupMember(ctx context.Context, groupID string, userIDs []string) ([]*relation2.GroupMemberModel, error)
	FindGroupMemberAll(ctx context.Context, groupID string) ([]*relation2.GroupMemberModel, error)
	FindGroupMemberFilterList(ctx context.Context, groupID string, filter int32, begin int32, maxNumber int32) ([]*relation2.GroupMemberModel, error) // relation.GetGroupMemberByGroupID(req.GroupID, req.Filter, req.NextSeq, 30)
	SearchGroupMember(ctx context.Context, groupID, name string, pageNumber, showNumber int32) (int32, []*relation2.GroupMemberModel, error)
	FindGroupMemberByUserID(ctx context.Context, groupID string, userIDs []string) ([]*relation2.GroupMemberModel, error)
	TakeGroupOwner(ctx context.Context, groupID string) (*relation2.GroupMemberModel, error)
	FindGroupOwnerUser(ctx context.Context, groupID []string) ([]*relation2.GroupMemberModel, error)
	CreateGroupMember(ctx context.Context, groupMember []*relation2.GroupMemberModel) error
	HandlerGroupRequest(ctx context.Context, groupID string, userID string, handledMsg string, handleResult int32, member *relation2.GroupMemberModel) error
	DeleteGroupMember(ctx context.Context, groupID string, userIDs []string) error
	MapGroupHash(ctx context.Context, groupIDs []string) (map[string]uint64, error)
	MapGroupMemberNum(ctx context.Context, groupIDs []string) (map[string]int, error)
	MapGroupOwnerUserID(ctx context.Context, groupIDs []string) (map[string]string, error)
	TransferGroupOwner(ctx context.Context, groupID string, oldOwnerUserID, newOwnerUserID string) error // 转让群
	UpdateGroupMember(ctx context.Context, groupID, userID string, data map[string]any) error

	CreateGroupRequest(ctx context.Context, requests []*relation2.GroupRequestModel) error
	GetGroupRecvApplicationList(ctx context.Context, userID string) ([]*relation2.GroupRequestModel, error) // ?
	TakeGroupRequest(ctx context.Context, groupID string, userID string) (*relation2.GroupRequestModel, error)
	FindUserGroupRequest(ctx context.Context, userID string, pageNumber, showNumber int32) (int32, []*relation2.GroupRequestModel, error)
	// superGroup
	TakeSuperGroup(ctx context.Context, groupID string) (superGroup *unrelation2.SuperGroupModel, err error)
	CreateSuperGroup(ctx context.Context, groupID string, initMemberIDList []string) error
	DeleteSuperGroup(ctx context.Context, groupID string) error
	DeleteSuperGroupMember(ctx context.Context, groupID string, userIDs []string) error
	AddUserToSuperGroup(ctx context.Context, groupID string, userIDs []string) error
	FindJoinSuperGroup(ctx context.Context, userID string, pageNumber, showNumber int32) (total int32, groupIDs []string, err error)
	MapSuperGroupMemberNum(ctx context.Context, groupIDs []string) (map[string]uint32, error)
}

var _ *GroupDataBase = (GroupDataBaseInterface)(nil)

type GroupDataBase struct {
	groupDB        *relation.GroupGorm
	groupMemberDB  *relation.GroupMemberGorm
	groupRequestDB *relation.GroupRequestGorm
	db             *gorm.DB

	cache   *cache.GroupCache
	mongoDB *unrelation.SuperGroupMongoDriver
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

//func (g *GroupDataBase) FindGroupsByID(ctx context.Context, groupIDs []string) (groups []*relation2.GroupModel, err error) {
//	return g.cache.GetGroupsInfo(ctx, groupIDs)
//}
//
//func (g *GroupDataBase) CreateGroup(ctx context.Context, groups []*relation2.GroupModel, groupMembers []*relation2.GroupMemberModel) error {
//	return g.db.Transaction(func(tx *gorm.DB) error {
//		if len(groups) > 0 {
//			if err := g.groupDB.Create(ctx, groups, tx); err != nil {
//				return err
//			}
//		}
//		if len(groupMembers) > 0 {
//			if err := g.groupMemberDB.Create(ctx, groupMembers, tx); err != nil {
//				return err
//			}
//		}
//		return nil
//	})
//}
//
//func (g *GroupDataBase) DeleteGroupByIDs(ctx context.Context, groupIDs []string) error {
//	return g.groupDB.DB.Transaction(func(tx *gorm.DB) error {
//		if err := g.groupDB.Delete(ctx, groupIDs, tx); err != nil {
//			return err
//		}
//		if err := g.cache.DelGroupsInfo(ctx, groupIDs); err != nil {
//			return err
//		}
//		return nil
//	})
//}
//
//func (g *GroupDataBase) TakeGroupByID(ctx context.Context, groupID string) (group *relation2.GroupModel, err error) {
//	return g.cache.GetGroupInfo(ctx, groupID)
//}
//
//func (g *GroupDataBase) Update(ctx context.Context, groups []*relation2.GroupModel) error {
//	return g.db.Transaction(func(tx *gorm.DB) error {
//		if err := g.groupDB.Update(ctx, groups, tx); err != nil {
//			return err
//		}
//		var groupIDs []string
//		for _, group := range groups {
//			groupIDs = append(groupIDs, group.GroupID)
//		}
//		if err := g.cache.DelGroupsInfo(ctx, groupIDs); err != nil {
//			return err
//		}
//		return nil
//	})
//}
//
//func (g *GroupDataBase) GetJoinedGroupList(ctx context.Context, userID string) ([]*relation2.GroupModel, error) {
//
//	return nil, nil
//}
//
//func (g *GroupDataBase) CreateSuperGroup(ctx context.Context, groupID string, initMemberIDList []string) error {
//	sess, err := g.mongoDB.MgoClient.StartSession()
//	if err != nil {
//		return err
//	}
//	defer sess.EndSession(ctx)
//	sCtx := mongo.NewSessionContext(ctx, sess)
//	if err = g.mongoDB.CreateSuperGroup(sCtx, groupID, initMemberIDList); err != nil {
//		_ = sess.AbortTransaction(ctx)
//		return err
//	}
//
//	if err = g.cache.BatchDelJoinedSuperGroupIDs(ctx, initMemberIDList); err != nil {
//		_ = sess.AbortTransaction(ctx)
//		return err
//	}
//	return sess.CommitTransaction(ctx)
//}
//
//func (g *GroupDataBase) GetSuperGroupByID(ctx context.Context, groupID string) (superGroup *unrelation.SuperGroup, err error) {
//	return g.mongoDB.GetSuperGroup(ctx, groupID)
//}

func (g *GroupDataBase) FindGroup(ctx context.Context, groupIDs []string) (groups []*relation2.GroupModel, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupDataBase) SearchGroup(ctx context.Context, name string, pageNumber, showNumber int32) (int32, []*relation2.GroupModel, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupDataBase) TakeGroup(ctx context.Context, groupID string) (group *relation2.GroupModel, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupDataBase) FindJoinedGroup(ctx context.Context, userID string, pageNumber, showNumber int32) ([]*relation2.GroupModel, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupDataBase) UpdateGroup(ctx context.Context, groupID string, data map[string]any) error {
	//TODO implement me
	panic("implement me")
}

func (g *GroupDataBase) DismissGroup(ctx context.Context, groupID string) error {
	//TODO implement me
	panic("implement me")
}

func (g *GroupDataBase) CreateGroup(ctx context.Context, groups []*relation2.GroupModel, groupMember []*relation2.GroupMemberModel) error {
	//TODO implement me
	panic("implement me")
}

func (g *GroupDataBase) TakeGroupMember(ctx context.Context, groupID string, userID string) (groupMember *relation2.GroupMemberModel, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupDataBase) FindGroupMember(ctx context.Context, groupID string, userIDs []string) ([]*relation2.GroupMemberModel, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupDataBase) FindGroupMemberAll(ctx context.Context, groupID string) ([]*relation2.GroupMemberModel, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupDataBase) FindGroupMemberFilterList(ctx context.Context, groupID string, filter int32, begin int32, maxNumber int32) ([]*relation2.GroupMemberModel, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupDataBase) FindGroupMembersByID(ctx context.Context, groupID string, userIDs []string) (groups []*relation2.GroupMemberModel, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupDataBase) SearchGroupMember(ctx context.Context, groupID, name string, pageNumber, showNumber int32) (int32, []*relation2.GroupMemberModel, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupDataBase) FindGroupMemberByUserID(ctx context.Context, groupID string, userIDs []string) ([]*relation2.GroupMemberModel, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupDataBase) TakeGroupOwner(ctx context.Context, groupID string) (*relation2.GroupMemberModel, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupDataBase) FindGroupOwnerUser(ctx context.Context, groupID []string) ([]*relation2.GroupMemberModel, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupDataBase) CreateGroupMember(ctx context.Context, groupMember []*relation2.GroupMemberModel) error {
	//TODO implement me
	panic("implement me")
}

func (g *GroupDataBase) HandlerGroupRequest(ctx context.Context, groupID string, userID string, handledMsg string, handleResult int32, member *relation2.GroupMemberModel) error {
	//TODO implement me
	panic("implement me")
}

func (g *GroupDataBase) DeleteGroupMember(ctx context.Context, groupID string, userIDs []string) error {
	//TODO implement me
	panic("implement me")
}

func (g *GroupDataBase) MapGroupHash(ctx context.Context, groupIDs []string) (map[string]uint64, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupDataBase) MapGroupMemberNum(ctx context.Context, groupIDs []string) (map[string]int, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupDataBase) MapGroupOwnerUserID(ctx context.Context, groupIDs []string) (map[string]string, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupDataBase) TransferGroupOwner(ctx context.Context, groupID string, oldOwnerUserID, newOwnerUserID string) error {
	//TODO implement me
	panic("implement me")
}

func (g *GroupDataBase) UpdateGroupMember(ctx context.Context, groupID, userID string, data map[string]any) error {
	//TODO implement me
	panic("implement me")
}

func (g *GroupDataBase) CreateGroupRequest(ctx context.Context, requests []*relation2.GroupRequestModel) error {
	//TODO implement me
	panic("implement me")
}

func (g *GroupDataBase) GetGroupRecvApplicationList(ctx context.Context, userID string) ([]*relation2.GroupRequestModel, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupDataBase) TakeGroupRequest(ctx context.Context, groupID string, userID string) (*relation2.GroupRequestModel, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupDataBase) FindUserGroupRequest(ctx context.Context, userID string, pageNumber, showNumber int32) (int32, []*relation2.GroupRequestModel, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupDataBase) TakeSuperGroup(ctx context.Context, groupID string) (superGroup *unrelation2.SuperGroupModel, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupDataBase) CreateSuperGroup(ctx context.Context, groupID string, initMemberIDList []string) error {
	//TODO implement me
	panic("implement me")
}

func (g *GroupDataBase) DeleteSuperGroup(ctx context.Context, groupID string) error {
	//TODO implement me
	panic("implement me")
}

func (g *GroupDataBase) DeleteSuperGroupMember(ctx context.Context, groupID string, userIDs []string) error {
	//TODO implement me
	panic("implement me")
}

func (g *GroupDataBase) AddUserToSuperGroup(ctx context.Context, groupID string, userIDs []string) error {
	//TODO implement me
	panic("implement me")
}

func (g *GroupDataBase) FindJoinSuperGroup(ctx context.Context, userID string, pageNumber, showNumber int32) (total int32, groupIDs []string, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupDataBase) MapSuperGroupMemberNum(ctx context.Context, groupIDs []string) (map[string]uint32, error) {
	//TODO implement me
	panic("implement me")
}
