package main

import (
	rpcCache "Open_IM/internal/rpc/cache"
	"Open_IM/pkg/common/config"
	promePkg "Open_IM/pkg/common/prometheus"

	"flag"
	"fmt"
)

func main() {
	defaultPorts := config.Config.RpcPort.OpenImCachePort
	rpcPort := flag.Int("port", defaultPorts[0], "RpcToken default listen port 10800")
	prometheusPort := flag.Int("promethus-port", config.Config.Prometheus.CachePrometheusPort[0], "cachePrometheusPort default listen port")
	flag.Parse()
	fmt.Println("start cache rpc server, port: ", *rpcPort)
	rpcServer := rpcCache.NewCacheServer(*rpcPort)
	go func() {
		err := promePkg.StartPromeSrv(*prometheusPort)
		if err != nil {
			panic(err)
		}
	}()
	rpcServer.Run()
}
