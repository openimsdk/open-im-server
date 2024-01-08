package cachekey

import (
	"strconv"
	"time"
)

const (
	groupExpireTime            = time.Second * 60 * 60 * 12
	groupInfoKey               = "GROUP_INFO:"
	groupMemberIDsKey          = "GROUP_MEMBER_IDS:"
	groupMembersHashKey        = "GROUP_MEMBERS_HASH2:"
	groupMemberInfoKey         = "GROUP_MEMBER_INFO:"
	joinedGroupsKey            = "JOIN_GROUPS_KEY:"
	groupMemberNumKey          = "GROUP_MEMBER_NUM_CACHE:"
	groupRoleLevelMemberIDsKey = "GROUP_ROLE_LEVEL_MEMBER_IDS:"
)

func GetGroupInfoKey(groupID string) string {
	return groupInfoKey + groupID
}

func GetJoinedGroupsKey(userID string) string {
	return joinedGroupsKey + userID
}

func GetGroupMembersHashKey(groupID string) string {
	return groupMembersHashKey + groupID
}

func GetGroupMemberIDsKey(groupID string) string {
	return groupMemberIDsKey + groupID
}

func GetGroupMemberInfoKey(groupID, userID string) string {
	return groupMemberInfoKey + groupID + "-" + userID
}

func GetGroupMemberNumKey(groupID string) string {
	return groupMemberNumKey + groupID
}

func GetGroupRoleLevelMemberIDsKey(groupID string, roleLevel int32) string {
	return groupRoleLevelMemberIDsKey + groupID + "-" + strconv.Itoa(int(roleLevel))
}
