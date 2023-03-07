package friend

import (
	"OpenIM/internal/common/check"
	"OpenIM/internal/common/convert"
	"OpenIM/internal/common/notification"
	"OpenIM/pkg/common/constant"
	"OpenIM/pkg/common/db/controller"
	"OpenIM/pkg/common/db/relation"
	tablerelation "OpenIM/pkg/common/db/table/relation"
	"OpenIM/pkg/common/db/tx"
	"OpenIM/pkg/common/tokenverify"
	"OpenIM/pkg/common/tracelog"
	registry "OpenIM/pkg/discoveryregistry"
	"OpenIM/pkg/errs"
	pbfriend "OpenIM/pkg/proto/friend"
	"OpenIM/pkg/utils"
	"context"
	"google.golang.org/grpc"
)

type friendServer struct {
	controller.FriendDatabase
	controller.BlackDatabase
	notification   *notification.Check
	userCheck      *check.UserCheck
	RegisterCenter registry.SvcDiscoveryRegistry
}

func Start(client registry.SvcDiscoveryRegistry, server *grpc.Server) error {
	db, err := relation.NewGormDB()
	if err != nil {
		return err
	}
	if err := db.AutoMigrate(&tablerelation.FriendModel{}, &tablerelation.FriendRequestModel{}, &tablerelation.BlackModel{}); err != nil {
		return err
	}
	pbfriend.RegisterFriendServer(server, &friendServer{
		FriendDatabase: controller.NewFriendDatabase(relation.NewFriendGorm(db), relation.NewFriendRequestGorm(db), tx.NewGorm(db)),
		BlackDatabase:  controller.NewBlackDatabase(relation.NewBlackGorm(db)),
		notification:   notification.NewCheck(client),
		userCheck:      check.NewUserCheck(client),
		RegisterCenter: client,
	})
	return nil
}

// ok
func (s *friendServer) ApplyToAddFriend(ctx context.Context, req *pbfriend.ApplyToAddFriendReq) (resp *pbfriend.ApplyToAddFriendResp, err error) {
	resp = &pbfriend.ApplyToAddFriendResp{}
	if err := tokenverify.CheckAccessV3(ctx, req.FromUserID); err != nil {
		return nil, err
	}
	if err := CallbackBeforeAddFriend(ctx, req); err != nil && err != errs.ErrCallbackContinue {
		return nil, err
	}
	if req.ToUserID == req.FromUserID {
		return nil, errs.ErrCanNotAddYourself.Wrap()
	}
	if _, err := s.userCheck.GetUsersInfoMap(ctx, []string{req.ToUserID, req.FromUserID}, true); err != nil {
		return nil, err
	}
	in1, in2, err := s.FriendDatabase.CheckIn(ctx, req.FromUserID, req.ToUserID)
	if err != nil {
		return nil, err
	}
	if in1 && in2 {
		return nil, errs.ErrRelationshipAlready.Wrap()
	}
	if err = s.FriendDatabase.AddFriendRequest(ctx, req.FromUserID, req.ToUserID, req.ReqMsg, req.Ex); err != nil {
		return nil, err
	}
	s.notification.FriendApplicationAddNotification(ctx, req)
	return resp, nil
}

// ok
func (s *friendServer) ImportFriends(ctx context.Context, req *pbfriend.ImportFriendReq) (resp *pbfriend.ImportFriendResp, err error) {
	resp = &pbfriend.ImportFriendResp{}
	if err := tokenverify.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	if _, err := s.userCheck.GetUsersInfos(ctx, append([]string{req.OwnerUserID}, req.FriendUserIDs...), true); err != nil {
		return nil, err
	}

	if utils.Contain(req.OwnerUserID, req.FriendUserIDs...) {
		return nil, errs.ErrCanNotAddYourself.Wrap()
	}
	if utils.Duplicate(req.FriendUserIDs) {
		return nil, errs.ErrArgs.Wrap("friend userID repeated")
	}

	if err := s.FriendDatabase.BecomeFriends(ctx, req.OwnerUserID, req.FriendUserIDs, constant.BecomeFriendByImport, tracelog.GetOpUserID(ctx)); err != nil {
		return nil, err
	}
	return resp, nil
}

// ok
func (s *friendServer) RespondFriendApply(ctx context.Context, req *pbfriend.RespondFriendApplyReq) (resp *pbfriend.RespondFriendApplyResp, err error) {
	resp = &pbfriend.RespondFriendApplyResp{}
	if err := s.userCheck.Access(ctx, req.ToUserID); err != nil {
		return nil, err
	}
	friendRequest := tablerelation.FriendRequestModel{FromUserID: req.FromUserID, ToUserID: req.ToUserID, HandleMsg: req.HandleMsg, HandleResult: req.HandleResult}
	if req.HandleResult == constant.FriendResponseAgree {
		err := s.AgreeFriendRequest(ctx, &friendRequest)
		if err != nil {
			return nil, err
		}
		s.notification.FriendApplicationAgreedNotification(ctx, req)
		return resp, nil
	}
	if req.HandleResult == constant.FriendResponseRefuse {
		err := s.RefuseFriendRequest(ctx, &friendRequest)
		if err != nil {
			return nil, err
		}
		s.notification.FriendApplicationRefusedNotification(ctx, req)
		return resp, nil
	}
	return nil, errs.ErrArgs.Wrap("req.HandleResult != -1/1")
}

