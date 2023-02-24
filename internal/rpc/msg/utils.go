package msg

import (
	"OpenIM/pkg/common/config"
	"OpenIM/pkg/common/constant"
	"OpenIM/pkg/proto/sdkws"
	"OpenIM/pkg/utils"
	"github.com/go-redis/redis/v8"
	"gorm.io/gorm"
)

func isMessageHasReadEnabled(msgData *sdkws.MsgData) bool {
	switch msgData.ContentType {
	case constant.HasReadReceipt:
		if config.Config.SingleMessageHasReadReceiptEnable {
			return true
		} else {
			return false
		}
	case constant.GroupHasReadReceipt:
		if config.Config.GroupMessageHasReadReceiptEnable {
			return true
		} else {
			return false
		}
	}
	return true
}

func IsNotFound(err error) bool {
	switch utils.Unwrap(err) {
	case redis.Nil, gorm.ErrRecordNotFound:
		return true
	default:
		return false
	}
}
