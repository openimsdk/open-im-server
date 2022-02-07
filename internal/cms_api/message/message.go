package message

import (
	openIMHttp "Open_IM/pkg/common/http"

	"Open_IM/pkg/common/constant"

	"github.com/gin-gonic/gin"
)

func BroadcastMessage(c *gin.Context) {
	openIMHttp.RespHttp200(c, constant.OK, nil)
}

func MassSendMassage(c *gin.Context) {
	openIMHttp.RespHttp200(c, constant.OK, nil)
}

func WithdrawMessage(c *gin.Context) {
	openIMHttp.RespHttp200(c, constant.OK, nil)
}
