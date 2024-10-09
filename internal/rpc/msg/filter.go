package msg

import (
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	pbchat "github.com/openimsdk/protocol/msg"
	"github.com/openimsdk/tools/utils/datautil"
	"strconv"
	"strings"
)

const (
	separator = "-"
)

func filterAfterMsg(msg *pbchat.SendMsgReq, after *config.AfterConfig) bool {
	return filterMsg(msg, after.AttentionIds, after.AllowedTypes, after.DeniedTypes)
}

func filterBeforeMsg(msg *pbchat.SendMsgReq, before *config.BeforeConfig) bool {
	return filterMsg(msg, nil, before.AllowedTypes, before.DeniedTypes)
}

func filterMsg(msg *pbchat.SendMsgReq, attentionIds, allowedTypes, deniedTypes []string) bool {
	// According to the attentionIds configuration, only some users are sent
	if len(attentionIds) != 0 && !datautil.Contains([]string{msg.MsgData.SendID, msg.MsgData.RecvID}, attentionIds...) {
		return false
	}
	if len(allowedTypes) != 0 && !isInInterval(msg.MsgData.ContentType, allowedTypes) {
		return false
	}
	if len(deniedTypes) != 0 && isInInterval(msg.MsgData.ContentType, deniedTypes) {
		return false
	}
	return true
}

func isInInterval(contentType int32, interval []string) bool {
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
			if datautil.BetweenEq(int(contentType), bottom, top) {
				return true
			}
		} else {
			iv, err := strconv.Atoi(v)
			if err != nil {
				continue
			}
			if int(contentType) == iv {
				return true
			}
		}
	}
	return false
}
