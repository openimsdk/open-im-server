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
	"github.com/openimsdk/open-im-server/v3/pkg/common/servererrs"
	"github.com/openimsdk/tools/utils/datautil"

	"github.com/golang-jwt/jwt/v4"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/tools/mcontext"
	"github.com/openimsdk/tools/tokenverify"
)

func Secret(secret string) jwt.Keyfunc {
	return func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	}
}

func CheckAccessV3(ctx context.Context, ownerUserID string, manager *config.Manager, imAdmin *config.IMAdmin) (err error) {
	opUserID := mcontext.GetOpUserID(ctx)
	if len(manager.UserID) > 0 && datautil.Contain(opUserID, manager.UserID...) {
		return nil
	}
	if datautil.Contain(opUserID, imAdmin.UserID...) {
		return nil
	}
	if opUserID == ownerUserID {
		return nil
	}
	return servererrs.ErrNoPermission.WrapMsg("ownerUserID", ownerUserID)
}

func IsAppManagerUid(ctx context.Context, manager *config.Manager, imAdmin *config.IMAdmin) bool {
	return (len(manager.UserID) > 0 && datautil.Contain(mcontext.GetOpUserID(ctx), manager.UserID...)) ||
		datautil.Contain(mcontext.GetOpUserID(ctx), imAdmin.UserID...)
}

func CheckAdmin(ctx context.Context, manager *config.Manager, imAdmin *config.IMAdmin) error {
	if len(manager.UserID) > 0 && datautil.Contain(mcontext.GetOpUserID(ctx), manager.UserID...) {
		return nil
	}
	if datautil.Contain(mcontext.GetOpUserID(ctx), imAdmin.UserID...) {
		return nil
	}
	return servererrs.ErrNoPermission.WrapMsg(fmt.Sprintf("user %s is not admin userID", mcontext.GetOpUserID(ctx)))
}

func CheckIMAdmin(ctx context.Context, config *config.GlobalConfig) error {
	if datautil.Contain(mcontext.GetOpUserID(ctx), config.IMAdmin.UserID...) {
		return nil
	}
	if len(config.Manager.UserID) > 0 && datautil.Contain(mcontext.GetOpUserID(ctx), config.Manager.UserID...) {
		return nil
	}
	return servererrs.ErrNoPermission.WrapMsg(fmt.Sprintf("user %s is not CheckIMAdmin userID", mcontext.GetOpUserID(ctx)))
}

func ParseRedisInterfaceToken(redisToken any, secret string) (*tokenverify.Claims, error) {
	return tokenverify.GetClaimFromToken(string(redisToken.([]uint8)), Secret(secret))
}

func IsManagerUserID(opUserID string, manager *config.Manager, imAdmin *config.IMAdmin) bool {
	return (len(manager.UserID) > 0 && datautil.Contain(opUserID, manager.UserID...)) || datautil.Contain(opUserID, imAdmin.UserID...)
}

func WsVerifyToken(token, userID, secret string, platformID int) error {
	claim, err := tokenverify.GetClaimFromToken(token, Secret(secret))
	if err != nil {
		return err
	}
	if claim.UserID != userID {
		return servererrs.ErrTokenInvalid.WrapMsg(fmt.Sprintf("token uid %s != userID %s", claim.UserID, userID))
	}
	if claim.PlatformID != platformID {
		return servererrs.ErrTokenInvalid.WrapMsg(fmt.Sprintf("token platform %d != %d", claim.PlatformID, platformID))
	}
	return nil
}
