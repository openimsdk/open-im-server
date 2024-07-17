package useronline

import (
	"errors"
	"strconv"
	"strings"
)

func ParseUserOnlineStatus(payload string) (string, []int32, error) {
	arr := strings.Split(payload, ":")
	if len(arr) == 0 {
		return "", nil, errors.New("invalid data")
	}
	userID := arr[len(arr)-1]
	if userID == "" {
		return "", nil, errors.New("userID is empty")
	}
	platformIDs := make([]int32, len(arr)-1)
	for i := range platformIDs {
		platformID, err := strconv.Atoi(arr[i])
		if err != nil {
			return "", nil, err
		}
		platformIDs[i] = int32(platformID)
	}
	return userID, platformIDs, nil
}
