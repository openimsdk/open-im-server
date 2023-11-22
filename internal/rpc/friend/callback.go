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

	"github.com/OpenIMSDK/protocol/constant"
	pbfriend "github.com/OpenIMSDK/protocol/friend"
	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/mcontext"

	cbapi "github.com/openimsdk/open-im-server/v3/pkg/callbackstruct"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/http"
)

func CallbackBeforeAddFriend(ctx context.Context, req *pbfriend.ApplyToAddFriendReq) error {
	if !config.Config.Callback.CallbackBeforeAddFriend.Enable {
		return nil
	}
	cbReq := &cbapi.CallbackBeforeAddFriendReq{
		CallbackCommand: constant.CallbackBeforeAddFriendCommand,
		FromUserID:      req.FromUserID,
		ToUserID:        req.ToUserID,
		ReqMsg:          req.ReqMsg,
		OperationID:     mcontext.GetOperationID(ctx),
	}
	resp := &cbapi.CallbackBeforeAddFriendResp{}
	if err := http.CallBackPostReturn(ctx, config.Config.Callback.CallbackUrl, cbReq, resp, config.Config.Callback.CallbackBeforeAddFriend); err != nil {
		if err == errs.ErrCallbackContinue {
			return nil
		}
		return err
	}
	return nil
}
func CallbackBeforeAddBlack(ctx context.Context, req *pbfriend.AddBlackReq) error {
	if !config.Config.Callback.CallbackBeforeAddBlack.Enable {
		return nil
	}
	cbReq := &cbapi.CallbackBeforeAddBlackReq{
		CallbackCommand: cbapi.CallbackBeforeAddBlackCommand,
		OwnerUserID:     req.OwnerUserID,
		BlackUserID:     req.BlackUserID,
		OperationID:     mcontext.GetOperationID(ctx),
	}
	resp := &cbapi.CallbackBeforeAddBlackResp{}
	if err := http.CallBackPostReturn(ctx, config.Config.Callback.CallbackUrl, cbReq, resp, config.Config.Callback.CallbackBeforeAddBlack); err != nil {
		if err == errs.ErrCallbackContinue {
			return nil
		}
		return err
	}
	utils.StructFieldNotNilReplace(req, resp)
	return nil
}
func CallbackAfterAddFriend(ctx context.Context, req *pbfriend.ApplyToAddFriendReq) error {
	if !config.Config.Callback.CallbackAfterAddFriend.Enable {
		return nil
	}
	cbReq := &cbapi.CallbackAfterAddFriendReq{
		CallbackCommand: cbapi.CallbackAfterAddFriendCommand,
		FromUserID:      req.FromUserID,
		ToUserID:        req.ToUserID,
		ReqMsg:          req.ReqMsg,
		OperationID:     mcontext.GetOperationID(ctx),
	}
	resp := &cbapi.CallbackAfterAddFriendResp{}
	if err := http.CallBackPostReturn(ctx, config.Config.Callback.CallbackUrl, cbReq, resp, config.Config.Callback.CallbackAfterAddFriend); err != nil {
		if err == errs.ErrCallbackContinue {
			return nil
		}
		return err
	}
	utils.StructFieldNotNilReplace(req, resp)
	return nil
}
func CallbackBeforeAddFriendAgree(ctx context.Context, req *pbfriend.RespondFriendApplyReq) error {
	if !config.Config.Callback.CallbackBeforeAddFriendAgree.Enable {
		return nil
	}
	cbReq := &cbapi.CallbackBeforeAddFriendAgreeReq{
		CallbackCommand: cbapi.CallbackBeforeAddFriendAgreeCommand,
		FromUserID:      req.FromUserID,
		ToUserID:        req.ToUserID,
		HandleMsg:       req.HandleMsg,
		HandleResult:    req.HandleResult,
		OperationID:     mcontext.GetOperationID(ctx),
	}
	resp := &cbapi.CallbackBeforeAddFriendAgreeResp{}
	if err := http.CallBackPostReturn(ctx, config.Config.Callback.CallbackUrl, cbReq, resp, config.Config.Callback.CallbackBeforeAddFriendAgree); err != nil {
		if err == errs.ErrCallbackContinue {
			return nil
		}
		return err
	}
	utils.StructFieldNotNilReplace(req, resp)
	return nil
}
