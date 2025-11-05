package cachekey

import (
	"strconv"
	"time"
)

const (
	groupExpireTime             = time.Second * 60 * 60 * 12
	GroupInfoKey                = "GROUP_INFO:"
	GroupMemberIDsKey           = "GROUP_MEMBER_IDS:"
	GroupMembersHashKey         = "GROUP_MEMBERS_HASH2:"
	GroupMemberInfoKey          = "GROUP_MEMBER_INFO:"
	JoinedGroupsKey             = "JOIN_GROUPS_KEY:"
	GroupMemberNumKey           = "GROUP_MEMBER_NUM_CACHE:"
	GroupRoleLevelMemberIDsKey  = "GROUP_ROLE_LEVEL_MEMBER_IDS:"
	GroupAdminLevelMemberIDsKey = "GROUP_ADMIN_LEVEL_MEMBER_IDS:"
	GroupMemberMaxVersionKey    = "GROUP_MEMBER_MAX_VERSION:"
	GroupJoinMaxVersionKey      = "GROUP_JOIN_MAX_VERSION:"
)

func GetGroupInfoKey(groupID string) string {
	return GroupInfoKey + groupID
}

func GetJoinedGroupsKey(userID string) string {
	return JoinedGroupsKey + userID
}

func GetGroupMembersHashKey(groupID string) string {
	return GroupMembersHashKey + groupID
}

func GetGroupMemberIDsKey(groupID string) string {
	return GroupMemberIDsKey + groupID
}

func GetGroupMemberInfoKey(groupID, userID string) string {
	return GroupMemberInfoKey + groupID + "-" + userID
}

func GetGroupMemberNumKey(groupID string) string {
	return GroupMemberNumKey + groupID
}

func GetGroupRoleLevelMemberIDsKey(groupID string, roleLevel int32) string {
	return GroupRoleLevelMemberIDsKey + groupID + "-" + strconv.Itoa(int(roleLevel))
}

func GetGroupMemberMaxVersionKey(groupID string) string {
	return GroupMemberMaxVersionKey + groupID
}

func GetJoinGroupMaxVersionKey(userID string) string {
	return GroupJoinMaxVersionKey + userID
}
