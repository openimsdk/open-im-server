package main

import (
	"Open_IM/internal/rpc/organization"
	"Open_IM/pkg/common/config"
	"flag"
	"fmt"
)

func main() {
	defaultPorts := config.Config.RpcPort.OpenImOrganizationPort
	rpcPort := flag.Int("port", defaultPorts[0], "get RpcOrganizationPort from cmd,default 11200 as port")
	flag.Parse()
	fmt.Println("start organization rpc server, port: ", *rpcPort)
	rpcServer := organization.NewServer(*rpcPort)
	rpcServer.Run()
}
