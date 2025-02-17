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

package getui

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"strconv"
	"sync"
	"time"

	"github.com/openimsdk/open-im-server/v3/internal/push/offlinepush/options"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mcontext"
	"github.com/openimsdk/tools/utils/httputil"
	"github.com/openimsdk/tools/utils/splitter"
	"github.com/redis/go-redis/v9"
)

var (
	ErrTokenExpire = errs.New("token expire")
	ErrUserIDEmpty = errs.New("userIDs is empty")
)

const (
	pushURL      = "/push/single/alias"
	authURL      = "/auth"
	taskURL      = "/push/list/message"
	batchPushURL = "/push/list/alias"

	// Codes.
	tokenExpireCode = 10001
	tokenExpireTime = 60 * 60 * 23
	taskIDTTL       = 1000 * 60 * 60 * 24
)

type Client struct {
	cache           cache.ThirdCache
	tokenExpireTime int64
	taskIDTTL       int64
	pushConf        *config.Push
	httpClient      *httputil.HTTPClient
}

func NewClient(pushConf *config.Push, cache cache.ThirdCache) *Client {
	return &Client{cache: cache,
		tokenExpireTime: tokenExpireTime,
		taskIDTTL:       taskIDTTL,
		pushConf:        pushConf,
		httpClient:      httputil.NewHTTPClient(httputil.NewClientConfig()),
	}
}

func (g *Client) Push(ctx context.Context, userIDs []string, title, content string, opts *options.Opts) error {
	token, err := g.cache.GetGetuiToken(ctx)
	if err != nil {
		if errors.Is(err, redis.Nil) {
			log.ZDebug(ctx, "getui token not exist in redis")
			token, err = g.getTokenAndSave2Redis(ctx)
			if err != nil {
				return err
			}
		} else {
			return err
		}
	}
	pushReq := newPushReq(g.pushConf, title, content)
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
					for i := 0; i < len(userIDs); i += maxNum {
						end := i + maxNum
						if end > len(userIDs) {
							end = len(userIDs)
						}
						if err = g.batchPush(ctx, token, userIDs[i:end], pushReq); err != nil {
							log.ZError(ctx, "batchPush failed", err, "index", index, "token", token, "req", pushReq)
						}
					}
					if err = g.batchPush(ctx, token, userIDs, pushReq); err != nil {
						log.ZError(ctx, "batchPush failed", err, "index", index, "token", token, "req", pushReq)
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
		return ErrUserIDEmpty
	}
	switch err {
	case ErrTokenExpire:
		token, err = g.getTokenAndSave2Redis(ctx)
	}
	return err
}

func (g *Client) Auth(ctx context.Context, timeStamp int64) (token string, expireTime int64, err error) {
	h := sha256.New()
	h.Write(
		[]byte(g.pushConf.GeTui.AppKey + strconv.Itoa(int(timeStamp)) + g.pushConf.GeTui.MasterSecret),
	)
	sign := hex.EncodeToString(h.Sum(nil))
	reqAuth := AuthReq{
		Sign:      sign,
		Timestamp: strconv.Itoa(int(timeStamp)),
		AppKey:    g.pushConf.GeTui.AppKey,
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
	pushReq.Settings = &Settings{TTL: &ttl, Strategy: defaultStrategy}
	err := g.request(ctx, taskURL, pushReq, token, &respTask)
	if err != nil {
		return "", errs.Wrap(err)
	}
	return respTask.TaskID, nil
}

// max num is 999.
func (g *Client) batchPush(ctx context.Context, token string, userIDs []string, pushReq PushReq) error {
	taskID, err := g.GetTaskID(ctx, token, pushReq)
	if err != nil {
		return err
	}
	pushReq = newBatchPushReq(userIDs, taskID)
	return g.request(ctx, batchPushURL, pushReq, token, nil)
}

func (g *Client) singlePush(ctx context.Context, token, userID string, pushReq PushReq) error {
	operationID := mcontext.GetOperationID(ctx)
	pushReq.RequestID = &operationID
	pushReq.Audience = &Audience{Alias: []string{userID}}
	return g.request(ctx, pushURL, pushReq, token, nil)
}

func (g *Client) request(ctx context.Context, url string, input any, token string, output any) error {
	header := map[string]string{"token": token}
	resp := &Resp{}
	resp.Data = output
	return g.postReturn(ctx, g.pushConf.GeTui.PushUrl+url, header, input, resp, 3)
}

func (g *Client) postReturn(
	ctx context.Context,
	url string,
	header map[string]string,
	input any,
	output RespI,
	timeout int,
) error {
	log.ZDebug(ctx, "url:", url, "header:", header, "input:", input, "timeout:", timeout)
	err := g.httpClient.PostReturn(ctx, url, header, input, output, timeout)
	if err != nil {
		return err
	}
	log.ZDebug(ctx, "output:", output)
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
	pushReq.Settings = &Settings{TTL: &g.taskIDTTL, Strategy: defaultStrategy}
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
