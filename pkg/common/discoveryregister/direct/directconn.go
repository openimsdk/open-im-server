package direct

import (
	"context"
	"errors"
	"fmt"
	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/log"
	config2 "github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/resolver"
	"math/rand"
	"net"
	"strings"
	"time"
)

type ServiceAddresses map[string][]int

func getServiceAddresses() ServiceAddresses {
	return ServiceAddresses{
		config2.Config.RpcRegisterName.OpenImUserName:           config2.Config.RpcPort.OpenImUserPort,
		config2.Config.RpcRegisterName.OpenImFriendName:         config2.Config.RpcPort.OpenImFriendPort,
		config2.Config.RpcRegisterName.OpenImMsgName:            config2.Config.RpcPort.OpenImMessagePort,
		config2.Config.RpcRegisterName.OpenImMessageGatewayName: config2.Config.LongConnSvr.OpenImMessageGatewayPort,
		config2.Config.RpcRegisterName.OpenImGroupName:          config2.Config.RpcPort.OpenImGroupPort,
		config2.Config.RpcRegisterName.OpenImAuthName:           config2.Config.RpcPort.OpenImAuthPort,
		config2.Config.RpcRegisterName.OpenImPushName:           config2.Config.RpcPort.OpenImPushPort,
		config2.Config.RpcRegisterName.OpenImConversationName:   config2.Config.RpcPort.OpenImConversationPort,
		config2.Config.RpcRegisterName.OpenImThirdName:          config2.Config.RpcPort.OpenImThirdPort,
	}
}

// fmt.Sprintf(config2.Config.Rpc.ListenIP+":%d", config2.Config.RpcPort.OpenImUserPort[0])
type ConnManager struct {
	additionalOpts        []grpc.DialOption
	currentServiceAddress string
	conns                 map[string][]*grpc.ClientConn
}

func (cm *ConnManager) GetClientLocalConns() map[string][]*grpc.ClientConn {
	return nil
}

func (cm *ConnManager) GetUserIdHashGatewayHost(ctx context.Context, userId string) (string, error) {
	return "", nil
}

func (cm *ConnManager) Register(serviceName, host string, port int, opts ...grpc.DialOption) error {
	return nil
}

func (cm *ConnManager) UnRegister() error {
	return nil
}

func (cm *ConnManager) CreateRpcRootNodes(serviceNames []string) error {
	return nil
}

func (cm *ConnManager) RegisterConf2Registry(key string, conf []byte) error {
	return nil
}

func (cm *ConnManager) GetConfFromRegistry(key string) ([]byte, error) {
	return nil, nil
}

func (cm *ConnManager) Close() {

}

func NewConnManager() (*ConnManager, error) {
	return &ConnManager{
		conns: make(map[string][]*grpc.ClientConn),
	}, nil
}

func (cm *ConnManager) GetConns(ctx context.Context,
	serviceName string, opts ...grpc.DialOption) ([]*grpc.ClientConn, error) {

	if conns, exists := cm.conns[serviceName]; exists {
		return conns, nil
	}
	ports := getServiceAddresses()[serviceName]
	var connections []*grpc.ClientConn
	var result string
	for _, port := range ports {
		if result != "" {
			result = result + "," + fmt.Sprintf(config2.Config.Rpc.ListenIP+":%d", port)
		} else {
			result = fmt.Sprintf(config2.Config.Rpc.ListenIP+":%d", port)
		}
	}
	conn, err := dialService(ctx, result, append(cm.additionalOpts, opts...)...)
	if err != nil {
		fmt.Errorf("connect to port %s failed,serviceName %s, IP %s", result, serviceName, config2.Config.Rpc.ListenIP)
	}
	connections = append(connections, conn)
	if len(connections) == 0 {
		return nil, fmt.Errorf("no connections found for service: %s", serviceName)
	}
	return connections, nil
}

func (cm *ConnManager) GetConn(ctx context.Context, serviceName string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	// Get service addresses
	addresses := getServiceAddresses()
	address, ok := addresses[serviceName]
	log.ZDebug(ctx, "getConn address", "address", address)
	if !ok {
		return nil, errs.Wrap(errors.New("unknown service name"), "serviceName", serviceName)
	}
	var result string
	for _, addr := range address {
		if result != "" {
			result = result + "," + fmt.Sprintf(config2.Config.Rpc.ListenIP+":%d", addr)
		} else {
			result = fmt.Sprintf(config2.Config.Rpc.ListenIP+":%d", addr)
		}
	}
	// Try to dial a new connection
	conn, err := dialService(ctx, result, append(cm.additionalOpts, opts...)...)
	if err != nil {
		return nil, errs.Wrap(err, "address", result)
	}

	// Store the new connection
	cm.conns[serviceName] = append(cm.conns[serviceName], conn)
	return conn, nil
}

func (cm *ConnManager) GetSelfConnTarget() string {
	return cm.currentServiceAddress
}

func (cm *ConnManager) AddOption(opts ...grpc.DialOption) {
	cm.additionalOpts = append(cm.additionalOpts, opts...)
}

func (cm *ConnManager) CloseConn(conn *grpc.ClientConn) {
	if conn != nil {
		conn.Close()
	}
}

func dialService(ctx context.Context, address string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	options := append(opts, grpc.WithInsecure()) // Replace WithInsecure with proper security options
	conn, err := grpc.DialContext(ctx, "mycustomscheme:///"+address, options...)

	if err != nil {
		return nil, err
	}
	return conn, nil
}

func checkServiceHealth(address string) bool {
	conn, err := net.DialTimeout("tcp", address, time.Second*3)
	if err != nil {
		return false
	}
	conn.Close()
	return true
}

const (
	slashSeparator = "/"
	// EndpointSepChar is the separator char in endpoints.
	EndpointSepChar = ','

	subsetSize = 32
)

// GetEndpoints returns the endpoints from the given target.
func GetEndpoints(target resolver.Target) string {
	return strings.Trim(target.URL.Path, slashSeparator)
}
func subset(set []string, sub int) []string {
	rand.Shuffle(len(set), func(i, j int) {
		set[i], set[j] = set[j], set[i]
	})
	if len(set) <= sub {
		return set
	}

	return set[:sub]
}

type nopResolver struct {
	cc resolver.ClientConn
}

func (n nopResolver) ResolveNow(options resolver.ResolveNowOptions) {

}

func (n nopResolver) Close() {

}

func (cm *ConnManager) Build(target resolver.Target, cc resolver.ClientConn, _ resolver.BuildOptions) (
	resolver.Resolver, error) {

	endpoints := strings.FieldsFunc(GetEndpoints(target), func(r rune) bool {
		return r == EndpointSepChar
	})
	log.ZDebug(context.Background(), "Build", "endpoints", endpoints)
	endpoints = subset(endpoints, subsetSize)
	addrs := make([]resolver.Address, 0, len(endpoints))

	for _, val := range endpoints {
		addrs = append(addrs, resolver.Address{
			Addr: val,
		})
	}
	log.ZDebug(context.Background(), "Build", "addrs", addrs)
	if err := cc.UpdateState(resolver.State{
		Addresses: addrs,
	}); err != nil {
		return nil, err
	}

	return &nopResolver{cc: cc}, nil
}
func init() {
	resolver.Register(&ConnManager{})
}
func (cm *ConnManager) Scheme() string {
	return "mycustomscheme" // return your custom scheme name
}
