package a2r

import (
	"OpenIM/internal/apiresp"
	"OpenIM/pkg/common/constant"
	"context"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"net/http"
)

func Call[A, B, C any](
	rpc func(client C, ctx context.Context, req *A, options ...grpc.CallOption) (*B, error),
	client func() (C, error),
	c *gin.Context,
) {
	var resp *apiresp.ApiResponse
	defer func() {
		c.JSON(http.StatusOK, resp)
	}()
	var req A
	if err := c.BindJSON(&req); err != nil {
		resp = apiresp.Error(constant.ErrArgs.Wrap(err.Error())) // 参数错误
		return
	}
	if check, ok := any(&req).(interface{ Check() error }); ok {
		if err := check.Check(); err != nil {
			resp = apiresp.Error(constant.ErrArgs.Wrap(err.Error())) // 参数校验失败
			return
		}
	}
	cli, err := client()
	if err != nil {
		resp = apiresp.Error(constant.ErrInternalServer.Wrap(err.Error())) // 获取RPC连接失败
		return
	}
	data, err := rpc(cli, c, &req)
	if err != nil {
		resp = apiresp.Error(err) // RPC调用失败
		return
	}
	resp = apiresp.Success(data) // 成功
}

func CallAny[A, B, C, D any](
	rpc func(client C, ctx context.Context, req *A, options ...grpc.CallOption) (*B, error),
	client func() (C, error),
	c *gin.Context,
	apiReq *D,
	cp func(apiReq *D, rpcReq *A) error,
) {
	var resp *apiresp.ApiResponse
	defer func() {
		c.JSON(http.StatusOK, resp)
	}()
	if err := c.BindJSON(apiReq); err != nil {
		resp = apiresp.Error(constant.ErrArgs.Wrap(err.Error())) // 参数错误
		return
	}
	var req A
	if err := cp(apiReq, &req); err != nil {
		resp = apiresp.Error(constant.ErrArgs.Wrap(err.Error())) // 参数错误
		return
	}
	if check, ok := any(&req).(interface{ Check() error }); ok {
		if err := check.Check(); err != nil {
			resp = apiresp.Error(constant.ErrArgs.Wrap(err.Error())) // 参数校验失败
			return
		}
	}
	cli, err := client()
	if err != nil {
		resp = apiresp.Error(constant.ErrInternalServer.Wrap(err.Error())) // 获取RPC连接失败
		return
	}
	data, err := rpc(cli, c, &req)
	if err != nil {
		resp = apiresp.Error(err) // RPC调用失败
		return
	}
	resp = apiresp.Success(data) // 成功
}
