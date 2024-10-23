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

package rpcclient

import (
	"context"

	"github.com/openimsdk/protocol/auth"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/system/program"
	"google.golang.org/grpc"
)

func NewAuth(discov discovery.SvcDiscoveryRegistry, rpcRegisterName string) *Auth {
	conn, err := discov.GetConn(context.Background(), rpcRegisterName)
	if err != nil {
		program.ExitWithError(err)
	}
	client := auth.NewAuthClient(conn)
	return &Auth{discov: discov, conn: conn, Client: client}
}

type Auth struct {
	conn   grpc.ClientConnInterface
	Client auth.AuthClient
	discov discovery.SvcDiscoveryRegistry
}

func (a *Auth) ParseToken(ctx context.Context, token string) (*auth.ParseTokenResp, error) {
	req := auth.ParseTokenReq{
		Token: token,
	}
	resp, err := a.Client.ParseToken(ctx, &req)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (a *Auth) InvalidateToken(ctx context.Context, preservedToken, userID string, platformID int) (*auth.InvalidateTokenResp, error) {
	req := auth.InvalidateTokenReq{
		PreservedToken: preservedToken,
		UserID:         userID,
		PlatformID:     int32(platformID),
	}
	resp, err := a.Client.InvalidateToken(ctx, &req)
	if err != nil {
		return nil, err
	}
	return resp, err
}

func (a *Auth) KickTokens(ctx context.Context, tokens []string) (*auth.KickTokensResp, error) {
	req := auth.KickTokensReq{
		Tokens: tokens,
	}
	resp, err := a.Client.KickTokens(ctx, &req)
	if err != nil {
		return nil, err
	}
	return resp, err
}
