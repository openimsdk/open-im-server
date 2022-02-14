package main

import (
	"Open_IM/internal/rpc/statistics"
	"flag"
)

func main() {
	rpcPort := flag.Int("port", 10800, "rpc listening port")
	flag.Parse()
	rpcServer := statistics.NewStatisticsServer(*rpcPort)
	rpcServer.Run()
}
