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

package user

import (
	"context"
	"github.com/openimsdk/tools/utils/datautil"

	cbapi "github.com/openimsdk/open-im-server/v3/pkg/callbackstruct"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/http"
	pbuser "github.com/openimsdk/protocol/user"
)

func CallbackBeforeUpdateUserInfo(ctx context.Context, callback *config.Webhooks, req *pbuser.UpdateUserInfoReq) error {
	if !callback.BeforeUpdateUserInfo.Enable {
		return nil
	}
	cbReq := &cbapi.CallbackBeforeUpdateUserInfoReq{
		CallbackCommand: cbapi.CallbackBeforeUpdateUserInfoCommand,
		UserID:          req.UserInfo.UserID,
		FaceURL:         &req.UserInfo.FaceURL,
		Nickname:        &req.UserInfo.Nickname,
	}
	resp := &cbapi.CallbackBeforeUpdateUserInfoResp{}
	if err := http.CallBackPostReturn(ctx, callback.URL, cbReq, resp, callback.BeforeUpdateUserInfo); err != nil {
		return err
	}
	datautil.NotNilReplace(&req.UserInfo.FaceURL, resp.FaceURL)
	datautil.NotNilReplace(&req.UserInfo.Ex, resp.Ex)
	datautil.NotNilReplace(&req.UserInfo.Nickname, resp.Nickname)
	return nil
}
func CallbackAfterUpdateUserInfo(ctx context.Context, callback *config.Webhooks, req *pbuser.UpdateUserInfoReq) error {
	if !callback.AfterUpdateUserInfo.Enable {
		return nil
	}
	cbReq := &cbapi.CallbackAfterUpdateUserInfoReq{
		CallbackCommand: cbapi.CallbackAfterUpdateUserInfoCommand,
		UserID:          req.UserInfo.UserID,
		FaceURL:         req.UserInfo.FaceURL,
		Nickname:        req.UserInfo.Nickname,
	}
	resp := &cbapi.CallbackAfterUpdateUserInfoResp{}
	if err := http.CallBackPostReturn(ctx, callback.URL, cbReq, resp, callback.AfterUpdateUserInfo); err != nil {
		return err
	}
	return nil
}
func CallbackBeforeUpdateUserInfoEx(ctx context.Context, callback *config.Webhooks, req *pbuser.UpdateUserInfoExReq) error {
	if !callback.BeforeUpdateUserInfoEx.Enable {
		return nil
	}
	cbReq := &cbapi.CallbackBeforeUpdateUserInfoExReq{
		CallbackCommand: cbapi.CallbackBeforeUpdateUserInfoExCommand,
		UserID:          req.UserInfo.UserID,
		FaceURL:         req.UserInfo.FaceURL,
		Nickname:        req.UserInfo.Nickname,
	}
	resp := &cbapi.CallbackBeforeUpdateUserInfoExResp{}
	if err := http.CallBackPostReturn(ctx, callback.URL, cbReq, resp, callback.BeforeUpdateUserInfoEx); err != nil {
		return err
	}
	datautil.NotNilReplace(req.UserInfo.FaceURL, resp.FaceURL)
	datautil.NotNilReplace(req.UserInfo.Ex, resp.Ex)
	datautil.NotNilReplace(req.UserInfo.Nickname, resp.Nickname)
	return nil
}
func CallbackAfterUpdateUserInfoEx(ctx context.Context, callback *config.Webhooks, req *pbuser.UpdateUserInfoExReq) error {
	if !callback.AfterUpdateUserInfoEx.Enable {
		return nil
	}
	cbReq := &cbapi.CallbackAfterUpdateUserInfoExReq{
		CallbackCommand: cbapi.CallbackAfterUpdateUserInfoExCommand,
		UserID:          req.UserInfo.UserID,
		FaceURL:         req.UserInfo.FaceURL,
		Nickname:        req.UserInfo.Nickname,
	}
	resp := &cbapi.CallbackAfterUpdateUserInfoExResp{}
	if err := http.CallBackPostReturn(ctx, callback.URL, cbReq, resp, callback.AfterUpdateUserInfoEx); err != nil {
		return err
	}
	return nil
}

func CallbackBeforeUserRegister(ctx context.Context, callback *config.Webhooks, req *pbuser.UserRegisterReq) error {
	if !callback.BeforeUserRegister.Enable {
		return nil
	}
	cbReq := &cbapi.CallbackBeforeUserRegisterReq{
		CallbackCommand: cbapi.CallbackBeforeUserRegisterCommand,
		Secret:          req.Secret,
		Users:           req.Users,
	}

	resp := &cbapi.CallbackBeforeUserRegisterResp{}
	if err := http.CallBackPostReturn(ctx, callback.URL, cbReq, resp, callback.BeforeUserRegister); err != nil {
		return err
	}
	if len(resp.Users) != 0 {
		req.Users = resp.Users
	}
	return nil
}

func CallbackAfterUserRegister(ctx context.Context, callback *config.Webhooks, req *pbuser.UserRegisterReq) error {
	if !callback.AfterUserRegister.Enable {
		return nil
	}
	cbReq := &cbapi.CallbackAfterUserRegisterReq{
		CallbackCommand: cbapi.CallbackAfterUserRegisterCommand,
		Secret:          req.Secret,
		Users:           req.Users,
	}

	resp := &cbapi.CallbackAfterUserRegisterResp{}
	if err := http.CallBackPostReturn(ctx, callback.URL, cbReq, resp, callback.AfterUserRegister); err != nil {
		return err
	}
	return nil
}
