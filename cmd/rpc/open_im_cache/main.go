package main

import (
	rpcCache "Open_IM/internal/rpc/cache"
	"flag"
	"fmt"
)

func main() {
	rpcPort := flag.Int("port", 10600, "RpcToken default listen port 10800")
	flag.Parse()
	fmt.Println("start auth rpc server, port: ", *rpcPort)
	rpcServer := rpcCache.NewOfficeServer(*rpcPort)
	rpcServer.Run()

}
