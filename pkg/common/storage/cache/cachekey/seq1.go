package cachekey

const (
	MallocSeq        = "MALLOC_SEQ:"
	MallocSeqLock    = "MALLOC_SEQ_LOCK:"
	MallocMinSeqLock = "MALLOC_MIN_SEQ:"
)

func GetMallocSeqKey(conversationID string) string {
	return MallocSeq + conversationID
}

func GetMallocSeqLockKey(conversationID string) string {
	return MallocSeqLock + conversationID
}

func GetMallocMinSeqKey(conversationID string) string {
	return MallocMinSeqLock + conversationID
}
