package friend

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbFriend "Open_IM/pkg/proto/friend"
	"Open_IM/pkg/utils"
	"context"
	"fmt"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"strings"
)

type friendServer struct {
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string
}

func NewFriendServer(port int) *friendServer {
	log.NewPrivateLog("friend")
	return &friendServer{
		rpcPort:         port,
		rpcRegisterName: config.Config.RpcRegisterName.OpenImFriendName,
		etcdSchema:      config.Config.Etcd.EtcdSchema,
		etcdAddr:        config.Config.Etcd.EtcdAddr,
	}
}

func (s *friendServer) Run() {
	log.Info("", "", fmt.Sprintf("rpc friend init...."))

	ip := utils.ServerIP
	registerAddress := ip + ":" + strconv.Itoa(s.rpcPort)
	//listener network
	listener, err := net.Listen("tcp", registerAddress)
	if err != nil {
		log.InfoByArgs(fmt.Sprintf("Failed to listen rpc friend network,err=%s", err.Error()))
		return
	}
	log.Info("", "", "listen network success, address = %s", registerAddress)
	defer listener.Close()
	//grpc server
	srv := grpc.NewServer()
	defer srv.GracefulStop()
	//User friend related services register to etcd
	pbFriend.RegisterFriendServer(srv, s)
	err = getcdv3.RegisterEtcd(s.etcdSchema, strings.Join(s.etcdAddr, ","), ip, s.rpcPort, s.rpcRegisterName, 10)
	if err != nil {
		log.ErrorByArgs("register rpc fiend service to etcd failed,err=%s", err.Error())
		return
	}
	err = srv.Serve(listener)
	if err != nil {
		log.ErrorByArgs("listen rpc friend error,err=%s", err.Error())
		return
	}
}

func (s *friendServer) GetFriendsInfo(ctx context.Context, req *pbFriend.GetFriendsInfoReq) (*pbFriend.GetFriendInfoResp, error) {
	log.Info(req.Token, req.OperationID, "rpc search user is server,args=%s", req.String())
	var (
		isInBlackList int32
		isFriend      int32
		comment       string
	)
	//Parse token, to find current user information
	claims, err := utils.ParseToken(req.Token)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,parse token failed", err.Error())
		return &pbFriend.GetFriendInfoResp{ErrorCode: config.ErrParseToken.ErrCode, ErrorMsg: config.ErrParseToken.ErrMsg}, nil
	}
	friendShip, err := im_mysql_model.FindFriendRelationshipFromFriend(claims.UID, req.Uid)
	if err == nil {
		isFriend = constant.FriendFlag
		comment = friendShip.Comment
	}
	friendUserInfo, err := im_mysql_model.FindUserByUID(req.Uid)
	if err != nil {
		log.Error(req.Token, req.OperationID, "err=%s,no this user", err.Error())
		return &pbFriend.GetFriendInfoResp{ErrorCode: config.ErrSearchUserInfo.ErrCode, ErrorMsg: config.ErrSearchUserInfo.ErrMsg}, nil
	}
	err = im_mysql_model.FindRelationshipFromBlackList(claims.UID, req.Uid)
	if err == nil {
		isInBlackList = constant.BlackListFlag
	}
	log.Info(req.Token, req.OperationID, "rpc search friend success return")
	return &pbFriend.GetFriendInfoResp{
		ErrorCode: 0,
		ErrorMsg:  "",
		Data: &pbFriend.GetFriendData{
			Uid:           friendUserInfo.UID,
			Icon:          friendUserInfo.Icon,
			Name:          friendUserInfo.Name,
			Gender:        friendUserInfo.Gender,
			Mobile:        friendUserInfo.Mobile,
			Birth:         friendUserInfo.Birth,
			Email:         friendUserInfo.Email,
			Ex:            friendUserInfo.Ex,
			Comment:       comment,
			IsFriend:      isFriend,
			IsInBlackList: isInBlackList,
		},
	}, nil

}
