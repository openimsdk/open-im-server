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

	cbapi "github.com/openimsdk/open-im-server/v3/pkg/callbackstruct"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/http"
	pbfriend "github.com/openimsdk/protocol/friend"
	"github.com/openimsdk/tools/utils/datautil"
)

func CallbackBeforeAddFriend(ctx context.Context, callback *config.Webhooks, req *pbfriend.ApplyToAddFriendReq) error {
	if !callback.BeforeAddFriend.Enable {
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
	if err := http.CallBackPostReturn(ctx, callback.URL, cbReq, resp, callback.BeforeAddFriend); err != nil {
		return err
	}
	return nil
}

func CallbackBeforeSetFriendRemark(ctx context.Context, callback *config.Webhooks, req *pbfriend.SetFriendRemarkReq) error {
	if !callback.BeforeSetFriendRemark.Enable {
		return nil
	}
	cbReq := &cbapi.CallbackBeforeSetFriendRemarkReq{
		CallbackCommand: cbapi.CallbackBeforeSetFriendRemark,
		OwnerUserID:     req.OwnerUserID,
		FriendUserID:    req.FriendUserID,
		Remark:          req.Remark,
	}
	resp := &cbapi.CallbackBeforeSetFriendRemarkResp{}
	if err := http.CallBackPostReturn(ctx, callback.URL, cbReq, resp, callback.BeforeAddFriend); err != nil {
		return err
	}
	datautil.NotNilReplace(&req.Remark, &resp.Remark)
	return nil
}

func CallbackAfterSetFriendRemark(ctx context.Context, callback *config.Webhooks, req *pbfriend.SetFriendRemarkReq) error {
	if !callback.AfterSetFriendRemark.Enable {
		return nil
	}
	cbReq := &cbapi.CallbackAfterSetFriendRemarkReq{
		CallbackCommand: cbapi.CallbackAfterSetFriendRemark,
		OwnerUserID:     req.OwnerUserID,
		FriendUserID:    req.FriendUserID,
		Remark:          req.Remark,
	}
	resp := &cbapi.CallbackAfterSetFriendRemarkResp{}
	if err := http.CallBackPostReturn(ctx, callback.URL, cbReq, resp, callback.BeforeAddFriend); err != nil {
		return err
	}
	return nil
}

func CallbackBeforeAddBlack(ctx context.Context, callback *config.Webhooks, req *pbfriend.AddBlackReq) error {
	if !callback.BeforeAddBlack.Enable {
		return nil
	}
	cbReq := &cbapi.CallbackBeforeAddBlackReq{
		CallbackCommand: cbapi.CallbackBeforeAddBlackCommand,
		OwnerUserID:     req.OwnerUserID,
		BlackUserID:     req.BlackUserID,
	}
	resp := &cbapi.CallbackBeforeAddBlackResp{}
	if err := http.CallBackPostReturn(ctx, callback.URL, cbReq, resp, callback.BeforeAddBlack); err != nil {
		return err
	}
	return nil
}

func CallbackAfterAddFriend(ctx context.Context, callback *config.Webhooks, req *pbfriend.ApplyToAddFriendReq) error {
	if !callback.AfterAddFriend.Enable {
		return nil
	}
	cbReq := &cbapi.CallbackAfterAddFriendReq{
		CallbackCommand: cbapi.CallbackAfterAddFriendCommand,
		FromUserID:      req.FromUserID,
		ToUserID:        req.ToUserID,
		ReqMsg:          req.ReqMsg,
	}
	resp := &cbapi.CallbackAfterAddFriendResp{}
	if err := http.CallBackPostReturn(ctx, callback.URL, cbReq, resp, callback.AfterAddFriend); err != nil {
		return err
	}

	return nil
}

func CallbackBeforeAddFriendAgree(ctx context.Context, callback *config.Webhooks, req *pbfriend.RespondFriendApplyReq) error {
	if !callback.BeforeAddFriendAgree.Enable {
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
	if err := http.CallBackPostReturn(ctx, callback.URL, cbReq, resp, callback.BeforeAddFriendAgree); err != nil {
		return err
	}
	return nil
}

func CallbackAfterDeleteFriend(ctx context.Context, callback *config.Webhooks, req *pbfriend.DeleteFriendReq) error {
	if !callback.AfterDeleteFriend.Enable {
		return nil
	}
	cbReq := &cbapi.CallbackAfterDeleteFriendReq{
		CallbackCommand: cbapi.CallbackAfterDeleteFriendCommand,
		OwnerUserID:     req.OwnerUserID,
		FriendUserID:    req.FriendUserID,
	}
	resp := &cbapi.CallbackAfterDeleteFriendResp{}
	if err := http.CallBackPostReturn(ctx, callback.URL, cbReq, resp, callback.AfterDeleteFriend); err != nil {
		return err
	}
	return nil
}

func CallbackBeforeImportFriends(ctx context.Context, callback *config.Webhooks, req *pbfriend.ImportFriendReq) error {
	if !callback.BeforeImportFriends.Enable {
		return nil
	}
	cbReq := &cbapi.CallbackBeforeImportFriendsReq{
		CallbackCommand: cbapi.CallbackBeforeImportFriendsCommand,
		OwnerUserID:     req.OwnerUserID,
		FriendUserIDs:   req.FriendUserIDs,
	}
	resp := &cbapi.CallbackBeforeImportFriendsResp{}
	if err := http.CallBackPostReturn(ctx, callback.URL, cbReq, resp, callback.BeforeImportFriends); err != nil {
		return err
	}
	if len(resp.FriendUserIDs) != 0 {
		req.FriendUserIDs = resp.FriendUserIDs
	}
	return nil
}

func CallbackAfterImportFriends(ctx context.Context, callback *config.Webhooks, req *pbfriend.ImportFriendReq) error {
	if !callback.AfterImportFriends.Enable {
		return nil
	}
	cbReq := &cbapi.CallbackAfterImportFriendsReq{
		CallbackCommand: cbapi.CallbackAfterImportFriendsCommand,
		OwnerUserID:     req.OwnerUserID,
		FriendUserIDs:   req.FriendUserIDs,
	}
	resp := &cbapi.CallbackAfterImportFriendsResp{}
	if err := http.CallBackPostReturn(ctx, callback.URL, cbReq, resp, callback.AfterImportFriends); err != nil {
		return err
	}
	return nil
}

func CallbackAfterRemoveBlack(ctx context.Context, callback *config.Webhooks, req *pbfriend.RemoveBlackReq) error {
	if !callback.AfterRemoveBlack.Enable {
		return nil
	}
	cbReq := &cbapi.CallbackAfterRemoveBlackReq{
		CallbackCommand: cbapi.CallbackAfterRemoveBlackCommand,
		OwnerUserID:     req.OwnerUserID,
		BlackUserID:     req.BlackUserID,
	}
	resp := &cbapi.CallbackAfterRemoveBlackResp{}
	if err := http.CallBackPostReturn(ctx, callback.URL, cbReq, resp, callback.AfterRemoveBlack); err != nil {
		return err
	}
	return nil
}
