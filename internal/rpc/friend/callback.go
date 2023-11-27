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

package friend

import (
	"context"
	"github.com/OpenIMSDK/tools/utils"

	pbfriend "github.com/OpenIMSDK/protocol/friend"
	cbapi "github.com/openimsdk/open-im-server/v3/pkg/callbackstruct"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/http"
)

func CallbackBeforeAddFriend(ctx context.Context, req *pbfriend.ApplyToAddFriendReq) error {
	if !config.Config.Callback.CallbackBeforeAddFriend.Enable {
		return nil
	}
	cbReq := &cbapi.CallbackBeforeAddFriendReq{
		CallbackCommand: cbapi.CallbackBeforeAddFriendCommand,
		FromUserID:      req.FromUserID,
		ToUserID:        req.ToUserID,
		ReqMsg:          req.ReqMsg,
	}
	resp := &cbapi.CallbackBeforeAddFriendResp{}
	if err := http.CallBackPostReturn(ctx, config.Config.Callback.CallbackUrl, cbReq, resp, config.Config.Callback.CallbackBeforeAddFriend); err != nil {
		return err
	}
	return nil
}

func CallbackBeforeSetFriendRemark(ctx context.Context, req *pbfriend.SetFriendRemarkReq) error {
	if !config.Config.Callback.CallbackBeforeSetFriendRemark.Enable {
		return nil
	}
	cbReq := &cbapi.CallbackBeforeSetFriendRemarkReq{
		CallbackCommand: cbapi.CallbackBeforeSetFriendRemark,
		OwnerUserID:     req.OwnerUserID,
		FriendUserID:    req.FriendUserID,
		Remark:          req.Remark,
	}
	resp := &cbapi.CallbackBeforeSetFriendRemarkResp{}
	if err := http.CallBackPostReturn(ctx, config.Config.Callback.CallbackUrl, cbReq, resp, config.Config.Callback.CallbackBeforeAddFriend); err != nil {
		return err
	}
	utils.NotNilReplace(&req.Remark, &resp.Remark)
	return nil
}

func CallbackAfterSetFriendRemark(ctx context.Context, req *pbfriend.SetFriendRemarkReq) error {
	if !config.Config.Callback.CallbackAfterSetFriendRemark.Enable {
		return nil
	}
	cbReq := &cbapi.CallbackAfterSetFriendRemarkReq{
		CallbackCommand: cbapi.CallbackAfterSetFriendRemark,
		OwnerUserID:     req.OwnerUserID,
		FriendUserID:    req.FriendUserID,
		Remark:          req.Remark,
	}
	resp := &cbapi.CallbackAfterSetFriendRemarkResp{}
	if err := http.CallBackPostReturn(ctx, config.Config.Callback.CallbackUrl, cbReq, resp, config.Config.Callback.CallbackBeforeAddFriend); err != nil {
		return err
	}
	return nil
}
