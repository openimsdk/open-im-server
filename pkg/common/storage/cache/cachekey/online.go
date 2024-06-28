package cachekey

import "time"

const (
	OnlineKey     = "ONLINE:"
	OnlineChannel = "online_change"
	OnlineExpire  = time.Hour / 2
)

func GetOnlineKey(userID string) string {
	return OnlineKey + userID
}
