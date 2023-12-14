// Copyright © 2023 OpenIM. All rights reserved.
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

package kubernetes

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"github.com/OpenIMSDK/tools/discoveryregistry"
)

// K8sDR represents the Kubernetes service discovery and registration client.
type K8sDR struct {
	options         []grpc.DialOption
	rpcRegisterAddr string
}

// NewK8sDiscoveryRegister creates a new instance of K8sDR for Kubernetes service discovery and registration.
func NewK8sDiscoveryRegister() (discoveryregistry.SvcDiscoveryRegistry, error) {

	return &K8sDR{}, nil
}

// Register registers a service with Kubernetes.
func (cli *K8sDR) Register(serviceName, host string, port int, opts ...grpc.DialOption) error {
	cli.rpcRegisterAddr = serviceName

	return nil
}

// UnRegister removes a service registration from Kubernetes.
func (cli *K8sDR) UnRegister() error {

	return nil
}

// CreateRpcRootNodes creates root nodes for RPC in Kubernetes.
func (cli *K8sDR) CreateRpcRootNodes(serviceNames []string) error {

	return nil
}

// RegisterConf2Registry registers a configuration to the registry.
func (cli *K8sDR) RegisterConf2Registry(key string, conf []byte) error {

	return nil
}

// GetConfFromRegistry retrieves a configuration from the registry.
func (cli *K8sDR) GetConfFromRegistry(key string) ([]byte, error) {
	return nil, nil
}

// GetConns returns a list of gRPC client connections for a given service.
func (cli *K8sDR) GetConns(ctx context.Context, serviceName string, opts ...grpc.DialOption) ([]*grpc.ClientConn, error) {
	conn, err := grpc.DialContext(ctx, serviceName, append(cli.options, opts...)...)
	return []*grpc.ClientConn{conn}, err
}

// GetConn returns a single gRPC client connection for a given service.
func (cli *K8sDR) GetConn(ctx context.Context, serviceName string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	return grpc.DialContext(ctx, serviceName, append(cli.options, opts...)...)
}

// GetSelfConnTarget returns the connection target of the client itself.
func (cli *K8sDR) GetSelfConnTarget() string {
	return cli.rpcRegisterAddr
}

// AddOption adds gRPC dial options to the client.
func (cli *K8sDR) AddOption(opts ...grpc.DialOption) {
	cli.options = append(cli.options, opts...)
}

// CloseConn closes a given gRPC client connection.
func (cli *K8sDR) CloseConn(conn *grpc.ClientConn) {
	conn.Close()
}

// do not use this method for call rpc.
func (cli *K8sDR) GetClientLocalConns() map[string][]*grpc.ClientConn {
	fmt.Println("should not call this function!!!!!!!!!!!!!!!!!!!!!!!!!")

	return nil
}

// Close closes the K8sDR client.
func (cli *K8sDR) Close() {

	// Close any open resources here (if applicable)
	return
}
