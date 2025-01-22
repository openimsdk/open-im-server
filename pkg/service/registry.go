package service

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
)

type Conn interface {
	GetConns(ctx context.Context, serviceName string, opts ...grpc.DialOption) ([]grpc.ClientConnInterface, error) //1
	GetConn(ctx context.Context, serviceName string, opts ...grpc.DialOption) (grpc.ClientConnInterface, error)    //2
	GetSelfConnTarget() string                                                                                     //3
	AddOption(opts ...grpc.DialOption)                                                                             //4
	CloseConn(conn *grpc.ClientConn)                                                                               //5
	// do not use this method for call rpc
}
type SvcDiscoveryRegistry interface {
	Conn
	Register(serviceName, host string, port int, opts ...grpc.DialOption) error //6
	UnRegister() error                                                          //7
	Close()
	GetUserIdHashGatewayHost(ctx context.Context, userId string) (string, error) //
}

var _ SvcDiscoveryRegistry = (*DiscoveryRegistry)(nil)

type DiscoveryRegistry struct {
}

func (x *DiscoveryRegistry) RegisterService(desc *grpc.ServiceDesc, impl any) {
	fmt.Println("RegisterService", desc, impl)
}

func (x *DiscoveryRegistry) GetConns(ctx context.Context, serviceName string, opts ...grpc.DialOption) ([]grpc.ClientConnInterface, error) {
	//TODO implement me
	panic("implement me")
}

func (x *DiscoveryRegistry) GetConn(ctx context.Context, serviceName string, opts ...grpc.DialOption) (grpc.ClientConnInterface, error) {
	//TODO implement me
	panic("implement me")
}

func (x *DiscoveryRegistry) GetSelfConnTarget() string {
	return ""
}

func (x *DiscoveryRegistry) AddOption(opts ...grpc.DialOption) {}

func (x *DiscoveryRegistry) CloseConn(conn *grpc.ClientConn) {
	_ = conn.Close()
}

func (x *DiscoveryRegistry) Register(serviceName, host string, port int, opts ...grpc.DialOption) error {
	return nil
}

func (x *DiscoveryRegistry) UnRegister() error {
	return nil
}

func (x *DiscoveryRegistry) Close() {}

func (x *DiscoveryRegistry) GetUserIdHashGatewayHost(ctx context.Context, userId string) (string, error) {
	return "", nil
}
