package main

import (
	"OpenIM/internal/rpc/user"
	"OpenIM/pkg/common/cmd"
	"OpenIM/pkg/common/config"
	"fmt"
	"os"
)

func main() {
	rpcCmd := cmd.NewRpcCmd("user")
	rpcCmd.AddPortFlag()
	rpcCmd.AddPrometheusPortFlag()
	if err := rpcCmd.Exec(); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
	if err := rpcCmd.StartSvr(config.Config.RpcRegisterName.OpenImUserName, user.Start); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}
