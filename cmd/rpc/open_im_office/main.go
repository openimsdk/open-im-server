package main

import (
	rpc "Open_IM/internal/rpc/office"
	"Open_IM/pkg/common/config"
	promePkg "Open_IM/pkg/common/prometheus"
	"flag"
	"fmt"
)

func main() {
	defaultPorts := config.Config.RpcPort.OpenImOfficePort
	rpcPort := flag.Int("port", defaultPorts[0], "rpc listening port")
	prometheusPort := flag.Int("prometheus_port", config.Config.Prometheus.OfficePrometheusPort[0], "officePrometheusPort default listen port")
	flag.Parse()
	fmt.Println("start office rpc server, port: ", *rpcPort, "\n")
	rpcServer := rpc.NewOfficeServer(*rpcPort)
	go func() {
		err := promePkg.StartPromeSrv(*prometheusPort)
		if err != nil {
			panic(err)
		}
	}()
	rpcServer.Run()
}
