package controller

import (
	"context"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	redis2 "github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/redis"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/common"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/tools/db/pagination"
	"github.com/openimsdk/tools/db/tx"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/utils/datautil"
	"github.com/redis/go-redis/v9"
)

type GroupDatabase interface {
	// CreateGroup creates new groups along with their members.
	CreateGroup(ctx context.Context, groups []*model.Group, groupMembers []*model.GroupMember) error
	// TakeGroup retrieves a single group by its ID.
	TakeGroup(ctx context.Context, groupID string) (group *model.Group, err error)
	// FindGroup retrieves multiple groups by their IDs.
	FindGroup(ctx context.Context, groupIDs []string) (groups []*model.Group, err error)
	// SearchGroup searches for groups based on a keyword and pagination settings, returns total count and groups.
	SearchGroup(ctx context.Context, keyword string, pagination pagination.Pagination) (int64, []*model.Group, error)
	// UpdateGroup updates the properties of a group identified by its ID.
	UpdateGroup(ctx context.Context, groupID string, data map[string]any) error
	// DismissGroup disbands a group and optionally removes its members based on the deleteMember flag.
	DismissGroup(ctx context.Context, groupID string, deleteMember bool) error

	// TakeGroupMember retrieves a specific group member by group ID and user ID.
	TakeGroupMember(ctx context.Context, groupID string, userID string) (groupMember *model.GroupMember, err error)
	// TakeGroupOwner retrieves the owner of a group by group ID.
	TakeGroupOwner(ctx context.Context, groupID string) (*model.GroupMember, error)
	// FindGroupMembers retrieves members of a group filtered by user IDs.
	FindGroupMembers(ctx context.Context, groupID string, userIDs []string) (groupMembers []*model.GroupMember, err error)
	// FindGroupMemberUser retrieves groups that a user is a member of, filtered by group IDs.
	FindGroupMemberUser(ctx context.Context, groupIDs []string, userID string) (groupMembers []*model.GroupMember, err error)
	// FindGroupMemberRoleLevels retrieves group members filtered by their role levels within a group.
	FindGroupMemberRoleLevels(ctx context.Context, groupID string, roleLevels []int32) (groupMembers []*model.GroupMember, err error)
	// FindGroupMemberAll retrieves all members of a group.
	FindGroupMemberAll(ctx context.Context, groupID string) (groupMembers []*model.GroupMember, err error)
	// FindGroupsOwner retrieves the owners for multiple groups.
	FindGroupsOwner(ctx context.Context, groupIDs []string) ([]*model.GroupMember, error)
	// FindGroupMemberUserID retrieves the user IDs of all members in a group.
	FindGroupMemberUserID(ctx context.Context, groupID string) ([]string, error)
	// FindGroupMemberNum retrieves the number of members in a group.
	FindGroupMemberNum(ctx context.Context, groupID string) (uint32, error)
	// FindUserManagedGroupID retrieves group IDs managed by a user.
	FindUserManagedGroupID(ctx context.Context, userID string) (groupIDs []string, err error)
	// PageGroupRequest paginates through group requests for specified groups.
	PageGroupRequest(ctx context.Context, groupIDs []string, handleResults []int, pagination pagination.Pagination) (int64, []*model.GroupRequest, error)
	// GetGroupRoleLevelMemberIDs retrieves user IDs of group members with a specific role level.
	GetGroupRoleLevelMemberIDs(ctx context.Context, groupID string, roleLevel int32) ([]string, error)

	// PageGetJoinGroup paginates through groups that a user has joined.
	PageGetJoinGroup(ctx context.Context, userID string, pagination pagination.Pagination) (total int64, totalGroupMembers []*model.GroupMember, err error)
	// PageGetGroupMember paginates through members of a group.
	PageGetGroupMember(ctx context.Context, groupID string, pagination pagination.Pagination) (total int64, totalGroupMembers []*model.GroupMember, err error)
	// SearchGroupMember searches for group members based on a keyword, group ID, and pagination settings.
	SearchGroupMember(ctx context.Context, keyword string, groupID string, pagination pagination.Pagination) (int64, []*model.GroupMember, error)
	// HandlerGroupRequest processes a group join request with a specified result.
	HandlerGroupRequest(ctx context.Context, groupID string, userID string, handledMsg string, handleResult int32, member *model.GroupMember) error
	// DeleteGroupMember removes specified users from a group.
	DeleteGroupMember(ctx context.Context, groupID string, userIDs []string) error
	// MapGroupMemberUserID maps group IDs to their members' simplified user IDs.
	MapGroupMemberUserID(ctx context.Context, groupIDs []string) (map[string]*common.GroupSimpleUserID, error)
	// MapGroupMemberNum maps group IDs to their member count.
	MapGroupMemberNum(ctx context.Context, groupIDs []string) (map[string]uint32, error)
	// TransferGroupOwner transfers the ownership of a group to another user.
	TransferGroupOwner(ctx context.Context, groupID string, oldOwnerUserID, newOwnerUserID string, roleLevel int32) error
	// UpdateGroupMember updates properties of a group member.
	UpdateGroupMember(ctx context.Context, groupID string, userID string, data map[string]any) error
	// UpdateGroupMembers batch updates properties of group members.
	UpdateGroupMembers(ctx context.Context, data []*common.BatchUpdateGroupMember) error

	// CreateGroupRequest creates new group join requests.
	CreateGroupRequest(ctx context.Context, requests []*model.GroupRequest) error
	// TakeGroupRequest retrieves a specific group join request.
	TakeGroupRequest(ctx context.Context, groupID string, userID string) (*model.GroupRequest, error)
	// FindGroupRequests retrieves multiple group join requests.
	FindGroupRequests(ctx context.Context, groupID string, userIDs []string) ([]*model.GroupRequest, error)
	// PageGroupRequestUser paginates through group join requests made by a user.
	PageGroupRequestUser(ctx context.Context, userID string, groupIDs []string, handleResults []int, pagination pagination.Pagination) (int64, []*model.GroupRequest, error)

	// CountTotal counts the total number of groups as of a certain date.
	CountTotal(ctx context.Context, before *time.Time) (count int64, err error)
	// CountRangeEverydayTotal counts the daily group creation total within a specified date range.
	CountRangeEverydayTotal(ctx context.Context, start time.Time, end time.Time) (map[string]int64, error)
	// DeleteGroupMemberHash deletes the hash entries for group members in specified groups.
	DeleteGroupMemberHash(ctx context.Context, groupIDs []string) error

	FindMemberIncrVersion(ctx context.Context, groupID string, version uint, limit int) (*model.VersionLog, error)
	BatchFindMemberIncrVersion(ctx context.Context, groupIDs []string, versions []uint64, limits []int) (map[string]*model.VersionLog, error)
	FindJoinIncrVersion(ctx context.Context, userID string, version uint, limit int) (*model.VersionLog, error)
	MemberGroupIncrVersion(ctx context.Context, groupID string, userIDs []string, state int32) error

	//FindSortGroupMemberUserIDs(ctx context.Context, groupID string) ([]string, error)
	//FindSortJoinGroupIDs(ctx context.Context, userID string) ([]string, error)

	FindMaxGroupMemberVersionCache(ctx context.Context, groupID string) (*model.VersionLog, error)
	BatchFindMaxGroupMemberVersionCache(ctx context.Context, groupIDs []string) (map[string]*model.VersionLog, error)
	FindMaxJoinGroupVersionCache(ctx context.Context, userID string) (*model.VersionLog, error)

	SearchJoinGroup(ctx context.Context, userID string, keyword string, pagination pagination.Pagination) (int64, []*model.Group, error)

	FindJoinGroupID(ctx context.Context, userID string) ([]string, error)

	GetGroupApplicationUnhandledCount(ctx context.Context, groupIDs []string, ts int64) (int64, error)
}

