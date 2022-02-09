package http

import (
	"Open_IM/pkg/common/constant"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"


	//"Open_IM/pkg/cms_api_struct"
	"net/http"

	"github.com/gin-gonic/gin"

)

type BaseResp struct {
	Code   int32       `json:"code"`
	ErrMsg string      `json:"err_msg"`
	Data   interface{} `json:"data"`
}

func RespHttp200(ctx *gin.Context, err error, data interface{}) {
	var resp BaseResp
	if v, ok := err.(constant.ErrInfo); ok {
		resp.Code = v.Code()
		resp.ErrMsg = v.Error()
	} else {
		resp.Code = constant.ErrServer.Code()
		resp.ErrMsg = constant.ErrServer.Error()
	}
	resp.Data=data
	ctx.JSON(http.StatusOK, resp)
}

// warp error
func WarpError(err constant.ErrInfo) error {
	return status.Error(codes.Code(err.ErrCode), err.ErrMsg)
}

// parse error from server
func RespHttp200S(ctx *gin.Context, err error, data interface{}) {
	var resp BaseResp
	switch e := err.(type) {
	case constant.ErrInfo:
		resp.Code = e.ErrCode
	default:
		s, ok := status.FromError(err)
		if !ok {
			return
		}
		resp.Code = int32(s.Code())
		resp.ErrMsg = s.Message()
	}
	resp.Data = data
	ctx.JSON(http.StatusOK, resp)
}

//func CheckErr(pb interface{}) constant.ErrInfo{
//
//}