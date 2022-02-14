package main

import (
	rpcMessageCMS "Open_IM/internal/rpc/admin_cms"
	"flag"
)

func main() {
	rpcPort := flag.Int("port", 11000, "rpc listening port")
	flag.Parse()
	rpcServer := rpcMessageCMS.NewAdminCMSServer(*rpcPort)
	rpcServer.Run()
}

