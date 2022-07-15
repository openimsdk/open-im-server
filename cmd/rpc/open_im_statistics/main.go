package main

import (
	"Open_IM/internal/rpc/statistics"
	"Open_IM/pkg/common/config"
	"flag"
	"fmt"
)

func main() {
	defaultPorts := config.Config.RpcPort.OpenImStatisticsPort
	rpcPort := flag.Int("port", defaultPorts[0], "rpc listening port")
	flag.Parse()
	fmt.Println("start statistics rpc server, port: ", *rpcPort)
	rpcServer := statistics.NewStatisticsServer(*rpcPort)
	rpcServer.Run()
}
