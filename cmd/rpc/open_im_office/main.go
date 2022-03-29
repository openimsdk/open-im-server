package main

import (
	rpc "Open_IM/internal/rpc/office"
	"flag"
)

func main() {
	rpcPort := flag.Int("port", 11100, "rpc listening port")
	flag.Parse()
	rpcServer := rpc.NewOfficeServer(*rpcPort)
	rpcServer.Run()
}
