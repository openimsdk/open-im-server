package common

import (
	tpns "Open_IM/internal/push/sdk/tpns-server-sdk-go/go"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func PushAndGetResult(pushReq *http.Request) {
	c := &http.Client{}
	rsp, err := c.Do(pushReq)
	fmt.Println()
	if err != nil {
		//fmt.Printf("http err:%v", err)
		return
	}
	defer rsp.Body.Close()
	body, err := ioutil.ReadAll(rsp.Body)
	//fmt.Printf("http ReadAll err:%v, body:%v  ", err, string(body))
	if err != nil {
		return
	}
	r := &tpns.CommonRsp{}
	json.Unmarshal(body, r)
	//fmt.Printf("push result: %+v", r)
}

func UploadFile(req *http.Request) (int, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return 0, err
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return 0, err
	}

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("response error, status: %s, body: %s", resp.Status, string(body))
	}

	type uploadResponse struct {
		RetCode  int    `json:"retCode"`
		ErrMsg   string `json:"errMsg"`
		UploadId int    `json:"uploadId"`
	}

	var ur uploadResponse
	if err := json.Unmarshal(body, &ur); err != nil {
		return 0, err
	}

	if ur.RetCode != 0 {
		return 0, fmt.Errorf("response with %d:%s", ur.RetCode, ur.ErrMsg)
	}
	return ur.UploadId, nil
}
