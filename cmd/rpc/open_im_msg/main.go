package main

import (
	rpcChat "Open_IM/internal/rpc/chat"
	"Open_IM/pkg/utils"
	"flag"
)

func main() {
	rpcPort := flag.String("port", "", "rpc listening port")
	flag.Parse()
	rpcServer := rpcChat.NewRpcChatServer(utils.StringToInt(*rpcPort))
	rpcServer.Run()
}
