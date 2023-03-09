package cmd

import (
	"OpenIM/internal/startrpc"
	"OpenIM/pkg/discoveryregistry"
	"google.golang.org/grpc"
)

type PushCmd struct {
	*RpcCmd
}

func NewPushCmd() *PushCmd {
	return &PushCmd{NewRpcCmd()}
}

func (r *PushCmd) StartSvr(name string, rpcFn func(client discoveryregistry.SvcDiscoveryRegistry, server *grpc.Server) error) error {
	return startrpc.Start(r.GetPortFlag(), name, r.GetPrometheusPortFlag(), rpcFn)
}
