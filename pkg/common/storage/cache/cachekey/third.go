package cachekey

import (
	"strconv"
)

const (
	getuiToken              = "GETUI_TOKEN"
	getuiTaskID             = "GETUI_TASK_ID"
	fmcToken                = "FCM_TOKEN:"
	userBadgeUnreadCountSum = "USER_BADGE_UNREAD_COUNT_SUM:"
)

func GetFcmAccountTokenKey(account string, platformID int) string {
	return fmcToken + account + ":" + strconv.Itoa(platformID)
}

func GetUserBadgeUnreadCountSumKey(userID string) string {
	return userBadgeUnreadCountSum + userID
}

func GetGetuiTokenKey() string {
	return getuiToken
}
func GetGetuiTaskIDKey() string {
	return getuiTaskID
}
