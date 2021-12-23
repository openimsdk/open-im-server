package push

import (
	"Open_IM/internal/push/jpush/common"
	"Open_IM/internal/push/jpush/requestBody"
	"Open_IM/pkg/common/config"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

type JPushResp struct {
}

func JGAccountListPush(accounts []string, alert, detailContent, platform string) ([]byte, error) {

	var pf requestBody.Platform
	_ = pf.SetPlatform(platform)
	var au requestBody.Audience
	au.SetAlias(accounts)
	var no requestBody.Notification
	no.SetAlert(alert, platform)
	var me requestBody.Message
	me.SetMsgContent(detailContent)
	var o requestBody.Options
	o.SetApnsProduction(false)
	var po requestBody.PushObj
	po.SetPlatform(&pf)
	po.SetAudience(&au)
	po.SetNotification(&no)
	po.SetMessage(&me)
	po.SetOptions(&o)

	con, err := json.Marshal(po)
	if err != nil {
		return nil, err
	}

	client := &http.Client{}

	req, err := http.NewRequest("POST", config.Config.Push.Jpns.PushUrl, bytes.NewBuffer(con))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", common.GetAuthorization(config.Config.Push.Jpns.AppKey, config.Config.Push.Jpns.MasterSecret))

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return result, nil
}
