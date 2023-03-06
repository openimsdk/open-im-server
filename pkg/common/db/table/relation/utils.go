package relation

import (
	"OpenIM/pkg/utils"
	"gorm.io/gorm"
)

type BatchUpdateGroupMember struct {
	GroupID string
	UserID  string
	Map     map[string]any
}

type GroupSimpleUserID struct {
	Hash      uint64
	MemberNum uint32
}

func IsNotFound(err error) bool {
	return utils.Unwrap(err) == gorm.ErrRecordNotFound
}
