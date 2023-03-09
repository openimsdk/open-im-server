package cmd

import (
	"OpenIM/internal/startrpc"
	"OpenIM/pkg/discoveryregistry"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

type AuthCmd struct {
	*RootCmd
}

func NewAuthCmd() *AuthCmd {
	authCmd := &AuthCmd{NewRootCmd()}
	return authCmd
}

func (a *AuthCmd) Exec() error {
	a.Command.RunE = func(cmd *cobra.Command, args []string) error {
		a.port = a.getPortFlag(cmd)
		a.prometheusPort = a.getPrometheusPortFlag(cmd)
		return a.getConfFromCmdAndInit(cmd)
	}
	return a.Execute()
}

func (a *AuthCmd) StartSvr(name string, rpcFn func(client discoveryregistry.SvcDiscoveryRegistry, server *grpc.Server) error) error {
	return startrpc.Start(a.GetPortFlag(), name, a.GetPrometheusPortFlag(), rpcFn)
}
