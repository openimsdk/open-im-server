package middleware

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/trace_log"
	"Open_IM/pkg/utils"
	"context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"path"
	"runtime/debug"
)

func RpcServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if r := recover(); r != nil {
			log.NewError("", info.FullMethod, "panic", r, "stack", string(debug.Stack()))
		}
	}()
	funcName := path.Base(info.FullMethod)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.New(codes.InvalidArgument, "missing metadata").Err()
	}
	var operationID string
	if opts := md.Get("operationID"); len(opts) != 1 || opts[0] == "" {
		return nil, status.New(codes.InvalidArgument, "operationID error").Err()
	} else {
		operationID = opts[0]
	}
	var opUserID string
	if opts := md.Get("opUserID"); len(opts) != 1 {
		return nil, status.New(codes.InvalidArgument, "opUserID error").Err()
	} else {
		opUserID = opts[0]
	}
	ctx = trace_log.NewRpcCtx(ctx, funcName, operationID)
	defer trace_log.ShowLog(ctx)
	trace_log.SetContextInfo(ctx, funcName, err, "opUserID", opUserID, "rpcReq", req.(interface{ String() string }).String())
	resp, err = handler(ctx, req)
	if err != nil {
		trace_log.SetContextInfo(ctx, funcName, err)
		errInfo := constant.ToAPIErrWithErr(err)
		var code codes.Code
		if errInfo.ErrCode == 0 {
			code = codes.Unknown
		} else {
			code = codes.Code(errInfo.ErrCode)
		}
		sta, err := status.New(code, errInfo.ErrMsg).WithDetails(wrapperspb.String(errInfo.DetailErrMsg))
		if err != nil {
			return nil, err
		}
		return nil, sta.Err()
	}
	trace_log.SetContextInfo(ctx, funcName, nil, "rpcResp", resp.(interface{ String() string }).String())
	return
}

func RpcClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
	if cc == nil {
		return utils.Wrap(constant.ErrRpcConn, "")
	}
	operationID, ok := ctx.Value("operationID").(string)
	if !ok {
		return utils.Wrap(constant.ErrArgs, "ctx missing operationID")
	}
	opUserID, ok := ctx.Value("opUserID").(string)
	if !ok {
		return utils.Wrap(constant.ErrArgs, "ctx missing opUserID")
	}
	md := metadata.Pairs("operationID", operationID, "opUserID", opUserID)
	return invoker(metadata.NewOutgoingContext(ctx, md), method, req, reply, cc, opts...)
}
