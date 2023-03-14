package apiresp

import (
	"OpenIM/pkg/errs"
	"fmt"
)

type apiResponse struct {
	ErrCode int    `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
	ErrDlt  string `json:"errDlt"`
	Data    any    `json:"data,omitempty"`
}

func apiSuccess(data any) *apiResponse {
	return &apiResponse{
		Data: data,
	}
}

func apiError(err error) *apiResponse {
	unwrap := errs.Unwrap(err)
	var dlt string
	if unwrap != err {
		dlt = fmt.Sprintf("%+v", dlt)
	}
	if codeErr, ok := unwrap.(errs.CodeError); ok {
		return &apiResponse{ErrCode: codeErr.Code(), ErrMsg: codeErr.Msg(), ErrDlt: dlt}
	}
	return &apiResponse{ErrCode: errs.ServerInternalError, ErrMsg: err.Error(), ErrDlt: dlt}
}
