package apiresp

import (
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"reflect"
)

type apiResponse struct {
	ErrCode int    `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
	ErrDlt  string `json:"errDlt"`
	Data    any    `json:"data,omitempty"`
}

func isEmptyStruct(v any) bool {
	typeOf := reflect.TypeOf(v)
	if typeOf.Kind() == reflect.Ptr {
		typeOf = typeOf.Elem()
	}
	if typeOf.Kind() != reflect.Struct {
		return false
	}
	num := typeOf.NumField()
	for i := 0; i < num; i++ {
		v := typeOf.Field(i).Name[0]
		if v >= 'A' && v <= 'Z' {
			return false
		}
	}
	return true
}

func apiSuccess(data any) *apiResponse {
	if isEmptyStruct(data) {
		return &apiResponse{}
	}
	return &apiResponse{
		Data: data,
	}
}

func apiError(err error) *apiResponse {
	unwrap := errs.Unwrap(err)
	if codeErr, ok := unwrap.(errs.CodeError); ok {
		return &apiResponse{ErrCode: codeErr.Code(), ErrMsg: codeErr.Msg(), ErrDlt: codeErr.Detail()}
	}
	return &apiResponse{ErrCode: errs.ServerInternalError, ErrMsg: err.Error()}
}
