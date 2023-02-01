package user

import (
	"Open_IM/internal/common/convert"
	chat "Open_IM/internal/rpc/msg"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db/controller"
	"Open_IM/pkg/common/db/relation"
	"Open_IM/pkg/common/log"
	promePkg "Open_IM/pkg/common/prometheus"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/common/tracelog"
	"Open_IM/pkg/getcdv3"
	pbFriend "Open_IM/pkg/proto/friend"
	pbGroup "Open_IM/pkg/proto/group"
	pbUser "Open_IM/pkg/proto/user"
	"Open_IM/pkg/utils"
	"context"
	"github.com/golang/protobuf/ptypes/wrappers"
	"net"
	"strconv"
	"strings"

	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"

	"google.golang.org/grpc"
)

type userServer struct {
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string
	controller.UserInterface
}

func NewUserServer(port int) *userServer {
	log.NewPrivateLog(constant.LogFileName)
	u := userServer{
		rpcPort:         port,
		rpcRegisterName: config.Config.RpcRegisterName.OpenImUserName,
		etcdSchema:      config.Config.Etcd.EtcdSchema,
		etcdAddr:        config.Config.Etcd.EtcdAddr,
	}
	//mysql init
	var mysql relation.Mysql
	var model relation.User
	err := mysql.InitConn().AutoMigrateModel(&model)
	if err != nil {
		panic("db init err:" + err.Error())
	}
	if mysql.GormConn() != nil {
		model.DB = mysql.GormConn()
	} else {
		panic("db init err:" + "conn is nil")
	}
	u.UserInterface = controller.NewUserController(model.DB)
	return &u
}

func (s *userServer) Run() {
	log.NewInfo("", "rpc user start...")

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
	log.NewInfo("", "listen network success, address ", address, listener)
	defer listener.Close()
	//grpc server
	var grpcOpts []grpc.ServerOption
	if config.Config.Prometheus.Enable {
		promePkg.NewGrpcRequestCounter()
		promePkg.NewGrpcRequestFailedCounter()
		promePkg.NewGrpcRequestSuccessCounter()
		grpcOpts = append(grpcOpts, []grpc.ServerOption{
			// grpc.UnaryInterceptor(promePkg.UnaryServerInterceptorProme),
			grpc.StreamInterceptor(grpcPrometheus.StreamServerInterceptor),
			grpc.UnaryInterceptor(grpcPrometheus.UnaryServerInterceptor),
		}...)
	}
	srv := grpc.NewServer(grpcOpts...)
	defer srv.GracefulStop()
	//Service registers with etcd
	pbUser.RegisterUserServer(srv, s)
	rpcRegisterIP := config.Config.RpcRegisterIP
	if config.Config.RpcRegisterIP == "" {
		rpcRegisterIP, err = utils.GetLocalIP()
		if err != nil {
			log.Error("", "GetLocalIP failed ", err.Error())
		}
	}
	err = getcdv3.RegisterEtcd(s.etcdSchema, strings.Join(s.etcdAddr, ","), rpcRegisterIP, s.rpcPort, s.rpcRegisterName, 10, "")
	if err != nil {
		log.NewError("", "RegisterEtcd failed ", err.Error(), s.etcdSchema, strings.Join(s.etcdAddr, ","), rpcRegisterIP, s.rpcPort, s.rpcRegisterName)
		panic(utils.Wrap(err, "register user module  rpc to etcd err"))
	}
	err = srv.Serve(listener)
	if err != nil {
		log.NewError("", "Serve failed ", err.Error())
		return
	}
	log.NewInfo("", "rpc  user success")
}

func (s *userServer) SyncJoinedGroupMemberFaceURL(ctx context.Context, userID string, faceURL string, operationID string, opUserID string) {
	etcdConn, err := getcdv3.GetConn(ctx, config.Config.RpcRegisterName.OpenImFriendName)
	if err != nil {
		return
	}
	client := pbGroup.NewGroupClient(etcdConn)
	newReq := &pbGroup.GetJoinedGroupListReq{FromUserID: userID}
	rpcResp, err := client.GetJoinedGroupList(ctx, newReq)
	if err != nil {
		return
	}

	for _, group := range rpcResp.Groups {
		req := &pbGroup.SetGroupMemberInfoReq{GroupID: group.GroupID, UserID: userID, FaceURL: &wrappers.StringValue{Value: faceURL}}
		_, err := client.SetGroupMemberInfo(ctx, req)
		if err != nil {
			return
		}
		chat.GroupMemberInfoSetNotification(operationID, opUserID, group.GroupID, userID)
	}
}

