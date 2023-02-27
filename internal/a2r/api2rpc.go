package a2r

import (
	"context"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

//// Call TEST
//func Call2[A, B, C, D, E any](
//	apiReq *A,
//	apiResp *B,
//	rpc func(client E, ctx context.Context, req C, options ...grpc.CallOption) (D, error),
//	client func() (E, error),
//	c *gin.Context,
//	before func(apiReq *A, rpcReq *C, bind func() error) error,
//	after func(rpcResp *D, apiResp *B, bind func() error) error,
//) {
//
//}

func Call[A, B, C any](
	rpc func(client C, ctx context.Context, req *A, options ...grpc.CallOption) (*B, error),
	client func() (C, error),
	c *gin.Context,
) {
	var req A
	if err := c.BindJSON(&req); err != nil {
		// todo 参数错误
		return
	}
	if check, ok := any(&req).(interface{ Check() error }); ok {
		if err := check.Check(); err != nil {
			// todo 参数校验失败
			return
		}
	}
	cli, err := client()
	if err != nil {
		// todo 获取rpc连接失败
		return
	}
	resp, err := rpc(cli, c, &req)
	if err != nil {
		// todo rpc请求失败
		return
	}
	_ = resp
}
