package main

import (
	"OpenIM/internal/rpc/auth"
	"OpenIM/internal/startrpc"
	"OpenIM/pkg/common/config"
)

func main() {
	startrpc.Start(config.Config.RpcPort.OpenImAuthPort, config.Config.RpcRegisterName.OpenImAuthName, config.Config.Prometheus.AuthPrometheusPort, auth.Start)
}
