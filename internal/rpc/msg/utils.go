package msg

import (
	"Open_IM/pkg/utils"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func IsNotFound(err error) bool {
	switch utils.Unwrap(err) {
	case gorm.ErrRecordNotFound, redis.Nil:
		return true
	default:
		return false
	}
}
