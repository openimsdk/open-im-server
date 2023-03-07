package main

import (
	"OpenIM/internal/rpc/auth"
	"OpenIM/internal/startrpc"
	"OpenIM/pkg/common/config"
)

func main() {
	if err := config.InitConfig(""); err != nil {
		panic(err.Error())
	}
	if err := startrpc.Start(config.Config.RpcPort.OpenImAuthPort[0], config.Config.RpcRegisterName.OpenImAuthName, config.Config.Prometheus.AuthPrometheusPort[0], auth.Start); err != nil {
		panic(err.Error())
	}
}
