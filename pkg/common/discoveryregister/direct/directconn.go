package direct

import (
	"context"
	"errors"
	"fmt"
	"github.com/OpenIMSDK/tools/errs"
	config2 "github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"google.golang.org/grpc"
	"net"
	"time"
)

type ServiceAddresses map[string]string

func getServiceAddresses() ServiceAddresses {
	fmt.Println("!!!!!!!!!!!!", fmt.Sprintf(config2.Config.RpcRegisterName.OpenImUserName, config2.Config.RpcPort.OpenImUserPort[0]))
	return ServiceAddresses{
		"User":           fmt.Sprintf(config2.Config.RpcRegisterName.OpenImUserName, config2.Config.RpcPort.OpenImUserPort[0]),
		"Friend":         fmt.Sprintf(config2.Config.RpcRegisterName.OpenImFriendName, config2.Config.RpcPort.OpenImFriendPort[0]),
		"Message":        fmt.Sprintf(config2.Config.RpcRegisterName.OpenImMsgName, config2.Config.RpcPort.OpenImMessagePort[0]),
		"MessageGateway": fmt.Sprintf(config2.Config.RpcRegisterName.OpenImMessageGatewayName, config2.Config.LongConnSvr.OpenImMessageGatewayPort[0]),
		"Group":          fmt.Sprintf(config2.Config.RpcRegisterName.OpenImGroupName, config2.Config.RpcPort.OpenImGroupPort[0]),
		"Auth":           fmt.Sprintf(config2.Config.RpcRegisterName.OpenImAuthName, config2.Config.RpcPort.OpenImAuthPort[0]),
		"Push":           fmt.Sprintf(config2.Config.RpcRegisterName.OpenImPushName, config2.Config.RpcPort.OpenImPushPort[0]),
		"Conversation":   fmt.Sprintf(config2.Config.RpcRegisterName.OpenImConversationName, config2.Config.RpcPort.OpenImConversationPort[0]),
		"Third":          fmt.Sprintf(config2.Config.RpcRegisterName.OpenImThirdName, config2.Config.RpcPort.OpenImThirdPort[0]),
	}
}

type ConnManager struct {
	additionalOpts        []grpc.DialOption
	currentServiceAddress string
	conns                 map[string]*grpc.ClientConn
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
		conns: make(map[string]*grpc.ClientConn),
	}, nil
}

func (cm *ConnManager) GetConns(ctx context.Context,
	serviceName string, opts ...grpc.DialOption) ([]*grpc.ClientConn, error) {
	var connections []*grpc.ClientConn
	for name, conn := range cm.conns {
		if name == serviceName {
			connections = append(connections, conn)
		}
	}
	if len(connections) == 0 {
		return nil, fmt.Errorf("no connections found for service: %s", serviceName)
	}
	return connections, nil
}

func (cm *ConnManager) GetConn(ctx context.Context, serviceName string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	if conn, exists := cm.conns[serviceName]; exists {
		return conn, nil
	}
	addresses := getServiceAddresses()
	address, ok := addresses[serviceName]
	if !ok {
		return nil, errs.Wrap(errors.New("unknown service name"), "serviceName", serviceName)
	}

	conn, err := dialService(address, opts...)
	if err != nil {
		return nil, err
	}
	cm.conns[serviceName] = conn
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

func dialService(address string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	options := append(opts, grpc.WithInsecure()) // Replace WithInsecure with proper security options
	conn, err := grpc.DialContext(context.Background(), address, options...)
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
