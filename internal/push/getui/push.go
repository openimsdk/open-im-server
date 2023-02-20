package getui

import (
	"Open_IM/internal/push"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/db/cache"
	//http2 "Open_IM/pkg/common/http"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/utils"
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"time"
)

var (
	TokenExpireError = errors.New("token expire")
)

const (
	pushURL      = "/push/single/alias"
	authURL      = "/auth"
	taskURL      = "/push/list/message"
	batchPushURL = "/push/list/alias"

	tokenExpire = 10001
	ttl         = 0
)

type Client struct {
	cache cache.Cache
}

func newClient(cache cache.Cache) *Client {
	return &Client{cache: cache}
}

func (g *Client) Push(ctx context.Context, userIDs []string, title, content, operationID string, opts *push.Opts) error {
	token, err := g.cache.GetGetuiToken(ctx)
	if err != nil {
		log.NewError(operationID, utils.GetSelfFuncName(), "GetGetuiToken failed", err.Error())
	}
	if token == "" || err != nil {
		token, err = g.getTokenAndSave2Redis(ctx)
		if err != nil {
			log.NewError(operationID, utils.GetSelfFuncName(), "getTokenAndSave2Redis failed", err.Error())
			return utils.Wrap(err, "")
		}
	}
	pushReq := newPushReq(title, content)
	pushReq.setPushChannel(title, content)
	pushResp := struct{}{}
	if len(userIDs) > 1 {
		taskID, err := g.GetTaskID(ctx, token, pushReq)
		if err != nil {
			return utils.Wrap(err, "GetTaskIDAndSave2Redis failed")
		}
		pushReq = PushReq{Audience: &Audience{Alias: userIDs}}
		var IsAsync = true
		pushReq.IsAsync = &IsAsync
		pushReq.TaskID = &taskID
		err = g.request(ctx, batchPushURL, pushReq, token, &pushResp)
	} else {
		reqID := utils.OperationIDGenerator()
		pushReq.RequestID = &reqID
		pushReq.Audience = &Audience{Alias: []string{userIDs[0]}}
		err = g.request(ctx, pushURL, pushReq, token, &pushResp)
	}
	switch err {
	case TokenExpireError:
		token, err = g.getTokenAndSave2Redis(ctx)
		if err != nil {
			log.NewError(operationID, utils.GetSelfFuncName(), "getTokenAndSave2Redis failed, ", err.Error())
		} else {
			log.NewInfo(operationID, utils.GetSelfFuncName(), "getTokenAndSave2Redis: ", token)
		}
	}
	if err != nil {
		return utils.Wrap(err, "push failed")
	}
	return utils.Wrap(err, "")
}

func (g *Client) Auth(ctx context.Context, timeStamp int64) (token string, expireTime int64, err error) {
	h := sha256.New()
	h.Write([]byte(config.Config.Push.Getui.AppKey + strconv.Itoa(int(timeStamp)) + config.Config.Push.Getui.MasterSecret))
	sign := hex.EncodeToString(h.Sum(nil))
	reqAuth := AuthReq{
		Sign:      sign,
		Timestamp: strconv.Itoa(int(timeStamp)),
		AppKey:    config.Config.Push.Getui.AppKey,
	}
	respAuth := AuthResp{}
	err = g.request(ctx, authURL, reqAuth, "", &respAuth)
	if err != nil {
		return "", 0, err
	}
	//log.NewInfo(operationID, utils.GetSelfFuncName(), "result: ", respAuth)
	expire, err := strconv.Atoi(respAuth.ExpireTime)
	return respAuth.Token, int64(expire), err
}

func (g *Client) GetTaskID(ctx context.Context, token string, pushReq PushReq) (string, error) {
	respTask := TaskResp{}
	ttl := int64(1000 * 60 * 5)
	pushReq.Settings = &Settings{TTL: &ttl}
	err := g.request(ctx, taskURL, pushReq, token, &respTask)
	if err != nil {
		return "", utils.Wrap(err, "")
	}
	return respTask.TaskID, nil
}

func (g *Client) request(ctx context.Context, url string, content interface{}, token string, output interface{}) error {
	con, err := json.Marshal(content)
	if err != nil {
		return err
	}
	client := &http.Client{}
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
	//log.NewDebug(operationID, "getui", utils.GetSelfFuncName(), "resp, ", string(result))
	commonResp := CommonResp{}
	commonResp.Data = output
	if err := json.Unmarshal(result, &commonResp); err != nil {
		return err
	}
	if commonResp.Code == tokenExpire {
		return TokenExpireError
	}
	return nil
}

func (g *Client) getTokenAndSave2Redis(ctx context.Context) (token string, err error) {
	token, _, err = g.Auth(ctx, time.Now().UnixNano()/1e6)
	if err != nil {
		return "", utils.Wrap(err, "Auth failed")
	}
	err = g.cache.SetGetuiTaskID(ctx, token, 60*60*23)
	if err != nil {
		return "", utils.Wrap(err, "Auth failed")
	}
	return token, nil
}

func (g *Client) GetTaskIDAndSave2Redis(ctx context.Context, token string, pushReq PushReq) (taskID string, err error) {
	ttl := int64(1000 * 60 * 60 * 24)
	pushReq.Settings = &Settings{TTL: &ttl}
	taskID, err = g.GetTaskID(ctx, token, pushReq)
	if err != nil {
		return "", utils.Wrap(err, "GetTaskIDAndSave2Redis failed")
	}
	err = g.cache.SetGetuiTaskID(ctx, taskID, 60*60*23)
	if err != nil {
		return "", utils.Wrap(err, "Auth failed")
	}
	return token, nil
}
