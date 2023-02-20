package main

import (
	"Open_IM/internal/rpc/conversation"
	"Open_IM/internal/startrpc"
	"Open_IM/pkg/common/config"
)

func main() {
	startrpc.Start(config.Config.RpcPort.OpenImConversationPort, config.Config.RpcRegisterName.OpenImConversationName, config.Config.Prometheus.ConversationPrometheusPort, conversation.Start)
}
