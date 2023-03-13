package main

import (
	"OpenIM/internal/rpc/auth"
	"OpenIM/pkg/common/cmd"
	"OpenIM/pkg/common/config"
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
