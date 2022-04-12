package organization

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	rpc "Open_IM/pkg/proto/organization"
	"Open_IM/pkg/utils"

	"context"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"strings"
)

type organizationServer struct {
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string
}

func NewGroupServer(port int) *organizationServer {
	log.NewPrivateLog(constant.LogFileName)
	return &organizationServer{
		rpcPort:         port,
		rpcRegisterName: config.Config.RpcRegisterName.OpenImGroupName,
		etcdSchema:      config.Config.Etcd.EtcdSchema,
		etcdAddr:        config.Config.Etcd.EtcdAddr,
	}
}

func (s *organizationServer) Run() {
	log.NewInfo("", "organization rpc start ")
	ip := utils.ServerIP
	registerAddress := ip + ":" + strconv.Itoa(s.rpcPort)
	//listener network
	listener, err := net.Listen("tcp", registerAddress)
	if err != nil {
		log.NewError("", "Listen failed ", err.Error(), registerAddress)
		return
	}
	log.NewInfo("", "listen network success, ", registerAddress, listener)
	defer listener.Close()
	//grpc server
	srv := grpc.NewServer()
	defer srv.GracefulStop()
	//Service registers with etcd
	rpc.RegisterOrganizationServer(srv, s)
	err = getcdv3.RegisterEtcd(s.etcdSchema, strings.Join(s.etcdAddr, ","), ip, s.rpcPort, s.rpcRegisterName, 10)
	if err != nil {
		log.NewError("", "RegisterEtcd failed ", err.Error())
		return
	}
	log.NewInfo("", "organization rpc RegisterEtcd success", ip, s.rpcPort, s.rpcRegisterName, 10)
	err = srv.Serve(listener)
	if err != nil {
		log.NewError("", "Serve failed ", err.Error())
		return
	}
	log.NewInfo("", "organization rpc success")
}

func (s *organizationServer) CreateDepartment(ctx context.Context, req *rpc.CreateDepartmentReq) (*rpc.CreateDepartmentResp, error) {
	return nil, nil
}
