package cache

import (
	"context"
	"math/big"
	"strings"
	"time"

	relationTb "github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	unrelationTb "github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/unrelation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"github.com/dtm-labs/rockscache"
	"github.com/redis/go-redis/v9"
)

const (
	groupExpireTime        = time.Second * 60 * 60 * 12
	groupInfoKey           = "GROUP_INFO:"
	groupMemberIDsKey      = "GROUP_MEMBER_IDS:"
	groupMembersHashKey    = "GROUP_MEMBERS_HASH:"
	groupMemberInfoKey     = "GROUP_MEMBER_INFO:"
	joinedSuperGroupsKey   = "JOIN_SUPER_GROUPS:"
	SuperGroupMemberIDsKey = "SUPER_GROUP_MEMBER_IDS:"
	joinedGroupsKey        = "JOIN_GROUPS_KEY:"
	groupMemberNumKey      = "GROUP_MEMBER_NUM_CACHE:"
)

type GroupCache interface {
	metaCache
	NewCache() GroupCache
	GetGroupsInfo(ctx context.Context, groupIDs []string) (groups []*relationTb.GroupModel, err error)
	GetGroupInfo(ctx context.Context, groupID string) (group *relationTb.GroupModel, err error)
	DelGroupsInfo(groupIDs ...string) GroupCache

	GetJoinedSuperGroupIDs(ctx context.Context, userID string) (joinedSuperGroupIDs []string, err error)
	DelJoinedSuperGroupIDs(userIDs ...string) GroupCache
	GetSuperGroupMemberIDs(ctx context.Context, groupIDs ...string) (models []*unrelationTb.SuperGroupModel, err error)
	DelSuperGroupMemberIDs(groupIDs ...string) GroupCache

	GetGroupMembersHash(ctx context.Context, groupID string) (hashCode uint64, err error)
	GetGroupMemberHashMap(ctx context.Context, groupIDs []string) (map[string]*relationTb.GroupSimpleUserID, error)
	DelGroupMembersHash(groupID string) GroupCache

	GetGroupMemberIDs(ctx context.Context, groupID string) (groupMemberIDs []string, err error)
	GetGroupsMemberIDs(ctx context.Context, groupIDs []string) (groupMemberIDs map[string][]string, err error)

	DelGroupMemberIDs(groupID string) GroupCache

	GetJoinedGroupIDs(ctx context.Context, userID string) (joinedGroupIDs []string, err error)
	DelJoinedGroupID(userID ...string) GroupCache

	GetGroupMemberInfo(ctx context.Context, groupID, userID string) (groupMember *relationTb.GroupMemberModel, err error)
	GetGroupMembersInfo(ctx context.Context, groupID string, userID []string) (groupMembers []*relationTb.GroupMemberModel, err error)
	GetAllGroupMembersInfo(ctx context.Context, groupID string) (groupMembers []*relationTb.GroupMemberModel, err error)
	GetGroupMembersPage(ctx context.Context, groupID string, userID []string, showNumber, pageNumber int32) (total uint32, groupMembers []*relationTb.GroupMemberModel, err error)

	DelGroupMembersInfo(groupID string, userID ...string) GroupCache

	GetGroupMemberNum(ctx context.Context, groupID string) (memberNum int64, err error)
	DelGroupsMemberNum(groupID ...string) GroupCache
}

type GroupCacheRedis struct {
	metaCache
	groupDB        relationTb.GroupModelInterface
	groupMemberDB  relationTb.GroupMemberModelInterface
	groupRequestDB relationTb.GroupRequestModelInterface
	mongoDB        unrelationTb.SuperGroupModelInterface
	expireTime     time.Duration
	rcClient       *rockscache.Client
}

func NewGroupCacheRedis(rdb redis.UniversalClient, groupDB relationTb.GroupModelInterface, groupMemberDB relationTb.GroupMemberModelInterface, groupRequestDB relationTb.GroupRequestModelInterface, mongoClient unrelationTb.SuperGroupModelInterface, opts rockscache.Options) GroupCache {
	rcClient := rockscache.NewClient(rdb, opts)
	return &GroupCacheRedis{rcClient: rcClient, expireTime: groupExpireTime,
		groupDB: groupDB, groupMemberDB: groupMemberDB, groupRequestDB: groupRequestDB,
		mongoDB: mongoClient, metaCache: NewMetaCacheRedis(rcClient),
	}
}

