package relation

import (
	"context"

	"github.com/openimsdk/open-im-server/v3/pkg/common/webhook"

	cbapi "github.com/openimsdk/open-im-server/v3/pkg/callbackstruct"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/protocol/relation"
)

func (s *friendServer) webhookAfterDeleteFriend(ctx context.Context, after *config.AfterConfig, req *relation.DeleteFriendReq) {
	cbReq := &cbapi.CallbackAfterDeleteFriendReq{
		CallbackCommand: cbapi.CallbackAfterDeleteFriendCommand,
		OwnerUserID:     req.OwnerUserID,
		FriendUserID:    req.FriendUserID,
	}
	s.webhookClient.AsyncPost(ctx, cbReq.GetCallbackCommand(), cbReq, &cbapi.CallbackAfterDeleteFriendResp{}, after)
}

func (s *friendServer) webhookBeforeAddFriend(ctx context.Context, before *config.BeforeConfig, req *relation.ApplyToAddFriendReq) error {
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

func (s *friendServer) webhookAfterAddFriend(ctx context.Context, after *config.AfterConfig, req *relation.ApplyToAddFriendReq) {
	cbReq := &cbapi.CallbackAfterAddFriendReq{
		CallbackCommand: cbapi.CallbackAfterAddFriendCommand,
		FromUserID:      req.FromUserID,
		ToUserID:        req.ToUserID,
		ReqMsg:          req.ReqMsg,
	}
	resp := &cbapi.CallbackAfterAddFriendResp{}
	s.webhookClient.AsyncPost(ctx, cbReq.GetCallbackCommand(), cbReq, resp, after)
}

func (s *friendServer) webhookAfterSetFriendRemark(ctx context.Context, after *config.AfterConfig, req *relation.SetFriendRemarkReq) {
	cbReq := &cbapi.CallbackAfterSetFriendRemarkReq{
		CallbackCommand: cbapi.CallbackAfterSetFriendRemarkCommand,
		OwnerUserID:     req.OwnerUserID,
		FriendUserID:    req.FriendUserID,
		Remark:          req.Remark,
	}
	resp := &cbapi.CallbackAfterSetFriendRemarkResp{}
	s.webhookClient.AsyncPost(ctx, cbReq.GetCallbackCommand(), cbReq, resp, after)
}

func (s *friendServer) webhookAfterImportFriends(ctx context.Context, after *config.AfterConfig, req *relation.ImportFriendReq) {
	cbReq := &cbapi.CallbackAfterImportFriendsReq{
		CallbackCommand: cbapi.CallbackAfterImportFriendsCommand,
		OwnerUserID:     req.OwnerUserID,
		FriendUserIDs:   req.FriendUserIDs,
	}
	resp := &cbapi.CallbackAfterImportFriendsResp{}
	s.webhookClient.AsyncPost(ctx, cbReq.GetCallbackCommand(), cbReq, resp, after)
}

func (s *friendServer) webhookAfterRemoveBlack(ctx context.Context, after *config.AfterConfig, req *relation.RemoveBlackReq) {
	cbReq := &cbapi.CallbackAfterRemoveBlackReq{
		CallbackCommand: cbapi.CallbackAfterRemoveBlackCommand,
		OwnerUserID:     req.OwnerUserID,
		BlackUserID:     req.BlackUserID,
	}
	resp := &cbapi.CallbackAfterRemoveBlackResp{}
	s.webhookClient.AsyncPost(ctx, cbReq.GetCallbackCommand(), cbReq, resp, after)
}

func (s *friendServer) webhookBeforeSetFriendRemark(ctx context.Context, before *config.BeforeConfig, req *relation.SetFriendRemarkReq) error {
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

func (s *friendServer) webhookBeforeAddBlack(ctx context.Context, before *config.BeforeConfig, req *relation.AddBlackReq) error {
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

func (s *friendServer) webhookBeforeAddFriendAgree(ctx context.Context, before *config.BeforeConfig, req *relation.RespondFriendApplyReq) error {
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

func (s *friendServer) webhookAfterAddFriendAgree(ctx context.Context, after *config.AfterConfig, req *relation.RespondFriendApplyReq) {
	cbReq := &cbapi.CallbackAfterAddFriendAgreeReq{
		CallbackCommand: cbapi.CallbackAfterAddFriendAgreeCommand,
		FromUserID:      req.FromUserID,
		ToUserID:        req.ToUserID,
		HandleMsg:       req.HandleMsg,
		HandleResult:    req.HandleResult,
	}
	resp := &cbapi.CallbackAfterAddFriendAgreeResp{}
	s.webhookClient.AsyncPost(ctx, cbReq.GetCallbackCommand(), cbReq, resp, after)
}

func (s *friendServer) webhookBeforeImportFriends(ctx context.Context, before *config.BeforeConfig, req *relation.ImportFriendReq) error {
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
