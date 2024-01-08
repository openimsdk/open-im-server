package cachekey

import "time"

const (
	friendExpireTime    = time.Second * 60 * 60 * 12
	friendIDsKey        = "FRIEND_IDS:"
	twoWayFriendsIDsKey = "COMMON_FRIENDS_IDS:"
	friendKey           = "FRIEND_INFO:"
	isFriendKey         = "IS_FRIEND:"
)

func GetFriendIDsKey(ownerUserID string) string {
	return friendIDsKey + ownerUserID
}

func GetTwoWayFriendsIDsKey(ownerUserID string) string {
	return twoWayFriendsIDsKey + ownerUserID
}

func GetFriendKey(ownerUserID, friendUserID string) string {
	return friendKey + ownerUserID + "-" + friendUserID
}

func GetIsFriendKey(possibleFriendUserID, userID string) string {
	return isFriendKey + possibleFriendUserID + "-" + userID
}
