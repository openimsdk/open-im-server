package main

import (
	"Open_IM/internal/push/logic"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	"flag"
	"fmt"
	"sync"
)

func main() {
	defaultPorts := config.Config.RpcPort.OpenImPushPort
	rpcPort := flag.Int("port", defaultPorts[0], "rpc listening port")
	flag.Parse()
	var wg sync.WaitGroup
	wg.Add(1)
	log.NewPrivateLog(constant.LogFileName)
	fmt.Println("start push rpc server, port: ", *rpcPort)
	logic.Init(*rpcPort)
	logic.Run()
	wg.Wait()
}
