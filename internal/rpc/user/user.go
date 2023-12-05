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

package user

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/OpenIMSDK/tools/tx"

	"github.com/openimsdk/open-im-server/v3/pkg/common/db/mgo"

	"github.com/OpenIMSDK/protocol/constant"
	"github.com/OpenIMSDK/protocol/sdkws"
	"github.com/OpenIMSDK/tools/errs"
	"github.com/OpenIMSDK/tools/log"

	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/unrelation"

	registry "github.com/OpenIMSDK/tools/discoveryregistry"

	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/convert"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/db/controller"
	tablerelation "github.com/openimsdk/open-im-server/v3/pkg/common/db/table/relation"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcclient/notification"

	pbuser "github.com/OpenIMSDK/protocol/user"
	"github.com/OpenIMSDK/tools/utils"
	"google.golang.org/grpc"
)

type userServer struct {
	controller.UserDatabase
	friendNotificationSender *notification.FriendNotificationSender
	userNotificationSender   *notification.UserNotificationSender
	friendRpcClient          *rpcclient.FriendRpcClient
	groupRpcClient           *rpcclient.GroupRpcClient
	RegisterCenter           registry.SvcDiscoveryRegistry
}

func Start(client registry.SvcDiscoveryRegistry, server *grpc.Server) error {
	rdb, err := cache.NewRedis()
	if err != nil {
		return err
	}
	mongo, err := unrelation.NewMongo()
	if err != nil {
		return err
	}
	users := make([]*tablerelation.UserModel, 0)
	if len(config.Config.Manager.UserID) != len(config.Config.Manager.Nickname) {
		return errors.New("len(config.Config.Manager.AppManagerUid) != len(config.Config.Manager.Nickname)")
	}
	for k, v := range config.Config.Manager.UserID {
		users = append(users, &tablerelation.UserModel{UserID: v, Nickname: config.Config.Manager.Nickname[k], AppMangerLevel: constant.AppAdmin})
	}
	userDB, err := mgo.NewUserMongo(mongo.GetDatabase())
	if err != nil {
		return err
	}
	cache := cache.NewUserCacheRedis(rdb, userDB, cache.GetDefaultOpt())
	userMongoDB := unrelation.NewUserMongoDriver(mongo.GetDatabase())
	database := controller.NewUserDatabase(userDB, cache, tx.NewMongo(mongo.GetClient()), userMongoDB)
	friendRpcClient := rpcclient.NewFriendRpcClient(client)
	groupRpcClient := rpcclient.NewGroupRpcClient(client)
	msgRpcClient := rpcclient.NewMessageRpcClient(client)
	u := &userServer{
		UserDatabase:             database,
		RegisterCenter:           client,
		friendRpcClient:          &friendRpcClient,
		groupRpcClient:           &groupRpcClient,
		friendNotificationSender: notification.NewFriendNotificationSender(&msgRpcClient, notification.WithDBFunc(database.FindWithError)),
		userNotificationSender:   notification.NewUserNotificationSender(&msgRpcClient, notification.WithUserFunc(database.FindWithError)),
	}
	pbuser.RegisterUserServer(server, u)
	return u.UserDatabase.InitOnce(context.Background(), users)
}

func (s *userServer) GetDesignateUsers(ctx context.Context, req *pbuser.GetDesignateUsersReq) (resp *pbuser.GetDesignateUsersResp, err error) {
	resp = &pbuser.GetDesignateUsersResp{}
	users, err := s.FindWithError(ctx, req.UserIDs)
	if err != nil {
		return nil, err
	}
	resp.UsersInfo = convert.UsersDB2Pb(users)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *userServer) UpdateUserInfo(ctx context.Context, req *pbuser.UpdateUserInfoReq) (resp *pbuser.UpdateUserInfoResp, err error) {
	resp = &pbuser.UpdateUserInfoResp{}
	err = authverify.CheckAccessV3(ctx, req.UserInfo.UserID)
	if err != nil {
		return nil, err
	}
	if err := CallbackBeforeUpdateUserInfo(ctx, req); err != nil {
		return nil, err
	}
	data := convert.UserPb2DBMap(req.UserInfo)
	if err := s.UpdateByMap(ctx, req.UserInfo.UserID, data); err != nil {
		return nil, err
	}
	_ = s.friendNotificationSender.UserInfoUpdatedNotification(ctx, req.UserInfo.UserID)
	friends, err := s.friendRpcClient.GetFriendIDs(ctx, req.UserInfo.UserID)
	if err != nil {
		return nil, err
	}
	if req.UserInfo.Nickname != "" || req.UserInfo.FaceURL != "" {
		if err := s.groupRpcClient.NotificationUserInfoUpdate(ctx, req.UserInfo.UserID); err != nil {
			log.ZError(ctx, "NotificationUserInfoUpdate", err)
		}
	}
	for _, friendID := range friends {
		s.friendNotificationSender.FriendInfoUpdatedNotification(ctx, req.UserInfo.UserID, friendID)
	}
	if err := CallbackAfterUpdateUserInfo(ctx, req); err != nil {
		return nil, err
	}
	if err := s.groupRpcClient.NotificationUserInfoUpdate(ctx, req.UserInfo.UserID); err != nil {
		log.ZError(ctx, "NotificationUserInfoUpdate", err, "userID", req.UserInfo.UserID)
	}
	return resp, nil
}

