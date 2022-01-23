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
	"net"
	"strconv"
	"strings"

	"google.golang.org/grpc"
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
		user.Birth = utils.UnixSecondToTime(int64(req.UserInfo.Birth))
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
		chat.UserInfoUpdatedNotification(req.OperationID, req.UserInfo.UserID, v.FriendUser.UserID)
	}
	chat.UserInfoUpdatedNotification(req.OperationID, req.UserInfo.UserID, req.OpUserID)
	return &pbUser.UpdateUserInfoResp{CommonResp: &pbUser.CommonResp{}}, nil
}

func (s *userServer) GetUser(ctx context.Context, req *pbUser.GetUserReq) (*pbUser.GetUserResp, error) {
	log.NewInfo(req.OperationID, "GetAllUserID args ", req.String())
	resp := &pbUser.GetUserResp{}
	user, err := imdb.GetUserByUserID(req.UserId)
	if err != nil {
		log.NewError(req.OperationID, "SelectAllUserID false ", err.Error())
		return resp, nil
	}
	resp.User.CreateTime = user.CreateTime.String()
	resp.User.ProfilePhoto = ""
	resp.User.Nickname = user.Nickname
	resp.User.UserID = user.UserID
	return resp, nil
}

func (s *userServer) GetUsers(ctx context.Context, req *pbUser.GetUsersReq) (*pbUser.GetUsersResp, error) {
	log.NewInfo(req.OperationID, "GetUsers args ", req.String())
	resp := &pbUser.GetUsersResp{}
	users, err := imdb.GetUsers(req.Pagination.ShowNumber, req.Pagination.PageNumber)
	if err != nil {
		log.NewError(req.OperationID, "SelectAllUserID false ", err.Error())
		return resp, nil
	}
	for _, v := range users {
		resp.User = append(resp.User, &pbUser.User{
			ProfilePhoto: "",
			UserID:       v.UserID,
			CreateTime:   v.CreateTime.String(),
			Nickname:     v.Nickname,
		})
	}
	return resp, nil
}

func (s *userServer) ResignUser(ctx context.Context, req *pbUser.ResignUserReq) (*pbUser.ResignUserResp, error) {
	log.NewInfo(req.OperationID, "ResignUser args ", req.String())

	return &pbUser.ResignUserResp{}, nil
}

func (s *userServer) AlterUser(ctx context.Context, req *pbUser.AlterUserReq) (*pbUser.AlterUserResp, error) {
	return &pbUser.AlterUserResp{}, nil
}

func (s *userServer) AddUser(ctx context.Context, req *pbUser.AddUserReq) (*pbUser.AddUserResp, error) {
	return &pbUser.AddUserResp{}, nil
}

func (s *userServer) BlockUser(ctx context.Context, req *pbUser.BlockUserReq) (*pbUser.BlockUserResp, error) {
	return &pbUser.BlockUserResp{}, nil
}

func (s *userServer) UnBlockUser(ctx context.Context, req *pbUser.UnBlockUserReq) (*pbUser.UnBlockUserResp, error) {
	return &pbUser.UnBlockUserResp{}, nil
}

func (s *userServer) GetBlockUsers(ctx context.Context, req *pbUser.GetBlockUsersReq) (*pbUser.GetBlockUsersResp, error) {
	return &pbUser.GetBlockUsersResp{}, nil
}
