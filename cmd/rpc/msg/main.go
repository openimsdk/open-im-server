package main

import (
	"OpenIM/internal/rpc/msg"
	"OpenIM/pkg/common/cmd"
	"OpenIM/pkg/common/config"
	"fmt"
	"os"
)

func main() {
	rpcCmd := cmd.NewRpcCmd(config.Config.RpcRegisterName.OpenImMsgName)
	rpcCmd.AddPortFlag()
	rpcCmd.AddPrometheusPortFlag()
	if err := rpcCmd.Exec(msg.Start); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
