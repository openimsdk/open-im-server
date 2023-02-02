package cache

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db/localcache"
	"Open_IM/pkg/common/db/relation"
	"Open_IM/pkg/common/db/unrelation"
	"Open_IM/pkg/common/tracelog"
	"Open_IM/pkg/utils"
	"context"
	"encoding/json"
	"github.com/dtm-labs/rockscache"
	"github.com/go-redis/redis/v8"
	"math/big"
	"sort"
	"strconv"
	"sync"
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

type GroupCache struct {
	group        *relation.GroupGorm
	groupMember  *relation.GroupMemberGorm
	groupRequest *relation.GroupRequestGorm
	mongoDB      *unrelation.SuperGroupMongoDriver
	expireTime   time.Duration
	redisClient  *RedisClient
	rcClient     *rockscache.Client

	//local cache
}

func NewGroupCache(rdb redis.UniversalClient, groupDB *relation.GroupGorm, groupMemberDB *relation.GroupMemberGorm, groupRequestDB *relation.GroupRequestGorm, mongoClient *unrelation.SuperGroupMongoDriver, opts rockscache.Options) *GroupCache {
	return &GroupCache{rcClient: rockscache.NewClient(rdb, opts), expireTime: groupExpireTime,
		group: groupDB, groupMember: groupMemberDB, groupRequest: groupRequestDB, redisClient: NewRedisClient(rdb),
		mongoDB: mongoClient,
	}
}

func (g *GroupCache) getRedisClient() *RedisClient {
	return g.redisClient
}

func (g *GroupCache) getGroupInfoKey(groupID string) string {
	return groupInfoKey + groupID
}

func (g *GroupCache) getJoinedSuperGroupsIDKey(userID string) string {
	return joinedSuperGroupsKey + userID
}

func (g *GroupCache) getJoinedGroupsKey(userID string) string {
	return joinedGroupsKey + userID
}

func (g *GroupCache) getGroupMembersHashKey(groupID string) string {
	return groupMembersHashKey + groupID
}

func (g *GroupCache) getGroupMemberIDsKey(groupID string) string {
	return groupMemberIDsKey + groupID
}

func (g *GroupCache) getGroupMemberInfoKey(groupID, userID string) string {
	return groupMemberInfoKey + groupID + "-" + userID
}

func (g *GroupCache) getGroupMemberNumKey(groupID string) string {
	return groupMemberNumKey + groupID
}

// / groupInfo
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

func (g *GroupCache) GetGroupInfo(ctx context.Context, groupID string) (group *relation.GroupGorm, err error) {
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
	group = &relation.GroupGorm{}
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupID", groupID, "group", *group)
	}()
	groupStr, err := g.rcClient.Fetch(g.getGroupInfoKey(groupID), g.expireTime, getGroup)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(groupStr), group)
	return group, utils.Wrap(err, "")
}

func (g *GroupCache) DelGroupInfo(ctx context.Context, groupID string) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupID", groupID)
	}()
	return g.rcClient.TagAsDeleted(g.getGroupInfoKey(groupID))
}

func (g *GroupCache) DelGroupsInfo(ctx context.Context, groupIDs []string) error {
	for _, groupID := range groupIDs {
		if err := g.DelGroupInfo(ctx, groupID); err != nil {
			return err
		}
	}
	return nil
}

// userJoinSuperGroup
func (g *GroupCache) BatchDelJoinedSuperGroupIDs(ctx context.Context, userIDs []string) (err error) {
	for _, userID := range userIDs {
		if err := g.DelJoinedSuperGroupIDs(ctx, userID); err != nil {
			return err
		}
	}
	return nil
}

func (g *GroupCache) DelJoinedSuperGroupIDs(ctx context.Context, userID string) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "userID", userID)
	}()
	return g.rcClient.TagAsDeleted(g.getJoinedSuperGroupsIDKey(userID))
}

