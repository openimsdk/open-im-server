package main

import (
	"OpenIM/internal/rpc/conversation"
	"OpenIM/pkg/common/cmd"
	"OpenIM/pkg/common/config"
	"fmt"
	"os"
)

func main() {
	rpcCmd := cmd.NewRpcCmd(config.Config.RpcRegisterName.OpenImConversationName)
	rpcCmd.AddPortFlag()
	rpcCmd.AddPrometheusPortFlag()
	if err := rpcCmd.Exec(conversation.Start); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
