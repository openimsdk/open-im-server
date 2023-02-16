package main

import (
	rpcConversation "Open_IM/internal/rpc/conversation"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	prome "Open_IM/pkg/common/prome"
	"flag"
	"fmt"
)

func main() {
	defaultPorts := config.Config.RpcPort.OpenImConversationPort
	rpcPort := flag.Int("port", defaultPorts[0], "RpcConversation default listen port 11300")
	prometheusPort := flag.Int("prometheus_port", config.Config.Prometheus.ConversationPrometheusPort[0], "conversationPrometheusPort default listen port")
	flag.Parse()
	fmt.Println("start conversation rpc server, port: ", *rpcPort, ", OpenIM version: ", constant.CurrentVersion, "\n")
	rpcServer := rpcConversation.NewRpcConversationServer(*rpcPort)
	go func() {
		err := prome.StartPromeSrv(*prometheusPort)
		if err != nil {
			panic(err)
		}
	}()
	rpcServer.Run()

}
