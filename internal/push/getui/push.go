package getui

import (
	"OpenIM/internal/push"
	"OpenIM/pkg/common/config"
	"OpenIM/pkg/common/db/cache"
	http2 "OpenIM/pkg/common/http"
	"OpenIM/pkg/common/log"
	"OpenIM/pkg/common/tracelog"
	"OpenIM/pkg/utils/splitter"
	"github.com/go-redis/redis/v8"
	"sync"

	"OpenIM/pkg/utils"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strconv"
	"time"
)

var (
	TokenExpireError = errors.New("token expire")
	UserIDEmptyError = errors.New("userIDs is empty")
)

const (
	pushURL      = "/push/single/alias"
	authURL      = "/auth"
	taskURL      = "/push/list/message"
	batchPushURL = "/push/list/alias"

	// codes
	tokenExpireCode = 10001
	tokenExpireTime = 60 * 60 * 23
	taskIDTTL       = 1000 * 60 * 60 * 24
)

type Client struct {
	cache           cache.Cache
	tokenExpireTime int64
	taskIDTTL       int64
}

func NewClient(cache cache.Cache) *Client {
	return &Client{cache: cache, tokenExpireTime: tokenExpireTime, taskIDTTL: taskIDTTL}
}

func (g *Client) Push(ctx context.Context, userIDs []string, title, content string, opts *push.Opts) error {
	token, err := g.cache.GetGetuiToken(ctx)
	if err != nil {
		if err == redis.Nil {
			token, err = g.getTokenAndSave2Redis(ctx)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	pushReq := newPushReq(title, content)
	pushReq.setPushChannel(title, content)
	if len(userIDs) > 1 {
		maxNum := 999
		if len(userIDs) > maxNum {
			s := splitter.NewSplitter(maxNum, userIDs)
			wg := sync.WaitGroup{}
			wg.Add(len(s.GetSplitResult()))
			for i, v := range s.GetSplitResult() {
				go func(index int, userIDs []string) {
					defer wg.Done()
					if err = g.batchPush(ctx, token, userIDs, pushReq); err != nil {
						log.NewError(tracelog.GetOperationID(ctx), "batchPush failed", i, token, pushReq)
					}
				}(i, v.Item)
			}
			wg.Wait()
		} else {
			err = g.batchPush(ctx, token, userIDs, pushReq)
		}
	} else if len(userIDs) == 1 {
		err = g.singlePush(ctx, token, userIDs[0], pushReq)
	} else {
		return UserIDEmptyError
	}
	switch err {
	case TokenExpireError:
		token, err = g.getTokenAndSave2Redis(ctx)
	}
	return err
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

// max num is 999
func (g *Client) batchPush(ctx context.Context, token string, userIDs []string, pushReq PushReq) error {
	taskID, err := g.GetTaskID(ctx, token, pushReq)
	if err != nil {
		return err
	}
	pushReq = newBatchPushReq(userIDs, taskID)
	return g.request(ctx, batchPushURL, pushReq, token, nil)
}

func (g *Client) singlePush(ctx context.Context, token, userID string, pushReq PushReq) error {
	operationID := tracelog.GetOperationID(ctx)
	pushReq.RequestID = &operationID
	pushReq.Audience = &Audience{Alias: []string{userID}}
	return g.request(ctx, pushURL, pushReq, token, nil)
}

func (g *Client) request(ctx context.Context, url string, input interface{}, token string, output interface{}) error {
	header := map[string]string{"token": token}
	resp := &Resp{}
	resp.Data = output
	return g.postReturn(config.Config.Push.Getui.PushUrl+url, header, input, resp, 3)
}

func (g *Client) postReturn(url string, header map[string]string, input interface{}, output RespI, timeout int) error {
	err := http2.PostReturn(url, header, input, output, timeout)
	if err != nil {
		return err
	}
	return output.parseError()
}

func (g *Client) getTokenAndSave2Redis(ctx context.Context) (token string, err error) {
	token, _, err = g.Auth(ctx, time.Now().UnixNano()/1e6)
	if err != nil {
		return
	}
	err = g.cache.SetGetuiToken(ctx, token, 60*60*23)
	if err != nil {
		return
	}
	return token, nil
}

func (g *Client) GetTaskIDAndSave2Redis(ctx context.Context, token string, pushReq PushReq) (taskID string, err error) {
	pushReq.Settings = &Settings{TTL: &g.taskIDTTL}
	taskID, err = g.GetTaskID(ctx, token, pushReq)
	if err != nil {
		return
	}
	err = g.cache.SetGetuiTaskID(ctx, taskID, g.tokenExpireTime)
	if err != nil {
		return
	}
	return token, nil
}
