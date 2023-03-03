package main

import (
	"OpenIM/internal/msgtransfer"
	"OpenIM/pkg/common/config"
	"OpenIM/pkg/common/constant"
	"OpenIM/pkg/common/log"
	"flag"
	"sync"
)

func main() {
	if err := config.InitConfig(); err != nil {
		panic(err.Error())
	}
	var wg sync.WaitGroup
	wg.Add(1)
	prometheusPort := flag.Int("prometheus_port", config.Config.Prometheus.MessageTransferPrometheusPort[0], "MessageTransferPrometheusPort default listen port")
	log.NewPrivateLog(constant.LogFileName)
	if err := msgtransfer.StartTransfer(*prometheusPort); err != nil {
		panic(err)
	}
	wg.Wait()
}
