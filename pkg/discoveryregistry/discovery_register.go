package discoveryregistry

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
)

type Conn interface {
	GetConns(ctx context.Context, serviceName string, opts ...grpc.DialOption) ([]*grpc.ClientConn, error)
	GetConn(ctx context.Context, serviceName string, opts ...grpc.DialOption) (*grpc.ClientConn, error)
	AddOption(opts ...grpc.DialOption)
	// do not use this method for call rpc
	GetClientLocalConns() map[string][]resolver.Address
}

type SvcDiscoveryRegistry interface {
	Conn
	Register(serviceName, host string, port int, opts ...grpc.DialOption) error
	UnRegister() error
	RegisterConf2Registry(key string, conf []byte) error
	GetConfFromRegistry(key string) ([]byte, error)
}
