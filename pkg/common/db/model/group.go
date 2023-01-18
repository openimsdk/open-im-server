package model

import (
	"Open_IM/pkg/common/db/cache"
	"Open_IM/pkg/common/db/mongoDB"
	"Open_IM/pkg/common/db/mysql"
	"context"
	_ "github.com/dtm-labs/rockscache"
	"github.com/go-redis/redis/v8"
	"go.mongodb.org/mongo-driver/mongo"
	"gorm.io/gorm"
	//"time"
)

type GroupInterface interface {
	Find(ctx context.Context, groupIDs []string) (groups []*mysql.Group, err error)
	Create(ctx context.Context, groups []*mysql.Group) error
	Delete(ctx context.Context, groupIDs []string) error
	Take(ctx context.Context, groupID string) (group *mysql.Group, err error)
}

type GroupController struct {
	db    DataBase
	cache *cache.GroupCache
	mongo *mongoDB.Client
}
type DataBase interface {
	Find(ctx context.Context, groupIDs []string) (groups []*mysql.Group, err error)
	Create(ctx context.Context, groups []*mysql.Group) error
	Delete(ctx context.Context, groupIDs []string) error
	Take(ctx context.Context, groupID string) (group *mysql.Group, err error)
	DeleteTx(ctx context.Context, groupIDs []string) error
}
type MySqlDatabase struct {
	mysql.GroupModelInterface
}

func (m *MySqlDatabase) Delete(ctx context.Context, groupIDs []string) error {
	panic("implement me")
}

func NewMySqlDatabase(db mysql.GroupModelInterface) DataBase {
	return &MySqlDatabase{db}
}
func (m *MySqlDatabase) DeleteTx(ctx context.Context, groupIDs []string) error {
	return nil
}

func NewGroupController(groupModel mysql.GroupModelInterface, rdb redis.UniversalClient, mdb *mongo.Client) *GroupController {
	return &GroupController{db: NewMySqlDatabase(groupModel)}
	//groupModel.cache = cache.NewGroupCache(rdb, db, rockscache.Options{
	//	DisableCacheRead:  false,
	//	StrongConsistency: true,
	//})
	//groupModel.mongo = mongoDB.NewMongoClient(mdb)
	//return &groupModel
}

func (g *GroupController) Find(ctx context.Context, groupIDs []string) (groups []*mysql.Group, err error) {
	return g.cache.GetGroupsInfoFromCache(ctx, groupIDs)
}

func (g *GroupController) Create(ctx context.Context, groups []*mysql.Group) error {
	return g.db.Create(ctx, groups)
}

func (g *GroupController) Delete(ctx context.Context, groupIDs []string) error {
	err := g.db.DB.Transaction(func(tx *gorm.DB) error {
		if err := g.db.Delete(ctx, groupIDs, tx); err != nil {
			return err
		}
		if err := g.cache.DelGroupsInfoFromCache(ctx, groupIDs); err != nil {
			return err
		}
		return nil
	})
	return err
}

func (g *GroupController) Take(ctx context.Context, groupID string) (group *mysql.Group, err error) {
	return g.cache.GetGroupInfoFromCache(ctx, groupID)
}

func (g *GroupController) Update(ctx context.Context, groups []*mysql.Group) error {
	err := g.db.DB.Transaction(func(tx *gorm.DB) error {
		if err := g.db.Update(ctx, groups, tx); err != nil {
			return err
		}
		var groupIDs []string
		for _, group := range groups {
			groupIDs = append(groupIDs, group.GroupID)
		}
		if err := g.cache.DelGroupsInfoFromCache(ctx, groupIDs); err != nil {
			return err
		}
		return nil
	})
	return err
}
