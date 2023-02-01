package friend

import (
	"Open_IM/internal/common/convert"
	chat "Open_IM/internal/rpc/msg"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db/controller"
	"Open_IM/pkg/common/db/relation"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/middleware"
	promePkg "Open_IM/pkg/common/prometheus"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/common/tracelog"
	"Open_IM/pkg/getcdv3"
	pbFriend "Open_IM/pkg/proto/friend"
	sdkws "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"net"
	"strconv"
	"strings"

	"Open_IM/internal/common/check"
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
	var model relation.FriendGorm
	err := mysql.InitConn().AutoMigrateModel(&relation.FriendModel{})
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

func (s *friendServer) AddFriend(ctx context.Context, req *pbFriend.AddFriendReq) (*pbFriend.AddFriendResp, error) {
	resp := &pbFriend.AddFriendResp{}
	if err := token_verify.CheckAccessV3(ctx, req.FromUserID); err != nil {
		return nil, err
	}
	if err := callbackBeforeAddFriendV1(req); err != nil {
		return nil, err
	}

	//检查toUserID fromUserID是否存在
	if _, err := check.GetUsersInfo(ctx, req.ToUserID, req.FromUserID); err != nil {
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

func (s *friendServer) ImportFriends(ctx context.Context, req *pbFriend.ImportFriendReq) (*pbFriend.ImportFriendResp, error) {
	resp := &pbFriend.ImportFriendResp{}
	if err := token_verify.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	if _, err := check.GetUsersInfo(ctx, req.OwnerUserID); err != nil {
		return nil, err
	}

	var friends []*relation.Friend
	for _, userID := range utils.RemoveDuplicateElement(req.FriendUserIDs) {
		friends = append(friends, &relation.Friend{OwnerUserID: userID, FriendUserID: req.OwnerUserID, AddSource: constant.BecomeFriendByImport, OperatorUserID: tracelog.GetOpUserID(ctx)})
	}
	if len(friends) > 0 {
		if err := s.FriendInterface.BecomeFriend(ctx, friends); err != nil {
			return nil, err
		}
	}
	return resp, nil
}

// process Friend application
func (s *friendServer) RespondFriendApply(ctx context.Context, req *pbFriend.RespondFriendApplyReq) (*pbFriend.RespondFriendApplyResp, error) {
	resp := &pbFriend.RespondFriendApplyResp{}
	if err := check.Access(ctx, req.ToUserID); err != nil {
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
	if err := check.Access(ctx, req.OwnerUserID); err != nil {
		return nil, err
	}
	if err := s.FriendInterface.Delete(ctx, req.OwnerUserID, req.FriendUserID); err != nil {
		return nil, err
	}
	chat.FriendDeletedNotification(req)
	return resp, nil
}

func (s *friendServer) SetFriendRemark(ctx context.Context, req *pbFriend.SetFriendRemarkReq) (*pbFriend.SetFriendRemarkResp, error) {
	resp := &pbFriend.SetFriendRemarkResp{}
	if err := check.Access(ctx, req.OwnerUserID); err != nil {
		return nil, err
	}
	if err := s.FriendInterface.UpdateRemark(ctx, req.OwnerUserID, req.FriendUserID, req.Remark); err != nil {
		return nil, err
	}
	chat.FriendRemarkSetNotification(tracelog.GetOperationID(ctx), tracelog.GetOpUserID(ctx), req.OwnerUserID, req.FriendUserID)
	return resp, nil
}

func (s *friendServer) GetFriends(ctx context.Context, req *pbFriend.GetFriendsReq) (*pbFriend.GetFriendsResp, error) {
	resp := &pbFriend.GetFriendsResp{}
	if err := check.Access(ctx, req.UserID); err != nil {
		return nil, err
	}
	friends, err := s.FriendInterface.FindOwnerFriends(ctx, req.UserID, req.Pagination.PageNumber, req.Pagination.ShowNumber)
	if err != nil {
		return nil, err
	}
	userIDList := make([]string, 0, len(friends))
	for _, f := range friends {
		userIDList = append(userIDList, f.FriendUserID)
	}
	users, err := check.GetUsersInfo(ctx, userIDList)
	if err != nil {
		return nil, err
	}
	userMap := make(map[string]*sdkws.UserInfo)
	for i, user := range users {
		userMap[user.UserID] = users[i]
	}
	for _, friendUser := range friends {

		friendUserInfo, err := (convert.NewDBFriend(friendUser)).Convert()
		if err != nil {
			return nil, err
		}
		resp.FriendsInfo = append(resp.FriendsInfo, friendUserInfo)
	}
	return resp, nil
}

// 获取接收到的好友申请（即别人主动申请的）
func (s *friendServer) GetToFriendsApply(ctx context.Context, req *pbFriend.GetToFriendsApplyReq) (*pbFriend.GetToFriendsApplyResp, error) {
	resp := &pbFriend.GetToFriendsApplyResp{}
	if err := check.Access(ctx, req.UserID); err != nil {
		return nil, err
	}
	friendRequests, err := s.FriendInterface.FindFriendRequestToMe(ctx, req.UserID, req.Pagination.PageNumber, req.Pagination.ShowNumber)
	if err != nil {
		return nil, err
	}
	for _, v := range friendRequests {
		fUser, err := convert.NewDBFriendRequest(v).Convert()
		if err != nil {
			return nil, err
		}
		resp.FriendRequests = append(resp.FriendRequests, fUser)
	}
	return resp, nil
}

// 获取主动发出去的好友申请列表
func (s *friendServer) GetFromFriendsApply(ctx context.Context, req *pbFriend.GetFromFriendsApplyReq) (*pbFriend.GetFromFriendsApplyResp, error) {
	resp := &pbFriend.GetFromFriendsApplyResp{}
	if err := check.Access(ctx, req.UserID); err != nil {
		return nil, err
	}
	friendRequests, err := s.FriendInterface.FindFriendRequestFromMe(ctx, req.UserID, req.Pagination.PageNumber, req.Pagination.ShowNumber)
	if err != nil {
		return nil, err
	}
	for _, v := range friendRequests {
		fUser, err := convert.NewDBFriendRequest(v).Convert()
		if err != nil {
			return nil, err
		}
		resp.FriendRequests = append(resp.FriendRequests, fUser)
	}
	return resp, nil
}

func (s *friendServer) IsFriend(ctx context.Context, req *pbFriend.IsFriendReq) (*pbFriend.IsFriendResp, error) {
	resp := &pbFriend.IsFriendResp{}
	err, in1, in2 := s.FriendInterface.CheckIn(ctx, req.UserID1, req.UserID2)
	if err != nil {
		return nil, err
	}
	resp.InUser1Friends = in1
	resp.InUser2Friends = in2
	return resp, nil
}

func (s *friendServer) GetFriendsInfo(ctx context.Context, req *pbFriend.GetFriendsInfoReq) (*pbFriend.GetFriendsInfoResp, error) {
	resp := pbFriend.GetFriendsInfoResp{}
	friends, err := s.FriendInterface.FindFriends(ctx, req.OwnerUserID, req.FriendUserIDs)
	if err != nil {
		return nil, err
	}
	for _, v := range friends {
		fUser, err := convert.NewDBFriend(v).Convert()
		if err != nil {
			return nil, err
		}
		resp.FriendsInfo = append(resp.FriendsInfo, fUser)
	}
	return &resp, nil
}
