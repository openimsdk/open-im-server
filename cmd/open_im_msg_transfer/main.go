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
	prometheusPort := flag.Int("prometheus_port", config.Config.Prometheus.MessageTransferPrometheusPort[0], "MessageTransferPrometheusPort default listen port")
	flag.Parse()
	log.NewPrivateLog(constant.LogFileName)
	logic.Init()
	fmt.Println("start msg_transfer server ", ", OpenIM version: ", constant.CurrentVersion, "\n")
	logic.Run(*prometheusPort)
	wg.Wait()
}
