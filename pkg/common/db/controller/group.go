package controller

import (
	"Open_IM/pkg/common/db/cache"
	"Open_IM/pkg/common/db/relation"
	"Open_IM/pkg/common/db/table"
	"Open_IM/pkg/common/db/unrelation"
	"context"
	"github.com/dtm-labs/rockscache"
	_ "github.com/dtm-labs/rockscache"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type GroupInterface interface {
	FindGroupsByID(ctx context.Context, groupIDs []string) (groups []*table.GroupModel, err error)
	CreateGroup(ctx context.Context, groups []*table.GroupModel, groupMember []*table.GroupMemberModel) error
	DeleteGroupByIDs(ctx context.Context, groupIDs []string) error
	TakeGroupByID(ctx context.Context, groupID string) (group *table.GroupModel, err error)
	TakeGroupMemberByID(ctx context.Context, groupID string, userID string) (groupMember *table.GroupModel, err error)
	GetJoinedGroupList(ctx context.Context, userID string) ([]*table.GroupModel, error)
	GetGroupMemberList(ctx context.Context, groupID string) ([]*table.GroupMemberModel, error)
	GetGroupMemberListByUserID(ctx context.Context, groupID string, userIDs []string) ([]*table.GroupMemberModel, error)
	GetGroupMemberFilterList(ctx context.Context, groupID string, filter int32, begin int32, maxNumber int32) ([]*table.GroupModel, error) // relation.GetGroupMemberByGroupID(req.GroupID, req.Filter, req.NextSeq, 30)
	FindGroupMembersByID(ctx context.Context, groupID string, userIDs []string) (groups []*table.GroupMemberModel, err error)
	DelGroupMember(ctx context.Context, groupID string, userIDs []string) error
	GetGroupMemberNum(ctx context.Context, groupIDs []string) (map[string]int, error)
	GetGroupOwnerUserID(ctx context.Context, groupIDs []string) (map[string]string, error)
	GetGroupRecvApplicationList(ctx context.Context, userID string) ([]*table.GroupRequestModel, error)

	CreateGroupMember(ctx context.Context, groupMember []*table.GroupMemberModel) error
	CreateGroupRequest(ctx context.Context, requests []*table.GroupRequestModel) error

	//mongo
	CreateSuperGroup(ctx context.Context, groupID string, initMemberIDList []string) error
	DelSuperGroupMember(ctx context.Context, groupID string, userIDs []string) error
	AddUserToSuperGroup(ctx context.Context, groupID string, userIDs []string) error
	GetSuperGroupByID(ctx context.Context, groupID string) (superGroup *unrelation.SuperGroup, err error)
}

var _ GroupInterface = (*GroupController)(nil)

type GroupController struct {
	database GroupDataBaseInterface
}

func (g *GroupController) TakeGroupMemberByID(ctx context.Context, groupID string, userID string) (groupMember *table.GroupModel, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupController) FindGroupMembersByID(ctx context.Context, groupID string, userIDs []string) (groups []*table.GroupModel, err error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupController) DelGroupMember(ctx context.Context, groupID string, userIDs []string) error {
	//TODO implement me
	panic("implement me")
}

func (g *GroupController) GetGroupRecvApplicationList(ctx context.Context, userID string) ([]*table.GroupRequestModel, error) {
	/*
		var groupRequestList []db.GroupRequest
			memberList, err := GetGroupMemberListByUserID(userID)
			if err != nil {
				return nil, err
			}
			for _, v := range memberList {
				if v.RoleLevel > constant.GroupOrdinaryUsers {
					list, err := GetGroupRequestByGroupID(v.GroupID)
					if err != nil {
						//		fmt.Println("111 GetGroupRequestByGroupID failed ", err.Error())
						continue
					}
					//	fmt.Println("222 GetGroupRequestByGroupID ok ", list)
					groupRequestList = append(groupRequestList, list...)
					//	fmt.Println("333 GetGroupRequestByGroupID ok ", groupRequestList)
				}
			}
			return groupRequestList, nil
	*/
	//TODO implement me
	panic("implement me")
}

func (g *GroupController) DelSuperGroupMember(ctx context.Context, groupID string, userIDs []string) error {
	//TODO implement me
	panic("implement me")
}

func (g *GroupController) GetJoinedGroupList(ctx context.Context, userID string) ([]*table.GroupModel, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupController) GetGroupMemberList(ctx context.Context, groupID string) ([]*table.GroupModel, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupController) GetGroupMemberListByUserID(ctx context.Context, groupID string, userIDs []string) ([]*table.GroupModel, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupController) GetGroupMemberFilterList(ctx context.Context, groupID string, filter int32, begin int32, maxNumber int32) ([]*table.GroupModel, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupController) GetGroupMemberNum(ctx context.Context, groupIDs []string) (map[string]int, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupController) GetGroupOwnerUserID(ctx context.Context, groupIDs []string) (map[string]string, error) {
	//TODO implement me
	panic("implement me")
}

func (g *GroupController) CreateGroupMember(ctx context.Context, groupMember []*table.GroupModel) error {
	//TODO implement me
	panic("implement me")
}

func (g *GroupController) CreateGroupRequest(ctx context.Context, requests []*table.GroupRequestModel) error {
	//TODO implement me
	panic("implement me")
}

func (g *GroupController) AddUserToSuperGroup(ctx context.Context, groupID string, userIDs []string) error {
	//TODO implement me
	panic("implement me")
}

func NewGroupController(db *gorm.DB, rdb redis.UniversalClient, mgoClient *mongo.Client) GroupInterface {
	groupController := &GroupController{database: newGroupDatabase(db, rdb, mgoClient)}
	return groupController
}

func (g *GroupController) FindGroupsByID(ctx context.Context, groupIDs []string) (groups []*table.GroupModel, err error) {
	return g.database.FindGroupsByID(ctx, groupIDs)
}

func (g *GroupController) CreateGroup(ctx context.Context, groups []*table.GroupModel, groupMember []*table.GroupModel) error {
	return g.database.CreateGroup(ctx, groups, groupMember)
}

func (g *GroupController) DeleteGroupByIDs(ctx context.Context, groupIDs []string) error {
	return g.database.DeleteGroupByIDs(ctx, groupIDs)
}

func (g *GroupController) TakeGroupByID(ctx context.Context, groupID string) (group *table.GroupModel, err error) {
	return g.database.TakeGroupByID(ctx, groupID)
}

func (g *GroupController) GetSuperGroupByID(ctx context.Context, groupID string) (superGroup *unrelation.SuperGroup, err error) {
	return g.database.GetSuperGroupByID(ctx, groupID)
}

func (g *GroupController) CreateSuperGroup(ctx context.Context, groupID string, initMemberIDList []string) error {
	return g.database.CreateSuperGroup(ctx, groupID, initMemberIDList)
}

type GroupDataBaseInterface interface {
	FindGroupsByID(ctx context.Context, groupIDs []string) (groups []*table.GroupModel, err error)
	CreateGroup(ctx context.Context, groups []*table.GroupModel, groupMember []*table.GroupModel) error
	DeleteGroupByIDs(ctx context.Context, groupIDs []string) error
	TakeGroupByID(ctx context.Context, groupID string) (group *table.GroupModel, err error)
	GetSuperGroupByID(ctx context.Context, groupID string) (superGroup *unrelation.SuperGroup, err error)
	CreateSuperGroup(ctx context.Context, groupID string, initMemberIDList []string) error
}

type GroupDataBase struct {
	groupDB        *relation.GroupGorm
	groupMemberDB  *relation.GroupMemberGorm
	groupRequestDB *relation.GroupRequestGorm
	db             *gorm.DB

	cache   *cache.GroupCache
	mongoDB *unrelation.SuperGroupMgoDB
}

func newGroupDatabase(db *gorm.DB, rdb redis.UniversalClient, mgoClient *mongo.Client) GroupDataBaseInterface {
	groupDB := relation.NewGroupDB(db)
	groupMemberDB := relation.NewGroupMemberDB(db)
	groupRequestDB := relation.NewGroupRequest(db)
	newDB := *db
	superGroupMgoDB := unrelation.NewSuperGroupMgoDB(mgoClient)
	database := &GroupDataBase{
		groupDB:        groupDB,
		groupMemberDB:  groupMemberDB,
		groupRequestDB: groupRequestDB,
		db:             &newDB,
		cache: cache.NewGroupCache(rdb, groupDB, groupMemberDB, groupRequestDB, superGroupMgoDB, rockscache.Options{
			RandomExpireAdjustment: 0.2,
			DisableCacheRead:       false,
			DisableCacheDelete:     false,
			StrongConsistency:      true,
		}),
		mongoDB: superGroupMgoDB,
	}
	return database
}

func (g *GroupDataBase) FindGroupsByID(ctx context.Context, groupIDs []string) (groups []*table.GroupModel, err error) {
	return g.cache.GetGroupsInfo(ctx, groupIDs)
}

func (g *GroupDataBase) CreateGroup(ctx context.Context, groups []*table.GroupModel, groupMembers []*table.GroupMemberModel) error {
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
		}
		return nil
	})
}

