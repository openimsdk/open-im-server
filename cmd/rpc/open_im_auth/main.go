package main

import (
	rpcAuth "Open_IM/internal/rpc/auth"
	"flag"
	"fmt"
)

func main() {
	rpcPort := flag.Int("port", 10600, "RpcToken default listen port 10800")
	flag.Parse()
	fmt.Println("start auth rpc server, port: ", *rpcPort)
	rpcServer := rpcAuth.NewRpcAuthServer(*rpcPort)
	rpcServer.Run()

}
