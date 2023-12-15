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
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"time"

	"github.com/OpenIMSDK/protocol/constant"
	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/log"

	"github.com/openimsdk/open-im-server/v3/pkg/callbackstruct"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
)

var (
	// define http client.
	client = &http.Client{
		Timeout: 15 * time.Second, // max timeout is 15s
	}
)

func init() {
	// reset http default transport
	http.DefaultTransport.(*http.Transport).MaxConnsPerHost = 100 // default: 2
}

func Get(url string) (response []byte, err error) {
	hclient := http.Client{Timeout: 5 * time.Second}
	resp, err := hclient.Get(url)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func Post(ctx context.Context, url string, header map[string]string, data any, timeout int) (content []byte, err error) {
	if timeout > 0 {
		var cancel func()
		ctx, cancel = context.WithTimeout(ctx, time.Second*time.Duration(timeout))
		defer cancel()
	}

	jsonStr, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}

	if operationID, _ := ctx.Value(constant.OperationID).(string); operationID != "" {
		req.Header.Set(constant.OperationID, operationID)
	}
	for k, v := range header {
		req.Header.Set(k, v)
	}
	req.Header.Add("content-type", "application/json; charset=utf-8")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func PostReturn(ctx context.Context, url string, header map[string]string, input, output any, timeOutSecond int) error {
	b, err := Post(ctx, url, header, input, timeOutSecond)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, output)
	return err
}

func callBackPostReturn(ctx context.Context, url, command string, input interface{}, output callbackstruct.CallbackResp, callbackConfig config.CallBackConfig) error {
	defer log.ZDebug(ctx, "callback", "url", url, "command", command, "input", input, "output", output, "callbackConfig", callbackConfig)
	//
	//v := urllib.Values{}
	//v.Set(constant.CallbackCommand, command)
	//url = url + "/" + v.Encode()
	url = url + "/" + command

	b, err := Post(ctx, url, nil, input, callbackConfig.CallbackTimeOut)
	if err != nil {
		if callbackConfig.CallbackFailedContinue != nil && *callbackConfig.CallbackFailedContinue {
			log.ZWarn(ctx, "callback failed but continue", err, "url", url)
			return nil
		}
		return errs.ErrNetwork.Wrap(err.Error())
	}

	if err = json.Unmarshal(b, output); err != nil {
		if callbackConfig.CallbackFailedContinue != nil && *callbackConfig.CallbackFailedContinue {
			log.ZWarn(ctx, "callback failed but continue", err, "url", url)
			return nil
		}
		return errs.ErrData.WithDetail(err.Error() + "response format error")
	}

	return output.Parse()
}

func CallBackPostReturn(ctx context.Context, url string, req callbackstruct.CallbackReq, resp callbackstruct.CallbackResp, callbackConfig config.CallBackConfig) error {
	return callBackPostReturn(ctx, url, req.GetCallbackCommand(), req, resp, callbackConfig)
}
