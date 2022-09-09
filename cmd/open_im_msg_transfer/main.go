package main

import (
	"Open_IM/internal/msg_transfer/logic"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	"flag"
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	rpcPort := flag.Int("port", config.Config.Prometheus.MessageTransferPrometheusPort[0], "MessageTransferPrometheusPort default listen port")
	log.NewPrivateLog(constant.LogFileName)
	logic.Init()
	fmt.Println("start msg_transfer server")
	logic.Run(*rpcPort)
	wg.Wait()
}
