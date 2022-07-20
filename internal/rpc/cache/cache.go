package cache

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db/rocks_cache"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbCache "Open_IM/pkg/proto/cache"
	"Open_IM/pkg/utils"
	"context"
	"google.golang.org/grpc"
	"net"
	"strconv"
	"strings"
)

type cacheServer struct {
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string
}

func NewCacheServer(port int) *cacheServer {
	log.NewPrivateLog(constant.LogFileName)
	return &cacheServer{
		rpcPort:         port,
		rpcRegisterName: config.Config.RpcRegisterName.OpenImCacheName,
		etcdSchema:      config.Config.Etcd.EtcdSchema,
		etcdAddr:        config.Config.Etcd.EtcdAddr,
	}
}

func (s *cacheServer) Run() {
	log.NewInfo("0", "cacheServer rpc start ")
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
	log.NewInfo("0", "listen network success, ", address, listener)
	defer listener.Close()
	//grpc server

	srv := grpc.NewServer()
	defer srv.GracefulStop()
	pbCache.RegisterCacheServer(srv, s)

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
		log.NewError("0", "RegisterEtcd failed ", err.Error())
		return
	}
	err = srv.Serve(listener)
	if err != nil {
		log.NewError("0", "Serve failed ", err.Error())
		return
	}
	log.NewInfo("0", "message cms rpc success")
}

func (s *cacheServer) GetFriendIDListFromCache(_ context.Context, req *pbCache.GetFriendIDListFromCacheReq) (resp *pbCache.GetFriendIDListFromCacheResp, err error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp = &pbCache.GetFriendIDListFromCacheResp{CommonResp: &pbCache.CommonResp{}}
	friendIDList, err := rocksCache.GetFriendIDListFromCache(req.UserID)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetFriendIDListFromCache", err.Error())
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = constant.ErrDB.ErrMsg
		return resp, nil
	}
	log.NewDebug(req.OperationID, utils.GetSelfFuncName(), friendIDList)
	resp.UserIDList = friendIDList
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp.String())
	return resp, nil
}

// this is for dtm call
func (s *cacheServer) DelFriendIDListFromCache(_ context.Context, req *pbCache.DelFriendIDListFromCacheReq) (resp *pbCache.DelFriendIDListFromCacheResp, err error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp = &pbCache.DelFriendIDListFromCacheResp{CommonResp: &pbCache.CommonResp{}}
	if err := rocksCache.DelFriendIDListFromCache(req.UserID); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "args: ", req.UserID, err.Error())
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = err.Error()
		return resp, nil
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", resp.String())
	return resp, nil
}

func (s *cacheServer) GetBlackIDListFromCache(_ context.Context, req *pbCache.GetBlackIDListFromCacheReq) (resp *pbCache.GetBlackIDListFromCacheResp, err error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp = &pbCache.GetBlackIDListFromCacheResp{CommonResp: &pbCache.CommonResp{}}
	blackUserIDList, err := rocksCache.GetBlackListFromCache(req.UserID)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "AddFriendToCache failed", err.Error())
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = constant.ErrDB.ErrMsg
		return resp, nil
	}
	log.NewDebug(req.OperationID, utils.GetSelfFuncName(), blackUserIDList)
	resp.UserIDList = blackUserIDList
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp.String())
	return resp, nil
}

func (s *cacheServer) DelBlackIDListFromCache(_ context.Context, req *pbCache.DelBlackIDListFromCacheReq) (resp *pbCache.DelBlackIDListFromCacheResp, err error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp = &pbCache.DelBlackIDListFromCacheResp{CommonResp: &pbCache.CommonResp{}}
	if err := rocksCache.DelBlackIDListFromCache(req.UserID); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "args: ", req.UserID, err.Error())
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = err.Error()
		return resp, nil
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", resp.String())
	return resp, nil
}

func (s *cacheServer) GetGroupMemberIDListFromCache(_ context.Context, req *pbCache.GetGroupMemberIDListFromCacheReq) (resp *pbCache.GetGroupMemberIDListFromCacheResp, err error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp = &pbCache.GetGroupMemberIDListFromCacheResp{
		CommonResp: &pbCache.CommonResp{},
	}
	userIDList, err := rocksCache.GetGroupMemberIDListFromCache(req.GroupID)
	if err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "GetGroupMemberIDListFromCache failed", err.Error())
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = constant.ErrDB.ErrMsg
		return resp, nil
	}
	resp.UserIDList = userIDList
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp.String())
	return resp, nil
}

func (s *cacheServer) DelGroupMemberIDListFromCache(_ context.Context, req *pbCache.DelGroupMemberIDListFromCacheReq) (resp *pbCache.DelGroupMemberIDListFromCacheResp, err error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp = &pbCache.DelGroupMemberIDListFromCacheResp{CommonResp: &pbCache.CommonResp{}}
	if err := rocksCache.DelGroupMemberIDListFromCache(req.GroupID); err != nil {
		log.NewError(req.OperationID, utils.GetSelfFuncName(), "args: ", req.GroupID, err.Error())
		resp.CommonResp.ErrCode = constant.ErrDB.ErrCode
		resp.CommonResp.ErrMsg = err.Error()
		return resp, nil
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", resp.String())
	return resp, nil
}
