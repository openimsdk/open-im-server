package main

import (
	"OpenIM/internal/rpc/msg"
	"OpenIM/internal/startrpc"
	"OpenIM/pkg/common/config"
)

func main() {
	startrpc.Start(config.Config.RpcPort.OpenImMessagePort, config.Config.RpcRegisterName.OpenImMsgName, config.Config.Prometheus.AuthPrometheusPort, msg.Start)
}
