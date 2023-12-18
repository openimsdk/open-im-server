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
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"google.golang.org/grpc"

	"github.com/OpenIMSDK/tools/discoveryregistry"
	"github.com/OpenIMSDK/tools/log"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
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
	if serviceName != config.Config.RpcRegisterName.OpenImMessageGatewayName {
		cli.rpcRegisterAddr = serviceName
	} else {
		cli.rpcRegisterAddr = cli.getSelfHost(context.Background())
	}

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

func (cli *K8sDR) getSelfHost(ctx context.Context) string {
	port := 88
	instance := "openimserver"
	selfPodName := os.Getenv("MY_POD_NAME")
	ns := os.Getenv("MY_POD_NAMESPACE")
	statefuleIndex := 0
	gatewayEnds := strings.Split(config.Config.RpcRegisterName.OpenImMessageGatewayName, ":")
	if len(gatewayEnds) != 2 {
		log.ZError(ctx, "msggateway RpcRegisterName is error:config.Config.RpcRegisterName.OpenImMessageGatewayName", errors.New("config error"))
	} else {
		port, _ = strconv.Atoi(gatewayEnds[1])
	}
	podInfo := strings.Split(selfPodName, "-")
	instance = podInfo[0]
	count := len(podInfo)
	statefuleIndex, _ = strconv.Atoi(podInfo[count-1])
	host := fmt.Sprintf("%s-openim-msggateway-%d.%s-openim-msggateway-headless.%s.svc.cluster.local:%d", instance, statefuleIndex, instance, ns, port)
	return host
}

// GetConns returns a list of gRPC client connections for a given service.
func (cli *K8sDR) GetConns(ctx context.Context, serviceName string, opts ...grpc.DialOption) ([]*grpc.ClientConn, error) {
	if serviceName != config.Config.RpcRegisterName.OpenImMessageGatewayName {
		conn, err := grpc.DialContext(ctx, serviceName, append(cli.options, opts...)...)
		return []*grpc.ClientConn{conn}, err
	}
	var ret []*grpc.ClientConn
	gatewayHosts := cli.getMsgGatewayHost(ctx)
	for _, host := range gatewayHosts {
		conn, err := grpc.DialContext(ctx, host, append(cli.options, opts...)...)
		if err != nil {
			return nil, err
		}
		ret = append(ret, conn)
	}
	return ret, nil
}

// like openimserver-openim-msggateway-0.openimserver-openim-msggateway-headless.openim-lin.svc.cluster.local:88
func (cli *K8sDR) getMsgGatewayHost(ctx context.Context) []string {
	port := 88
	instance := "openimserver"
	selfPodName := os.Getenv("MY_POD_NAME")
	replicas := os.Getenv("MY_MSGGATEWAY_REPLICACOUNT")
	ns := os.Getenv("MY_POD_NAMESPACE")
	gatewayEnds := strings.Split(config.Config.RpcRegisterName.OpenImMessageGatewayName, ":")
	if len(gatewayEnds) != 2 {
		log.ZError(ctx, "msggateway RpcRegisterName is error:config.Config.RpcRegisterName.OpenImMessageGatewayName", errors.New("config error"))
	} else {
		port, _ = strconv.Atoi(gatewayEnds[1])
	}
	nReplicas, _ := strconv.Atoi(replicas)
	podInfo := strings.Split(selfPodName, "-")
	instance = podInfo[0]
	var ret []string
	for i := 0; i < nReplicas; i++ {
		host := fmt.Sprintf("%s-openim-msggateway-%d.%s-openim-msggateway-headless.%s.svc.cluster.local:%d", instance, i, instance, ns, port)
		ret = append(ret, host)
	}
	log.ZInfo(ctx, "getMsgGatewayHost", "instance", instance, "selfPodName", selfPodName, "replicas", replicas, "ns", ns, "ret", ret)

	return ret
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
