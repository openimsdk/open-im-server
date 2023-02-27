package api2rpc

import (
	"OpenIM/pkg/errs"
	"context"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type rpcFunc[E, C, D any] func(client E, ctx context.Context, req *C, options ...grpc.CallOption) (*D, error)

func New[A, B, C, D any, E any](apiReq *A, apiResp *B, rpc func(client E, ctx context.Context, req *C, options ...grpc.CallOption) (*D, error)) RpcCall[A, B, C, D, E] {
	return &rpcCall[A, B, C, D, E]{
		apiReq:  apiReq,
		apiResp: apiResp,
		rpcFn:   rpc,
	}
}

type rpcCall[A, B, C, D any, E any] struct {
	apiReq  *A
	apiResp *B
	rpcFn   func(client E, ctx context.Context, req *C, options ...grpc.CallOption) (*D, error)
	before  func(apiReq *A, rpcReq *C, bind func() error) error
	after   func(rpcResp *D, apiResp *B, bind func() error) error
}

func (r *rpcCall[A, B, C, D, E]) Before(fn func(apiReq *A, rpcReq *C, bind func() error) error) RpcCall[A, B, C, D, E] {
	r.before = fn
	return r
}

func (r *rpcCall[A, B, C, D, E]) After(fn func(rpcResp *D, apiResp *B, bind func() error) error) RpcCall[A, B, C, D, E] {
	r.after = fn
	return r
}

func (r *rpcCall[A, B, C, D, E]) Call(c *gin.Context, client func() (E, error)) {
	var resp baseResp
	err := r.call(c, client)
	if err == nil {
		resp.Data = r.apiResp
	} else {
		cerr, ok := err.(errs.Coderr)
		if ok {
			resp.ErrCode = int32(cerr.Code())
			resp.ErrMsg = cerr.Msg()
			resp.ErrDtl = cerr.Detail()
		} else {
			resp.ErrCode = 10000
			resp.ErrMsg = err.Error()
		}
	}
}

func (r *rpcCall[A, B, C, D, E]) defaultCopyReq(rpcReq *C) error {
	if r.apiReq != nil {
		CopyAny(r.apiReq, rpcReq)
	}
	return nil
}

func (r *rpcCall[A, B, C, D, E]) defaultCopyResp(rpcResp *D) error {
	if r.apiResp != nil {
		CopyAny(rpcResp, r.apiResp)
	}
	return nil
}

func (r *rpcCall[A, B, C, D, E]) call(c *gin.Context, client func() (E, error)) error {
	if err := c.BindJSON(r.apiReq); err != nil {
		return err
	}
	var err error
	var rpcReq C
	if r.before == nil {
		err = r.defaultCopyReq(&rpcReq)
	} else {
		err = r.before(r.apiReq, &rpcReq, func() error { return r.defaultCopyReq(&rpcReq) })
	}
	if err != nil {
		return err
	}
	cli, err := client()
	if err != nil {
		return err
	}
	rpcResp, err := r.rpcFn(cli, c, &rpcReq)
	if err != nil {
		return err
	}
	var apiResp B
	if r.after == nil {
		return r.defaultCopyResp(rpcResp)
	} else {
		return r.after(rpcResp, &apiResp, func() error { return r.defaultCopyResp(rpcResp) })
	}
}

type RpcCall[A, B, C, D, E any] interface {
	Before(fn func(apiReq *A, rpcReq *C, bind func() error) error) RpcCall[A, B, C, D, E]
	After(fn func(rpcResp *D, apiResp *B, bind func() error) error) RpcCall[A, B, C, D, E]
	Call(c *gin.Context, client func() (E, error))
}
