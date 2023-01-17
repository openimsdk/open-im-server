package friend

import (
	chat "Open_IM/internal/rpc/msg"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db/mysql"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/middleware"
	promePkg "Open_IM/pkg/common/prometheus"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/common/tools"
	"Open_IM/pkg/common/trace_log"
	cp "Open_IM/pkg/common/utils"
	"Open_IM/pkg/getcdv3"
	pbCache "Open_IM/pkg/proto/cache"
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
	friendModel        *mysql.Friend
	friendRequestModel *mysql.FriendRequest
	blackModel         *mysql.Black
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
	db := mysql.ConnectToDB()
	s.friendModel = mysql.NewFriend(db)
	s.friendRequestModel = mysql.NewFriendRequest(db)
	s.blackModel = mysql.NewBlack(db)

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
	black := imdb.Black{OwnerUserID: req.FromUserID, BlockUserID: req.ToUserID, OperatorUserID: tools.OpUserID(ctx)}
	if err := s.blackModel.Create(ctx, []*imdb.Black{&black}); err != nil {
		return nil, err
	}
	etcdConn, err := getcdv3.GetConn(ctx, config.Config.RpcRegisterName.OpenImCacheName)
	if err != nil {
		return nil, err
	}
	_, err = pbCache.NewCacheClient(etcdConn).DelBlackIDListFromCache(ctx, &pbCache.DelBlackIDListFromCacheReq{UserID: req.FromUserID})
	if err != nil {
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
	userIDList, err := rocksCache.GetFriendIDListFromCache(ctx, req.ToUserID)
	if err != nil {
		return nil, err
	}
	userIDList2, err := rocksCache.GetFriendIDListFromCache(ctx, req.FromUserID)
	if err != nil {
		return nil, err
	}
	var isSend = true
	for _, v := range userIDList {
		if v == req.FromUserID {
			for _, v2 := range userIDList2 {
				if v2 == req.ToUserID {
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
		friendRequest := imdb.FriendRequest{
			FromUserID:   req.FromUserID,
			ToUserID:     req.ToUserID,
			HandleResult: 0,
			ReqMsg:       req.ReqMsg,
			CreateTime:   time.Now(),
		}
		if err := s.friendRequestModel.Create(ctx, []*imdb.FriendRequest{&friendRequest}); err != nil {
			return nil, err
		}
		chat.FriendApplicationNotification(req)
	}
	return resp, nil
}

func (s *friendServer) ImportFriend(ctx context.Context, req *pbFriend.ImportFriendReq) (*pbFriend.ImportFriendResp, error) {
	resp := &pbFriend.ImportFriendResp{}
	if !utils.IsContain(tools.OpUserID(ctx), config.Config.Manager.AppManagerUid) {
		return nil, constant.ErrNoPermission.Wrap()
	}
	if _, err := GetUserInfo(ctx, req.FromUserID); err != nil {
		return nil, err
	}
	for _, userID := range req.FriendUserIDList {
		if _, err := GetUserInfo(ctx, userID); err != nil {
			return nil, err
		}
		fs, err := s.friendModel.FindUserState(ctx, req.FromUserID, userID)
		if err != nil {
			return nil, err
		}
		var friends []*imdb.Friend
		switch len(fs) {
		case 1:
			if fs[0].OwnerUserID == req.FromUserID {
				friends = append(friends, &imdb.Friend{OwnerUserID: userID, FriendUserID: req.FromUserID})
			} else {
				friends = append(friends, &imdb.Friend{OwnerUserID: req.FromUserID, FriendUserID: userID})
			}
		case 0:
			friends = append(friends, &imdb.Friend{OwnerUserID: userID, FriendUserID: req.FromUserID}, &imdb.Friend{OwnerUserID: req.FromUserID, FriendUserID: userID})
		default:
			continue
		}
		if err := s.friendModel.Create(ctx, friends); err != nil {
			return nil, err
		}
	}
	etcdConn, err := getcdv3.GetConn(ctx, config.Config.RpcRegisterName.OpenImCacheName)
	if err != nil {
		return nil, err
	}
	cacheClient := pbCache.NewCacheClient(etcdConn)
	if _, err := cacheClient.DelFriendIDListFromCache(ctx, &pbCache.DelFriendIDListFromCacheReq{UserID: req.FromUserID}); err != nil {
		return nil, err
	}
	if err := rocksCache.DelAllFriendsInfoFromCache(ctx, req.FromUserID); err != nil {
		trace_log.SetCtxInfo(ctx, "DelAllFriendsInfoFromCache", err, "userID", req.FromUserID)
	}
	for _, userID := range req.FriendUserIDList {
		if _, err = cacheClient.DelFriendIDListFromCache(ctx, &pbCache.DelFriendIDListFromCacheReq{UserID: userID}); err != nil {
			return nil, err
		}
		if err := rocksCache.DelAllFriendsInfoFromCache(ctx, userID); err != nil {
			trace_log.SetCtxInfo(ctx, "DelAllFriendsInfoFromCache", err, "userID", userID)
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
	err = imdb.UpdateFriendApplication(friendRequest)
	if err != nil {
		return nil, err
	}

	//Change the status of the friend request form
	if req.HandleResult == constant.FriendFlag {
		var isInsert bool
		//Establish friendship after find friend relationship not exists
		_, err := s.friendModel.Take(ctx, req.FromUserID, req.ToUserID)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if err := s.friendModel.Create(ctx, []*imdb.Friend{{OwnerUserID: req.FromUserID, FriendUserID: req.ToUserID, OperatorUserID: tools.OpUserID(ctx)}}); err != nil {
				return nil, err
			}
			isInsert = true
		} else if err != nil {
			return nil, err
		}

		// cache rpc
		if isInsert {
			etcdConn, err := getcdv3.GetConn(ctx, config.Config.RpcRegisterName.OpenImCacheName)
			if err != nil {
				return nil, err
			}
			client := pbCache.NewCacheClient(etcdConn)

			if _, err := client.DelFriendIDListFromCache(context.Background(), &pbCache.DelFriendIDListFromCacheReq{UserID: req.ToUserID}); err != nil {
				return nil, err
			}
			if _, err := client.DelFriendIDListFromCache(context.Background(), &pbCache.DelFriendIDListFromCacheReq{UserID: req.FromUserID}); err != nil {
				return nil, err
			}
			if err := rocksCache.DelAllFriendsInfoFromCache(ctx, req.ToUserID); err != nil {
				trace_log.SetCtxInfo(ctx, "DelAllFriendsInfoFromCache", err, "userID", req.ToUserID)
			}
			if err := rocksCache.DelAllFriendsInfoFromCache(ctx, req.FromUserID); err != nil {
				trace_log.SetCtxInfo(ctx, "DelAllFriendsInfoFromCache", err, "userID", req.FromUserID)
			}
			chat.FriendAddedNotification(tools.OperationID(ctx), tools.OpUserID(ctx), req.FromUserID, req.ToUserID)
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
	etcdConn, err := getcdv3.GetConn(ctx, config.Config.RpcRegisterName.OpenImCacheName)
	if err != nil {
		return nil, err
	}
	_, err = pbCache.NewCacheClient(etcdConn).DelFriendIDListFromCache(context.Background(), &pbCache.DelFriendIDListFromCacheReq{UserID: req.FromUserID})
	if err != nil {
		return nil, err
	}
	if err := rocksCache.DelAllFriendsInfoFromCache(ctx, req.FromUserID); err != nil {
		trace_log.SetCtxInfo(ctx, "DelAllFriendsInfoFromCache", err, "DelAllFriendsInfoFromCache", req.FromUserID)
	}
	if err := rocksCache.DelAllFriendsInfoFromCache(ctx, req.ToUserID); err != nil {
		trace_log.SetCtxInfo(ctx, "DelAllFriendsInfoFromCache", err, "DelAllFriendsInfoFromCache", req.ToUserID)
	}
	chat.FriendDeletedNotification(req)
	return resp, nil
}

func (s *friendServer) GetBlacklist(ctx context.Context, req *pbFriend.GetBlacklistReq) (*pbFriend.GetBlacklistResp, error) {
	resp := &pbFriend.GetBlacklistResp{}
	if err := token_verify.CheckAccessV3(ctx, req.FromUserID); err != nil {
		return nil, err
	}
	blackIDList, err := rocksCache.GetBlackListFromCache(ctx, req.FromUserID)
	if err != nil {
		return nil, err
	}
	for _, userID := range blackIDList {
		user, err := rocksCache.GetUserInfoFromCache(ctx, userID)
		if err != nil {
			trace_log.SetCtxInfo(ctx, "GetUserInfoFromCache", err, "userID", userID)
			continue
		}
		var blackUserInfo sdkws.PublicUserInfo
		utils.CopyStructFields(&blackUserInfo, user)
		resp.BlackUserInfoList = append(resp.BlackUserInfoList, &blackUserInfo)
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
	if err := rocksCache.DelAllFriendsInfoFromCache(ctx, req.FromUserID); err != nil {
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
	if err := s.blackModel.Delete(ctx, []*imdb.Black{{OwnerUserID: req.FromUserID, BlockUserID: req.ToUserID}}); err != nil {
		return nil, err
	}
	etcdConn, err := getcdv3.GetConn(ctx, config.Config.RpcRegisterName.OpenImCacheName)
	if err != nil {
		return nil, err
	}
	_, err = pbCache.NewCacheClient(etcdConn).DelBlackIDListFromCache(context.Background(), &pbCache.DelBlackIDListFromCacheReq{UserID: req.FromUserID})
	if err != nil {
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
	blackIDList, err := rocksCache.GetBlackListFromCache(ctx, req.FromUserID)
	if err != nil {
		return nil, err
	}
	resp.Response = utils.IsContain(req.ToUserID, blackIDList)
	return resp, nil
}

func (s *friendServer) IsFriend(ctx context.Context, req *pbFriend.IsFriendReq) (*pbFriend.IsFriendResp, error) {
	resp := &pbFriend.IsFriendResp{}
	if err := token_verify.CheckAccessV3(ctx, req.FromUserID); err != nil {
		return nil, err
	}
	friendIDList, err := rocksCache.GetFriendIDListFromCache(ctx, req.FromUserID)
	if err != nil {
		return nil, err
	}
	resp.Response = utils.IsContain(req.ToUserID, friendIDList)
	return resp, nil
}

func (s *friendServer) GetFriendList(ctx context.Context, req *pbFriend.GetFriendListReq) (*pbFriend.GetFriendListResp, error) {
	resp := &pbFriend.GetFriendListResp{}
	if err := token_verify.CheckAccessV3(ctx, req.FromUserID); err != nil {
		return nil, err
	}
	friendList, err := rocksCache.GetAllFriendsInfoFromCache(ctx, req.FromUserID)
	if err != nil {
		return nil, err
	}
	var userInfoList []*sdkws.FriendInfo
	for _, friendUser := range friendList {
		friendUserInfo := sdkws.FriendInfo{FriendUser: &sdkws.UserInfo{}}
		cp.FriendDBCopyOpenIM(&friendUserInfo, friendUser)
		userInfoList = append(userInfoList, &friendUserInfo)
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
	applyUsersInfo, err := s.friendRequestModel.FindToUserID(ctx, req.FromUserID)
	if err != nil {
		return nil, err
	}
	for _, applyUserInfo := range applyUsersInfo {
		var userInfo sdkws.FriendRequest
		cp.FriendRequestDBCopyOpenIM(&userInfo, applyUserInfo)
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
	usersInfo, err := s.friendRequestModel.FindFromUserID(ctx, req.FromUserID)
	if err != nil {
		return nil, err
	}
	for _, selfApplyOtherUserInfo := range usersInfo {
		var userInfo sdkws.FriendRequest // pbFriend.ApplyUserInfo
		cp.FriendRequestDBCopyOpenIM(&userInfo, selfApplyOtherUserInfo)
		resp.FriendRequestList = append(resp.FriendRequestList, &userInfo)
	}
	return resp, nil
}
