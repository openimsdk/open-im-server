package apiresp

import (
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"reflect"
)

type ApiCodeResponse struct {
	ErrCode int    `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
	ErrDlt  string `json:"errDlt"`
	Data    any    `json:"data,omitempty"`
}

func isAllFieldsPrivate(v any) bool {
	typeOf := reflect.TypeOf(v)
	if typeOf.Kind() == reflect.Ptr {
		typeOf = typeOf.Elem()
	}
	if typeOf.Kind() != reflect.Struct {
		return false
	}
	num := typeOf.NumField()
	for i := 0; i < num; i++ {
		c := typeOf.Field(i).Name[0]
		if c >= 'A' && c <= 'Z' {
			return false
		}
	}
	return true
}

func ApiSuccess(data any) *ApiCodeResponse {
	if isAllFieldsPrivate(data) {
		return &ApiCodeResponse{}
	}
	return &ApiCodeResponse{
		Data: data,
	}
}

func ApiError(err error) *ApiCodeResponse {
	if err == nil {
		return ApiSuccess(nil)
	}
	unwrap := errs.Unwrap(err)
	if codeErr, ok := unwrap.(errs.CodeError); ok {
		resp := ApiCodeResponse{ErrCode: codeErr.Code(), ErrMsg: codeErr.Msg(), ErrDlt: codeErr.Detail()}
		if resp.ErrDlt == "" {
			resp.ErrDlt = err.Error()
		}
		return &resp
	}
	return &ApiCodeResponse{ErrCode: errs.ServerInternalError, ErrMsg: err.Error()}
}
