package common

type BatchUpdateGroupMember struct {
	GroupID string
	UserID  string
	Map     map[string]any
}

type GroupSimpleUserID struct {
	Hash      uint64
	MemberNum uint32
}
