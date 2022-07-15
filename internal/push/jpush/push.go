package push

import (
	"Open_IM/internal/push"
	"Open_IM/internal/push/jpush/common"
	"Open_IM/internal/push/jpush/requestBody"
	"Open_IM/pkg/common/config"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
)

var (
	JPushClient *JPush
)

func init() {
	JPushClient = newGetuiClient()
}

type JPush struct{}

func newGetuiClient() *JPush {
	return &JPush{}
}

func (j *JPush) Auth(apiKey, secretKey string, timeStamp int64) (token string, err error) {
	return token, nil
}

func (j *JPush) SetAlias(cid, alias string) (resp string, err error) {
	return resp, nil
}

func (j *JPush) Push(accounts []string, alert, detailContent, operationID string, opts push.PushOpts) (string, error) {

	var pf requestBody.Platform
	pf.SetAll()
	var au requestBody.Audience
	au.SetAlias(accounts)
	var no requestBody.Notification

	var extras requestBody.Extras
	if opts.Signal.ClientMsgID != "" {
		extras.ClientMsgID = opts.Signal.ClientMsgID
	}
	no.IOSEnableMutableContent()
	no.SetExtras(extras)
	no.SetAlert(alert)
	var me requestBody.Message
	me.SetMsgContent(detailContent)
	var o requestBody.Options
	o.SetApnsProduction(config.Config.IOSPush.Production)
	var po requestBody.PushObj
	po.SetPlatform(&pf)
	po.SetAudience(&au)
	po.SetNotification(&no)
	po.SetMessage(&me)
	po.SetOptions(&o)

	con, err := json.Marshal(po)
	if err != nil {
		return "", err
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", config.Config.Push.Jpns.PushUrl, bytes.NewBuffer(con))
	if err != nil {
		return "", err
	}
	req.Header.Set("Authorization", common.GetAuthorization(config.Config.Push.Jpns.AppKey, config.Config.Push.Jpns.MasterSecret))

	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(result), nil
}
