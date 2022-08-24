package callback

import (
	"Open_IM/pkg/common/constant"
	server_api_params "Open_IM/pkg/proto/sdk_ws"
	"github.com/golang/protobuf/proto"
)

func GetContent(msg *server_api_params.MsgData) string {
	if msg.ContentType >= constant.NotificationBegin && msg.ContentType <= constant.NotificationEnd {
		var tips server_api_params.TipsComm
		_ = proto.Unmarshal(msg.Content, &tips)
		//marshaler := jsonpb.Marshaler{
		//	OrigName:     true,
		//	EnumsAsInts:  false,
		//	EmitDefaults: false,
		//}
		content := tips.JsonDetail
		return content
	} else {
		return string(msg.Content)
	}
}
