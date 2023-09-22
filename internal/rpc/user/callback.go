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

	"github.com/OpenIMSDK/protocol/constant"
	pbuser "github.com/OpenIMSDK/protocol/user"
	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/mcontext"
	"github.com/OpenIMSDK/tools/utils"

	cbapi "github.com/openimsdk/open-im-server/v3/pkg/callbackstruct"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/http"
)

func CallbackBeforeUpdateUserInfo(ctx context.Context, req *pbuser.UpdateUserInfoReq) error {
	if !config.Config.Callback.CallbackBeforeUpdateUserInfo.Enable {
		return nil
	}
	cbReq := &cbapi.CallbackBeforeUpdateUserInfoReq{
		CallbackCommand: constant.CallbackBeforeUpdateUserInfoCommand,
		OperationID:     mcontext.GetOperationID(ctx),
		UserID:          req.UserInfo.UserID,
		FaceURL:         &req.UserInfo.FaceURL,
		Nickname:        &req.UserInfo.Nickname,
	}
	resp := &cbapi.CallbackBeforeUpdateUserInfoResp{}
	if err := http.CallBackPostReturn(ctx, config.Config.Callback.CallbackUrl, cbReq, resp, config.Config.Callback.CallbackBeforeUpdateUserInfo); err != nil {
		if err == errs.ErrCallbackContinue {
			return nil
		}
		return err
	}
	utils.NotNilReplace(&req.UserInfo.FaceURL, resp.FaceURL)
	utils.NotNilReplace(&req.UserInfo.Ex, resp.Ex)
	utils.NotNilReplace(&req.UserInfo.Nickname, resp.Nickname)
	return nil
}
