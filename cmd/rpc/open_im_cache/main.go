package main

import (
	rpcCache "Open_IM/internal/rpc/cache"
	"Open_IM/pkg/common/config"
	"flag"
	"fmt"
)

func main() {
	defaultPorts := config.Config.RpcPort.OpenImCachePort
	rpcPort := flag.Int("port", defaultPorts[0], "RpcToken default listen port 10800")
	flag.Parse()
	fmt.Println("start auth rpc server, port: ", *rpcPort)
	rpcServer := rpcCache.NewCacheServer(*rpcPort)
	rpcServer.Run()

}
