package main

import (
	"OpenIM/internal/api"
	"OpenIM/pkg/common/cmd"
	"OpenIM/pkg/common/config"
	"OpenIM/pkg/common/log"
	log2 "OpenIM/pkg/common/logger"
	"OpenIM/pkg/common/mw"
	"context"
	"fmt"
	"github.com/OpenIMSDK/openKeeper"
	"os"
	"strconv"

	"OpenIM/pkg/common/constant"
)

func main() {
	apiCmd := cmd.NewApiCmd()
	apiCmd.AddPortFlag()
	apiCmd.AddApi(run)
	if err := apiCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func run(port int) error {
	if port == 0 {
		port = config.Config.Api.GinPort[0]
	}
	zk, err := openKeeper.NewClient(config.Config.Zookeeper.ZkAddr, config.Config.Zookeeper.Schema, 10, config.Config.Zookeeper.UserName, config.Config.Zookeeper.Password)
	if err != nil {
		return err
	}
	log.NewPrivateLog(constant.LogFileName)
	zk.AddOption(mw.GrpcClient())
	router := api.NewGinRouter(zk)
	address := constant.LocalHost + ":" + strconv.Itoa(port)
	if config.Config.Api.ListenIP != "" {
		address = config.Config.Api.ListenIP + ":" + strconv.Itoa(port)
	}
	fmt.Println("start api server, address: ", address, ", OpenIM version: ", config.Version)
	log2.Info(context.Background(), "start server success", "address", address, "version", config.Version)
	log.Info("s", "start server")
	err = router.Run(address)
	if err != nil {
		log.Error("", "api run failed ", address, err.Error())
		return err
	}
	return nil
}
