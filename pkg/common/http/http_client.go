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
	"io/ioutil"
	"net/http"
	urlLib "net/url"
	"time"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/callbackstruct"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
)

var client http.Client

func Get(url string) (response []byte, err error) {
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func Post(
	ctx context.Context,
	url string,
	header map[string]string,
	data interface{},
	timeout int,
) (content []byte, err error) {
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

func PostReturn(
	ctx context.Context,
	url string,
	header map[string]string,
	input, output interface{},
	timeOutSecond int,
) error {
	b, err := Post(ctx, url, header, input, timeOutSecond)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, output)
	return err
}

func callBackPostReturn(
	ctx context.Context,
	url, command string,
	input interface{},
	output callbackstruct.CallbackResp,
	callbackConfig config.CallBackConfig,
) error {
	defer log.ZDebug(ctx, "callback", "url", url, "command", command, "input", input, "callbackConfig", callbackConfig)
	v := urlLib.Values{}
	v.Set(constant.CallbackCommand, command)
	url = url + "?" + v.Encode()
	b, err := Post(ctx, url, nil, input, callbackConfig.CallbackTimeOut)
	if err != nil {
		if callbackConfig.CallbackFailedContinue != nil && *callbackConfig.CallbackFailedContinue {
			log.ZWarn(ctx, "callback failed but continue", err, "url", url)
			return errs.ErrCallbackContinue
		}
		return errs.ErrNetwork.Wrap(err.Error())
	}
	if err = json.Unmarshal(b, output); err != nil {
		if callbackConfig.CallbackFailedContinue != nil && *callbackConfig.CallbackFailedContinue {
			log.ZWarn(ctx, "callback failed but continue", err, "url", url)
			return errs.ErrCallbackContinue
		}
		return errs.ErrData.Wrap(err.Error())
	}
	return output.Parse()
}

func CallBackPostReturn(
	ctx context.Context,
	url string,
	req callbackstruct.CallbackReq,
	resp callbackstruct.CallbackResp,
	callbackConfig config.CallBackConfig,
) error {
	return callBackPostReturn(ctx, url, req.GetCallbackCommand(), req, resp, callbackConfig)
}
