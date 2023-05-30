package mw

import (
	"context"
	"fmt"
	"math"
	"runtime"
	"strings"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/mw/specialerror"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/errinfo"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func rpcString(v interface{}) string {
	if s, ok := v.(interface{ String() string }); ok {
		return s.String()
	}
	return fmt.Sprintf("%+v", v)
}

func RpcServerInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	log.ZDebug(ctx, "rpc server req", "req", rpcString(req))

	//defer func() {
	//	if r := recover(); r != nil {
	//		log.ZError(ctx, "rpc panic", nil, "FullMethod", info.FullMethod, "type:", fmt.Sprintf("%T", r), "panic:", r)
	//		fmt.Printf("panic: %+v\nstack info: %s\n", r, string(debug.Stack()))
	//		pc, file, line, ok := runtime.Caller(4)
	//		if !ok {
	//			panic("get runtime.Caller failed")
	//		}
	//		errInfo := &errinfo.ErrorInfo{
	//			Path:  file,
	//			Line:  uint32(line),
	//			Name:  runtime.FuncForPC(pc).Name(),
	//			Cause: fmt.Sprintf("%s", r),
	//			Warp:  nil,
	//		}
	//		sta, err_ := status.New(codes.Code(errs.ErrInternalServer.Code()), errs.ErrInternalServer.Msg()).WithDetails(errInfo)
	//		if err_ != nil {
	//			panic(err_)
	//		}
	//		err = sta.Err()
	//	}
	//}()
	funcName := info.FullMethod
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.New(codes.InvalidArgument, "missing metadata").Err()
	}
	if keys := md.Get(constant.RpcCustomHeader); len(keys) > 0 {
		for _, key := range keys {
			values := md.Get(key)
			if len(values) == 0 {
				return nil, status.New(codes.InvalidArgument, fmt.Sprintf("missing metadata key %s", key)).Err()
			}
			ctx = context.WithValue(ctx, key, values)
		}
	}
	args := make([]string, 0, 4)
	if opts := md.Get(constant.OperationID); len(opts) != 1 || opts[0] == "" {
		return nil, status.New(codes.InvalidArgument, "operationID error").Err()
	} else {
		args = append(args, constant.OperationID, opts[0])
		ctx = context.WithValue(ctx, constant.OperationID, opts[0])
	}
	if opts := md.Get(constant.OpUserID); len(opts) == 1 {
		args = append(args, constant.OpUserID, opts[0])
		ctx = context.WithValue(ctx, constant.OpUserID, opts[0])
	}
	if opts := md.Get(constant.OpUserPlatform); len(opts) == 1 {
		ctx = context.WithValue(ctx, constant.OpUserPlatform, opts[0])
	}
	if opts := md.Get(constant.ConnID); len(opts) == 1 {
		ctx = context.WithValue(ctx, constant.ConnID, opts[0])
	}
	if opts := md.Get(constant.CheckKey); len(opts) != 1 || opts[0] == "" {
		return nil, status.New(codes.InvalidArgument, "check key empty").Err()
	} else {
		if err := verifyReqKey(args, opts[0]); err != nil {
			return nil, status.New(codes.InvalidArgument, err.Error()).Err()
		}
	}
	log.ZInfo(ctx, "rpc server req", "funcName", funcName, "req", rpcString(req))
	resp, err = handler(ctx, req)
	if err == nil {
		log.ZInfo(ctx, "rpc server resp", "funcName", funcName, "resp", rpcString(resp))
		return resp, nil
	}
	log.ZError(ctx, "rpc server resp", err, "funcName", funcName)
	unwrap := errs.Unwrap(err)
	codeErr := specialerror.ErrCode(unwrap)
	if codeErr == nil {
		log.ZError(ctx, "rpc InternalServer error", err, "req", req)
		codeErr = errs.ErrInternalServer
	}
	code := codeErr.Code()
	if code <= 0 || code > math.MaxUint32 {
		log.ZError(ctx, "rpc UnknownError", err, "rpc UnknownCode:", code)
		code = errs.ServerInternalError
	}
	grpcStatus := status.New(codes.Code(code), codeErr.Msg())
	var errInfo *errinfo.ErrorInfo
	if config.Config.Log.WithStack {
		if unwrap != err {
			sti, ok := err.(interface{ StackTrace() errors.StackTrace })
			if ok {
				log.ZWarn(ctx, "rpc server resp", err, "funcName", funcName, "unwrap", unwrap.Error(), "stack", fmt.Sprintf("%+v", err))
				if fs := sti.StackTrace(); len(fs) > 0 {
					pc := uintptr(fs[0])
					fn := runtime.FuncForPC(pc)
					file, line := fn.FileLine(pc)
					errInfo = &errinfo.ErrorInfo{
						Path:  file,
						Line:  uint32(line),
						Name:  fn.Name(),
						Cause: unwrap.Error(),
						Warp:  nil,
					}
					if arr := strings.Split(err.Error(), ": "); len(arr) > 1 {
						errInfo.Warp = arr[:len(arr)-1]
					}
				}
			}
		}
	}
	if errInfo == nil {
		errInfo = &errinfo.ErrorInfo{Cause: err.Error()}
	}
	details, err := grpcStatus.WithDetails(errInfo)
	if err != nil {
		panic(err)
	}
	log.ZWarn(ctx, "rpc server resp", err, "funcName", funcName)
	return nil, details.Err()
}

func GrpcServer() grpc.ServerOption {
	return grpc.UnaryInterceptor(RpcServerInterceptor)
}
