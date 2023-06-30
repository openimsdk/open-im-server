package friend

import (
	"context"
	"time"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/convert"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/mcontext"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/tokenverify"
	pbFriend "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/friend"
)

func (s *friendServer) GetPaginationBlacks(ctx context.Context, req *pbFriend.GetPaginationBlacksReq) (resp *pbFriend.GetPaginationBlacksResp, err error) {
	if err := s.userRpcClient.Access(ctx, req.UserID); err != nil {
		return nil, err
	}
	var pageNumber, showNumber int32
	if req.Pagination != nil {
		pageNumber = req.Pagination.PageNumber
		showNumber = req.Pagination.ShowNumber
	}
	blacks, total, err := s.blackDatabase.FindOwnerBlacks(ctx, req.UserID, pageNumber, showNumber)
	if err != nil {
		return nil, err
	}
	resp = &pbFriend.GetPaginationBlacksResp{}
	resp.Blacks, err = convert.BlackDB2Pb(ctx, blacks, s.userRpcClient.GetUsersInfoMap)
	if err != nil {
		return nil, err
	}
	resp.Total = int32(total)
	return resp, nil
}

func (s *friendServer) IsBlack(ctx context.Context, req *pbFriend.IsBlackReq) (*pbFriend.IsBlackResp, error) {
	in1, in2, err := s.blackDatabase.CheckIn(ctx, req.UserID1, req.UserID2)
	if err != nil {
		return nil, err
	}
	resp := &pbFriend.IsBlackResp{}
	resp.InUser1Blacks = in1
	resp.InUser2Blacks = in2
	return resp, nil
}

func (s *friendServer) RemoveBlack(ctx context.Context, req *pbFriend.RemoveBlackReq) (*pbFriend.RemoveBlackResp, error) {
	if err := s.userRpcClient.Access(ctx, req.OwnerUserID); err != nil {
		return nil, err
	}
	if err := s.blackDatabase.Delete(ctx, []*relation.BlackModel{{OwnerUserID: req.OwnerUserID, BlockUserID: req.BlackUserID}}); err != nil {
		return nil, err
	}
	s.notificationSender.BlackDeletedNotification(ctx, req)
	return &pbFriend.RemoveBlackResp{}, nil
}

func (s *friendServer) AddBlack(ctx context.Context, req *pbFriend.AddBlackReq) (*pbFriend.AddBlackResp, error) {
	if err := tokenverify.CheckAccessV3(ctx, req.OwnerUserID); err != nil {
		return nil, err
	}
	_, err := s.userRpcClient.GetUsersInfo(ctx, []string{req.OwnerUserID, req.BlackUserID})
	if err != nil {
		return nil, err
	}
	black := relation.BlackModel{OwnerUserID: req.OwnerUserID, BlockUserID: req.BlackUserID, OperatorUserID: mcontext.GetOpUserID(ctx), CreateTime: time.Now()}
	if err := s.blackDatabase.Create(ctx, []*relation.BlackModel{&black}); err != nil {
		return nil, err
	}
	s.notificationSender.BlackAddedNotification(ctx, req)
	return &pbFriend.AddBlackResp{}, nil
}
