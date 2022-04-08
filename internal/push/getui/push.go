package getui

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/utils"
	"bytes"
	"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var (
	GetuiClient *Getui
)

func init() {
	GetuiClient = newGetuiClient()
}

type Getui struct{}

type GetuiCommonResp struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

type AuthReq struct {
	Sign      string `json:"sign"`
	Timestamp string `json:"timestamp"`
	Appkey    string `json:"appkey"`
}

type AuthResp struct {
	ExpireTime string `json:"expire_time"`
	Token      string `json:"token"`
}

type PushReq struct {
	RequestID string `json:"request_id"`
	Audience  struct {
		Cid []string `json:"cid"`
	} `json:"audience"`
	PushMssage struct {
		Notification Notification `json:"notification"`
	} `json:"push_message"`
}

type Notification struct {
	Title     string `json:"title"`
	Body      string `json:"body"`
	ClickType string `json:"click_type"`
	Intent    string `json:"intent"`
	Url       string `json:"url"`
}

type PushResp struct {
}

func newGetuiClient() *Getui {
	return &Getui{}
}

func (g *Getui) Push(userIDList []string, alert, detailContent, platform, operationID string) (resp string, err error) {
	token, err := db.DB.GetGetuiToken()
	if err != nil {
		log.NewError(operationID, utils.OperationIDGenerator(), "GetGetuiToken", err.Error())
	}
	if token == "" || err != nil {
		token, expireTime, err := g.Auth(operationID, time.Now().Unix())
		if err != nil {
			return "", utils.Wrap(err, "Auth failed")
		}
		log.NewInfo(operationID, "getui", utils.GetSelfFuncName(), token, expireTime, err)
		err = db.DB.SetGetuiToken(token, expireTime-time.Now().Unix()-20)
		if err != nil {
			return "", utils.Wrap(err, "Auth failed")
		}
	}
	pushReq := PushReq{
		RequestID: utils.OperationIDGenerator(),
		Audience: struct {
			Cid []string `json:"cid"`
		}{Cid: []string{userIDList[0]}},
	}
	pushReq.PushMssage.Notification = Notification{
		Title:     alert,
		Body:      alert,
		ClickType: "none",
	}
	if config.Config.Push.Getui.Intent != "" {
		pushReq.PushMssage.Notification.Intent = config.Config.Push.Getui.Intent
		pushReq.PushMssage.Notification.ClickType = "intent"
	}
	pushResp := PushResp{}
	err = g.request(pushReq, token, &pushResp, operationID)
	if err != nil {
		return "", utils.Wrap(err, "")
	}
	respBytes, err := json.Marshal(pushResp)
	return string(respBytes), err
}

func (g *Getui) Auth(operationID string, timeStamp int64) (token string, expireTime int64, err error) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), timeStamp)
	h := sha512.New()
	h.Write([]byte(config.Config.Push.Getui.AppKey + strconv.Itoa(int(timeStamp)) + config.Config.Push.Getui.MasterSecret))
	sum := h.Sum(nil)
	sign := hex.EncodeToString(sum)
	log.NewInfo(operationID, utils.GetSelfFuncName(), "sha256 result", sign)
	reqAuth := AuthReq{
		Sign:      sign,
		Timestamp: strconv.Itoa(int(timeStamp)),
		Appkey:    config.Config.Push.Getui.AppKey,
	}
	respAuth := AuthResp{}
	err = g.request(reqAuth, "", &respAuth, operationID)
	expire, err := strconv.Atoi(respAuth.ExpireTime)
	return respAuth.Token, int64(expire), err
}

func (g *Getui) request(content interface{}, token string, returnStruct interface{}, operationID string) error {
	con, err := json.Marshal(content)
	if err != nil {
		return err
	}
	client := &http.Client{}
	req, err := http.NewRequest("POST", config.Config.Push.Getui.PushUrl, bytes.NewBuffer(con))
	if err != nil {
		return err
	}
	req.Header.Set(token, token)
	req.Header.Set("content-type", "application/json")
	req.Header.Set("content-type", "charset=utf-8")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.NewInfo(operationID, "getui", utils.GetSelfFuncName(), string(result))
	if err := json.Unmarshal(result, returnStruct); err != nil {
		return err
	}
	return nil
}
