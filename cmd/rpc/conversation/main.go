package main

import (
	"OpenIM/internal/rpc/conversation"
	"OpenIM/internal/startrpc"
	"OpenIM/pkg/common/config"
)

func main() {
	if err := config.InitConfig(); err != nil {
		panic(err.Error())
	}
	if err := startrpc.Start(config.Config.RpcPort.OpenImConversationPort[0], config.Config.RpcRegisterName.OpenImConversationName, config.Config.Prometheus.ConversationPrometheusPort[0], conversation.Start); err != nil {
		panic(err.Error())
	}
}
