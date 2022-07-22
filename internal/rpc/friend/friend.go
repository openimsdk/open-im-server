package friend

import (
	chat "Open_IM/internal/rpc/msg"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	imdb "Open_IM/pkg/common/db/mysql_model/im_mysql_model"
	"Open_IM/pkg/common/db/rocks_cache"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/token_verify"
	cp "Open_IM/pkg/common/utils"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbCache "Open_IM/pkg/proto/cache"
	pbFriend "Open_IM/pkg/proto/friend"
	sdkws "Open_IM/pkg/proto/sdk_ws"
	"Open_IM/pkg/utils"
	"context"
	"net"
	"strconv"
	"strings"
	"time"

	"google.golang.org/grpc"
)

type friendServer struct {
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string
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
	srv := grpc.NewServer()
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
	err = getcdv3.RegisterEtcd(s.etcdSchema, strings.Join(s.etcdAddr, ","), rpcRegisterIP, s.rpcPort, s.rpcRegisterName, 10)
	if err != nil {
		log.NewError("0", "RegisterEtcd failed ", err.Error(), s.etcdSchema, strings.Join(s.etcdAddr, ","), rpcRegisterIP, s.rpcPort, s.rpcRegisterName)
		return
	}
	err = srv.Serve(listener)
	if err != nil {
		log.NewError("0", "Serve failed ", err.Error(), listener)
		return
	}
}

func (s *friendServer) AddBlacklist(ctx context.Context, req *pbFriend.AddBlacklistReq) (*pbFriend.AddBlacklistResp, error) {
	log.NewInfo(req.CommID.OperationID, "AddBlacklist args ", req.String())
	ok := token_verify.CheckAccess(req.CommID.OpUserID, req.CommID.FromUserID)
	if !ok {
		log.NewError(req.CommID.OperationID, "CheckAccess false ", req.CommID.OpUserID, req.CommID.FromUserID)
		return &pbFriend.AddBlacklistResp{CommonResp: &pbFriend.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}}, nil
	}
	black := db.Black{OwnerUserID: req.CommID.FromUserID, BlockUserID: req.CommID.ToUserID, OperatorUserID: req.CommID.OpUserID}

	err := imdb.InsertInToUserBlackList(black)
	if err != nil {
		log.NewError(req.CommID.OperationID, "InsertInToUserBlackList failed ", err.Error())
		return &pbFriend.AddBlacklistResp{CommonResp: &pbFriend.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: err.Error()}}, nil
	}
	log.NewInfo(req.CommID.OperationID, "AddBlacklist rpc ok ", req.CommID.FromUserID, req.CommID.ToUserID)

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImCacheName, req.CommID.OperationID)
	if etcdConn == nil {
		errMsg := req.CommID.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.CommID.OperationID, errMsg)
		return &pbFriend.AddBlacklistResp{CommonResp: &pbFriend.CommonResp{ErrCode: constant.ErrInternal.ErrCode, ErrMsg: errMsg}}, nil
	}
	cacheClient := pbCache.NewCacheClient(etcdConn)
	cacheResp, err := cacheClient.DelBlackIDListFromCache(context.Background(), &pbCache.DelBlackIDListFromCacheReq{UserID: req.CommID.FromUserID, OperationID: req.CommID.OperationID})
	if err != nil {
		log.NewError(req.CommID.OperationID, "DelBlackIDListFromCache rpc call failed ", err.Error())
		return &pbFriend.AddBlacklistResp{CommonResp: &pbFriend.CommonResp{ErrCode: 500, ErrMsg: "AddBlackUserToCache rpc call failed"}}, nil
	}
	if cacheResp.CommonResp.ErrCode != 0 {
		log.NewError(req.CommID.OperationID, "DelBlackIDListFromCache rpc logic call failed ", cacheResp.String())
		return &pbFriend.AddBlacklistResp{CommonResp: &pbFriend.CommonResp{ErrCode: cacheResp.CommonResp.ErrCode, ErrMsg: cacheResp.CommonResp.ErrMsg}}, nil
	}

	chat.BlackAddedNotification(req)
	return &pbFriend.AddBlacklistResp{CommonResp: &pbFriend.CommonResp{}}, nil
}

