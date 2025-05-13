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

	"github.com/golang-jwt/jwt/v4"
	"github.com/openimsdk/open-im-server/v3/pkg/common/servererrs"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/tools/mcontext"
	"github.com/openimsdk/tools/utils/datautil"
)

func Secret(secret string) jwt.Keyfunc {
	return func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	}
}

func CheckAccessV3(ctx context.Context, ownerUserID string, imAdminUserID []string) (err error) {
	opUserID := mcontext.GetOpUserID(ctx)
	if datautil.Contain(opUserID, imAdminUserID...) {
		return nil
	}
	if opUserID == ownerUserID {
		return nil
	}
	return servererrs.ErrNoPermission.WrapMsg("ownerUserID", ownerUserID)
}

func IsAppManagerUid(ctx context.Context, imAdminUserID []string) bool {
	return datautil.Contain(mcontext.GetOpUserID(ctx), imAdminUserID...)
}

func CheckAdmin(ctx context.Context) error {
	if IsAdmin(ctx) {
		return nil
	}
	return servererrs.ErrNoPermission.WrapMsg(fmt.Sprintf("user %s is not admin userID", mcontext.GetOpUserID(ctx)))
}

func IsManagerUserID(opUserID string, imAdminUserID []string) bool {
	return datautil.Contain(opUserID, imAdminUserID...)
}

func CheckSystemAccount(ctx context.Context, level int32) bool {
	return level >= constant.AppAdmin
}

type ctxAuthKey struct{}

func WithIMAdminUserIDs(ctx context.Context, imAdminUserID []string) context.Context {
	return context.WithValue(ctx, ctxAuthKey{}, imAdminUserID)
}

func GetIMAdminUserIDs(ctx context.Context) []string {
	imAdminUserID, _ := ctx.Value(ctxAuthKey{}).([]string)
	return imAdminUserID
}

func IsAdmin(ctx context.Context) bool {
	return datautil.Contain(mcontext.GetOpUserID(ctx), GetIMAdminUserIDs(ctx)...)
}

func CheckAccess(ctx context.Context, ownerUserID string) error {
	opUserID := mcontext.GetOpUserID(ctx)
	if opUserID == ownerUserID {
		return nil
	}
	if !datautil.Contain(mcontext.GetOpUserID(ctx), GetIMAdminUserIDs(ctx)...) {
		return servererrs.ErrNoPermission.WrapMsg("ownerUserID", ownerUserID)
	}
	return nil
}
