package cachekey

const (
	MallocSeq        = "MALLOC_SEQ:"
	MallocMinSeqLock = "MALLOC_MIN_SEQ:"

	SeqUserMaxSeq  = "SEQ_USER_MAX:"
	SeqUserMinSeq  = "SEQ_USER_MIN:"
	SeqUserReadSeq = "SEQ_USER_READ:"
)

func GetMallocSeqKey(conversationID string) string {
	return MallocSeq + conversationID
}

func GetMallocMinSeqKey(conversationID string) string {
	return MallocMinSeqLock + conversationID
}

func GetSeqUserMaxSeqKey(conversationID string, userID string) string {
	return SeqUserMaxSeq + conversationID + ":" + userID
}

func GetSeqUserMinSeqKey(conversationID string, userID string) string {
	return SeqUserMinSeq + conversationID + ":" + userID
}

func GetSeqUserReadSeqKey(conversationID string, userID string) string {
	return SeqUserReadSeq + conversationID + ":" + userID
}
