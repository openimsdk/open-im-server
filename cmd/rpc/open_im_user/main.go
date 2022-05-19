package main

import (
	"Open_IM/internal/rpc/user"
	"Open_IM/pkg/common/config"
	"flag"
	"fmt"
)

func main() {
	defaultPorts := config.Config.RpcPort.OpenImUserPort
	rpcPort := flag.Int("port", defaultPorts[0], "rpc listening port")
	flag.Parse()
	fmt.Println("start user rpc server, port: ", *rpcPort)
	rpcServer := user.NewUserServer(*rpcPort)
	rpcServer.Run()
}
