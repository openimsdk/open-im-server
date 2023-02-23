package main

import (
	"OpenIM/internal/rpc/user"
	"OpenIM/internal/startrpc"
	"OpenIM/pkg/common/config"
)

func main() {

	startrpc.Start(config.Config.RpcPort.OpenImUserPort[0], config.Config.RpcRegisterName.OpenImUserName, config.Config.Prometheus.UserPrometheusPort[0], user.Start)

}
