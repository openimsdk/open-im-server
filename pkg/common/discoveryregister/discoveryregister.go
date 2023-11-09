package discoveryregister

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/OpenIMSDK/tools/discoveryregistry"
	openkeeper "github.com/OpenIMSDK/tools/discoveryregistry/zookeeper"
	"github.com/OpenIMSDK/tools/log"
	"google.golang.org/grpc"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
)

func NewDiscoveryRegister(envType string) (discoveryregistry.SvcDiscoveryRegistry, error) {
	var client discoveryregistry.SvcDiscoveryRegistry
	var err error
	switch envType {
	case "zookeeper":
		client, err = openkeeper.NewClient(config.Config.Zookeeper.ZkAddr, config.Config.Zookeeper.Schema,
			openkeeper.WithFreq(time.Hour), openkeeper.WithUserNameAndPassword(
				config.Config.Zookeeper.Username,
				config.Config.Zookeeper.Password,
			), openkeeper.WithRoundRobin(), openkeeper.WithTimeout(10), openkeeper.WithLogger(log.NewZkLogger()))
	case "k8s":
		client, err = NewK8sDiscoveryRegister()
	default:
		client = nil
		err = errors.New("envType not correct")
	}
	return client, err
}

type K8sDR struct {
	options         []grpc.DialOption
	rpcRegisterAddr string
}

func NewK8sDiscoveryRegister() (discoveryregistry.SvcDiscoveryRegistry, error) {
	return &K8sDR{}, nil
}

func (cli *K8sDR) Register(serviceName, host string, port int, opts ...grpc.DialOption) error {
	cli.rpcRegisterAddr = serviceName
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
func (cli *K8sDR) GetConns(ctx context.Context, serviceName string, opts ...grpc.DialOption) ([]*grpc.ClientConn, error) {

	conn, err := grpc.DialContext(ctx, serviceName, append(cli.options, opts...)...)
	return []*grpc.ClientConn{conn}, err
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

// do not use this method for call rpc
func (cli *K8sDR) GetClientLocalConns() map[string][]*grpc.ClientConn {
	fmt.Println("should not call this function!!!!!!!!!!!!!!!!!!!!!!!!!")
	return nil
}
func (cli *K8sDR) Close() {
	return
}
