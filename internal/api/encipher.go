package api

import (
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/gin-gonic/gin"
	"net/http"
)

func IsEncipher(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"errCode": 0, "errMsg": "", "data": config.Config.Encipher})
}