func (g *GroupCache) GetJoinedSuperGroupIDs(ctx context.Context, userID string) (joinedSuperGroupIDs []string, err error) {
	getJoinedSuperGroupIDList := func() (string, error) {
		userToSuperGroup, err := g.mongoDB.GetSuperGroupByUserID(ctx, userID)
		if err != nil {
			return "", err
		}
		bytes, err := json.Marshal(userToSuperGroup.GroupIDList)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		return string(bytes), nil
	}
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "userID", userID, "joinedSuperGroupIDs", joinedSuperGroupIDs)
	}()
	joinedSuperGroupListStr, err := g.rcClient.Fetch(g.getJoinedSuperGroupsIDKey(userID), time.Second*30*60, getJoinedSuperGroupIDList)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(joinedSuperGroupListStr), &joinedSuperGroupIDs)
	return joinedSuperGroupIDs, utils.Wrap(err, "")
}

// groupMembersHash
func (g *GroupCache) GetGroupMembersHash(ctx context.Context, groupID string) (hashCodeUint64 uint64, err error) {
	generateHash := func() (string, error) {
		groupInfo, err := g.GetGroupInfo(ctx, groupID)
		if err != nil {
			return "", err
		}
		if groupInfo.Status == constant.GroupStatusDismissed {
			return "0", nil
		}
		groupMemberIDList, err := g.GetGroupMemberIDs(ctx, groupID)
		if err != nil {
			return "", err
		}
		sort.Strings(groupMemberIDList)
		var all string
		for _, v := range groupMemberIDList {
			all += v
		}
		bi := big.NewInt(0)
		bi.SetString(utils.Md5(all)[0:8], 16)
		return strconv.Itoa(int(bi.Uint64())), nil
	}
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupID", groupID, "hashCodeUint64", hashCodeUint64)
	}()
	hashCodeStr, err := g.rcClient.Fetch(g.getGroupMembersHashKey(groupID), time.Second*30*60, generateHash)
	if err != nil {
		return 0, utils.Wrap(err, "fetch failed")
	}
	hashCode, err := strconv.Atoi(hashCodeStr)
	return uint64(hashCode), err
}

func (g *GroupCache) DelGroupMembersHash(ctx context.Context, groupID string) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupID", groupID)
	}()
	return g.rcClient.TagAsDeleted(g.getGroupMembersHashKey(groupID))
}

// groupMemberIDs
// from redis
func (g *GroupCache) GetGroupMemberIDs(ctx context.Context, groupID string) (groupMemberIDs []string, err error) {
	f := func() (string, error) {
		groupInfo, err := g.GetGroupInfo(ctx, groupID)
		if err != nil {
			return "", err
		}
		var groupMemberIDList []string
		if groupInfo.GroupType == constant.SuperGroup {
			superGroup, err := g.mongoDB.GetSuperGroup(ctx, groupID)
			if err != nil {
				return "", err
			}
			groupMemberIDList = superGroup.MemberIDList
		} else {
			groupMemberIDList, err = relation.GetGroupMemberIDListByGroupID(groupID)
			if err != nil {
				return "", err
			}
		}
		bytes, err := json.Marshal(groupMemberIDList)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		return string(bytes), nil
	}
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupID", groupID, "groupMemberIDList", groupMemberIDs)
	}()
	groupIDListStr, err := g.rcClient.Fetch(g.getGroupMemberIDsKey(groupID), time.Second*30*60, f)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(groupIDListStr), &groupMemberIDs)
	return groupMemberIDs, nil
}

func (g *GroupCache) DelGroupMemberIDs(ctx context.Context, groupID string) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupID", groupID)
	}()
	return g.rcClient.TagAsDeleted(g.getGroupMemberIDsKey(groupID))
}

// JoinedGroups
func (g *GroupCache) GetJoinedGroupIDs(ctx context.Context, userID string) (joinedGroupIDs []string, err error) {
	getJoinedGroupIDList := func() (string, error) {
		joinedGroupList, err := relation.GetJoinedGroupIDListByUserID(userID)
		if err != nil {
			return "", err
		}
		bytes, err := json.Marshal(joinedGroupList)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		return string(bytes), nil
	}
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "userID", userID, "joinedGroupIDs", joinedGroupIDs)
	}()
	joinedGroupIDListStr, err := g.rcClient.Fetch(g.getJoinedGroupsKey(userID), time.Second*30*60, getJoinedGroupIDList)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(joinedGroupIDListStr), &joinedGroupIDs)
	return joinedGroupIDs, utils.Wrap(err, "")
}

