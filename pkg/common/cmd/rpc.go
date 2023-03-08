package cmd

import (
	"OpenIM/pkg/discoveryregistry"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

type RpcCmd struct {
	*RootCmd
}

func NewRpcCmd() *RpcCmd {
	return &RpcCmd{NewRootCmd()}
}

func (r *RpcCmd) AddRpc(f func(port, rpcRegisterName string, prometheusPort int, rpcFn func(client discoveryregistry.SvcDiscoveryRegistry, server *grpc.Server) error, options ...grpc.ServerOption) error) {
	r.Command.RunE = func(cmd *cobra.Command, args []string) error {
		return f(r.port, r.prometheusPort)
	}
}
