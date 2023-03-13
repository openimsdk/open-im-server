package auth

import (
	"OpenIM/internal/common/check"
	"OpenIM/pkg/common/config"
	"OpenIM/pkg/common/constant"
	"OpenIM/pkg/common/db/cache"
	"OpenIM/pkg/common/db/controller"
	"OpenIM/pkg/common/log"
	"OpenIM/pkg/common/tokenverify"
	"OpenIM/pkg/common/tracelog"
	"OpenIM/pkg/discoveryregistry"
	"OpenIM/pkg/errs"
	pbAuth "OpenIM/pkg/proto/auth"
	"OpenIM/pkg/proto/msggateway"
	"OpenIM/pkg/utils"
	"context"
	"google.golang.org/grpc"
)

type authServer struct {
	authDatabase   controller.AuthDatabase
	userCheck      *check.UserCheck
	RegisterCenter discoveryregistry.SvcDiscoveryRegistry
}

func Start(client discoveryregistry.SvcDiscoveryRegistry, server *grpc.Server) error {
	rdb, err := cache.NewRedis()
	if err != nil {
		return err
	}
	pbAuth.RegisterAuthServer(server, &authServer{
		userCheck:      check.NewUserCheck(client),
		RegisterCenter: client,
		authDatabase:   controller.NewAuthDatabase(cache.NewCacheModel(rdb), config.Config.TokenPolicy.AccessSecret, config.Config.TokenPolicy.AccessExpire),
	})
	return nil
}

func (s *authServer) UserToken(ctx context.Context, req *pbAuth.UserTokenReq) (*pbAuth.UserTokenResp, error) {
	log.Info("", "rpc come UserToken")
	resp := pbAuth.UserTokenResp{}
	if _, err := s.userCheck.GetUserInfo(ctx, req.UserID); err != nil {
		return nil, err
	}
	token, err := s.authDatabase.CreateToken(ctx, req.UserID, constant.PlatformIDToName(int(req.PlatformID)))
	if err != nil {
		return nil, err
	}
	resp.Token = token
	resp.ExpireTimeSeconds = config.Config.TokenPolicy.AccessExpire
	return &resp, nil
}

func (s *authServer) parseToken(ctx context.Context, tokensString string) (claims *tokenverify.Claims, err error) {
	claims, err = tokenverify.GetClaimFromToken(tokensString)
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	m, err := s.authDatabase.GetTokensWithoutError(ctx, claims.UID, claims.Platform)
	if err != nil {
		return nil, err
	}
	if len(m) == 0 {
		return nil, errs.ErrTokenNotExist.Wrap()
	}
	if v, ok := m[tokensString]; ok {
		switch v {
		case constant.NormalToken:
			return claims, nil
		case constant.KickedToken:
			return nil, errs.ErrTokenKicked.Wrap()
		default:
			return nil, utils.Wrap(errs.ErrTokenUnknown, "")
		}
	}
	return nil, errs.ErrTokenNotExist.Wrap()
}

func (s *authServer) ParseToken(ctx context.Context, req *pbAuth.ParseTokenReq) (resp *pbAuth.ParseTokenResp, err error) {
	resp = &pbAuth.ParseTokenResp{}
	claims, err := s.parseToken(ctx, req.Token)
	if err != nil {
		return nil, err
	}
	resp.UserID = claims.UID
	resp.Platform = claims.Platform
	resp.ExpireTimeSeconds = claims.ExpiresAt.Unix()
	return resp, nil
}

func (s *authServer) ForceLogout(ctx context.Context, req *pbAuth.ForceLogoutReq) (*pbAuth.ForceLogoutResp, error) {
	resp := pbAuth.ForceLogoutResp{}
	if err := tokenverify.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	if err := s.forceKickOff(ctx, req.UserID, req.PlatformID, tracelog.GetOperationID(ctx)); err != nil {
		return nil, err
	}
	return &resp, nil
}

func (s *authServer) forceKickOff(ctx context.Context, userID string, platformID int32, operationID string) error {
	conns, err := s.RegisterCenter.GetConns(config.Config.RpcRegisterName.OpenImMessageGatewayName)
	if err != nil {
		return err
	}
	for _, v := range conns {
		client := msggateway.NewMsgGatewayClient(v)
		kickReq := &msggateway.KickUserOfflineReq{KickUserIDList: []string{userID}, PlatformID: platformID}
		log.NewInfo(operationID, "KickUserOffline ", client, kickReq.String())
		_, err := client.KickUserOffline(ctx, kickReq)
		return utils.Wrap(err, "")
	}
	return errs.ErrInternalServer.Wrap()
}
