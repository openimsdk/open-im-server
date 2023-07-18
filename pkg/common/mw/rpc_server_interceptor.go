// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package mw

import (
	"context"
	"fmt"
	"math"
	"runtime"
	"strings"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/checker"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/mw/specialerror"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/errinfo"
)

func rpcString(v interface{}) string {
	if s, ok := v.(interface{ String() string }); ok {
		return s.String()
	}
	return fmt.Sprintf("%+v", v)
}

func RpcServerInterceptor(
	ctx context.Context,
	req interface{},
	info *grpc.UnaryServerInfo,
	handler grpc.UnaryHandler,
) (resp interface{}, err error) {
	log.ZDebug(ctx, "rpc server req", "req", rpcString(req))
	//defer func() {
	//	if r := recover(); r != nil {
	// 		log.ZError(ctx, "rpc panic", nil, "FullMethod", info.FullMethod, "type:", fmt.Sprintf("%T", r), "panic:", r)
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
	// 		sta, err_ := status.New(codes.Code(errs.ErrInternalServer.Code()),
	// errs.ErrInternalServer.Msg()).WithDetails(errInfo)
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
	log.ZInfo(ctx, "rpc server req", "funcName", funcName, "req", rpcString(req))
	resp, err = func() (interface{}, error) {
		if err := checker.Validate(req); err != nil {
			return nil, err
		}
		return handler(ctx, req)
	}()
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
				log.ZWarn(
					ctx,
					"rpc server resp",
					err,
					"funcName",
					funcName,
					"unwrap",
					unwrap.Error(),
					"stack",
					fmt.Sprintf("%+v", err),
				)
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
		log.ZWarn(ctx, "rpc server resp WithDetails error", err, "funcName", funcName)
		return nil, errs.Wrap(err)
	}
	log.ZWarn(ctx, "rpc server resp", err, "funcName", funcName)
	return nil, details.Err()
}

func GrpcServer() grpc.ServerOption {
	return grpc.ChainUnaryInterceptor(RpcServerInterceptor)
}
