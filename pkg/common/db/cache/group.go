// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/dtm-labs/rockscache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/cachekey"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	relationtb "github.com/openimsdk/open-im-server/v3/pkg/common/db/table/relation"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/datautil"
	"github.com/redis/go-redis/v9"
)

const (
	groupExpireTime = time.Second * 60 * 60 * 12
)

type GroupHash interface {
	GetGroupHash(ctx context.Context, groupID string) (uint64, error)
}

type GroupCache interface {
	metaCache
	NewCache() GroupCache
	GetGroupsInfo(ctx context.Context, groupIDs []string) (groups []*relationtb.GroupModel, err error)
	GetGroupInfo(ctx context.Context, groupID string) (group *relationtb.GroupModel, err error)
	DelGroupsInfo(groupIDs ...string) GroupCache

	GetGroupMembersHash(ctx context.Context, groupID string) (hashCode uint64, err error)
	GetGroupMemberHashMap(ctx context.Context, groupIDs []string) (map[string]*relationtb.GroupSimpleUserID, error)
	DelGroupMembersHash(groupID string) GroupCache

	GetGroupMemberIDs(ctx context.Context, groupID string) (groupMemberIDs []string, err error)
	GetGroupsMemberIDs(ctx context.Context, groupIDs []string) (groupMemberIDs map[string][]string, err error)

	DelGroupMemberIDs(groupID string) GroupCache

	GetJoinedGroupIDs(ctx context.Context, userID string) (joinedGroupIDs []string, err error)
	DelJoinedGroupID(userID ...string) GroupCache

	GetGroupMemberInfo(ctx context.Context, groupID, userID string) (groupMember *relationtb.GroupMemberModel, err error)
	GetGroupMembersInfo(ctx context.Context, groupID string, userID []string) (groupMembers []*relationtb.GroupMemberModel, err error)
	GetAllGroupMembersInfo(ctx context.Context, groupID string) (groupMembers []*relationtb.GroupMemberModel, err error)
	GetGroupMembersPage(ctx context.Context, groupID string, userID []string, showNumber, pageNumber int32) (total uint32, groupMembers []*relationtb.GroupMemberModel, err error)
	FindGroupMemberUser(ctx context.Context, groupIDs []string, userID string) ([]*relationtb.GroupMemberModel, error)

	GetGroupRoleLevelMemberIDs(ctx context.Context, groupID string, roleLevel int32) ([]string, error)
	GetGroupOwner(ctx context.Context, groupID string) (*relationtb.GroupMemberModel, error)
	GetGroupsOwner(ctx context.Context, groupIDs []string) ([]*relationtb.GroupMemberModel, error)
	DelGroupRoleLevel(groupID string, roleLevel []int32) GroupCache
	DelGroupAllRoleLevel(groupID string) GroupCache
	DelGroupMembersInfo(groupID string, userID ...string) GroupCache
	GetGroupRoleLevelMemberInfo(ctx context.Context, groupID string, roleLevel int32) ([]*relationtb.GroupMemberModel, error)
	GetGroupRolesLevelMemberInfo(ctx context.Context, groupID string, roleLevels []int32) ([]*relationtb.GroupMemberModel, error)
	GetGroupMemberNum(ctx context.Context, groupID string) (memberNum int64, err error)
	DelGroupsMemberNum(groupID ...string) GroupCache
}

type GroupCacheRedis struct {
	metaCache
	groupDB        relationtb.GroupModelInterface
	groupMemberDB  relationtb.GroupMemberModelInterface
	groupRequestDB relationtb.GroupRequestModelInterface
	expireTime     time.Duration
	rcClient       *rockscache.Client
	groupHash      GroupHash
}

func NewGroupCacheRedis(
	rdb redis.UniversalClient,
	localCache *config.LocalCache,
	groupDB relationtb.GroupModelInterface,
	groupMemberDB relationtb.GroupMemberModelInterface,
	groupRequestDB relationtb.GroupRequestModelInterface,
	hashCode GroupHash,
	opts rockscache.Options,
) GroupCache {
	rcClient := rockscache.NewClient(rdb, opts)
	mc := NewMetaCacheRedis(rcClient)
	g := localCache.Group
	mc.SetTopic(g.Topic)
	log.ZDebug(context.Background(), "group local cache init", "Topic", g.Topic, "SlotNum", g.SlotNum, "SlotSize", g.SlotSize, "enable", g.Enable())
	mc.SetRawRedisClient(rdb)
	return &GroupCacheRedis{
		rcClient: rcClient, expireTime: groupExpireTime,
		groupDB: groupDB, groupMemberDB: groupMemberDB, groupRequestDB: groupRequestDB,
		groupHash: hashCode,
		metaCache: mc,
	}
}

