package cache

import (
	"Open_IM/pkg/common/db/mysql"
	"Open_IM/pkg/common/trace_log"
	"Open_IM/pkg/utils"
	"context"
	"encoding/json"
	"github.com/dtm-labs/rockscache"
	"github.com/go-redis/redis/v8"
	"time"
)

const GroupExpireTime = time.Second * 60 * 60 * 12
const groupInfoCacheKey = "GROUP_INFO_CACHE:"

type GroupCache struct {
	db          *mysql.Group
	expireTime  time.Duration
	redisClient *RedisClient
	rcClient    *rockscache.Client
}

func NewGroupCache(rdb redis.UniversalClient, db *mysql.Group, opts rockscache.Options) *GroupCache {
	redisClient := NewRedisClient(rdb)
	return &GroupCache{rcClient: rockscache.NewClient(rdb, opts), expireTime: GroupExpireTime, db: db, redisClient: redisClient}
}

func (g *GroupCache) getRedisClient() *RedisClient {
	return g.redisClient
}

func (g *GroupCache) GetGroupsInfoFromCache(ctx context.Context, groupIDs []string) (groups []*mysql.Group, err error) {
	for _, groupID := range groupIDs {
		group, err := g.GetGroupInfoFromCache(ctx, groupID)
		if err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}
	return groups, nil
}

func (g *GroupCache) GetGroupInfoFromCache(ctx context.Context, groupID string) (group *mysql.Group, err error) {
	getGroup := func() (string, error) {
		groupInfo, err := g.db.Take(ctx, groupID)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		bytes, err := json.Marshal(groupInfo)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		return string(bytes), nil
	}
	group = &mysql.Group{}
	defer func() {
		trace_log.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupID", groupID, "group", *group)
	}()
	groupStr, err := g.rcClient.Fetch(g.getGroupInfoCacheKey(groupID), g.expireTime, getGroup)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	err = json.Unmarshal([]byte(groupStr), group)
	return group, utils.Wrap(err, "")
}

func (g *GroupCache) DelGroupInfoFromCache(ctx context.Context, groupID string) (err error) {
	defer func() {
		trace_log.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupID", groupID)
	}()
	return g.rcClient.TagAsDeleted(g.getGroupInfoCacheKey(groupID))
}

func (g *GroupCache) DelGroupsInfoFromCache(ctx context.Context, groupIDs []string) error {
	for _, groupID := range groupIDs {
		if err := g.DelGroupInfoFromCache(ctx, groupID); err != nil {
			return err
		}
	}
	return nil
}

func (g *GroupCache) getGroupInfoCacheKey(groupID string) string {
	return groupInfoCacheKey + groupID
}
