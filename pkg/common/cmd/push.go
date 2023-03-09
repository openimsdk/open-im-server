package cmd

import (
	"OpenIM/internal/push"
	"OpenIM/internal/startrpc"
	"OpenIM/pkg/common/config"
	"github.com/spf13/cobra"
)

type PushCmd struct {
	*AuthCmd
}

func NewPushCmd() *PushCmd {
	return &PushCmd{NewAuthCmd()}
}

func (r *PushCmd) AddPush() {
	r.Command.RunE = func(cmd *cobra.Command, args []string) error {
		return startrpc.Start(r.getPortFlag(cmd), config.Config.RpcRegisterName.OpenImPushName, r.getPrometheusPortFlag(cmd), push.Start)
	}
}
