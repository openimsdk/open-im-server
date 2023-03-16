package main

import (
	"github.com/OpenIMSDK/Open-IM-Server/internal/rpc/third"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/cmd"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
)

func main() {
	rpcCmd := cmd.NewRpcCmd("third")
	rpcCmd.AddPortFlag()
	rpcCmd.AddPrometheusPortFlag()
	if err := rpcCmd.Exec(); err != nil {
		panic(err.Error())
	}
	if err := rpcCmd.StartSvr(config.Config.RpcRegisterName.OpenImThirdName, third.Start); err != nil {
		panic(err.Error())
	}
}
