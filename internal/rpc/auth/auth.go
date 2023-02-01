package auth

import (
	"Open_IM/internal/common/check"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db/controller"
	"Open_IM/pkg/common/log"
	promePkg "Open_IM/pkg/common/prometheus"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/common/tracelog"
	"Open_IM/pkg/getcdv3"
	pbAuth "Open_IM/pkg/proto/auth"
	pbRelay "Open_IM/pkg/proto/relay"
	"Open_IM/pkg/utils"
	"context"
	"net"
	"strconv"
	"strings"

	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"

	"Open_IM/pkg/common/config"

	"google.golang.org/grpc"
)

func (s *rpcAuth) UserToken(ctx context.Context, req *pbAuth.UserTokenReq) (*pbAuth.UserTokenResp, error) {
	resp := pbAuth.UserTokenResp{}
	if _, err := check.GetUsersInfo(ctx, req.UserID); err != nil {
		return nil, err
	}
	token, err := s.CreateToken(ctx, req.UserID, int(req.PlatformID), config.Config.TokenPolicy.AccessExpire)
	if err != nil {
		return nil, err
	}
	resp.Token = token
	resp.ExpireTimeSeconds = config.Config.TokenPolicy.AccessExpire
	return &resp, nil
}

func (s *rpcAuth) parseToken(ctx context.Context, tokensString, operationID string) (claims *token_verify.Claims, err error) {
	claims, err = token_verify.GetClaimFromToken(tokensString)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	m, err := s.GetTokens(ctx, claims.UID, claims.Platform)
	if err != nil {
		return nil, err
	}

	if v, ok := m[tokensString]; ok {
		switch v {
		case constant.NormalToken:
			return claims, nil
		case constant.KickedToken:
			return nil, utils.Wrap(constant.ErrTokenKicked, "this token has been kicked by other same terminal ")
		default:
			return nil, utils.Wrap(constant.ErrTokenUnknown, "")
		}
	}
	return nil, utils.Wrap(constant.ErrTokenNotExist, "redis token map not find")
}

func (s *rpcAuth) ParseToken(ctx context.Context, req *pbAuth.ParseTokenReq) (*pbAuth.ParseTokenResp, error) {
	resp := pbAuth.ParseTokenResp{}
	claims, err := s.parseToken(ctx, req.Token, req.OperationID)
	if err != nil {
		return nil, err
	}
	resp.UserID = claims.UID
	resp.Platform = claims.Platform
	resp.ExpireTimeSeconds = claims.ExpiresAt.Unix()
	return &resp, nil
}

func (s *rpcAuth) ForceLogout(ctx context.Context, req *pbAuth.ForceLogoutReq) (*pbAuth.ForceLogoutResp, error) {
	resp := pbAuth.ForceLogoutResp{}
	if err := token_verify.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	if err := s.forceKickOff(ctx, req.UserID, req.PlatformID, tracelog.GetOperationID(ctx)); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *rpcAuth) forceKickOff(ctx context.Context, userID string, platformID int32, operationID string) error {
	grpcCons := getcdv3.GetDefaultGatewayConn4Unique(config.Config.Etcd.EtcdSchema, strings.Join(config.Config.Etcd.EtcdAddr, ","), operationID)
	for _, v := range grpcCons {
		client := pbRelay.NewRelayClient(v)
		kickReq := &pbRelay.KickUserOfflineReq{OperationID: operationID, KickUserIDList: []string{userID}, PlatformID: platformID}
		log.NewInfo(operationID, "KickUserOffline ", client, kickReq.String())
		_, err := client.KickUserOffline(ctx, kickReq)
		return utils.Wrap(err, "")
	}
	return constant.ErrInternalServer.Wrap()
}

type rpcAuth struct {
	rpcPort         int
	rpcRegisterName string
	etcdSchema      string
	etcdAddr        []string
	controller.AuthInterface
}

func NewRpcAuthServer(port int) *rpcAuth {
	log.NewPrivateLog(constant.LogFileName)
	return &rpcAuth{
		rpcPort:         port,
		rpcRegisterName: config.Config.RpcRegisterName.OpenImAuthName,
		etcdSchema:      config.Config.Etcd.EtcdSchema,
		etcdAddr:        config.Config.Etcd.EtcdAddr,
	}
}

func (s *rpcAuth) Run() {
	operationID := utils.OperationIDGenerator()
	log.NewInfo(operationID, "rpc auth start...")

	listenIP := ""
	if config.Config.ListenIP == "" {
		listenIP = "0.0.0.0"
	} else {
		listenIP = config.Config.ListenIP
	}
	address := listenIP + ":" + strconv.Itoa(s.rpcPort)
	listener, err := net.Listen("tcp", address)
	if err != nil {
		panic("listening err:" + err.Error() + s.rpcRegisterName)
	}
	log.NewInfo(operationID, "listen network success, ", address, listener)
	var grpcOpts []grpc.ServerOption
	if config.Config.Prometheus.Enable {
		promePkg.NewGrpcRequestCounter()
		promePkg.NewGrpcRequestFailedCounter()
		promePkg.NewGrpcRequestSuccessCounter()
		promePkg.NewUserRegisterCounter()
		promePkg.NewUserLoginCounter()
		grpcOpts = append(grpcOpts, []grpc.ServerOption{
			// grpc.UnaryInterceptor(promePkg.UnaryServerInterceptorProme),
			grpc.StreamInterceptor(grpcPrometheus.StreamServerInterceptor),
			grpc.UnaryInterceptor(grpcPrometheus.UnaryServerInterceptor),
		}...)
	}
	srv := grpc.NewServer(grpcOpts...)
	defer srv.GracefulStop()

	//service registers with etcd
	pbAuth.RegisterAuthServer(srv, s)
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
		log.NewError(operationID, "RegisterEtcd failed ", err.Error(),
			s.etcdSchema, strings.Join(s.etcdAddr, ","), rpcRegisterIP, s.rpcPort, s.rpcRegisterName)
		panic(utils.Wrap(err, "register auth module  rpc to etcd err"))

	}
	log.NewInfo(operationID, "RegisterAuthServer ok ", s.etcdSchema, strings.Join(s.etcdAddr, ","), rpcRegisterIP, s.rpcPort, s.rpcRegisterName)
	err = srv.Serve(listener)
	if err != nil {
		log.NewError(operationID, "Serve failed ", err.Error())
		return
	}
	log.NewInfo(operationID, "rpc auth ok")
}
