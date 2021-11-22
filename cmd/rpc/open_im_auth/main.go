package main

import (
	rpcAuth "Open_IM/internal/rpc/auth"
	"flag"
)

func main() {
	rpcPort := flag.Int("port", 10600, "RpcToken default listen port 10600")
	flag.Parse()
	rpcServer := rpcAuth.NewRpcAuthServer(*rpcPort)
	rpcServer.Run()
}
