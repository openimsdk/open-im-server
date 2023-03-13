package main

import (
	"OpenIM/internal/rpc/group"
	"OpenIM/pkg/common/cmd"
	"OpenIM/pkg/common/config"
)

func main() {
	rpcCmd := cmd.NewRpcCmd("group")
	rpcCmd.AddPortFlag()
	rpcCmd.AddPrometheusPortFlag()
	if err := rpcCmd.Exec(); err != nil {
		panic(err.Error())
	}
	if err := rpcCmd.StartSvr(config.Config.RpcRegisterName.OpenImGroupName, group.Start); err != nil {
		panic(err.Error())
	}
}
