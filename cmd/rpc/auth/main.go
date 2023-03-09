package main

import (
	"OpenIM/internal/rpc/auth"
	"OpenIM/pkg/common/cmd"
	"OpenIM/pkg/common/config"
	"fmt"
	"os"
)

func main() {
	rpcCmd := cmd.NewRpcCmd()
	rpcCmd.AddPortFlag()
	rpcCmd.AddPrometheusPortFlag()
	if err := rpcCmd.Exec2(config.Config.RpcRegisterName.OpenImAuthName, auth.Start); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
