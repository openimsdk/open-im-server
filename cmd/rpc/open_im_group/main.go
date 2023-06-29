package main

import (
	"Open_IM/internal/rpc/group"
	"flag"
	"fmt"
)

func main() {
	rpcPort := flag.Int("port", 10500, "get RpcGroupPort from cmd,default 16000 as port")
	flag.Parse()
	fmt.Println("start group rpc server, port: ", *rpcPort)
	rpcServer := group.NewGroupServer(*rpcPort)
	rpcServer.Run()
}
