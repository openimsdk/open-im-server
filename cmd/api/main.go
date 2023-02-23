package main

import (
	"OpenIM/internal/api"
	"OpenIM/internal/api/third"
	"OpenIM/pkg/common/config"
	"OpenIM/pkg/common/log"
	"flag"
	"fmt"

	"strconv"

	"OpenIM/pkg/common/constant"
)

func main() {
	log.NewPrivateLog(constant.LogFileName)
	router := api.NewGinRouter()
	go third.MinioInit()
	ginPort := flag.Int("port", config.Config.Api.GinPort[0], "get ginServerPort from cmd,default 10002 as port")
	flag.Parse()
	config.InitConfig()
	address := "0.0.0.0:" + strconv.Itoa(*ginPort)
	if config.Config.Api.ListenIP != "" {
		address = config.Config.Api.ListenIP + ":" + strconv.Itoa(*ginPort)
	}
	fmt.Println("start api server, address: ", address, ", OpenIM version: ", constant.CurrentVersion)
	err := router.Run(address)
	if err != nil {
		log.Error("", "api run failed ", address, err.Error())
		panic("api start failed " + err.Error())
	}
}