func (s *friendServer) AddFriend(ctx context.Context, req *pbFriend.AddFriendReq) (*pbFriend.AddFriendResp, error) {
	log.NewInfo(req.CommID.OperationID, "AddFriend args ", req.String())
	ok := token_verify.CheckAccess(req.CommID.OpUserID, req.CommID.FromUserID)
	if !ok {
		log.NewError(req.CommID.OperationID, "CheckAccess false ", req.CommID.OpUserID, req.CommID.FromUserID)
		return &pbFriend.AddFriendResp{CommonResp: &pbFriend.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}}, nil
	}
	//Cannot add non-existent users
	if _, err := imdb.GetUserByUserID(req.CommID.ToUserID); err != nil {
		log.NewError(req.CommID.OperationID, "GetUserByUserID failed ", err.Error(), req.CommID.ToUserID)
		return &pbFriend.AddFriendResp{CommonResp: &pbFriend.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}

	//Establish a latest relationship in the friend request table
	friendRequest := db.FriendRequest{
		HandleResult: 0, ReqMsg: req.ReqMsg, CreateTime: time.Now()}
	utils.CopyStructFields(&friendRequest, req.CommID)
	// {openIM001 openIM002 0 test add friend 0001-01-01 00:00:00 +0000 UTC   0001-01-01 00:00:00 +0000 UTC }]
	log.NewDebug(req.CommID.OperationID, "UpdateFriendApplication args ", friendRequest)
	//err := imdb.InsertFriendApplication(&friendRequest)
	err := imdb.InsertFriendApplication(&friendRequest,
		map[string]interface{}{"handle_result": 0, "req_msg": friendRequest.ReqMsg, "create_time": friendRequest.CreateTime,
			"handler_user_id": "", "handle_msg": "", "handle_time": utils.UnixSecondToTime(0), "ex": ""})
	if err != nil {
		log.NewError(req.CommID.OperationID, "UpdateFriendApplication failed ", err.Error(), friendRequest)
		return &pbFriend.AddFriendResp{CommonResp: &pbFriend.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}

	chat.FriendApplicationNotification(req)
	return &pbFriend.AddFriendResp{CommonResp: &pbFriend.CommonResp{}}, nil
}

func (s *friendServer) ImportFriend(ctx context.Context, req *pbFriend.ImportFriendReq) (*pbFriend.ImportFriendResp, error) {
	log.NewInfo(req.OperationID, "ImportFriend args ", req.String())
	resp := pbFriend.ImportFriendResp{CommonResp: &pbFriend.CommonResp{}}
	var c pbFriend.CommonResp

	if !utils.IsContain(req.OpUserID, config.Config.Manager.AppManagerUid) {
		log.NewError(req.OperationID, "not authorized", req.OpUserID, config.Config.Manager.AppManagerUid)
		c.ErrCode = constant.ErrAccess.ErrCode
		c.ErrMsg = constant.ErrAccess.ErrMsg
		for _, v := range req.FriendUserIDList {
			resp.UserIDResultList = append(resp.UserIDResultList, &pbFriend.UserIDResult{UserID: v, Result: -1})
		}
		resp.CommonResp = &c
		return &resp, nil
	}
	if _, err := imdb.GetUserByUserID(req.FromUserID); err != nil {
		log.NewError(req.OperationID, "GetUserByUserID failed ", err.Error(), req.FromUserID)
		c.ErrCode = constant.ErrDB.ErrCode
		c.ErrMsg = "this user not exists,cant not add friend"
		for _, v := range req.FriendUserIDList {
			resp.UserIDResultList = append(resp.UserIDResultList, &pbFriend.UserIDResult{UserID: v, Result: -1})
		}
		resp.CommonResp = &c
		return &resp, nil
	}

	for _, v := range req.FriendUserIDList {
		log.NewDebug(req.OperationID, "FriendUserIDList ", v)
		if _, fErr := imdb.GetUserByUserID(v); fErr != nil {
			log.NewError(req.OperationID, "GetUserByUserID failed ", fErr.Error(), v)
			resp.UserIDResultList = append(resp.UserIDResultList, &pbFriend.UserIDResult{UserID: v, Result: -1})
		} else {
			if _, err := imdb.GetFriendRelationshipFromFriend(req.FromUserID, v); err != nil {
				//Establish two single friendship
				toInsertFollow := db.Friend{OwnerUserID: req.FromUserID, FriendUserID: v}
				err1 := imdb.InsertToFriend(&toInsertFollow)
				if err1 != nil {
					log.NewError(req.OperationID, "InsertToFriend failed ", err1.Error(), toInsertFollow)
					resp.UserIDResultList = append(resp.UserIDResultList, &pbFriend.UserIDResult{UserID: v, Result: -1})
					continue
				}
				toInsertFollow = db.Friend{OwnerUserID: v, FriendUserID: req.FromUserID}
				err2 := imdb.InsertToFriend(&toInsertFollow)
				if err2 != nil {
					log.NewError(req.OperationID, "InsertToFriend failed ", err2.Error(), toInsertFollow)
					resp.UserIDResultList = append(resp.UserIDResultList, &pbFriend.UserIDResult{UserID: v, Result: -1})
					continue
				}
				resp.UserIDResultList = append(resp.UserIDResultList, &pbFriend.UserIDResult{UserID: v, Result: 0})
				log.NewDebug(req.OperationID, "UserIDResultList ", resp.UserIDResultList)
				chat.FriendAddedNotification(req.OperationID, req.OpUserID, req.FromUserID, v)
			} else {
				log.NewWarn(req.OperationID, "GetFriendRelationshipFromFriend ok", req.FromUserID, v)
				resp.UserIDResultList = append(resp.UserIDResultList, &pbFriend.UserIDResult{UserID: v, Result: 0})
			}
		}
	}

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImCacheName, req.OperationID)
	if etcdConn == nil {
		errMsg := req.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.OperationID, errMsg)
		resp.CommonResp.ErrMsg = errMsg
		resp.CommonResp.ErrCode = 500
		return &resp, nil
	}
	cacheClient := pbCache.NewCacheClient(etcdConn)
	cacheResp, err := cacheClient.DelFriendIDListFromCache(context.Background(), &pbCache.DelFriendIDListFromCacheReq{UserID: req.FromUserID, OperationID: req.OperationID})
	if err != nil {
		log.NewError(req.OperationID, "DelBlackIDListFromCache rpc call failed ", err.Error())
		resp.CommonResp.ErrCode = 500
		resp.CommonResp.ErrMsg = err.Error()
		return &resp, nil
	}
	if cacheResp.CommonResp.ErrCode != 0 {
		log.NewError(req.OperationID, "DelBlackIDListFromCache rpc logic call failed ", cacheResp.String())
		resp.CommonResp.ErrCode = 500
		resp.CommonResp.ErrMsg = cacheResp.CommonResp.ErrMsg
		return &resp, nil
	}
	if err := rocksCache.DelAllFriendsInfoFromCache(req.FromUserID); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), req.FromUserID)
	}

	for _, userID := range req.FriendUserIDList {
		cacheResp, err := cacheClient.DelFriendIDListFromCache(context.Background(), &pbCache.DelFriendIDListFromCacheReq{UserID: userID, OperationID: req.OperationID})
		if err != nil {
			log.NewError(req.OperationID, "DelBlackIDListFromCache rpc call failed ", err.Error())
		}
		if cacheResp != nil && cacheResp.CommonResp != nil {
			if cacheResp.CommonResp.ErrCode != 0 {
				log.NewError(req.OperationID, "DelBlackIDListFromCache rpc logic call failed ", cacheResp.String())
				resp.CommonResp.ErrCode = 500
				resp.CommonResp.ErrMsg = cacheResp.CommonResp.ErrMsg
				return &resp, nil
			}
		}
		if err := rocksCache.DelAllFriendsInfoFromCache(userID); err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), err.Error(), userID)
		}
	}

	resp.CommonResp.ErrCode = 0
	log.NewInfo(req.OperationID, "ImportFriend rpc ok ", resp.String())
	return &resp, nil
}

