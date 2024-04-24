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
	"github.com/openimsdk/open-im-server/v3/pkg/common/webhook"

	cbapi "github.com/openimsdk/open-im-server/v3/pkg/callbackstruct"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	pbfriend "github.com/openimsdk/protocol/friend"
)

func (s *friendServer) webhookAfterDeleteFriend(ctx context.Context, after *config.AfterConfig, req *pbfriend.DeleteFriendReq) {
	cbReq := &cbapi.CallbackAfterDeleteFriendReq{
		CallbackCommand: cbapi.CallbackAfterDeleteFriendCommand,
		OwnerUserID:     req.OwnerUserID,
		FriendUserID:    req.FriendUserID,
	}
	s.webhookClient.AsyncPost(ctx, cbReq.GetCallbackCommand(), cbReq, &cbapi.CallbackAfterDeleteFriendResp{}, after)
}

func (s *friendServer) webhookBeforeAddFriend(ctx context.Context, before *config.BeforeConfig, req *pbfriend.ApplyToAddFriendReq) error {
	return webhook.WithCondition(ctx, before, func(ctx context.Context) error {
		cbReq := &cbapi.CallbackBeforeAddFriendReq{
			CallbackCommand: cbapi.CallbackBeforeAddFriendCommand,
			FromUserID:      req.FromUserID,
			ToUserID:        req.ToUserID,
			ReqMsg:          req.ReqMsg,
			Ex:              req.Ex,
		}
		resp := &cbapi.CallbackBeforeAddFriendResp{}

		if err := s.webhookClient.SyncPost(ctx, cbReq.GetCallbackCommand(), cbReq, resp, before); err != nil {
			return err
		}
		return nil
	})
}

func (s *friendServer) webhookAfterAddFriend(ctx context.Context, after *config.AfterConfig, req *pbfriend.ApplyToAddFriendReq) {
	cbReq := &cbapi.CallbackAfterAddFriendReq{
		CallbackCommand: cbapi.CallbackAfterAddFriendCommand,
		FromUserID:      req.FromUserID,
		ToUserID:        req.ToUserID,
		ReqMsg:          req.ReqMsg,
	}
	resp := &cbapi.CallbackAfterAddFriendResp{}
	s.webhookClient.AsyncPost(ctx, cbReq.GetCallbackCommand(), cbReq, resp, after)
}

func (s *friendServer) webhookAfterSetFriendRemark(ctx context.Context, after *config.AfterConfig, req *pbfriend.SetFriendRemarkReq) {

	cbReq := &cbapi.CallbackAfterSetFriendRemarkReq{
		CallbackCommand: cbapi.CallbackAfterSetFriendRemarkCommand,
		OwnerUserID:     req.OwnerUserID,
		FriendUserID:    req.FriendUserID,
		Remark:          req.Remark,
	}
	resp := &cbapi.CallbackAfterSetFriendRemarkResp{}
	s.webhookClient.AsyncPost(ctx, cbReq.GetCallbackCommand(), cbReq, resp, after)
}

func (s *friendServer) webhookAfterImportFriends(ctx context.Context, after *config.AfterConfig, req *pbfriend.ImportFriendReq) {
	cbReq := &cbapi.CallbackAfterImportFriendsReq{
		CallbackCommand: cbapi.CallbackAfterImportFriendsCommand,
		OwnerUserID:     req.OwnerUserID,
		FriendUserIDs:   req.FriendUserIDs,
	}
	resp := &cbapi.CallbackAfterImportFriendsResp{}
	s.webhookClient.AsyncPost(ctx, cbReq.GetCallbackCommand(), cbReq, resp, after)
}

func (s *friendServer) webhookAfterRemoveBlack(ctx context.Context, after *config.AfterConfig, req *pbfriend.RemoveBlackReq) {
	cbReq := &cbapi.CallbackAfterRemoveBlackReq{
		CallbackCommand: cbapi.CallbackAfterRemoveBlackCommand,
		OwnerUserID:     req.OwnerUserID,
		BlackUserID:     req.BlackUserID,
	}
	resp := &cbapi.CallbackAfterRemoveBlackResp{}
	s.webhookClient.AsyncPost(ctx, cbReq.GetCallbackCommand(), cbReq, resp, after)
}

func (s *friendServer) webhookBeforeSetFriendRemark(ctx context.Context, before *config.BeforeConfig, req *pbfriend.SetFriendRemarkReq) error {
	return webhook.WithCondition(ctx, before, func(ctx context.Context) error {
		cbReq := &cbapi.CallbackBeforeSetFriendRemarkReq{
			CallbackCommand: cbapi.CallbackBeforeSetFriendRemarkCommand,
			OwnerUserID:     req.OwnerUserID,
			FriendUserID:    req.FriendUserID,
			Remark:          req.Remark,
		}
		resp := &cbapi.CallbackBeforeSetFriendRemarkResp{}
		if err := s.webhookClient.SyncPost(ctx, cbReq.GetCallbackCommand(), cbReq, resp, before); err != nil {
			return err
		}
		if resp.Remark != "" {
			req.Remark = resp.Remark
		}
		return nil
	})
}

func (s *friendServer) webhookBeforeAddBlack(ctx context.Context, before *config.BeforeConfig, req *pbfriend.AddBlackReq) error {
	return webhook.WithCondition(ctx, before, func(ctx context.Context) error {
		cbReq := &cbapi.CallbackBeforeAddBlackReq{
			CallbackCommand: cbapi.CallbackBeforeAddBlackCommand,
			OwnerUserID:     req.OwnerUserID,
			BlackUserID:     req.BlackUserID,
		}
		resp := &cbapi.CallbackBeforeAddBlackResp{}
		return s.webhookClient.SyncPost(ctx, cbReq.GetCallbackCommand(), cbReq, resp, before)
	})
}

func (s *friendServer) webhookBeforeAddFriendAgree(ctx context.Context, before *config.BeforeConfig, req *pbfriend.RespondFriendApplyReq) error {
	return webhook.WithCondition(ctx, before, func(ctx context.Context) error {
		cbReq := &cbapi.CallbackBeforeAddFriendAgreeReq{
			CallbackCommand: cbapi.CallbackBeforeAddFriendAgreeCommand,
			FromUserID:      req.FromUserID,
			ToUserID:        req.ToUserID,
			HandleMsg:       req.HandleMsg,
			HandleResult:    req.HandleResult,
		}
		resp := &cbapi.CallbackBeforeAddFriendAgreeResp{}
		return s.webhookClient.SyncPost(ctx, cbReq.GetCallbackCommand(), cbReq, resp, before)
	})
}

func (s *friendServer) webhookBeforeImportFriends(ctx context.Context, before *config.BeforeConfig, req *pbfriend.ImportFriendReq) error {
	return webhook.WithCondition(ctx, before, func(ctx context.Context) error {
		cbReq := &cbapi.CallbackBeforeImportFriendsReq{
			CallbackCommand: cbapi.CallbackBeforeImportFriendsCommand,
			OwnerUserID:     req.OwnerUserID,
			FriendUserIDs:   req.FriendUserIDs,
		}
		resp := &cbapi.CallbackBeforeImportFriendsResp{}
		if err := s.webhookClient.SyncPost(ctx, cbReq.GetCallbackCommand(), cbReq, resp, before); err != nil {
			return err
		}
		if len(resp.FriendUserIDs) > 0 {
			req.FriendUserIDs = resp.FriendUserIDs
		}
		return nil
	})
}
