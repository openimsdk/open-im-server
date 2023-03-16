package apiresp

import (
	"context"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
)

type apiResponse struct {
	ErrCode int    `json:"errCode"`
	ErrMsg  string `json:"errMsg"`
	ErrDlt  string `json:"errDlt"`
	Data    any    `json:"data,omitempty"`
}

func apiSuccess(data any) *apiResponse {
	log.ZDebug(context.Background(), "apiSuccess", "resp", &apiResponse{
		Data: data,
	})
	return &apiResponse{
		Data: data,
	}
}

func apiError(err error) *apiResponse {
	unwrap := errs.Unwrap(err)
	var dlt string
	if unwrap != err {
		dlt = err.Error()
	}
	if codeErr, ok := unwrap.(errs.CodeError); ok {
		return &apiResponse{ErrCode: codeErr.Code(), ErrMsg: codeErr.Msg(), ErrDlt: dlt}
	}
	return &apiResponse{ErrCode: errs.ServerInternalError, ErrMsg: err.Error(), ErrDlt: dlt}
}
