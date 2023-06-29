package main

import (
	"Open_IM/internal/rpc/statistics"
	"flag"
	"fmt"
)

func main() {
	rpcPort := flag.Int("port", 10800, "rpc listening port")
	flag.Parse()
	fmt.Println("start statistics rpc server, port: ", *rpcPort)
	rpcServer := statistics.NewStatisticsServer(*rpcPort)
	rpcServer.Run()
}
