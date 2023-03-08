package main

import (
	"OpenIM/internal/rpc/third"
	"OpenIM/pkg/common/cmd"
	"OpenIM/pkg/common/config"
	"fmt"
	"os"
)

func main() {
	rpcCmd := cmd.NewRpcCmd(config.Config.RpcRegisterName.OpenImThirdName)
	rpcCmd.AddPortFlag()
	rpcCmd.AddPrometheusPortFlag()
	if err := rpcCmd.Exec(third.Start); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
