package cmd

import (
	"OpenIM/internal/push"
	"OpenIM/internal/startrpc"
	"OpenIM/pkg/common/config"
	"github.com/spf13/cobra"
)

type PushCmd struct {
	*RpcCmd
}

func NewPushCmd() *PushCmd {
	return &PushCmd{NewRpcCmd()}
}

func (r *RpcCmd) AddPush() {
	r.Command.RunE = func(cmd *cobra.Command, args []string) error {
		return startrpc.Start(r.getPortFlag(cmd), r.rpcRegisterName, r.getPrometheusPortFlag(cmd), push.Start)
	}
}