func (g *GroupCacheRedis) NewCache() GroupCache {
	return &GroupCacheRedis{rcClient: g.rcClient, expireTime: g.expireTime, groupDB: g.groupDB, groupMemberDB: g.groupMemberDB, groupRequestDB: g.groupRequestDB, mongoDB: g.mongoDB, metaCache: NewMetaCacheRedis(g.rcClient, g.metaCache.GetPreDelKeys()...)}
}

func (g *GroupCacheRedis) getGroupInfoKey(groupID string) string {
	return groupInfoKey + groupID
}

func (g *GroupCacheRedis) getJoinedSuperGroupsIDKey(userID string) string {
	return joinedSuperGroupsKey + userID
}

func (g *GroupCacheRedis) getJoinedGroupsKey(userID string) string {
	return joinedGroupsKey + userID
}

func (g *GroupCacheRedis) getSuperGroupMemberIDsKey(groupID string) string {
	return SuperGroupMemberIDsKey + groupID
}

func (g *GroupCacheRedis) getGroupMembersHashKey(groupID string) string {
	return groupMembersHashKey + groupID
}

func (g *GroupCacheRedis) getGroupMemberIDsKey(groupID string) string {
	return groupMemberIDsKey + groupID
}

func (g *GroupCacheRedis) getGroupMemberInfoKey(groupID, userID string) string {
	return groupMemberInfoKey + groupID + "-" + userID
}

func (g *GroupCacheRedis) getGroupMemberNumKey(groupID string) string {
	return groupMemberNumKey + groupID
}

func (g *GroupCacheRedis) GetGroupIndex(group *relationTb.GroupModel, keys []string) (int, error) {
	key := g.getGroupInfoKey(group.GroupID)
	for i, _key := range keys {
		if _key == key {
			return i, nil
		}
	}
	return 0, errIndex
}

func (g *GroupCacheRedis) GetGroupMemberIndex(groupMember *relationTb.GroupMemberModel, keys []string) (int, error) {
	key := g.getGroupMemberInfoKey(groupMember.GroupID, groupMember.UserID)
	for i, _key := range keys {
		if _key == key {
			return i, nil
		}
	}
	return 0, errIndex
}

// / groupInfo
func (g *GroupCacheRedis) GetGroupsInfo(ctx context.Context, groupIDs []string) (groups []*relationTb.GroupModel, err error) {
	var keys []string
	for _, group := range groupIDs {
		keys = append(keys, g.getGroupInfoKey(group))
	}
	return batchGetCache(ctx, g.rcClient, keys, g.expireTime, g.GetGroupIndex, func(ctx context.Context) ([]*relationTb.GroupModel, error) {
		return g.groupDB.Find(ctx, groupIDs)
	})
}

func (g *GroupCacheRedis) GetGroupInfo(ctx context.Context, groupID string) (group *relationTb.GroupModel, err error) {
	return getCache(ctx, g.rcClient, g.getGroupInfoKey(groupID), g.expireTime, func(ctx context.Context) (*relationTb.GroupModel, error) {
		return g.groupDB.Take(ctx, groupID)
	})
}

func (g *GroupCacheRedis) DelGroupsInfo(groupIDs ...string) GroupCache {
	new := g.NewCache()
	var keys []string
	for _, groupID := range groupIDs {
		keys = append(keys, g.getGroupInfoKey(groupID))
	}
	new.AddKeys(keys...)
	return new
}

func (g *GroupCacheRedis) GetJoinedSuperGroupIDs(ctx context.Context, userID string) (joinedSuperGroupIDs []string, err error) {
	return getCache(ctx, g.rcClient, g.getJoinedSuperGroupsIDKey(userID), g.expireTime, func(ctx context.Context) ([]string, error) {
		userGroup, err := g.mongoDB.GetSuperGroupByUserID(ctx, userID)
		if err != nil {
			return nil, err
		}
		return userGroup.GroupIDs, nil
	})
}