//process Friend application
func (s *friendServer) AddFriendResponse(ctx context.Context, req *pbFriend.AddFriendResponseReq) (*pbFriend.AddFriendResponseResp, error) {
	log.NewInfo(req.CommID.OperationID, "AddFriendResponse args ", req.String())

	if !token_verify.CheckAccess(req.CommID.OpUserID, req.CommID.FromUserID) {
		log.NewError(req.CommID.OperationID, "CheckAccess false ", req.CommID.OpUserID, req.CommID.FromUserID)
		return &pbFriend.AddFriendResponseResp{CommonResp: &pbFriend.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}}, nil
	}

	//Check there application before agreeing or refuse to a friend's application
	//req.CommID.FromUserID process req.CommID.ToUserID
	friendRequest, err := imdb.GetFriendApplicationByBothUserID(req.CommID.ToUserID, req.CommID.FromUserID)
	if err != nil {
		log.NewError(req.CommID.OperationID, "GetFriendApplicationByBothUserID failed ", err.Error(), req.CommID.ToUserID, req.CommID.FromUserID)
		return &pbFriend.AddFriendResponseResp{CommonResp: &pbFriend.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}
	friendRequest.HandleResult = req.HandleResult
	friendRequest.HandleTime = time.Now()
	//friendRequest.HandleTime.Unix()
	friendRequest.HandleMsg = req.HandleMsg
	friendRequest.HandlerUserID = req.CommID.OpUserID
	err = imdb.UpdateFriendApplication(friendRequest)
	if err != nil {
		log.NewError(req.CommID.OperationID, "UpdateFriendApplication failed ", err.Error(), friendRequest)
		return &pbFriend.AddFriendResponseResp{CommonResp: &pbFriend.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}

	//Change the status of the friend request form
	if req.HandleResult == constant.FriendFlag {
		//Establish friendship after find friend relationship not exists
		_, err := imdb.GetFriendRelationshipFromFriend(req.CommID.FromUserID, req.CommID.ToUserID)
		if err == nil {
			log.NewWarn(req.CommID.OperationID, "GetFriendRelationshipFromFriend exist", req.CommID.FromUserID, req.CommID.ToUserID)
		} else {
			//Establish two single friendship
			toInsertFollow := db.Friend{OwnerUserID: req.CommID.FromUserID, FriendUserID: req.CommID.ToUserID, OperatorUserID: req.CommID.OpUserID}
			err = imdb.InsertToFriend(&toInsertFollow)
			if err != nil {
				log.NewError(req.CommID.OperationID, "InsertToFriend failed ", err.Error(), toInsertFollow)
				return &pbFriend.AddFriendResponseResp{CommonResp: &pbFriend.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
			}
		}

		_, err = imdb.GetFriendRelationshipFromFriend(req.CommID.ToUserID, req.CommID.FromUserID)
		if err == nil {
			log.NewWarn(req.CommID.OperationID, "GetFriendRelationshipFromFriend exist", req.CommID.ToUserID, req.CommID.FromUserID)
		} else {
			toInsertFollow := db.Friend{OwnerUserID: req.CommID.ToUserID, FriendUserID: req.CommID.FromUserID, OperatorUserID: req.CommID.OpUserID}
			err = imdb.InsertToFriend(&toInsertFollow)
			if err != nil {
				log.NewError(req.CommID.OperationID, "InsertToFriend failed ", err.Error(), toInsertFollow)
				return &pbFriend.AddFriendResponseResp{CommonResp: &pbFriend.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
			}
			// cache rpc
			delFriendIDListFromCacheReq := &pbCache.DelFriendIDListFromCacheReq{OperationID: req.CommID.OperationID}
			etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImCacheName, req.CommID.OperationID)
			if etcdConn == nil {
				errMsg := req.CommID.OperationID + "getcdv3.GetConn == nil"
				log.NewError(req.CommID.OperationID, errMsg)
				return &pbFriend.AddFriendResponseResp{CommonResp: &pbFriend.CommonResp{ErrCode: constant.ErrInternal.ErrCode, ErrMsg: errMsg}}, nil
			}
			client := pbCache.NewCacheClient(etcdConn)
			delFriendIDListFromCacheReq.UserID = req.CommID.ToUserID
			respPb, err := client.DelFriendIDListFromCache(context.Background(), delFriendIDListFromCacheReq)
			if err != nil {
				log.NewError(req.CommID.OperationID, utils.GetSelfFuncName(), "DelFriendIDListFromCache failed", err.Error())
				return &pbFriend.AddFriendResponseResp{CommonResp: &pbFriend.CommonResp{ErrCode: constant.ErrServer.ErrCode, ErrMsg: constant.ErrServer.ErrMsg}}, nil
			}
			if respPb.CommonResp.ErrCode != 0 {
				log.NewError(req.CommID.OperationID, utils.GetSelfFuncName(), "DelFriendIDListFromCache failed", respPb.CommonResp.String())
				return &pbFriend.AddFriendResponseResp{CommonResp: &pbFriend.CommonResp{ErrCode: respPb.CommonResp.ErrCode, ErrMsg: respPb.CommonResp.ErrMsg}}, nil
			}
			delFriendIDListFromCacheReq.UserID = req.CommID.FromUserID
			respPb, err = client.DelFriendIDListFromCache(context.Background(), delFriendIDListFromCacheReq)
			if err != nil {
				log.NewError(req.CommID.OperationID, utils.GetSelfFuncName(), "DelFriendIDListFromCache failed", err.Error())
				return &pbFriend.AddFriendResponseResp{CommonResp: &pbFriend.CommonResp{ErrCode: constant.ErrServer.ErrCode, ErrMsg: constant.ErrServer.ErrMsg}}, nil
			}
			if respPb.CommonResp.ErrCode != 0 {
				log.NewError(req.CommID.OperationID, utils.GetSelfFuncName(), "DelFriendIDListFromCache failed", respPb.CommonResp.String())
				return &pbFriend.AddFriendResponseResp{CommonResp: &pbFriend.CommonResp{ErrCode: respPb.CommonResp.ErrCode, ErrMsg: respPb.CommonResp.ErrMsg}}, nil
			}
			if err := rocksCache.DelAllFriendsInfoFromCache(req.CommID.ToUserID); err != nil {
				log.NewError(req.CommID.OperationID, utils.GetSelfFuncName(), err.Error(), req.CommID.ToUserID)
			}
			if err := rocksCache.DelAllFriendsInfoFromCache(req.CommID.FromUserID); err != nil {
				log.NewError(req.CommID.OperationID, utils.GetSelfFuncName(), err.Error(), req.CommID.FromUserID)
			}
			chat.FriendAddedNotification(req.CommID.OperationID, req.CommID.OpUserID, req.CommID.FromUserID, req.CommID.ToUserID)
		}
	}

	if req.HandleResult == constant.FriendResponseAgree {
		chat.FriendApplicationApprovedNotification(req)
	} else if req.HandleResult == constant.FriendResponseRefuse {
		chat.FriendApplicationRejectedNotification(req)
	} else {
		log.Error(req.CommID.OperationID, "HandleResult failed ", req.HandleResult)
	}

	log.NewInfo(req.CommID.OperationID, "rpc AddFriendResponse ok")
	return &pbFriend.AddFriendResponseResp{CommonResp: &pbFriend.CommonResp{}}, nil
}

func (s *friendServer) DeleteFriend(ctx context.Context, req *pbFriend.DeleteFriendReq) (*pbFriend.DeleteFriendResp, error) {
	log.NewInfo(req.CommID.OperationID, "DeleteFriend args ", req.String())
	//Parse token, to find current user information
	if !token_verify.CheckAccess(req.CommID.OpUserID, req.CommID.FromUserID) {
		log.NewError(req.CommID.OperationID, "CheckAccess false ", req.CommID.OpUserID, req.CommID.FromUserID)
		return &pbFriend.DeleteFriendResp{CommonResp: &pbFriend.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}}, nil
	}
	err := imdb.DeleteSingleFriendInfo(req.CommID.FromUserID, req.CommID.ToUserID)
	if err != nil {
		log.NewError(req.CommID.OperationID, "DeleteSingleFriendInfo failed", err.Error(), req.CommID.FromUserID, req.CommID.ToUserID)
		return &pbFriend.DeleteFriendResp{CommonResp: &pbFriend.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}}, nil
	}
	log.NewInfo(req.CommID.OperationID, "DeleteFriend rpc ok")

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImCacheName, req.CommID.OperationID)
	if etcdConn == nil {
		errMsg := req.CommID.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.CommID.OperationID, errMsg)
		return &pbFriend.DeleteFriendResp{CommonResp: &pbFriend.CommonResp{ErrCode: constant.ErrInternal.ErrCode, ErrMsg: errMsg}}, nil
	}
	client := pbCache.NewCacheClient(etcdConn)
	respPb, err := client.DelFriendIDListFromCache(context.Background(), &pbCache.DelFriendIDListFromCacheReq{OperationID: req.CommID.OperationID, UserID: req.CommID.FromUserID})
	if err != nil {
		log.NewError(req.CommID.OperationID, utils.GetSelfFuncName(), "DelFriendIDListFromCache rpc failed", err.Error())
		return &pbFriend.DeleteFriendResp{CommonResp: &pbFriend.CommonResp{ErrCode: constant.ErrServer.ErrCode, ErrMsg: constant.ErrServer.ErrMsg}}, nil
	}
	if respPb.CommonResp.ErrCode != 0 {
		log.NewError(req.CommID.OperationID, utils.GetSelfFuncName(), "DelFriendIDListFromCache failed", respPb.CommonResp.String())
		return &pbFriend.DeleteFriendResp{CommonResp: &pbFriend.CommonResp{ErrCode: respPb.CommonResp.ErrCode, ErrMsg: respPb.CommonResp.ErrMsg}}, nil
	}
	if err := rocksCache.DelAllFriendsInfoFromCache(req.CommID.FromUserID); err != nil {
		log.NewError(req.CommID.OperationID, utils.GetSelfFuncName(), err.Error(), req.CommID.FromUserID)
	}
	if err := rocksCache.DelAllFriendsInfoFromCache(req.CommID.ToUserID); err != nil {
		log.NewError(req.CommID.OperationID, utils.GetSelfFuncName(), err.Error(), req.CommID.FromUserID)
	}
	chat.FriendDeletedNotification(req)
	return &pbFriend.DeleteFriendResp{CommonResp: &pbFriend.CommonResp{}}, nil
}

