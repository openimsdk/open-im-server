package cmd

import (
	"OpenIM/internal/startrpc"
	"OpenIM/pkg/discoveryregistry"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

type RpcCmd struct {
	*RootCmd
	rpcRegisterName string
}

func NewRpcCmd(rpcRegisterName string) *RpcCmd {
	return &RpcCmd{NewRootCmd(), rpcRegisterName}
}

func (r *RpcCmd) AddRpc(rpcFn func(client discoveryregistry.SvcDiscoveryRegistry, server *grpc.Server) error) {
	r.Command.RunE = func(cmd *cobra.Command, args []string) error {
		return startrpc.Start(r.getPortFlag(cmd), r.rpcRegisterName, r.getPrometheusPortFlag(cmd), rpcFn)
	}
}

func (r *RpcCmd) Exec(rpcFn func(client discoveryregistry.SvcDiscoveryRegistry, server *grpc.Server) error) error {
	r.AddRpc(rpcFn)
	return r.Execute()
}
