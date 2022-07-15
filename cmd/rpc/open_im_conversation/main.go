package main

import (
	rpcConversation "Open_IM/internal/rpc/conversation"
	"Open_IM/pkg/common/config"
	"flag"
	"fmt"
)

func main() {
	defaultPorts := config.Config.RpcPort.OpenImConversationPort
	rpcPort := flag.Int("port", defaultPorts[0], "RpcConversation default listen port 11300")
	flag.Parse()
	fmt.Println("start conversation rpc server, port: ", *rpcPort)
	rpcServer := rpcConversation.NewRpcConversationServer(*rpcPort)
	rpcServer.Run()

}
