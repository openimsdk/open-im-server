package auth

import (
	"Open_IM/internal/common/check"
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	"Open_IM/pkg/common/db/cache"
	"Open_IM/pkg/common/db/controller"
	"Open_IM/pkg/common/db/relation"
	relationTb "Open_IM/pkg/common/db/table/relation"
	"Open_IM/pkg/common/log"
	"Open_IM/pkg/common/tokenverify"
	"Open_IM/pkg/common/tracelog"
	discoveryRegistry "Open_IM/pkg/discoveryregistry"
	pbAuth "Open_IM/pkg/proto/auth"
	pbRelay "Open_IM/pkg/proto/relay"
	"Open_IM/pkg/utils"
	"context"
	"github.com/OpenIMSDK/openKeeper"
	"google.golang.org/grpc"
)

func Start(client *openKeeper.ZkClient, server *grpc.Server) error {
	mysql, err := relation.NewGormDB()
	if err != nil {
		return err
	}
	if err := mysql.AutoMigrate(&relationTb.FriendModel{}, &relationTb.FriendRequestModel{}, &relationTb.BlackModel{}); err != nil {
		return err
	}
	redis, err := cache.NewRedis()
	if err != nil {
		return err
	}
	pbAuth.RegisterAuthServer(server, &authServer{
		userCheck:      check.NewUserCheck(client),
		RegisterCenter: client,
		AuthInterface:  controller.NewAuthController(redis.GetClient(), config.Config.TokenPolicy.AccessSecret, config.Config.TokenPolicy.AccessExpire),
	})
	return nil
}

func (s *authServer) UserToken(ctx context.Context, req *pbAuth.UserTokenReq) (*pbAuth.UserTokenResp, error) {
	resp := pbAuth.UserTokenResp{}
	if _, err := s.userCheck.GetUsersInfo(ctx, req.UserID); err != nil {
		return nil, err
	}
	token, err := s.CreateToken(ctx, req.UserID, constant.PlatformIDToName(int(req.PlatformID)))
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
	m, err := s.GetTokensWithoutError(ctx, claims.UID, claims.Platform)
	if err != nil {
		return nil, err
	}
	if len(m) == 0 {
		return nil, constant.ErrTokenNotExist.Wrap()
	}
	if v, ok := m[tokensString]; ok {
		switch v {
		case constant.NormalToken:
			return claims, nil
		case constant.KickedToken:
			return nil, constant.ErrTokenKicked.Wrap()
		default:
			return nil, utils.Wrap(constant.ErrTokenUnknown, "")
		}
	}
	return nil, constant.ErrTokenNotExist.Wrap()
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

type authServer struct {
	controller.AuthInterface
	userCheck      *check.UserCheck
	RegisterCenter discoveryRegistry.SvcDiscoveryRegistry
}
