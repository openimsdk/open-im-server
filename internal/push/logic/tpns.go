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

package logic

import (
	tpns "Open_IM/internal/push/sdk/tpns-server-sdk-go/go"
	"Open_IM/internal/push/sdk/tpns-server-sdk-go/go/auth"
	"Open_IM/internal/push/sdk/tpns-server-sdk-go/go/common"
	"Open_IM/internal/push/sdk/tpns-server-sdk-go/go/req"
	"Open_IM/pkg/common/config"
)

var badgeType = -2
var iosAcceptId = auth.Auther{AccessID: config.Config.Push.Tpns.Ios.AccessID, SecretKey: config.Config.Push.Tpns.Ios.SecretKey}

func IOSAccountListPush(accounts []string, title, content, jsonCustomContent string) {
	var iosMessage = tpns.Message{
		Title:   title,
		Content: content,
		IOS: &tpns.IOSParams{
			Aps: &tpns.Aps{
				BadgeType: &badgeType,
				Sound:     "default",
				Category:  "INVITE_CATEGORY",
			},
			CustomContent: jsonCustomContent,
			//CustomContent: `"{"key\":\"value\"}"`,
		},
	}
	pushReq, reqBody, err := req.NewListAccountPush(accounts, iosMessage)
	if err != nil {
		return
	}
	iosAcceptId.Auth(pushReq, auth.UseSignAuthored, iosAcceptId, reqBody)
	common.PushAndGetResult(pushReq)
}
