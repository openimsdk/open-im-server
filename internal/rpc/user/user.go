package user

import (
	chat "Open_IM/internal/rpc/msg"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbFriend "Open_IM/pkg/proto/friend"
	sdkws "Open_IM/pkg/proto/sdk_ws"
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
	log.NewInfo("0", "", "rpc user start...")

	ip := utils.ServerIP
	registerAddress := ip + ":" + strconv.Itoa(s.rpcPort)
	//listener network
	listener, err := net.Listen("tcp", registerAddress)
	if err != nil {
		log.NewError("0", "listen network failed ", err.Error(), registerAddress)
		return
	}
	log.NewInfo("0", "listen network success, address ", registerAddress, listener)
	defer listener.Close()
	//grpc server
	srv := grpc.NewServer()
	defer srv.GracefulStop()
	//Service registers with etcd
	pbUser.RegisterUserServer(srv, s)
	err = getcdv3.RegisterEtcd(s.etcdSchema, strings.Join(s.etcdAddr, ","), ip, s.rpcPort, s.rpcRegisterName, 10)
	if err != nil {
		log.NewError("0", "RegisterEtcd failed ", err.Error(), s.etcdSchema, strings.Join(s.etcdAddr, ","), ip, s.rpcPort, s.rpcRegisterName)
		return
	}
	err = srv.Serve(listener)
	if err != nil {
		log.NewError("0", "Serve failed ", err.Error())
		return
	}
	log.NewInfo("0", "rpc  user success")
}

func (s *userServer) GetUserInfo(ctx context.Context, req *pbUser.GetUserInfoReq) (*pbUser.GetUserInfoResp, error) {
	log.NewInfo(req.OperationID, "GetUserInfo args ", req.String())
	var userInfoList []*sdkws.UserInfo
	if len(req.UserIDList) > 0 {
		for _, userID := range req.UserIDList {
			var userInfo sdkws.UserInfo
			user, err := imdb.GetUserByUserID(userID)
			if err != nil {
				log.NewError(req.OperationID, "GetUserByUserID failed ", err.Error(), userID)
				continue
			}
			utils.CopyStructFields(&userInfo, user)
			userInfoList = append(userInfoList, &userInfo)
		}
	} else {

		return &pbUser.GetUserInfoResp{CommonResp: &pbUser.CommonResp{ErrCode: constant.ErrArgs.ErrCode, ErrMsg: constant.ErrArgs.ErrMsg}}, nil
	}
	log.NewInfo(req.OperationID, "GetUserInfo rpc return ", pbUser.GetUserInfoResp{CommonResp: &pbUser.CommonResp{}, UserInfoList: userInfoList})
	return &pbUser.GetUserInfoResp{CommonResp: &pbUser.CommonResp{}, UserInfoList: userInfoList}, nil
}