func (g *GroupCache) DelJoinedGroupIDs(ctx context.Context, userID string) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "userID", userID)
	}()
	return g.rcClient.TagAsDeleted(g.getJoinedGroupsKey(userID))
}

// GetGroupMemberInfo
func (g *GroupCache) GetGroupMemberInfo(ctx context.Context, groupID, userID string) (groupMember *relation.GroupMember, err error) {
	getGroupMemberInfo := func() (string, error) {
		groupMemberInfo, err := relation.GetGroupMemberInfoByGroupIDAndUserID(groupID, userID)
		if err != nil {
			return "", err
		}
		bytes, err := json.Marshal(groupMemberInfo)
		if err != nil {
			return "", utils.Wrap(err, "")
		}
		return string(bytes), nil
	}
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupID", groupID, "userID", userID, "groupMember", *groupMember)
	}()
	groupMemberInfoStr, err := g.rcClient.Fetch(g.getGroupMemberInfoKey(groupID, userID), time.Second*30*60, getGroupMemberInfo)
	if err != nil {
		return nil, err
	}
	groupMember = &relation.GroupMember{}
	err = json.Unmarshal([]byte(groupMemberInfoStr), groupMember)
	return groupMember, utils.Wrap(err, "")
}

func (g *GroupCache) GetGroupMembersInfo(ctx context.Context, count, offset int32, groupID string) (groupMembers []*relation.GroupMember, err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "count", count, "offset", offset, "groupID", groupID, "groupMember", groupMembers)
	}()
	groupMemberIDList, err := g.GetGroupMemberIDs(ctx, groupID)
	if err != nil {
		return nil, err
	}
	if count < 0 || offset < 0 {
		return nil, nil
	}
	var groupMemberList []*relation.GroupMember
	var start, stop int32
	start = offset
	stop = offset + count
	l := int32(len(groupMemberIDList))
	if start > stop {
		return nil, nil
	}
	if start >= l {
		return nil, nil
	}
	if count != 0 {
		if stop >= l {
			stop = l
		}
		groupMemberIDList = groupMemberIDList[start:stop]
	} else {
		if l < 1000 {
			stop = l
		} else {
			stop = 1000
		}
		groupMemberIDList = groupMemberIDList[start:stop]
	}
	for _, userID := range groupMemberIDList {
		groupMember, err := g.GetGroupMemberInfo(ctx, groupID, userID)
		if err != nil {
			return
		}
		groupMembers = append(groupMembers, groupMember)
	}
	return groupMemberList, nil
}

func (g *GroupCache) DelGroupMemberInfo(ctx context.Context, groupID, userID string) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupID", groupID, "userID", userID)
	}()
	return g.rcClient.TagAsDeleted(g.getGroupMemberInfoKey(groupID, userID))
}

// groupMemberNum
func (g *GroupCache) GetGroupMemberNum(ctx context.Context, groupID string) (num int, err error) {
	getGroupMemberNum := func() (string, error) {
		num, err := relation.GetGroupMemberNumByGroupID(groupID)
		if err != nil {
			return "", err
		}
		return strconv.Itoa(int(num)), nil
	}
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupID", groupID, "num", num)
	}()
	groupMember, err := g.rcClient.Fetch(g.getGroupMemberNumKey(groupID), time.Second*30*60, getGroupMemberNum)
	if err != nil {
		return 0, err
	}
	return strconv.Atoi(groupMember)
}

func (g *GroupCache) DelGroupMemberNum(ctx context.Context, groupID string) (err error) {
	defer func() {
		tracelog.SetCtxDebug(ctx, utils.GetFuncName(1), err, "groupID", groupID)
	}()
	return g.rcClient.TagAsDeleted(g.getGroupMemberNumKey(groupID))
}
