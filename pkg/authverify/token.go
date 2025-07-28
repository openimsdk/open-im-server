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

func CheckAdmin(ctx context.Context) error {
	if IsAdmin(ctx) {
		return nil
	}
	return servererrs.ErrNoPermission.WrapMsg(fmt.Sprintf("user %s is not admin userID", mcontext.GetOpUserID(ctx)))
}

//func IsManagerUserID(opUserID string, imAdminUserID []string) bool {
//	return datautil.Contain(opUserID, imAdminUserID...)
//}

func CheckUserIsAdmin(ctx context.Context, userID string) bool {
	return datautil.Contain(userID, GetIMAdminUserIDs(ctx)...)
}

func CheckSystemAccount(ctx context.Context, level int32) bool {
	return level >= constant.AppAdmin
}

const (
	CtxAdminUserIDsKey = "CtxAdminUserIDsKey"
)

func WithIMAdminUserIDs(ctx context.Context, imAdminUserID []string) context.Context {
	return context.WithValue(ctx, CtxAdminUserIDsKey, imAdminUserID)
}

func GetIMAdminUserIDs(ctx context.Context) []string {
	imAdminUserID, _ := ctx.Value(CtxAdminUserIDsKey).([]string)
	return imAdminUserID
}

func IsAdmin(ctx context.Context) bool {
	return IsTempAdmin(ctx) || IsSystemAdmin(ctx)
}

func CheckAccess(ctx context.Context, ownerUserID string) error {
	if mcontext.GetOpUserID(ctx) == ownerUserID {
		return nil
	}
	if IsAdmin(ctx) {
		return nil
	}
	return servererrs.ErrNoPermission.WrapMsg("ownerUserID", ownerUserID)
}

func CheckAccessIn(ctx context.Context, ownerUserIDs ...string) error {
	opUserID := mcontext.GetOpUserID(ctx)
	for _, userID := range ownerUserIDs {
		if opUserID == userID {
			return nil
		}
	}
	if IsAdmin(ctx) {
		return nil
	}
	return servererrs.ErrNoPermission.WrapMsg("opUser in ownerUserIDs")
}

var tempAdminValue = []string{"1"}

const ctxTempAdminKey = "ctxImTempAdminKey"

func WithTempAdmin(ctx context.Context) context.Context {
	keys, _ := ctx.Value(constant.RpcCustomHeader).([]string)
	if datautil.Contain(ctxTempAdminKey, keys...) {
		return ctx
	}
	if len(keys) > 0 {
		temp := make([]string, 0, len(keys)+1)
		temp = append(temp, keys...)
		keys = append(temp, ctxTempAdminKey)
	} else {
		keys = []string{ctxTempAdminKey}
	}
	ctx = context.WithValue(ctx, constant.RpcCustomHeader, keys)
	return context.WithValue(ctx, ctxTempAdminKey, tempAdminValue)
}

func IsTempAdmin(ctx context.Context) bool {
	values, _ := ctx.Value(ctxTempAdminKey).([]string)
	return datautil.Equal(tempAdminValue, values)
}

func IsSystemAdmin(ctx context.Context) bool {
	return datautil.Contain(mcontext.GetOpUserID(ctx), GetIMAdminUserIDs(ctx)...)
}
