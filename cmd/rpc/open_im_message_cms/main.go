package main

import (
	rpcMessageCMS "Open_IM/internal/rpc/message_cms"
	"Open_IM/pkg/common/config"
	"flag"
	"fmt"
)

func main() {
	defaultPorts := config.Config.RpcPort.OpenImMessageCmsPort[0]
	rpcPort := flag.Int("port", defaultPorts, "rpc listening port")
	flag.Parse()
	fmt.Println("start msg cms rpc server, port: ", *rpcPort)
	rpcServer := rpcMessageCMS.NewMessageCMSServer(*rpcPort)
	rpcServer.Run()
}
