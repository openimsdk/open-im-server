package cache

import (
	relationTb "OpenIM/pkg/common/db/table/relation"
	unrelation2 "OpenIM/pkg/common/db/table/unrelation"
	"OpenIM/pkg/common/tracelog"
	"OpenIM/pkg/utils"
	"context"
	"github.com/dtm-labs/rockscache"
	"github.com/go-redis/redis/v8"
	"math/big"
	"strings"
	"time"
)

const (
	groupExpireTime      = time.Second * 60 * 60 * 12
	groupInfoKey         = "GROUP_INFO:"
	groupMemberIDsKey    = "GROUP_MEMBER_IDS:"
	groupMembersHashKey  = "GROUP_MEMBERS_HASH:"
	groupMemberInfoKey   = "GROUP_MEMBER_INFO:"
	joinedSuperGroupsKey = "JOIN_SUPER_GROUPS:"
	joinedGroupsKey      = "JOIN_GROUPS_KEY:"
	groupMemberNumKey    = "GROUP_MEMBER_NUM_CACHE:"
)

type GroupCache interface {
	GetGroupsInfo(ctx context.Context, groupIDs []string, fn func(ctx context.Context, groupIDs []string) (groups []*relationTb.GroupModel, err error)) (groups []*relationTb.GroupModel, err error)
	DelGroupsInfo(ctx context.Context, groupID string) (err error)
	GetGroupInfo(ctx context.Context, groupID string, fn func(ctx context.Context, groupID string) (group *relationTb.GroupModel, err error)) (group *relationTb.GroupModel, err error)
	DelGroupInfo(ctx context.Context, groupID string) (err error)

	BatchDelJoinedSuperGroupIDs(ctx context.Context, userIDs []string, fn func(ctx context.Context, userIDs []string) error) (err error)

	GetJoinedSuperGroupIDs(ctx context.Context, userID string, fn func(ctx context.Context, userID string) (joinedSuperGroupIDs []string, err error)) (joinedSuperGroupIDs []string, err error)
	DelJoinedSuperGroupIDs(ctx context.Context, userID string) (err error)

	GetGroupMembersHash(ctx context.Context, groupID string, fn func(ctx context.Context, groupID string) (hashCodeUint64 uint64, err error)) (hashCodeUint64 uint64, err error)
	DelGroupMembersHash(ctx context.Context, groupID string) (err error)

	GetGroupMemberIDs(ctx context.Context, groupID string, fn func(ctx context.Context, groupID string) (groupMemberIDs []string, err error)) (groupMemberIDs []string, err error)
	DelGroupMemberIDs(ctx context.Context, groupID string) error

	GetJoinedGroupIDs(ctx context.Context, userID string, fn func(ctx context.Context, userID string) (joinedGroupIDs []string, err error)) (joinedGroupIDs []string, err error)
	DelJoinedGroupIDs(ctx context.Context, userID string) (err error)

	GetGroupMemberInfo(ctx context.Context, groupID, userID string, fn func(ctx context.Context, groupID, userID string) (groupMember *relationTb.GroupMemberModel, err error)) (groupMember *relationTb.GroupMemberModel, err error)
	DelGroupMemberInfo(ctx context.Context, groupID, userID string) (err error)

	GetGroupMemberNum(ctx context.Context, groupID string, fn func(ctx context.Context, groupID string) (num int, err error)) (num int, err error)
	DelGroupMemberNum(ctx context.Context, groupID string) (err error)
}

type GroupCacheRedisInterface interface {
	GetGroupsInfo(ctx context.Context, groupIDs []string) (groups []*relationTb.GroupModel, err error)
	GetGroupInfo(ctx context.Context, groupID string) (group *relationTb.GroupModel, err error)
	BatchDelJoinedSuperGroupIDs(ctx context.Context, userIDs []string) (err error)
	DelJoinedSuperGroupIDs(ctx context.Context, userID string) (err error)
	GetJoinedSuperGroupIDs(ctx context.Context, userID string) (joinedSuperGroupIDs []string, err error)
	GetGroupMembersHash(ctx context.Context, groupID string) (hashCodeUint64 uint64, err error)
	GetGroupMemberHash1(ctx context.Context, groupIDs []string) (map[string]*relationTb.GroupSimpleUserID, error)
	DelGroupMembersHash(ctx context.Context, groupID string) (err error)
	GetGroupMemberIDs(ctx context.Context, groupID string) (groupMemberIDs []string, err error)
	DelGroupMemberIDs(ctx context.Context, groupID string) (err error)
	DelJoinedGroupID(ctx context.Context, userID string) (err error)
	GetGroupMemberInfo(ctx context.Context, groupID, userID string) (groupMember *relationTb.GroupMemberModel, err error)
	DelGroupMemberInfo(ctx context.Context, groupID, userID string) (err error)
	DelGroupMemberNum(ctx context.Context, groupID string) (err error)
	DelGroupInfo(ctx context.Context, groupID string) (err error)
	DelGroupsInfo(ctx context.Context, groupIDs []string) error
}

