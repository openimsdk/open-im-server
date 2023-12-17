package discoveryregister

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
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
	if serviceName != config.Config.RpcRegisterName.OpenImMessageGatewayName {
		cli.rpcRegisterAddr = serviceName
	} else {
		cli.rpcRegisterAddr = cli.getSelfHost(context.Background())
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
func (cli *K8sDR) GetConns(ctx context.Context, serviceName string, opts ...grpc.DialOption) ([]*grpc.ClientConn, error) {

	if serviceName != config.Config.RpcRegisterName.OpenImMessageGatewayName {
		conn, err := grpc.DialContext(ctx, serviceName, append(cli.options, opts...)...)
		return []*grpc.ClientConn{conn}, err
	} else {
		var ret []*grpc.ClientConn
		gatewayHosts := cli.getMsgGatewayHost(ctx)
		for _, host := range gatewayHosts {
			conn, err := grpc.DialContext(ctx, host, append(cli.options, opts...)...)
			if err != nil {
				return nil, err
			} else {
				ret = append(ret, conn)
			}
		}
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

// do not use this method for call rpc
func (cli *K8sDR) GetClientLocalConns() map[string][]*grpc.ClientConn {
	fmt.Println("should not call this function!!!!!!!!!!!!!!!!!!!!!!!!!")
	return nil
}
func (cli *K8sDR) Close() {
	return
}
