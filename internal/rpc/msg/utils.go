package msg

import (
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

func isMessageHasReadEnabled(msgData *sdkws.MsgData) bool {
	switch {
	case msgData.ContentType == constant.HasReadReceipt && msgData.SessionType == constant.SingleChatType:
		if config.Config.SingleMessageHasReadReceiptEnable {
			return true
		} else {
			return false
		}
	case msgData.ContentType == constant.HasReadReceipt && msgData.SessionType == constant.SuperGroupChatType:
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
