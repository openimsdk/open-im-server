package main

import (
	rpcAuth "Open_IM/internal/rpc/auth"
	"Open_IM/pkg/common/config"
	"flag"
	"fmt"
)

func main() {
	defaultPorts := config.Config.RpcPort.OpenImAuthPort
	rpcPort := flag.Int("port", defaultPorts[0], "RpcToken default listen port 10800")
	flag.Parse()
	fmt.Println("start auth rpc server, port: ", *rpcPort)
	rpcServer := rpcAuth.NewRpcAuthServer(*rpcPort)
	rpcServer.Run()

}
