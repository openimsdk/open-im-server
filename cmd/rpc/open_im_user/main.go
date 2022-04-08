package main

import (
	"Open_IM/internal/rpc/user"
	"flag"
	"fmt"
)

func main() {
	rpcPort := flag.Int("port", 10100, "rpc listening port")
	flag.Parse()
	fmt.Println("start user rpc server, port: ", *rpcPort)
	rpcServer := user.NewUserServer(*rpcPort)
	rpcServer.Run()
}
