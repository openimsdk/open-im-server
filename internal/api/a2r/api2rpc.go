package a2r

import (
	"OpenIM/internal/apiresp"
	"OpenIM/pkg/common/log"
	"OpenIM/pkg/errs"
	"context"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func Call[A, B, C any](
	rpc func(client C, ctx context.Context, req *A, options ...grpc.CallOption) (*B, error),
	client func() (C, error),
	c *gin.Context,
) {
	var req A
	if err := c.BindJSON(&req); err != nil {
		log.ZWarn(c, "gin bind json error", err, "req", req)
		apiresp.GinError(c, errs.ErrArgs.Wrap(err.Error())) // 参数错误
		return
	}
	if check, ok := any(&req).(interface{ Check() error }); ok {
		if err := check.Check(); err != nil {
			log.ZWarn(c, "custom check  error", err, "req", req)
			apiresp.GinError(c, errs.ErrArgs.Wrap(err.Error())) // 参数校验失败
			return
		}
	}
	cli, err := client()
	if err != nil {
		apiresp.GinError(c, errs.ErrInternalServer.Wrap(err.Error())) // 获取RPC连接失败
		return
	}
	data, err := rpc(cli, c, &req)
	if err != nil {
		apiresp.GinError(c, err) // RPC调用失败
		return
	}
	apiresp.GinSuccess(c, data) // 成功
}