func NewGroupDatabase(
	rdb redis.UniversalClient,
	localCache *config.LocalCache,
	groupDB database.Group,
	groupMemberDB database.GroupMember,
	groupRequestDB database.GroupRequest,
	ctxTx tx.Tx,
	groupHash cache.GroupHash,
) GroupDatabase {
	return &groupDatabase{
		groupDB:        groupDB,
		groupMemberDB:  groupMemberDB,
		groupRequestDB: groupRequestDB,
		ctxTx:          ctxTx,
		cache:          redis2.NewGroupCacheRedis(rdb, localCache, groupDB, groupMemberDB, groupRequestDB, groupHash),
	}
}

type groupDatabase struct {
	groupDB        database.Group
	groupMemberDB  database.GroupMember
	groupRequestDB database.GroupRequest
	ctxTx          tx.Tx
	cache          cache.GroupCache
}

func (g *groupDatabase) FindJoinGroupID(ctx context.Context, userID string) ([]string, error) {
	return g.cache.GetJoinedGroupIDs(ctx, userID)
}

func (g *groupDatabase) FindGroupMembers(ctx context.Context, groupID string, userIDs []string) ([]*model.GroupMember, error) {
	return g.cache.GetGroupMembersInfo(ctx, groupID, userIDs)
}

func (g *groupDatabase) FindGroupMemberUser(ctx context.Context, groupIDs []string, userID string) ([]*model.GroupMember, error) {
	return g.cache.FindGroupMemberUser(ctx, groupIDs, userID)
}

