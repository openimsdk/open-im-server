package main

import (
	rpcAuth "Open_IM/internal/rpc/auth"
	"Open_IM/pkg/common/config"
	promePkg "Open_IM/pkg/common/prometheus"
	"flag"
	"fmt"
)

func main() {
	defaultPorts := config.Config.RpcPort.OpenImAuthPort
	rpcPort := flag.Int("port", defaultPorts[0], "RpcToken default listen port 10800")
	prometheusPort := flag.Int("promethus-port", config.Config.Prometheus.AuthPrometheusPort[0], "authPrometheusPort default listen port")
	flag.Parse()
	fmt.Println("start auth rpc server, port: ", *rpcPort)
	rpcServer := rpcAuth.NewRpcAuthServer(*rpcPort)
	go func() {
		err := promePkg.StartPromeSrv(*prometheusPort)
		if err != nil {
			panic(err)
		}
	}()
	rpcServer.Run()
}
