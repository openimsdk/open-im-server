package rpc_server

import (
	"Open_IM/internal/common/network"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	discoveryRegistry "Open_IM/pkg/discoveryregistry"
	"github.com/OpenIMSDK/openKeeper"
	"net"
	"strconv"
)

type RpcServer struct {
	Port           int
	RegisterName   string
	RegisterCenter discoveryRegistry.SvcDiscoveryRegistry
}

func NewRpcServer(registerIPInConfig string, port int, registerName string, zkServers []string, zkRoot string) (*RpcServer, error) {
	log.NewPrivateLog(constant.LogFileName)
	s := &RpcServer{
		Port:         port,
		RegisterName: registerName,
	}

	zkClient, err := openKeeper.NewClient(zkServers, zkRoot, 10, "", "")
	if err != nil {
		return nil, err
	}
	registerIP, err := network.GetRpcRegisterIP(registerIPInConfig)
	if err != nil {
		return nil, err
	}
	err = zkClient.Register(s.RegisterName, registerIP, s.Port)
	if err != nil {
		return nil, err
	}
	s.RegisterCenter = zkClient
	return s, nil
}

func GetTcpListen(listenIPInConfig string, port int) (net.Listener, string, error) {
	address := network.GetListenIP(listenIPInConfig) + ":" + strconv.Itoa(port)
	listener, err := net.Listen("tcp", address)
	return listener, address, err
}
