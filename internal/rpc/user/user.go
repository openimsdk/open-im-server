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
	"math/rand"
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
	if len(config.Config.IMAdmin.UserID) != len(config.Config.IMAdmin.Nickname) {
		return errors.New("len(config.Config.AppNotificationAdmin.AppManagerUid) != len(config.Config.AppNotificationAdmin.Nickname)")
	}
	for k, v := range config.Config.IMAdmin.UserID {
		users = append(users, &tablerelation.UserModel{UserID: v, Nickname: config.Config.IMAdmin.Nickname[k], AppMangerLevel: constant.AppNotificationAdmin})
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

// ProcessUserCommandAdd user general function add
func (s *userServer) ProcessUserCommandAdd(ctx context.Context, req *pbuser.ProcessUserCommandAddReq) (*pbuser.ProcessUserCommandAddResp, error) {
	// Assuming you have a method in s.UserDatabase to add a user command
	err := s.UserDatabase.AddUserCommand(ctx, req.UserID, req.Type, req.Uuid, req.Value)
	if err != nil {
		return nil, err
	}

	return &pbuser.ProcessUserCommandAddResp{}, nil
}

// ProcessUserCommandDelete user general function delete
func (s *userServer) ProcessUserCommandDelete(ctx context.Context, req *pbuser.ProcessUserCommandDeleteReq) (*pbuser.ProcessUserCommandDeleteResp, error) {
	// Assuming you have a method in s.UserDatabase to delete a user command
	err := s.UserDatabase.DeleteUserCommand(ctx, req.UserID, req.Type, req.Uuid)
	if err != nil {
		return nil, err
	}

	return &pbuser.ProcessUserCommandDeleteResp{}, nil
}

// ProcessUserCommandUpdate user general function update
func (s *userServer) ProcessUserCommandUpdate(ctx context.Context, req *pbuser.ProcessUserCommandUpdateReq) (*pbuser.ProcessUserCommandUpdateResp, error) {
	// Assuming you have a method in s.UserDatabase to update a user command
	err := s.UserDatabase.UpdateUserCommand(ctx, req.UserID, req.Type, req.Uuid, req.Value)
	if err != nil {
		return nil, err
	}

	return &pbuser.ProcessUserCommandUpdateResp{}, nil
}

func (s *userServer) ProcessUserCommandGet(ctx context.Context, req *pbuser.ProcessUserCommandGetReq) (*pbuser.ProcessUserCommandGetResp, error) {
	// Fetch user commands from the database
	commands, err := s.UserDatabase.GetUserCommands(ctx, req.UserID, req.Type)
	if err != nil {
		return nil, err
	}

	// Initialize commandInfoSlice as an empty slice
	commandInfoSlice := make([]*pbuser.CommandInfoResp, 0, len(commands))

	for _, command := range commands {
		// No need to use index since command is already a pointer
		commandInfoSlice = append(commandInfoSlice, &pbuser.CommandInfoResp{
			Uuid:       command.Uuid,
			Value:      command.Value,
			CreateTime: command.CreateTime,
		})
	}

	// Return the response with the slice
	return &pbuser.ProcessUserCommandGetResp{KVArray: commandInfoSlice}, nil
}

func (s *userServer) AddNotificationAccount(ctx context.Context, req *pbuser.AddNotificationAccountReq) (*pbuser.AddNotificationAccountResp, error) {
	if err := authverify.CheckIMAdmin(ctx); err != nil {
		return nil, err
	}

	var userID string
	for i := 0; i < 20; i++ {
		userId := s.genUserID()
		_, err := s.UserDatabase.FindWithError(ctx, []string{userId})
		if err == nil {
			continue
		}
		userID = userId
		break
	}
	if userID == "" {
		return nil, errs.ErrInternalServer.Wrap("gen user id failed")
	}

	user := &tablerelation.UserModel{
		UserID:         userID,
		Nickname:       req.NickName,
		FaceURL:        req.FaceURL,
		CreateTime:     time.Now(),
		AppMangerLevel: constant.AppNotificationAdmin,
	}
	if err := s.UserDatabase.Create(ctx, []*tablerelation.UserModel{user}); err != nil {
		return nil, err
	}

	return &pbuser.AddNotificationAccountResp{}, nil
}

func (s *userServer) UpdateNotificationAccountInfo(ctx context.Context, req *pbuser.UpdateNotificationAccountInfoReq) (*pbuser.UpdateNotificationAccountInfoResp, error) {
	if err := authverify.CheckIMAdmin(ctx); err != nil {
		return nil, err
	}

	if _, err := s.UserDatabase.FindWithError(ctx, []string{req.UserID}); err != nil {
		return nil, errs.ErrArgs.Wrap()
	}

	user := map[string]interface{}{}

	if req.NickName != "" {
		user["nickname"] = req.NickName
	}

	if req.FaceURL != "" {
		user["face_url"] = req.FaceURL
	}

	if err := s.UserDatabase.UpdateByMap(ctx, req.UserID, user); err != nil {
		return nil, err
	}

	return &pbuser.UpdateNotificationAccountInfoResp{}, nil
}

func (s *userServer) SearchNotificationAccount(ctx context.Context, req *pbuser.SearchNotificationAccountReq) (*pbuser.SearchNotificationAccountResp, error) {
	if err := authverify.CheckIMAdmin(ctx); err != nil {
		return nil, err
	}

	_, users, err := s.UserDatabase.Page(ctx, req.Pagination)
	if err != nil {
		return nil, err
	}

	var total int64
	accounts := make([]*pbuser.NotificationAccountInfo, 0, len(users))
	for _, v := range users {
		if v.AppMangerLevel != constant.AppNotificationAdmin {
			continue
		}
		temp := &pbuser.NotificationAccountInfo{
			UserID:   v.UserID,
			FaceURL:  v.FaceURL,
			NickName: v.Nickname,
		}
		accounts = append(accounts, temp)
		total += 1
	}
	return &pbuser.SearchNotificationAccountResp{Total: total, NotificationAccounts: accounts}, nil
}

func (s *userServer) UpdateUserInfoEx(ctx context.Context, req *pbuser.UpdateUserInfoExReq) (*pbuser.UpdateUserInfoExResp, error) {
	//TODO implement me
	panic("implement me")
}

func (s *userServer) GetNotificationAccount(ctx context.Context, req *pbuser.GetNotificationAccountReq) (*pbuser.GetNotificationAccountResp, error) {
	if req.UserID == "" {
		return nil, errs.ErrArgs.Wrap("userID is empty")
	}
	user, err := s.UserDatabase.GetUserByID(ctx, req.UserID)
	if err != nil {
		return nil, errs.ErrUserIDNotFound.Wrap()
	}
	if user.AppMangerLevel == constant.AppAdmin || user.AppMangerLevel == constant.AppNotificationAdmin {
		return &pbuser.GetNotificationAccountResp{}, nil
	}

	return nil, errs.ErrNoPermission.Wrap("notification messages cannot be sent for this ID")
}

func (s *userServer) genUserID() string {
	const l = 10
	data := make([]byte, l)
	rand.Read(data)
	chars := []byte("0123456789")
	for i := 0; i < len(data); i++ {
		if i == 0 {
			data[i] = chars[1:][data[i]%9]
		} else {
			data[i] = chars[data[i]%10]
		}
	}
	return string(data)
}
