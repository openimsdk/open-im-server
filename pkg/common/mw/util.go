package mw

import (
	"OpenIM/pkg/errs"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/wrapperspb"
	"math"
)

func rpcString(v interface{}) string {
	if s, ok := v.(interface{ String() string }); ok {
		return s.String()
	}
	return fmt.Sprintf("%+v", v)
}

func rpcErrorToCode(err error) *status.Status {
	unwrap := errs.Unwrap(err)
	var (
		code codes.Code
		msg  string
	)
	if unwrap.(errs.CodeError) != nil {
		c := unwrap.(errs.CodeError).Code()
		if c <= 0 || c > math.MaxUint32 {
			code = codes.OutOfRange // 错误码超出范围
		} else {
			code = codes.Code(c)
		}
		msg = unwrap.(errs.CodeError).Msg()
	} else {
		code = codes.Unknown
		msg = unwrap.Error()
	}
	sta := status.New(code, msg)
	if unwrap == err {
		return sta
	}
	details, err := sta.WithDetails(wrapperspb.String(fmt.Sprintf("%+v", err)))
	if err != nil {
		return sta
	}
	return details
}
