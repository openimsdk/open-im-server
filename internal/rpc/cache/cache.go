package cache

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/grpc-etcdv3/getcdv3"
	pbCache "Open_IM/pkg/proto/cache"
	commonPb "Open_IM/pkg/proto/sdk_ws"
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
	ip := utils.ServerIP
	registerAddress := ip + ":" + strconv.Itoa(s.rpcPort)
	//listener network
	listener, err := net.Listen("tcp", registerAddress)
	if err != nil {
		log.NewError("0", "Listen failed ", err.Error(), registerAddress)
		return
	}
	log.NewInfo("0", "listen network success, ", registerAddress, listener)
	defer listener.Close()
	//grpc server
	srv := grpc.NewServer()
	defer srv.GracefulStop()
	pbCache.RegisterCacheServer(srv, s)
	err = getcdv3.RegisterEtcd(s.etcdSchema, strings.Join(s.etcdAddr, ","), ip, s.rpcPort, s.rpcRegisterName, 10)
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

func (s *cacheServer) GetUserInfo(_ context.Context, req *pbCache.GetUserInfoReq) (resp *pbCache.GetUserInfoResp, err error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp = &pbCache.GetUserInfoResp{
		UserInfoList: []*commonPb.UserInfo{},
		CommonResp:   &pbCache.CommonResp{},
	}
	for _, userID := range req.UserIDList {
		userInfo, err := db.DB.GetUserInfo(userID)
		if err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "get userInfo from cache failed", err.Error())
			continue
		}
		resp.UserInfoList = append(resp.UserInfoList, userInfo)
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp.String())
	return resp, nil
}

func (s *cacheServer) UpdateUserInfo(_ context.Context, req *pbCache.UpdateUserInfoReq) (resp *pbCache.UpdateUserInfoResp, err error) {
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", req.String())
	resp = &pbCache.UpdateUserInfoResp{
		CommonResp: &pbCache.CommonResp{},
	}
	for _, userInfo := range req.UserInfoList {
		if err := db.DB.SetUserInfo(userInfo); err != nil {
			log.NewError(req.OperationID, utils.GetSelfFuncName(), "set userInfo to cache failed", err.Error())
			return resp, nil
		}
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "resp: ", resp.String())
	return resp, nil
}

func (s *cacheServer) GetBlackList(_ context.Context, req *pbCache.GetBlackListReq) (resp *pbCache.GetBlackListResp, err error) {
	return nil, nil
}

func (s *cacheServer) UpdateBlackList(_ context.Context, req *pbCache.UpdateBlackListReq) (resp *pbCache.UpdateBlackListResp, err error) {
	return nil, nil
}

func (s *cacheServer) GetFriendInfo(_ context.Context, req *pbCache.GetFriendInfoReq) (resp *pbCache.GetFriendInfoResp, err error) {
	return nil, nil
}

func (s *cacheServer) UpdateFriendInfo(_ context.Context, req *pbCache.UpdateFriendInfoReq) (resp *pbCache.UpdateFriendInfoResp, err error) {
	return nil, nil
}
