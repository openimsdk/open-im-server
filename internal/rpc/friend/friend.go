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

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/convert"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/rpcclient"

	"google.golang.org/grpc"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/cache"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/controller"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/relation"
	tablerelation "github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/tx"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/tokenverify"
	registry "github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	pbfriend "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/friend"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/rpcclient/notification"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
)

type friendServer struct {
	friendDatabase     controller.FriendDatabase
	blackDatabase      controller.BlackDatabase
	userRpcClient      *rpcclient.UserRpcClient
	notificationSender *notification.FriendNotificationSender
	RegisterCenter     registry.SvcDiscoveryRegistry
}

func Start(client registry.SvcDiscoveryRegistry, server *grpc.Server) error {
	db, err := relation.NewGormDB()
	if err != nil {
		return err
	}
	if err := db.AutoMigrate(&tablerelation.FriendModel{}, &tablerelation.FriendRequestModel{}, &tablerelation.BlackModel{}); err != nil {
		return err
	}
	rdb, err := cache.NewRedis()
	if err != nil {
		return err
	}
	blackDB := relation.NewBlackGorm(db)
	friendDB := relation.NewFriendGorm(db)
	userRpcClient := rpcclient.NewUserRpcClient(client)
	msgRpcClient := rpcclient.NewMessageRpcClient(client)
	notificationSender := notification.NewFriendNotificationSender(
		&msgRpcClient,
		notification.WithRpcFunc(userRpcClient.GetUsersInfo),
	)
	pbfriend.RegisterFriendServer(server, &friendServer{
		friendDatabase: controller.NewFriendDatabase(
			friendDB,
			relation.NewFriendRequestGorm(db),
			cache.NewFriendCacheRedis(rdb, friendDB, cache.GetDefaultOpt()),
			tx.NewGorm(db),
		),
		blackDatabase: controller.NewBlackDatabase(
			blackDB,
			cache.NewBlackCacheRedis(rdb, blackDB, cache.GetDefaultOpt()),
		),
		userRpcClient:      &userRpcClient,
		notificationSender: notificationSender,
		RegisterCenter:     client,
	})
	return nil
}

// ok
func (s *friendServer) ApplyToAddFriend(
	ctx context.Context,
	req *pbfriend.ApplyToAddFriendReq,
) (resp *pbfriend.ApplyToAddFriendResp, err error) {
	defer log.ZInfo(ctx, utils.GetFuncName()+" Return")
	resp = &pbfriend.ApplyToAddFriendResp{}
	if err := tokenverify.CheckAccessV3(ctx, req.FromUserID); err != nil {
		return nil, err
	}
	if req.ToUserID == req.FromUserID {
		return nil, errs.ErrCanNotAddYourself.Wrap()
	}
	if err := CallbackBeforeAddFriend(ctx, req); err != nil && err != errs.ErrCallbackContinue {
		return nil, err
	}
	if _, err := s.userRpcClient.GetUsersInfoMap(ctx, []string{req.ToUserID, req.FromUserID}); err != nil {
		return nil, err
	}
	in1, in2, err := s.friendDatabase.CheckIn(ctx, req.FromUserID, req.ToUserID)
	if err != nil {
		return nil, err
	}
	if in1 && in2 {
		return nil, errs.ErrRelationshipAlready.Wrap()
	}
	if err = s.friendDatabase.AddFriendRequest(ctx, req.FromUserID, req.ToUserID, req.ReqMsg, req.Ex); err != nil {
		return nil, err
	}
	s.notificationSender.FriendApplicationAddNotification(ctx, req)
	return resp, nil
}

// ok
func (s *friendServer) ImportFriends(
	ctx context.Context,
	req *pbfriend.ImportFriendReq,
) (resp *pbfriend.ImportFriendResp, err error) {
	defer log.ZInfo(ctx, utils.GetFuncName()+" Return")
	if err := tokenverify.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	if _, err := s.userRpcClient.GetUsersInfo(ctx, append([]string{req.OwnerUserID}, req.FriendUserIDs...)); err != nil {
		return nil, err
	}

	if utils.Contain(req.OwnerUserID, req.FriendUserIDs...) {
		return nil, errs.ErrCanNotAddYourself.Wrap()
	}
	if utils.Duplicate(req.FriendUserIDs) {
		return nil, errs.ErrArgs.Wrap("friend userID repeated")
	}

	if err := s.friendDatabase.BecomeFriends(ctx, req.OwnerUserID, req.FriendUserIDs, constant.BecomeFriendByImport); err != nil {
		return nil, err
	}
	return &pbfriend.ImportFriendResp{}, nil
}

