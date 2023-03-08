package main

import (
	"OpenIM/internal/api"
	"OpenIM/pkg/common/cmd"
	"OpenIM/pkg/common/config"
	"OpenIM/pkg/common/log"
	"OpenIM/pkg/common/mw"
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
	zk, err := openKeeper.NewClient(config.Config.Zookeeper.ZkAddr, "", 10, "", "")
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
	err = router.Run(address)
	if err != nil {
		log.Error("", "api run failed ", address, err.Error())
		return err
	}
	return nil
}
