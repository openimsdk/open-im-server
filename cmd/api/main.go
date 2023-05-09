package main

import (
	"bytes"
	"context"
	"fmt"
	"net"
	"strconv"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/OpenIMSDK/Open-IM-Server/internal/api"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/cmd"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/cache"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/openKeeper"
)

func main() {
	apiCmd := cmd.NewApiCmd()
	apiCmd.AddPortFlag()
	apiCmd.AddApi(run)
	if err := apiCmd.Execute(); err != nil {
		panic(err.Error())
	}
}

func run(port int) error {
	if port == 0 {
		port = config.Config.Api.GinPort[0]
	}
	var err error
	rdb, err := cache.NewRedis()
	if err != nil {
		return err
	}
	var client discoveryregistry.SvcDiscoveryRegistry
	client, err = openKeeper.NewClient(config.Config.Zookeeper.ZkAddr, config.Config.Zookeeper.Schema,
		openKeeper.WithFreq(time.Hour), openKeeper.WithUserNameAndPassword(config.Config.Zookeeper.UserName,
			config.Config.Zookeeper.Password), openKeeper.WithTimeout(10))
	if err != nil {
		return err
	}
	buf := bytes.NewBuffer(nil)
	if err := yaml.NewEncoder(buf).Encode(config.Config); err != nil {
		return err
	}
	if err := client.RegisterConf2Registry(constant.OpenIMCommonConfigKey, buf.Bytes()); err != nil {
		return err
	}
	log.NewPrivateLog(constant.LogFileName)
	router := api.NewGinRouter(client, rdb)
	var address string
	if config.Config.Api.ListenIP != "" {
		address = net.JoinHostPort(config.Config.Api.ListenIP, strconv.Itoa(port))
	} else {
		address = net.JoinHostPort("0.0.0.0", strconv.Itoa(port))
	}
	fmt.Println("start api server, address: ", address, ", OpenIM version: ", config.Version)
	log.ZInfo(context.Background(), "start server success", "address", address, "version", config.Version)
	err = router.Run(address)
	if err != nil {
		log.ZError(context.Background(), "api run failed ", err, "address", address)
		return err
	}
	return nil
}
