package main

import (
	"OpenIM/internal/rpc/group"
	"OpenIM/internal/startrpc"
	"OpenIM/pkg/common/config"
)

func main() {
	startrpc.Start(config.Config.RpcPort.OpenImGroupPort, config.Config.RpcRegisterName.OpenImGroupName, config.Config.Prometheus.GroupPrometheusPort, group.Start)
}
