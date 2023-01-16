package model

import (
	"Open_IM/pkg/common/db/cache"
	"Open_IM/pkg/common/db/mysql"
	"Open_IM/pkg/common/trace_log"
	"Open_IM/pkg/utils"
	"context"
	"encoding/json"
	"github.com/dtm-labs/rockscache"
	"gorm.io/gorm"
	"time"
)

type GroupModel struct {
	strongRc *cache.RcClient
	weakRc   *cache.RcClient
	db       *mysql.Group
	rdb      *cache.RedisClient
}

const GroupExpireTime = time.Second * 300 * 60
const RandomExpireAdjustment = 0.2

//cache key
const groupInfoCache = "GROUP_INFO_CACHE:"

func NewGroupModel(ctx context.Context) {
	var groupModel GroupModel
	redisClient := cache.InitRedis(ctx)
	rdb := cache.NewRedisClient(redisClient)
	groupModel.rdb = rdb
	groupModel.db = mysql.NewGroupDB()
	groupModel.strongRc = cache.NewRcClient(redisClient, GroupExpireTime, rockscache.Options{
		RandomExpireAdjustment: RandomExpireAdjustment,
		DisableCacheRead:       false,
		DisableCacheDelete:     false,
		StrongConsistency:      true,
	})
	groupModel.weakRc = cache.NewRcClient(redisClient, GroupExpireTime, rockscache.Options{
		RandomExpireAdjustment: RandomExpireAdjustment,
		DisableCacheRead:       false,
		DisableCacheDelete:     false,
		StrongConsistency:      false,
	})
}

func (g *GroupModel) Find(ctx context.Context, groupIDs []string) (groups []*mysql.Group, err error) {
	for _, groupID := range groupIDs {
		group, err := g.getGroupInfoFromCache(ctx, groupID)
		if err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}
	return groups, nil
}

func (g *GroupModel) Create(ctx context.Context, groups []*mysql.Group) error {
	return g.db.Create(ctx, groups)
}

func (g *GroupModel) Delete(ctx context.Context, groupIDs []string) error {
	tx := g.db.DB.Begin()
	if err := g.db.Delete(ctx, groupIDs); err != nil {
		tx.Commit()
		return err
	}
	if err := g.deleteGroupsInCache(ctx, groupIDs); err != nil {
		tx.Rollback()
		return err
	}
	return nil
}

func (g *GroupModel) getGroupCacheKey(groupID string) string {
	return groupInfoCache + groupID
}

func (g *GroupModel) deleteGroupsInCache(ctx context.Context, groupIDs []string) error {
	for _, groupID := range groupIDs {
		if err := g.weakRc.Cache.TagAsDeleted(g.getGroupCacheKey(groupID)); err != nil {
			return err
		}
	}
	return nil
}

func (g *GroupModel) getGroupInfoFromCache(ctx context.Context, groupID string) (groupInfo *mysql.Group, err error) {
	getGroupInfo := func() (string, error) {
		groupInfo, err := mysql.GetGroupInfoByGroupID(groupID)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		bytes, err := json.Marshal(groupInfo)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		return string(bytes), nil
	}
	groupInfo = &mysql.Group{}
	defer func() {
		trace_log.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupID", groupID, "groupInfo", groupInfo)
	}()
	groupInfoStr, err := g.weakRc.Cache.Fetch(groupInfoCache+groupID, GroupExpireTime, getGroupInfo)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	err = json.Unmarshal([]byte(groupInfoStr), groupInfo)
	return groupInfo, utils.Wrap(err, "")
}
