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

	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/convert"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/controller"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/mgo"
	tablerelation "github.com/openimsdk/open-im-server/v3/pkg/common/db/table/relation"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/unrelation"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient/notification"
	"github.com/openimsdk/protocol/constant"
	pbfriend "github.com/openimsdk/protocol/friend"
	"github.com/openimsdk/protocol/sdkws"
	registry "github.com/openimsdk/tools/discoveryregistry"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/tx"
	"github.com/openimsdk/tools/utils"
	"google.golang.org/grpc"
)

type friendServer struct {
	friendDatabase        controller.FriendDatabase
	blackDatabase         controller.BlackDatabase
	userRpcClient         *rpcclient.UserRpcClient
	notificationSender    *notification.FriendNotificationSender
	conversationRpcClient rpcclient.ConversationRpcClient
	RegisterCenter        registry.SvcDiscoveryRegistry
	config                *config.GlobalConfig
}

func Start(ctx context.Context, config *config.GlobalConfig, client registry.SvcDiscoveryRegistry, server *grpc.Server) error {
	// Initialize MongoDB
	mongo, err := unrelation.NewMongoDB(ctx, &config.Mongo)
	if err != nil {
		return err
	}

	// Initialize Redis
	rdb, err := cache.NewRedis(ctx, &config.Redis)
	if err != nil {
		return err
	}

	friendMongoDB, err := mgo.NewFriendMongo(mongo.GetDatabase(config.Mongo.Database))
	if err != nil {
		return err
	}

	friendRequestMongoDB, err := mgo.NewFriendRequestMongo(mongo.GetDatabase(config.Mongo.Database))
	if err != nil {
		return err
	}

	blackMongoDB, err := mgo.NewBlackMongo(mongo.GetDatabase(config.Mongo.Database))
	if err != nil {
		return err
	}

	// Initialize RPC clients
	userRpcClient := rpcclient.NewUserRpcClient(client, config.RpcRegisterName.OpenImUserName, &config.Manager, &config.IMAdmin)
	msgRpcClient := rpcclient.NewMessageRpcClient(client, config.RpcRegisterName.OpenImMsgName)

	// Initialize notification sender
	notificationSender := notification.NewFriendNotificationSender(
		&config.Notification,
		&msgRpcClient,
		notification.WithRpcFunc(userRpcClient.GetUsersInfo),
	)

	// Register Friend server with refactored MongoDB and Redis integrations
	pbfriend.RegisterFriendServer(server, &friendServer{
		friendDatabase: controller.NewFriendDatabase(
			friendMongoDB,
			friendRequestMongoDB,
			cache.NewFriendCacheRedis(rdb, friendMongoDB, cache.GetDefaultOpt()),
			tx.NewMongo(mongo.GetClient()),
		),
		blackDatabase: controller.NewBlackDatabase(
			blackMongoDB,
			cache.NewBlackCacheRedis(rdb, blackMongoDB, cache.GetDefaultOpt()),
		),
		userRpcClient:         &userRpcClient,
		notificationSender:    notificationSender,
		RegisterCenter:        client,
		conversationRpcClient: rpcclient.NewConversationRpcClient(client, config.RpcRegisterName.OpenImConversationName),
		config:                config,
	})

	return nil
}

// ok.
func (s *friendServer) ApplyToAddFriend(ctx context.Context, req *pbfriend.ApplyToAddFriendReq) (resp *pbfriend.ApplyToAddFriendResp, err error) {
	resp = &pbfriend.ApplyToAddFriendResp{}
	if err := authverify.CheckAccessV3(ctx, req.FromUserID, &s.config.Manager, &s.config.IMAdmin); err != nil {
		return nil, err
	}
	if req.ToUserID == req.FromUserID {
		return nil, errs.ErrCanNotAddYourself.WrapMsg("req.ToUserID", req.ToUserID)
	}
	if err = CallbackBeforeAddFriend(ctx, &s.config.Callback, req); err != nil && err != errs.ErrCallbackContinue {
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
		return nil, errs.ErrRelationshipAlready.WrapMsg("already friends has f")
	}
	if err = s.friendDatabase.AddFriendRequest(ctx, req.FromUserID, req.ToUserID, req.ReqMsg, req.Ex); err != nil {
		return nil, err
	}
	s.notificationSender.FriendApplicationAddNotification(ctx, req)
	if err = CallbackAfterAddFriend(ctx, &s.config.Callback, req); err != nil && err != errs.ErrCallbackContinue {
		return nil, err
	}
	return resp, nil
}

