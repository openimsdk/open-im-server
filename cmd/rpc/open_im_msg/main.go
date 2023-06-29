package main

import (
	rpcChat "Open_IM/internal/rpc/msg"
	"flag"
	"fmt"
)

func main() {
	rpcPort := flag.Int("port", 10300, "rpc listening port")
	flag.Parse()
	fmt.Println("start msg rpc server, port: ", *rpcPort)
	rpcServer := rpcChat.NewRpcChatServer(*rpcPort)
	rpcServer.Run()
}
