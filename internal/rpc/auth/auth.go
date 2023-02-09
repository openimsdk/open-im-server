package auth

import (
	"Open_IM/internal/common/check"
	"Open_IM/internal/common/rpc_server"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db/controller"
	"Open_IM/pkg/common/log"
	promePkg "Open_IM/pkg/common/prometheus"
	"Open_IM/pkg/common/tokenverify"
	"Open_IM/pkg/common/tracelog"
	pbAuth "Open_IM/pkg/proto/auth"
	pbRelay "Open_IM/pkg/proto/relay"
	"Open_IM/pkg/utils"
	"context"
	grpcPrometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"google.golang.org/grpc"
)

func NewRpcAuthServer(port int) *rpcAuth {
	r, err := rpc_server.NewRpcServer(config.Config.RpcRegisterIP, port, config.Config.RpcRegisterName.OpenImAuthName, config.Config.Zookeeper.ZkAddr, config.Config.Zookeeper.Schema)
	if err != nil {
		panic(err)
	}
	return &rpcAuth{
		RpcServer: r,
	}
}

func (s *rpcAuth) Run() {
	operationID := utils.OperationIDGenerator()
	log.NewInfo(operationID, "rpc auth start...")
	listener, address, err := rpc_server.GetTcpListen(config.Config.ListenIP, s.Port)
	if err != nil {
		panic(err)
	}
	log.NewInfo(operationID, "listen network success ", listener, address)
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
	pbAuth.RegisterAuthServer(srv, s)
	err = srv.Serve(listener)
	if err != nil {
		panic(err)
	}
	log.NewInfo(operationID, "rpc auth ok")
}

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

func (s *rpcAuth) parseToken(ctx context.Context, tokensString, operationID string) (claims *tokenverify.Claims, err error) {
	claims, err = tokenverify.GetClaimFromToken(tokensString)
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
	if err := tokenverify.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	if err := s.forceKickOff(ctx, req.UserID, req.PlatformID, tracelog.GetOperationID(ctx)); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *rpcAuth) forceKickOff(ctx context.Context, userID string, platformID int32, operationID string) error {
	grpcCons, err := s.RegisterCenter.GetConns(config.Config.RpcRegisterName.OpenImRelayName)
	if err != nil {
		return err
	}
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
	*rpc_server.RpcServer
	controller.AuthInterface
}