func (s *userServer) SyncJoinedGroupMemberNickname(ctx context.Context, userID string, newNickname, oldNickname string, operationID string, opUserID string) {
	etcdConn, err := getcdv3.GetConn(ctx, config.Config.RpcRegisterName.OpenImFriendName)
	if err != nil {
		return
	}
	client := pbGroup.NewGroupClient(etcdConn)
	newReq := &pbGroup.GetJoinedGroupListReq{FromUserID: userID}
	rpcResp, err := client.GetJoinedGroupList(ctx, newReq)
	if err != nil {
		return
	}
	req := pbGroup.GetUserInGroupMembersReq{UserID: userID}
	for _, group := range rpcResp.Groups {
		req.GroupIDs = append(req.GroupIDs, group.GroupID)
	}
	resp, err := client.GetUserInGroupMembers(ctx, &req)
	if err != nil {
		return
	}
	for _, v := range resp.Members {
		if v.Nickname == oldNickname {
			req := pbGroup.SetGroupMemberNicknameReq{Nickname: newNickname, GroupID: v.GroupID, UserID: v.UserID}
			_, err := client.SetGroupMemberNickname(ctx, &req)
			if err != nil {
				return
			}
			chat.GroupMemberInfoSetNotification(operationID, opUserID, v.GroupID, userID)
		}
	}
}

func (s *userServer) GetUsersInfo(ctx context.Context, req *pbUser.GetUsersInfoReq) (*pbUser.GetUsersInfoResp, error) {
	resp := &pbUser.GetUsersInfoResp{}
	users, err := s.Find(ctx, req.UserIDs)
	if err != nil {
		return nil, err
	}
	resp.UsersInfo, err = (*convert.DBUser)(nil).DB2PB(users)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *userServer) UpdateUserInfo(ctx context.Context, req *pbUser.UpdateUserInfoReq) (*pbUser.UpdateUserInfoResp, error) {
	resp := pbUser.UpdateUserInfoResp{}
	err := token_verify.CheckAccessV3(ctx, req.UserInfo.UserID)
	if err != nil {
		return nil, err
	}
	oldNickname := ""
	if req.UserInfo.Nickname != "" {
		u, err := s.Find(ctx, []string{req.UserInfo.UserID})
		if err != nil {
			return nil, err
		}
		oldNickname = u[0].Nickname
	}
	user, err := convert.NewPBUser(req.UserInfo).Convert()
	if err != nil {
		return nil, err
	}
	err = s.Update(ctx, []*relation.User{user})
	if err != nil {
		return nil, err
	}
	etcdConn, err := getcdv3.GetConn(ctx, config.Config.RpcRegisterName.OpenImFriendName)
	if err != nil {
		return nil, err
	}
	client := pbFriend.NewFriendClient(etcdConn)
	newReq := &pbFriend.GetFriendsReq{UserID: req.UserInfo.UserID}
	rpcResp, err := client.GetFriends(context.Background(), newReq)
	if err != nil {
		return nil, err
	}
	go func() {
		for _, v := range rpcResp.FriendsInfo {
			chat.FriendInfoUpdatedNotification(tracelog.GetOperationID(ctx), req.UserInfo.UserID, v.FriendUser.UserID, tracelog.GetOpUserID(ctx))
		}
	}()

	chat.UserInfoUpdatedNotification(tracelog.GetOperationID(ctx), tracelog.GetOpUserID(ctx), req.UserInfo.UserID)
	if req.UserInfo.FaceURL != "" {
		s.SyncJoinedGroupMemberFaceURL(ctx, req.UserInfo.UserID, req.UserInfo.FaceURL, tracelog.GetOperationID(ctx), tracelog.GetOpUserID(ctx))
	}
	if req.UserInfo.Nickname != "" {
		s.SyncJoinedGroupMemberNickname(ctx, req.UserInfo.UserID, req.UserInfo.Nickname, oldNickname, tracelog.GetOperationID(ctx), tracelog.GetOpUserID(ctx))
	}
	return &resp, nil
}

