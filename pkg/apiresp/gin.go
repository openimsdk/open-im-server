package apiresp

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func GinError(c *gin.Context, err error) {
	if err == nil {
		GinSuccess(c, nil)
		return
	}
	c.JSON(http.StatusOK, apiError(err))
}

func GinSuccess(c *gin.Context, data any) {
	c.JSON(http.StatusOK, apiSuccess(data))
}
