package mw

import (
	"OpenIM/pkg/common/constant"
	"OpenIM/pkg/common/log"
	"OpenIM/pkg/common/tracelog"
	"OpenIM/pkg/errs"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"runtime/debug"
)

const OperationID = "operationID"
const OpUserID = "opUserID"

func rpcServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	var operationID string
	defer func() {
		if r := recover(); r != nil {
			log.NewError(operationID, info.FullMethod, "type:", fmt.Sprintf("%T", r), "panic:", r, "stack:", string(debug.Stack()))
		}
	}()
	log.Info("", "rpc come here,in rpc call")
	funcName := info.FullMethod
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.New(codes.InvalidArgument, "missing metadata").Err()
	}
	if opts := md.Get(OperationID); len(opts) != 1 || opts[0] == "" {
		return nil, status.New(codes.InvalidArgument, "operationID error").Err()
	} else {
		operationID = opts[0]
	}
	var opUserID string
	if opts := md.Get(OpUserID); len(opts) == 1 {
		opUserID = opts[0]
	}
	ctx = tracelog.SetFuncInfos(ctx, funcName, operationID)
	tracelog.SetCtxInfo(ctx, funcName, err, "opUserID", opUserID, "rpcReq", rpcString(req))
	resp, err = handler(ctx, req)
	if err != nil {
		tracelog.SetCtxInfo(ctx, funcName, err)
		return nil, rpcErrorToCode(err).Err()
	}
	tracelog.SetCtxInfo(ctx, funcName, nil, "rpcResp", rpcString(resp))
	return
}

func rpcClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
	if ctx == nil {
		return errs.ErrInternalServer.Wrap("call rpc request context is nil")
	}
	operationID, ok := ctx.Value(constant.OperationID).(string)
	if !ok {
		return errs.ErrArgs.Wrap("ctx missing operationID")
	}
	md := metadata.Pairs(constant.OperationID, operationID)
	opUserID, ok := ctx.Value(constant.OpUserID).(string)
	if ok {
		md.Append(constant.OpUserID, opUserID)
	}
	log.Info("", "rpc come here before")
	err = invoker(metadata.NewOutgoingContext(ctx, md), method, req, reply, cc, opts...)
	if err == nil {
		return nil
	}
	log.Info("", "rpc come here err", err.Error())

	rpcErr, ok := err.(interface{ GRPCStatus() *status.Status })
	if !ok {
		return errs.NewCodeError(errs.DefaultOtherError, err.Error()).Wrap()
	}
	sta := rpcErr.GRPCStatus()
	if sta.Code() == 0 {
		return errs.NewCodeError(errs.DefaultOtherError, err.Error()).Wrap()
	}
	details := sta.Details()
	if len(details) == 0 {
		return errs.NewCodeError(int(sta.Code()), sta.Message()).Wrap()
	}
	if v, ok := details[0].(*wrapperspb.StringValue); ok {
		return errs.NewCodeError(int(sta.Code()), sta.Message()).Wrap(v.String())
	}
	return errs.NewCodeError(int(sta.Code()), sta.Message()).Wrap()
}

func GrpcServer() grpc.ServerOption {
	return grpc.UnaryInterceptor(rpcServerInterceptor)
}

func GrpcClient() grpc.DialOption {
	return grpc.WithUnaryInterceptor(rpcClientInterceptor)
}
