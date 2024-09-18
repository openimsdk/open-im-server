package cachekey

import (
	"strings"
	"time"
)

const (
	OnlineKey     = "ONLINE:"
	OnlineChannel = "online_change"
	OnlineExpire  = time.Hour / 2
)

func GetOnlineKey(userID string) string {
	return OnlineKey + userID
}

func GetOnlineKeyUserID(key string) string {
	return strings.TrimPrefix(key, OnlineKey)
}
