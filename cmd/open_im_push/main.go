package main

import (
	"Open_IM/internal/push/logic"
	"flag"
	"sync"
)

func main() {
	rpcPort := flag.Int("port", -1, "rpc listening port")
	flag.Parse()
	var wg sync.WaitGroup
	wg.Add(1)
	logic.Init(*rpcPort)
	logic.Run()
	wg.Wait()
}
