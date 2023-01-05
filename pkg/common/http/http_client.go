/*
** description("").
** copyright('open-im,www.open-im.io').
** author("fg,Gordon@tuoyun.net").
** time(2021/5/27 10:31).
 */
package http

import (
	cbApi "Open_IM/pkg/call_back_struct"
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

//application/json; charset=utf-8
func Post(url string, data interface{}, timeOutSecond int) (content []byte, err error) {
	jsonStr, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonStr))
	if err != nil {
		return nil, err
	}
	req.Close = true
	req.Header.Add("content-type", "application/json; charset=utf-8")

	client := &http.Client{Timeout: time.Duration(timeOutSecond) * time.Second}
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

func CallBackPostReturn(url, callbackCommand string, input interface{}, output cbApi.CallbackResp, timeOut int, failedContinue *bool) (bool, error) {
	v := urlLib.Values{}
	v.Set("callbackCommand", callbackCommand)
	url = url + "?" + v.Encode()
	b, err := Post(url, input, timeOut)
	if err != nil {
		if failedContinue != nil {
			return *failedContinue, err
		}
		return true, err
	}
	if err = json.Unmarshal(b, output); err != nil {
		if failedContinue != nil {
			return *failedContinue, err
		}
		return true, err
	}
	return output.Parse()
}