func (s *friendServer) GetBlacklist(ctx context.Context, req *pbFriend.GetBlacklistReq) (*pbFriend.GetBlacklistResp, error) {
	log.NewInfo(req.CommID.OperationID, "GetBlacklist args ", req.String())

	if !token_verify.CheckAccess(req.CommID.OpUserID, req.CommID.FromUserID) {
		log.NewError(req.CommID.OperationID, "CheckAccess false ", req.CommID.OpUserID, req.CommID.FromUserID)
		return &pbFriend.GetBlacklistResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}, nil
	}

	blackIDList, err := rocksCache.GetBlackListFromCache(req.CommID.FromUserID)
	if err != nil {
		log.NewError(req.CommID.OperationID, "GetBlackListByUID failed ", err.Error(), req.CommID.FromUserID)
		return &pbFriend.GetBlacklistResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}, nil
	}

	var (
		userInfoList []*sdkws.PublicUserInfo
	)
	for _, userID := range blackIDList {
		user, err := rocksCache.GetUserInfoFromCache(userID)
		if err != nil {
			log.NewError(req.CommID.OperationID, "GetUserByUserID failed ", err.Error(), userID)
			continue
		}
		var blackUserInfo sdkws.PublicUserInfo
		utils.CopyStructFields(&blackUserInfo, user)
		userInfoList = append(userInfoList, &blackUserInfo)
	}
	log.NewInfo(req.CommID.OperationID, "rpc GetBlacklist ok ", pbFriend.GetBlacklistResp{BlackUserInfoList: userInfoList})
	return &pbFriend.GetBlacklistResp{BlackUserInfoList: userInfoList}, nil
}

