// Copyright © 2024 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package direct

//import (
//	"context"
//	"fmt"
//
//	config2 "github.com/openimsdk/open-im-server/v3/pkg/common/config"
//	"github.com/openimsdk/tools/errs"
//	"google.golang.org/grpc"
//	"google.golang.org/grpc/credentials/insecure"
//)
//
//type ServiceAddresses map[string][]int
//
//func getServiceAddresses(rpcRegisterName *config2.RpcRegisterName,
//	rpcPort *config2.RpcPort, longConnSvrPort []int) ServiceAddresses {
//	return ServiceAddresses{
//		rpcRegisterName.OpenImUserName:           rpcPort.OpenImUserPort,
//		rpcRegisterName.OpenImFriendName:         rpcPort.OpenImFriendPort,
//		rpcRegisterName.OpenImMsgName:            rpcPort.OpenImMessagePort,
//		rpcRegisterName.OpenImMessageGatewayName: longConnSvrPort,
//		rpcRegisterName.OpenImGroupName:          rpcPort.OpenImGroupPort,
//		rpcRegisterName.OpenImAuthName:           rpcPort.OpenImAuthPort,
//		rpcRegisterName.OpenImPushName:           rpcPort.OpenImPushPort,
//		rpcRegisterName.OpenImConversationName:   rpcPort.OpenImConversationPort,
//		rpcRegisterName.OpenImThirdName:          rpcPort.OpenImThirdPort,
//	}
//}
//
//type ConnDirect struct {
//	additionalOpts        []grpc.DialOption
//	currentServiceAddress string
//	conns                 map[string][]*grpc.ClientConn
//	resolverDirect        *ResolverDirect
//	config                *config2.GlobalConfig
//}
//
//func (cd *ConnDirect) GetClientLocalConns() map[string][]*grpc.ClientConn {
//	return nil
//}
//
//func (cd *ConnDirect) GetUserIdHashGatewayHost(ctx context.Context, userId string) (string, error) {
//	return "", nil
//}
//
//func (cd *ConnDirect) Register(serviceName, host string, port int, opts ...grpc.DialOption) error {
//	return nil
//}
//
//func (cd *ConnDirect) UnRegister() error {
//	return nil
//}
//
//func (cd *ConnDirect) CreateRpcRootNodes(serviceNames []string) error {
//	return nil
//}
//
//func (cd *ConnDirect) RegisterConf2Registry(key string, conf []byte) error {
//	return nil
//}
//
//func (cd *ConnDirect) GetConfFromRegistry(key string) ([]byte, error) {
//	return nil, nil
//}
//
//func (cd *ConnDirect) Close() {
//
//}
//
//func NewConnDirect(config *config2.GlobalConfig) (*ConnDirect, error) {
//	return &ConnDirect{
//		conns:          make(map[string][]*grpc.ClientConn),
//		resolverDirect: NewResolverDirect(),
//		config:         config,
//	}, nil
//}
//
//func (cd *ConnDirect) GetConns(ctx context.Context,
//	serviceName string, opts ...grpc.DialOption) ([]*grpc.ClientConn, error) {
//
//	if conns, exists := cd.conns[serviceName]; exists {
//		return conns, nil
//	}
//	ports := getServiceAddresses(&cd.config.RpcRegisterName,
//		&cd.config.RpcPort, cd.config.LongConnSvr.OpenImMessageGatewayPort)[serviceName]
//	var connections []*grpc.ClientConn
//	for _, port := range ports {
//		conn, err := cd.dialServiceWithoutResolver(ctx, fmt.Sprintf(cd.config.Rpc.ListenIP+":%d", port), append(cd.additionalOpts, opts...)...)
//		if err != nil {
//			return nil, errs.Wrap(fmt.Errorf("connect to port %d failed,serviceName %s, IP %s", port, serviceName, cd.config.Rpc.ListenIP))
//		}
//		connections = append(connections, conn)
//	}
//
//	if len(connections) == 0 {
//		return nil, errs.New("no connections found for service", "serviceName", serviceName).Wrap()
//	}
//	return connections, nil
//}
//
//func (cd *ConnDirect) GetConn(ctx context.Context, serviceName string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
//	// Get service addresses
//	addresses := getServiceAddresses(&cd.config.RpcRegisterName,
//		&cd.config.RpcPort, cd.config.LongConnSvr.OpenImMessageGatewayPort)
//	address, ok := addresses[serviceName]
//	if !ok {
//		return nil, errs.New("unknown service name", "serviceName", serviceName).Wrap()
//	}
//	var result string
//	for _, addr := range address {
//		if result != "" {
//			result = result + "," + fmt.Sprintf(cd.config.Rpc.ListenIP+":%d", addr)
//		} else {
//			result = fmt.Sprintf(cd.config.Rpc.ListenIP+":%d", addr)
//		}
//	}
//	// Try to dial a new connection
//	conn, err := cd.dialService(ctx, result, append(cd.additionalOpts, opts...)...)
//	if err != nil {
//		return nil, errs.WrapMsg(err, "address", result)
//	}
//
//	// Store the new connection
//	cd.conns[serviceName] = append(cd.conns[serviceName], conn)
//	return conn, nil
//}
//
//func (cd *ConnDirect) GetSelfConnTarget() string {
//	return cd.currentServiceAddress
//}
//
//func (cd *ConnDirect) AddOption(opts ...grpc.DialOption) {
//	cd.additionalOpts = append(cd.additionalOpts, opts...)
//}
//
//func (cd *ConnDirect) CloseConn(conn *grpc.ClientConn) {
//	if conn != nil {
//		conn.Close()
//	}
//}
//
//func (cd *ConnDirect) dialService(ctx context.Context, address string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
//	options := append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
//	conn, err := grpc.DialContext(ctx, cd.resolverDirect.Scheme()+":///"+address, options...)
//
//	if err != nil {
//		return nil, errs.WrapMsg(err, "address", address)
//	}
//	return conn, nil
//}
//
//func (cd *ConnDirect) dialServiceWithoutResolver(ctx context.Context, address string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
//	options := append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
//	conn, err := grpc.DialContext(ctx, address, options...)
//
//	if err != nil {
//		return nil, errs.Wrap(err)
//	}
//	return conn, nil
//}