// ok.
func (s *friendServer) ImportFriends(ctx context.Context, req *pbfriend.ImportFriendReq) (resp *pbfriend.ImportFriendResp, err error) {
	if err := authverify.CheckAdmin(ctx, &s.config.Manager, &s.config.IMAdmin); err != nil {
		return nil, err
	}
	if _, err := s.userRpcClient.GetUsersInfo(ctx, append([]string{req.OwnerUserID}, req.FriendUserIDs...)); err != nil {
		return nil, err
	}
	if utils.Contain(req.OwnerUserID, req.FriendUserIDs...) {
		return nil, errs.ErrCanNotAddYourself.WrapMsg("can not add yourself")
	}
	if utils.Duplicate(req.FriendUserIDs) {
		return nil, errs.ErrArgs.WrapMsg("friend userID repeated")
	}

	if err := CallbackBeforeImportFriends(ctx, &s.config.Callback, req); err != nil {
		return nil, err
	}

	if err := s.friendDatabase.BecomeFriends(ctx, req.OwnerUserID, req.FriendUserIDs, constant.BecomeFriendByImport); err != nil {
		return nil, err
	}
	for _, userID := range req.FriendUserIDs {
		s.notificationSender.FriendApplicationAgreedNotification(ctx, &pbfriend.RespondFriendApplyReq{
			FromUserID:   req.OwnerUserID,
			ToUserID:     userID,
			HandleResult: constant.FriendResponseAgree,
		})
	}
	if err := CallbackAfterImportFriends(ctx, &s.config.Callback, req); err != nil {
		return nil, err
	}
	return &pbfriend.ImportFriendResp{}, nil
}

