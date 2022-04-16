package main

import (
	"Open_IM/internal/rpc/group"
	"flag"
	"fmt"
)

func main() {
	rpcPort := flag.Int("port", 11200, "get RpcOrganizationPort from cmd,default 11200 as port")
	flag.Parse()
	fmt.Println("start organization rpc server, port: ", *rpcPort)
	rpcServer := group.NewGroupServer(*rpcPort)
	rpcServer.Run()
}
