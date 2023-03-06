package main

import (
	"OpenIM/internal/rpc/group"
	"OpenIM/internal/startrpc"
	"OpenIM/pkg/common/config"
)

func main() {
	if err := config.InitConfig(); err != nil {
		panic(err.Error())
	}
	if err := startrpc.Start(config.Config.RpcPort.OpenImGroupPort[0], config.Config.RpcRegisterName.OpenImGroupName, config.Config.Prometheus.GroupPrometheusPort[0], group.Start); err != nil {
		panic(err.Error())
	}
}
