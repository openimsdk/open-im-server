package main

import (
	"OpenIM/internal/rpc/friend"
	"OpenIM/pkg/common/cmd"
	"OpenIM/pkg/common/config"
)

func main() {
	rpcCmd := cmd.NewRpcCmd("friend")
	rpcCmd.AddPortFlag()
	rpcCmd.AddPrometheusPortFlag()
	if err := rpcCmd.Exec(); err != nil {
		panic(err.Error())
	}
	if err := rpcCmd.StartSvr(config.Config.RpcRegisterName.OpenImFriendName, friend.Start); err != nil {
		panic(err.Error())
	}
}
