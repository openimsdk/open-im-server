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

	"github.com/stathat/consistent"
	"google.golang.org/grpc"

	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/log"
)

// K8sDR represents the Kubernetes service discovery and registration client.
type K8sDR struct {
	options               []grpc.DialOption
	rpcRegisterAddr       string
	gatewayHostConsistent *consistent.Consistent
	gatewayName           string
}

func NewK8sDiscoveryRegister(gatewayName string) (discovery.SvcDiscoveryRegistry, error) {
	gatewayConsistent := consistent.New()
	gatewayHosts := getMsgGatewayHost(context.Background(), gatewayName)
	for _, v := range gatewayHosts {
		gatewayConsistent.Add(v)
	}
	return &K8sDR{gatewayHostConsistent: gatewayConsistent}, nil
}

func (cli *K8sDR) Register(serviceName, host string, port int, opts ...grpc.DialOption) error {
	if serviceName != cli.gatewayName {
		cli.rpcRegisterAddr = serviceName
	} else {
		cli.rpcRegisterAddr = getSelfHost(context.Background(), cli.gatewayName)
	}

	return nil
}

func (cli *K8sDR) UnRegister() error {

	return nil
}

func (cli *K8sDR) CreateRpcRootNodes(serviceNames []string) error {

	return nil
}

func (cli *K8sDR) RegisterConf2Registry(key string, conf []byte) error {

	return nil
}

func (cli *K8sDR) GetConfFromRegistry(key string) ([]byte, error) {

	return nil, nil
}

func (cli *K8sDR) GetUserIdHashGatewayHost(ctx context.Context, userId string) (string, error) {
	host, err := cli.gatewayHostConsistent.Get(userId)
	if err != nil {
		log.ZError(ctx, "GetUserIdHashGatewayHost error", err)
	}
	return host, err
}

func getSelfHost(ctx context.Context, gatewayName string) string {
	port := 88
	instance := "openimserver"
	selfPodName := os.Getenv("MY_POD_NAME")
	ns := os.Getenv("MY_POD_NAMESPACE")
	statefuleIndex := 0
	gatewayEnds := strings.Split(gatewayName, ":")
	if len(gatewayEnds) != 2 {
		log.ZError(ctx, "msggateway RpcRegisterName is error:config.RpcRegisterName.OpenImMessageGatewayName", errors.New("config error"))
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

// like openimserver-openim-msggateway-0.openimserver-openim-msggateway-headless.openim-lin.svc.cluster.local:88.
// Replica set in kubernetes environment
func getMsgGatewayHost(ctx context.Context, gatewayName string) []string {
	port := 88
	instance := "openimserver"
	selfPodName := os.Getenv("MY_POD_NAME")
	replicas := os.Getenv("MY_MSGGATEWAY_REPLICACOUNT")
	ns := os.Getenv("MY_POD_NAMESPACE")
	gatewayEnds := strings.Split(gatewayName, ":")
	if len(gatewayEnds) != 2 {
		log.ZError(ctx, "msggateway RpcRegisterName is error:config.RpcRegisterName.OpenImMessageGatewayName", errors.New("config error"))
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
	log.ZDebug(ctx, "getMsgGatewayHost", "instance", instance, "selfPodName", selfPodName, "replicas", replicas, "ns", ns, "ret", ret)
	return ret
}

// GetConns returns the gRPC client connections to the specified service.
func (cli *K8sDR) GetConns(ctx context.Context, serviceName string, opts ...grpc.DialOption) ([]*grpc.ClientConn, error) {

	// This conditional checks if the serviceName is not the OpenImMessageGatewayName.
	// It seems to handle a special case for the OpenImMessageGateway.
	if serviceName != cli.gatewayName {
		// DialContext creates a client connection to the given target (serviceName) using the specified context.
		// 'cli.options' are likely default or common options for all connections in this struct.
		// 'opts...' allows for additional gRPC dial options to be passed and used.
		conn, err := grpc.DialContext(ctx, serviceName, append(cli.options, opts...)...)

		// The function returns a slice of client connections with the new connection, or an error if occurred.
		return []*grpc.ClientConn{conn}, err
	} else {
		// This block is executed if the serviceName is OpenImMessageGatewayName.
		// 'ret' will accumulate the connections to return.
		var ret []*grpc.ClientConn

		// getMsgGatewayHost presumably retrieves hosts for the message gateway service.
		// The context is passed, likely for cancellation and timeout control.
		gatewayHosts := getMsgGatewayHost(ctx, cli.gatewayName)

		// Iterating over the retrieved gateway hosts.
		for _, host := range gatewayHosts {
			// Establishes a connection to each host.
			// Again, appending cli.options with any additional opts provided.
			conn, err := grpc.DialContext(ctx, host, append(cli.options, opts...)...)

			// If there's an error while dialing any host, the function returns immediately with the error.
			if err != nil {
				return nil, err
			} else {
				// If the connection is successful, it is added to the 'ret' slice.
				ret = append(ret, conn)
			}
		}
		// After all hosts are processed, the slice of connections is returned.
		return ret, nil
	}
}

func (cli *K8sDR) GetConn(ctx context.Context, serviceName string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {

	return grpc.DialContext(ctx, serviceName, append(cli.options, opts...)...)
}

func (cli *K8sDR) GetSelfConnTarget() string {

	return cli.rpcRegisterAddr
}

func (cli *K8sDR) AddOption(opts ...grpc.DialOption) {
	cli.options = append(cli.options, opts...)
}

func (cli *K8sDR) CloseConn(conn *grpc.ClientConn) {
	conn.Close()
}

// do not use this method for call rpc.
func (cli *K8sDR) GetClientLocalConns() map[string][]*grpc.ClientConn {
	log.ZError(context.Background(), "should not call this function!", nil)
	return nil
}

func (cli *K8sDR) Close() {

}
