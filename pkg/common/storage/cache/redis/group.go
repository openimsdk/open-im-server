package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/cachekey"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/common"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/redis/go-redis/v9"
)

const (
	groupExpireTime = time.Second * 60 * 60 * 12
)

type GroupCacheRedis struct {
	cache.BatchDeleter
	groupDB        database.Group
	groupMemberDB  database.GroupMember
	groupRequestDB database.GroupRequest
	expireTime     time.Duration
	rcClient       *rocksCacheClient
	groupHash      cache.GroupHash
}

func NewGroupCacheRedis(rdb redis.UniversalClient, localCache *config.LocalCache, groupDB database.Group, groupMemberDB database.GroupMember, groupRequestDB database.GroupRequest, hashCode cache.GroupHash) cache.GroupCache {
	rc := newRocksCacheClient(rdb)
	return &GroupCacheRedis{
		BatchDeleter:   rc.GetBatchDeleter(localCache.Group.Topic),
		rcClient:       rc,
		expireTime:     groupExpireTime,
		groupDB:        groupDB,
		groupMemberDB:  groupMemberDB,
		groupRequestDB: groupRequestDB,
		groupHash:      hashCode,
	}
}

func (g *GroupCacheRedis) CloneGroupCache() cache.GroupCache {
	return &GroupCacheRedis{
		BatchDeleter:   g.BatchDeleter.Clone(),
		rcClient:       g.rcClient,
		expireTime:     g.expireTime,
		groupDB:        g.groupDB,
		groupMemberDB:  g.groupMemberDB,
		groupRequestDB: g.groupRequestDB,
	}
}

func (g *GroupCacheRedis) getGroupInfoKey(groupID string) string {
	return cachekey.GetGroupInfoKey(groupID)
}

func (g *GroupCacheRedis) getJoinedGroupsKey(userID string) string {
	return cachekey.GetJoinedGroupsKey(userID)
}

func (g *GroupCacheRedis) getGroupMembersHashKey(groupID string) string {
	return cachekey.GetGroupMembersHashKey(groupID)
}

func (g *GroupCacheRedis) getGroupMemberIDsKey(groupID string) string {
	return cachekey.GetGroupMemberIDsKey(groupID)
}

func (g *GroupCacheRedis) getGroupMemberInfoKey(groupID, userID string) string {
	return cachekey.GetGroupMemberInfoKey(groupID, userID)
}

func (g *GroupCacheRedis) getGroupMemberNumKey(groupID string) string {
	return cachekey.GetGroupMemberNumKey(groupID)
}

func (g *GroupCacheRedis) getGroupRoleLevelMemberIDsKey(groupID string, roleLevel int32) string {
	return cachekey.GetGroupRoleLevelMemberIDsKey(groupID, roleLevel)
}

func (g *GroupCacheRedis) getGroupMemberMaxVersionKey(groupID string) string {
	return cachekey.GetGroupMemberMaxVersionKey(groupID)
}

func (g *GroupCacheRedis) getJoinGroupMaxVersionKey(userID string) string {
	return cachekey.GetJoinGroupMaxVersionKey(userID)
}

func (g *GroupCacheRedis) getGroupID(group *model.Group) string {
	return group.GroupID
}

func (g *GroupCacheRedis) GetGroupsInfo(ctx context.Context, groupIDs []string) (groups []*model.Group, err error) {
	return batchGetCache2(ctx, g.rcClient, g.expireTime, groupIDs, g.getGroupInfoKey, g.getGroupID, g.groupDB.Find)
}

func (g *GroupCacheRedis) GetGroupInfo(ctx context.Context, groupID string) (group *model.Group, err error) {
	return getCache(ctx, g.rcClient, g.getGroupInfoKey(groupID), g.expireTime, func(ctx context.Context) (*model.Group, error) {
		return g.groupDB.Take(ctx, groupID)
	})
}

func (g *GroupCacheRedis) DelGroupsInfo(groupIDs ...string) cache.GroupCache {
	newGroupCache := g.CloneGroupCache()
	keys := make([]string, 0, len(groupIDs))
	for _, groupID := range groupIDs {
		keys = append(keys, g.getGroupInfoKey(groupID))
	}
	newGroupCache.AddKeys(keys...)

	return newGroupCache
}

