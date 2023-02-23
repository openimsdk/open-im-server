package main

import (
	"Open_IM/internal/push"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	"flag"
	"fmt"
	"sync"
)

func main() {
	defaultPorts := config.Config.RpcPort.OpenImPushPort
	rpcPort := flag.Int("port", defaultPorts[0], "rpc listening port")
	prometheusPort := flag.Int("prometheus_port", config.Config.Prometheus.MessageTransferPrometheusPort[0], "PushrometheusPort default listen port")
	flag.Parse()
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
