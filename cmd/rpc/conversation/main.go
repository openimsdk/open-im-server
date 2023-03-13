package main

import (
	"OpenIM/internal/rpc/conversation"
	"OpenIM/pkg/common/cmd"
	"OpenIM/pkg/common/config"
)

func main() {
	rpcCmd := cmd.NewRpcCmd("conversation")
	rpcCmd.AddPortFlag()
	rpcCmd.AddPrometheusPortFlag()
	if err := rpcCmd.Exec(); err != nil {
		panic(err.Error())
	}
	if err := rpcCmd.StartSvr(config.Config.RpcRegisterName.OpenImConversationName, conversation.Start); err != nil {
		panic(err.Error())
	}
}
