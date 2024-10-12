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

package controller

import (
	"context"
	"github.com/openimsdk/tools/log"

	"github.com/golang-jwt/jwt/v4"
	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/tokenverify"
)

type AuthDatabase interface {
	// If the result is empty, no error is returned.
	GetTokensWithoutError(ctx context.Context, userID string, platformID int) (map[string]int, error)
	// Create token
	CreateToken(ctx context.Context, userID string, platformID int) (string, error)

	SetTokenMapByUidPid(ctx context.Context, userID string, platformID int, m map[string]int) error
}

type authDatabase struct {
	cache            cache.TokenModel
	accessSecret     string
	accessExpire     int64
	multiLoginPolicy int
}

func NewAuthDatabase(cache cache.TokenModel, accessSecret string, accessExpire int64, policy int) AuthDatabase {
	return &authDatabase{cache: cache, accessSecret: accessSecret, accessExpire: accessExpire, multiLoginPolicy: policy}
}

// If the result is empty.
func (a *authDatabase) GetTokensWithoutError(ctx context.Context, userID string, platformID int) (map[string]int, error) {
	return a.cache.GetTokensWithoutError(ctx, userID, platformID)
}

func (a *authDatabase) SetTokenMapByUidPid(ctx context.Context, userID string, platformID int, m map[string]int) error {
	return a.cache.SetTokenMapByUidPid(ctx, userID, platformID, m)
}

// Create Token.
func (a *authDatabase) CreateToken(ctx context.Context, userID string, platformID int) (string, error) {
	// todo: get all platform token
	tokens, err := a.cache.GetTokensWithoutError(ctx, userID, platformID)
	if err != nil {
		return "", err
	}
	var deleteTokenKey []string
	var kickedTokenKey []string
	for k, v := range tokens {
		t, err := tokenverify.GetClaimFromToken(k, authverify.Secret(a.accessSecret))
		if err != nil || v != constant.NormalToken {
			deleteTokenKey = append(deleteTokenKey, k)
		} else if a.checkKickToken(ctx, platformID, t) {
			kickedTokenKey = append(kickedTokenKey, k)
		}
	}
	if len(deleteTokenKey) != 0 {
		err = a.cache.DeleteTokenByUidPid(ctx, userID, platformID, deleteTokenKey)
		if err != nil {
			return "", err
		}
	}

	const adminTokenMaxNum = 30
	if platformID == constant.AdminPlatformID {
		if len(kickedTokenKey) > adminTokenMaxNum {
			kickedTokenKey = kickedTokenKey[:len(kickedTokenKey)-adminTokenMaxNum]
		} else {
			kickedTokenKey = nil
		}
	}

	if len(kickedTokenKey) != 0 {
		for _, k := range kickedTokenKey {
			err := a.cache.SetTokenFlagEx(ctx, userID, platformID, k, constant.KickedToken)
			if err != nil {
				return "", err
			}
			log.ZDebug(ctx, "kicked token in create token", "token", k)
		}
	}

	claims := tokenverify.BuildClaims(userID, platformID, a.accessExpire)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(a.accessSecret))
	if err != nil {
		return "", errs.WrapMsg(err, "token.SignedString")
	}

	if err = a.cache.SetTokenFlagEx(ctx, userID, platformID, tokenString, constant.NormalToken); err != nil {
		return "", err
	}
	return tokenString, nil
}

func (a *authDatabase) checkKickToken(ctx context.Context, platformID int, token *tokenverify.Claims) bool {
	switch a.multiLoginPolicy {
	case constant.DefalutNotKick:
		return false
	case constant.PCAndOther:
		if constant.PlatformIDToClass(platformID) == constant.TerminalPC ||
			constant.PlatformIDToClass(token.PlatformID) == constant.TerminalPC {
			return false
		}
		return true
	case constant.AllLoginButSameTermKick:
		if platformID == token.PlatformID {
			return true
		}
		return false
	default:
		return false
	}
}