func (s *userServer) SetGlobalRecvMessageOpt(ctx context.Context, req *pbuser.SetGlobalRecvMessageOptReq) (resp *pbuser.SetGlobalRecvMessageOptResp, err error) {
	resp = &pbuser.SetGlobalRecvMessageOptResp{}
	if _, err := s.FindWithError(ctx, []string{req.UserID}); err != nil {
		return nil, err
	}
	m := make(map[string]any, 1)
	m["global_recv_msg_opt"] = req.GlobalRecvMsgOpt
	if err := s.UpdateByMap(ctx, req.UserID, m); err != nil {
		return nil, err
	}
	s.friendNotificationSender.UserInfoUpdatedNotification(ctx, req.UserID)
	return resp, nil
}

func (s *userServer) AccountCheck(ctx context.Context, req *pbuser.AccountCheckReq) (resp *pbuser.AccountCheckResp, err error) {
	resp = &pbuser.AccountCheckResp{}
	if utils.Duplicate(req.CheckUserIDs) {
		return nil, errs.ErrArgs.Wrap("userID repeated")
	}
	err = authverify.CheckAdmin(ctx)
	if err != nil {
		return nil, err
	}
	users, err := s.Find(ctx, req.CheckUserIDs)
	if err != nil {
		return nil, err
	}
	userIDs := make(map[string]any, 0)
	for _, v := range users {
		userIDs[v.UserID] = nil
	}
	for _, v := range req.CheckUserIDs {
		temp := &pbuser.AccountCheckRespSingleUserStatus{UserID: v}
		if _, ok := userIDs[v]; ok {
			temp.AccountStatus = constant.Registered
		} else {
			temp.AccountStatus = constant.UnRegistered
		}
		resp.Results = append(resp.Results, temp)
	}
	return resp, nil
}

func (s *userServer) GetPaginationUsers(ctx context.Context, req *pbuser.GetPaginationUsersReq) (resp *pbuser.GetPaginationUsersResp, err error) {
	total, users, err := s.Page(ctx, req.Pagination)
	if err != nil {
		return nil, err
	}
	return &pbuser.GetPaginationUsersResp{Total: int32(total), Users: convert.UsersDB2Pb(users)}, err
}

func (s *userServer) UserRegister(ctx context.Context, req *pbuser.UserRegisterReq) (resp *pbuser.UserRegisterResp, err error) {
	resp = &pbuser.UserRegisterResp{}
	if len(req.Users) == 0 {
		return nil, errs.ErrArgs.Wrap("users is empty")
	}
	if req.Secret != config.Config.Secret {
		log.ZDebug(ctx, "UserRegister", config.Config.Secret, req.Secret)
		return nil, errs.ErrNoPermission.Wrap("secret invalid")
	}
	if utils.DuplicateAny(req.Users, func(e *sdkws.UserInfo) string { return e.UserID }) {
		return nil, errs.ErrArgs.Wrap("userID repeated")
	}
	userIDs := make([]string, 0)
	for _, user := range req.Users {
		if user.UserID == "" {
			return nil, errs.ErrArgs.Wrap("userID is empty")
		}
		if strings.Contains(user.UserID, ":") {
			return nil, errs.ErrArgs.Wrap("userID contains ':' is invalid userID")
		}
		userIDs = append(userIDs, user.UserID)
	}
	exist, err := s.IsExist(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, errs.ErrRegisteredAlready.Wrap("userID registered already")
	}
	if err := CallbackBeforeUserRegister(ctx, req); err != nil {
		return nil, err
	}
	now := time.Now()
	users := make([]*tablerelation.UserModel, 0, len(req.Users))
	for _, user := range req.Users {
		users = append(users, &tablerelation.UserModel{
			UserID:           user.UserID,
			Nickname:         user.Nickname,
			FaceURL:          user.FaceURL,
			Ex:               user.Ex,
			CreateTime:       now,
			AppMangerLevel:   user.AppMangerLevel,
			GlobalRecvMsgOpt: user.GlobalRecvMsgOpt,
		})
	}
	if err := s.Create(ctx, users); err != nil {
		return nil, err
	}

	if err := CallbackAfterUserRegister(ctx, req); err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *userServer) GetGlobalRecvMessageOpt(ctx context.Context, req *pbuser.GetGlobalRecvMessageOptReq) (resp *pbuser.GetGlobalRecvMessageOptResp, err error) {
	user, err := s.FindWithError(ctx, []string{req.UserID})
	if err != nil {
		return nil, err
	}
	return &pbuser.GetGlobalRecvMessageOptResp{GlobalRecvMsgOpt: user[0].GlobalRecvMsgOpt}, nil
}

