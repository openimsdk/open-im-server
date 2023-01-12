package cache

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	rocksCache "Open_IM/pkg/common/db/rocks_cache"
	"Open_IM/pkg/common/log"
	promePkg "Open_IM/pkg/common/prometheus"
	pbCache "Open_IM/pkg/proto/cache"
	"Open_IM/pkg/utils"
	"context"
	"github.com/OpenIMSDK/getcdv3"
	"net"
	"strconv"
	"strings"

	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
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
	pbCache.RegisterCacheServer(srv, s)

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
		log.NewError("0", "RegisterEtcd failed ", err.Error())
		panic(utils.Wrap(err, "register cache module  rpc to etcd err"))
	}
	go rocksCache.DelKeys()
	err = srv.Serve(listener)
	if err != nil {
		log.NewError("0", "Serve failed ", err.Error())
		return
	}
	log.NewInfo("0", "message cms rpc success")
}

func (s *cacheServer) GetFriendIDListFromCache(ctx context.Context, req *pbCache.GetFriendIDListFromCacheReq) (resp *pbCache.GetFriendIDListFromCacheResp, err error) {
	resp = &pbCache.GetFriendIDListFromCacheResp{}
	friendIDList, err := rocksCache.GetFriendIDListFromCache(ctx, req.UserID)
	if err != nil {
		return
	}
	resp.UserIDList = friendIDList
	return
}

// this is for dtm call
func (s *cacheServer) DelFriendIDListFromCache(ctx context.Context, req *pbCache.DelFriendIDListFromCacheReq) (resp *pbCache.DelFriendIDListFromCacheResp, err error) {
	resp = &pbCache.DelFriendIDListFromCacheResp{}
	if err := rocksCache.DelFriendIDListFromCache(ctx, req.UserID); err != nil {
		return
	}
	log.NewInfo(req.OperationID, utils.GetSelfFuncName(), "req: ", resp.String())
	return
}

func (s *cacheServer) GetBlackIDListFromCache(ctx context.Context, req *pbCache.GetBlackIDListFromCacheReq) (resp *pbCache.GetBlackIDListFromCacheResp, err error) {
	resp = &pbCache.GetBlackIDListFromCacheResp{}
	blackUserIDList, err := rocksCache.GetBlackListFromCache(ctx, req.UserID)
	if err != nil {
		return
	}
	resp.UserIDList = blackUserIDList
	return
}

func (s *cacheServer) DelBlackIDListFromCache(ctx context.Context, req *pbCache.DelBlackIDListFromCacheReq) (resp *pbCache.DelBlackIDListFromCacheResp, err error) {
	resp = &pbCache.DelBlackIDListFromCacheResp{}
	if err := rocksCache.DelBlackIDListFromCache(ctx, req.UserID); err != nil {
		return
	}
	return resp, nil
}

func (s *cacheServer) GetGroupMemberIDListFromCache(ctx context.Context, req *pbCache.GetGroupMemberIDListFromCacheReq) (resp *pbCache.GetGroupMemberIDListFromCacheResp, err error) {
	resp = &pbCache.GetGroupMemberIDListFromCacheResp{}
	userIDList, err := rocksCache.GetGroupMemberIDListFromCache(ctx, req.GroupID)
	if err != nil {
		return
	}
	resp.UserIDList = userIDList
	return
}

func (s *cacheServer) DelGroupMemberIDListFromCache(ctx context.Context, req *pbCache.DelGroupMemberIDListFromCacheReq) (resp *pbCache.DelGroupMemberIDListFromCacheResp, err error) {
	resp = &pbCache.DelGroupMemberIDListFromCacheResp{}
	if err := rocksCache.DelGroupMemberIDListFromCache(ctx, req.GroupID); err != nil {
		return resp, nil
	}
	return resp, nil
}
