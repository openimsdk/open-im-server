package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"runtime"
	"strconv"
	"time"

	"net/http"
	_ "net/http/pprof"

	"github.com/OpenIMSDK/Open-IM-Server/internal/api"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/cmd"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/cache"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	openKeeper "github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry/zookeeper"
)

func main() {
	apiCmd := cmd.NewApiCmd()
	apiCmd.AddPortFlag()
	apiCmd.AddApi(run)
	if err := apiCmd.Execute(); err != nil {
		panic(err.Error())
	}
}

func startPprof() {
	runtime.GOMAXPROCS(1)
	runtime.SetMutexProfileFraction(1)
	runtime.SetBlockProfileRate(1)
	if err := http.ListenAndServe(":6060", nil); err != nil {
		panic(err)
	}
	os.Exit(0)
}

func run(port int) error {
	if port == 0 {
		return fmt.Errorf("port is empty")
	}
	rdb, err := cache.NewRedis()
	if err != nil {
		return err
	}
	fmt.Println("api start init discov client")
	var client discoveryregistry.SvcDiscoveryRegistry
	client, err = openKeeper.NewClient(config.Config.Zookeeper.ZkAddr, config.Config.Zookeeper.Schema,
		openKeeper.WithFreq(time.Hour), openKeeper.WithUserNameAndPassword(config.Config.Zookeeper.Username,
			config.Config.Zookeeper.Password), openKeeper.WithRoundRobin(), openKeeper.WithTimeout(10), openKeeper.WithLogger(log.NewZkLogger()))
	if err != nil {
		return err
	}
	if client.CreateRpcRootNodes(config.GetServiceNames()); err != nil {
		return err
	}
	fmt.Println("api init discov client success")
	fmt.Println("api register public config to discov")
	if err := client.RegisterConf2Registry(constant.OpenIMCommonConfigKey, config.EncodeConfig()); err != nil {
		return err
	}
	fmt.Println("api register public config to discov success")
	router := api.NewGinRouter(client, rdb)
	fmt.Println("api init router success")
	var address string
	if config.Config.Api.ListenIP != "" {
		address = net.JoinHostPort(config.Config.Api.ListenIP, strconv.Itoa(port))
	} else {
		address = net.JoinHostPort("0.0.0.0", strconv.Itoa(port))
	}
	fmt.Println("start api server, address: ", address, ", OpenIM version: ", config.Version)
	log.ZInfo(context.Background(), "start server success", "address", address, "version", config.Version)
	go startPprof()
	err = router.Run(address)
	if err != nil {
		log.ZError(context.Background(), "api run failed ", err, "address", address)
		return err
	}
	return nil
}