func (g *GroupCacheRedis) DelGroupsOwner(groupIDs ...string) cache.GroupCache {
	newGroupCache := g.CloneGroupCache()
	keys := make([]string, 0, len(groupIDs))
	for _, groupID := range groupIDs {
		keys = append(keys, g.getGroupRoleLevelMemberIDsKey(groupID, constant.GroupOwner))
	}
	newGroupCache.AddKeys(keys...)

	return newGroupCache
}

func (g *GroupCacheRedis) DelGroupRoleLevel(groupID string, roleLevels []int32) cache.GroupCache {
	newGroupCache := g.CloneGroupCache()
	keys := make([]string, 0, len(roleLevels))
	for _, roleLevel := range roleLevels {
		keys = append(keys, g.getGroupRoleLevelMemberIDsKey(groupID, roleLevel))
	}
	newGroupCache.AddKeys(keys...)
	return newGroupCache
}

func (g *GroupCacheRedis) DelGroupAllRoleLevel(groupID string) cache.GroupCache {
	return g.DelGroupRoleLevel(groupID, []int32{constant.GroupOwner, constant.GroupAdmin, constant.GroupOrdinaryUsers})
}

func (g *GroupCacheRedis) GetGroupMembersHash(ctx context.Context, groupID string) (hashCode uint64, err error) {
	if g.groupHash == nil {
		return 0, errs.ErrInternalServer.WrapMsg("group hash is nil")
	}
	return getCache(ctx, g.rcClient, g.getGroupMembersHashKey(groupID), g.expireTime, func(ctx context.Context) (uint64, error) {
		return g.groupHash.GetGroupHash(ctx, groupID)
	})
}

func (g *GroupCacheRedis) GetGroupMemberHashMap(ctx context.Context, groupIDs []string) (map[string]*common.GroupSimpleUserID, error) {
	if g.groupHash == nil {
		return nil, errs.ErrInternalServer.WrapMsg("group hash is nil")
	}
	res := make(map[string]*common.GroupSimpleUserID)
	for _, groupID := range groupIDs {
		hash, err := g.GetGroupMembersHash(ctx, groupID)
		if err != nil {
			return nil, err
		}
		log.ZDebug(ctx, "GetGroupMemberHashMap", "groupID", groupID, "hash", hash)
		num, err := g.GetGroupMemberNum(ctx, groupID)
		if err != nil {
			return nil, err
		}
		res[groupID] = &common.GroupSimpleUserID{Hash: hash, MemberNum: uint32(num)}
	}

	return res, nil
}

func (g *GroupCacheRedis) DelGroupMembersHash(groupID string) cache.GroupCache {
	cache := g.CloneGroupCache()
	cache.AddKeys(g.getGroupMembersHashKey(groupID))

	return cache
}

func (g *GroupCacheRedis) GetGroupMemberIDs(ctx context.Context, groupID string) (groupMemberIDs []string, err error) {
	return getCache(ctx, g.rcClient, g.getGroupMemberIDsKey(groupID), g.expireTime, func(ctx context.Context) ([]string, error) {
		return g.groupMemberDB.FindMemberUserID(ctx, groupID)
	})
}

func (g *GroupCacheRedis) DelGroupMemberIDs(groupID string) cache.GroupCache {
	cache := g.CloneGroupCache()
	cache.AddKeys(g.getGroupMemberIDsKey(groupID))

	return cache
}

func (g *GroupCacheRedis) findUserJoinedGroupID(ctx context.Context, userID string) ([]string, error) {
	groupIDs, err := g.groupMemberDB.FindUserJoinedGroupID(ctx, userID)
	if err != nil {
		return nil, err
	}
	return g.groupDB.FindJoinSortGroupID(ctx, groupIDs)
}

func (g *GroupCacheRedis) GetJoinedGroupIDs(ctx context.Context, userID string) (joinedGroupIDs []string, err error) {
	return getCache(ctx, g.rcClient, g.getJoinedGroupsKey(userID), g.expireTime, func(ctx context.Context) ([]string, error) {
		return g.findUserJoinedGroupID(ctx, userID)
	})
}

