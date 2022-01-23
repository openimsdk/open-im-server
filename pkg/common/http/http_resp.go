package http

import (
	"Open_IM/pkg/common/constant"
	"net/http"

	"github.com/gin-gonic/gin"
)

type BaseResp struct {
	Code   int32       `json:"code"`
	ErrMsg string      `json:"err_msg"`
	Data   interface{} `json:"data"`
}

func RespHttp200(ctx *gin.Context, err constant.ErrInfo, data interface{}) {
	resp := BaseResp{
		Code:   err.Code(),
		ErrMsg: err.Error(),
		Data:   data,
	}
	ctx.JSON(http.StatusOK, resp)
}