type GroupCacheRedis struct {
	group        relationTb.GroupModelInterface
	groupMember  relationTb.GroupMemberModelInterface
	groupRequest relationTb.GroupRequestModelInterface
	mongoDB      unrelation2.SuperGroupModelInterface
	expireTime   time.Duration
	redisClient  *RedisClient
	rcClient     *rockscache.Client
}

func NewGroupCacheRedis(rdb redis.UniversalClient, groupDB relationTb.GroupModelInterface, groupMemberDB relationTb.GroupMemberModelInterface, groupRequestDB relationTb.GroupRequestModelInterface, mongoClient unrelation2.SuperGroupModelInterface, opts rockscache.Options) GroupCacheRedisInterface {
	return &GroupCacheRedis{rcClient: rockscache.NewClient(rdb, opts), expireTime: groupExpireTime,
		group: groupDB, groupMember: groupMemberDB, groupRequest: groupRequestDB, redisClient: NewRedisClient(rdb),
		mongoDB: mongoClient,
	}
}

func (g *GroupCacheRedis) getRedisClient() *RedisClient {
	return g.redisClient
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

// / groupInfo
func (g *GroupCacheRedis) GetGroupsInfo(ctx context.Context, groupIDs []string) (groups []*relationTb.GroupModel, err error) {
	return GetCacheFor(ctx, groupIDs, func(ctx context.Context, groupID string) (*relationTb.GroupModel, error) {
		return g.GetGroupInfo(ctx, groupID)
	})
}

func (g *GroupCacheRedis) GetGroupInfo(ctx context.Context, groupID string) (group *relationTb.GroupModel, err error) {
	return GetCache(ctx, g.rcClient, g.getGroupInfoKey(groupID), g.expireTime, func(ctx context.Context) (*relationTb.GroupModel, error) {
		return g.group.Take(ctx, groupID)
	})
}

// userJoinSuperGroup
func (g *GroupCacheRedis) BatchDelJoinedSuperGroupIDs(ctx context.Context, userIDs []string) (err error) {
	for _, userID := range userIDs {
		if err := g.DelJoinedSuperGroupIDs(ctx, userID); err != nil {
			return err
		}
	}
	return nil
}

func (g *GroupCacheRedis) DelJoinedSuperGroupIDs(ctx context.Context, userID string) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "userID", userID)
	}()
	return g.rcClient.TagAsDeleted(g.getJoinedSuperGroupsIDKey(userID))
}

func (g *GroupCacheRedis) GetJoinedSuperGroupIDs(ctx context.Context, userID string) (joinedSuperGroupIDs []string, err error) {
	return GetCache(ctx, g.rcClient, g.getJoinedSuperGroupsIDKey(userID), g.expireTime, func(ctx context.Context) ([]string, error) {
		userGroup, err := g.mongoDB.GetSuperGroupByUserID(ctx, userID)
		if err != nil {
			return nil, err
		}
		return userGroup.GroupIDs, nil
	})
}

// groupMembersHash
func (g *GroupCacheRedis) GetGroupMembersHash(ctx context.Context, groupID string) (hashCodeUint64 uint64, err error) {
	return GetCache(ctx, g.rcClient, g.getGroupMembersHashKey(groupID), g.expireTime, func(ctx context.Context) (uint64, error) {
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

func (g *GroupCacheRedis) GetGroupMemberHash1(ctx context.Context, groupIDs []string) (map[string]*relationTb.GroupSimpleUserID, error) {
	// todo
	mapGroupUserIDs, err := g.groupMember.FindJoinUserID(ctx, groupIDs)
	if err != nil {
		return nil, err
	}
	res := make(map[string]*relationTb.GroupSimpleUserID)
	for _, groupID := range groupIDs {
		userIDs := mapGroupUserIDs[groupID]
		users := &relationTb.GroupSimpleUserID{}
		if len(userIDs) > 0 {
			utils.Sort(userIDs, true)
			bi := big.NewInt(0)
			bi.SetString(utils.Md5(strings.Join(userIDs, ";"))[0:8], 16)
			users.Hash = bi.Uint64()
		}
		res[groupID] = users
	}
	return res, nil
}

func (g *GroupCacheRedis) DelGroupMembersHash(ctx context.Context, groupID string) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupID", groupID)
	}()
	return g.rcClient.TagAsDeleted(g.getGroupMembersHashKey(groupID))
}

// groupMemberIDs
func (g *GroupCacheRedis) GetGroupMemberIDs(ctx context.Context, groupID string) (groupMemberIDs []string, err error) {
	return GetCache(ctx, g.rcClient, g.getGroupMemberIDsKey(groupID), g.expireTime, func(ctx context.Context) ([]string, error) {
		return g.groupMember.FindMemberUserID(ctx, groupID)
	})
}

func (g *GroupCacheRedis) DelGroupMemberIDs(ctx context.Context, groupID string) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupID", groupID)
	}()
	return g.rcClient.TagAsDeleted(g.getGroupMemberIDsKey(groupID))
}

