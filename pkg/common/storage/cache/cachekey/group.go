// Copyright © 2024 OpenIM. All rights reserved.
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
