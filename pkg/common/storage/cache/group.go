package cache

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/common"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
)

type GroupHash interface {
	GetGroupHash(ctx context.Context, groupID string) (uint64, error)
}

type GroupCache interface {
	BatchDeleter
	CloneGroupCache() GroupCache
	GetGroupsInfo(ctx context.Context, groupIDs []string) (groups []*model.Group, err error)
	GetGroupInfo(ctx context.Context, groupID string) (group *model.Group, err error)
	DelGroupsInfo(groupIDs ...string) GroupCache

	GetGroupMembersHash(ctx context.Context, groupID string) (hashCode uint64, err error)
	GetGroupMemberHashMap(ctx context.Context, groupIDs []string) (map[string]*common.GroupSimpleUserID, error)
	DelGroupMembersHash(groupID string) GroupCache

	GetGroupMemberIDs(ctx context.Context, groupID string) (groupMemberIDs []string, err error)

	DelGroupMemberIDs(groupID string) GroupCache

	GetJoinedGroupIDs(ctx context.Context, userID string) (joinedGroupIDs []string, err error)
	DelJoinedGroupID(userID ...string) GroupCache

	GetGroupMemberInfo(ctx context.Context, groupID, userID string) (groupMember *model.GroupMember, err error)
	GetGroupMembersInfo(ctx context.Context, groupID string, userID []string) (groupMembers []*model.GroupMember, err error)
	GetAllGroupMembersInfo(ctx context.Context, groupID string) (groupMembers []*model.GroupMember, err error)
	FindGroupMemberUser(ctx context.Context, groupIDs []string, userID string) ([]*model.GroupMember, error)

	GetGroupRoleLevelMemberIDs(ctx context.Context, groupID string, roleLevel int32) ([]string, error)
	GetGroupOwner(ctx context.Context, groupID string) (*model.GroupMember, error)
	GetGroupsOwner(ctx context.Context, groupIDs []string) ([]*model.GroupMember, error)
	DelGroupRoleLevel(groupID string, roleLevel []int32) GroupCache
	DelGroupAllRoleLevel(groupID string) GroupCache
	DelGroupMembersInfo(groupID string, userID ...string) GroupCache
	GetGroupRoleLevelMemberInfo(ctx context.Context, groupID string, roleLevel int32) ([]*model.GroupMember, error)
	GetGroupRolesLevelMemberInfo(ctx context.Context, groupID string, roleLevels []int32) ([]*model.GroupMember, error)
	GetGroupMemberNum(ctx context.Context, groupID string) (memberNum int64, err error)
	DelGroupsMemberNum(groupID ...string) GroupCache

	//FindSortGroupMemberUserIDs(ctx context.Context, groupID string) ([]string, error)
	//FindSortJoinGroupIDs(ctx context.Context, userID string) ([]string, error)

	DelMaxGroupMemberVersion(groupIDs ...string) GroupCache
	DelMaxJoinGroupVersion(userIDs ...string) GroupCache
	FindMaxGroupMemberVersion(ctx context.Context, groupID string) (*model.VersionLog, error)
	BatchFindMaxGroupMemberVersion(ctx context.Context, groupIDs []string) ([]*model.VersionLog, error)
	FindMaxJoinGroupVersion(ctx context.Context, userID string) (*model.VersionLog, error)
}
