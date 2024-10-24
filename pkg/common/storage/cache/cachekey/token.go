package cachekey

import (
	"github.com/openimsdk/protocol/constant"
	"strings"
)

const (
	UidPidToken = "UID_PID_TOKEN_STATUS:"
)

func GetTokenKey(userID string, platformID int) string {
	return UidPidToken + userID + ":" + constant.PlatformIDToName(platformID)
}

func GetAllPlatformTokenKey(userID string) []string {
	res := make([]string, len(constant.PlatformID2Name))
	for k := range constant.PlatformID2Name {
		res[k-1] = GetTokenKey(userID, k)
	}
	return res
}

func GetPlatformIDByTokenKey(key string) int {
	splitKey := strings.Split(key, ":")
	platform := splitKey[len(splitKey)-1]
	return constant.PlatformNameToID(platform)
}