func (g *groupDatabase) FindGroupMemberRoleLevels(ctx context.Context, groupID string, roleLevels []int32) ([]*model.GroupMember, error) {
	return g.cache.GetGroupRolesLevelMemberInfo(ctx, groupID, roleLevels)
}

func (g *groupDatabase) FindGroupMemberAll(ctx context.Context, groupID string) ([]*model.GroupMember, error) {
	return g.cache.GetAllGroupMembersInfo(ctx, groupID)
}

func (g *groupDatabase) FindGroupsOwner(ctx context.Context, groupIDs []string) ([]*model.GroupMember, error) {
	return g.cache.GetGroupsOwner(ctx, groupIDs)
}

func (g *groupDatabase) GetGroupRoleLevelMemberIDs(ctx context.Context, groupID string, roleLevel int32) ([]string, error) {
	return g.cache.GetGroupRoleLevelMemberIDs(ctx, groupID, roleLevel)
}

func (g *groupDatabase) CreateGroup(ctx context.Context, groups []*model.Group, groupMembers []*model.GroupMember) error {
	if len(groups)+len(groupMembers) == 0 {
		return nil
	}
	return g.ctxTx.Transaction(ctx, func(ctx context.Context) error {
		c := g.cache.CloneGroupCache()
		if len(groups) > 0 {
			if err := g.groupDB.Create(ctx, groups); err != nil {
				return err
			}
			for _, group := range groups {
				c = c.DelGroupsInfo(group.GroupID).
					DelGroupMembersHash(group.GroupID).
					DelGroupsMemberNum(group.GroupID).
					DelGroupMemberIDs(group.GroupID).
					DelGroupAllRoleLevel(group.GroupID).
					DelMaxGroupMemberVersion(group.GroupID)
			}
		}
		if len(groupMembers) > 0 {
			if err := g.groupMemberDB.Create(ctx, groupMembers); err != nil {
				return err
			}
			for _, groupMember := range groupMembers {
				c = c.DelGroupMembersHash(groupMember.GroupID).
					DelGroupsMemberNum(groupMember.GroupID).
					DelGroupMemberIDs(groupMember.GroupID).
					DelJoinedGroupID(groupMember.UserID).
					DelGroupMembersInfo(groupMember.GroupID, groupMember.UserID).
					DelGroupAllRoleLevel(groupMember.GroupID).
					DelMaxJoinGroupVersion(groupMember.UserID).
					DelMaxGroupMemberVersion(groupMember.GroupID)
			}
		}
		return c.ChainExecDel(ctx)
	})
}

func (g *groupDatabase) FindGroupMemberUserID(ctx context.Context, groupID string) ([]string, error) {
	return g.cache.GetGroupMemberIDs(ctx, groupID)
}

