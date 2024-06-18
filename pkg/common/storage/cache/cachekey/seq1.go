package cachekey

const (
	MallocSeq        = "MALLOC_SEQ:"
	MallocMinSeqLock = "MALLOC_MIN_SEQ:"
)

func GetMallocSeqKey(conversationID string) string {
	return MallocSeq + conversationID
}

func GetMallocMinSeqKey(conversationID string) string {
	return MallocMinSeqLock + conversationID
}
