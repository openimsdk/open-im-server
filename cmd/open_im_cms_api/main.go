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
	ginPort := flag.Int("port", 10006, "get ginServerPort from cmd,default 8000 as port")
	flag.Parse()
	fmt.Println("start cms api server, port: ", ginPort)
	router.Run(":" + strconv.Itoa(*ginPort))
}

//
