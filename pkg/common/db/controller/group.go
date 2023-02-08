package controller

import (
	"Open_IM/pkg/common/constant"
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

type GroupInterface GroupDataBaseInterface

//type GroupInterface interface {
//	// Group
//	CreateGroup(ctx context.Context, groups []*relation2.GroupModel, groupMembers []*relation2.GroupMemberModel) error
//	TakeGroup(ctx context.Context, groupID string) (group *relation2.GroupModel, err error)
//	FindGroup(ctx context.Context, groupIDs []string) (groups []*relation2.GroupModel, err error)
//	SearchGroup(ctx context.Context, keyword string, pageNumber, showNumber int32) (int32, []*relation2.GroupModel, error)
//	UpdateGroup(ctx context.Context, groupID string, data map[string]any) error
//	DismissGroup(ctx context.Context, groupID string) error // 解散群，并删除群成员
//	// GroupMember
//	TakeGroupMember(ctx context.Context, groupID string, userID string) (groupMember *relation2.GroupMemberModel, err error)
//	TakeGroupOwner(ctx context.Context, groupID string) (*relation2.GroupMemberModel, error)
//	FindGroupMember(ctx context.Context, groupIDs []string, userIDs []string, roleLevels []int32) ([]*relation2.GroupMemberModel, error)
//	PageGroupMember(ctx context.Context, groupIDs []string, userIDs []string, roleLevels []int32, pageNumber, showNumber int32) (int32, []*relation2.GroupMemberModel, error)
//	SearchGroupMember(ctx context.Context, name string, groupIDs []string, userIDs []string, roleLevels []int32, pageNumber, showNumber int32) (int32, []*relation2.GroupMemberModel, error)
//	HandlerGroupRequest(ctx context.Context, groupID string, userID string, handledMsg string, handleResult int32, member *relation2.GroupMemberModel) error
//	DeleteGroupMember(ctx context.Context, groupID string, userIDs []string) error
//	MapGroupHash(ctx context.Context, groupIDs []string) (map[string]uint64, error)
//	MapGroupMemberNum(ctx context.Context, groupIDs []string) (map[string]int, error)
//	TransferGroupOwner(ctx context.Context, groupID string, oldOwnerUserID, newOwnerUserID string) error // 转让群
//	UpdateGroupMember(ctx context.Context, groupID, userID string, data map[string]any) error
//	// GroupRequest
//	CreateGroupRequest(ctx context.Context, requests []*relation2.GroupRequestModel) error
//	TakeGroupRequest(ctx context.Context, groupID string, userID string) (*relation2.GroupRequestModel, error)
//	PageGroupRequestUser(ctx context.Context, userID string, pageNumber, showNumber int32) (int32, []*relation2.GroupRequestModel, error)
//	// SuperGroup
//	TakeSuperGroup(ctx context.Context, groupID string) (superGroup *unrelation2.SuperGroupModel, err error)
//	FindJoinSuperGroup(ctx context.Context, userID string, pageNumber, showNumber int32) (total int32, groupIDs []string, err error)
//	CreateSuperGroup(ctx context.Context, groupID string, initMemberIDList []string) error
//	DeleteSuperGroup(ctx context.Context, groupID string) error
//	DeleteSuperGroupMember(ctx context.Context, groupID string, userIDs []string) error
//	CreateSuperGroupMember(ctx context.Context, groupID string, userIDs []string) error
//	MapSuperGroupMemberNum(ctx context.Context, groupIDs []string) (map[string]uint32, error)
//}
//
//var _ GroupInterface = (*GroupController)(nil)
//
//type GroupController struct {
//	database GroupDataBaseInterface
//}

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
	MapGroupHash(ctx context.Context, groupIDs []string) (map[string]uint64, error)
	MapGroupMemberNum(ctx context.Context, groupIDs []string) (map[string]int, error)
	TransferGroupOwner(ctx context.Context, groupID string, oldOwnerUserID, newOwnerUserID string) error // 转让群
	UpdateGroupMember(ctx context.Context, groupID, userID string, data map[string]any) error
	// GroupRequest
	CreateGroupRequest(ctx context.Context, requests []*relation2.GroupRequestModel) error
	TakeGroupRequest(ctx context.Context, groupID string, userID string) (*relation2.GroupRequestModel, error)
	PageGroupRequestUser(ctx context.Context, userID string, pageNumber, showNumber int32) (int32, []*relation2.GroupRequestModel, error)
	// SuperGroup
	TakeSuperGroup(ctx context.Context, groupID string) (superGroup *unrelation2.SuperGroupModel, err error)
	FindJoinSuperGroup(ctx context.Context, userID string, pageNumber, showNumber int32) (total int32, groupIDs []string, err error)
	CreateSuperGroup(ctx context.Context, groupID string, initMemberIDList []string) error
	DeleteSuperGroup(ctx context.Context, groupID string) error
	DeleteSuperGroupMember(ctx context.Context, groupID string, userIDs []string) error
	CreateSuperGroupMember(ctx context.Context, groupID string, userIDs []string) error
	MapSuperGroupMemberNum(ctx context.Context, groupIDs []string) (map[string]uint32, error)
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

//func (g *GroupDataBase) FindGroup(ctx context.Context, groupIDs []string) (groups []*relation2.GroupModel, err error) {
//	return g.groupDB.Find(ctx, groupIDs)
//}
//
//func (g *GroupDataBase) SearchGroup(ctx context.Context, name string, pageNumber, showNumber int32) (int32, []*relation2.GroupModel, error) {
//	return g.groupDB.Search(ctx, name, pageNumber, showNumber)
//}
//
//func (g *GroupDataBase) TakeGroup(ctx context.Context, groupID string) (group *relation2.GroupModel, err error) {
//	return g.groupDB.Take(ctx, groupID)
//}
//
//func (g *GroupDataBase) FindJoinedGroup(ctx context.Context, userID string, pageNumber, showNumber int32) (int32, []*relation2.GroupModel, error) {
//	total, members, err := g.groupMemberDB.PageByUser(ctx, userID, pageNumber, showNumber)
//	if err != nil {
//		return 0, nil, err
//	}
//	if len(members) == 0 {
//		return total, []*relation2.GroupModel{}, nil
//	}
//	groupIDs := utils.Slice(members, func(e *relation2.GroupMemberModel) string {
//		return e.GroupID
//	})
//	groups, err := g.groupDB.Find(ctx, groupIDs)
//	if err != nil {
//		return 0, nil, err
//	}
//	utils.OrderPtr(groupIDs, &groups, func(e *relation2.GroupModel) string {
//		return e.GroupID
//	})
//	return total, groups, nil
//}
//
//func (g *GroupDataBase) UpdateGroup(ctx context.Context, groupID string, data map[string]any) error {
//	return g.groupDB.UpdateMap(ctx, groupID, data)
//}
//
//func (g *GroupDataBase) DismissGroup(ctx context.Context, groupID string) error {
//	return utils.Wrap(g.db.Transaction(func(tx *gorm.DB) error {
//		if err := g.groupDB.UpdateStatus(ctx, groupID, constant.GroupStatusDismissed, tx); err != nil {
//			return err
//		}
//		return g.groupMemberDB.DeleteGroup(ctx, []string{groupID}, tx)
//	}), "")
//}
//
//func (g *GroupDataBase) CreateGroup(ctx context.Context, groups []*relation2.GroupModel, groupMembers []*relation2.GroupMemberModel) error {
//	if len(groups) > 0 && len(groupMembers) > 0 {
//		return g.db.Transaction(func(tx *gorm.DB) error {
//			if err := g.groupDB.Create(ctx, groups, tx); err != nil {
//				return err
//			}
//			return g.groupMemberDB.Create(ctx, groupMembers, tx)
//		})
//	}
//	if len(groups) > 0 {
//		return g.groupDB.Create(ctx, groups)
//	}
//	if len(groupMembers) > 0 {
//		return g.groupMemberDB.Create(ctx, groupMembers)
//	}
//	return nil
//}
//
//func (g *GroupDataBase) TakeGroupMember(ctx context.Context, groupID string, userID string) (groupMember *relation2.GroupMemberModel, err error) {
//	return g.groupMemberDB.Take(ctx, groupID, userID)
//}
//
//func (g *GroupDataBase) FindGroupMember(ctx context.Context, groupID string, userIDs []string) ([]*relation2.GroupMemberModel, error) {
//	return g.groupMemberDB.FindGroupUser(ctx, []string{groupID}, userIDs, nil)
//}
//
//func (g *GroupDataBase) FindGroupMemberAll(ctx context.Context, groupID string) ([]*relation2.GroupMemberModel, error) {
//	return g.groupMemberDB.FindGroupUser(ctx, []string{groupID}, nil, nil)
//}
//
//func (g *GroupDataBase) SearchGroupMember(ctx context.Context, groupID string, name string, pageNumber, showNumber int32) (int32, []*relation2.GroupMemberModel, error) {
//	return g.groupMemberDB.SearchMember(ctx, groupID, name, pageNumber, showNumber)
//}
//
//func (g *GroupDataBase) TakeGroupOwner(ctx context.Context, groupID string) (*relation2.GroupMemberModel, error) {
//	return g.groupMemberDB.TakeOwner(ctx, groupID)
//}
//
//func (g *GroupDataBase) FindGroupOwnerUser(ctx context.Context, groupIDs []string) ([]*relation2.GroupMemberModel, error) {
//	return g.groupMemberDB.FindGroupUser(ctx, groupIDs, nil, []int32{constant.GroupOwner})
//}

var _ GroupDataBaseInterface = (*GroupDataBase)(nil)

type GroupDataBase struct {
	groupDB        *relation.GroupGorm
	groupMemberDB  *relation.GroupMemberGorm
	groupRequestDB *relation.GroupRequestGorm
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

func (g *GroupDataBase) TakeGroupRequest(ctx context.Context, groupID string, userID string) (*relation2.GroupRequestModel, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupDataBase) PageGroupRequestUser(ctx context.Context, userID string, pageNumber, showNumber int32) (int32, []*relation2.GroupRequestModel, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupDataBase) TakeSuperGroup(ctx context.Context, groupID string) (superGroup *unrelation2.SuperGroupModel, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupDataBase) FindJoinSuperGroup(ctx context.Context, userID string, pageNumber, showNumber int32) (total int32, groupIDs []string, err error) {
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

func (g *GroupDataBase) CreateSuperGroupMember(ctx context.Context, groupID string, userIDs []string) error {
	//TODO implement me
	panic("implement me")
}

func (g *GroupDataBase) MapSuperGroupMemberNum(ctx context.Context, groupIDs []string) (map[string]uint32, error) {
	//TODO implement me
	panic("implement me")
}
