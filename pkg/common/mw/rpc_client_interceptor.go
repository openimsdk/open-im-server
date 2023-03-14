package mw

import (
	"OpenIM/pkg/common/constant"
	"OpenIM/pkg/common/log"
	"OpenIM/pkg/errs"
	"context"
	"errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
)

func GrpcClient() grpc.DialOption {
	return grpc.WithUnaryInterceptor(rpcClientInterceptor)
}

func rpcClientInterceptor(ctx context.Context, method string, req, resp interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
	if ctx == nil {
		return errs.ErrInternalServer.Wrap("call rpc request context is nil")
	}
	log.ZInfo(ctx, "rpc client req", "req", "funcName", method, rpcString(req))
	operationID, ok := ctx.Value(constant.OperationID).(string)
	if !ok {
		log.ZWarn(ctx, "ctx missing operationID", errors.New("ctx missing operationID"), "funcName", method)
		return errs.ErrArgs.Wrap("ctx missing operationID")
	}
	md := metadata.Pairs(constant.OperationID, operationID)
	opUserID, ok := ctx.Value(constant.OpUserID).(string)
	if ok {
		md.Append(constant.OpUserID, opUserID)
	}
	err = invoker(metadata.NewOutgoingContext(ctx, md), method, req, resp, cc, opts...)
	if err == nil {
		log.ZInfo(ctx, "rpc client resp", "funcName", method, rpcString(resp))
		return nil
	}
	log.ZError(ctx, "rpc result error:", err)
	rpcErr, ok := err.(interface{ GRPCStatus() *status.Status })
	if !ok {
		return errs.ErrInternalServer.Wrap(err.Error())
	}
	sta := rpcErr.GRPCStatus()
	if sta.Code() == 0 {
		return errs.NewCodeError(errs.ServerInternalError, err.Error()).Wrap()
	}
	if details := sta.Details(); len(details) > 0 {
		if v, ok := details[0].(*wrapperspb.StringValue); ok {
			return errs.NewCodeError(int(sta.Code()), sta.Message()).Wrap(v.String())
		}
	}
	return errs.NewCodeError(int(sta.Code()), sta.Message()).Wrap()
}
