package msg

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/proto/sdkws"
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
