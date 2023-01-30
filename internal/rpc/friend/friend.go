package friend

import (
	chat "Open_IM/internal/rpc/msg"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db/controller"

	"Open_IM/pkg/common/db/relation"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/middleware"
	promePkg "Open_IM/pkg/common/prometheus"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/common/tools"
	"Open_IM/pkg/getcdv3"
	pbFriend "Open_IM/pkg/proto/friend"
	sdkws "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"net"
	"strconv"
	"strings"

	"google.golang.org/grpc"
)

type friendServer struct {
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string
	controller.FriendInterface
	controller.BlackInterface
}

func NewFriendServer(port int) *friendServer {
	log.NewPrivateLog(constant.LogFileName)
	f := friendServer{
		rpcPort:         port,
		rpcRegisterName: config.Config.RpcRegisterName.OpenImFriendName,
		etcdSchema:      config.Config.Etcd.EtcdSchema,
		etcdAddr:        config.Config.Etcd.EtcdAddr,
	}
	//mysql init
	var mysql relation.Mysql
	var model relation.Friend
	err := mysql.InitConn().AutoMigrateModel(&model)
	if err != nil {
		panic("db init err:" + err.Error())
	}
	if mysql.GormConn() != nil {
		model.DB = mysql.GormConn()
	} else {
		panic("db init err:" + "conn is nil")
	}
	f.FriendInterface = controller.NewFriendController(model.DB)
	f.BlackInterface = controller.NewBlackController(model.DB)
	return &f
}

func (s *friendServer) Run() {
	log.NewInfo("0", "friendServer run...")

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
	log.NewInfo("0", "listen ok ", address)
	defer listener.Close()
	//grpc server
	var grpcOpts []grpc.ServerOption
	grpcOpts = append(grpcOpts, grpc.UnaryInterceptor(middleware.RpcServerInterceptor))
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
	//User friend related services register to etcd
	pbFriend.RegisterFriendServer(srv, s)
	rpcRegisterIP := config.Config.RpcRegisterIP
	if config.Config.RpcRegisterIP == "" {
		rpcRegisterIP, err = utils.GetLocalIP()
		if err != nil {
			log.Error("", "GetLocalIP failed ", err.Error())
		}
	}
	log.NewInfo("", "rpcRegisterIP", rpcRegisterIP)
	err = getcdv3.RegisterEtcd(s.etcdSchema, strings.Join(s.etcdAddr, ","), rpcRegisterIP, s.rpcPort, s.rpcRegisterName, 10, "")
	if err != nil {
		log.NewError("0", "RegisterEtcd failed ", err.Error(), s.etcdSchema, strings.Join(s.etcdAddr, ","), rpcRegisterIP, s.rpcPort, s.rpcRegisterName)
		panic(utils.Wrap(err, "register friend module  rpc to etcd err"))
	}
	err = srv.Serve(listener)
	if err != nil {
		log.NewError("0", "Serve failed ", err.Error(), listener)
		return
	}
}

func (s *friendServer) AddBlacklist(ctx context.Context, req *pbFriend.AddBlacklistReq) (*pbFriend.AddBlacklistResp, error) {
	resp := &pbFriend.AddBlacklistResp{}
	if err := token_verify.CheckAccessV3(ctx, req.FromUserID); err != nil {
		return nil, err
	}
	black := relation.Black{OwnerUserID: req.FromUserID, BlockUserID: req.ToUserID, OperatorUserID: tools.OpUserID(ctx)}
	if err := s.BlackInterface.Create(ctx, []*relation.Black{&black}); err != nil {
		return nil, err
	}
	chat.BlackAddedNotification(req)
	return resp, nil
}

func (s *friendServer) AddFriend(ctx context.Context, req *pbFriend.AddFriendReq) (*pbFriend.AddFriendResp, error) {
	resp := &pbFriend.AddFriendResp{}
	if err := token_verify.CheckAccessV3(ctx, req.FromUserID); err != nil {
		return nil, err
	}
	if err := callbackBeforeAddFriendV1(req); err != nil {
		return nil, err
	}

	//检查toUserID是否存在
	if _, err := GetUsersInfo(ctx, []string{req.ToUserID}); err != nil {
		return nil, err
	}
	//from是否在to的好友列表里面
	err, in1, in2 := s.FriendInterface.CheckIn(ctx, req.FromUserID, req.ToUserID)
	if err != nil {
		return nil, err
	}
	if in1 && in2 {
		return nil, constant.ErrRelationshipAlready.Wrap()
	}
	if err = s.FriendInterface.AddFriendRequest(ctx, req.FromUserID, req.ToUserID, req.ReqMsg, req.Ex); err != nil {
		return nil, err
	}
	chat.FriendApplicationNotification(req)
	return resp, nil
}

