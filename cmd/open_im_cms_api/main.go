package main

import (
	"Open_IM/internal/cms_api"
	"Open_IM/pkg/utils"
	"fmt"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := cms_api.NewGinRouter()
	router.Use(utils.CorsHandler())
	fmt.Println("start cms api server, port: ", 8000)
	router.Run(":" + "8000")
}
