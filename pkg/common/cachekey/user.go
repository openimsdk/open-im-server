package cachekey

const (
	UserInfoKey             = "USER_INFO:"
	UserGlobalRecvMsgOptKey = "USER_GLOBAL_RECV_MSG_OPT_KEY:"
)

func GetUserInfoKey(userID string) string {
	return UserInfoKey + userID
}

func GetUserGlobalRecvMsgOptKey(userID string) string {
	return UserGlobalRecvMsgOptKey + userID
}