func (g *GroupCacheRedis) DelJoinedGroupID(userIDs ...string) cache.GroupCache {
	keys := make([]string, 0, len(userIDs))
	for _, userID := range userIDs {
		keys = append(keys, g.getJoinedGroupsKey(userID))
	}
	cache := g.CloneGroupCache()
	cache.AddKeys(keys...)

	return cache
}

func (g *GroupCacheRedis) GetGroupMemberInfo(ctx context.Context, groupID, userID string) (groupMember *model.GroupMember, err error) {
	return getCache(ctx, g.rcClient, g.getGroupMemberInfoKey(groupID, userID), g.expireTime, func(ctx context.Context) (*model.GroupMember, error) {
		return g.groupMemberDB.Take(ctx, groupID, userID)
	})
}

func (g *GroupCacheRedis) GetGroupMembersInfo(ctx context.Context, groupID string, userIDs []string) ([]*model.GroupMember, error) {
	return batchGetCache2(ctx, g.rcClient, g.expireTime, userIDs, func(userID string) string {
		return g.getGroupMemberInfoKey(groupID, userID)
	}, func(member *model.GroupMember) string {
		return member.UserID
	}, func(ctx context.Context, userIDs []string) ([]*model.GroupMember, error) {
		return g.groupMemberDB.Find(ctx, groupID, userIDs)
	})
}

func (g *GroupCacheRedis) GetAllGroupMembersInfo(ctx context.Context, groupID string) (groupMembers []*model.GroupMember, err error) {
	groupMemberIDs, err := g.GetGroupMemberIDs(ctx, groupID)
	if err != nil {
		return nil, err
	}

	return g.GetGroupMembersInfo(ctx, groupID, groupMemberIDs)
}

func (g *GroupCacheRedis) DelGroupMembersInfo(groupID string, userIDs ...string) cache.GroupCache {
	keys := make([]string, 0, len(userIDs))
	for _, userID := range userIDs {
		keys = append(keys, g.getGroupMemberInfoKey(groupID, userID))
	}
	cache := g.CloneGroupCache()
	cache.AddKeys(keys...)

	return cache
}

func (g *GroupCacheRedis) GetGroupMemberNum(ctx context.Context, groupID string) (memberNum int64, err error) {
	return getCache(ctx, g.rcClient, g.getGroupMemberNumKey(groupID), g.expireTime, func(ctx context.Context) (int64, error) {
		return g.groupMemberDB.TakeGroupMemberNum(ctx, groupID)
	})
}

func (g *GroupCacheRedis) DelGroupsMemberNum(groupID ...string) cache.GroupCache {
	keys := make([]string, 0, len(groupID))
	for _, groupID := range groupID {
		keys = append(keys, g.getGroupMemberNumKey(groupID))
	}
	cache := g.CloneGroupCache()
	cache.AddKeys(keys...)

	return cache
}

func (g *GroupCacheRedis) GetGroupOwner(ctx context.Context, groupID string) (*model.GroupMember, error) {
	members, err := g.GetGroupRoleLevelMemberInfo(ctx, groupID, constant.GroupOwner)
	if err != nil {
		return nil, err
	}
	if len(members) == 0 {
		return nil, errs.ErrRecordNotFound.WrapMsg(fmt.Sprintf("group %s owner not found", groupID))
	}
	return members[0], nil
}

func (g *GroupCacheRedis) GetGroupsOwner(ctx context.Context, groupIDs []string) ([]*model.GroupMember, error) {
	members := make([]*model.GroupMember, 0, len(groupIDs))
	for _, groupID := range groupIDs {
		items, err := g.GetGroupRoleLevelMemberInfo(ctx, groupID, constant.GroupOwner)
		if err != nil {
			return nil, err
		}
		if len(items) > 0 {
			members = append(members, items[0])
		}
	}
	return members, nil
}

func (g *GroupCacheRedis) GetGroupRoleLevelMemberIDs(ctx context.Context, groupID string, roleLevel int32) ([]string, error) {
	return getCache(ctx, g.rcClient, g.getGroupRoleLevelMemberIDsKey(groupID, roleLevel), g.expireTime, func(ctx context.Context) ([]string, error) {
		return g.groupMemberDB.FindRoleLevelUserIDs(ctx, groupID, roleLevel)
	})
}

