package main

import (
	rpcMessageCMS "Open_IM/internal/rpc/admin_cms"
	"flag"
	"fmt"
)

func main() {
	rpcPort := flag.Int("port", 11000, "rpc listening port")
	flag.Parse()
	fmt.Println("start cms rpc server, port: ", *rpcPort)
	rpcServer := rpcMessageCMS.NewAdminCMSServer(*rpcPort)
	rpcServer.Run()
}
