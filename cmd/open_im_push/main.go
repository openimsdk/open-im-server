package main

import (
	"Open_IM/internal/push/logic"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	"flag"
	"fmt"
	"sync"
)

func main() {
	rpcPort := flag.Int("port", 10700, "rpc listening port")
	flag.Parse()
	var wg sync.WaitGroup
	wg.Add(1)
	log.NewPrivateLog(constant.LogFileName)
	fmt.Println("start push rpc server, port: ", *rpcPort)
	logic.Init(*rpcPort)
	logic.Run()
	wg.Wait()
}
