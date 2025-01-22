package service

//
//import (
//	"context"
//	"fmt"
//	"sync"
//
//	"google.golang.org/grpc"
//)
//
//type DiscoveryRegistry struct {
//	lock     sync.RWMutex
//	services map[string]grpc.ClientConnInterface
//}
//
//func (x *DiscoveryRegistry) RegisterService(desc *grpc.ServiceDesc, impl any) {
//	fmt.Println("RegisterService", desc, impl)
//}
//
//func (x *DiscoveryRegistry) GetConns(ctx context.Context, serviceName string, opts ...grpc.DialOption) ([]grpc.ClientConnInterface, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (x *DiscoveryRegistry) GetConn(ctx context.Context, serviceName string, opts ...grpc.DialOption) (grpc.ClientConnInterface, error) {
//	//TODO implement me
//	panic("implement me")
//}
//
//func (x *DiscoveryRegistry) IsSelfNode(cc grpc.ClientConnInterface) bool {
//
//	return false
//}