func (s *friendServer) SetFriendRemark(ctx context.Context, req *pbFriend.SetFriendRemarkReq) (*pbFriend.SetFriendRemarkResp, error) {
	log.NewInfo(req.CommID.OperationID, "SetFriendComment args ", req.String())
	//Parse token, to find current user information
	if !token_verify.CheckAccess(req.CommID.OpUserID, req.CommID.FromUserID) {
		log.NewError(req.CommID.OperationID, "CheckAccess false ", req.CommID.OpUserID, req.CommID.FromUserID)
		return &pbFriend.SetFriendRemarkResp{CommonResp: &pbFriend.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}}, nil
	}

	err := imdb.UpdateFriendComment(req.CommID.FromUserID, req.CommID.ToUserID, req.Remark)
	if err != nil {
		log.NewError(req.CommID.OperationID, "UpdateFriendComment failed ", req.CommID.FromUserID, req.CommID.ToUserID, req.Remark, err.Error())
		return &pbFriend.SetFriendRemarkResp{CommonResp: &pbFriend.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}
	log.NewInfo(req.CommID.OperationID, "rpc SetFriendComment ok")
	err = rocksCache.DelAllFriendsInfoFromCache(req.CommID.FromUserID)
	if err != nil {
		log.NewError(req.CommID.OperationID, "DelAllFriendInfoFromCache failed ", req.CommID.FromUserID, err.Error())
		return &pbFriend.SetFriendRemarkResp{CommonResp: &pbFriend.CommonResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}}, nil
	}
	chat.FriendRemarkSetNotification(req.CommID.OperationID, req.CommID.OpUserID, req.CommID.FromUserID, req.CommID.ToUserID)
	return &pbFriend.SetFriendRemarkResp{CommonResp: &pbFriend.CommonResp{}}, nil
}

