package friend

import (
	"Open_IM/internal/common/convert"
	chat "Open_IM/internal/rpc/msg"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db/controller"
	"Open_IM/pkg/common/db/relation"
	relationTb "Open_IM/pkg/common/db/table/relation"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/middleware"
	promePkg "Open_IM/pkg/common/prometheus"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/common/tracelog"
	pbFriend "Open_IM/pkg/proto/friend"
	pbUser "Open_IM/pkg/proto/user"
	"Open_IM/pkg/utils"
	"context"
	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"net"
	"strconv"
	"strings"

	"Open_IM/internal/common/check"
	"github.com/OpenIMSDK/getcdv3"
	"google.golang.org/grpc"
)

type friendServer struct {
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string
	controller.FriendInterface
	controller.BlackInterface

	userRpc pbUser.UserClient
}

func NewFriendServer(port int) *friendServer {
	log.NewPrivateLog(constant.LogFileName)
	f := friendServer{
		rpcPort:         port,
		rpcRegisterName: config.Config.RpcRegisterName.OpenImFriendName,
		etcdSchema:      config.Config.Etcd.EtcdSchema,
		etcdAddr:        config.Config.Etcd.EtcdAddr,
	}
	ttl := 10
	etcdClient, err := getcdv3.NewEtcdConn(config.Config.Etcd.EtcdSchema, strings.Join(f.etcdAddr, ","), config.Config.RpcRegisterIP, config.Config.Etcd.UserName, config.Config.Etcd.Password, port, ttl)
	if err != nil {
		panic("NewEtcdConn failed" + err.Error())
	}
	err = etcdClient.RegisterEtcd("", f.rpcRegisterName)
	if err != nil {
		panic("NewEtcdConn failed" + err.Error())
	}

	etcdClient.SetDefaultEtcdConfig(config.Config.RpcRegisterName.OpenImUserName, config.Config.RpcPort.OpenImUserPort)
	conn := etcdClient.GetConn("", config.Config.RpcRegisterName.OpenImUserName)
	f.userRpc = pbUser.NewUserClient(conn)

	//mysql init
	var mysql relation.Mysql
	var model relation.FriendGorm
	err = mysql.InitConn().AutoMigrateModel(&relationTb.FriendModel{})
	if err != nil {
		panic("db init err:" + err.Error())
	}
	err = mysql.InitConn().AutoMigrateModel(&relationTb.FriendRequestModel{})
	if err != nil {
		panic("db init err:" + err.Error())
	}

	err = mysql.InitConn().AutoMigrateModel(&relationTb.BlackModel{})
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
	pbFriend.RegisterFriendServer(srv, s)
	err = srv.Serve(listener)
	if err != nil {
		log.NewError("0", "Serve failed ", err.Error(), listener)
		return
	}
}

// ok
func (s *friendServer) ApplyToAddFriend(ctx context.Context, req *pbFriend.ApplyToAddFriendReq) (resp *pbFriend.ApplyToAddFriendResp, err error) {
	resp = &pbFriend.ApplyToAddFriendResp{}
	if err := token_verify.CheckAccessV3(ctx, req.FromUserID); err != nil {
		return nil, err
	}
	if err := callbackBeforeAddFriendV1(ctx, req); err != nil {
		return nil, err
	}
	if req.ToUserID == req.FromUserID {
		return nil, constant.ErrCanNotAddYourself.Wrap()
	}
	if _, err := check.GetUsersInfo(ctx, req.ToUserID, req.FromUserID); err != nil {
		return nil, err
	}
	in1, in2, err := s.FriendInterface.CheckIn(ctx, req.FromUserID, req.ToUserID)
	if err != nil {
		return nil, err
	}
	if in1 && in2 {
		return nil, constant.ErrRelationshipAlready.Wrap()
	}
	if err = s.FriendInterface.AddFriendRequest(ctx, req.FromUserID, req.ToUserID, req.ReqMsg, req.Ex); err != nil {
		return nil, err
	}
	chat.FriendApplicationAddNotification(ctx, req)
	return resp, nil
}

// ok
func (s *friendServer) ImportFriends(ctx context.Context, req *pbFriend.ImportFriendReq) (resp *pbFriend.ImportFriendResp, err error) {
	resp = &pbFriend.ImportFriendResp{}
	if err := token_verify.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	if _, err := check.GetUsersInfo(ctx, req.OwnerUserID, req.FriendUserIDs); err != nil {
		return nil, err
	}

	if utils.Contain(req.FriendUserIDs, req.OwnerUserID) {
		return nil, constant.ErrCanNotAddYourself.Wrap()
	}
	if utils.Duplicate(req.FriendUserIDs) {
		return nil, constant.ErrArgs.Wrap("friend userID repeated")
	}

	if err := s.FriendInterface.BecomeFriends(ctx, req.OwnerUserID, req.FriendUserIDs, constant.BecomeFriendByImport, tracelog.GetOpUserID(ctx)); err != nil {
		return nil, err
	}
	return resp, nil
}

// ok
func (s *friendServer) RespondFriendApply(ctx context.Context, req *pbFriend.RespondFriendApplyReq) (resp *pbFriend.RespondFriendApplyResp, err error) {
	resp = &pbFriend.RespondFriendApplyResp{}
	if err := check.Access(ctx, req.ToUserID); err != nil {
		return nil, err
	}
	friendRequest := relationTb.FriendRequestModel{FromUserID: req.FromUserID, ToUserID: req.ToUserID, HandleMsg: req.HandleMsg, HandleResult: req.HandleResult}
	if req.HandleResult == constant.FriendResponseAgree {
		err := s.AgreeFriendRequest(ctx, &friendRequest)
		if err != nil {
			return nil, err
		}
		chat.FriendApplicationAgreedNotification(ctx, req)
		return resp, nil
	}
	if req.HandleResult == constant.FriendResponseRefuse {
		err := s.RefuseFriendRequest(ctx, &friendRequest)
		if err != nil {
			return nil, err
		}
		chat.FriendApplicationRefusedNotification(ctx, req)
		return resp, nil
	}
	return nil, constant.ErrArgs.Wrap("req.HandleResult != -1/1")
}

