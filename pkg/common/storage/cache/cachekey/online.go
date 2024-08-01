package cachekey

import "time"

const (
	OnlineKey     = "ONLINE:"
	OnlineChannel = "online_change"
	//OnlineExpire  = time.Hour / 2
	OnlineExpire = time.Minute / 2 // test
)

func GetOnlineKey(userID string) string {
	return OnlineKey + userID
}
