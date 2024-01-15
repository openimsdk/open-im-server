package cachekey

const (
	userInfoKey             = "USER_INFO:"
	userGlobalRecvMsgOptKey = "USER_GLOBAL_RECV_MSG_OPT_KEY:"
)

func GetUserInfoKey(userID string) string {
	return userInfoKey + userID
}

func GetUserGlobalRecvMsgOptKey(userID string) string {
	return userGlobalRecvMsgOptKey + userID
}
