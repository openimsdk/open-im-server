package push

import (
	"OpenIM/pkg/common/config"
	"OpenIM/pkg/common/constant"
	"OpenIM/pkg/common/db/cache"
	"OpenIM/pkg/common/db/controller"
	"OpenIM/pkg/common/log"
	"OpenIM/pkg/common/prome"
	pbPush "OpenIM/pkg/proto/push"
	"OpenIM/pkg/utils"
	"context"
	"net"
	"strconv"
	"strings"

	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
)

type RPCServer struct {
	rpcPort         int
	rpcRegisterName string
	pushInterface   controller.PushInterface
	pusher          Pusher
}

func (r *RPCServer) Init(rpcPort int, cache cache.Cache) {
	r.rpcPort = rpcPort
	r.rpcRegisterName = config.Config.RpcRegisterName.OpenImPushName
}

func (r *RPCServer) run() {
	listenIP := ""
	if config.Config.ListenIP == "" {
		listenIP = "0.0.0.0"
	} else {
		listenIP = config.Config.ListenIP
	}
	address := listenIP + ":" + strconv.Itoa(r.rpcPort)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic("listening err:" + err.Error() + r.rpcRegisterName)
	}
	defer listener.Close()
	var grpcOpts []grpc.ServerOption
	if config.Config.Prometheus.Enable {
		prome.NewGrpcRequestCounter()
		prome.NewGrpcRequestFailedCounter()
		prome.NewGrpcRequestSuccessCounter()
		grpcOpts = append(grpcOpts, []grpc.ServerOption{
			// grpc.UnaryInterceptor(prome.UnaryServerInterceptorProme),
			grpc.StreamInterceptor(grpcPrometheus.StreamServerInterceptor),
			grpc.UnaryInterceptor(grpcPrometheus.UnaryServerInterceptor),
		}...)
	}
	srv := grpc.NewServer(grpcOpts...)
	defer srv.GracefulStop()
	pbPush.RegisterPushMsgServiceServer(srv, r)
	rpcRegisterIP := config.Config.RpcRegisterIP
	if config.Config.RpcRegisterIP == "" {
		rpcRegisterIP, err = utils.GetLocalIP()
		if err != nil {
			log.Error("", "GetLocalIP failed ", err.Error())
		}
	}

	err = rpc.RegisterEtcd(r.etcdSchema, strings.Join(r.etcdAddr, ","), rpcRegisterIP, r.rpcPort, r.rpcRegisterName, 10)
	if err != nil {
		log.Error("", "register push module  rpc to etcd err", err.Error(), r.etcdSchema, strings.Join(r.etcdAddr, ","), rpcRegisterIP, r.rpcPort, r.rpcRegisterName)
		panic(utils.Wrap(err, "register push module  rpc to etcd err"))
	}
	err = srv.Serve(listener)
	if err != nil {
		log.Error("", "push module rpc start err", err.Error())
		return
	}
}

func (r *RPCServer) PushMsg(ctx context.Context, pbData *pbPush.PushMsgReq) (resp *pbPush.PushMsgResp, err error) {
	switch pbData.MsgData.SessionType {
	case constant.SuperGroupChatType:
		err = r.pusher.MsgToSuperGroupUser(ctx, pbData.SourceID, pbData.MsgData)
	default:
		err = r.pusher.MsgToUser(ctx, pbData.SourceID, pbData.MsgData)
	}
	return &pbPush.PushMsgResp{}, err
}

func (r *RPCServer) DelUserPushToken(ctx context.Context, req *pbPush.DelUserPushTokenReq) (resp *pbPush.DelUserPushTokenResp, err error) {
	return &pbPush.DelUserPushTokenResp{}, r.pushInterface.DelFcmToken(ctx, req.UserID, int(req.PlatformID))
}
