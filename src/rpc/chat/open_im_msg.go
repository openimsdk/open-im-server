package main

import (
	rpcChat "Open_IM/src/rpc/chat/chat"
	"Open_IM/src/utils"
	"flag"
)

func main() {
	rpcPort := flag.String("port", "", "rpc listening port")
	flag.Parse()
	rpcServer := rpcChat.NewRpcChatServer(utils.StringToInt(*rpcPort))
	rpcServer.Run()
}
