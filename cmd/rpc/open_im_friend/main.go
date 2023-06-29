package main

import (
	"Open_IM/internal/rpc/friend"
	"flag"
	"fmt"
)

func main() {

	rpcPort := flag.Int("port", 10200, "get RpcFriendPort from cmd,default 12000 as port")
	flag.Parse()
	fmt.Println("start friend rpc server, port: ", *rpcPort)
	rpcServer := friend.NewFriendServer(*rpcPort)
	rpcServer.Run()
}
