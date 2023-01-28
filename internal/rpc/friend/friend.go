package friend

import (
	chat "Open_IM/internal/rpc/msg"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db/model"
	"Open_IM/pkg/common/db/mysql"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/middleware"
	promePkg "Open_IM/pkg/common/prometheus"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/common/tools"
	"Open_IM/pkg/common/trace_log"
	"Open_IM/pkg/getcdv3"
	pbFriend "Open_IM/pkg/proto/friend"
	sdkws "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
	"errors"
	"gorm.io/gorm"
	"net"
	"strconv"
	"strings"
	"time"

	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"

	"google.golang.org/grpc"
)

type friendServer struct {
	rpcPort            int
	rpcRegisterName    string
	etcdSchema         string
	etcdAddr           []string
	friendModel        *controller.FriendModel
	friendRequestModel *controller.FriendRequestModel
	blackModel         *controller.BlackModel
}

func NewFriendServer(port int) *friendServer {
	log.NewPrivateLog(constant.LogFileName)
	return &friendServer{
		rpcPort:         port,
		rpcRegisterName: config.Config.RpcRegisterName.OpenImFriendName,
		etcdSchema:      config.Config.Etcd.EtcdSchema,
		etcdAddr:        config.Config.Etcd.EtcdAddr,
	}
}

func (s *friendServer) Run() {
	db := relation.ConnectToDB()
	//s.friendModel = mysql.NewFriend(db)
	//s.friendRequestModel = mysql.NewFriendRequest(db)
	//s.blackModel = mysql.NewBlack(db)

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
	if err := s.blackModel.Create(ctx, []*relation.Black{&black}); err != nil {
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
	friends1, err := s.friendModel.FindOwnerUserID(ctx, req.ToUserID)
	if err != nil {
		return nil, err
	}
	friends2, err := s.friendModel.FindOwnerUserID(ctx, req.FromUserID)
	if err != nil {
		return nil, err
	}
	var isSend = true
	for _, v1 := range friends1 {
		if v1.FriendUserID == req.FromUserID {
			for _, v2 := range friends2 {
				if v2.FriendUserID == req.ToUserID {
					isSend = false
					break
				}
			}
			break
		}
	}
	//Cannot add non-existent users
	if isSend {
		if _, err := GetUserInfo(ctx, req.ToUserID); err != nil {
			return nil, err
		}
		friendRequest := relation.FriendRequest{
			FromUserID:   req.FromUserID,
			ToUserID:     req.ToUserID,
			HandleResult: 0,
			ReqMsg:       req.ReqMsg,
			CreateTime:   time.Now(),
		}
		if err := s.friendRequestModel.Create(ctx, []*relation.FriendRequest{&friendRequest}); err != nil {
			return nil, err
		}
		chat.FriendApplicationNotification(req)
	}
	return resp, nil
}

func (s *friendServer) ImportFriend(ctx context.Context, req *pbFriend.ImportFriendReq) (*pbFriend.ImportFriendResp, error) {
	resp := &pbFriend.ImportFriendResp{}
	if err := token_verify.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	if _, err := GetUserInfo(ctx, req.FromUserID); err != nil {
		return nil, err
	}

	var friends []*relation.Friend
	for _, userID := range utils.RemoveDuplicateElement(req.FriendUserIDList) {
		if _, err := GetUserInfo(ctx, userID); err != nil {
			return nil, err
		}
		fs, err := s.friendModel.FindUserState(ctx, req.FromUserID, userID)
		if err != nil {
			return nil, err
		}
		switch len(fs) {
		case 1:
			if fs[0].OwnerUserID == req.FromUserID {
				friends = append(friends, &relation.Friend{OwnerUserID: userID, FriendUserID: req.FromUserID})
			} else {
				friends = append(friends, &relation.Friend{OwnerUserID: req.FromUserID, FriendUserID: userID})
			}
		case 0:
			friends = append(friends, &relation.Friend{OwnerUserID: userID, FriendUserID: req.FromUserID}, &relation.Friend{OwnerUserID: req.FromUserID, FriendUserID: userID})
		default:
			continue
		}
	}
	if len(friends) > 0 {
		if err := s.friendModel.Create(ctx, friends); err != nil {
			return nil, err
		}
	}
	return resp, nil
}

// process Friend application
func (s *friendServer) AddFriendResponse(ctx context.Context, req *pbFriend.AddFriendResponseReq) (*pbFriend.AddFriendResponseResp, error) {
	resp := &pbFriend.AddFriendResponseResp{}
	if err := token_verify.CheckAccessV3(ctx, req.FromUserID); err != nil {
		return nil, err
	}
	friendRequest, err := s.friendRequestModel.Take(ctx, req.ToUserID, req.FromUserID)
	if err != nil {
		return nil, err
	}
	friendRequest.HandleResult = req.HandleResult
	friendRequest.HandleTime = time.Now()
	friendRequest.HandleMsg = req.HandleMsg
	friendRequest.HandlerUserID = tools.OpUserID(ctx)
	err = relation.UpdateFriendApplication(friendRequest)
	if err != nil {
		return nil, err
	}

	//Change the status of the friend request form
	if req.HandleResult == constant.FriendFlag {
		//Establish friendship after find friend relationship not exists
		_, err := s.friendModel.Take(ctx, req.FromUserID, req.ToUserID)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := s.friendModel.Create(ctx, []*relation.Friend{{OwnerUserID: req.FromUserID, FriendUserID: req.ToUserID, OperatorUserID: tools.OpUserID(ctx)}}); err != nil {
				return nil, err
			}
			chat.FriendAddedNotification(tools.OperationID(ctx), tools.OpUserID(ctx), req.FromUserID, req.ToUserID)
		} else if err != nil {
			return nil, err
		}
	}

	if req.HandleResult == constant.FriendResponseAgree {
		chat.FriendApplicationApprovedNotification(req)
	} else if req.HandleResult == constant.FriendResponseRefuse {
		chat.FriendApplicationRejectedNotification(req)
	} else {
		trace_log.SetCtxInfo(ctx, utils.GetSelfFuncName(), nil, "handleResult", req.HandleResult)
	}
	return resp, nil
}

