package main

import (
	"OpenIM/internal/rpc/auth"
	"OpenIM/pkg/common/cmd"
	"OpenIM/pkg/common/config"
	"fmt"
	"os"
)

func main() {
	rpcCmd := cmd.NewRpcCmd(config.Config.RpcRegisterName.OpenImAuthName)
	rpcCmd.AddPortFlag()
	rpcCmd.AddPrometheusPortFlag()
	if err := rpcCmd.Exec(auth.Start); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