func (g *GroupCacheRedis) GetGroupRoleLevelMemberInfo(ctx context.Context, groupID string, roleLevel int32) ([]*model.GroupMember, error) {
	userIDs, err := g.GetGroupRoleLevelMemberIDs(ctx, groupID, roleLevel)
	if err != nil {
		return nil, err
	}
	return g.GetGroupMembersInfo(ctx, groupID, userIDs)
}

func (g *GroupCacheRedis) GetGroupRolesLevelMemberInfo(ctx context.Context, groupID string, roleLevels []int32) ([]*model.GroupMember, error) {
	var userIDs []string
	for _, roleLevel := range roleLevels {
		ids, err := g.GetGroupRoleLevelMemberIDs(ctx, groupID, roleLevel)
		if err != nil {
			return nil, err
		}
		userIDs = append(userIDs, ids...)
	}
	return g.GetGroupMembersInfo(ctx, groupID, userIDs)
}

func (g *GroupCacheRedis) FindGroupMemberUser(ctx context.Context, groupIDs []string, userID string) ([]*model.GroupMember, error) {
	if len(groupIDs) == 0 {
		var err error
		groupIDs, err = g.GetJoinedGroupIDs(ctx, userID)
		if err != nil {
			return nil, err
		}
	}
	return batchGetCache2(ctx, g.rcClient, g.expireTime, groupIDs, func(groupID string) string {
		return g.getGroupMemberInfoKey(groupID, userID)
	}, func(member *model.GroupMember) string {
		return member.GroupID
	}, func(ctx context.Context, groupIDs []string) ([]*model.GroupMember, error) {
		return g.groupMemberDB.FindInGroup(ctx, userID, groupIDs)
	})
}

func (g *GroupCacheRedis) DelMaxGroupMemberVersion(groupIDs ...string) cache.GroupCache {
	keys := make([]string, 0, len(groupIDs))
	for _, groupID := range groupIDs {
		keys = append(keys, g.getGroupMemberMaxVersionKey(groupID))
	}
	cache := g.CloneGroupCache()
	cache.AddKeys(keys...)
	return cache
}

func (g *GroupCacheRedis) DelMaxJoinGroupVersion(userIDs ...string) cache.GroupCache {
	keys := make([]string, 0, len(userIDs))
	for _, userID := range userIDs {
		keys = append(keys, g.getJoinGroupMaxVersionKey(userID))
	}
	cache := g.CloneGroupCache()
	cache.AddKeys(keys...)
	return cache
}

func (g *GroupCacheRedis) FindMaxGroupMemberVersion(ctx context.Context, groupID string) (*model.VersionLog, error) {
	return getCache(ctx, g.rcClient, g.getGroupMemberMaxVersionKey(groupID), g.expireTime, func(ctx context.Context) (*model.VersionLog, error) {
		return g.groupMemberDB.FindMemberIncrVersion(ctx, groupID, 0, 0)
	})
}

func (g *GroupCacheRedis) BatchFindMaxGroupMemberVersion(ctx context.Context, groupIDs []string) ([]*model.VersionLog, error) {
	return batchGetCache2(ctx, g.rcClient, g.expireTime, groupIDs,
		func(groupID string) string {
			return g.getGroupMemberMaxVersionKey(groupID)
		}, func(versionLog *model.VersionLog) string {
			return versionLog.DID
		}, func(ctx context.Context, groupIDs []string) ([]*model.VersionLog, error) {
			// create two slices with len is groupIDs, just need 0
			versions := make([]uint, len(groupIDs))
			limits := make([]int, len(groupIDs))

			return g.groupMemberDB.BatchFindMemberIncrVersion(ctx, groupIDs, versions, limits)
		})
}

func (g *GroupCacheRedis) FindMaxJoinGroupVersion(ctx context.Context, userID string) (*model.VersionLog, error) {
	return getCache(ctx, g.rcClient, g.getJoinGroupMaxVersionKey(userID), g.expireTime, func(ctx context.Context) (*model.VersionLog, error) {
		return g.groupMemberDB.FindJoinIncrVersion(ctx, userID, 0, 0)
	})
}
