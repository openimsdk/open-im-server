package main

import (
	"github.com/OpenIMSDK/Open-IM-Server/internal/rpc/auth"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/cmd"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
)

func main() {
	authCmd := cmd.NewRpcCmd("auth")
	authCmd.AddPortFlag()
	authCmd.AddPrometheusPortFlag()
	if err := authCmd.Exec(); err != nil {
		panic(err.Error())
	}
	if err := authCmd.StartSvr(config.Config.RpcRegisterName.OpenImAuthName, auth.Start); err != nil {
		panic(err.Error())
	}
}
