package main

import (
	"OpenIM/internal/push"
	"OpenIM/pkg/common/config"
	"OpenIM/pkg/common/constant"
	"OpenIM/pkg/common/log"
	"flag"
	"fmt"
	"sync"
)

func main() {
	defaultPorts := config.Config.RpcPort.OpenImPushPort
	rpcPort := flag.Int("port", defaultPorts[0], "rpc listening port")
	prometheusPort := flag.Int("prometheus_port", config.Config.Prometheus.MessageTransferPrometheusPort[0], "PushrometheusPort default listen port")
	configPath := flag.String("config_path", "../config/", "config folder")
	flag.Parse()
	if err := config.InitConfig(*configPath); err != nil {
		panic(err.Error())
	}
	var wg sync.WaitGroup
	wg.Add(1)
	log.NewPrivateLog(constant.LogFileName)
	fmt.Println("start push rpc server, port: ", *rpcPort, ", OpenIM version: ", constant.CurrentVersion, "\n")
	pusher := push.Push{}
	if err := pusher.Init(*rpcPort); err != nil {
		panic(err.Error())
	}
	pusher.Run(*prometheusPort)
	wg.Wait()
}
