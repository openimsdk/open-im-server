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

package http

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/openimsdk/open-im-server/v3/pkg/callbackstruct"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/servererrs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/utils/httputil"
)

var (
	// Define http client.
	client = httputil.NewHTTPClient(httputil.NewClientConfig())
)

func init() {
	// reset http default transport
	http.DefaultTransport.(*http.Transport).MaxConnsPerHost = 100 // default: 2
}

func callBackPostReturn(ctx context.Context, url, command string, input interface{}, output callbackstruct.CallbackResp, callbackConfig config.WebhookConfig) error {
	url = url + "/" + command
	log.ZInfo(ctx, "callback", "url", url, "input", input, "config", callbackConfig)
	b, err := client.Post(ctx, url, nil, input, callbackConfig.Timeout)
	if err != nil {
		if callbackConfig.FailedContinue {
			log.ZInfo(ctx, "callback failed but continue", err, "url", url)
			return nil
		}
		log.ZWarn(ctx, "callback network failed", err, "url", url, "input", input)
		return servererrs.ErrNetwork.WrapMsg(err.Error())
	}
	if err = json.Unmarshal(b, output); err != nil {
		if callbackConfig.FailedContinue {
			log.ZWarn(ctx, "callback failed but continue", err, "url", url)
			return nil
		}
		log.ZWarn(ctx, "callback json unmarshal failed", err, "url", url, "input", input, "response", string(b))
		return servererrs.ErrData.WithDetail(err.Error() + "response format error")
	}
	if err := output.Parse(); err != nil {
		log.ZWarn(ctx, "callback parse failed", err, "url", url, "input", input, "response", string(b))
	}
	log.ZInfo(ctx, "callback success", "url", url, "input", input, "response", string(b))
	return nil
}

func CallBackPostReturn(ctx context.Context, url string, req callbackstruct.CallbackReq, resp callbackstruct.CallbackResp, callbackConfig config.WebhookConfig) error {
	return callBackPostReturn(ctx, url, req.GetCallbackCommand(), req, resp, callbackConfig)
}
