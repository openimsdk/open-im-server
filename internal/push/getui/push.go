package getui

import (
	"Open_IM/internal/push"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/utils"
	"bytes"
	"crypto/sha256"
	"errors"

	//"crypto/sha512"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var (
	GetuiClient *Getui

	TokenExpireError = errors.New("token expire")
)

const (
	PushURL      = "/push/single/alias"
	AuthURL      = "/auth"
	TaskURL      = "/push/list/message"
	BatchPushURL = "/push/list/alias"
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

type TaskResp struct {
	TaskID string `json:"taskID"`
}

type Settings struct {
	TTL *int64 `json:"ttl"`
}

type Audience struct {
	Alias []string `json:"alias"`
}

type PushMessage struct {
	Notification *Notification `json:"notification,omitempty"`
	Transmission *string       `json:"transmission,omitempty"`
}

type PushChannel struct {
	Ios     *Ios     `json:"ios"`
	Android *Android `json:"android"`
}

type PushReq struct {
	RequestID   *string      `json:"request_id"`
	Settings    *Settings    `json:"settings"`
	Audience    *Audience    `json:"audience"`
	PushMessage *PushMessage `json:"push_message"`
	PushChannel *PushChannel `json:"push_channel"`
	IsAsync     *bool        `json:"is_async"`
	Taskid      *string      `json:"taskid"`
}

type Ios struct {
	NotiType  *string `json:"type"`
	AutoBadge *string `json:"auto_badge"`
	Aps       struct {
		Sound string `json:"sound"`
		Alert Alert  `json:"alert"`
	} `json:"aps"`
}

type Alert struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

type Android struct {
	Ups struct {
		Notification Notification `json:"notification"`
		Options      Options      `json:"options"`
	} `json:"ups"`
}

type Notification struct {
	Title       string `json:"title"`
	Body        string `json:"body"`
	ChannelID   string `json:"channel_id"`
	ChannelName string `json:"channel_name"`
	ClickType   string `json:"click_type"`
}

type Options struct {
	HW *HW `json:"HW"`
	XM *XM `json:"XM"`
	VV *VV `json:"VV"`
	OP *OP `json:"OP"`
	HO *HO `json:"HO"`
}

type HW struct {
	Category     string `json:"/message/android/category"`
	DefaultSound bool   `json:"/message/android/notification/default_sound"`
	ChannelID    string `json:"/message/android/notification/channel_id"`
	Sound        string `json:"/message/android/notification/sound"`
	Importance   string `json:"/message/android/notification/importance"`
}

type XM struct {
	ChannelID string `json:"/extra.channel_id"`
}

type VV struct {
	Classification int `json:"/classification"`
}

type OP struct {
	ChannelID string `json:"/channel_id"`
}

type HO struct {
	Importance string `json:"/android/notification/importance"`
}

type PushResp struct {
}

func newGetuiClient() *Getui {
	return &Getui{}
}

func (g *Getui) Push(userIDList []string, title, detailContent, operationID string, opts push.PushOpts) (resp string, err error) {
	token, err := db.DB.GetGetuiToken()
	log.NewDebug(operationID, utils.GetSelfFuncName(), "token：", token, userIDList)
	if err != nil {
		log.NewError(operationID, utils.GetSelfFuncName(), "GetGetuiToken failed", err.Error())
	}
	if token == "" || err != nil {
		token, err = g.getTokenAndSave2Redis(operationID)
		if err != nil {
			log.NewError(operationID, utils.GetSelfFuncName(), "getTokenAndSave2Redis failed", err.Error())
			return "", utils.Wrap(err, "")
		}
	}

	pushReq := PushReq{PushMessage: &PushMessage{Notification: &Notification{
		Title:       title,
		Body:        detailContent,
		ClickType:   "startapp",
		ChannelID:   config.Config.Push.Getui.ChannelID,
		ChannelName: config.Config.Push.Getui.ChannelName,
	}}}
	pushReq.setPushChannel(title, detailContent)
	pushResp := PushResp{}
	if len(userIDList) > 1 {
		taskID, err := g.GetTaskID(operationID, token, pushReq)
		if err != nil {
			return "", utils.Wrap(err, "GetTaskIDAndSave2Redis failed")
		}
		pushReq = PushReq{Audience: &Audience{Alias: userIDList}}
		var IsAsync = true
		pushReq.IsAsync = &IsAsync
		pushReq.Taskid = &taskID
		err = g.request(BatchPushURL, pushReq, token, &pushResp, operationID)
	} else {
		reqID := utils.OperationIDGenerator()
		pushReq.RequestID = &reqID
		pushReq.Audience = &Audience{Alias: []string{userIDList[0]}}
		err = g.request(PushURL, pushReq, token, &pushResp, operationID)
	}
	switch err {
	case TokenExpireError:
		token, err = g.getTokenAndSave2Redis(operationID)
		if err != nil {
			log.NewError(operationID, utils.GetSelfFuncName(), "getTokenAndSave2Redis failed, ", err.Error())
		} else {
			log.NewInfo(operationID, utils.GetSelfFuncName(), "getTokenAndSave2Redis: ", token)
		}
	}
	if err != nil {
		return "", utils.Wrap(err, "push failed")
	}
	respBytes, err := json.Marshal(pushResp)
	return string(respBytes), utils.Wrap(err, "")
}

func (g *Getui) Auth(operationID string, timeStamp int64) (token string, expireTime int64, err error) {
	log.NewInfo(operationID, utils.GetSelfFuncName(), config.Config.Push.Getui.AppKey, timeStamp, config.Config.Push.Getui.MasterSecret)
	h := sha256.New()
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
	err = g.request(AuthURL, reqAuth, "", &respAuth, operationID)
	if err != nil {
		return "", 0, err
	}
	log.NewInfo(operationID, utils.GetSelfFuncName(), "result: ", respAuth)
	expire, err := strconv.Atoi(respAuth.ExpireTime)
	return respAuth.Token, int64(expire), err
}

func (g *Getui) GetTaskID(operationID, token string, pushReq PushReq) (string, error) {
	respTask := TaskResp{}
	ttl := int64(1000 * 60 * 5)
	pushReq.Settings = &Settings{TTL: &ttl}
	err := g.request(TaskURL, pushReq, token, &respTask, operationID)
	if err != nil {
		return "", utils.Wrap(err, "")
	}
	return respTask.TaskID, nil
}

func (g *Getui) request(url string, content interface{}, token string, returnStruct interface{}, operationID string) error {
	con, err := json.Marshal(content)
	if err != nil {
		return err
	}
	client := &http.Client{}
	log.Debug(operationID, utils.GetSelfFuncName(), "json:", string(con), "token:", token)
	req, err := http.NewRequest("POST", config.Config.Push.Getui.PushUrl+url, bytes.NewBuffer(con))
	if err != nil {
		return err
	}
	if token != "" {
		req.Header.Set("token", token)
	}
	req.Header.Set("content-type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	result, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.NewDebug(operationID, "getui", utils.GetSelfFuncName(), "resp, ", string(result))
	commonResp := GetuiCommonResp{}
	commonResp.Data = returnStruct
	if err := json.Unmarshal(result, &commonResp); err != nil {
		return err
	}
	if commonResp.Code == 10001 {
		return TokenExpireError
	}
	return nil
}

func (pushReq *PushReq) setPushChannel(title string, body string) {
	pushReq.PushChannel = &PushChannel{}
	autoBadge := "+1"
	pushReq.PushChannel.Ios = &Ios{AutoBadge: &autoBadge}
	notify := "notify"
	pushReq.PushChannel.Ios.NotiType = &notify
	pushReq.PushChannel.Ios.Aps.Sound = "default"
	pushReq.PushChannel.Ios.Aps.Alert = Alert{
		Title: title,
		Body:  body,
	}
	pushReq.PushChannel.Android = &Android{}
	pushReq.PushChannel.Android.Ups.Notification = Notification{
		Title:     title,
		Body:      body,
		ClickType: "startapp",
	}
	pushReq.PushChannel.Android.Ups.Options = Options{
		HW: &HW{Category: config.Config.Push.Getui.Channel.HW.Category, ChannelID: "RingRing4", Sound: "/raw/ring001", Importance: "NORMAL"},
		HO: &HO{Importance: "NORMAL"},
		XM: &XM{ChannelID: config.Config.Push.Getui.Channel.XM.ChannelID},
		OP: &OP{ChannelID: config.Config.Push.Getui.Channel.OPPO.ChannelID},
		VV: &VV{Classification: 1},
	}
}

func (g *Getui) getTokenAndSave2Redis(operationID string) (token string, err error) {
	token, expireTime, err := g.Auth(operationID, time.Now().UnixNano()/1e6)
	if err != nil {
		return "", utils.Wrap(err, "Auth failed")
	}
	log.NewDebug(operationID, "getui", utils.GetSelfFuncName(), token, expireTime, err)
	err = db.DB.SetGetuiToken(token, 60*60*23)
	if err != nil {
		return "", utils.Wrap(err, "Auth failed")
	}
	return token, nil
}

func (g *Getui) GetTaskIDAndSave2Redis(operationID, token string, pushReq PushReq) (taskID string, err error) {
	ttl := int64(1000 * 60 * 60 * 24)
	pushReq.Settings = &Settings{TTL: &ttl}
	taskID, err = g.GetTaskID(operationID, token, pushReq)
	if err != nil {
		return "", utils.Wrap(err, "GetTaskIDAndSave2Redis failed")
	}
	err = db.DB.SetGetuiTaskID(taskID, 60*60*23)
	if err != nil {
		return "", utils.Wrap(err, "Auth failed")
	}
	return token, nil
}
