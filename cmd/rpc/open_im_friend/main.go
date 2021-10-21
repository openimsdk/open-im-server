package main

import (
	"Open_IM/internal/rpc/friend"
	"flag"
)

func main() {

	rpcPort := flag.Int("port", 10200, "get RpcFriendPort from cmd,default 12000 as port")
	flag.Parse()
	rpcServer := friend.NewFriendServer(*rpcPort)
	rpcServer.Run()
}
