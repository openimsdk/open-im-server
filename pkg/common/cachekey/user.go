package cachekey

import "time"

const (
	userExpireTime            = time.Second * 60 * 60 * 12
	userInfoKey               = "USER_INFO:"
	userGlobalRecvMsgOptKey   = "USER_GLOBAL_RECV_MSG_OPT_KEY:"
	olineStatusKey            = "ONLINE_STATUS:"
	userOlineStatusExpireTime = time.Second * 60 * 60 * 24
	statusMod                 = 501
	platformID                = "_PlatformIDSuffix"
)

func GetUserInfoKey(userID string) string {
	return userInfoKey + userID
}

func GetUserGlobalRecvMsgOptKey(userID string) string {
	return userGlobalRecvMsgOptKey + userID
}
