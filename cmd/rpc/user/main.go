package main

import (
	"OpenIM/internal/rpc/user"
	"OpenIM/pkg/common/cmd"
	"OpenIM/pkg/common/config"
	"fmt"
	"os"
)

func main() {
	rpcCmd := cmd.NewRpcCmd(config.Config.RpcRegisterName.OpenImUserName)
	rpcCmd.AddPortFlag()
	rpcCmd.AddPrometheusPortFlag()
	if err := rpcCmd.Exec(user.Start); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
