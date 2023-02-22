/*
** description("").
** copyright('open-im,www.open-im.io').
** author("fg,Gordon@tuoyun.net").
** time(2021/5/27 10:31).
 */
package http

import (
	"Open_IM/pkg/callbackstruct"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	urlLib "net/url"
	"time"
)

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

// application/json; charset=utf-8
func Post(url string, header map[string]string, data interface{}, timeout int) (content []byte, err error) {
	jsonStr, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}
	for k, v := range header {
		req.Header.Set(k, v)
	}
	req.Close = true
	req.Header.Add("content-type", "application/json; charset=utf-8")
	client := &http.Client{Timeout: time.Duration(timeout) * time.Second}
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

func PostReturn(url string, header map[string]string, input, output interface{}, timeOutSecond int) error {
	b, err := Post(url, header, input, timeOutSecond)
	if err != nil {
		return err
	}
	err = json.Unmarshal(b, output)
	return err
}

func callBackPostReturn(url, command string, input interface{}, output callbackstruct.CallbackResp, callbackConfig config.CallBackConfig) error {
	v := urlLib.Values{}
	v.Set("callbackCommand", command)
	url = url + "?" + v.Encode()
	b, err := Post(url, nil, input, callbackConfig.CallbackTimeOut)
	if err != nil {
		if callbackConfig.CallbackFailedContinue != nil && *callbackConfig.CallbackFailedContinue {
			return constant.ErrCallbackContinue
		}
		return constant.NewErrNetwork(err)
	}
	if err = json.Unmarshal(b, output); err != nil {
		if callbackConfig.CallbackFailedContinue != nil && *callbackConfig.CallbackFailedContinue {
			return constant.ErrCallbackContinue
		}
		return constant.NewErrData(err)
	}
	return output.Parse()
}

func CallBackPostReturn(url string, req callbackstruct.CallbackReq, resp callbackstruct.CallbackResp, callbackConfig config.CallBackConfig) error {
	return callBackPostReturn(url, req.GetCallbackCommand(), req, resp, callbackConfig)
}
