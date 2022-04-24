package main

import (
	"Open_IM/internal/cms_api"
	"Open_IM/pkg/utils"
	"flag"
	"fmt"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {
	gin.SetMode(gin.ReleaseMode)
	router := cms_api.NewGinRouter()
	router.Use(utils.CorsHandler())
	ginPort := flag.Int("port", 8000, "get ginServerPort from cmd,default 10000 as port")
	fmt.Println("start cms api server, port: ", *ginPort)
	router.Run(":" + strconv.Itoa(*ginPort))
}