func (g *groupDatabase) FindGroupMemberNum(ctx context.Context, groupID string) (uint32, error) {
	num, err := g.cache.GetGroupMemberNum(ctx, groupID)
	if err != nil {
		return 0, err
	}
	return uint32(num), nil
}

func (g *groupDatabase) TakeGroup(ctx context.Context, groupID string) (*model.Group, error) {
	return g.cache.GetGroupInfo(ctx, groupID)
}

func (g *groupDatabase) FindGroup(ctx context.Context, groupIDs []string) ([]*model.Group, error) {
	return g.cache.GetGroupsInfo(ctx, groupIDs)
}

func (g *groupDatabase) SearchGroup(ctx context.Context, keyword string, pagination pagination.Pagination) (int64, []*model.Group, error) {
	return g.groupDB.Search(ctx, keyword, pagination)
}

func (g *groupDatabase) UpdateGroup(ctx context.Context, groupID string, data map[string]any) error {
	return g.ctxTx.Transaction(ctx, func(ctx context.Context) error {
		if err := g.groupDB.UpdateMap(ctx, groupID, data); err != nil {
			return err
		}
		if err := g.groupMemberDB.MemberGroupIncrVersion(ctx, groupID, []string{""}, model.VersionStateUpdate); err != nil {
			return err
		}
		return g.cache.CloneGroupCache().DelGroupsInfo(groupID).DelMaxGroupMemberVersion(groupID).ChainExecDel(ctx)
	})
}

func (g *groupDatabase) DismissGroup(ctx context.Context, groupID string, deleteMember bool) error {
	return g.ctxTx.Transaction(ctx, func(ctx context.Context) error {
		c := g.cache.CloneGroupCache()
		if err := g.groupDB.UpdateStatus(ctx, groupID, constant.GroupStatusDismissed); err != nil {
			return err
		}
		if deleteMember {
			userIDs, err := g.cache.GetGroupMemberIDs(ctx, groupID)
			if err != nil {
				return err
			}
			if err := g.groupMemberDB.Delete(ctx, groupID, nil); err != nil {
				return err
			}
			c = c.DelJoinedGroupID(userIDs...).
				DelGroupMemberIDs(groupID).
				DelGroupsMemberNum(groupID).
				DelGroupMembersHash(groupID).
				DelGroupAllRoleLevel(groupID).
				DelGroupMembersInfo(groupID, userIDs...).
				DelMaxGroupMemberVersion(groupID).
				DelMaxJoinGroupVersion(userIDs...)
			for _, userID := range userIDs {
				if err := g.groupMemberDB.JoinGroupIncrVersion(ctx, userID, []string{groupID}, model.VersionStateDelete); err != nil {
					return err
				}
			}
		} else {
			if err := g.groupMemberDB.MemberGroupIncrVersion(ctx, groupID, []string{""}, model.VersionStateUpdate); err != nil {
				return err
			}
			c = c.DelMaxGroupMemberVersion(groupID)
		}
		return c.DelGroupsInfo(groupID).ChainExecDel(ctx)
	})
}

func (g *groupDatabase) TakeGroupMember(ctx context.Context, groupID string, userID string) (*model.GroupMember, error) {
	return g.cache.GetGroupMemberInfo(ctx, groupID, userID)
}

func (g *groupDatabase) TakeGroupOwner(ctx context.Context, groupID string) (*model.GroupMember, error) {
	return g.cache.GetGroupOwner(ctx, groupID)
}

func (g *groupDatabase) FindUserManagedGroupID(ctx context.Context, userID string) (groupIDs []string, err error) {
	return g.groupMemberDB.FindUserManagedGroupID(ctx, userID)
}

func (g *groupDatabase) PageGroupRequest(ctx context.Context, groupIDs []string, handleResults []int, pagination pagination.Pagination) (int64, []*model.GroupRequest, error) {
	return g.groupRequestDB.PageGroup(ctx, groupIDs, handleResults, pagination)
}

