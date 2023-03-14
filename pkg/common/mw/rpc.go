package mw

import (
	"OpenIM/pkg/common/constant"
	"OpenIM/pkg/common/log"
	"OpenIM/pkg/common/mw/specialerror"
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

	"errors"
)

const OperationID = "operationID"
const OpUserID = "opUserID"

func rpcString(v interface{}) string {
	if s, ok := v.(interface{ String() string }); ok {
		return s.String()
	}
	return fmt.Sprintf("%+v", v)
}

func rpcServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	var operationID string
	defer func() {
		if r := recover(); r != nil {
			log.ZError(ctx, "rpc panic", nil, "FullMethod", info.FullMethod, "type:", fmt.Sprintf("%T", r), "panic:", r, string(debug.Stack()))
		}
	}()
	log.Info("", "rpc come here,in rpc call")
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
	if err == nil {
		log.Info(operationID, "opUserID", opUserID, "RPC", funcName, "Resp", rpcString(resp))
		return resp, nil
	}
	log.ZError(ctx, "rpc InternalServer:", err, "req", req)
	unwrap := errs.Unwrap(err)
	codeErr := specialerror.ErrCode(unwrap)
	if codeErr == nil {
		log.ZError(ctx, "rpc InternalServer:", err, "req", req)
		codeErr = errs.ErrInternalServer
	}
	var stack string
	if unwrap != err {
		stack = fmt.Sprintf("%+v", err)
		log.ZError(ctx, "rpc error stack:", err)
	}
	code := codeErr.Code()
	if code <= 0 || code > math.MaxUint32 {
		log.ZError(ctx, "rpc UnknownError", err, "rpc UnknownCode:", code)
		code = errs.ServerInternalError
	}
	grpcStatus := status.New(codes.Code(code), codeErr.Msg())
	if errs.Unwrap(err) != err {
		if details, err := grpcStatus.WithDetails(wrapperspb.String(stack)); err == nil {
			grpcStatus = details
		}
	}
	log.ZInfo(ctx, "rpc resp", "funcName", funcName, "Resp", rpcString(resp))
	return nil, grpcStatus.Err()
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
	log.Info(operationID, "OpUserID", "RPC", method, "Req", rpcString(req))
	err = invoker(metadata.NewOutgoingContext(ctx, md), method, req, reply, cc, opts...)
	if err == nil {
		log.Info(operationID, "Resp", rpcString(reply))
		return nil
	}
	log.Info(operationID, "rpc error:", err.Error())
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

func GrpcServer() grpc.ServerOption {
	return grpc.UnaryInterceptor(rpcServerInterceptor)
}

func GrpcClient() grpc.DialOption {
	return grpc.WithUnaryInterceptor(rpcClientInterceptor)
}
