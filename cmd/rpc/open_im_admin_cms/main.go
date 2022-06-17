package main

import (
	rpcMessageCMS "Open_IM/internal/rpc/admin_cms"
	"Open_IM/pkg/common/config"
	"flag"
	"fmt"
)

func main() {
	defaultPorts := config.Config.RpcPort.OpenImAdminCmsPort
	rpcPort := flag.Int("port", defaultPorts[0], "rpc listening port")
	flag.Parse()
	fmt.Println("start cms rpc server, port: ", *rpcPort)
	rpcServer := rpcMessageCMS.NewAdminCMSServer(*rpcPort)
	rpcServer.Run()
}
