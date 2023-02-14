package main

import (
	"Open_IM/internal/rpc/friend"
	"Open_IM/internal/startrpc"
	"Open_IM/pkg/common/config"
)

func main() {
	startrpc.Start(config.Config.RpcPort.OpenImFriendPort[0], config.Config.RpcRegisterName.OpenImFriendName, config.Config.Prometheus.FriendPrometheusPort[0], friend.Start)
}