func (s *friendServer) ImportFriend(ctx context.Context, req *pbFriend.ImportFriendReq) (*pbFriend.ImportFriendResp, error) {
	resp := &pbFriend.ImportFriendResp{}
	if err := token_verify.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	if _, err := GetUsersInfo(ctx, []string{req.FromUserID}); err != nil {
		return nil, err
	}

	var friends []*relation.Friend
	for _, userID := range utils.RemoveDuplicateElement(req.FriendUserIDList) {
		friends = append(friends, &relation.Friend{OwnerUserID: userID, FriendUserID: req.FromUserID, AddSource: constant.BecomeFriendByImport, OperatorUserID: tools.OpUserID(ctx)})
	}
	if len(friends) > 0 {
		if err := s.FriendInterface.BecomeFriend(ctx, friends); err != nil {
			return nil, err
		}
	}
	return resp, nil
}

// process Friend application
func (s *friendServer) friendApplyResponse(ctx context.Context, req *pbFriend.FriendApplyResponseReq) (*pbFriend.FriendApplyResponseResp, error) {
	resp := &pbFriend.FriendApplyResponseResp{}
	if err := token_verify.CheckAccessV3(ctx, req.ToUserID); err != nil {
		return nil, err
	}

	friendRequest := relation.FriendRequest{FromUserID: req.FromUserID, ToUserID: req.ToUserID, HandleMsg: req.HandleMsg, HandleResult: req.HandleResult}
	if req.HandleResult == constant.FriendResponseAgree {
		err := s.AgreeFriendRequest(ctx, &friendRequest)
		if err != nil {
			return nil, err
		}
		chat.FriendApplicationApprovedNotification(req)
		return resp, nil
	}
	if req.HandleResult == constant.FriendResponseRefuse {
		err := s.RefuseFriendRequest(ctx, &friendRequest)
		if err != nil {
			return nil, err
		}
		chat.FriendApplicationRejectedNotification(req)
		return resp, nil
	}
	return nil, constant.ErrArgs.Wrap("req.HandleResult != -1/1")
}

func (s *friendServer) DeleteFriend(ctx context.Context, req *pbFriend.DeleteFriendReq) (*pbFriend.DeleteFriendResp, error) {
	resp := &pbFriend.DeleteFriendResp{}
	if err := token_verify.CheckAccessV3(ctx, req.FromUserID); err != nil {
		return nil, err
	}
	if err := s.FriendInterface.Delete(ctx, req.FromUserID, req.ToUserID); err != nil {
		return nil, err
	}
	chat.FriendDeletedNotification(req)
	return resp, nil
}

func (s *friendServer) SetFriendRemark(ctx context.Context, req *pbFriend.SetFriendRemarkReq) (*pbFriend.SetFriendRemarkResp, error) {
	resp := &pbFriend.SetFriendRemarkResp{}
	if err := token_verify.CheckAccessV3(ctx, req.FromUserID); err != nil {
		return nil, err
	}
	if err := s.FriendInterface.UpdateRemark(ctx, req.FromUserID, req.ToUserID, req.Remark); err != nil {
		return nil, err
	}
	chat.FriendRemarkSetNotification(tools.OperationID(ctx), tools.OpUserID(ctx), req.FromUserID, req.ToUserID)
	return resp, nil
}

func (s *friendServer) RemoveBlacklist(ctx context.Context, req *pbFriend.RemoveBlacklistReq) (*pbFriend.RemoveBlacklistResp, error) {
	resp := &pbFriend.RemoveBlacklistResp{}
	//Parse token, to find current user information
	if err := token_verify.CheckAccessV3(ctx, req.FromUserID); err != nil {
		return nil, err
	}
	if err := s.BlackInterface.Delete(ctx, []*relation.Black{{OwnerUserID: req.FromUserID, BlockUserID: req.ToUserID}}); err != nil {
		return nil, err
	}
	chat.BlackDeletedNotification(req)
	return resp, nil
}

func (s *friendServer) GetFriendList(ctx context.Context, req *pbFriend.GetFriendListReq) (*pbFriend.GetFriendListResp, error) {
	resp := &pbFriend.GetFriendListResp{}
	if err := token_verify.CheckAccessV3(ctx, req.FromUserID); err != nil {
		return nil, err
	}
	friends, err := s.FriendInterface.FindOwnerUserID(ctx, req.FromUserID)
	if err != nil {
		return nil, err
	}
	userIDList := make([]string, 0, len(friends))
	for _, f := range friends {
		userIDList = append(userIDList, f.FriendUserID)
	}
	users, err := GetUsersInfo(ctx, userIDList)
	if err != nil {
		return nil, err
	}
	userMap := make(map[string]*sdkws.UserInfo)
	for i, user := range users {
		userMap[user.UserID] = users[i]
	}
	for _, friendUser := range friends {
		friendUserInfo := sdkws.FriendInfo{FriendUser: userMap[friendUser.FriendUserID]}
		utils.CopyStructFields(&friendUserInfo, friendUser)
		resp.FriendInfoList = append(resp.FriendInfoList, &friendUserInfo)
	}
	return resp, nil
}

