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
