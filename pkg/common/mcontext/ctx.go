package mcontext

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
)

func NewCtx(operationID string) context.Context {
	c := context.Background()
	ctx := context.WithValue(c, constant.OperationID, operationID)
	SetOperationID(ctx, operationID)
	return ctx
}

func SetOperationID(ctx context.Context, operationID string) {
	ctx = context.WithValue(ctx, constant.OperationID, operationID)
}

func SetOpUserID(ctx context.Context, opUserID string) {
	ctx = context.WithValue(ctx, constant.OpUserID, opUserID)
}

func SetConnID(ctx context.Context, connID string) {
	ctx = context.WithValue(ctx, constant.ConnID, connID)
}

func GetOperationID(ctx context.Context) string {
	if ctx.Value(constant.OperationID) != nil {
		s, ok := ctx.Value(constant.OperationID).(string)
		if ok {
			return s
		}
	}
	return ""
}
func GetMustCtxInfo(ctx context.Context) (operationID, opUserID, platform, connID string, err error) {
	operationID, ok := ctx.Value(constant.OperationID).(string)
	if !ok {
		err = errs.ErrArgs.Wrap("ctx missing operationID")
		return
	}
	opUserID, ok1 := ctx.Value(constant.OpUserID).(string)
	if !ok1 {
		err = errs.ErrArgs.Wrap("ctx missing opUserID")
		return
	}
	platform, ok2 := ctx.Value(constant.OpUserPlatform).(string)
	if !ok2 {
		err = errs.ErrArgs.Wrap("ctx missing platform")
		return
	}
	connID, _ = ctx.Value(constant.ConnID).(string)
	return

}

func GetOpUserID(ctx context.Context) string {
	if ctx.Value(constant.OpUserID) != "" {
		s, ok := ctx.Value(constant.OpUserID).(string)
		if ok {
			return s
		}
	}
	return ""
}

func GetConnID(ctx context.Context) string {
	if ctx.Value(constant.ConnID) != "" {
		s, ok := ctx.Value(constant.ConnID).(string)
		if ok {
			return s
		}
	}
	return ""
}
