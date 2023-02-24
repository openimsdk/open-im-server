package api2rpc

import (
	"context"
	"google.golang.org/grpc"
)

// NewRpc A: apiReq B: apiResp  C: rpcReq  D: rpcResp  Z: rpcClient (group.GroupClient)
func NewRpc[A, B any, C, D any, Z any](bind ApiBind[A, B], client func(conn *grpc.ClientConn) Z, rpc func(client Z, ctx context.Context, req *C, options ...grpc.CallOption) (*D, error)) *Rpc[A, B, C, D, Z] {
	return &Rpc[A, B, C, D, Z]{
		bind:   bind,
		client: client,
		rpc:    rpc,
	}
}

type Rpc[A, B any, C, D any, Z any] struct {
	bind   ApiBind[A, B]
	name   string
	client func(conn *grpc.ClientConn) Z
	rpc    func(client Z, ctx context.Context, req *C, options ...grpc.CallOption) (*D, error)
	before func(apiReq *A, rpcReq *C, bind func() error) error
	after  func(rpcResp *D, apiResp *B, bind func() error) error
}

func (a *Rpc[A, B, C, D, Z]) Name(name string) *Rpc[A, B, C, D, Z] {
	a.name = name
	return a
}

func (a *Rpc[A, B, C, D, Z]) Before(fn func(apiReq *A, rpcReq *C, bind func() error) error) *Rpc[A, B, C, D, Z] {
	a.before = fn
	return a
}

func (a *Rpc[A, B, C, D, Z]) After(fn func(rpcResp *D, apiResp *B, bind func() error) error) *Rpc[A, B, C, D, Z] {
	a.after = fn
	return a
}

func (a *Rpc[A, B, C, D, Z]) defaultCopyReq(apiReq *A, rpcReq *C) error {
	CopyAny(apiReq, rpcReq)
	return nil
}

func (a *Rpc[A, B, C, D, Z]) defaultCopyResp(rpcResp *D, apiResp *B) error {
	CopyAny(rpcResp, apiResp)
	return nil
}

func (a *Rpc[A, B, C, D, Z]) GetGrpcConn() (*grpc.ClientConn, error) {
	return nil, nil // todo
}

func (a *Rpc[A, B, C, D, Z]) execute() (*B, error) {
	var apiReq A
	if err := a.bind.Bind(&apiReq); err != nil {
		return nil, err
	}
	opID := a.bind.OperationID()
	userID, err := a.bind.OpUserID()
	if err != nil {
		return nil, err
	}
	_, _ = opID, userID
	var rpcReq C
	if a.before == nil {
		err = a.defaultCopyReq(&apiReq, &rpcReq)
	} else {
		err = a.before(&apiReq, &rpcReq, func() error { return a.defaultCopyReq(&apiReq, &rpcReq) })
	}
	if err != nil {
		return nil, err
	}
	conn, err := a.GetGrpcConn()
	if err != nil {
		return nil, err
	}
	rpcResp, err := a.rpc(a.client(conn), a.bind.Context(), &rpcReq)
	if err != nil {
		return nil, err
	}
	var apiResp B
	if a.after == nil {
		err = a.defaultCopyResp(rpcResp, &apiResp)
	} else {
		err = a.after(rpcResp, &apiResp, func() error { return a.defaultCopyResp(rpcResp, &apiResp) })
	}
	if err != nil {
		return nil, err
	}
	return &apiResp, nil
}

func (a *Rpc[A, B, C, D, Z]) Call() {
	a.bind.Resp(a.execute())
}
