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

package webhook

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/tools/log"
	"github.com/openimsdk/tools/mcontext"
	"github.com/openimsdk/tools/mq/memamq"
	"github.com/openimsdk/tools/utils/httputil"

	"github.com/openimsdk/open-im-server/v3/pkg/callbackstruct"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/servererrs"
)

type Client struct {
	client *httputil.HTTPClient
	url    string
	queue  *memamq.MemoryQueue
}

const (
	webhookWorkerCount = 2
	webhookBufferSize  = 100
)

func NewWebhookClient(url string, options ...*memamq.MemoryQueue) *Client {
	var queue *memamq.MemoryQueue
	if len(options) > 0 && options[0] != nil {
		queue = options[0]
	} else {
		queue = memamq.NewMemoryQueue(webhookWorkerCount, webhookBufferSize)
	}

	http.DefaultTransport.(*http.Transport).MaxConnsPerHost = 100 // Enhance the default number of max connections per host

	return &Client{
		client: httputil.NewHTTPClient(httputil.NewClientConfig()),
		url:    url,
		queue:  queue,
	}
}

func (c *Client) SyncPost(ctx context.Context, command string, req callbackstruct.CallbackReq, resp callbackstruct.CallbackResp, before *config.BeforeConfig) error {
	return c.post(ctx, command, req, resp, before.Timeout)
}

func (c *Client) AsyncPost(ctx context.Context, command string, req callbackstruct.CallbackReq, resp callbackstruct.CallbackResp, after *config.AfterConfig) {
	if after.Enable {
		c.queue.Push(func() { c.post(ctx, command, req, resp, after.Timeout) })
	}
}

func (c *Client) post(ctx context.Context, command string, input interface{}, output callbackstruct.CallbackResp, timeout int) error {
	ctx = mcontext.WithMustInfoCtx([]string{mcontext.GetOperationID(ctx), mcontext.GetOpUserID(ctx), mcontext.GetOpUserPlatform(ctx), mcontext.GetConnID(ctx)})
	fullURL := c.url + "/" + command
	log.ZInfo(ctx, "webhook", "url", fullURL, "input", input, "config", timeout)
	operationID, _ := ctx.Value(constant.OperationID).(string)
	b, err := c.client.Post(ctx, fullURL, map[string]string{constant.OperationID: operationID}, input, timeout)
	if err != nil {
		return servererrs.ErrNetwork.WrapMsg(err.Error(), "post url", fullURL)
	}
	if err = json.Unmarshal(b, output); err != nil {
		return servererrs.ErrData.WithDetail(err.Error() + " response format error")
	}
	if err := output.Parse(); err != nil {
		return err
	}
	log.ZInfo(ctx, "webhook success", "url", fullURL, "input", input, "response", string(b))
	return nil
}