func (g *GroupCacheRedis) GetSuperGroupMemberIDs(ctx context.Context, groupIDs ...string) (models []*unrelationTb.SuperGroupModel, err error) {
	var keys []string
	for _, group := range groupIDs {
		keys = append(keys, g.getSuperGroupMemberIDsKey(group))
	}
	return batchGetCache(ctx, g.rcClient, keys, g.expireTime, func(model *unrelationTb.SuperGroupModel, keys []string) (int, error) {
		for i, key := range keys {
			if g.getSuperGroupMemberIDsKey(model.GroupID) == key {
				return i, nil
			}
		}
		return 0, errIndex
	}, func(ctx context.Context) ([]*unrelationTb.SuperGroupModel, error) {
		return g.mongoDB.FindSuperGroup(ctx, groupIDs)
	})
}

// userJoinSuperGroup
func (g *GroupCacheRedis) DelJoinedSuperGroupIDs(userIDs ...string) GroupCache {
	new := g.NewCache()
	var keys []string
	for _, userID := range userIDs {
		keys = append(keys, g.getJoinedSuperGroupsIDKey(userID))
	}
	new.AddKeys(keys...)
	return new
}

func (g *GroupCacheRedis) DelSuperGroupMemberIDs(groupIDs ...string) GroupCache {
	new := g.NewCache()
	var keys []string
	for _, groupID := range groupIDs {
		keys = append(keys, g.getSuperGroupMemberIDsKey(groupID))
	}
	new.AddKeys(keys...)
	return new
}

// groupMembersHash
func (g *GroupCacheRedis) GetGroupMembersHash(ctx context.Context, groupID string) (hashCode uint64, err error) {
	return getCache(ctx, g.rcClient, g.getGroupMembersHashKey(groupID), g.expireTime, func(ctx context.Context) (uint64, error) {
		userIDs, err := g.GetGroupMemberIDs(ctx, groupID)
		if err != nil {
			return 0, err
		}
		utils.Sort(userIDs, true)
		bi := big.NewInt(0)
		bi.SetString(utils.Md5(strings.Join(userIDs, ";"))[0:8], 16)
		return bi.Uint64(), nil
	})
}

func (g *GroupCacheRedis) GetGroupMemberHashMap(ctx context.Context, groupIDs []string) (map[string]*relationTb.GroupSimpleUserID, error) {
	res := make(map[string]*relationTb.GroupSimpleUserID)
	for _, groupID := range groupIDs {
		hash, err := g.GetGroupMembersHash(ctx, groupID)
		if err != nil {
			return nil, err
		}
		num, err := g.GetGroupMemberNum(ctx, groupID)
		if err != nil {
			return nil, err
		}
		res[groupID] = &relationTb.GroupSimpleUserID{Hash: hash, MemberNum: uint32(num)}
	}
	return res, nil
}

func (g *GroupCacheRedis) DelGroupMembersHash(groupID string) GroupCache {
	cache := g.NewCache()
	cache.AddKeys(g.getGroupMembersHashKey(groupID))
	return cache
}

// groupMemberIDs
func (g *GroupCacheRedis) GetGroupMemberIDs(ctx context.Context, groupID string) (groupMemberIDs []string, err error) {
	return getCache(ctx, g.rcClient, g.getGroupMemberIDsKey(groupID), g.expireTime, func(ctx context.Context) ([]string, error) {
		return g.groupMemberDB.FindMemberUserID(ctx, groupID)
	})
}

func (g *GroupCacheRedis) GetGroupsMemberIDs(ctx context.Context, groupIDs []string) (map[string][]string, error) {
	m := make(map[string][]string)
	for _, groupID := range groupIDs {
		userIDs, err := g.GetGroupMemberIDs(ctx, groupID)
		if err != nil {
			return nil, err
		}
		m[groupID] = userIDs
	}
	return m, nil
}

func (g *GroupCacheRedis) DelGroupMemberIDs(groupID string) GroupCache {
	cache := g.NewCache()
	cache.AddKeys(g.getGroupMemberIDsKey(groupID))
	return cache
}

func (g *GroupCacheRedis) GetJoinedGroupIDs(ctx context.Context, userID string) (joinedGroupIDs []string, err error) {
	return getCache(ctx, g.rcClient, g.getJoinedGroupsKey(userID), g.expireTime, func(ctx context.Context) ([]string, error) {
		return g.groupMemberDB.FindUserJoinedGroupID(ctx, userID)
	})
}

