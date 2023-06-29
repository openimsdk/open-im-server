package main

import (
	"Open_IM/internal/msg_transfer/logic"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	"fmt"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	wg.Add(1)
	log.NewPrivateLog(constant.LogFileName)
	logic.Init()
	fmt.Println("start msg_transfer server")
	logic.Run()
	wg.Wait()
}
