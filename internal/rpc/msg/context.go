package msg

import (
	pbChat "Open_IM/pkg/proto/chat"
	"context"
	"net/http"
)

// SendContext is the most important part of RPC SendMsg. It allows us to pass variables between middleware
type SendContext struct {
	ctx context.Context
	rpc *rpcChat
	// beforeFilters are filters which will be triggered before send msg
	beforeFilters []BeforeSendFilter
}

func NewSendContext(ctx context.Context, rpc *rpcChat) *SendContext {
	return &SendContext{
		ctx:           ctx,
		rpc:           rpc,
		beforeFilters: rpc.beforeSenders,
	}
}

func (c *SendContext) SetCtx(ctx context.Context) {
	c.ctx = ctx
}

func (c *SendContext) Value(key interface{}) interface{} {
	return c.ctx.Value(key)
}

func (c *SendContext) WithValue(key, val interface{}) {
	ctx := context.WithValue(c.ctx, key, val)
	c.SetCtx(ctx)
}

// doBeforeFilters executes the pending filters in the chain inside the calling handler.
func (c *SendContext) doBeforeFilters(pb *pbChat.SendMsgReq) (*pbChat.SendMsgResp, bool, error) {
	for _, handler := range c.beforeFilters {
		res, ok, err := handler(c, pb)
		if err != nil {
			return nil, false, err
		}
		if !ok {
			return res, ok, nil
		}
	}

	return nil, true, nil
}

func (c *SendContext) SendMsg(pb *pbChat.SendMsgReq) (*pbChat.SendMsgResp, error) {
	replay := pbChat.SendMsgResp{}
	res, ok, err := c.doBeforeFilters(pb)
	if err != nil {
		return returnMsg(&replay, pb, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError), err.Error(), 0)
	}
	if !ok {
		return res, nil
	}

	// fmt.Println("SEND_MSG:before send filters do over")

	res, err = c.rpc.doSendMsg(c.ctx, pb)
	if err != nil {
		return res, err
	}

	return nil, nil
}
