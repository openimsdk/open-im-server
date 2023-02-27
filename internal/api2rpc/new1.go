package api2rpc

import (
	"context"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

type FUNC[E, C, D any] func(client E, ctx context.Context, req *C, options ...grpc.CallOption) (*D, error)

// Call1 TEST
func Call1[A, B, C, D, E any](
	apiReq *A,
	apiResp *B,
	rpc FUNC[E, C, D],
	//client func() (E, error),
	c *gin.Context,
	before func(apiReq *A, rpcReq *C, bind func() error) error,
	after func(rpcResp *D, apiResp *B, bind func() error) error,
) {

}

func Call2[C, D, E any](
	rpc func(client E, ctx context.Context, req *C, options ...grpc.CallOption) (*D, error),
	client func() (E, error),
	c *gin.Context,
) {

}

func Call3[C, D, E any](
	rpc func(client E, ctx context.Context, req *C, options ...grpc.CallOption) (*D, error),
	//client func() (E, error),
	c *gin.Context,
) {

}

func Call4[C, D, E any](
	rpc func(client E, ctx context.Context, req *C, options ...grpc.CallOption) (*D, error),
	c *gin.Context,
) {

}

func Call10[A, B, C, D, E any](apiReq A, apiResp B, client func() (E, error), call func(client E, rpcReq C) (D, error)) {

}
