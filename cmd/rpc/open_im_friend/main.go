package main

import (
	"Open_IM/internal/rpc/friend"
	"Open_IM/pkg/common/config"
	"flag"
	"fmt"
)

func main() {
	defaultPorts := config.Config.RpcPort.OpenImFriendPort
	rpcPort := flag.Int("port", defaultPorts[0], "get RpcFriendPort from cmd,default 12000 as port")
	flag.Parse()
	fmt.Println("start friend rpc server, port: ", *rpcPort)
	rpcServer := friend.NewFriendServer(*rpcPort)
	rpcServer.Run()
}