func (s *userServer) SetReceiveMessageOpt(ctx context.Context, req *pbUser.SetReceiveMessageOptReq) (*pbUser.SetReceiveMessageOptResp, error) {
	log.NewInfo(req.OperationID, "SetReceiveMessageOpt args ", req.String())
	m := make(map[string]int, len(req.ConversationIDList))
	for _, v := range req.ConversationIDList {
		m[v] = int(req.Opt)
	}
	err := db.DB.SetMultiConversationMsgOpt(req.FromUserID, m)
	if err != nil {
		log.NewError(req.OperationID, "SetMultiConversationMsgOpt failed ", err.Error(), req)
		return &pbUser.SetReceiveMessageOptResp{CommonResp: &pbUser.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}
	resp := pbUser.SetReceiveMessageOptResp{CommonResp: &pbUser.CommonResp{}}

	for _, v := range req.ConversationIDList {
		resp.ConversationOptResultList = append(resp.ConversationOptResultList, &pbUser.OptResult{ConversationID: v, Result: req.Opt})
	}
	log.NewInfo(req.OperationID, "SetReceiveMessageOpt rpc return ", resp.String())
	return &resp, nil
}

func (s *userServer) GetReceiveMessageOpt(ctx context.Context, req *pbUser.GetReceiveMessageOptReq) (*pbUser.GetReceiveMessageOptResp, error) {
	log.NewInfo(req.OperationID, "GetReceiveMessageOpt args ", req.String())
	m, err := db.DB.GetMultiConversationMsgOpt(req.FromUserID, req.ConversationIDList)
	if err != nil {
		log.NewError(req.OperationID, "GetMultiConversationMsgOpt failed ", err.Error(), req.FromUserID, req.ConversationIDList)
		return &pbUser.GetReceiveMessageOptResp{CommonResp: &pbUser.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}
	resp := pbUser.GetReceiveMessageOptResp{CommonResp: &pbUser.CommonResp{}}
	for k, v := range m {
		resp.ConversationOptResultList = append(resp.ConversationOptResultList, &pbUser.OptResult{ConversationID: k, Result: int32(v)})
	}
	log.NewInfo(req.OperationID, "GetReceiveMessageOpt rpc return ", resp.String())
	return &resp, nil
}

func (s *userServer) GetAllConversationMsgOpt(ctx context.Context, req *pbUser.GetAllConversationMsgOptReq) (*pbUser.GetAllConversationMsgOptResp, error) {
	log.NewInfo(req.OperationID, "GetAllConversationMsgOpt args ", req.String())
	m, err := db.DB.GetAllConversationMsgOpt(req.FromUserID)
	if err != nil {
		log.NewError(req.OperationID, "GetAllConversationMsgOpt failed ", err.Error(), req.FromUserID)
		return &pbUser.GetAllConversationMsgOptResp{CommonResp: &pbUser.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}
	resp := pbUser.GetAllConversationMsgOptResp{CommonResp: &pbUser.CommonResp{}}
	for k, v := range m {
		resp.ConversationOptResultList = append(resp.ConversationOptResultList, &pbUser.OptResult{ConversationID: k, Result: int32(v)})
	}
	log.NewInfo(req.OperationID, "GetAllConversationMsgOpt rpc return ", resp.String())
	return &resp, nil
}
func (s *userServer) DeleteUsers(_ context.Context, req *pbUser.DeleteUsersReq) (*pbUser.DeleteUsersResp, error) {
	log.NewInfo(req.OperationID, "DeleteUsers args ", req.String())
	if !token_verify.IsMangerUserID(req.OpUserID) {
		log.NewError(req.OperationID, "IsMangerUserID false ", req.OpUserID)
		return &pbUser.DeleteUsersResp{CommonResp: &pbUser.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}, FailedUserIDList: req.DeleteUserIDList}, nil
	}
	var common pbUser.CommonResp
	resp := pbUser.DeleteUsersResp{CommonResp: &common}
	for _, userID := range req.DeleteUserIDList {
		i := imdb.DeleteUser(userID)
		if i == 0 {
			log.NewError(req.OperationID, "delete user error", userID)
			common.ErrCode = 201
			common.ErrMsg = "some uid deleted failed"
			resp.FailedUserIDList = append(resp.FailedUserIDList, userID)
		}
	}
	log.NewInfo(req.OperationID, "DeleteUsers rpc return ", resp.String())
	return &resp, nil
}

func (s *userServer) GetAllUserID(_ context.Context, req *pbUser.GetAllUserIDReq) (*pbUser.GetAllUserIDResp, error) {
	log.NewInfo(req.OperationID, "GetAllUserID args ", req.String())
	if !token_verify.IsMangerUserID(req.OpUserID) {
		log.NewError(req.OperationID, "IsMangerUserID false ", req.OpUserID)
		return &pbUser.GetAllUserIDResp{CommonResp: &pbUser.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}}, nil
	}
	uidList, err := imdb.SelectAllUserID()
	if err != nil {
		log.NewError(req.OperationID, "SelectAllUserID false ", err.Error())
		return &pbUser.GetAllUserIDResp{CommonResp: &pbUser.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	} else {
		log.NewInfo(req.OperationID, "GetAllUserID rpc return ", pbUser.GetAllUserIDResp{CommonResp: &pbUser.CommonResp{}, UserIDList: uidList})
		return &pbUser.GetAllUserIDResp{CommonResp: &pbUser.CommonResp{}, UserIDList: uidList}, nil
	}
}

func (s *userServer) AccountCheck(_ context.Context, req *pbUser.AccountCheckReq) (*pbUser.AccountCheckResp, error) {
	log.NewInfo(req.OperationID, "AccountCheck args ", req.String())
	if !token_verify.IsMangerUserID(req.OpUserID) {
		log.NewError(req.OperationID, "IsMangerUserID false ", req.OpUserID)
		return &pbUser.AccountCheckResp{CommonResp: &pbUser.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}}, nil
	}
	uidList, err := imdb.SelectSomeUserID(req.CheckUserIDList)
	log.NewDebug(req.OperationID, "from db uid list is:", uidList)
	if err != nil {
		log.NewError(req.OperationID, "SelectSomeUserID failed ", err.Error(), req.CheckUserIDList)
		return &pbUser.AccountCheckResp{CommonResp: &pbUser.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	} else {
		var r []*pbUser.AccountCheckResp_SingleUserStatus
		for _, v := range req.CheckUserIDList {
			temp := new(pbUser.AccountCheckResp_SingleUserStatus)
			temp.UserID = v
			if utils.IsContain(v, uidList) {
				temp.AccountStatus = constant.Registered
			} else {
				temp.AccountStatus = constant.UnRegistered
			}
			r = append(r, temp)
		}
		resp := pbUser.AccountCheckResp{CommonResp: &pbUser.CommonResp{ErrCode: 0, ErrMsg: ""}, ResultList: r}
		log.NewInfo(req.OperationID, "AccountCheck rpc return ", resp.String())
		return &resp, nil
	}

}

func (s *userServer) UpdateUserInfo(ctx context.Context, req *pbUser.UpdateUserInfoReq) (*pbUser.UpdateUserInfoResp, error) {
	log.NewInfo(req.OperationID, "UpdateUserInfo args ", req.String())
	if !token_verify.CheckAccess(req.OpUserID, req.UserInfo.UserID) {
		log.NewError(req.OperationID, "CheckAccess false ", req.OpUserID, req.UserInfo.UserID)
		return &pbUser.UpdateUserInfoResp{CommonResp: &pbUser.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}}, nil
	}

	var user db.User
	utils.CopyStructFields(&user, req.UserInfo)
	if req.UserInfo.Birth != 0 {
		user.Birth = utils.UnixSecondToTime(req.UserInfo.Birth)
	}
	err := imdb.UpdateUserInfo(user)
	if err != nil {
		log.NewError(req.OperationID, "UpdateUserInfo failed ", err.Error(), user)
		return &pbUser.UpdateUserInfoResp{CommonResp: &pbUser.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}
	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImFriendName)
	client := pbFriend.NewFriendClient(etcdConn)
	newReq := &pbFriend.GetFriendListReq{
		CommID: &pbFriend.CommID{OperationID: req.OperationID, FromUserID: req.UserInfo.UserID, OpUserID: req.OpUserID},
	}

	RpcResp, err := client.GetFriendList(context.Background(), newReq)
	if err != nil {
		log.NewError(req.OperationID, "GetFriendList failed ", err.Error(), newReq)
		return &pbUser.UpdateUserInfoResp{CommonResp: &pbUser.CommonResp{}}, nil
	}
	for _, v := range RpcResp.FriendInfoList {
		chat.FriendInfoChangedNotification(req.OperationID, req.OpUserID, req.UserInfo.UserID, v.FriendUser.UserID)
	}
	chat.SelfInfoUpdatedNotification(req.OperationID, req.UserInfo.UserID)
	return &pbUser.UpdateUserInfoResp{CommonResp: &pbUser.CommonResp{}}, nil
}