func (g *GroupDataBase) DeleteGroupByIDs(ctx context.Context, groupIDs []string) error {
	return g.groupDB.DB.Transaction(func(tx *gorm.DB) error {
		if err := g.groupDB.Delete(ctx, groupIDs, tx); err != nil {
			return err
		}
		if err := g.cache.DelGroupsInfo(ctx, groupIDs); err != nil {
			return err
		}
		return nil
	})
}

func (g *GroupDataBase) TakeGroupByID(ctx context.Context, groupID string) (group *table.GroupModel, err error) {
	return g.cache.GetGroupInfo(ctx, groupID)
}

func (g *GroupDataBase) Update(ctx context.Context, groups []*table.GroupModel) error {
	return g.db.Transaction(func(tx *gorm.DB) error {
		if err := g.groupDB.Update(ctx, groups, tx); err != nil {
			return err
		}
		var groupIDs []string
		for _, group := range groups {
			groupIDs = append(groupIDs, group.GroupID)
		}
		if err := g.cache.DelGroupsInfo(ctx, groupIDs); err != nil {
			return err
		}
		return nil
	})
}

func (g *GroupDataBase) GetJoinedGroupList(ctx context.Context, userID string) ([]*table.GroupModel, error) {

	return nil, nil
}

func (g *GroupDataBase) CreateSuperGroup(ctx context.Context, groupID string, initMemberIDList []string) error {
	sess, err := g.mongoDB.MgoClient.StartSession()
	if err != nil {
		return err
	}
	defer sess.EndSession(ctx)
	sCtx := mongo.NewSessionContext(ctx, sess)
	if err = g.mongoDB.CreateSuperGroup(sCtx, groupID, initMemberIDList); err != nil {
		_ = sess.AbortTransaction(ctx)
		return err
	}

	if err = g.cache.BatchDelJoinedSuperGroupIDs(ctx, initMemberIDList); err != nil {
		_ = sess.AbortTransaction(ctx)
		return err
	}
	return sess.CommitTransaction(ctx)
}

func (g *GroupDataBase) GetSuperGroupByID(ctx context.Context, groupID string) (superGroup *unrelation.SuperGroup, err error) {
	return g.mongoDB.GetSuperGroup(ctx, groupID)
}
