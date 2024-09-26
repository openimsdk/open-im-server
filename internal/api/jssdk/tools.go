package jssdk

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/openimsdk/tools/apiresp"
	"google.golang.org/grpc"
)

func field[A, B, C any](ctx context.Context, fn func(ctx context.Context, req *A, opts ...grpc.CallOption) (*B, error), req *A, get func(*B) C) (C, error) {
	resp, err := fn(ctx, req)
	if err != nil {
		var c C
		return c, err
	}
	return get(resp), nil
}

func call[R any](c *gin.Context, fn func(ctx *gin.Context) (R, error)) {
	resp, err := fn(c)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}
	apiresp.GinSuccess(c, resp)
}
