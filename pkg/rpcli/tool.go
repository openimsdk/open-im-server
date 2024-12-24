package rpcli

import (
	"context"
	"github.com/openimsdk/tools/errs"
	"google.golang.org/grpc"
)

func extractField[A, B, C any](ctx context.Context, fn func(ctx context.Context, req *A, opts ...grpc.CallOption) (*B, error), req *A, get func(*B) C) (C, error) {
	resp, err := fn(ctx, req)
	if err != nil {
		var c C
		return c, err
	}
	return get(resp), nil
}

func firstValue[A any](val []A, err error) (A, error) {
	if err != nil {
		var a A
		return a, err
	}
	if len(val) == 0 {
		var a A
		return a, errs.ErrRecordNotFound.WrapMsg("record not found")
	}
	return val[0], nil
}

func ignoreResp(_ any, err error) error {
	return err
}
