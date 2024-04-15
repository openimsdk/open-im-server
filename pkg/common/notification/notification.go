package notification

import (
	"github.com/openimsdk/open-im-server/v3/pkg/callbackstruct"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/util/memAsyncQueue"
	"github.com/openimsdk/tools/utils/httputil"
	"net/http"
)

package webhook

import (
"context"
"encoding/json"
"github.com/openimsdk/open-im-server/v3/pkg/callbackstruct"
"github.com/openimsdk/open-im-server/v3/pkg/common/config"
"github.com/openimsdk/open-im-server/v3/pkg/common/servererrs"
"github.com/openimsdk/open-im-server/v3/pkg/util/memAsyncQueue"
"github.com/openimsdk/protocol/constant"
"github.com/openimsdk/tools/log"
"github.com/openimsdk/tools/utils/httputil"
"net/http"
)

type Client struct {
	url    string
	queue  *memAsyncQueue.MemoryQueue
}

func NewWebhookClient(url string, queue *memAsyncQueue.MemoryQueue) *Client {
	http.DefaultTransport.(*http.Transport).MaxConnsPerHost = 100 // Enhance the default number of max connections per host
	return &Client{
		
		url:    url,
		queue:  queue,
	}
}

func (c *Client) SyncPost(ctx context.Context, command string, req callbackstruct.CallbackReq, resp callbackstruct.CallbackResp, before *config.BeforeConfig) error {
	if before.Enable {
		return c.post(ctx, command, req, resp, before.Timeout)
	}
	return nil
}

func (c *Client) AsyncPost(ctx context.Context, command string, req callbackstruct.CallbackReq, resp callbackstruct.CallbackResp, after *config.AfterConfig) {
	if after.Enable {
		c.queue.Push(func() { c.post(ctx, command, req, resp, after.Timeout) })
	}
}