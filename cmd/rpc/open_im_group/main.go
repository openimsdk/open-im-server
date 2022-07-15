package main

import (
	"Open_IM/internal/rpc/group"
	"Open_IM/pkg/common/config"
	"flag"
	"fmt"
)

func main() {
	defaultPorts := config.Config.RpcPort.OpenImGroupPort
	rpcPort := flag.Int("port", defaultPorts[0], "get RpcGroupPort from cmd,default 16000 as port")
	flag.Parse()
	fmt.Println("start group rpc server, port: ", *rpcPort)
	rpcServer := group.NewGroupServer(*rpcPort)
	rpcServer.Run()
}
