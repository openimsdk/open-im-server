package auth

import (
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbAuth "Open_IM/pkg/proto/auth"
	pbRelay "Open_IM/pkg/proto/relay"
	"Open_IM/pkg/utils"
	"context"
	"errors"
	"net"
	"strconv"
	"strings"

	"Open_IM/pkg/common/config"

	"google.golang.org/grpc"
)

func (rpc *rpcAuth) UserRegister(_ context.Context, req *pbAuth.UserRegisterReq) (*pbAuth.UserRegisterResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " rpc args ", req.String())
	var user db.User
	utils.CopyStructFields(&user, req.UserInfo)
	if req.UserInfo.Birth != 0 {
		user.Birth = utils.UnixSecondToTime(int64(req.UserInfo.Birth))
	}
	log.Debug(req.OperationID, "copy ", user, req.UserInfo)
	err := imdb.UserRegister(user)
	if err != nil {
		errMsg := req.OperationID + " imdb.UserRegister failed " + err.Error() + user.UserID
		log.NewError(req.OperationID, errMsg, user)
		return &pbAuth.UserRegisterResp{CommonResp: &pbAuth.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}}, nil
	}

	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " rpc return ", pbAuth.UserRegisterResp{CommonResp: &pbAuth.CommonResp{}})
	return &pbAuth.UserRegisterResp{CommonResp: &pbAuth.CommonResp{}}, nil
}

func (rpc *rpcAuth) UserToken(_ context.Context, req *pbAuth.UserTokenReq) (*pbAuth.UserTokenResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " rpc args ", req.String())
	_, err := imdb.GetUserByUserID(req.FromUserID)
	if err != nil {
		errMsg := req.OperationID + " imdb.GetUserByUserID failed " + err.Error() + req.FromUserID
		log.NewError(req.OperationID, errMsg)
		return &pbAuth.UserTokenResp{CommonResp: &pbAuth.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}}, nil
	}

	tokens, expTime, err := token_verify.CreateToken(req.FromUserID, int(req.Platform))
	if err != nil {
		errMsg := req.OperationID + " token_verify.CreateToken failed " + err.Error() + req.FromUserID + utils.Int32ToString(req.Platform)
		log.NewError(req.OperationID, errMsg)
		return &pbAuth.UserTokenResp{CommonResp: &pbAuth.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}}, nil
	}

	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " rpc return ", pbAuth.UserTokenResp{CommonResp: &pbAuth.CommonResp{}, Token: tokens, ExpiredTime: expTime})
	return &pbAuth.UserTokenResp{CommonResp: &pbAuth.CommonResp{}, Token: tokens, ExpiredTime: expTime}, nil
}

func (rpc *rpcAuth) ForceLogout(_ context.Context, req *pbAuth.ForceLogoutReq) (*pbAuth.ForceLogoutResp, error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " rpc args ", req.String())
	if !token_verify.IsManagerUserID(req.OpUserID) {
		errMsg := req.OperationID + " IsManagerUserID false " + req.OpUserID
		log.NewError(req.OperationID, errMsg)
		return &pbAuth.ForceLogoutResp{CommonResp: &pbAuth.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: errMsg}}, nil
	}
	if err := token_verify.DeleteToken(req.FromUserID, int(req.Platform)); err != nil {
		errMsg := req.OperationID + " DeleteToken failed " + err.Error() + req.FromUserID + utils.Int32ToString(req.Platform)
		log.NewError(req.OperationID, errMsg)
		return &pbAuth.ForceLogoutResp{CommonResp: &pbAuth.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}}, nil
	}
	if err := rpc.forceKickOff(req.FromUserID, req.Platform, req.OperationID); err != nil {
		errMsg := req.OperationID + " forceKickOff failed " + err.Error() + req.FromUserID + utils.Int32ToString(req.Platform)
		log.NewError(req.OperationID, errMsg)
		return &pbAuth.ForceLogoutResp{CommonResp: &pbAuth.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: errMsg}}, nil
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), " rpc return ", pbAuth.UserTokenResp{CommonResp: &pbAuth.CommonResp{}})
	return &pbAuth.ForceLogoutResp{CommonResp: &pbAuth.CommonResp{}}, nil
}

func (rpc *rpcAuth) forceKickOff(userID string, platformID int32, operationID string) error {

	grpcCons := getcdv3.GetConn4Unique(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImRelayName)
	for _, v := range grpcCons {
		client := pbRelay.NewRelayClient(v)
		kickReq := &pbRelay.KickUserOfflineReq{OperationID: operationID, KickUserIDList: []string{userID}, PlatformID: platformID}
		log.NewInfo(operationID, "KickUserOffline ", client, kickReq.String())
		_, err := client.KickUserOffline(context.Background(), kickReq)
		return utils.Wrap(err, "")
	}

	return errors.New("no rpc node ")
}

type rpcAuth struct {
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string
}

func NewRpcAuthServer(port int) *rpcAuth {
	log.NewPrivateLog(constant.LogFileName)
	return &rpcAuth{
		rpcPort:         port,
		rpcRegisterName: config.Config.RpcRegisterName.OpenImAuthName,
		etcdSchema:      config.Config.Etcd.EtcdSchema,
		etcdAddr:        config.Config.Etcd.EtcdAddr,
	}
}

func (rpc *rpcAuth) Run() {
	operationID := utils.OperationIDGenerator()
	log.NewInfo(operationID, "rpc auth start...")

	listenIP := ""
	if config.Config.ListenIP == "" {
		listenIP = "0.0.0.0"
	} else {
		listenIP = config.Config.ListenIP
	}
	address := listenIP + ":" + strconv.Itoa(rpc.rpcPort)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic("listening err:" + err.Error() + rpc.rpcRegisterName)
	}
	log.NewInfo(operationID, "listen network success, ", address, listener)
	//grpc server
	srv := grpc.NewServer()
	defer srv.GracefulStop()

	//service registers with etcd
	pbAuth.RegisterAuthServer(srv, rpc)
	rpcRegisterIP := config.Config.RpcRegisterIP
	if config.Config.RpcRegisterIP == "" {
		rpcRegisterIP, err = utils.GetLocalIP()
		if err != nil {
			log.Error("", "GetLocalIP failed ", err.Error())
		}
	}
	log.NewInfo("", "rpcRegisterIP", rpcRegisterIP)

	err = getcdv3.RegisterEtcd(rpc.etcdSchema, strings.Join(rpc.etcdAddr, ","), rpcRegisterIP, rpc.rpcPort, rpc.rpcRegisterName, 10)
	if err != nil {
		log.NewError(operationID, "RegisterEtcd failed ", err.Error(),
			rpc.etcdSchema, strings.Join(rpc.etcdAddr, ","), rpcRegisterIP, rpc.rpcPort, rpc.rpcRegisterName)
		return
	}
	log.NewInfo(operationID, "RegisterAuthServer ok ", rpc.etcdSchema, strings.Join(rpc.etcdAddr, ","), rpcRegisterIP, rpc.rpcPort, rpc.rpcRegisterName)
	err = srv.Serve(listener)
	if err != nil {
		log.NewError(operationID, "Serve failed ", err.Error())
		return
	}
	log.NewInfo(operationID, "rpc auth ok")
}