// ok
func (s *friendServer) DeleteFriend(ctx context.Context, req *pbfriend.DeleteFriendReq) (resp *pbfriend.DeleteFriendResp, err error) {
	resp = &pbfriend.DeleteFriendResp{}
	if err := s.userCheck.Access(ctx, req.OwnerUserID); err != nil {
		return nil, err
	}
	_, err = s.FindFriendsWithError(ctx, req.OwnerUserID, []string{req.FriendUserID})
	if err != nil {
		return nil, err
	}
	if err := s.FriendDatabase.Delete(ctx, req.OwnerUserID, []string{req.FriendUserID}); err != nil {
		return nil, err
	}
	s.notification.FriendDeletedNotification(ctx, req)
	return resp, nil
}

// ok
func (s *friendServer) SetFriendRemark(ctx context.Context, req *pbfriend.SetFriendRemarkReq) (resp *pbfriend.SetFriendRemarkResp, err error) {
	resp = &pbfriend.SetFriendRemarkResp{}
	if err := s.userCheck.Access(ctx, req.OwnerUserID); err != nil {
		return nil, err
	}
	_, err = s.FindFriendsWithError(ctx, req.OwnerUserID, []string{req.FriendUserID})
	if err != nil {
		return nil, err
	}
	if err := s.FriendDatabase.UpdateRemark(ctx, req.OwnerUserID, req.FriendUserID, req.Remark); err != nil {
		return nil, err
	}
	s.notification.FriendRemarkSetNotification(ctx, req.OwnerUserID, req.FriendUserID)
	return resp, nil
}

// ok
func (s *friendServer) GetDesignatedFriends(ctx context.Context, req *pbfriend.GetDesignatedFriendsReq) (resp *pbfriend.GetDesignatedFriendsResp, err error) {

	resp = &pbfriend.GetDesignatedFriendsResp{}

	if utils.Duplicate(req.FriendUserIDs) {
		return nil, errs.ErrArgs.Wrap("friend userID repeated")
	}
	friends, err := s.FriendDatabase.FindFriendsWithError(ctx, req.OwnerUserID, req.FriendUserIDs)
	if err != nil {
		return nil, err
	}
	if resp.FriendsInfo, err = (*convert.NewDBFriend(nil, s.RegisterCenter)).DB2PB(ctx, friends); err != nil {
		return nil, err
	}
	return resp, nil

}

// ok 获取接收到的好友申请（即别人主动申请的）
func (s *friendServer) GetPaginationFriendsApplyTo(ctx context.Context, req *pbfriend.GetPaginationFriendsApplyToReq) (resp *pbfriend.GetPaginationFriendsApplyToResp, err error) {
	resp = &pbfriend.GetPaginationFriendsApplyToResp{}
	if err := s.userCheck.Access(ctx, req.UserID); err != nil {
		return nil, err
	}
	friendRequests, total, err := s.FriendDatabase.PageFriendRequestToMe(ctx, req.UserID, req.Pagination.PageNumber, req.Pagination.ShowNumber)
	if err != nil {
		return nil, err
	}
	resp.FriendRequests, err = (*convert.NewDBFriendRequest(nil, s.RegisterCenter)).DB2PB(ctx, friendRequests)
	if err != nil {
		return nil, err
	}
	resp.Total = int32(total)
	return resp, nil
}

// ok 获取主动发出去的好友申请列表
func (s *friendServer) GetPaginationFriendsApplyFrom(ctx context.Context, req *pbfriend.GetPaginationFriendsApplyFromReq) (resp *pbfriend.GetPaginationFriendsApplyFromResp, err error) {
	resp = &pbfriend.GetPaginationFriendsApplyFromResp{}
	if err := s.userCheck.Access(ctx, req.UserID); err != nil {
		return nil, err
	}
	friendRequests, total, err := s.FriendDatabase.PageFriendRequestFromMe(ctx, req.UserID, req.Pagination.PageNumber, req.Pagination.ShowNumber)
	if err != nil {
		return nil, err
	}
	resp.FriendRequests, err = (*convert.NewDBFriendRequest(nil, s.RegisterCenter)).DB2PB(ctx, friendRequests)
	if err != nil {
		return nil, err
	}
	resp.Total = int32(total)
	return resp, nil
}

// ok
func (s *friendServer) IsFriend(ctx context.Context, req *pbfriend.IsFriendReq) (resp *pbfriend.IsFriendResp, err error) {
	resp = &pbfriend.IsFriendResp{}
	resp.InUser1Friends, resp.InUser2Friends, err = s.FriendDatabase.CheckIn(ctx, req.UserID1, req.UserID2)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// ok
func (s *friendServer) GetPaginationFriends(ctx context.Context, req *pbfriend.GetPaginationFriendsReq) (resp *pbfriend.GetPaginationFriendsResp, err error) {
	if err := s.userCheck.Access(ctx, req.UserID); err != nil {
		return nil, err
	}
	friends, total, err := s.FriendDatabase.PageOwnerFriends(ctx, req.UserID, req.Pagination.PageNumber, req.Pagination.ShowNumber)
	if err != nil {
		return nil, err
	}
	resp = &pbfriend.GetPaginationFriendsResp{}
	resp.FriendsInfo, err = (*convert.NewDBFriend(nil, s.RegisterCenter)).DB2PB(ctx, friends)
	if err != nil {
		return nil, err
	}
	resp.Total = int32(total)
	return resp, nil
}

func (s *friendServer) GetFriendIDs(ctx context.Context, req *pbfriend.GetFriendIDsReq) (resp *pbfriend.GetFriendIDsResp, err error) {
	if err := s.userCheck.Access(ctx, req.UserID); err != nil {
		return nil, err
	}
	resp = &pbfriend.GetFriendIDsResp{}
	resp.FriendIDs, err = s.FriendDatabase.FindFriendUserIDs(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