func (g *GroupCacheRedis) NewCache() GroupCache {
	return &GroupCacheRedis{
		rcClient:       g.rcClient,
		expireTime:     g.expireTime,
		groupDB:        g.groupDB,
		groupMemberDB:  g.groupMemberDB,
		groupRequestDB: g.groupRequestDB,
		metaCache:      g.Copy(),
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

func (g *GroupCacheRedis) GetGroupIndex(group *relationtb.GroupModel, keys []string) (int, error) {
	key := g.getGroupInfoKey(group.GroupID)
	for i, _key := range keys {
		if _key == key {
			return i, nil
		}
	}

	return 0, errIndex
}

func (g *GroupCacheRedis) GetGroupMemberIndex(groupMember *relationtb.GroupMemberModel, keys []string) (int, error) {
	key := g.getGroupMemberInfoKey(groupMember.GroupID, groupMember.UserID)
	for i, _key := range keys {
		if _key == key {
			return i, nil
		}
	}

	return 0, errIndex
}

func (g *GroupCacheRedis) GetGroupsInfo(ctx context.Context, groupIDs []string) (groups []*relationtb.GroupModel, err error) {
	return batchGetCache2(ctx, g.rcClient, g.expireTime, groupIDs, func(groupID string) string {
		return g.getGroupInfoKey(groupID)
	}, func(ctx context.Context, groupID string) (*relationtb.GroupModel, error) {
		return g.groupDB.Take(ctx, groupID)
	})
}

func (g *GroupCacheRedis) GetGroupInfo(ctx context.Context, groupID string) (group *relationtb.GroupModel, err error) {
	return getCache(ctx, g.rcClient, g.getGroupInfoKey(groupID), g.expireTime, func(ctx context.Context) (*relationtb.GroupModel, error) {
		return g.groupDB.Take(ctx, groupID)
	})
}

func (g *GroupCacheRedis) DelGroupsInfo(groupIDs ...string) GroupCache {
	newGroupCache := g.NewCache()
	keys := make([]string, 0, len(groupIDs))
	for _, groupID := range groupIDs {
		keys = append(keys, g.getGroupInfoKey(groupID))
	}
	newGroupCache.AddKeys(keys...)

	return newGroupCache
}

func (g *GroupCacheRedis) DelGroupsOwner(groupIDs ...string) GroupCache {
	newGroupCache := g.NewCache()
	keys := make([]string, 0, len(groupIDs))
	for _, groupID := range groupIDs {
		keys = append(keys, g.getGroupRoleLevelMemberIDsKey(groupID, constant.GroupOwner))
	}
	newGroupCache.AddKeys(keys...)

	return newGroupCache
}

func (g *GroupCacheRedis) DelGroupRoleLevel(groupID string, roleLevels []int32) GroupCache {
	newGroupCache := g.NewCache()
	keys := make([]string, 0, len(roleLevels))
	for _, roleLevel := range roleLevels {
		keys = append(keys, g.getGroupRoleLevelMemberIDsKey(groupID, roleLevel))
	}
	newGroupCache.AddKeys(keys...)
	return newGroupCache
}

func (g *GroupCacheRedis) DelGroupAllRoleLevel(groupID string) GroupCache {
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

func (g *GroupCacheRedis) GetGroupMemberHashMap(ctx context.Context, groupIDs []string) (map[string]*relationtb.GroupSimpleUserID, error) {
	if g.groupHash == nil {
		return nil, errs.ErrInternalServer.WrapMsg("group hash is nil")
	}
	res := make(map[string]*relationtb.GroupSimpleUserID)
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
		res[groupID] = &relationtb.GroupSimpleUserID{Hash: hash, MemberNum: uint32(num)}
	}

	return res, nil
}

func (g *GroupCacheRedis) DelGroupMembersHash(groupID string) GroupCache {
	cache := g.NewCache()
	cache.AddKeys(g.getGroupMembersHashKey(groupID))

	return cache
}

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
	keys := make([]string, 0, len(userIDs))
	for _, userID := range userIDs {
		keys = append(keys, g.getJoinedGroupsKey(userID))
	}
	cache := g.NewCache()
	cache.AddKeys(keys...)

	return cache
}

func (g *GroupCacheRedis) GetGroupMemberInfo(ctx context.Context, groupID, userID string) (groupMember *relationtb.GroupMemberModel, err error) {
	return getCache(ctx, g.rcClient, g.getGroupMemberInfoKey(groupID, userID), g.expireTime, func(ctx context.Context) (*relationtb.GroupMemberModel, error) {
		return g.groupMemberDB.Take(ctx, groupID, userID)
	})
}

func (g *GroupCacheRedis) GetGroupMembersInfo(ctx context.Context, groupID string, userIDs []string) ([]*relationtb.GroupMemberModel, error) {
	return batchGetCache2(ctx, g.rcClient, g.expireTime, userIDs, func(userID string) string {
		return g.getGroupMemberInfoKey(groupID, userID)
	}, func(ctx context.Context, userID string) (*relationtb.GroupMemberModel, error) {
		return g.groupMemberDB.Take(ctx, groupID, userID)
	})
}

func (g *GroupCacheRedis) GetGroupMembersPage(
	ctx context.Context,
	groupID string,
	userIDs []string,
	showNumber, pageNumber int32,
) (total uint32, groupMembers []*relationtb.GroupMemberModel, err error) {
	groupMemberIDs, err := g.GetGroupMemberIDs(ctx, groupID)
	if err != nil {
		return 0, nil, err
	}
	if userIDs != nil {
		userIDs = datautil.BothExist(userIDs, groupMemberIDs)
	} else {
		userIDs = groupMemberIDs
	}
	groupMembers, err = g.GetGroupMembersInfo(ctx, groupID, datautil.Paginate(userIDs, int(showNumber), int(showNumber)))

	return uint32(len(userIDs)), groupMembers, err
}

func (g *GroupCacheRedis) GetAllGroupMembersInfo(ctx context.Context, groupID string) (groupMembers []*relationtb.GroupMemberModel, err error) {
	groupMemberIDs, err := g.GetGroupMemberIDs(ctx, groupID)
	if err != nil {
		return nil, err
	}

	return g.GetGroupMembersInfo(ctx, groupID, groupMemberIDs)
}

func (g *GroupCacheRedis) GetAllGroupMemberInfo(ctx context.Context, groupID string) ([]*relationtb.GroupMemberModel, error) {
	groupMemberIDs, err := g.GetGroupMemberIDs(ctx, groupID)
	if err != nil {
		return nil, err
	}
	return g.GetGroupMembersInfo(ctx, groupID, groupMemberIDs)
}

func (g *GroupCacheRedis) DelGroupMembersInfo(groupID string, userIDs ...string) GroupCache {
	keys := make([]string, 0, len(userIDs))
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
	keys := make([]string, 0, len(groupID))
	for _, groupID := range groupID {
		keys = append(keys, g.getGroupMemberNumKey(groupID))
	}
	cache := g.NewCache()
	cache.AddKeys(keys...)

	return cache
}

func (g *GroupCacheRedis) GetGroupOwner(ctx context.Context, groupID string) (*relationtb.GroupMemberModel, error) {
	members, err := g.GetGroupRoleLevelMemberInfo(ctx, groupID, constant.GroupOwner)
	if err != nil {
		return nil, err
	}
	if len(members) == 0 {
		return nil, errs.ErrRecordNotFound.WrapMsg(fmt.Sprintf("group %s owner not found", groupID))
	}
	return members[0], nil
}

func (g *GroupCacheRedis) GetGroupsOwner(ctx context.Context, groupIDs []string) ([]*relationtb.GroupMemberModel, error) {
	members := make([]*relationtb.GroupMemberModel, 0, len(groupIDs))
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

func (g *GroupCacheRedis) GetGroupRoleLevelMemberInfo(ctx context.Context, groupID string, roleLevel int32) ([]*relationtb.GroupMemberModel, error) {
	userIDs, err := g.GetGroupRoleLevelMemberIDs(ctx, groupID, roleLevel)
	if err != nil {
		return nil, err
	}
	return g.GetGroupMembersInfo(ctx, groupID, userIDs)
}

func (g *GroupCacheRedis) GetGroupRolesLevelMemberInfo(ctx context.Context, groupID string, roleLevels []int32) ([]*relationtb.GroupMemberModel, error) {
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

func (g *GroupCacheRedis) FindGroupMemberUser(ctx context.Context, groupIDs []string, userID string) (_ []*relationtb.GroupMemberModel, err error) {
	if len(groupIDs) == 0 {
		groupIDs, err = g.GetJoinedGroupIDs(ctx, userID)
		if err != nil {
			return nil, err
		}
	}
	return batchGetCache2(ctx, g.rcClient, g.expireTime, groupIDs, func(groupID string) string {
		return g.getGroupMemberInfoKey(groupID, userID)
	}, func(ctx context.Context, groupID string) (*relationtb.GroupMemberModel, error) {
		return g.groupMemberDB.Take(ctx, groupID, userID)
	})
}
