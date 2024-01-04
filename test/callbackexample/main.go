package main

import (
	"call-back-http/control"
	"github.com/gin-gonic/gin"
)

func main() {
	engine := gin.Default()
	router := engine.Group("/callback")
	router.POST("/callbackBeforeSendSingleMsgCommand", control.CallbackBeforeSendSingleMsgCommand)

	if err := engine.Run("0.0.0.0:18889"); err != nil {
		panic(err)
	}
}
