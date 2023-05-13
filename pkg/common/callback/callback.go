// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

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
