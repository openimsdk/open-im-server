package discoveryregistry

import (
	"google.golang.org/grpc"
)

type SvcDiscoveryRegistry interface {
	Register(serviceName, host string, port int, opts ...grpc.DialOption) error
	UnRegister() error
	GetConns(serviceName string, opts ...grpc.DialOption) ([]*grpc.ClientConn, error)
	GetConn(serviceName string, opts ...grpc.DialOption) (*grpc.ClientConn, error)
	Re
}
