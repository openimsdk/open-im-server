package logic

import (
	"Open_IM/src/common/config"
	tpns "Open_IM/src/push/sdk/tpns-server-sdk-go/go"
	"Open_IM/src/push/sdk/tpns-server-sdk-go/go/auth"
	"Open_IM/src/push/sdk/tpns-server-sdk-go/go/common"
	"Open_IM/src/push/sdk/tpns-server-sdk-go/go/req"
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