func (g *groupDatabase) PageGetJoinGroup(ctx context.Context, userID string, pagination pagination.Pagination) (total int64, totalGroupMembers []*model.GroupMember, err error) {
	groupIDs, err := g.cache.GetJoinedGroupIDs(ctx, userID)
	if err != nil {
		return 0, nil, err
	}
	for _, groupID := range datautil.Paginate(groupIDs, int(pagination.GetPageNumber()), int(pagination.GetShowNumber())) {
		groupMembers, err := g.cache.GetGroupMembersInfo(ctx, groupID, []string{userID})
		if err != nil {
			return 0, nil, err
		}
		totalGroupMembers = append(totalGroupMembers, groupMembers...)
	}
	return int64(len(groupIDs)), totalGroupMembers, nil
}

func (g *groupDatabase) PageGetGroupMember(ctx context.Context, groupID string, pagination pagination.Pagination) (total int64, totalGroupMembers []*model.GroupMember, err error) {
	groupMemberIDs, err := g.cache.GetGroupMemberIDs(ctx, groupID)
	if err != nil {
		return 0, nil, err
	}
	pageIDs := datautil.Paginate(groupMemberIDs, int(pagination.GetPageNumber()), int(pagination.GetShowNumber()))
	if len(pageIDs) == 0 {
		return int64(len(groupMemberIDs)), nil, nil
	}
	members, err := g.cache.GetGroupMembersInfo(ctx, groupID, pageIDs)
	if err != nil {
		return 0, nil, err
	}
	return int64(len(groupMemberIDs)), members, nil
}

func (g *groupDatabase) SearchGroupMember(ctx context.Context, keyword string, groupID string, pagination pagination.Pagination) (int64, []*model.GroupMember, error) {
	return g.groupMemberDB.SearchMember(ctx, keyword, groupID, pagination)
}

func (g *groupDatabase) HandlerGroupRequest(ctx context.Context, groupID string, userID string, handledMsg string, handleResult int32, member *model.GroupMember) error {
	return g.ctxTx.Transaction(ctx, func(ctx context.Context) error {
		if err := g.groupRequestDB.UpdateHandler(ctx, groupID, userID, handledMsg, handleResult); err != nil {
			return err
		}
		if member != nil {
			c := g.cache.CloneGroupCache()
			if err := g.groupMemberDB.Create(ctx, []*model.GroupMember{member}); err != nil {
				return err
			}
			c = c.DelGroupMembersHash(groupID).
				DelGroupMembersInfo(groupID, member.UserID).
				DelGroupMemberIDs(groupID).
				DelGroupsMemberNum(groupID).
				DelJoinedGroupID(member.UserID).
				DelGroupRoleLevel(groupID, []int32{member.RoleLevel}).
				DelMaxJoinGroupVersion(userID).
				DelMaxGroupMemberVersion(groupID)
			if err := c.ChainExecDel(ctx); err != nil {
				return err
			}
		}
		return nil
	})
}

func (g *groupDatabase) DeleteGroupMember(ctx context.Context, groupID string, userIDs []string) error {
	return g.ctxTx.Transaction(ctx, func(ctx context.Context) error {
		if err := g.groupMemberDB.Delete(ctx, groupID, userIDs); err != nil {
			return err
		}
		c := g.cache.CloneGroupCache()
		return c.DelGroupMembersHash(groupID).
			DelGroupMemberIDs(groupID).
			DelGroupsMemberNum(groupID).
			DelJoinedGroupID(userIDs...).
			DelGroupMembersInfo(groupID, userIDs...).
			DelGroupAllRoleLevel(groupID).
			DelMaxGroupMemberVersion(groupID).
			DelMaxJoinGroupVersion(userIDs...).
			ChainExecDel(ctx)
	})
}

func (g *groupDatabase) MapGroupMemberUserID(ctx context.Context, groupIDs []string) (map[string]*common.GroupSimpleUserID, error) {
	return g.cache.GetGroupMemberHashMap(ctx, groupIDs)
}

