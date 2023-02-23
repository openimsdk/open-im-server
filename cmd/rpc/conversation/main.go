package main

import (
	"OpenIM/internal/rpc/conversation"
	"OpenIM/internal/startrpc"
	"OpenIM/pkg/common/config"
)

func main() {
	startrpc.Start(config.Config.RpcPort.OpenImConversationPort, config.Config.RpcRegisterName.OpenImConversationName, config.Config.Prometheus.ConversationPrometheusPort, conversation.Start)
}
