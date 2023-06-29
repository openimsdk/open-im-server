package main

import (
	rpc "Open_IM/internal/rpc/office"
	"flag"
	"fmt"
)

func main() {
	rpcPort := flag.Int("port", 11100, "rpc listening port")
	flag.Parse()
	fmt.Println("start office rpc server, port: ", *rpcPort)
	rpcServer := rpc.NewOfficeServer(*rpcPort)
	rpcServer.Run()
}
