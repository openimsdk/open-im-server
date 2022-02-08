package main

import (
	rpcMessageCMS "Open_IM/internal/rpc/message_cms"
	"flag"
)

func main() {
	rpcPort := flag.Int("port", 11200, "rpc listening port")
	flag.Parse()
	rpcServer := rpcMessageCMS.NewMessageCMSServer(*rpcPort)
	rpcServer.Run()
}
