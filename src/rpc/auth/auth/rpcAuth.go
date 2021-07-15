package rpcAuth

import (
	"Open_IM/src/common/config"
	log2 "Open_IM/src/common/log"
	"Open_IM/src/grpc-etcdv3/getcdv3"
	pbAuth "Open_IM/src/proto/auth"
	"Open_IM/src/utils"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"strings"
)

type rpcAuth struct {
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string
}

func NewRpcAuthServer(port int) *rpcAuth {
	return &rpcAuth{
		rpcPort:         port,
		rpcRegisterName: config.Config.RpcRegisterName.RpcGetTokenName,
		etcdSchema:      config.Config.Etcd.EtcdSchema,
		etcdAddr:        config.Config.Etcd.EtcdAddr,
	}
}

func (rpc *rpcAuth) Run() {
	log2.Info("", "", "rpc get_token init...")

	address := utils.ServerIP + ":" + strconv.Itoa(rpc.rpcPort)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		log2.Error("", "", "listen network failed, err = %s, address = %s", err.Error(), address)
		return
	}
	log2.Info("", "", "listen network success, address = %s", address)

	//grpc server
	srv := grpc.NewServer()
	defer srv.GracefulStop()

	//service registers with etcd

	pbAuth.RegisterAuthServer(srv, rpc)
	err = getcdv3.RegisterEtcd(rpc.etcdSchema, strings.Join(rpc.etcdAddr, ","), utils.ServerIP, rpc.rpcPort, rpc.rpcRegisterName, 10)
	if err != nil {
		log2.Error("", "", "register rpc get_token to etcd failed, err = %s", err.Error())
		return
	}

	err = srv.Serve(listener)
	if err != nil {
		log2.Info("", "", "rpc get_token fail, err = %s", err.Error())
		return
	}
	log2.Info("", "", "rpc get_token init success")
}
