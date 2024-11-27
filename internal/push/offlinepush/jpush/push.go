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

package jpush

import (
	"context"
	"encoding/base64"
	"fmt"

	"github.com/openimsdk/tools/utils/httputil"

	"github.com/openimsdk/open-im-server/v3/internal/push/offlinepush/jpush/body"
	"github.com/openimsdk/open-im-server/v3/internal/push/offlinepush/options"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
)

type JPush struct {
	pushConf   *config.Push
	httpClient *httputil.HTTPClient
}

func NewClient(pushConf *config.Push) *JPush {
	return &JPush{pushConf: pushConf,
		httpClient: httputil.NewHTTPClient(httputil.NewClientConfig()),
	}
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

func (j *JPush) Push(ctx context.Context, userIDs []string, title, content string, opts *options.Opts) error {
	var pf body.Platform
	pf.SetAll()
	var au body.Audience
	au.SetAlias(userIDs)
	var no body.Notification
	extras := make(map[string]string)
	extras["ex"] = opts.Ex
	if opts.Signal.ClientMsgID != "" {
		extras["ClientMsgID"] = opts.Signal.ClientMsgID
	}
	no.IOSEnableMutableContent()
	no.SetExtras(extras)
	no.SetAlert(content, title, opts)
	no.SetAndroidIntent(j.pushConf)

	var msg body.Message
	msg.SetMsgContent(content)
	msg.SetTitle(title)
	if opts.Signal.ClientMsgID != "" {
		msg.SetExtras("ClientMsgID", opts.Signal.ClientMsgID)
	}
	msg.SetExtras("ex", opts.Ex)
	var opt body.Options
	opt.SetApnsProduction(j.pushConf.IOSPush.Production)
	var pushObj body.PushObj
	pushObj.SetPlatform(&pf)
	pushObj.SetAudience(&au)
	pushObj.SetNotification(&no)
	pushObj.SetMessage(&msg)
	pushObj.SetOptions(&opt)
	var resp map[string]any
	return j.request(ctx, pushObj, &resp, 5)
}

func (j *JPush) request(ctx context.Context, po body.PushObj, resp *map[string]any, timeout int) error {
	err := j.httpClient.PostReturn(
		ctx,
		j.pushConf.JPush.PushURL,
		map[string]string{
			"Authorization": j.getAuthorization(j.pushConf.JPush.AppKey, j.pushConf.JPush.MasterSecret),
		},
		po,
		resp,
		timeout,
	)
	if err != nil {
		return err
	}
	if (*resp)["sendno"] != "0" {
		return fmt.Errorf("jpush push failed %v", resp)
	}
	return nil
}
