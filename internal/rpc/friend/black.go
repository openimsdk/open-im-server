package friend

import (
	"Open_IM/internal/common/check"
	"Open_IM/internal/common/convert"
	chat "Open_IM/internal/rpc/msg"
	"Open_IM/pkg/common/db/table/relation"
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
	resp.Blacks, err = (*convert.DBBlack)(nil).DB2PB(blacks)
	if err != nil {
		return nil, err
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
	if err := s.BlackInterface.Delete(ctx, []*relation.BlackModel{{OwnerUserID: req.OwnerUserID, BlockUserID: req.BlackUserID}}); err != nil {
		return nil, err
	}
	chat.BlackDeletedNotification(ctx, req)
	return resp, nil
}

func (s *friendServer) AddBlack(ctx context.Context, req *pbFriend.AddBlackReq) (*pbFriend.AddBlackResp, error) {
	resp := &pbFriend.AddBlackResp{}
	if err := token_verify.CheckAccessV3(ctx, req.OwnerUserID); err != nil {
		return nil, err
	}
	black := relation.BlackModel{OwnerUserID: req.OwnerUserID, BlockUserID: req.BlackUserID, OperatorUserID: tracelog.GetOpUserID(ctx)}
	if err := s.BlackInterface.Create(ctx, []*relation.BlackModel{&black}); err != nil {
		return nil, err
	}
	chat.BlackAddedNotification(ctx, req)
	return resp, nil
}
