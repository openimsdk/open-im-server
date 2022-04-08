package main

import (
	rpcMessageCMS "Open_IM/internal/rpc/message_cms"
	"flag"
	"fmt"
)

func main() {
	rpcPort := flag.Int("port", 10900, "rpc listening port")
	flag.Parse()
	fmt.Println("start msg cms rpc server, port: ", *rpcPort)
	rpcServer := rpcMessageCMS.NewMessageCMSServer(*rpcPort)
	rpcServer.Run()
}