// ok
func (s *friendServer) RespondFriendApply(
	ctx context.Context,
	req *pbfriend.RespondFriendApplyReq,
) (resp *pbfriend.RespondFriendApplyResp, err error) {
	defer log.ZInfo(ctx, utils.GetFuncName()+" Return")
	resp = &pbfriend.RespondFriendApplyResp{}
	if err := tokenverify.CheckAccessV3(ctx, req.ToUserID); err != nil {
		return nil, err
	}

	friendRequest := tablerelation.FriendRequestModel{
		FromUserID:   req.FromUserID,
		ToUserID:     req.ToUserID,
		HandleMsg:    req.HandleMsg,
		HandleResult: req.HandleResult,
	}
	if req.HandleResult == constant.FriendResponseAgree {
		err := s.friendDatabase.AgreeFriendRequest(ctx, &friendRequest)
		if err != nil {
			return nil, err
		}
		s.notificationSender.FriendApplicationAgreedNotification(ctx, req)
		return resp, nil
	}
	if req.HandleResult == constant.FriendResponseRefuse {
		err := s.friendDatabase.RefuseFriendRequest(ctx, &friendRequest)
		if err != nil {
			return nil, err
		}
		s.notificationSender.FriendApplicationRefusedNotification(ctx, req)
		return resp, nil
	}
	return nil, errs.ErrArgs.Wrap("req.HandleResult != -1/1")
}

// ok
func (s *friendServer) DeleteFriend(
	ctx context.Context,
	req *pbfriend.DeleteFriendReq,
) (resp *pbfriend.DeleteFriendResp, err error) {
	defer log.ZInfo(ctx, utils.GetFuncName()+" Return")
	resp = &pbfriend.DeleteFriendResp{}
	if err := s.userRpcClient.Access(ctx, req.OwnerUserID); err != nil {
		return nil, err
	}
	_, err = s.friendDatabase.FindFriendsWithError(ctx, req.OwnerUserID, []string{req.FriendUserID})
	if err != nil {
		return nil, err
	}
	if err := s.friendDatabase.Delete(ctx, req.OwnerUserID, []string{req.FriendUserID}); err != nil {
		return nil, err
	}
	s.notificationSender.FriendDeletedNotification(ctx, req)
	return resp, nil
}

// ok
func (s *friendServer) SetFriendRemark(
	ctx context.Context,
	req *pbfriend.SetFriendRemarkReq,
) (resp *pbfriend.SetFriendRemarkResp, err error) {
	defer log.ZInfo(ctx, utils.GetFuncName()+" Return")
	resp = &pbfriend.SetFriendRemarkResp{}
	if err := s.userRpcClient.Access(ctx, req.OwnerUserID); err != nil {
		return nil, err
	}
	_, err = s.friendDatabase.FindFriendsWithError(ctx, req.OwnerUserID, []string{req.FriendUserID})
	if err != nil {
		return nil, err
	}
	if err := s.friendDatabase.UpdateRemark(ctx, req.OwnerUserID, req.FriendUserID, req.Remark); err != nil {
		return nil, err
	}
	s.notificationSender.FriendRemarkSetNotification(ctx, req.OwnerUserID, req.FriendUserID)
	return resp, nil
}

// ok
func (s *friendServer) GetDesignatedFriends(
	ctx context.Context,
	req *pbfriend.GetDesignatedFriendsReq,
) (resp *pbfriend.GetDesignatedFriendsResp, err error) {
	defer log.ZInfo(ctx, utils.GetFuncName()+" Return")
	resp = &pbfriend.GetDesignatedFriendsResp{}
	if utils.Duplicate(req.FriendUserIDs) {
		return nil, errs.ErrArgs.Wrap("friend userID repeated")
	}
	friends, err := s.friendDatabase.FindFriendsWithError(ctx, req.OwnerUserID, req.FriendUserIDs)
	if err != nil {
		return nil, err
	}
	if resp.FriendsInfo, err = convert.FriendsDB2Pb(ctx, friends, s.userRpcClient.GetUsersInfoMap); err != nil {
		return nil, err
	}
	return resp, nil
}

