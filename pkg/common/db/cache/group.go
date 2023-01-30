package cache

import (
	"Open_IM/pkg/common/db/relation"
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
	group        *relation.Group
	groupMember  *relation.GroupMember
	groupRequest *relation.GroupRequest
	expireTime   time.Duration
	redisClient  *RedisClient
	rcClient     *rockscache.Client
}

func NewGroupCache(rdb redis.UniversalClient, groupDB *relation.Group, groupMemberDB *relation.GroupMember, groupRequestDB *relation.GroupRequest, opts rockscache.Options) *GroupCache {
	return &GroupCache{rcClient: rockscache.NewClient(rdb, opts), expireTime: GroupExpireTime, group: groupDB, groupMember: groupMemberDB, groupRequest: groupRequestDB, redisClient: NewRedisClient(rdb)}
}

func (g *GroupCache) getRedisClient() *RedisClient {
	return g.redisClient
}

func (g *GroupCache) GetGroupsInfo(ctx context.Context, groupIDs []string) (groups []*relation.Group, err error) {
	for _, groupID := range groupIDs {
		group, err := g.GetGroupInfo(ctx, groupID)
		if err != nil {
			return nil, err
		}
		groups = append(groups, group)
	}
	return groups, nil
}

func (g *GroupCache) GetGroupInfo(ctx context.Context, groupID string) (group *relation.Group, err error) {
	getGroup := func() (string, error) {
		groupInfo, err := g.group.Take(ctx, groupID)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		bytes, err := json.Marshal(groupInfo)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		return string(bytes), nil
	}
	group = &relation.Group{}
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupID", groupID, "group", *group)
	}()
	groupStr, err := g.rcClient.Fetch(g.getGroupInfoCacheKey(groupID), g.expireTime, getGroup)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	err = json.Unmarshal([]byte(groupStr), group)
	return group, utils.Wrap(err, "")
}

func (g *GroupCache) DelGroupInfo(ctx context.Context, groupID string) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupID", groupID)
	}()
	return g.rcClient.TagAsDeleted(g.getGroupInfoCacheKey(groupID))
}

func (g *GroupCache) DelGroupsInfo(ctx context.Context, groupIDs []string) error {
	for _, groupID := range groupIDs {
		if err := g.DelGroupInfo(ctx, groupID); err != nil {
			return err
		}
	}
	return nil
}

func (g *GroupCache) getGroupInfoCacheKey(groupID string) string {
	return groupInfoCacheKey + groupID
}

func (g *GroupCache) DelJoinedSuperGroupIDs(ctx context.Context, userIDs []string) (err error) {
	for _, userID := range userIDs {
		if err := g.rcClient.TagAsDeleted(joinedSuperGroupListCache + userID); err != nil {
			return err
		}
	}
}

func (g *GroupCache) DelJoinedSuperGroupID(ctx context.Context, userID string) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "userID", userID)
	}()
	return g.rcClient.TagAsDeleted(joinedSuperGroupListCache + userID)
}
