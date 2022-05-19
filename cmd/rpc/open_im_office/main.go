package main

import (
	rpc "Open_IM/internal/rpc/office"
	"Open_IM/pkg/common/config"
	"flag"
	"fmt"
)

func main() {
	defaultPorts := config.Config.RpcPort.OpenImOfficePort
	rpcPort := flag.Int("port", defaultPorts[0], "rpc listening port")
	flag.Parse()
	fmt.Println("start office rpc server, port: ", *rpcPort)
	rpcServer := rpc.NewOfficeServer(*rpcPort)
	rpcServer.Run()
}
