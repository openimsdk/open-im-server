package cachekey

import (
	"strconv"
)

const (
	sendMsgFailedFlag = "SEND_MSG_FAILED_FLAG:"
	messageCache      = "MSG_CACHE:"
)

func GetMsgCacheKey(conversationID string, seq int64) string {
	return messageCache + conversationID + ":" + strconv.Itoa(int(seq))
}

func GetSendMsgKey(id string) string {
	return sendMsgFailedFlag + id
}