func (s *friendServer) DeleteFriend(ctx context.Context, req *pbFriend.DeleteFriendReq) (*pbFriend.DeleteFriendResp, error) {
	resp := &pbFriend.DeleteFriendResp{}
	if err := token_verify.CheckAccessV3(ctx, req.FromUserID); err != nil {
		return nil, err
	}
	if err := s.friendModel.Delete(ctx, req.FromUserID, req.ToUserID); err != nil {
		return nil, err
	}
	chat.FriendDeletedNotification(req)
	return resp, nil
}

func (s *friendServer) GetBlacklist(ctx context.Context, req *pbFriend.GetBlacklistReq) (*pbFriend.GetBlacklistResp, error) {
	resp := &pbFriend.GetBlacklistResp{}
	if err := token_verify.CheckAccessV3(ctx, req.FromUserID); err != nil {
		return nil, err
	}
	blacks, err := s.blackModel.FindByOwnerUserID(ctx, req.FromUserID)
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

func (s *friendServer) SetFriendRemark(ctx context.Context, req *pbFriend.SetFriendRemarkReq) (*pbFriend.SetFriendRemarkResp, error) {
	resp := &pbFriend.SetFriendRemarkResp{}
	if err := token_verify.CheckAccessV3(ctx, req.FromUserID); err != nil {
		return nil, err
	}
	if err := s.friendModel.UpdateRemark(ctx, req.FromUserID, req.ToUserID, req.Remark); err != nil {
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
	if err := s.blackModel.Delete(ctx, []*relation.Black{{OwnerUserID: req.FromUserID, BlockUserID: req.ToUserID}}); err != nil {
		return nil, err
	}
	chat.BlackDeletedNotification(req)
	return resp, nil
}

func (s *friendServer) IsInBlackList(ctx context.Context, req *pbFriend.IsInBlackListReq) (*pbFriend.IsInBlackListResp, error) {
	resp := &pbFriend.IsInBlackListResp{}
	if err := token_verify.CheckAccessV3(ctx, req.FromUserID); err != nil {
		return nil, err
	}
	exist, err := s.blackModel.IsExist(ctx, req.FromUserID, req.ToUserID)
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
	exist, err := s.friendModel.IsExist(ctx, req.FromUserID, req.ToUserID)
	if err != nil {
		return nil, err
	}
	resp.Response = exist
	return resp, nil
}

func (s *friendServer) GetFriendList(ctx context.Context, req *pbFriend.GetFriendListReq) (*pbFriend.GetFriendListResp, error) {
	resp := &pbFriend.GetFriendListResp{}
	if err := token_verify.CheckAccessV3(ctx, req.FromUserID); err != nil {
		return nil, err
	}
	friends, err := s.friendModel.FindOwnerUserID(ctx, req.FromUserID)
	if err != nil {
		return nil, err
	}
	userIDList := make([]string, 0, len(friends))
	for _, f := range friends {
		userIDList = append(userIDList, f.FriendUserID)
	}
	users, err := GetUserInfoList(ctx, userIDList)
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
	friendRequests, err := s.friendRequestModel.FindFromUserID(ctx, req.FromUserID)
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
