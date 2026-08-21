package webhook

import (
	"strconv"
	"strings"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/tools/utils/datautil"
)

func IsContentTypeInIntervals(contentType int32, intervals []string) bool {
	for _, interval := range intervals {
		if strings.Contains(interval, "-") {
			bounds := strings.Split(interval, "-")
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
			contentTypeValue, err := strconv.Atoi(interval)
			if err != nil {
				continue
			}
			if int(contentType) == contentTypeValue {
				return true
			}
		}
	}
	return false
}

func FilterBeforeMsg(msg *sdkws.MsgData, before *config.BeforeConfig) bool {
	return filterContentType(msg.ContentType, before.AllowedTypes, before.DeniedTypes)
}

func FilterAfterMsg(msg *sdkws.MsgData, after *config.AfterConfig, attentionTargetID string) bool {
	if len(after.AttentionIds) != 0 && !datautil.Contains([]string{msg.SendID, attentionTargetID}, after.AttentionIds...) {
		return false
	}
	return filterContentType(msg.ContentType, after.AllowedTypes, after.DeniedTypes)
}

func filterContentType(contentType int32, allowedTypes, deniedTypes []string) bool {
	if len(allowedTypes) != 0 && !IsContentTypeInIntervals(contentType, allowedTypes) {
		return false
	}
	if len(deniedTypes) != 0 && IsContentTypeInIntervals(contentType, deniedTypes) {
		return false
	}
	return true
}
