package main

import (
	"OpenIM/internal/rpc/user"
	"OpenIM/pkg/common/cmd"
	"OpenIM/pkg/common/config"
)

func main() {
	rpcCmd := cmd.NewRpcCmd("user")
	rpcCmd.AddPortFlag()
	rpcCmd.AddPrometheusPortFlag()
	if err := rpcCmd.Exec(); err != nil {
		panic(err.Error())
	}
	if err := rpcCmd.StartSvr(config.Config.RpcRegisterName.OpenImUserName, user.Start); err != nil {
		panic(err.Error())
	}
}