func (s *userServer) SetGlobalRecvMessageOpt(ctx context.Context, req *pbUser.SetGlobalRecvMessageOptReq) (*pbUser.SetGlobalRecvMessageOptResp, error) {
	resp := pbUser.SetGlobalRecvMessageOptResp{}
	m := make(map[string]interface{}, 1)
	m["global_recv_msg_opt"] = req.GlobalRecvMsgOpt
	err := s.UpdateByMap(ctx, req.UserID, m)
	if err != nil {
		return nil, err
	}
	chat.UserInfoUpdatedNotification(tracelog.GetOperationID(ctx), req.UserID, req.UserID)
	return &resp, nil
}

func (s *userServer) AccountCheck(ctx context.Context, req *pbUser.AccountCheckReq) (*pbUser.AccountCheckResp, error) {
	resp := pbUser.AccountCheckResp{}
	err := token_verify.CheckManagerUserID(ctx, tracelog.GetOpUserID(ctx))
	if err != nil {
		return nil, err
	}
	user, err := s.Find(ctx, req.CheckUserIDs)
	if err != nil {
		return nil, err
	}
	uidList := make([]string, 0)
	for _, v := range user {
		uidList = append(uidList, v.UserID)
	}
	var r []*pbUser.AccountCheckRespSingleUserStatus
	for _, v := range req.CheckUserIDs {
		temp := new(pbUser.AccountCheckRespSingleUserStatus)
		temp.UserID = v
		if utils.IsContain(v, uidList) {
			temp.AccountStatus = constant.Registered
		} else {
			temp.AccountStatus = constant.UnRegistered
		}
		r = append(r, temp)
	}
	return &resp, nil
}

func (s *userServer) GetUsers(ctx context.Context, req *pbUser.GetUsersReq) (*pbUser.GetUsersResp, error) {
	resp := pbUser.GetUsersResp{}
	var err error
	if req.UserID != "" {
		u, err := s.Find(ctx, []string{req.UserID})
		if err != nil {
			return nil, err
		}
		resp.Total = 1
		u1, err := convert.NewDBUser(u[0]).Convert()
		if err != nil {
			return nil, err
		}
		resp.Users = append(resp.Users, u1)
		return &resp, nil
	}

	if req.UserName != "" {
		usersDB, total, err := s.GetByName(ctx, req.UserName, req.Pagination.ShowNumber, req.Pagination.PageNumber)
		if err != nil {
			return nil, err
		}
		resp.Total = int32(total)
		for _, v := range usersDB {
			u1, err := convert.NewDBUser(v).Convert()
			if err != nil {
				return nil, err
			}
			resp.Users = append(resp.Users, u1)
		}
		return &resp, nil
	} else if req.Content != "" {
		usersDB, total, err := s.GetByNameAndID(ctx, req.UserName, req.Pagination.ShowNumber, req.Pagination.PageNumber)
		if err != nil {
			return nil, err
		}
		resp.Total = int32(total)
		for _, v := range usersDB {
			u1, err := convert.NewDBUser(v).Convert()
			if err != nil {
				return nil, err
			}
			resp.Users = append(resp.Users, u1)
		}
		return &resp, nil
	}

	usersDB, total, err := s.Get(ctx, req.Pagination.ShowNumber, req.Pagination.PageNumber)
	if err != nil {
		return nil, err
	}

	resp.Total = int32(total)

	for _, userDB := range usersDB {
		u, err := convert.NewDBUser(userDB).Convert()
		if err != nil {
			return nil, err
		}
		resp.Users = append(resp.Users, u)
	}
	return &resp, nil
}

func (s *userServer) UserRegister(ctx context.Context, req *pbUser.UserRegisterReq) (*pbUser.UserRegisterResp, error) {
	resp := pbUser.UserRegisterResp{}
	userIDs := make([]string, 0)
	for _, v := range req.Users {
		userIDs = append(userIDs, v.UserID)
	}
	exist, err := s.IsExist(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, constant.ErrRegisteredAlready.Wrap("exist")
	}
	users, err := (*convert.PBUser)(nil).PB2DB(req.Users)
	if err != nil {
		return nil, err
	}
	err = s.Create(ctx, users)
	if err != nil {
		return nil, err
	}
	return &resp, nil
}