// received
func (s *friendServer) GetFriendApplyList(ctx context.Context, req *pbFriend.GetFriendApplyListReq) (*pbFriend.GetFriendApplyListResp, error) {
	resp := &pbFriend.GetFriendApplyListResp{}
	//Parse token, to find current user information
	if err := token_verify.CheckAccessV3(ctx, req.FromUserID); err != nil {
		return nil, err
	}
	//	Find the  current user friend applications received
	friendRequests, err := s.friendRequestModel.FindToUserID(ctx, req.FromUserID)
	if err != nil {
		return nil, err
	}
	userIDList := make([]string, 0, len(friendRequests))
	for _, f := range friendRequests {
		userIDList = append(userIDList, f.FromUserID)
	}
	users, err := GetPublicUserInfoBatch(ctx, userIDList)
	if err != nil {
		return nil, err
	}
	userMap := make(map[string]*sdkws.PublicUserInfo)
	for i, user := range users {
		userMap[user.UserID] = users[i]
	}
	for _, friendRequest := range friendRequests {
		var userInfo sdkws.FriendRequest
		if u, ok := userMap[friendRequest.FromUserID]; ok {
			utils.CopyStructFields(&userInfo, u)
		}
		utils.CopyStructFields(&userInfo, friendRequest)
		resp.FriendRequestList = append(resp.FriendRequestList, &userInfo)
	}
	return resp, nil
}

func (s *friendServer) GetSelfApplyList(ctx context.Context, req *pbFriend.GetSelfApplyListReq) (*pbFriend.GetSelfApplyListResp, error) {
	resp := &pbFriend.GetSelfApplyListResp{}
	//Parse token, to find current user information
	if err := token_verify.CheckAccessV3(ctx, req.FromUserID); err != nil {
		return nil, err
	}
	//	Find the self add other userinfo
	friendRequests, err := s.FriendRequestInterface.FindFromUserID(ctx, req.FromUserID)
	if err != nil {
		return nil, err
	}
	userIDList := make([]string, 0, len(friendRequests))
	for _, f := range friendRequests {
		userIDList = append(userIDList, f.ToUserID)
	}
	users, err := GetPublicUserInfoBatch(ctx, userIDList)
	if err != nil {
		return nil, err
	}
	userMap := make(map[string]*sdkws.PublicUserInfo)
	for i, user := range users {
		userMap[user.UserID] = users[i]
	}
	for _, friendRequest := range friendRequests {
		var userInfo sdkws.FriendRequest
		if u, ok := userMap[friendRequest.ToUserID]; ok {
			utils.CopyStructFields(&userInfo, u)
		}
		utils.CopyStructFields(&userInfo, friendRequest)
		resp.FriendRequestList = append(resp.FriendRequestList, &userInfo)
	}
	return resp, nil
}

func (s *friendServer) GetBlacklist(ctx context.Context, req *pbFriend.GetBlacklistReq) (*pbFriend.GetBlacklistResp, error) {
	resp := &pbFriend.GetBlacklistResp{}
	if err := token_verify.CheckAccessV3(ctx, req.FromUserID); err != nil {
		return nil, err
	}
	blacks, err := s.BlackInterface.FindByOwnerUserID(ctx, req.FromUserID)
	if err != nil {
		return nil, err
	}
	blackIDList := make([]string, 0, len(blacks))
	for _, black := range blacks {
		blackIDList = append(blackIDList, black.BlockUserID)
	}
	resp.BlackUserInfoList, err = GetPublicUserInfoBatch(ctx, blackIDList)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *friendServer) IsInBlackList(ctx context.Context, req *pbFriend.IsInBlackListReq) (*pbFriend.IsInBlackListResp, error) {
	resp := &pbFriend.IsInBlackListResp{}
	if err := token_verify.CheckAccessV3(ctx, req.FromUserID); err != nil {
		return nil, err
	}
	exist, err := s.BlackInterface.IsExist(ctx, req.FromUserID, req.ToUserID)
	if err != nil {
		return nil, err
	}
	resp.Response = exist
	return resp, nil
}

func (s *friendServer) IsFriend(ctx context.Context, req *pbFriend.IsFriendReq) (*pbFriend.IsFriendResp, error) {
	resp := &pbFriend.IsFriendResp{}
	if err := token_verify.CheckAccessV3(ctx, req.FromUserID); err != nil {
		return nil, err
	}
	exist, err := s.FriendInterface.IsExist(ctx, req.FromUserID, req.ToUserID)
	if err != nil {
		return nil, err
	}
	resp.Response = exist
	return resp, nil
}
