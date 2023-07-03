package relation

import (
	"gorm.io/gorm"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
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
