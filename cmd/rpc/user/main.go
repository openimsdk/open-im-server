package main

import (
	"OpenIM/internal/rpc/user"
	"OpenIM/internal/startrpc"
	"OpenIM/pkg/common/config"
)

func main() {
	if err := config.InitConfig(); err != nil {
		panic(err.Error())
	}
	if err := startrpc.Start(config.Config.RpcPort.OpenImUserPort[0], config.Config.RpcRegisterName.OpenImUserName, config.Config.Prometheus.UserPrometheusPort[0], user.Start); err != nil {
		panic(err.Error())
	}
}
