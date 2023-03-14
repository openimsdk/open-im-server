package mw

import (
	"OpenIM/pkg/common/constant"
	"OpenIM/pkg/common/log"
	"OpenIM/pkg/common/mw/specialerror"
	"OpenIM/pkg/errs"
	"context"
	"errors"
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
			log.ZError(ctx, "rpc panic", nil, "FullMethod", info.FullMethod, "type:", fmt.Sprintf("%T", r), "panic:", r, string(debug.Stack()))
		}
	}()
	funcName := info.FullMethod
	log.ZInfo(ctx, "rpc req", "funcName", funcName, "req", rpcString(req))
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
	ctx = context.WithValue(ctx, OperationID, operationID)
	ctx = context.WithValue(ctx, OpUserID, opUserID)
	resp, err = handler(ctx, req)
	if err != nil {
		log.ZError(ctx, "handler rpc error", err, "req", req)
		unwrap := errs.Unwrap(err)
		codeErr := specialerror.ErrCode(unwrap)
		if codeErr == nil {
			log.ZError(ctx, "rpc InternalServer:", err, "req", req)
			codeErr = errs.ErrInternalServer
		}
		if unwrap != err {
			log.ZError(ctx, "rpc error stack:", err)
		}
		code := codeErr.Code()
		if code <= 0 || code > math.MaxUint32 {
			log.ZError(ctx, "rpc UnknownError", err, "rpc UnknownCode:", code)
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
	log.ZInfo(ctx, "rpc resp", "funcName", funcName, "Resp", rpcString(resp))
	return resp, nil
}

func rpcClientInterceptor(ctx context.Context, method string, req, reply interface{}, cc *grpc.ClientConn, invoker grpc.UnaryInvoker, opts ...grpc.CallOption) (err error) {
	if ctx == nil {
		return errs.ErrInternalServer.Wrap("call rpc request context is nil")
	}
	operationID, ok := ctx.Value(constant.OperationID).(string)
	if !ok {
		log.ZError(ctx, "ctx missing operationID", errors.New("ctx missing operationID"))
		return errs.ErrArgs.Wrap("ctx missing operationID")
	}
	md := metadata.Pairs(constant.OperationID, operationID)
	opUserID, ok := ctx.Value(constant.OpUserID).(string)
	if ok {
		md.Append(constant.OpUserID, opUserID)
	}
	err = invoker(metadata.NewOutgoingContext(ctx, md), method, req, reply, cc, opts...)
	if err == nil {
		return nil
	}
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
