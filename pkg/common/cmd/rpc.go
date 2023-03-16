package cmd

import (
	"github.com/OpenIMSDK/Open-IM-Server/internal/startrpc"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

type RpcCmd struct {
	*RootCmd
}

func NewRpcCmd(name string) *RpcCmd {
	authCmd := &RpcCmd{NewRootCmd(name)}
	return authCmd
}

func (a *RpcCmd) Exec() error {
	a.Command.Run = func(cmd *cobra.Command, args []string) {
		a.port = a.getPortFlag(cmd)
		a.prometheusPort = a.getPrometheusPortFlag(cmd)
	}
	return a.Execute()
}

func (a *RpcCmd) StartSvr(name string, rpcFn func(client discoveryregistry.SvcDiscoveryRegistry, server *grpc.Server) error) error {
	return startrpc.Start(a.GetPortFlag(), name, a.GetPrometheusPortFlag(), rpcFn)
}
