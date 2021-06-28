package main

import (
	"Open_IM/src/rpc/group/group"
	"flag"
)

func main() {
	rpcPort := flag.Int("port", 10500, "get RpcGroupPort from cmd,default 16000 as port")
	flag.Parse()
	rpcServer := group.NewGroupServer(*rpcPort)
	rpcServer.Run()
}
