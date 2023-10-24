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

package authverify

import (
	"context"
	"fmt"

	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/mcontext"
	"github.com/OpenIMSDK/tools/tokenverify"
	"github.com/OpenIMSDK/tools/utils"
	"github.com/golang-jwt/jwt/v4"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
)

func Secret() jwt.Keyfunc {
	return func(token *jwt.Token) (interface{}, error) {
		return []byte(config.Config.Secret), nil
	}
}

func CheckAccessV3(ctx context.Context, ownerUserID string) (err error) {
	opUserID := mcontext.GetOpUserID(ctx)
	if utils.IsContain(opUserID, config.Config.Manager.UserID) {
		return nil
	}
	if opUserID == ownerUserID {
		return nil
	}
	return errs.ErrNoPermission.Wrap(utils.GetSelfFuncName())
}

func IsAppManagerUid(ctx context.Context) bool {
	return utils.IsContain(mcontext.GetOpUserID(ctx), config.Config.Manager.UserID)
}

func CheckAdmin(ctx context.Context) error {
	if utils.IsContain(mcontext.GetOpUserID(ctx), config.Config.Manager.UserID) {
		return nil
	}
	return errs.ErrNoPermission.Wrap(fmt.Sprintf("user %s is not admin userID", mcontext.GetOpUserID(ctx)))
}

func ParseRedisInterfaceToken(redisToken interface{}) (*tokenverify.Claims, error) {
	return tokenverify.GetClaimFromToken(string(redisToken.([]uint8)), Secret())
}

func IsManagerUserID(opUserID string) bool {
	return utils.IsContain(opUserID, config.Config.Manager.UserID)
}

func WsVerifyToken(token, userID string, platformID int) error {
	claim, err := tokenverify.GetClaimFromToken(token, Secret())
	if err != nil {
		return err
	}
	if claim.UserID != userID {
		return errs.ErrTokenInvalid.Wrap(fmt.Sprintf("token uid %s != userID %s", claim.UserID, userID))
	}
	if claim.PlatformID != platformID {
		return errs.ErrTokenInvalid.Wrap(fmt.Sprintf("token platform %d != %d", claim.PlatformID, platformID))
	}
	return nil
}
