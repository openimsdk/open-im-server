// Copyright © 2023 OpenIM. All rights reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package auth

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/pkg/authverify"

	"google.golang.org/grpc"

	pbauth "github.com/OpenIMSDK/protocol/auth"
	"github.com/OpenIMSDK/protocol/constant"
	"github.com/OpenIMSDK/protocol/msggateway"
	"github.com/OpenIMSDK/tools/discoveryregistry"
	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/log"
	"github.com/OpenIMSDK/tools/mcontext"
	"github.com/OpenIMSDK/tools/tokenverify"
	"github.com/OpenIMSDK/tools/utils"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/controller"
	"github.com/openimsdk/open-im-server/v3/pkg/common/prommetrics"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient"
)

type authServer struct {
	authDatabase   controller.AuthDatabase
	userRpcClient  *rpcclient.UserRpcClient
	RegisterCenter discoveryregistry.SvcDiscoveryRegistry
}

func Start(client discoveryregistry.SvcDiscoveryRegistry, server *grpc.Server) error {
	rdb, err := cache.NewRedis()
	if err != nil {
		return err
	}
	userRpcClient := rpcclient.NewUserRpcClient(client)
	pbauth.RegisterAuthServer(server, &authServer{
		userRpcClient:  &userRpcClient,
		RegisterCenter: client,
		authDatabase: controller.NewAuthDatabase(
			cache.NewMsgCacheModel(rdb),
			config.Config.Secret,
			config.Config.TokenPolicy.Expire,
		),
	})
	return nil
}

func (s *authServer) UserToken(ctx context.Context, req *pbauth.UserTokenReq) (*pbauth.UserTokenResp, error) {
	resp := pbauth.UserTokenResp{}
	if req.Secret != config.Config.Secret {
		return nil, errs.ErrNoPermission.Wrap("secret invalid")
	}
	if _, err := s.userRpcClient.GetUserInfo(ctx, req.UserID); err != nil {
		return nil, err
	}
	token, err := s.authDatabase.CreateToken(ctx, req.UserID, int(req.PlatformID))
	if err != nil {
		return nil, err
	}
	prommetrics.UserLoginCounter.Inc()
	resp.Token = token
	resp.ExpireTimeSeconds = config.Config.TokenPolicy.Expire * 24 * 60 * 60
	return &resp, nil
}

func (s *authServer) parseToken(ctx context.Context, tokensString string) (claims *tokenverify.Claims, err error) {
	claims, err = tokenverify.GetClaimFromToken(tokensString, authverify.Secret())
	if err != nil {
		return nil, utils.Wrap(err, "")
	}
	m, err := s.authDatabase.GetTokensWithoutError(ctx, claims.UserID, claims.PlatformID)
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

func (s *authServer) ParseToken(
	ctx context.Context,
	req *pbauth.ParseTokenReq,
) (resp *pbauth.ParseTokenResp, err error) {
	resp = &pbauth.ParseTokenResp{}
	claims, err := s.parseToken(ctx, req.Token)
	if err != nil {
		return nil, err
	}
	resp.UserID = claims.UserID
	resp.Platform = constant.PlatformIDToName(claims.PlatformID)
	resp.ExpireTimeSeconds = claims.ExpiresAt.Unix()
	return resp, nil
}

func (s *authServer) ForceLogout(ctx context.Context, req *pbauth.ForceLogoutReq) (*pbauth.ForceLogoutResp, error) {
	if err := authverify.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	if err := s.forceKickOff(ctx, req.UserID, req.PlatformID, mcontext.GetOperationID(ctx)); err != nil {
		return nil, err
	}
	return &pbauth.ForceLogoutResp{}, nil
}

func (s *authServer) forceKickOff(ctx context.Context, userID string, platformID int32, operationID string) error {
	conns, err := s.RegisterCenter.GetConns(ctx, config.Config.RpcRegisterName.OpenImMessageGatewayName)
	if err != nil {
		return err
	}
	for _, v := range conns {
		log.ZDebug(ctx, "forceKickOff", "conn", v.Target())
	}
	for _, v := range conns {
		client := msggateway.NewMsgGatewayClient(v)
		kickReq := &msggateway.KickUserOfflineReq{KickUserIDList: []string{userID}, PlatformID: platformID}
		_, err := client.KickUserOffline(ctx, kickReq)
		if err != nil {
			log.ZError(ctx, "forceKickOff", err, "kickReq", kickReq)
		}
	}
	return nil
}
