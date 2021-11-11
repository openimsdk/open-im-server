package user

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbUser "Open_IM/pkg/proto/user"
	"Open_IM/pkg/utils"
	"context"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"strings"
)

type userServer struct {
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string
}

func NewUserServer(port int) *userServer {
	log.NewPrivateLog("user")
	return &userServer{
		rpcPort:         port,
		rpcRegisterName: config.Config.RpcRegisterName.OpenImUserName,
		etcdSchema:      config.Config.Etcd.EtcdSchema,
		etcdAddr:        config.Config.Etcd.EtcdAddr,
	}
}

func (s *userServer) Run() {
	log.Info("", "", "rpc user init....")

	ip := utils.ServerIP
	registerAddress := ip + ":" + strconv.Itoa(s.rpcPort)
	//listener network
	listener, err := net.Listen("tcp", registerAddress)
	if err != nil {
		log.InfoByArgs("listen network failed,err=%s", err.Error())
		return
	}
	log.Info("", "", "listen network success, address = %s", registerAddress)
	defer listener.Close()
	//grpc server
	srv := grpc.NewServer()
	defer srv.GracefulStop()
	//Service registers with etcd
	pbUser.RegisterUserServer(srv, s)
	err = getcdv3.RegisterEtcd(s.etcdSchema, strings.Join(s.etcdAddr, ","), ip, s.rpcPort, s.rpcRegisterName, 10)
	if err != nil {
		log.ErrorByArgs("register rpc token to etcd failed,err=%s", err.Error())
		return
	}
	err = srv.Serve(listener)
	if err != nil {
		log.ErrorByArgs("listen token failed,err=%s", err.Error())
		return
	}
	log.Info("", "", "rpc token init success")
}

func (s *userServer) GetUserInfo(ctx context.Context, req *pbUser.GetUserInfoReq) (*pbUser.GetUserInfoResp, error) {
	log.InfoByKv("rpc get_user_info is server", req.OperationID)

	var userInfoList []*pbUser.UserInfo
	//Obtain user information according to userID
	if len(req.UserIDList) > 0 {
		for _, userID := range req.UserIDList {
			var userInfo pbUser.UserInfo
			user, err := im_mysql_model.FindUserByUID(userID)
			if err != nil {
				log.ErrorByKv("search userinfo failed", req.OperationID, "userID", userID, "err=%s", err.Error())
				continue
			}
			userInfo.Uid = user.UID
			userInfo.Icon = user.Icon
			userInfo.Name = user.Name
			userInfo.Gender = user.Gender
			userInfo.Mobile = user.Mobile
			userInfo.Birth = user.Birth
			userInfo.Email = user.Email
			userInfo.Ex = user.Ex
			userInfoList = append(userInfoList, &userInfo)
		}
	} else {
		return &pbUser.GetUserInfoResp{ErrorCode: 999, ErrorMsg: "uidList is nil"}, nil
	}
	log.InfoByKv("rpc get userInfo return success", req.OperationID, "token", req.Token)
	return &pbUser.GetUserInfoResp{
		ErrorCode: 0,
		ErrorMsg:  "",
		Data:      userInfoList,
	}, nil
}