// ok.
func (s *friendServer) RespondFriendApply(ctx context.Context, req *pbfriend.RespondFriendApplyReq) (resp *pbfriend.RespondFriendApplyResp, err error) {
	resp = &pbfriend.RespondFriendApplyResp{}
	if err := authverify.CheckAccessV3(ctx, req.ToUserID, &s.config.Manager, &s.config.IMAdmin); err != nil {
		return nil, err
	}

	friendRequest := tablerelation.FriendRequestModel{
		FromUserID:   req.FromUserID,
		ToUserID:     req.ToUserID,
		HandleMsg:    req.HandleMsg,
		HandleResult: req.HandleResult,
	}
	if req.HandleResult == constant.FriendResponseAgree {
		if err := CallbackBeforeAddFriendAgree(ctx, &s.config.Callback, req); err != nil && err != errs.ErrCallbackContinue {
			return nil, err
		}
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
	return nil, errs.ErrArgs.WrapMsg("req.HandleResult != -1/1")
}

// ok.
func (s *friendServer) DeleteFriend(ctx context.Context, req *pbfriend.DeleteFriendReq) (resp *pbfriend.DeleteFriendResp, err error) {
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
	if err := CallbackAfterDeleteFriend(ctx, &s.config.Callback, req); err != nil {
		return nil, err
	}
	return resp, nil
}

// ok.
func (s *friendServer) SetFriendRemark(ctx context.Context, req *pbfriend.SetFriendRemarkReq) (resp *pbfriend.SetFriendRemarkResp, err error) {

	if err = CallbackBeforeSetFriendRemark(ctx, &s.config.Callback, req); err != nil && err != errs.ErrCallbackContinue {
		return nil, err
	}
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
	if err := CallbackAfterSetFriendRemark(ctx, &s.config.Callback, req); err != nil && err != errs.ErrCallbackContinue {
		return nil, err
	}
	s.notificationSender.FriendRemarkSetNotification(ctx, req.OwnerUserID, req.FriendUserID)
	return resp, nil
}

// ok.
func (s *friendServer) GetDesignatedFriends(ctx context.Context, req *pbfriend.GetDesignatedFriendsReq) (resp *pbfriend.GetDesignatedFriendsResp, err error) {
	resp = &pbfriend.GetDesignatedFriendsResp{}
	if utils.Duplicate(req.FriendUserIDs) {
		return nil, errs.ErrArgs.WrapMsg("friend userID repeated")
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

// Get the list of friend requests sent out proactively.
func (s *friendServer) GetDesignatedFriendsApply(ctx context.Context,
	req *pbfriend.GetDesignatedFriendsApplyReq) (resp *pbfriend.GetDesignatedFriendsApplyResp, err error) {
	friendRequests, err := s.friendDatabase.FindBothFriendRequests(ctx, req.FromUserID, req.ToUserID)
	if err != nil {
		return nil, err
	}
	resp = &pbfriend.GetDesignatedFriendsApplyResp{}
	resp.FriendRequests, err = convert.FriendRequestDB2Pb(ctx, friendRequests, s.userRpcClient.GetUsersInfoMap)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// Get received friend requests (i.e., those initiated by others).
func (s *friendServer) GetPaginationFriendsApplyTo(ctx context.Context, req *pbfriend.GetPaginationFriendsApplyToReq) (resp *pbfriend.GetPaginationFriendsApplyToResp, err error) {
	if err := s.userRpcClient.Access(ctx, req.UserID); err != nil {
		return nil, err
	}
	total, friendRequests, err := s.friendDatabase.PageFriendRequestToMe(ctx, req.UserID, req.Pagination)
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

func (s *friendServer) GetPaginationFriendsApplyFrom(ctx context.Context, req *pbfriend.GetPaginationFriendsApplyFromReq) (resp *pbfriend.GetPaginationFriendsApplyFromResp, err error) {
	resp = &pbfriend.GetPaginationFriendsApplyFromResp{}
	if err := s.userRpcClient.Access(ctx, req.UserID); err != nil {
		return nil, err
	}
	total, friendRequests, err := s.friendDatabase.PageFriendRequestFromMe(ctx, req.UserID, req.Pagination)
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

// ok.
func (s *friendServer) IsFriend(ctx context.Context, req *pbfriend.IsFriendReq) (resp *pbfriend.IsFriendResp, err error) {
	resp = &pbfriend.IsFriendResp{}
	resp.InUser1Friends, resp.InUser2Friends, err = s.friendDatabase.CheckIn(ctx, req.UserID1, req.UserID2)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *friendServer) GetPaginationFriends(ctx context.Context, req *pbfriend.GetPaginationFriendsReq) (resp *pbfriend.GetPaginationFriendsResp, err error) {
	if err := s.userRpcClient.Access(ctx, req.UserID); err != nil {
		return nil, err
	}
	total, friends, err := s.friendDatabase.PageOwnerFriends(ctx, req.UserID, req.Pagination)
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

func (s *friendServer) GetFriendIDs(ctx context.Context, req *pbfriend.GetFriendIDsReq) (resp *pbfriend.GetFriendIDsResp, err error) {
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

func (s *friendServer) GetSpecifiedFriendsInfo(ctx context.Context, req *pbfriend.GetSpecifiedFriendsInfoReq) (*pbfriend.GetSpecifiedFriendsInfoResp, error) {
	if len(req.UserIDList) == 0 {
		return nil, errs.ErrArgs.WrapMsg("userIDList is empty")
	}
	if utils.Duplicate(req.UserIDList) {
		return nil, errs.ErrArgs.WrapMsg("userIDList repeated")
	}
	userMap, err := s.userRpcClient.GetUsersInfoMap(ctx, req.UserIDList)
	if err != nil {
		return nil, err
	}
	friends, err := s.friendDatabase.FindFriendsWithError(ctx, req.OwnerUserID, req.UserIDList)
	if err != nil {
		return nil, err
	}
	blacks, err := s.blackDatabase.FindBlackInfos(ctx, req.OwnerUserID, req.UserIDList)
	if err != nil {
		return nil, err
	}
	friendMap := utils.SliceToMap(friends, func(e *tablerelation.FriendModel) string {
		return e.FriendUserID
	})
	blackMap := utils.SliceToMap(blacks, func(e *tablerelation.BlackModel) string {
		return e.BlockUserID
	})
	resp := &pbfriend.GetSpecifiedFriendsInfoResp{
		Infos: make([]*pbfriend.GetSpecifiedFriendsInfoInfo, 0, len(req.UserIDList)),
	}
	for _, userID := range req.UserIDList {
		user := userMap[userID]
		if user == nil {
			continue
		}
		var friendInfo *sdkws.FriendInfo
		if friend := friendMap[userID]; friend != nil {

			friendInfo = &sdkws.FriendInfo{
				OwnerUserID:    friend.OwnerUserID,
				Remark:         friend.Remark,
				CreateTime:     friend.CreateTime.UnixMilli(),
				AddSource:      friend.AddSource,
				OperatorUserID: friend.OperatorUserID,
				Ex:             friend.Ex,
				IsPinned:       friend.IsPinned,
			}
		}
		var blackInfo *sdkws.BlackInfo
		if black := blackMap[userID]; black != nil {
			blackInfo = &sdkws.BlackInfo{
				OwnerUserID:    black.OwnerUserID,
				CreateTime:     black.CreateTime.UnixMilli(),
				AddSource:      black.AddSource,
				OperatorUserID: black.OperatorUserID,
				Ex:             black.Ex,
			}
		}
		resp.Infos = append(resp.Infos, &pbfriend.GetSpecifiedFriendsInfoInfo{
			UserInfo:   user,
			FriendInfo: friendInfo,
			BlackInfo:  blackInfo,
		})
	}
	return resp, nil
}
func (s *friendServer) UpdateFriends(
	ctx context.Context,
	req *pbfriend.UpdateFriendsReq,
) (*pbfriend.UpdateFriendsResp, error) {
	if len(req.FriendUserIDs) == 0 {
		return nil, errs.ErrArgs.WrapMsg("friendIDList is empty")
	}
	if utils.Duplicate(req.FriendUserIDs) {
		return nil, errs.ErrArgs.WrapMsg("friendIDList repeated")
	}

	_, err := s.friendDatabase.FindFriendsWithError(ctx, req.OwnerUserID, req.FriendUserIDs)
	if err != nil {
		return nil, err
	}

	val := make(map[string]any)

	if req.IsPinned != nil {
		val["is_pinned"] = req.IsPinned.Value
	}
	if req.Remark != nil {
		val["remark"] = req.Remark.Value
	}
	if req.Ex != nil {
		val["ex"] = req.Ex.Value
	}
	if err = s.friendDatabase.UpdateFriends(ctx, req.OwnerUserID, req.FriendUserIDs, val); err != nil {
		return nil, err
	}

	resp := &pbfriend.UpdateFriendsResp{}

	err = s.notificationSender.FriendsInfoUpdateNotification(ctx, req.OwnerUserID, req.FriendUserIDs)
	if err != nil {
		return nil, errs.WrapMsg(err, "FriendsInfoUpdateNotification Error")
	}
	return resp, nil
}
