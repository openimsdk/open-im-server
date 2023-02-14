package main

import (
	"Open_IM/internal/rpc/group"
	"Open_IM/internal/startrpc"
	"Open_IM/pkg/common/config"
)

func main() {
	startrpc.Start(config.Config.RpcPort.OpenImGroupPort[0], config.Config.RpcRegisterName.OpenImGroupName, config.Config.Prometheus.GroupPrometheusPort[0], group.Start)
}
