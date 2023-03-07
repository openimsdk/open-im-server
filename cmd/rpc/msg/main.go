package main

import (
	"OpenIM/internal/rpc/msg"
	"OpenIM/internal/startrpc"
	"OpenIM/pkg/common/config"
)

func main() {
	if err := config.InitConfig(""); err != nil {
		panic(err.Error())
	}
	if err := startrpc.Start(config.Config.RpcPort.OpenImMessagePort[0], config.Config.RpcRegisterName.OpenImMsgName, config.Config.Prometheus.AuthPrometheusPort[0], msg.Start); err != nil {
		panic(err.Error())
	}
}
