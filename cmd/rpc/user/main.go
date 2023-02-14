package main

import (
	"Open_IM/internal/rpc/user"
	"Open_IM/internal/startrpc"
	"Open_IM/pkg/common/config"
)

func main() {

	startrpc.Start(config.Config.RpcPort.OpenImUserPort[0], config.Config.RpcRegisterName.OpenImUserName, config.Config.Prometheus.UserPrometheusPort[0], user.Start)

}
