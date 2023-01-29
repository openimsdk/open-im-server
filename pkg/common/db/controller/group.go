package controller

import (
	"Open_IM/pkg/common/db/cache"
	"Open_IM/pkg/common/db/relation"
	"Open_IM/pkg/common/db/unrelation"
	"context"
	"github.com/dtm-labs/rockscache"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type GroupInterface interface {
	FindGroupsByID(ctx context.Context, groupIDs []string) (groups []*relation.Group, err error)
	CreateGroup(ctx context.Context, groups []*relation.Group) error
	DeleteGroupByIDs(ctx context.Context, groupIDs []string) error
	TakeGroupByID(ctx context.Context, groupID string) (group *relation.Group, err error)

	//mongo
	CreateSuperGroup(ctx context.Context, groupID string, initMemberIDList []string, memberNumCount int) error
	GetSuperGroupByID(ctx context.Context, groupID string) (superGroup *unrelation.SuperGroup, err error)
}

type GroupController struct {
	database DataBase
}

func NewGroupController(db *gorm.DB, rdb redis.UniversalClient, mgoDB *mongo.Database) GroupInterface {
	groupController := &GroupController{database: newGroupDatabase(db, rdb, mgoDB)}
	return groupController
}

func (g *GroupController) FindGroupsByID(ctx context.Context, groupIDs []string) (groups []*relation.Group, err error) {
	return g.database.Find(ctx, groupIDs)
}

func (g *GroupController) CreateGroup(ctx context.Context, groups []*relation.Group) error {
	return g.database.Create(ctx, groups)
}

func (g *GroupController) DeleteGroupByIDs(ctx context.Context, groupIDs []string) error {
	return g.database.Delete(ctx, groupIDs)
}

func (g *GroupController) TakeGroupByID(ctx context.Context, groupID string) (group *relation.Group, err error) {
	return g.database.Take(ctx, groupID)
}

func (g *GroupController) GetSuperGroupByID(ctx context.Context, groupID string) (superGroup *unrelation.SuperGroup, err error) {
	return g.database.GetSuperGroup(ctx, groupID)
}

func (g *GroupController) CreateSuperGroup(ctx context.Context, groupID string, initMemberIDList []string, memberNumCount int) error {
	return g.database.CreateSuperGroup(ctx, groupID, initMemberIDList, memberNumCount)
}

type DataBase interface {
	Find(ctx context.Context, groupIDs []string) (groups []*relation.Group, err error)
	Create(ctx context.Context, groups []*relation.Group) error
	Delete(ctx context.Context, groupIDs []string) error
	Take(ctx context.Context, groupID string) (group *relation.Group, err error)
	GetSuperGroup(ctx context.Context, groupID string) (superGroup *unrelation.SuperGroup, err error)
	CreateSuperGroup(ctx context.Context, groupID string, initMemberIDList []string, memberNumCount int) error
}

type GroupDataBase struct {
	groupDB        *relation.Group
	groupMemberDB  *relation.GroupMember
	groupRequestDB *relation.GroupRequest

	cache   *cache.GroupCache
	mongoDB *unrelation.SuperGroupMgoDB
}

func newGroupDatabase(db *gorm.DB, rdb redis.UniversalClient, mgoDB *mongo.Database) DataBase {
	groupDB := relation.NewGroupDB(db)
	groupMemberDB := relation.NewGroupMemberDB(db)
	groupRequestDB := relation.NewGroupRequest(db)
	database := &GroupDataBase{
		groupDB:        groupDB,
		groupMemberDB:  groupMemberDB,
		groupRequestDB: groupRequestDB,
		cache: cache.NewGroupCache(rdb, groupDB, groupMemberDB, groupRequestDB, rockscache.Options{
			RandomExpireAdjustment: 0.2,
			DisableCacheRead:       false,
			DisableCacheDelete:     false,
			StrongConsistency:      true,
		}),
		mongoDB: unrelation.NewSuperGroupMgoDB(mgoDB),
	}
	return database
}

func (g *GroupDataBase) Find(ctx context.Context, groupIDs []string) (groups []*relation.Group, err error) {
	return g.cache.GetGroupsInfo(ctx, groupIDs)
}

func (g *GroupDataBase) Create(ctx context.Context, groups []*relation.Group) error {
	return g.groupDB.Create(ctx, groups)
}

func (g *GroupDataBase) Delete(ctx context.Context, groupIDs []string) error {
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

func (g *GroupDataBase) Take(ctx context.Context, groupID string) (group *relation.Group, err error) {
	return g.cache.GetGroupInfo(ctx, groupID)
}

func (g *GroupDataBase) Update(ctx context.Context, groups []*relation.Group) error {
	return g.groupDB.DB.Transaction(func(tx *gorm.DB) error {
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

func (g *GroupDataBase) CreateSuperGroup(ctx context.Context, groupID string, initMemberIDList []string, memberNumCount int) error {
	return g.mongoDB.CreateSuperGroup(ctx, groupID, initMemberIDList, memberNumCount, g.cache.DelJoinedSuperGroupIDs)
}

func (g *GroupDataBase) GetSuperGroup(ctx context.Context, groupID string) (superGroup *unrelation.SuperGroup, err error) {
	return g.mongoDB.GetSuperGroup(ctx, groupID)
}
