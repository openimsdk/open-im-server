package main

import (
	"Open_IM/internal/demo/register"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/utils"
	"flag"
	"strconv"

	"github.com/gin-gonic/gin"
)

func main() {

	r := gin.Default()
	r.Use(utils.CorsHandler())

	authRouterGroup := r.Group("/auth")
	{
		authRouterGroup.POST("/code", register.SendVerificationCode)
		authRouterGroup.POST("/verify", register.Verify)
		authRouterGroup.POST("/password", register.SetPassword)
		authRouterGroup.POST("/login", register.Login)
		authRouterGroup.POST("/reset_password", register.ResetPassword)
	}
	log.NewPrivateLog("demo")
	ginPort := flag.Int("port", 42233, "get ginServerPort from cmd,default 42233 as port")
	flag.Parse()
	r.Run(":" + strconv.Itoa(*ginPort))
}
