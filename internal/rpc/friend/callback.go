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

	pbfriend "github.com/OpenIMSDK/protocol/friend"
	"github.com/OpenIMSDK/tools/utils"

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
		Ex:              req.Ex,
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
func CallbackBeforeAddBlack(ctx context.Context, req *pbfriend.AddBlackReq) error {
	if !config.Config.Callback.CallbackBeforeAddBlack.Enable {
		return nil
	}
	cbReq := &cbapi.CallbackBeforeAddBlackReq{
		CallbackCommand: cbapi.CallbackBeforeAddBlackCommand,
		OwnerUserID:     req.OwnerUserID,
		BlackUserID:     req.BlackUserID,
	}
	resp := &cbapi.CallbackBeforeAddBlackResp{}
	if err := http.CallBackPostReturn(ctx, config.Config.Callback.CallbackUrl, cbReq, resp, config.Config.Callback.CallbackBeforeAddBlack); err != nil {
		return err
	}
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
	}
	resp := &cbapi.CallbackAfterAddFriendResp{}
	if err := http.CallBackPostReturn(ctx, config.Config.Callback.CallbackUrl, cbReq, resp, config.Config.Callback.CallbackAfterAddFriend); err != nil {
		return err
	}

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
	}
	resp := &cbapi.CallbackBeforeAddFriendAgreeResp{}
	if err := http.CallBackPostReturn(ctx, config.Config.Callback.CallbackUrl, cbReq, resp, config.Config.Callback.CallbackBeforeAddFriendAgree); err != nil {
		return err
	}
	return nil
}
func CallbackAfterDeleteFriend(ctx context.Context, req *pbfriend.DeleteFriendReq) error {
	if !config.Config.Callback.CallbackAfterDeleteFriend.Enable {
		return nil
	}
	cbReq := &cbapi.CallbackAfterDeleteFriendReq{
		CallbackCommand: cbapi.CallbackAfterDeleteFriendCommand,
		OwnerUserID:     req.OwnerUserID,
		FriendUserID:    req.FriendUserID,
	}
	resp := &cbapi.CallbackAfterDeleteFriendResp{}
	if err := http.CallBackPostReturn(ctx, config.Config.Callback.CallbackUrl, cbReq, resp, config.Config.Callback.CallbackAfterDeleteFriend); err != nil {
		return err
	}
	return nil
}
func CallbackBeforeImportFriends(ctx context.Context, req *pbfriend.ImportFriendReq) error {
	if !config.Config.Callback.CallbackBeforeImportFriends.Enable {
		return nil
	}
	cbReq := &cbapi.CallbackBeforeImportFriendsReq{
		CallbackCommand: cbapi.CallbackBeforeImportFriendsCommand,
		OwnerUserID:     req.OwnerUserID,
		FriendUserIDs:   req.FriendUserIDs,
	}
	resp := &cbapi.CallbackBeforeImportFriendsResp{}
	if err := http.CallBackPostReturn(ctx, config.Config.Callback.CallbackUrl, cbReq, resp, config.Config.Callback.CallbackBeforeImportFriends); err != nil {
		return err
	}
	if len(resp.FriendUserIDs) != 0 {
		req.FriendUserIDs = resp.FriendUserIDs
	}
	return nil
}
func CallbackAfterImportFriends(ctx context.Context, req *pbfriend.ImportFriendReq) error {
	if !config.Config.Callback.CallbackAfterImportFriends.Enable {
		return nil
	}
	cbReq := &cbapi.CallbackAfterImportFriendsReq{
		CallbackCommand: cbapi.CallbackAfterImportFriendsCommand,
		OwnerUserID:     req.OwnerUserID,
		FriendUserIDs:   req.FriendUserIDs,
	}
	resp := &cbapi.CallbackAfterImportFriendsResp{}
	if err := http.CallBackPostReturn(ctx, config.Config.Callback.CallbackUrl, cbReq, resp, config.Config.Callback.CallbackAfterImportFriends); err != nil {
		return err
	}
	return nil
}

func CallbackAfterRemoveBlack(ctx context.Context, req *pbfriend.RemoveBlackReq) error {
	if !config.Config.Callback.CallbackAfterRemoveBlack.Enable {
		return nil
	}
	cbReq := &cbapi.CallbackAfterRemoveBlackReq{
		CallbackCommand: cbapi.CallbackAfterRemoveBlackCommand,
		OwnerUserID:     req.OwnerUserID,
		BlackUserID:     req.BlackUserID,
	}
	resp := &cbapi.CallbackAfterRemoveBlackResp{}
	if err := http.CallBackPostReturn(ctx, config.Config.Callback.CallbackUrl, cbReq, resp, config.Config.Callback.CallbackAfterRemoveBlack); err != nil {
		return err
	}
	return nil
}
