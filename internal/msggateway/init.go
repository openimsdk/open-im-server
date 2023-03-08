package msggateway

import (
	"OpenIM/pkg/common/config"
	"OpenIM/pkg/common/constant"
	"OpenIM/pkg/common/log"
	"fmt"
	"sync"
	"time"
)

func RunWsAndServer(rpcPort, wsPort, prometheusPort int) error {
	var wg sync.WaitGroup
	wg.Add(1)
	log.NewPrivateLog(constant.LogFileName)
	fmt.Println("start rpc/msg_gateway server, port: ", rpcPort, wsPort, prometheusPort, ", OpenIM version: ", config.Version)
	longServer, err := NewWsServer(
		WithPort(wsPort),
		WithMaxConnNum(int64(config.Config.LongConnSvr.WebsocketMaxConnNum)),
		WithHandshakeTimeout(time.Duration(config.Config.LongConnSvr.WebsocketTimeOut)*time.Second),
		WithMessageMaxMsgLength(config.Config.LongConnSvr.WebsocketMaxMsgLen))
	if err != nil {
		return err
	}
	hubServer := NewServer(rpcPort, longServer)
	go hubServer.Start()
	go hubServer.LongConnServer.Run()
	wg.Wait()
	return nil
}
