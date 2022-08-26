package admin_cms

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	openIMHttp "Open_IM/pkg/common/http"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbAdminCMS "Open_IM/pkg/proto/admin_cms"
	server_api_params "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
	"net"
	"strconv"
	"strings"

	"google.golang.org/grpc"
)

type adminCMSServer struct {
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string
}

func NewAdminCMSServer(port int) *adminCMSServer {
	log.NewPrivateLog(constant.LogFileName)
	return &adminCMSServer{
		rpcPort:         port,
		rpcRegisterName: config.Config.RpcRegisterName.OpenImAdminCMSName,
		etcdSchema:      config.Config.Etcd.EtcdSchema,
		etcdAddr:        config.Config.Etcd.EtcdAddr,
	}
}

func (s *adminCMSServer) Run() {
	log.NewInfo("0", "AdminCMS rpc start ")
	listenIP := ""
	if config.Config.ListenIP == "" {
		listenIP = "0.0.0.0"
	} else {
		listenIP = config.Config.ListenIP
	}
	address := listenIP + ":" + strconv.Itoa(s.rpcPort)

	//listener network
	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic("listening err:" + err.Error() + s.rpcRegisterName)
	}
	log.NewInfo("0", "listen network success, ", address, listener)
	defer listener.Close()
	//grpc server
	srv := grpc.NewServer()
	defer srv.GracefulStop()
	//Service registers with etcd
	pbAdminCMS.RegisterAdminCMSServer(srv, s)
	rpcRegisterIP := config.Config.RpcRegisterIP
	if config.Config.RpcRegisterIP == "" {
		rpcRegisterIP, err = utils.GetLocalIP()
		if err != nil {
			log.Error("", "GetLocalIP failed ", err.Error())
		}
	}
	log.NewInfo("", "rpcRegisterIP", rpcRegisterIP)
	err = getcdv3.RegisterEtcd(s.etcdSchema, strings.Join(s.etcdAddr, ","), rpcRegisterIP, s.rpcPort, s.rpcRegisterName, 10)
	if err != nil {
		log.NewError("0", "RegisterEtcd failed ", err.Error())
		panic(utils.Wrap(err, "register admin module  rpc to etcd err"))
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
	for i, adminID := range config.Config.Manager.AppManagerUid {
		if adminID == req.AdminID && config.Config.Manager.Secrets[i] == req.Secret {
			token, expTime, err := token_verify.CreateToken(adminID, constant.SingleChatType)
			log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "generate token success", "token: ", token, "expTime:", expTime)
			if err != nil {
				log.NewError(req.OperationID, utils.GetSelfFuncName(), "generate token failed", "adminID: ", adminID, err.Error())
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
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp.String())
	return resp, nil
}

func (s *adminCMSServer) AddUserRegisterAddFriendIDList(_ context.Context, req *pbAdminCMS.AddUserRegisterAddFriendIDListReq) (*pbAdminCMS.AddUserRegisterAddFriendIDListResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp := &pbAdminCMS.AddUserRegisterAddFriendIDListResp{}
	if err := imdb.AddUserRegisterAddFriendIDList(req.UserIDList...); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), req.UserIDList)
		return resp, openIMHttp.WrapError(constant.ErrDB)
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", req.String())
	return resp, nil
}

func (s *adminCMSServer) ReduceUserRegisterAddFriendIDList(_ context.Context, req *pbAdminCMS.ReduceUserRegisterAddFriendIDListReq) (*pbAdminCMS.ReduceUserRegisterAddFriendIDListResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp := &pbAdminCMS.ReduceUserRegisterAddFriendIDListResp{}
	if req.Operation == 0 {
		if err := imdb.ReduceUserRegisterAddFriendIDList(req.UserIDList...); err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), req.UserIDList)
			return resp, openIMHttp.WrapError(constant.ErrDB)
		}
	} else {
		if err := imdb.DeleteAllRegisterAddFriendIDList(); err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), req.UserIDList)
			return resp, openIMHttp.WrapError(constant.ErrDB)
		}
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", req.String())
	return resp, nil
}

func (s *adminCMSServer) GetUserRegisterAddFriendIDList(_ context.Context, req *pbAdminCMS.GetUserRegisterAddFriendIDListReq) (*pbAdminCMS.GetUserRegisterAddFriendIDListResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp := &pbAdminCMS.GetUserRegisterAddFriendIDListResp{UserInfoList: []*server_api_params.UserInfo{}}
	userIDList, err := imdb.GetRegisterAddFriendList(req.Pagination.ShowNumber, req.Pagination.PageNumber)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error())
		return resp, openIMHttp.WrapError(constant.ErrDB)
	}
	userList, err := imdb.GetUsersByUserIDList(userIDList)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), userIDList)
		return resp, openIMHttp.WrapError(constant.ErrDB)
	}
	log.NewDebug(req.OperationID, utils.GetSelfFuncName(), userList, userIDList)
	resp.Pagination = &server_api_params.ResponsePagination{
		CurrentPage: req.Pagination.PageNumber,
		ShowNumber:  req.Pagination.ShowNumber,
	}
	utils.CopyStructFields(&resp.UserInfoList, userList)
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", req.String())
	return resp, nil
}
