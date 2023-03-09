package cmd

import (
	"OpenIM/internal/startrpc"
	"OpenIM/pkg/discoveryregistry"
	"github.com/spf13/cobra"
	"google.golang.org/grpc"
)

type RpcCmd struct {
	*RootCmd
}

func NewRpcCmd() *RpcCmd {
	rpcCmd := &RpcCmd{NewRootCmd()}
	return rpcCmd
}

func (r *RpcCmd) addRpc(rpcRegisterName string, rpcFn func(client discoveryregistry.SvcDiscoveryRegistry, server *grpc.Server) error) {
	r.Command.RunE = func(cmd *cobra.Command, args []string) error {
		return startrpc.Start(r.getPortFlag(cmd), rpcRegisterName, r.getPrometheusPortFlag(cmd), rpcFn)
	}
}

func (r *RpcCmd) Exec(rpcRegisterName string, rpcFn func(client discoveryregistry.SvcDiscoveryRegistry, server *grpc.Server) error) error {
	r.addRpc(rpcRegisterName, rpcFn)
	return r.Execute()
}

func (r *RpcCmd) addRpc2(rpcRegisterName *string, rpcFn func(client discoveryregistry.SvcDiscoveryRegistry, server *grpc.Server) error) {
	r.Command.RunE = func(cmd *cobra.Command, args []string) error {
		return startrpc.Start(r.getPortFlag(cmd), *rpcRegisterName, r.getPrometheusPortFlag(cmd), rpcFn)
	}
}

func (r *RpcCmd) Exec2(rpcRegisterName *string, rpcFn func(client discoveryregistry.SvcDiscoveryRegistry, server *grpc.Server) error) error {
	r.addRpc2(rpcRegisterName, rpcFn)
	return r.Execute()
}