// ok 获取接收到的好友申请（即别人主动申请的）
func (s *friendServer) GetPaginationFriendsApplyTo(
	ctx context.Context,
	req *pbfriend.GetPaginationFriendsApplyToReq,
) (resp *pbfriend.GetPaginationFriendsApplyToResp, err error) {
	defer log.ZInfo(ctx, utils.GetFuncName()+" Return")
	if err := s.userRpcClient.Access(ctx, req.UserID); err != nil {
		return nil, err
	}
	pageNumber, showNumber := utils.GetPage(req.Pagination)
	friendRequests, total, err := s.friendDatabase.PageFriendRequestToMe(ctx, req.UserID, pageNumber, showNumber)
	if err != nil {
		return nil, err
	}
	resp = &pbfriend.GetPaginationFriendsApplyToResp{}
	resp.FriendRequests, err = convert.FriendRequestDB2Pb(ctx, friendRequests, s.userRpcClient.GetUsersInfoMap)
	if err != nil {
		return nil, err
	}
	resp.Total = int32(total)
	return resp, nil
}

// ok 获取主动发出去的好友申请列表
func (s *friendServer) GetPaginationFriendsApplyFrom(
	ctx context.Context,
	req *pbfriend.GetPaginationFriendsApplyFromReq,
) (resp *pbfriend.GetPaginationFriendsApplyFromResp, err error) {
	defer log.ZInfo(ctx, utils.GetFuncName()+" Return")
	resp = &pbfriend.GetPaginationFriendsApplyFromResp{}
	if err := s.userRpcClient.Access(ctx, req.UserID); err != nil {
		return nil, err
	}
	pageNumber, showNumber := utils.GetPage(req.Pagination)
	friendRequests, total, err := s.friendDatabase.PageFriendRequestFromMe(ctx, req.UserID, pageNumber, showNumber)
	if err != nil {
		return nil, err
	}
	resp.FriendRequests, err = convert.FriendRequestDB2Pb(ctx, friendRequests, s.userRpcClient.GetUsersInfoMap)
	if err != nil {
		return nil, err
	}
	resp.Total = int32(total)
	return resp, nil
}

// ok
func (s *friendServer) IsFriend(
	ctx context.Context,
	req *pbfriend.IsFriendReq,
) (resp *pbfriend.IsFriendResp, err error) {
	defer log.ZInfo(ctx, utils.GetFuncName()+" Return")
	resp = &pbfriend.IsFriendResp{}
	resp.InUser1Friends, resp.InUser2Friends, err = s.friendDatabase.CheckIn(ctx, req.UserID1, req.UserID2)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// ok
func (s *friendServer) GetPaginationFriends(
	ctx context.Context,
	req *pbfriend.GetPaginationFriendsReq,
) (resp *pbfriend.GetPaginationFriendsResp, err error) {
	defer log.ZInfo(ctx, utils.GetFuncName()+" Return")
	if err := s.userRpcClient.Access(ctx, req.UserID); err != nil {
		return nil, err
	}
	pageNumber, showNumber := utils.GetPage(req.Pagination)
	friends, total, err := s.friendDatabase.PageOwnerFriends(ctx, req.UserID, pageNumber, showNumber)
	if err != nil {
		return nil, err
	}
	resp = &pbfriend.GetPaginationFriendsResp{}
	resp.FriendsInfo, err = convert.FriendsDB2Pb(ctx, friends, s.userRpcClient.GetUsersInfoMap)
	if err != nil {
		return nil, err
	}
	resp.Total = int32(total)
	return resp, nil
}

func (s *friendServer) GetFriendIDs(
	ctx context.Context,
	req *pbfriend.GetFriendIDsReq,
) (resp *pbfriend.GetFriendIDsResp, err error) {
	defer log.ZInfo(ctx, utils.GetFuncName()+" Return")
	if err := s.userRpcClient.Access(ctx, req.UserID); err != nil {
		return nil, err
	}
	resp = &pbfriend.GetFriendIDsResp{}
	resp.FriendIDs, err = s.friendDatabase.FindFriendUserIDs(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
