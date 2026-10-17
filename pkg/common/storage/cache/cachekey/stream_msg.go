package cachekey

const streamMessageCache = "STREAM_MSG:"

func GetStreamMsgKey(conversationID string, clientMsgID string) string {
	return streamMessageCache + conversationID + ":" + clientMsgID
}
