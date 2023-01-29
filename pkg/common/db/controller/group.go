package controller

import (
	"Open_IM/pkg/common/db/cache"
	"Open_IM/pkg/common/db/relation"
	"Open_IM/pkg/common/db/unrelation"
	"context"
	"github.com/dtm-labs/rockscache"
	_ "github.com/dtm-labs/rockscache"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
)

type GroupInterface interface {
	FindGroupsByID(ctx context.Context, groupIDs []string) (groups []*relation.Group, err error)
	CreateGroup(ctx context.Context, groups []*relation.Group, groupMember []*relation.GroupMember) error
	DeleteGroupByIDs(ctx context.Context, groupIDs []string) error
	TakeGroupByID(ctx context.Context, groupID string) (group *relation.Group, err error)
	GetJoinedGroupList(ctx context.Context, userID string) ([]*relation.Group, error)
	GetGroupMemberNum(ctx context.Context, groupIDs []string) (map[string]int, error)
	GetGroupOwnerUserID(ctx context.Context, groupIDs []string) (map[string]string, error)
	//mongo
	CreateSuperGroup(ctx context.Context, groupID string, initMemberIDList []string) error
	GetSuperGroupByID(ctx context.Context, groupID string) (superGroup *unrelation.SuperGroup, err error)
}

type GroupController struct {
	database DataBase
}

func NewGroupController(db *gorm.DB, rdb redis.UniversalClient, mgoDB *mongo.Client) GroupInterface {
	groupController := &GroupController{database: newGroupDatabase(db, rdb, mgoDB)}
	return groupController
}

func (g *GroupController) FindGroupsByID(ctx context.Context, groupIDs []string) (groups []*relation.Group, err error) {
	return g.database.FindGroupsByID(ctx, groupIDs)
}

func (g *GroupController) CreateGroup(ctx context.Context, groups []*relation.Group, groupMember []*relation.GroupMember) error {
	return g.database.CreateGroup(ctx, groups, groupMember)
}

func (g *GroupController) DeleteGroupByIDs(ctx context.Context, groupIDs []string) error {
	return g.database.DeleteGroupByIDs(ctx, groupIDs)
}

func (g *GroupController) TakeGroupByID(ctx context.Context, groupID string) (group *relation.Group, err error) {
	return g.database.TakeGroupByID(ctx, groupID)
}

func (g *GroupController) GetSuperGroupByID(ctx context.Context, groupID string) (superGroup *unrelation.SuperGroup, err error) {
	return g.database.GetSuperGroupByID(ctx, groupID)
}

func (g *GroupController) CreateSuperGroup(ctx context.Context, groupID string, initMemberIDList []string) error {
	return g.database.CreateSuperGroup(ctx, groupID, initMemberIDList)
}

func (g *GroupController) GetJoinedGroupList(ctx context.Context, userID string) ([]*relation.Group, error) {
	return g.database.GetJoinedGroupList(ctx, userID)
}

type DataBase interface {
	FindGroupsByID(ctx context.Context, groupIDs []string) (groups []*relation.Group, err error)
	CreateGroup(ctx context.Context, groups []*relation.Group, groupMember []*relation.GroupMember) error
	DeleteGroupByIDs(ctx context.Context, groupIDs []string) error
	GetJoinedGroupList(ctx context.Context, userID string) ([]*relation.Group, error)

	TakeGroupByID(ctx context.Context, groupID string) (group *relation.Group, err error)
	GetSuperGroupByID(ctx context.Context, groupID string) (superGroup *unrelation.SuperGroup, err error)
	CreateSuperGroup(ctx context.Context, groupID string, initMemberIDList []string) error
}

type GroupDataBase struct {
	groupDB        *relation.Group
	groupMemberDB  *relation.GroupMember
	groupRequestDB *relation.GroupRequest
	db             *gorm.DB

	cache   *cache.GroupCache
	mongoDB *unrelation.SuperGroupMgoDB
}

func newGroupDatabase(db *gorm.DB, rdb redis.UniversalClient, mgoDB *mongo.Client) DataBase {
	groupDB := relation.NewGroupDB(db)
	groupMemberDB := relation.NewGroupMemberDB(db)
	groupRequestDB := relation.NewGroupRequest(db)
	newDB := db
	database := &GroupDataBase{
		groupDB:        groupDB,
		groupMemberDB:  groupMemberDB,
		groupRequestDB: groupRequestDB,
		db:             newDB,
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

func (g *GroupDataBase) FindGroupsByID(ctx context.Context, groupIDs []string) (groups []*relation.Group, err error) {
	return g.cache.GetGroupsInfo(ctx, groupIDs)
}

func (g *GroupDataBase) CreateGroup(ctx context.Context, groups []*relation.Group, groupMember []*relation.GroupMember) error {
	return g.db.Transaction(func(tx *gorm.DB) error {
		if err := g.groupDB.Create(ctx, groups, tx); err != nil {
			return err
		}
		if len(groupMember) > 0 {
			if err := g.groupMemberDB.Create(ctx, groupMember, tx); err != nil {
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

func (g *GroupDataBase) TakeGroupByID(ctx context.Context, groupID string) (group *relation.Group, err error) {
	return g.cache.GetGroupInfo(ctx, groupID)
}

func (g *GroupDataBase) Update(ctx context.Context, groups []*relation.Group) error {
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

func (g *GroupDataBase) GetJoinedGroupList(ctx context.Context, userID string) ([]*relation.Group, error) {

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

	if err = g.cache.DelJoinedSuperGroupIDs(ctx, initMemberIDList); err != nil {
		_ = sess.AbortTransaction(ctx)
		return err
	}
	return sess.CommitTransaction(ctx)
}

func (g *GroupDataBase) GetSuperGroupByID(ctx context.Context, groupID string) (superGroup *unrelation.SuperGroup, err error) {
	return g.mongoDB.GetSuperGroup(ctx, groupID)
}