//// JoinedGroups
//func (g *GroupCacheRedis) GetJoinedGroupIDs(ctx context.Context, userID string) (joinedGroupIDs []string, err error) {
//	getJoinedGroupIDList := func() (string, error) {
//		joinedGroupList, err := relation.GetJoinedGroupIDListByUserID(userID)
//		if err != nil {
//			return "", err
//		}
//		bytes, err := json.Marshal(joinedGroupList)
//		if err != nil {
//			return "", utils.Wrap(err, "")
//		}
//		return string(bytes), nil
//	}
//	defer func() {
//		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "userID", userID, "joinedGroupIDs", joinedGroupIDs)
//	}()
//	joinedGroupIDListStr, err := g.rcClient.Fetch(g.getJoinedGroupsKey(userID), time.Second*30*60, getJoinedGroupIDList)
//	if err != nil {
//		return nil, err
//	}
//	err = json.Unmarshal([]byte(joinedGroupIDListStr), &joinedGroupIDs)
//	return joinedGroupIDs, utils.Wrap(err, "")
//}

func (g *GroupCacheRedis) DelJoinedGroupID(ctx context.Context, userID string) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "userID", userID)
	}()
	return g.rcClient.TagAsDeleted(g.getJoinedGroupsKey(userID))
}

//func (g *GroupCacheRedis) DelJoinedGroupIDs(ctx context.Context, userIDs []string) (err error) {
//	defer func() {
//		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "userID", userID)
//	}()
//	for _, userID := range userIDs {
//		if err := g.DelJoinedGroupID(ctx, userID); err != nil {
//			return err
//		}
//	}
//	return nil
//}

func (g *GroupCacheRedis) GetGroupMemberInfo(ctx context.Context, groupID, userID string) (groupMember *relationTb.GroupMemberModel, err error) {
	return GetCache(ctx, g.rcClient, g.getGroupMemberInfoKey(groupID, userID), g.expireTime, func(ctx context.Context) (*relationTb.GroupMemberModel, error) {
		return g.groupMember.Take(ctx, groupID, userID)
	})
}

//func (g *GroupCacheRedis) GetGroupMembersInfo(ctx context.Context, groupID, userIDs []string) (groupMember *relationTb.GroupMemberModel, err error) {
//
//	return nil, err
//}

//func (g *GroupCacheRedis) GetGroupMembersInfo(ctx context.Context, count, offset int32, groupID string) (groupMembers []*relation.GroupMember, err error) {
//	defer func() {
//		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "count", count, "offset", offset, "groupID", groupID, "groupMember", groupMembers)
//	}()
//	groupMemberIDList, err := g.GetGroupMemberIDs(ctx, groupID)
//	if err != nil {
//		return nil, err
//	}
//	if count < 0 || offset < 0 {
//		return nil, nil
//	}
//	var groupMemberList []*relation.GroupMember
//	var start, stop int32
//	start = offset
//	stop = offset + count
//	l := int32(len(groupMemberIDList))
//	if start > stop {
//		return nil, nil
//	}
//	if start >= l {
//		return nil, nil
//	}
//	if count != 0 {
//		if stop >= l {
//			stop = l
//		}
//		groupMemberIDList = groupMemberIDList[start:stop]
//	} else {
//		if l < 1000 {
//			stop = l
//		} else {
//			stop = 1000
//		}
//		groupMemberIDList = groupMemberIDList[start:stop]
//	}
//	for _, userID := range groupMemberIDList {
//		groupMember, err := g.GetGroupMemberInfo(ctx, groupID, userID)
//		if err != nil {
//			return
//		}
//		groupMembers = append(groupMembers, groupMember)
//	}
//	return groupMemberList, nil
//}

func (g *GroupCacheRedis) DelGroupMemberInfo(ctx context.Context, groupID, userID string) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupID", groupID, "userID", userID)
	}()
	return g.rcClient.TagAsDeleted(g.getGroupMemberInfoKey(groupID, userID))
}

// groupMemberNum
//func (g *GroupCacheRedis) GetGroupMemberNum(ctx context.Context, groupID string) (num int, err error) {
//	getGroupMemberNum := func() (string, error) {
//		num, err := relation.GetGroupMemberNumByGroupID(groupID)
//		if err != nil {
//			return "", err
//		}
//		return strconv.Itoa(int(num)), nil
//	}
//	defer func() {
//		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupID", groupID, "num", num)
//	}()
//	groupMember, err := g.rcClient.Fetch(g.getGroupMemberNumKey(groupID), time.Second*30*60, getGroupMemberNum)
//	if err != nil {
//		return 0, err
//	}
//	return strconv.Atoi(groupMember)
//}

func (g *GroupCacheRedis) DelGroupMemberNum(ctx context.Context, groupID string) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupID", groupID)
	}()
	return g.rcClient.TagAsDeleted(g.getGroupMemberNumKey(groupID))
}

func (g *GroupCacheRedis) DelGroupInfo(ctx context.Context, groupID string) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupID", groupID)
	}()
	return g.rcClient.TagAsDeleted(g.getGroupInfoKey(groupID))
}

func (g *GroupCacheRedis) DelGroupsInfo(ctx context.Context, groupIDs []string) error {
	for _, groupID := range groupIDs {
		if err := g.DelGroupInfo(ctx, groupID); err != nil {
			return err
		}
	}
	return nil
}
