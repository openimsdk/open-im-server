package mcontext

import (
	"context"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
)

var mapper = []string{constant.OperationID, constant.OpUserID, constant.OpUserPlatform, constant.ConnID}

func WithOpUserIDContext(ctx context.Context, opUserID string) context.Context {
	return context.WithValue(ctx, constant.OpUserID, opUserID)
}
func WithOpUserPlatformContext(ctx context.Context, platform string) context.Context {
	return context.WithValue(ctx, constant.OpUserPlatform, platform)
}
func WithTriggerIDContext(ctx context.Context, triggerID string) context.Context {
	return context.WithValue(ctx, constant.TriggerID, triggerID)
}
func NewCtx(operationID string) context.Context {
	c := context.Background()
	ctx := context.WithValue(c, constant.OperationID, operationID)
	return SetOperationID(ctx, operationID)
}

func SetOperationID(ctx context.Context, operationID string) context.Context {
	return context.WithValue(ctx, constant.OperationID, operationID)
}

func SetOpUserID(ctx context.Context, opUserID string) context.Context {
	return context.WithValue(ctx, constant.OpUserID, opUserID)
}

func SetConnID(ctx context.Context, connID string) context.Context {
	return context.WithValue(ctx, constant.ConnID, connID)
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

func GetTriggerID(ctx context.Context) string {
	if ctx.Value(constant.TriggerID) != "" {
		s, ok := ctx.Value(constant.TriggerID).(string)
		if ok {
			return s
		}
	}
	return ""
}
func GetOpUserPlatform(ctx context.Context) string {
	if ctx.Value(constant.OpUserPlatform) != "" {
		s, ok := ctx.Value(constant.OpUserPlatform).(string)
		if ok {
			return s
		}
	}
	return ""
}
func GetRemoteAddr(ctx context.Context) string {
	if ctx.Value(constant.RemoteAddr) != "" {
		s, ok := ctx.Value(constant.RemoteAddr).(string)
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

func GetCtxInfos(ctx context.Context) (operationID, opUserID, platform, connID string, err error) {
	operationID, ok := ctx.Value(constant.OperationID).(string)
	if !ok {
		err = errs.ErrArgs.Wrap("ctx missing operationID")
		return
	}
	opUserID, _ = ctx.Value(constant.OpUserID).(string)
	platform, _ = ctx.Value(constant.OpUserPlatform).(string)
	connID, _ = ctx.Value(constant.ConnID).(string)
	return
}

func WithMustInfoCtx(values []string) context.Context {
	ctx := context.Background()
	for i, v := range values {
		ctx = context.WithValue(ctx, mapper[i], v)

	}
	return ctx

}