func (g *groupDatabase) MapGroupMemberNum(ctx context.Context, groupIDs []string) (m map[string]uint32, err error) {
	m = make(map[string]uint32)
	for _, groupID := range groupIDs {
		num, err := g.cache.GetGroupMemberNum(ctx, groupID)
		if err != nil {
			return nil, err
		}
		m[groupID] = uint32(num)
	}
	return m, nil
}

func (g *groupDatabase) TransferGroupOwner(ctx context.Context, groupID string, oldOwnerUserID, newOwnerUserID string, roleLevel int32) error {
	return g.ctxTx.Transaction(ctx, func(ctx context.Context) error {
		if err := g.groupMemberDB.UpdateUserRoleLevels(ctx, groupID, oldOwnerUserID, roleLevel, newOwnerUserID, constant.GroupOwner); err != nil {
			return err
		}
		c := g.cache.CloneGroupCache()
		return c.DelGroupMembersInfo(groupID, oldOwnerUserID, newOwnerUserID).
			DelGroupAllRoleLevel(groupID).
			DelGroupMembersHash(groupID).
			DelMaxGroupMemberVersion(groupID).
			DelGroupMemberIDs(groupID).
			ChainExecDel(ctx)
	})
}

func (g *groupDatabase) UpdateGroupMember(ctx context.Context, groupID string, userID string, data map[string]any) error {
	if len(data) == 0 {
		return nil
	}
	return g.ctxTx.Transaction(ctx, func(ctx context.Context) error {
		if err := g.groupMemberDB.Update(ctx, groupID, userID, data); err != nil {
			return err
		}
		c := g.cache.CloneGroupCache()
		c = c.DelGroupMembersInfo(groupID, userID)
		if g.groupMemberDB.IsUpdateRoleLevel(data) {
			c = c.DelGroupAllRoleLevel(groupID).DelGroupMemberIDs(groupID)
		}
		c = c.DelMaxGroupMemberVersion(groupID)
		return c.ChainExecDel(ctx)
	})
}

func (g *groupDatabase) UpdateGroupMembers(ctx context.Context, data []*common.BatchUpdateGroupMember) error {
	return g.ctxTx.Transaction(ctx, func(ctx context.Context) error {
		c := g.cache.CloneGroupCache()
		for _, item := range data {
			if err := g.groupMemberDB.Update(ctx, item.GroupID, item.UserID, item.Map); err != nil {
				return err
			}
			if g.groupMemberDB.IsUpdateRoleLevel(item.Map) {
				c = c.DelGroupAllRoleLevel(item.GroupID).DelGroupMemberIDs(item.GroupID)
			}
			c = c.DelGroupMembersInfo(item.GroupID, item.UserID).DelMaxGroupMemberVersion(item.GroupID).DelGroupMembersHash(item.GroupID)
		}
		return c.ChainExecDel(ctx)
	})
}

func (g *groupDatabase) CreateGroupRequest(ctx context.Context, requests []*model.GroupRequest) error {
	return g.ctxTx.Transaction(ctx, func(ctx context.Context) error {
		for _, request := range requests {
			if err := g.groupRequestDB.Delete(ctx, request.GroupID, request.UserID); err != nil {
				return err
			}
		}
		return g.groupRequestDB.Create(ctx, requests)
	})
}

func (g *groupDatabase) TakeGroupRequest(ctx context.Context, groupID string, userID string) (*model.GroupRequest, error) {
	return g.groupRequestDB.Take(ctx, groupID, userID)
}

func (g *groupDatabase) PageGroupRequestUser(ctx context.Context, userID string, groupIDs []string, handleResults []int, pagination pagination.Pagination) (int64, []*model.GroupRequest, error) {
	return g.groupRequestDB.Page(ctx, userID, groupIDs, handleResults, pagination)
}

func (g *groupDatabase) CountTotal(ctx context.Context, before *time.Time) (count int64, err error) {
	return g.groupDB.CountTotal(ctx, before)
}

func (g *groupDatabase) CountRangeEverydayTotal(ctx context.Context, start time.Time, end time.Time) (map[string]int64, error) {
	return g.groupDB.CountRangeEverydayTotal(ctx, start, end)
}

