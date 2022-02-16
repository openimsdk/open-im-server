package main

import (
	rpcChat "Open_IM/internal/rpc/msg"
	"flag"
)

func main() {
	rpcPort := flag.Int("port", 10300, "rpc listening port")
	flag.Parse()
	rpcServer := rpcChat.NewRpcChatServer(*rpcPort)
	rpcServer.Run()
}
