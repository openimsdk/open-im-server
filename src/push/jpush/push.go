package push

import (
	"Open_IM/src/common/config"
	"Open_IM/src/push/jpush/common"
	"Open_IM/src/push/jpush/requestBody"
	"bytes"
	"encoding/json"
	"net/http"
)

func JGAccountListPush(accounts []string, jsonCustomContent string, Platform string) (*http.Response, error) {

	var pf requestBody.Platform
	_ = pf.SetAndroid()
	var au requestBody.Audience
	au.SetAlias(accounts)
	var no requestBody.Notification
	no.SetAlert(jsonCustomContent)
	var me requestBody.Message
	me.SetMsgContent(jsonCustomContent)
	var po requestBody.PushObj
	po.SetPlatform(&pf)
	po.SetAudience(&au)
	po.SetNotification(&no)
	po.SetMessage(&me)

	con, err := json.Marshal(po)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	req, err := http.NewRequest("POST", common.PushUrl, bytes.NewBuffer(con))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", common.GetAuthorization(config.Config.Push.Jpns.AppKey, config.Config.Push.Jpns.MasterSecret))

	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return resp, nil
}
