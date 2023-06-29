package main

import (
	"Open_IM/internal/msg_gateway/gate"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	"flag"
	"fmt"
	"sync"
)

func main() {
	log.NewPrivateLog(constant.LogFileName)
	rpcPort := flag.Int("rpc_port", 10400, "rpc listening port")
	wsPort := flag.Int("ws_port", 17778, "ws listening port")
	flag.Parse()
	var wg sync.WaitGroup
	wg.Add(1)
	fmt.Println("start rpc/msg_gateway server, port: ", *rpcPort, *wsPort)
	gate.Init(*rpcPort, *wsPort)
	gate.Run()
	wg.Wait()
}
