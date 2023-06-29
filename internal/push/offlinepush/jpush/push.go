package jpush

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/OpenIMSDK/Open-IM-Server/internal/push/offlinepush"
	"github.com/OpenIMSDK/Open-IM-Server/internal/push/offlinepush/jpush/body"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	http2 "github.com/OpenIMSDK/Open-IM-Server/pkg/common/http"
)

type JPush struct{}

func NewClient() *JPush {
	return &JPush{}
}

func (j *JPush) Auth(apiKey, secretKey string, timeStamp int64) (token string, err error) {
	return token, nil
}

func (j *JPush) SetAlias(cid, alias string) (resp string, err error) {
	return resp, nil
}

func (j *JPush) getAuthorization(appKey string, masterSecret string) string {
	str := fmt.Sprintf("%s:%s", appKey, masterSecret)
	buf := []byte(str)
	Authorization := fmt.Sprintf("Basic %s", base64.StdEncoding.EncodeToString(buf))
	return Authorization
}

func (j *JPush) Push(ctx context.Context, userIDs []string, title, content string, opts *offlinepush.Opts) error {
	var pf body.Platform
	pf.SetAll()
	var au body.Audience
	au.SetAlias(userIDs)
	var no body.Notification
	var extras body.Extras
	if opts.Signal.ClientMsgID != "" {
		extras.ClientMsgID = opts.Signal.ClientMsgID
	}
	no.IOSEnableMutableContent()
	no.SetExtras(extras)
	no.SetAlert(title)
	var msg body.Message
	msg.SetMsgContent(content)
	var opt body.Options
	opt.SetApnsProduction(config.Config.IOSPush.Production)
	var pushObj body.PushObj
	pushObj.SetPlatform(&pf)
	pushObj.SetAudience(&au)
	pushObj.SetNotification(&no)
	pushObj.SetMessage(&msg)
	pushObj.SetOptions(&opt)
	var resp interface{}
	return j.request(ctx, pushObj, resp, 5)
}

func (j *JPush) request(ctx context.Context, po body.PushObj, resp interface{}, timeout int) error {
	return http2.PostReturn(ctx, config.Config.Push.Jpns.PushUrl, map[string]string{"Authorization": j.getAuthorization(config.Config.Push.Jpns.AppKey, config.Config.Push.Jpns.MasterSecret)}, po, resp, timeout)
}
