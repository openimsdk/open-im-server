package main

import (
	"OpenIM/internal/push"
	"OpenIM/pkg/common/cmd"
	"OpenIM/pkg/common/config"
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
