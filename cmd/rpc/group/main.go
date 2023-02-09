package main

import (
	"Open_IM/internal/rpc/group"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	promePkg "Open_IM/pkg/common/prometheus"
	"flag"
	"fmt"
)

func main() {
	defaultPorts := config.Config.RpcPort.OpenImGroupPort
	rpcPort := flag.Int("port", defaultPorts[0], "get RpcGroupPort from cmd,default 16000 as port")
	prometheusPort := flag.Int("prometheus_port", config.Config.Prometheus.GroupPrometheusPort[0], "groupPrometheusPort default listen port")
	flag.Parse()
	fmt.Println("start group rpc server, port: ", *rpcPort, ", OpenIM version: ", constant.CurrentVersion, "\n")
	rpcServer := group.NewGroupServer(*rpcPort)
	go func() {
		err := promePkg.StartPromeSrv(*prometheusPort)
		if err != nil {
			panic(err)
		}
	}()
	rpcServer.Run()
}
