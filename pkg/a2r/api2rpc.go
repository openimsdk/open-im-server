package a2r

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/apiresp"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func Call[A, B, C any](
	rpc func(client C, ctx context.Context, req *A, options ...grpc.CallOption) (*B, error),
	client func(ctx context.Context) (C, error),
	c *gin.Context,
) {
	log.ZDebug(c, "before bind")
	var req A
	if err := c.BindJSON(&req); err != nil {
		log.ZWarn(c, "gin bind json error", err, "req", req)
		apiresp.GinError(c, errs.ErrArgs.WithDetail(err.Error()).Wrap()) // 参数错误
		return
	}
	if check, ok := any(&req).(interface{ Check() error }); ok {
		if err := check.Check(); err != nil {
			log.ZWarn(c, "custom check error", err, "req", req)
			apiresp.GinError(c, errs.ErrArgs.Wrap(err.Error())) // 参数校验失败
			return
		}
	}
	log.ZDebug(c, "before get grpc conn")
	cli, err := client(c)
	if err != nil {
		log.ZError(c, "get conn error", err, "req", req)
		apiresp.GinError(c, errs.ErrInternalServer.Wrap(err.Error())) // 获取RPC连接失败
		return
	}
	log.ZDebug(c, "before call rpc")
	data, err := rpc(cli, c, &req)
	if err != nil {
		log.ZError(c, "rpc call error", err, "req", req)
		apiresp.GinError(c, err) // RPC调用失败
		return
	}
	apiresp.GinSuccess(c, data) // 成功
}
