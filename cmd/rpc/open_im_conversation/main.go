package main

import (
	rpcConversation "Open_IM/internal/rpc/conversation"
	"flag"
	"fmt"
)

func main() {
	rpcPort := flag.Int("port", 11300, "RpcConversation default listen port 11300")
	flag.Parse()
	fmt.Println("start conversation rpc server, port: ", *rpcPort)
	rpcServer := rpcConversation.NewRpcConversationServer(*rpcPort)
	rpcServer.Run()

}
