package main

import (
	"Open_IM/internal/cms_api"
	"Open_IM/pkg/utils"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := cms_api.NewGinRouter()
	router.Use(utils.CorsHandler())
	router.Run(":" + "8000")
}
