package msg

import (
	"strconv"
	"strings"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/protocol/constant"
	pbchat "github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/tools/utils/datautil"
)

const (
	separator = "-"
)

func filterAfterMsg(msg *pbchat.SendMsgReq, after *config.AfterConfig) bool {
	return filterMsg(msg, after.AttentionIds, after.DeniedTypes)
}

func filterBeforeMsg(msg *pbchat.SendMsgReq, before *config.BeforeConfig) bool {
	return filterMsg(msg, nil, before.DeniedTypes)
}

func filterMsg(msg *pbchat.SendMsgReq, attentionIds []string, deniedTypes []int32) bool {
	// According to the attentionIds configuration, only some users are sent
	if len(attentionIds) != 0 && !datautil.Contain(msg.MsgData.RecvID, attentionIds...) {
		return false
	}

	if defaultDeniedTypes(msg.MsgData.ContentType) {
		return false
	}

	if len(deniedTypes) != 0 && datautil.Contain(msg.MsgData.ContentType, deniedTypes...) {
		return false
	}
	//if len(allowedTypes) != 0 && !isInInterval(msg.MsgData.ContentType, allowedTypes) {
	//	return false
	//}
	//if len(deniedTypes) != 0 && isInInterval(msg.MsgData.ContentType, deniedTypes) {
	//	return false
	//}
	return true
}

func defaultDeniedTypes(contentType int32) bool {
	if contentType >= constant.NotificationBegin && contentType <= constant.NotificationEnd {
		return true
	}
	if contentType == constant.Typing {
		return true
	}
	return false
}

// isInInterval if data is in interval
// Supports two formats: a single type or a range. The range is defined by the lower and upper bounds connected with a hyphen ("-")
// e.g. [1, 100, 200-500, 600-700] means that only data within the range
// {1, 100} ∪ [200, 500] ∪ [600, 700] will return true.
func isInInterval(data int32, interval []string) bool {
	for _, v := range interval {
		if strings.Contains(v, separator) {
			// is interval
			bounds := strings.Split(v, separator)
			if len(bounds) != 2 {
				continue
			}
			bottom, err := strconv.Atoi(bounds[0])
			if err != nil {
				continue
			}
			top, err := strconv.Atoi(bounds[1])
			if err != nil {
				continue
			}
			if datautil.BetweenEq(int(data), bottom, top) {
				return true
			}
		} else {
			iv, err := strconv.Atoi(v)
			if err != nil {
				continue
			}
			if int(data) == iv {
				return true
			}
		}
	}
	return false
}
