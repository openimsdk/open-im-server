package model

import (
	"Open_IM/pkg/common/db/cache"
	"Open_IM/pkg/common/db/mongo"
	"Open_IM/pkg/common/db/mysql"
	"Open_IM/pkg/common/trace_log"
	"Open_IM/pkg/utils"
	"context"
	"encoding/json"
	"github.com/dtm-labs/rockscache"
	"gorm.io/gorm"
	//"time"
)

type GroupModel struct {
	db    *mysql.Group
	cache *cache.GroupCache
	mongo *mongo.Client
}

func NewGroupModel() *GroupModel {
	var groupModel GroupModel
	groupModel.db = mysql.NewGroupDB()
	groupModel.cache = cache.NewGroupCache(cache.InitRedis(), groupModel.db, rockscache.Options{
		DisableCacheRead:  false,
		StrongConsistency: true,
	})
	groupModel.mongo = mongo.NewMongoClient()
	return &groupModel
}

func (g *GroupModel) Find(ctx context.Context, groupIDs []string) (groups []*mysql.Group, err error) {
	return g.cache.GetGroupsInfoFromCache(ctx, groupIDs)
}

func (g *GroupModel) Create(ctx context.Context, groups []*mysql.Group) error {
	return g.db.Create(ctx, groups)
}

func (g *GroupModel) Delete(ctx context.Context, groupIDs []string) error {
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

func (g *GroupModel) Take(ctx context.Context, groupID string) (group *mysql.Group, err error) {
	return g.cache.GetGroupInfoFromCache(ctx, groupID)
}
