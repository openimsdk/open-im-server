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

const (
	FriendIDsKey        = "FRIEND_IDS:"
	TwoWayFriendsIDsKey = "COMMON_FRIENDS_IDS:"
	FriendKey           = "FRIEND_INFO:"
	IsFriendKey         = "IS_FRIEND:" // local cache key
	//FriendSyncSortUserIDsKey = "FRIEND_SYNC_SORT_USER_IDS:"
	FriendMaxVersionKey = "FRIEND_MAX_VERSION:"
)

func GetFriendIDsKey(ownerUserID string) string {
	return FriendIDsKey + ownerUserID
}

func GetTwoWayFriendsIDsKey(ownerUserID string) string {
	return TwoWayFriendsIDsKey + ownerUserID
}

func GetFriendKey(ownerUserID, friendUserID string) string {
	return FriendKey + ownerUserID + "-" + friendUserID
}

func GetFriendMaxVersionKey(ownerUserID string) string {
	return FriendMaxVersionKey + ownerUserID
}

func GetIsFriendKey(possibleFriendUserID, userID string) string {
	return IsFriendKey + possibleFriendUserID + "-" + userID
}

//func GetFriendSyncSortUserIDsKey(ownerUserID string, count int) string {
//	return FriendSyncSortUserIDsKey + strconv.Itoa(count) + ":" + ownerUserID
//}