func (s *friendServer) RemoveBlacklist(ctx context.Context, req *pbFriend.RemoveBlacklistReq) (*pbFriend.RemoveBlacklistResp, error) {
	log.NewInfo(req.CommID.OperationID, "RemoveBlacklist args ", req.String())
	//Parse token, to find current user information
	if !token_verify.CheckAccess(req.CommID.OpUserID, req.CommID.FromUserID) {
		log.NewError(req.CommID.OperationID, "CheckAccess false ", req.CommID.OpUserID, req.CommID.FromUserID)
		return &pbFriend.RemoveBlacklistResp{CommonResp: &pbFriend.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}}, nil
	}
	err := imdb.RemoveBlackList(req.CommID.FromUserID, req.CommID.ToUserID)
	if err != nil {
		log.NewError(req.CommID.OperationID, "RemoveBlackList failed", err.Error(), req.CommID.FromUserID, req.CommID.ToUserID)
		return &pbFriend.RemoveBlacklistResp{CommonResp: &pbFriend.CommonResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}}, nil

	}
	log.NewInfo(req.CommID.OperationID, "rpc RemoveBlacklist ok ")

	etcdConn := getcdv3.GetConn(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), config.Config.RpcRegisterName.OpenImCacheName, req.CommID.OperationID)
	if etcdConn == nil {
		errMsg := req.CommID.OperationID + "getcdv3.GetConn == nil"
		log.NewError(req.CommID.OperationID, errMsg)
		return &pbFriend.RemoveBlacklistResp{CommonResp: &pbFriend.CommonResp{ErrCode: constant.ErrInternal.ErrCode, ErrMsg: errMsg}}, nil
	}

	cacheClient := pbCache.NewCacheClient(etcdConn)
	cacheResp, err := cacheClient.DelBlackIDListFromCache(context.Background(), &pbCache.DelBlackIDListFromCacheReq{UserID: req.CommID.FromUserID, OperationID: req.CommID.OperationID})
	if err != nil {
		log.NewError(req.CommID.OperationID, "DelBlackIDListFromCache rpc call failed ", err.Error())
		return &pbFriend.RemoveBlacklistResp{CommonResp: &pbFriend.CommonResp{ErrCode: 500, ErrMsg: "ReduceBlackUserFromCache rpc call failed"}}, nil
	}
	if cacheResp.CommonResp.ErrCode != 0 {
		log.NewError(req.CommID.OperationID, "DelBlackIDListFromCache rpc logic call failed ", cacheResp.CommonResp.String())
		return &pbFriend.RemoveBlacklistResp{CommonResp: &pbFriend.CommonResp{ErrCode: cacheResp.CommonResp.ErrCode, ErrMsg: cacheResp.CommonResp.ErrMsg}}, nil
	}
	chat.BlackDeletedNotification(req)
	return &pbFriend.RemoveBlacklistResp{CommonResp: &pbFriend.CommonResp{}}, nil
}

func (s *friendServer) IsInBlackList(ctx context.Context, req *pbFriend.IsInBlackListReq) (*pbFriend.IsInBlackListResp, error) {
	log.NewInfo("IsInBlackList args ", req.String())
	if !token_verify.CheckAccess(req.CommID.OpUserID, req.CommID.FromUserID) {
		log.NewError(req.CommID.OperationID, "CheckAccess false ", req.CommID.OpUserID, req.CommID.FromUserID)
		return &pbFriend.IsInBlackListResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}, nil
	}
	blackIDList, err := rocksCache.GetBlackListFromCache(req.CommID.FromUserID)
	if err != nil {
		log.NewError(req.CommID.OperationID, utils.GetSelfFuncName(), err.Error(), req.CommID.FromUserID)
		return &pbFriend.IsInBlackListResp{ErrMsg: err.Error(), ErrCode: constant.ErrDB.ErrCode}, nil
	}
	var isInBlacklist bool
	if utils.IsContain(req.CommID.ToUserID, blackIDList) {
		isInBlacklist = true
	}
	log.NewInfo(req.CommID.OperationID, "IsInBlackList rpc ok ", pbFriend.IsInBlackListResp{Response: isInBlacklist})
	return &pbFriend.IsInBlackListResp{Response: isInBlacklist}, nil
}

