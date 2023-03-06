package main

import (
	"OpenIM/internal/push"
	"OpenIM/internal/startrpc"
	"OpenIM/pkg/common/config"
	"OpenIM/pkg/common/constant"
	"OpenIM/pkg/common/log"
)

func main() {
	if err := config.InitConfig(); err != nil {
		panic(err.Error())
	}
	log.NewPrivateLog(constant.LogFileName)
	if err := startrpc.Start(config.Config.RpcPort.OpenImAuthPort[0], config.Config.RpcRegisterName.OpenImAuthName, config.Config.Prometheus.AuthPrometheusPort[0], push.Start); err != nil {
		panic(err.Error())
	}
}
