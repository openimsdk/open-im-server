package main

import (
	"OpenIM/internal/msggateway"
	"OpenIM/pkg/common/config"
	"OpenIM/pkg/common/constant"
	"OpenIM/pkg/common/log"
	"flag"
	"fmt"
	"sync"
)

func main() {
	log.NewPrivateLog(constant.LogFileName)
	defaultRpcPorts := config.Config.RpcPort.OpenImMessageGatewayPort
	defaultWsPorts := config.Config.LongConnSvr.WebsocketPort
	defaultPromePorts := config.Config.Prometheus.MessageGatewayPrometheusPort
	rpcPort := flag.Int("rpc_port", defaultRpcPorts[0], "rpc listening port")
	wsPort := flag.Int("ws_port", defaultWsPorts[0], "ws listening port")
	prometheusPort := flag.Int("prometheus_port", defaultPromePorts[0], "PushrometheusPort default listen port")
	configPath := flag.String("config_path", "../config/", "config folder")
	flag.Parse()
	if err := config.InitConfig(*configPath); err != nil {
		panic(err.Error())
	}
	var wg sync.WaitGroup
	wg.Add(1)
	fmt.Println("start rpc/msg_gateway server, port: ", *rpcPort, *wsPort, *prometheusPort, ", OpenIM version: ", constant.CurrentVersion, "\n")
	msggateway.Init(*rpcPort, *wsPort)
	msggateway.Run(*prometheusPort)
	wg.Wait()
}