func (s *friendServer) IsFriend(ctx context.Context, req *pbFriend.IsFriendReq) (*pbFriend.IsFriendResp, error) {
	log.NewInfo(req.CommID.OperationID, req.String())
	if !token_verify.CheckAccess(req.CommID.OpUserID, req.CommID.FromUserID) {
		log.NewError(req.CommID.OperationID, "CheckAccess false ", req.CommID.OpUserID, req.CommID.FromUserID)
		return &pbFriend.IsFriendResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}, nil
	}
	friendIDList, err := rocksCache.GetFriendIDListFromCache(req.CommID.FromUserID)
	if err != nil {
		log.NewError(req.CommID.OperationID, utils.GetSelfFuncName(), err.Error(), req.CommID.FromUserID)
		return &pbFriend.IsFriendResp{ErrMsg: err.Error(), ErrCode: constant.ErrDB.ErrCode}, nil
	}
	var isFriend bool
	if utils.IsContain(req.CommID.ToUserID, friendIDList) {
		isFriend = true
	}
	log.NewInfo(req.CommID.OperationID, pbFriend.IsFriendResp{Response: isFriend})
	return &pbFriend.IsFriendResp{Response: isFriend}, nil
}

func (s *friendServer) GetFriendList(ctx context.Context, req *pbFriend.GetFriendListReq) (*pbFriend.GetFriendListResp, error) {
	log.NewInfo("GetFriendList args ", req.String())
	if !token_verify.CheckAccess(req.CommID.OpUserID, req.CommID.FromUserID) {
		log.NewError(req.CommID.OperationID, "CheckAccess false ", req.CommID.OpUserID, req.CommID.FromUserID)
		return &pbFriend.GetFriendListResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}, nil
	}
	friendList, err := rocksCache.GetAllFriendsInfoFromCache(req.CommID.FromUserID)
	if err != nil {
		log.NewError(req.CommID.OperationID, utils.GetSelfFuncName(), "GetFriendIDListFromCache failed", err.Error(), req.CommID.FromUserID)
		return &pbFriend.GetFriendListResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}, nil
	}

	var userInfoList []*sdkws.FriendInfo
	for _, friendUser := range friendList {
		friendUserInfo := sdkws.FriendInfo{FriendUser: &sdkws.UserInfo{}}
		cp.FriendDBCopyOpenIM(&friendUserInfo, friendUser)
		log.NewDebug(req.CommID.OperationID, "friends : ", friendUser, "openim friends: ", friendUserInfo)
		userInfoList = append(userInfoList, &friendUserInfo)
	}
	log.NewInfo(req.CommID.OperationID, "rpc GetFriendList ok", pbFriend.GetFriendListResp{FriendInfoList: userInfoList})
	return &pbFriend.GetFriendListResp{FriendInfoList: userInfoList}, nil
}

//received
func (s *friendServer) GetFriendApplyList(ctx context.Context, req *pbFriend.GetFriendApplyListReq) (*pbFriend.GetFriendApplyListResp, error) {
	log.NewInfo(req.CommID.OperationID, "GetFriendApplyList args ", req.String())
	//Parse token, to find current user information
	if !token_verify.CheckAccess(req.CommID.OpUserID, req.CommID.FromUserID) {
		log.NewError(req.CommID.OperationID, "CheckAccess false ", req.CommID.OpUserID, req.CommID.FromUserID)
		return &pbFriend.GetFriendApplyListResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}, nil
	}
	//	Find the  current user friend applications received
	ApplyUsersInfo, err := imdb.GetReceivedFriendsApplicationListByUserID(req.CommID.FromUserID)
	if err != nil {
		log.NewError(req.CommID.OperationID, "GetReceivedFriendsApplicationListByUserID ", err.Error(), req.CommID.FromUserID)
		return &pbFriend.GetFriendApplyListResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}, nil
	}

	var appleUserList []*sdkws.FriendRequest
	for _, applyUserInfo := range ApplyUsersInfo {
		var userInfo sdkws.FriendRequest
		cp.FriendRequestDBCopyOpenIM(&userInfo, &applyUserInfo)
		//	utils.CopyStructFields(&userInfo, applyUserInfo)
		//	u, err := imdb.GetUserByUserID(userInfo.FromUserID)
		//	if err != nil {
		//		log.Error(req.CommID.OperationID, "GetUserByUserID", userInfo.FromUserID)
		//		continue
		//	}
		//	userInfo.FromNickname = u.Nickname
		//	userInfo.FromFaceURL = u.FaceURL
		//	userInfo.FromGender = u.Gender
		//
		//	u, err = imdb.GetUserByUserID(userInfo.ToUserID)
		//	if err != nil {
		//		log.Error(req.CommID.OperationID, "GetUserByUserID", userInfo.ToUserID)
		//		continue
		//	}
		//	userInfo.ToNickname = u.Nickname
		//	userInfo.ToFaceURL = u.FaceURL
		//	userInfo.ToGender = u.Gender
		appleUserList = append(appleUserList, &userInfo)
	}

	log.NewInfo(req.CommID.OperationID, "rpc GetFriendApplyList ok", pbFriend.GetFriendApplyListResp{FriendRequestList: appleUserList})
	return &pbFriend.GetFriendApplyListResp{FriendRequestList: appleUserList}, nil
}

