package main

import (
	"Open_IM/internal/rpc/auth"
	"Open_IM/internal/startrpc"
	"Open_IM/pkg/common/config"
)

func main() {
	startrpc.Start(config.Config.RpcPort.OpenImAuthPort, config.Config.RpcRegisterName.OpenImAuthName, config.Config.Prometheus.AuthPrometheusPort, auth.Start)
}