// ok
func (s *friendServer) DeleteFriend(ctx context.Context, req *pbFriend.DeleteFriendReq) (resp *pbFriend.DeleteFriendResp, err error) {
	resp = &pbFriend.DeleteFriendResp{}
	if err := check.Access(ctx, req.OwnerUserID); err != nil {
		return nil, err
	}
	_, err = s.FindFriends(ctx, req.OwnerUserID, []string{req.FriendUserID})
	if err != nil {
		return nil, err
	}
	if err := s.FriendInterface.Delete(ctx, req.OwnerUserID, []string{req.FriendUserID}); err != nil {
		return nil, err
	}
	chat.FriendDeletedNotification(ctx, req)
	return resp, nil
}

// ok
func (s *friendServer) SetFriendRemark(ctx context.Context, req *pbFriend.SetFriendRemarkReq) (resp *pbFriend.SetFriendRemarkResp, err error) {
	resp = &pbFriend.SetFriendRemarkResp{}
	if err := check.Access(ctx, req.OwnerUserID); err != nil {
		return nil, err
	}
	_, err = s.FindFriends(ctx, req.OwnerUserID, []string{req.FriendUserID})
	if err != nil {
		return nil, err
	}
	if err := s.FriendInterface.UpdateRemark(ctx, req.OwnerUserID, req.FriendUserID, req.Remark); err != nil {
		return nil, err
	}
	chat.FriendRemarkSetNotification(ctx, req.OwnerUserID, req.FriendUserID)
	return resp, nil
}

// ok
func (s *friendServer) GetDesignatedFriendsReq(ctx context.Context, req *pbFriend.GetDesignatedFriendsReq) (resp *pbFriend.GetDesignatedFriendsResp, err error) {
	resp = &pbFriend.GetDesignatedFriendsResp{}
	if err := check.Access(ctx, req.UserID); err != nil {
		return nil, err
	}
	friends, total, err := s.FriendInterface.FindOwnerFriends(ctx, req.UserID, req.Pagination.PageNumber, req.Pagination.ShowNumber)
	if err != nil {
		return nil, err
	}
	resp.FriendsInfo, err = (*convert.DBFriend)(nil).DB2PB(friends)
	if err != nil {
		return nil, err
	}
	resp.Total = int32(total)
	return resp, nil
}

// ok 获取接收到的好友申请（即别人主动申请的）
func (s *friendServer) GetPaginationFriendsApplyTo(ctx context.Context, req *pbFriend.GetPaginationFriendsApplyToReq) (resp *pbFriend.GetPaginationFriendsApplyToResp, err error) {
	resp = &pbFriend.GetPaginationFriendsApplyToResp{}
	if err := check.Access(ctx, req.UserID); err != nil {
		return nil, err
	}
	friendRequests, total, err := s.FriendInterface.FindFriendRequestToMe(ctx, req.UserID, req.Pagination.PageNumber, req.Pagination.ShowNumber)
	if err != nil {
		return nil, err
	}
	resp.FriendRequests, err = (*convert.DBFriendRequest)(nil).DB2PB(friendRequests)
	if err != nil {
		return nil, err
	}
	resp.Total = int32(total)
	return resp, nil
}

// ok 获取主动发出去的好友申请列表
func (s *friendServer) GetPaginationFriendsApplyFrom(ctx context.Context, req *pbFriend.GetPaginationFriendsApplyFromReq) (resp *pbFriend.GetPaginationFriendsApplyFromResp, err error) {
	resp = &pbFriend.GetPaginationFriendsApplyFromResp{}
	if err := check.Access(ctx, req.UserID); err != nil {
		return nil, err
	}
	friendRequests, total, err := s.FriendInterface.FindFriendRequestFromMe(ctx, req.UserID, req.Pagination.PageNumber, req.Pagination.ShowNumber)
	if err != nil {
		return nil, err
	}
	resp.FriendRequests, err = (*convert.DBFriendRequest)(nil).DB2PB(friendRequests)
	if err != nil {
		return nil, err
	}
	resp.Total = int32(total)
	return resp, nil
}

// ok
func (s *friendServer) IsFriend(ctx context.Context, req *pbFriend.IsFriendReq) (resp *pbFriend.IsFriendResp, err error) {
	resp = &pbFriend.IsFriendResp{}
	resp.InUser1Friends, resp.InUser2Friends, err = s.FriendInterface.CheckIn(ctx, req.UserID1, req.UserID2)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// ok
func (s *friendServer) GetPaginationFriends(ctx context.Context, req *pbFriend.GetPaginationFriendsReq) (resp *pbFriend.GetPaginationFriendsResp, err error) {
	resp = &pbFriend.GetPaginationFriendsResp{}
	if utils.Duplicate(req.FriendUserIDs) {
		return nil, constant.ErrArgs.Wrap("friend userID repeated")
	}
	friends, err := s.FriendInterface.FindFriends(ctx, req.OwnerUserID, req.FriendUserIDs)
	if err != nil {
		return nil, err
	}
	if resp.FriendsInfo, err = (*convert.DBFriend)(nil).DB2PB(friends); err != nil {
		return nil, err
	}
	return resp, nil
}