func (s *friendServer) GetSelfApplyList(ctx context.Context, req *pbFriend.GetSelfApplyListReq) (*pbFriend.GetSelfApplyListResp, error) {
	log.NewInfo(req.CommID.OperationID, "GetSelfApplyList args ", req.String())

	//Parse token, to find current user information
	if !token_verify.CheckAccess(req.CommID.OpUserID, req.CommID.FromUserID) {
		log.NewError(req.CommID.OperationID, "CheckAccess false ", req.CommID.OpUserID, req.CommID.FromUserID)
		return &pbFriend.GetSelfApplyListResp{ErrCode: constant.ErrAccess.ErrCode, ErrMsg: constant.ErrAccess.ErrMsg}, nil
	}
	//	Find the self add other userinfo
	usersInfo, err := imdb.GetSendFriendApplicationListByUserID(req.CommID.FromUserID)
	if err != nil {
		log.NewError(req.CommID.OperationID, "GetSendFriendApplicationListByUserID failed ", err.Error(), req.CommID.FromUserID)
		return &pbFriend.GetSelfApplyListResp{ErrCode: constant.ErrDB.ErrCode, ErrMsg: constant.ErrDB.ErrMsg}, nil
	}
	var selfApplyOtherUserList []*sdkws.FriendRequest
	for _, selfApplyOtherUserInfo := range usersInfo {
		var userInfo sdkws.FriendRequest // pbFriend.ApplyUserInfo
		cp.FriendRequestDBCopyOpenIM(&userInfo, &selfApplyOtherUserInfo)
		//u, err := imdb.GetUserByUserID(userInfo.FromUserID)
		//if err != nil {
		//	log.Error(req.CommID.OperationID, "GetUserByUserID", userInfo.FromUserID)
		//	continue
		//}
		//userInfo.FromNickname = u.Nickname
		//userInfo.FromFaceURL = u.FaceURL
		//userInfo.FromGender = u.Gender
		//
		//u, err = imdb.GetUserByUserID(userInfo.ToUserID)
		//if err != nil {
		//	log.Error(req.CommID.OperationID, "GetUserByUserID", userInfo.ToUserID)
		//	continue
		//}
		//userInfo.ToNickname = u.Nickname
		//userInfo.ToFaceURL = u.FaceURL
		//userInfo.ToGender = u.Gender

		selfApplyOtherUserList = append(selfApplyOtherUserList, &userInfo)
	}
	log.NewInfo(req.CommID.OperationID, "rpc GetSelfApplyList ok", pbFriend.GetSelfApplyListResp{FriendRequestList: selfApplyOtherUserList})
	return &pbFriend.GetSelfApplyListResp{FriendRequestList: selfApplyOtherUserList}, nil
}

////
//func (s *friendServer) GetFriendsInfo(ctx context.Context, req *pbFriend.GetFriendsInfoReq) (*pbFriend.GetFriendInfoResp, error) {
//	return nil, nil
////	log.NewInfo(req.CommID.OperationID, "GetFriendsInfo args ", req.String())
////	var (
////		isInBlackList int32
////		//	isFriend      int32
////		comment string
////	)
////
////	friendShip, err := imdb.FindFriendRelationshipFromFriend(req.CommID.FromUserID, req.CommID.ToUserID)
////	if err != nil {
////		log.NewError(req.CommID.OperationID, "FindFriendRelationshipFromFriend failed ", err.Error())
////		return &pbFriend.GetFriendInfoResp{ErrCode: constant.ErrSearchUserInfo.ErrCode, ErrMsg: constant.ErrSearchUserInfo.ErrMsg}, nil
////		//	isFriend = constant.FriendFlag
////	}
////	comment = friendShip.Remark
////
////	friendUserInfo, err := imdb.FindUserByUID(req.CommID.ToUserID)
////	if err != nil {
////		log.NewError(req.CommID.OperationID, "FindUserByUID failed ", err.Error())
////		return &pbFriend.GetFriendInfoResp{ErrCode: constant.ErrSearchUserInfo.ErrCode, ErrMsg: constant.ErrSearchUserInfo.ErrMsg}, nil
////	}
////
////	err = imdb.FindRelationshipFromBlackList(req.CommID.FromUserID, req.CommID.ToUserID)
////	if err == nil {
////		isInBlackList = constant.BlackListFlag
////	}
////
////	resp := pbFriend.GetFriendInfoResp{ErrCode: 0, ErrMsg:  "",}
////
////	utils.CopyStructFields(resp.FriendInfoList, friendUserInfo)
////	resp.Data.IsBlack = isInBlackList
////	resp.Data.OwnerUserID = req.CommID.FromUserID
////	resp.Data.Remark = comment
////	resp.Data.CreateTime = friendUserInfo.CreateTime
////
////	log.NewInfo(req.CommID.OperationID, "GetFriendsInfo ok ", resp)
////	return &resp, nil
////
//}
