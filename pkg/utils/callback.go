package utils

import (
	"OpenIM/pkg/common/constant"
	sdkws "OpenIM/pkg/proto/sdkws"
	"github.com/golang/protobuf/proto"
)

func GetContent(msg *sdkws.MsgData) string {
	if msg.ContentType >= constant.NotificationBegin && msg.ContentType <= constant.NotificationEnd {
		var tips sdkws.TipsComm
		_ = proto.Unmarshal(msg.Content, &tips)
		content := tips.JsonDetail
		return content
	} else {
		return string(msg.Content)
	}
}
