package cachekey

const (
	FriendIDsKey        = "FRIEND_IDS:"
	TwoWayFriendsIDsKey = "COMMON_FRIENDS_IDS:"
	FriendKey           = "FRIEND_INFO:"
	IsFriendKey         = "IS_FRIEND:" // local cache key
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

func GetIsFriendKey(possibleFriendUserID, userID string) string {
	return IsFriendKey + possibleFriendUserID + "-" + userID
}
