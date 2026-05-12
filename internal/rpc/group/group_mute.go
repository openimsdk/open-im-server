package group

import (
	"context"
	"time"

	"github.com/openimsdk/open-im-server/v3/pkg/common/servererrs"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	pbgroup "github.com/openimsdk/protocol/group"
	"github.com/openimsdk/tools/mcontext"
)

func (s *groupServer) SetGroupMute(ctx context.Context, req *pbgroup.SetGroupMuteReq) (*pbgroup.SetGroupMuteResp, error) {
	opUserID := mcontext.GetOpUserID(ctx)
	if opUserID == "" {
		return nil, servererrs.ErrNoPermission.WrapMsg("op user id is empty")
	}
	if _, err := s.db.TakeGroupMember(ctx, req.GroupID, opUserID); err != nil {
		return nil, err
	}
	if req.Duration == 0 {
		return &pbgroup.SetGroupMuteResp{}, s.groupMuteDB.Delete(ctx, opUserID, req.GroupID)
	}
	var muteEnd int64
	if req.Duration != -1 {
		muteEnd = time.Now().Unix() + req.Duration
	}
	return &pbgroup.SetGroupMuteResp{}, s.groupMuteDB.Upsert(ctx, &model.GroupMute{
		OwnerUserID:  opUserID,
		GroupID:      req.GroupID,
		MuteEndTime:  muteEnd,
		MuteDuration: req.Duration,
		CreateTime:   time.Now(),
	})
}

func (s *groupServer) GetGroupMute(ctx context.Context, req *pbgroup.GetGroupMuteReq) (*pbgroup.GetGroupMuteResp, error) {
	opUserID := mcontext.GetOpUserID(ctx)
	if opUserID == "" {
		return nil, servererrs.ErrNoPermission.WrapMsg("op user id is empty")
	}
	if _, err := s.db.TakeGroupMember(ctx, req.GroupID, opUserID); err != nil {
		return nil, err
	}
	rec, err := s.groupMuteDB.Get(ctx, opUserID, req.GroupID)
	if err != nil {
		return nil, err
	}
	if rec == nil {
		return &pbgroup.GetGroupMuteResp{}, nil
	}
	now := time.Now().Unix()
	if rec.MuteEndTime != 0 && rec.MuteEndTime <= now {
		return &pbgroup.GetGroupMuteResp{}, nil
	}
	duration := rec.MuteDuration
	if duration == 0 && rec.MuteEndTime == 0 {
		duration = -1
	}
	return &pbgroup.GetGroupMuteResp{
		Muted:       true,
		MuteEndTime: rec.MuteEndTime,
		Duration:    duration,
	}, nil
}
