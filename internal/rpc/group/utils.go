package group

import (
	"Open_IM/pkg/common/tracelog"
	"gorm.io/gorm"
)

func IsNotFound(err error) bool {
	if err == nil {
		return false
	}
	return tracelog.Unwrap(err) == gorm.ErrRecordNotFound
}
