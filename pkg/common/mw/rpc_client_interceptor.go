package mw

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/errinfo"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func GrpcClient() grpc.DialOption {
	return grpc.WithUnaryInterceptor(RpcClientInterceptor)
}

func RpcClientInterceptor(ctx context.Context, method string, req, resp interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
	if ctx == nil {
		return errs.ErrInternalServer.Wrap("call rpc request context is nil")
	}
	ctx, err = getRpcContext(ctx, method)
	if err != nil {
		return err
	}
	log.ZDebug(ctx, "get rpc ctx success", "conn target", cc.Target())
	err = invoker(ctx, method, req, resp, cc, opts...)
	if err == nil {
		log.ZInfo(ctx, "rpc client resp", "funcName", method, "resp", rpcString(resp))
		return nil
	}
	log.ZError(ctx, "rpc resp error", err)
	rpcErr, ok := err.(interface{ GRPCStatus() *status.Status })
	if !ok {
		return errs.ErrInternalServer.Wrap(err.Error())
	}
	sta := rpcErr.GRPCStatus()
	if sta.Code() == 0 {
		return errs.NewCodeError(errs.ServerInternalError, err.Error()).Wrap()
	}
	if details := sta.Details(); len(details) > 0 {
		errInfo, ok := details[0].(*errinfo.ErrorInfo)
		if ok {
			s := strings.Join(errInfo.Warp, "->") + errInfo.Cause
			return errs.NewCodeError(int(sta.Code()), sta.Message()).WithDetail(s).Wrap()
		}
	}
	return errs.NewCodeError(int(sta.Code()), sta.Message()).Wrap()
}

func getRpcContext(ctx context.Context, method string) (context.Context, error) {
	md := metadata.Pairs()
	if keys, _ := ctx.Value(constant.RpcCustomHeader).([]string); len(keys) > 0 {
		for _, key := range keys {
			val, ok := ctx.Value(key).([]string)
			if !ok {
				return nil, errs.ErrInternalServer.Wrap(fmt.Sprintf("ctx missing key %s", key))
			}
			if len(val) == 0 {
				return nil, errs.ErrInternalServer.Wrap(fmt.Sprintf("ctx key %s value is empty", key))
			}
			md.Set(key, val...)
		}
		md.Set(constant.RpcCustomHeader, keys...)
	}
	operationID, ok := ctx.Value(constant.OperationID).(string)
	if !ok {
		log.ZWarn(ctx, "ctx missing operationID", errors.New("ctx missing operationID"), "funcName", method)
		return nil, errs.ErrArgs.Wrap("ctx missing operationID")
	}
	md.Set(constant.OperationID, operationID)
	var checkArgs []string
	checkArgs = append(checkArgs, constant.OperationID, operationID)
	opUserID, ok := ctx.Value(constant.OpUserID).(string)
	if ok {
		md.Set(constant.OpUserID, opUserID)
		checkArgs = append(checkArgs, constant.OpUserID, opUserID)
	}
	opUserIDPlatformID, ok := ctx.Value(constant.OpUserPlatform).(string)
	if ok {
		md.Set(constant.OpUserPlatform, opUserIDPlatformID)
	}
	connID, ok := ctx.Value(constant.ConnID).(string)
	if ok {
		md.Set(constant.ConnID, connID)
	}
	md.Set(constant.CheckKey, genReqKey(checkArgs))
	return metadata.NewOutgoingContext(ctx, md), nil
}
