package api2rpc

import (
	"context"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type rpcFunc[E, C, D any] func(client E, ctx context.Context, req *C, options ...grpc.CallOption) (*D, error)

func Rpc[A, B, C, D any, E any](apiReq *A, apiResp *B, rpc func(client E, ctx context.Context, req *C, options ...grpc.CallOption) (*D, error)) RpcCall[A, B, C, D, E] {
	return &rpcCall[A, B, C, D, E]{
		apiReq:  apiReq,
		apiResp: apiResp,
		rpcFn:   rpc,
	}
}

type rpcCall[A, B, C, D any, E any] struct {
	apiReq  *A
	apiResp *B
	client  func() (E, error)
	rpcFn   func(client E, ctx context.Context, req *C, options ...grpc.CallOption) (*D, error)
	api     Api
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

func (r *rpcCall[A, B, C, D, E]) Must(c *gin.Context, client func() (E, error)) RpcCall[A, B, C, D, E] {
	r.api = NewGin1(c)
	r.client = client
	return r
}

func (r *rpcCall[A, B, C, D, E]) Call() {
	r.api.Resp(r.apiResp, r.call())
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

func (r *rpcCall[A, B, C, D, E]) call() error {
	if err := r.api.Bind(r.apiReq); err != nil {
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
	client, err := r.client()
	if err != nil {
		return err
	}
	rpcResp, err := r.rpcFn(client, r.api.Context(), &rpcReq)
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
	Must(c *gin.Context, client func() (E, error)) RpcCall[A, B, C, D, E]
	Call()
}
