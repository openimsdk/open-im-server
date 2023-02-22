package friend

import (
	"Open_IM/internal/common/convert"
	"Open_IM/pkg/common/db/table/relation"
	"Open_IM/pkg/common/tokenverify"
	"Open_IM/pkg/common/tracelog"
	pbFriend "Open_IM/pkg/proto/friend"
	"context"
)

func (s *friendServer) GetPaginationBlacks(ctx context.Context, req *pbFriend.GetPaginationBlacksReq) (resp *pbFriend.GetPaginationBlacksResp, err error) {
	resp = &pbFriend.GetPaginationBlacksResp{}
	if err := s.userCheck.Access(ctx, req.UserID); err != nil {
		return nil, err
	}
	blacks, total, err := s.BlackDatabase.FindOwnerBlacks(ctx, req.UserID, req.Pagination.PageNumber, req.Pagination.ShowNumber)
	if err != nil {
		return nil, err
	}
	resp.Blacks, err = (*convert.NewDBBlack(nil, s.RegisterCenter)).DB2PB(ctx, blacks)
	if err != nil {
		return nil, err
	}
	resp.Total = int32(total)
	return resp, nil
}

func (s *friendServer) IsBlack(ctx context.Context, req *pbFriend.IsBlackReq) (*pbFriend.IsBlackResp, error) {
	resp := &pbFriend.IsBlackResp{}
	in1, in2, err := s.BlackDatabase.CheckIn(ctx, req.UserID1, req.UserID2)
	if err != nil {
		return nil, err
	}
	resp.InUser1Blacks = in1
	resp.InUser2Blacks = in2
	return resp, nil
}

func (s *friendServer) RemoveBlack(ctx context.Context, req *pbFriend.RemoveBlackReq) (*pbFriend.RemoveBlackResp, error) {
	resp := &pbFriend.RemoveBlackResp{}
	if err := s.userCheck.Access(ctx, req.OwnerUserID); err != nil {
		return nil, err
	}
	if err := s.BlackDatabase.Delete(ctx, []*relation.BlackModel{{OwnerUserID: req.OwnerUserID, BlockUserID: req.BlackUserID}}); err != nil {
		return nil, err
	}
	s.notification.BlackDeletedNotification(ctx, req)
	return resp, nil
}

func (s *friendServer) AddBlack(ctx context.Context, req *pbFriend.AddBlackReq) (*pbFriend.AddBlackResp, error) {
	resp := &pbFriend.AddBlackResp{}
	if err := tokenverify.CheckAccessV3(ctx, req.OwnerUserID); err != nil {
		return nil, err
	}
	black := relation.BlackModel{OwnerUserID: req.OwnerUserID, BlockUserID: req.BlackUserID, OperatorUserID: tracelog.GetOpUserID(ctx)}
	if err := s.BlackDatabase.Create(ctx, []*relation.BlackModel{&black}); err != nil {
		return nil, err
	}
	s.notification.BlackAddedNotification(ctx, req)
	return resp, nil
}
