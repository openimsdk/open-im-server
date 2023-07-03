package utils

import (
	"google.golang.org/protobuf/proto"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	sdkws "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
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