// GetAllUserID Get user account by page.
func (s *userServer) GetAllUserID(ctx context.Context, req *pbuser.GetAllUserIDReq) (resp *pbuser.GetAllUserIDResp, err error) {
	total, userIDs, err := s.UserDatabase.GetAllUserID(ctx, req.Pagination)
	if err != nil {
		return nil, err
	}
	return &pbuser.GetAllUserIDResp{Total: int32(total), UserIDs: userIDs}, nil
}

// SubscribeOrCancelUsersStatus Subscribe online or cancel online users.
func (s *userServer) SubscribeOrCancelUsersStatus(ctx context.Context, req *pbuser.SubscribeOrCancelUsersStatusReq) (resp *pbuser.SubscribeOrCancelUsersStatusResp, err error) {
	if req.Genre == constant.SubscriberUser {
		err = s.UserDatabase.SubscribeUsersStatus(ctx, req.UserID, req.UserIDs)
		if err != nil {
			return nil, err
		}
		var status []*pbuser.OnlineStatus
		status, err = s.UserDatabase.GetUserStatus(ctx, req.UserIDs)
		if err != nil {
			return nil, err
		}
		return &pbuser.SubscribeOrCancelUsersStatusResp{StatusList: status}, nil
	} else if req.Genre == constant.Unsubscribe {
		err = s.UserDatabase.UnsubscribeUsersStatus(ctx, req.UserID, req.UserIDs)
		if err != nil {
			return nil, err
		}
	}
	return &pbuser.SubscribeOrCancelUsersStatusResp{}, nil
}

// GetUserStatus Get the online status of the user.
func (s *userServer) GetUserStatus(ctx context.Context, req *pbuser.GetUserStatusReq) (resp *pbuser.GetUserStatusResp,
	err error) {
	onlineStatusList, err := s.UserDatabase.GetUserStatus(ctx, req.UserIDs)
	if err != nil {
		return nil, err
	}
	return &pbuser.GetUserStatusResp{StatusList: onlineStatusList}, nil
}

// SetUserStatus Synchronize user's online status.
func (s *userServer) SetUserStatus(ctx context.Context, req *pbuser.SetUserStatusReq) (resp *pbuser.SetUserStatusResp,
	err error) {
	err = s.UserDatabase.SetUserStatus(ctx, req.UserID, req.Status, req.PlatformID)
	if err != nil {
		return nil, err
	}
	list, err := s.UserDatabase.GetSubscribedList(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	for _, userID := range list {
		tips := &sdkws.UserStatusChangeTips{
			FromUserID: req.UserID,
			ToUserID:   userID,
			Status:     req.Status,
			PlatformID: req.PlatformID,
		}
		s.userNotificationSender.UserStatusChangeNotification(ctx, tips)
	}

	return &pbuser.SetUserStatusResp{}, nil
}

// GetSubscribeUsersStatus Get the online status of subscribers.
func (s *userServer) GetSubscribeUsersStatus(ctx context.Context,
	req *pbuser.GetSubscribeUsersStatusReq) (*pbuser.GetSubscribeUsersStatusResp, error) {
	userList, err := s.UserDatabase.GetAllSubscribeList(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	onlineStatusList, err := s.UserDatabase.GetUserStatus(ctx, userList)
	if err != nil {
		return nil, err
	}
	return &pbuser.GetSubscribeUsersStatusResp{StatusList: onlineStatusList}, nil
}
