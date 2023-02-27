package a2r

import (
	"context"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

// Call TEST
func Call[A, B, C, D, E any](
	apiReq *A,
	apiResp *B,
	rpc func(client E, ctx context.Context, req C, options ...grpc.CallOption) (D, error),
	client func() (E, error),
	c *gin.Context,
	before func(apiReq *A, rpcReq *C, bind func() error) error,
	after func(rpcResp *D, apiResp *B, bind func() error) error,
) {

}

func Call1[C, D, E any](
	rpc func(client E, ctx context.Context, req C, options ...grpc.CallOption) (D, error),
	client func() (E, error),
	c *gin.Context,
) {

}
