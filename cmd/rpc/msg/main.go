package main

import (
	"Open_IM/internal/rpc/msg"
	"Open_IM/internal/startrpc"
	"Open_IM/pkg/common/config"
)

func main() {
	startrpc.Start(config.Config.RpcPort.OpenImMessagePort, config.Config.RpcRegisterName.OpenImMsgName, config.Config.Prometheus.AuthPrometheusPort, msg.Start)
}
