package friend

import (
	"Open_IM/internal/common/check"
	"Open_IM/internal/common/convert"
	chat "Open_IM/internal/rpc/msg"
	"Open_IM/pkg/common/db/relation"
	"Open_IM/pkg/common/token_verify"
	"Open_IM/pkg/common/tracelog"
	pbFriend "Open_IM/pkg/proto/friend"
	"context"
)

func (s *friendServer) GetBlacks(ctx context.Context, req *pbFriend.GetBlacksReq) (*pbFriend.GetBlacksResp, error) {
	resp := &pbFriend.GetBlacksResp{}
	if err := check.Access(ctx, req.UserID); err != nil {
		return nil, err
	}
	blacks, total, err := s.BlackInterface.FindOwnerBlacks(ctx, req.UserID, req.Pagination.PageNumber, req.Pagination.ShowNumber)
	if err != nil {
		return nil, err
	}
	blackIDList := make([]string, 0, len(blacks))
	for _, black := range blacks {
		b, err := convert.NewDBBlack(black).Convert()
		if err != nil {
			return nil, err
		}
		resp.Blacks = append(resp.Blacks, b)
		blackIDList = append(blackIDList, black.BlockUserID)
	}
	resp.Total = int32(total)
	return resp, nil
}

func (s *friendServer) IsBlack(ctx context.Context, req *pbFriend.IsBlackReq) (*pbFriend.IsBlackResp, error) {
	resp := &pbFriend.IsBlackResp{}
	in1, in2, err := s.BlackInterface.CheckIn(ctx, req.UserID1, req.UserID2)
	if err != nil {
		return nil, err
	}
	resp.InUser1Blacks = in1
	resp.InUser2Blacks = in2
	return resp, nil
}

func (s *friendServer) RemoveBlack(ctx context.Context, req *pbFriend.RemoveBlackReq) (*pbFriend.RemoveBlackResp, error) {
	resp := &pbFriend.RemoveBlackResp{}
	if err := check.Access(ctx, req.OwnerUserID); err != nil {
		return nil, err
	}
	if err := s.BlackInterface.Delete(ctx, []*relation.Black{{OwnerUserID: req.OwnerUserID, BlockUserID: req.BlackUserID}}); err != nil {
		return nil, err
	}
	chat.BlackDeletedNotification(req)
	return resp, nil
}

func (s *friendServer) AddBlack(ctx context.Context, req *pbFriend.AddBlackReq) (*pbFriend.AddBlackResp, error) {
	resp := &pbFriend.AddBlackResp{}
	if err := token_verify.CheckAccessV3(ctx, req.OwnerUserID); err != nil {
		return nil, err
	}
	black := relation.Black{OwnerUserID: req.OwnerUserID, BlockUserID: req.BlackUserID, OperatorUserID: tracelog.GetOpUserID(ctx)}
	if err := s.BlackInterface.Create(ctx, []*relation.Black{&black}); err != nil {
		return nil, err
	}
	chat.BlackAddedNotification(tracelog.GetOperationID(ctx), req)
	return resp, nil
}