func (g *groupDatabase) FindGroupRequests(ctx context.Context, groupID string, userIDs []string) ([]*model.GroupRequest, error) {
	return g.groupRequestDB.FindGroupRequests(ctx, groupID, userIDs)
}

func (g *groupDatabase) DeleteGroupMemberHash(ctx context.Context, groupIDs []string) error {
	if len(groupIDs) == 0 {
		return nil
	}
	c := g.cache.CloneGroupCache()
	for _, groupID := range groupIDs {
		c = c.DelGroupMembersHash(groupID)
	}
	return c.ChainExecDel(ctx)
}

func (g *groupDatabase) FindMemberIncrVersion(ctx context.Context, groupID string, version uint, limit int) (*model.VersionLog, error) {
	return g.groupMemberDB.FindMemberIncrVersion(ctx, groupID, version, limit)
}

func (g *groupDatabase) BatchFindMemberIncrVersion(ctx context.Context, groupIDs []string, versions []uint64, limits []int) (map[string]*model.VersionLog, error) {
	if len(groupIDs) == 0 {
		return nil, errs.Wrap(errs.New("groupIDs is nil."))
	}

	// convert []uint64 to []uint
	var uintVersions []uint
	for _, version := range versions {
		uintVersions = append(uintVersions, uint(version))
	}

	versionLogs, err := g.groupMemberDB.BatchFindMemberIncrVersion(ctx, groupIDs, uintVersions, limits)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	groupMemberIncrVersionsMap := datautil.SliceToMap(versionLogs, func(e *model.VersionLog) string {
		return e.DID
	})

	return groupMemberIncrVersionsMap, nil
}

func (g *groupDatabase) FindJoinIncrVersion(ctx context.Context, userID string, version uint, limit int) (*model.VersionLog, error) {
	return g.groupMemberDB.FindJoinIncrVersion(ctx, userID, version, limit)
}

func (g *groupDatabase) FindMaxGroupMemberVersionCache(ctx context.Context, groupID string) (*model.VersionLog, error) {
	return g.cache.FindMaxGroupMemberVersion(ctx, groupID)
}

func (g *groupDatabase) BatchFindMaxGroupMemberVersionCache(ctx context.Context, groupIDs []string) (map[string]*model.VersionLog, error) {
	if len(groupIDs) == 0 {
		return nil, errs.Wrap(errs.New("groupIDs is nil in Cache."))
	}
	versionLogs, err := g.cache.BatchFindMaxGroupMemberVersion(ctx, groupIDs)
	if err != nil {
		return nil, errs.Wrap(err)
	}
	maxGroupMemberVersionsMap := datautil.SliceToMap(versionLogs, func(e *model.VersionLog) string {
		return e.DID
	})
	return maxGroupMemberVersionsMap, nil
}

func (g *groupDatabase) FindMaxJoinGroupVersionCache(ctx context.Context, userID string) (*model.VersionLog, error) {
	return g.cache.FindMaxJoinGroupVersion(ctx, userID)
}

func (g *groupDatabase) SearchJoinGroup(ctx context.Context, userID string, keyword string, pagination pagination.Pagination) (int64, []*model.Group, error) {
	groupIDs, err := g.cache.GetJoinedGroupIDs(ctx, userID)
	if err != nil {
		return 0, nil, err
	}
	return g.groupDB.SearchJoin(ctx, groupIDs, keyword, pagination)
}

func (g *groupDatabase) MemberGroupIncrVersion(ctx context.Context, groupID string, userIDs []string, state int32) error {
	if err := g.groupMemberDB.MemberGroupIncrVersion(ctx, groupID, userIDs, state); err != nil {
		return err
	}
	return g.cache.DelMaxGroupMemberVersion(groupID).ChainExecDel(ctx)
}

func (g *groupDatabase) GetGroupApplicationUnhandledCount(ctx context.Context, groupIDs []string, ts int64) (int64, error) {
	return g.groupRequestDB.GetUnhandledCount(ctx, groupIDs, ts)
}
