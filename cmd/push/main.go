package main

import (
	"github.com/OpenIMSDK/Open-IM-Server/internal/push"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/cmd"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
)

func main() {
	pushCmd := cmd.NewRpcCmd("push")
	pushCmd.AddPortFlag()
	pushCmd.AddPrometheusPortFlag()
	if err := pushCmd.Exec(); err != nil {
		panic(err.Error())
	}
	if err := pushCmd.StartSvr(config.Config.RpcRegisterName.OpenImPushName, push.Start); err != nil {
		panic(err.Error())
	}
}
