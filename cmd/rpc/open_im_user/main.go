package main

import (
	"Open_IM/internal/rpc/user"
	"flag"
)

func main() {
	rpcPort := flag.Int("port", 10100, "rpc listening port")
	flag.Parse()
	rpcServer := user.NewUserServer(*rpcPort)
	rpcServer.Run()
}
