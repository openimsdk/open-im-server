package main

import (
	rpcConversation "Open_IM/internal/rpc/conversation"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	promePkg "Open_IM/pkg/common/prometheus"
	"flag"
	"fmt"
)

func main() {
	defaultPorts := config.Config.RpcPort.OpenImConversationPort
	rpcPort := flag.Int("port", defaultPorts[0], "RpcConversation default listen port 11300")
	prometheusPort := flag.Int("prometheus_port", config.Config.Prometheus.ConversationPrometheusPort[0], "conversationPrometheusPort default listen port")
	flag.Parse()
	fmt.Println("start conversation rpc server, port: ", *rpcPort, ", OpenIM version: ", constant.CurrentVersion)
	rpcServer := rpcConversation.NewRpcConversationServer(*rpcPort)
	go func() {
		err := promePkg.StartPromeSrv(*prometheusPort)
		if err != nil {
			panic(err)
		}
	}()
	rpcServer.Run()

}
