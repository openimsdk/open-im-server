package mw

import (
	"OpenIM/pkg/common/constant"
	"OpenIM/pkg/common/log"
	"OpenIM/pkg/common/mw/specialerror"
	"OpenIM/pkg/common/tracelog"
	"OpenIM/pkg/errs"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"math"
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
	log.Info(OperationID, "opUserID", opUserID, "RPC", funcName, "Req", rpcString(req))
	ctx = tracelog.SetFuncInfos(ctx, funcName, operationID)
	ctx = context.WithValue(ctx, OperationID, operationID)
	ctx = context.WithValue(ctx, OpUserID, opUserID)
	resp, err = handler(ctx, req)
	if err != nil {
		log.Info(operationID, "rpc error:", err.Error())
		unwrap := errs.Unwrap(err)
		codeErr := specialerror.ErrCode(unwrap)
		if codeErr == nil {
			log.Error(operationID, "rpc InternalServer:", err.Error())
			codeErr = errs.ErrInternalServer
		}
		if unwrap != err {
			log.Info(operationID, "rpc error stack:", fmt.Sprintf("%+v", err))
		}
		code := codeErr.Code()
		if code <= 0 || code > math.MaxUint32 {
			log.Error(operationID, "rpc UnknownCode:", code, "err:", err.Error())
			code = errs.UnknownCode
		}
		grpcStatus := status.New(codes.Code(code), codeErr.Msg())
		if errs.Unwrap(err) != err {
			stack := fmt.Sprintf("%+v", err)
			log.Info(operationID, "rpc stack:", stack)
			if details, err := grpcStatus.WithDetails(wrapperspb.String(stack)); err == nil {
				grpcStatus = details
			}
		}
		return nil, grpcStatus.Err()
	}
	log.Info(OperationID, "opUserID", opUserID, "RPC", funcName, "Resp", rpcString(resp))
	return resp, nil
}

func rpcClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
	if ctx == nil {
		return errs.ErrInternalServer.Wrap("call rpc request context is nil")
	}
	operationID, ok := ctx.Value(constant.OperationID).(string)
	if !ok {
		log.Error("1111", "ctx missing operationID")
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
		return errs.ErrInternalServer.Wrap(err.Error())
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