func (g *GroupCacheRedis) DelJoinedGroupID(userIDs ...string) GroupCache {
	var keys []string
	for _, userID := range userIDs {
		keys = append(keys, g.getJoinedGroupsKey(userID))
	}
	cache := g.NewCache()
	cache.AddKeys(keys...)
	return cache
}

func (g *GroupCacheRedis) GetGroupMemberInfo(ctx context.Context, groupID, userID string) (groupMember *relationTb.GroupMemberModel, err error) {
	return getCache(ctx, g.rcClient, g.getGroupMemberInfoKey(groupID, userID), g.expireTime, func(ctx context.Context) (*relationTb.GroupMemberModel, error) {
		return g.groupMemberDB.Take(ctx, groupID, userID)
	})
}

func (g *GroupCacheRedis) GetGroupMembersInfo(ctx context.Context, groupID string, userIDs []string) ([]*relationTb.GroupMemberModel, error) {
	var keys []string
	for _, userID := range userIDs {
		keys = append(keys, g.getGroupMemberInfoKey(groupID, userID))
	}
	return batchGetCache(ctx, g.rcClient, keys, g.expireTime, g.GetGroupMemberIndex, func(ctx context.Context) ([]*relationTb.GroupMemberModel, error) {
		return g.groupMemberDB.Find(ctx, []string{groupID}, userIDs, nil)
	})
}

func (g *GroupCacheRedis) GetGroupMembersPage(ctx context.Context, groupID string, userIDs []string, showNumber, pageNumber int32) (total uint32, groupMembers []*relationTb.GroupMemberModel, err error) {
	groupMemberIDs, err := g.GetGroupMemberIDs(ctx, groupID)
	if err != nil {
		return 0, nil, err
	}
	if userIDs != nil {
		userIDs = utils.BothExist(userIDs, groupMemberIDs)
	} else {
		userIDs = groupMemberIDs
	}
	groupMembers, err = g.GetGroupMembersInfo(ctx, groupID, utils.Paginate(userIDs, int(showNumber), int(showNumber)))
	return uint32(len(userIDs)), groupMembers, err
}

func (g *GroupCacheRedis) GetAllGroupMembersInfo(ctx context.Context, groupID string) (groupMembers []*relationTb.GroupMemberModel, err error) {
	groupMemberIDs, err := g.GetGroupMemberIDs(ctx, groupID)
	if err != nil {
		return nil, err
	}
	return g.GetGroupMembersInfo(ctx, groupID, groupMemberIDs)
}

func (g *GroupCacheRedis) GetAllGroupMemberInfo(ctx context.Context, groupID string) ([]*relationTb.GroupMemberModel, error) {
	groupMemberIDs, err := g.GetGroupMemberIDs(ctx, groupID)
	if err != nil {
		return nil, err
	}
	var keys []string
	for _, groupMemberID := range groupMemberIDs {
		keys = append(keys, g.getGroupMemberInfoKey(groupID, groupMemberID))
	}
	return batchGetCache(ctx, g.rcClient, keys, g.expireTime, g.GetGroupMemberIndex, func(ctx context.Context) ([]*relationTb.GroupMemberModel, error) {
		return g.groupMemberDB.Find(ctx, []string{groupID}, groupMemberIDs, nil)
	})
}

func (g *GroupCacheRedis) DelGroupMembersInfo(groupID string, userIDs ...string) GroupCache {
	var keys []string
	for _, userID := range userIDs {
		keys = append(keys, g.getGroupMemberInfoKey(groupID, userID))
	}
	cache := g.NewCache()
	cache.AddKeys(keys...)
	return cache
}

func (g *GroupCacheRedis) GetGroupMemberNum(ctx context.Context, groupID string) (memberNum int64, err error) {
	return getCache(ctx, g.rcClient, g.getGroupMemberNumKey(groupID), g.expireTime, func(ctx context.Context) (int64, error) {
		return g.groupMemberDB.TakeGroupMemberNum(ctx, groupID)
	})
}

func (g *GroupCacheRedis) DelGroupsMemberNum(groupID ...string) GroupCache {
	var keys []string
	for _, groupID := range groupID {
		keys = append(keys, g.getGroupMemberNumKey(groupID))
	}
	cache := g.NewCache()
	cache.AddKeys(keys...)
	return cache
}
