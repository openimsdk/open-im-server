package admin_cms

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	openIMHttp "Open_IM/pkg/common/http"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbAdminCMS "Open_IM/pkg/proto/admin_cms"
	"Open_IM/pkg/utils"
	"context"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"strings"
)

type adminCMSServer struct {
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string
}

func NewAdminCMSServer(port int) *adminCMSServer {
	log.NewPrivateLog("AdminCMS")
	return &adminCMSServer{
		rpcPort:         port,
		rpcRegisterName: config.Config.RpcRegisterName.OpenImAdminCMSName,
		etcdSchema:      config.Config.Etcd.EtcdSchema,
		etcdAddr:        config.Config.Etcd.EtcdAddr,
	}
}

func (s *adminCMSServer) Run() {
	log.NewInfo("0", "AdminCMS rpc start ")
	ip := utils.ServerIP
	registerAddress := ip + ":" + strconv.Itoa(s.rpcPort)
	//listener network
	listener, err := net.Listen("tcp", registerAddress)
	if err != nil {
		log.NewError("0", "Listen failed ", err.Error(), registerAddress)
		return
	}
	log.NewInfo("0", "listen network success, ", registerAddress, listener)
	defer listener.Close()
	//grpc server
	srv := grpc.NewServer()
	defer srv.GracefulStop()
	//Service registers with etcd
	pbAdminCMS.RegisterAdminCMSServer(srv, s)
	err = getcdv3.RegisterEtcd(s.etcdSchema, strings.Join(s.etcdAddr, ","), ip, s.rpcPort, s.rpcRegisterName, 10)
	if err != nil {
		log.NewError("0", "RegisterEtcd failed ", err.Error())
		return
	}
	err = srv.Serve(listener)
	if err != nil {
		log.NewError("0", "Serve failed ", err.Error())
		return
	}
	log.NewInfo("0", "message cms rpc success")
}

func (s *adminCMSServer) AdminLogin(_ context.Context, req *pbAdminCMS.AdminLoginReq) (*pbAdminCMS.AdminLoginResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp := &pbAdminCMS.AdminLoginResp{}
	for i, adminID := range config.Config.Manager.AppManagerUid{
		if adminID == req.AdminID && config.Config.Manager.Secrets[i] == req.Secret {
			token, expTime, err := token_verify.CreateToken(adminID, constant.SingleChatType)
			log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "generate token success", "token: ", token, "expTime:", expTime)
			if err != nil {
				log.NewError(req.OperationID, utils.GetSelfFuncName(), "generate token failed", "adminID: ", adminID,  err.Error())
				return resp, openIMHttp.WrapError(constant.ErrTokenUnknown)
			}
			resp.Token = token
			break
		}
	}

	if resp.Token == "" {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "failed")
		return resp, openIMHttp.WrapError(constant.ErrTokenMalformed)
	}
	return resp, nil
}