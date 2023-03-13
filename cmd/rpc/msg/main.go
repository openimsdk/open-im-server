package main

import (
	"OpenIM/internal/rpc/msg"
	"OpenIM/pkg/common/cmd"
	"OpenIM/pkg/common/config"
)

func main() {
	rpcCmd := cmd.NewRpcCmd("msg")
	rpcCmd.AddPortFlag()
	rpcCmd.AddPrometheusPortFlag()
	if err := rpcCmd.Exec(); err != nil {
		panic(err.Error())
	}
	if err := rpcCmd.StartSvr(config.Config.RpcRegisterName.OpenImMsgName, msg.Start); err != nil {
		panic(err.Error())
	}
}
