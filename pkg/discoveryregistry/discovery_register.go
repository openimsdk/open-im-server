package discoveryregistry

import (
	"google.golang.org/grpc"
)

type Conn interface {
	GetConns(serviceName string, opts ...grpc.DialOption) ([]*grpc.ClientConn, error)
	GetConn(serviceName string, opts ...grpc.DialOption) (*grpc.ClientConn, error)
	AddOption(opts ...grpc.DialOption)
}

type SvcDiscoveryRegistry interface {
	Conn
	Register(serviceName, host string, port int, opts ...grpc.DialOption) error
	UnRegister() error
	RegisterConf2Registry(key string, conf []byte) error
	GetConfFromRegistry(key string) ([]byte, error)
}
