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
	"sync"
	"time"

	"github.com/openimsdk/open-im-server/v3/internal/rpc/relation"
	"github.com/openimsdk/open-im-server/v3/pkg/authverify"
	"github.com/openimsdk/open-im-server/v3/pkg/common/config"
	"github.com/openimsdk/open-im-server/v3/pkg/common/convert"
	"github.com/openimsdk/open-im-server/v3/pkg/common/prommetrics"
	"github.com/openimsdk/open-im-server/v3/pkg/common/servererrs"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/cache/redis"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/controller"
	"github.com/openimsdk/open-im-server/v3/pkg/common/storage/database/mgo"
	tablerelation "github.com/openimsdk/open-im-server/v3/pkg/common/storage/model"
	"github.com/openimsdk/open-im-server/v3/pkg/common/webhook"
	"github.com/openimsdk/open-im-server/v3/pkg/dbbuild"
	"github.com/openimsdk/open-im-server/v3/pkg/localcache"
	"github.com/openimsdk/open-im-server/v3/pkg/rpcli"
	"github.com/openimsdk/protocol/constant"
	"github.com/openimsdk/protocol/group"
	friendpb "github.com/openimsdk/protocol/relation"
	"github.com/openimsdk/protocol/sdkws"
	pbuser "github.com/openimsdk/protocol/user"
	"github.com/openimsdk/tools/db/pagination"
	"github.com/openimsdk/tools/discovery"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/utils/datautil"
	"google.golang.org/grpc"
)

const (
	defaultSecret = "openIM123"
)

type userServer struct {
	pbuser.UnimplementedUserServer
	online                   cache.OnlineCache
	db                       controller.UserDatabase
	friendNotificationSender *relation.FriendNotificationSender
	userNotificationSender   *UserNotificationSender
	RegisterCenter           discovery.Conn
	config                   *Config
	webhookClient            *webhook.Client
	groupClient              *rpcli.GroupClient
	relationClient           *rpcli.RelationClient
	clientConfig             controller.ClientConfigDatabase
}

type Config struct {
	RpcConfig          config.User
	RedisConfig        config.Redis
	MongodbConfig      config.Mongo
	KafkaConfig        config.Kafka
	NotificationConfig config.Notification
	Share              config.Share
	WebhooksConfig     config.Webhooks
	LocalCacheConfig   config.LocalCache
	Discovery          config.Discovery
}

func Start(ctx context.Context, config *Config, client discovery.SvcDiscoveryRegistry, server grpc.ServiceRegistrar) error {
	dbb := dbbuild.NewBuilder(&config.MongodbConfig, &config.RedisConfig)
	mgocli, err := dbb.Mongo(ctx)
	if err != nil {
		return err
	}
	rdb, err := dbb.Redis(ctx)
	if err != nil {
		return err
	}

	users := make([]*tablerelation.User, 0)

	for _, v := range config.Share.IMAdminUserID {
		users = append(users, &tablerelation.User{UserID: v, Nickname: v, AppMangerLevel: constant.AppAdmin})
	}
	userDB, err := mgo.NewUserMongo(mgocli.GetDB())
	if err != nil {
		return err
	}
	clientConfigDB, err := mgo.NewClientConfig(mgocli.GetDB())
	if err != nil {
		return err
	}
	msgConn, err := client.GetConn(ctx, config.Discovery.RpcService.Msg)
	if err != nil {
		return err
	}
	groupConn, err := client.GetConn(ctx, config.Discovery.RpcService.Group)
	if err != nil {
		return err
	}
	friendConn, err := client.GetConn(ctx, config.Discovery.RpcService.Friend)
	if err != nil {
		return err
	}
	msgClient := rpcli.NewMsgClient(msgConn)
	userCache := redis.NewUserCacheRedis(rdb, &config.LocalCacheConfig, userDB, redis.GetRocksCacheOptions())
	database := controller.NewUserDatabase(userDB, userCache, mgocli.GetTx())
	localcache.InitLocalCache(&config.LocalCacheConfig)
	u := &userServer{
		online:                   redis.NewUserOnline(rdb),
		db:                       database,
		RegisterCenter:           client,
		friendNotificationSender: relation.NewFriendNotificationSender(&config.NotificationConfig, msgClient, relation.WithDBFunc(database.FindWithError)),
		userNotificationSender:   NewUserNotificationSender(config, msgClient, WithUserFunc(database.FindWithError)),
		config:                   config,
		webhookClient:            webhook.NewWebhookClient(config.WebhooksConfig.URL),
		clientConfig:             controller.NewClientConfigDatabase(clientConfigDB, redis.NewClientConfigCache(rdb, clientConfigDB), mgocli.GetTx()),
		groupClient:              rpcli.NewGroupClient(groupConn),
		relationClient:           rpcli.NewRelationClient(friendConn),
	}
	pbuser.RegisterUserServer(server, u)
	return u.db.InitOnce(context.Background(), users)
}

func (s *userServer) GetDesignateUsers(ctx context.Context, req *pbuser.GetDesignateUsersReq) (resp *pbuser.GetDesignateUsersResp, err error) {
	resp = &pbuser.GetDesignateUsersResp{}
	users, err := s.db.Find(ctx, req.UserIDs)
	if err != nil {
		return nil, err
	}

	resp.UsersInfo = convert.UsersDB2Pb(users)
	return resp, nil
}

// deprecated:
// UpdateUserInfo
func (s *userServer) UpdateUserInfo(ctx context.Context, req *pbuser.UpdateUserInfoReq) (resp *pbuser.UpdateUserInfoResp, err error) {
	resp = &pbuser.UpdateUserInfoResp{}
	err = authverify.CheckAccess(ctx, req.UserInfo.UserID)
	if err != nil {
		return nil, err
	}

	if err := s.webhookBeforeUpdateUserInfo(ctx, &s.config.WebhooksConfig.BeforeUpdateUserInfo, req); err != nil {
		return nil, err
	}
	data := convert.UserPb2DBMap(req.UserInfo)
	oldUser, err := s.db.GetUserByID(ctx, req.UserInfo.UserID)
	if err != nil {
		return nil, err
	}
	if err := s.db.UpdateByMap(ctx, req.UserInfo.UserID, data); err != nil {
		return nil, err
	}
	s.friendNotificationSender.UserInfoUpdatedNotification(ctx, req.UserInfo.UserID)

	s.webhookAfterUpdateUserInfo(ctx, &s.config.WebhooksConfig.AfterUpdateUserInfo, req)
	if err = s.NotificationUserInfoUpdate(ctx, req.UserInfo.UserID, oldUser); err != nil {
		return nil, err
	}
	return resp, nil
}

func (s *userServer) UpdateUserInfoEx(ctx context.Context, req *pbuser.UpdateUserInfoExReq) (resp *pbuser.UpdateUserInfoExResp, err error) {
	resp = &pbuser.UpdateUserInfoExResp{}
	err = authverify.CheckAccess(ctx, req.UserInfo.UserID)
	if err != nil {
		return nil, err
	}

	if err = s.webhookBeforeUpdateUserInfoEx(ctx, &s.config.WebhooksConfig.BeforeUpdateUserInfoEx, req); err != nil {
		return nil, err
	}

	oldUser, err := s.db.GetUserByID(ctx, req.UserInfo.UserID)
	if err != nil {
		return nil, err
	}

	data := convert.UserPb2DBMapEx(req.UserInfo)
	if err = s.db.UpdateByMap(ctx, req.UserInfo.UserID, data); err != nil {
		return nil, err
	}

	s.friendNotificationSender.UserInfoUpdatedNotification(ctx, req.UserInfo.UserID)

	//friends, err := s.friendRpcClient.GetFriendIDs(ctx, req.UserInfo.UserID)
	//if err != nil {
	//	return nil, err
	//}
	//if req.UserInfo.Nickname != nil || req.UserInfo.FaceURL != nil {
	//	if err := s.NotificationUserInfoUpdate(ctx, req.UserInfo.UserID); err != nil {
	//		return nil, err
	//	}
	//}
	//for _, friendID := range friends {
	//	s.friendNotificationSender.FriendInfoUpdatedNotification(ctx, req.UserInfo.UserID, friendID)
	//}

	s.webhookAfterUpdateUserInfoEx(ctx, &s.config.WebhooksConfig.AfterUpdateUserInfoEx, req)
	if err := s.NotificationUserInfoUpdate(ctx, req.UserInfo.UserID, oldUser); err != nil {
		return nil, err
	}

	return resp, nil
}
func (s *userServer) SetGlobalRecvMessageOpt(ctx context.Context, req *pbuser.SetGlobalRecvMessageOptReq) (resp *pbuser.SetGlobalRecvMessageOptResp, err error) {
	resp = &pbuser.SetGlobalRecvMessageOptResp{}
	if _, err := s.db.FindWithError(ctx, []string{req.UserID}); err != nil {
		return nil, err
	}
	m := make(map[string]any, 1)
	m["global_recv_msg_opt"] = req.GlobalRecvMsgOpt
	if err := s.db.UpdateByMap(ctx, req.UserID, m); err != nil {
		return nil, err
	}
	s.friendNotificationSender.UserInfoUpdatedNotification(ctx, req.UserID)
	return resp, nil
}

func (s *userServer) AccountCheck(ctx context.Context, req *pbuser.AccountCheckReq) (resp *pbuser.AccountCheckResp, err error) {
	resp = &pbuser.AccountCheckResp{}
	if datautil.Duplicate(req.CheckUserIDs) {
		return nil, errs.ErrArgs.WrapMsg("userID repeated")
	}
	if err = authverify.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	users, err := s.db.Find(ctx, req.CheckUserIDs)
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
	if req.UserID == "" && req.NickName == "" {
		total, users, err := s.db.PageFindUser(ctx, constant.IMOrdinaryUser, constant.AppOrdinaryUsers, req.Pagination)
		if err != nil {
			return nil, err
		}
		return &pbuser.GetPaginationUsersResp{Total: int32(total), Users: convert.UsersDB2Pb(users)}, err
	} else {
		total, users, err := s.db.PageFindUserWithKeyword(ctx, constant.IMOrdinaryUser, constant.AppOrdinaryUsers, req.UserID, req.NickName, req.Pagination)
		if err != nil {
			return nil, err
		}
		return &pbuser.GetPaginationUsersResp{Total: int32(total), Users: convert.UsersDB2Pb(users)}, err

	}

}

func (s *userServer) UserRegister(ctx context.Context, req *pbuser.UserRegisterReq) (resp *pbuser.UserRegisterResp, err error) {
	resp = &pbuser.UserRegisterResp{}
	if len(req.Users) == 0 {
		return nil, errs.ErrArgs.WrapMsg("users is empty")
	}
	// check if secret is changed
	//if s.config.Share.Secret == defaultSecret {
	//	return nil, servererrs.ErrSecretNotChanged.Wrap()
	//}
	if err = authverify.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	if datautil.DuplicateAny(req.Users, func(e *sdkws.UserInfo) string { return e.UserID }) {
		return nil, errs.ErrArgs.WrapMsg("userID repeated")
	}
	userIDs := make([]string, 0)
	for _, user := range req.Users {
		if user.UserID == "" {
			return nil, errs.ErrArgs.WrapMsg("userID is empty")
		}
		if strings.Contains(user.UserID, ":") {
			return nil, errs.ErrArgs.WrapMsg("userID contains ':' is invalid userID")
		}
		userIDs = append(userIDs, user.UserID)
	}
	exist, err := s.db.IsExist(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, servererrs.ErrRegisteredAlready.WrapMsg("userID registered already")
	}
	if err := s.webhookBeforeUserRegister(ctx, &s.config.WebhooksConfig.BeforeUserRegister, req); err != nil {
		return nil, err
	}
	now := time.Now()
	users := make([]*tablerelation.User, 0, len(req.Users))
	for _, user := range req.Users {
		users = append(users, &tablerelation.User{
			UserID:           user.UserID,
			Nickname:         user.Nickname,
			FaceURL:          user.FaceURL,
			Ex:               user.Ex,
			CreateTime:       now,
			AppMangerLevel:   user.AppMangerLevel,
			GlobalRecvMsgOpt: user.GlobalRecvMsgOpt,
		})
	}
	if err := s.db.Create(ctx, users); err != nil {
		return nil, err
	}

	prommetrics.UserRegisterCounter.Add(float64(len(users)))

	s.webhookAfterUserRegister(ctx, &s.config.WebhooksConfig.AfterUserRegister, req)
	return resp, nil
}

func (s *userServer) GetGlobalRecvMessageOpt(ctx context.Context, req *pbuser.GetGlobalRecvMessageOptReq) (resp *pbuser.GetGlobalRecvMessageOptResp, err error) {
	user, err := s.db.FindWithError(ctx, []string{req.UserID})
	if err != nil {
		return nil, err
	}
	return &pbuser.GetGlobalRecvMessageOptResp{GlobalRecvMsgOpt: user[0].GlobalRecvMsgOpt}, nil
}

// GetAllUserID Get user account by page.
func (s *userServer) GetAllUserID(ctx context.Context, req *pbuser.GetAllUserIDReq) (resp *pbuser.GetAllUserIDResp, err error) {
	total, userIDs, err := s.db.GetAllUserID(ctx, req.Pagination)
	if err != nil {
		return nil, err
	}
	return &pbuser.GetAllUserIDResp{Total: int32(total), UserIDs: userIDs}, nil
}

// ProcessUserCommandAdd user general function add.
func (s *userServer) ProcessUserCommandAdd(ctx context.Context, req *pbuser.ProcessUserCommandAddReq) (*pbuser.ProcessUserCommandAddResp, error) {
	err := authverify.CheckAccess(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	var value string
	if req.Value != nil {
		value = req.Value.Value
	}
	var ex string
	if req.Ex != nil {
		value = req.Ex.Value
	}
	// Assuming you have a method in s.storage to add a user command
	err = s.db.AddUserCommand(ctx, req.UserID, req.Type, req.Uuid, value, ex)
	if err != nil {
		return nil, err
	}
	tips := &sdkws.UserCommandAddTips{
		FromUserID: req.UserID,
		ToUserID:   req.UserID,
	}
	s.userNotificationSender.UserCommandAddNotification(ctx, tips)
	return &pbuser.ProcessUserCommandAddResp{}, nil
}

// ProcessUserCommandDelete user general function delete.
func (s *userServer) ProcessUserCommandDelete(ctx context.Context, req *pbuser.ProcessUserCommandDeleteReq) (*pbuser.ProcessUserCommandDeleteResp, error) {
	err := authverify.CheckAccess(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	err = s.db.DeleteUserCommand(ctx, req.UserID, req.Type, req.Uuid)
	if err != nil {
		return nil, err
	}
	tips := &sdkws.UserCommandDeleteTips{
		FromUserID: req.UserID,
		ToUserID:   req.UserID,
	}
	s.userNotificationSender.UserCommandDeleteNotification(ctx, tips)
	return &pbuser.ProcessUserCommandDeleteResp{}, nil
}

// ProcessUserCommandUpdate user general function update.
func (s *userServer) ProcessUserCommandUpdate(ctx context.Context, req *pbuser.ProcessUserCommandUpdateReq) (*pbuser.ProcessUserCommandUpdateResp, error) {
	err := authverify.CheckAccess(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	val := make(map[string]any)

	// Map fields from eax to val
	if req.Value != nil {
		val["value"] = req.Value.Value
	}
	if req.Ex != nil {
		val["ex"] = req.Ex.Value
	}

	// Assuming you have a method in s.storage to update a user command
	err = s.db.UpdateUserCommand(ctx, req.UserID, req.Type, req.Uuid, val)
	if err != nil {
		return nil, err
	}
	tips := &sdkws.UserCommandUpdateTips{
		FromUserID: req.UserID,
		ToUserID:   req.UserID,
	}
	s.userNotificationSender.UserCommandUpdateNotification(ctx, tips)
	return &pbuser.ProcessUserCommandUpdateResp{}, nil
}

func (s *userServer) ProcessUserCommandGet(ctx context.Context, req *pbuser.ProcessUserCommandGetReq) (*pbuser.ProcessUserCommandGetResp, error) {

	err := authverify.CheckAccess(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	// Fetch user commands from the database
	commands, err := s.db.GetUserCommands(ctx, req.UserID, req.Type)
	if err != nil {
		return nil, err
	}

	// Initialize commandInfoSlice as an empty slice
	commandInfoSlice := make([]*pbuser.CommandInfoResp, 0, len(commands))

	for _, command := range commands {
		// No need to use index since command is already a pointer
		commandInfoSlice = append(commandInfoSlice, &pbuser.CommandInfoResp{
			Type:       command.Type,
			Uuid:       command.Uuid,
			Value:      command.Value,
			CreateTime: command.CreateTime,
			Ex:         command.Ex,
		})
	}

	// Return the response with the slice
	return &pbuser.ProcessUserCommandGetResp{CommandResp: commandInfoSlice}, nil
}

func (s *userServer) ProcessUserCommandGetAll(ctx context.Context, req *pbuser.ProcessUserCommandGetAllReq) (*pbuser.ProcessUserCommandGetAllResp, error) {
	err := authverify.CheckAccess(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	// Fetch user commands from the database
	commands, err := s.db.GetAllUserCommands(ctx, req.UserID)
	if err != nil {
		return nil, err
	}

	// Initialize commandInfoSlice as an empty slice
	commandInfoSlice := make([]*pbuser.AllCommandInfoResp, 0, len(commands))

	for _, command := range commands {
		// No need to use index since command is already a pointer
		commandInfoSlice = append(commandInfoSlice, &pbuser.AllCommandInfoResp{
			Type:       command.Type,
			Uuid:       command.Uuid,
			Value:      command.Value,
			CreateTime: command.CreateTime,
			Ex:         command.Ex,
		})
	}

	// Return the response with the slice
	return &pbuser.ProcessUserCommandGetAllResp{CommandResp: commandInfoSlice}, nil
}

func (s *userServer) AddNotificationAccount(ctx context.Context, req *pbuser.AddNotificationAccountReq) (*pbuser.AddNotificationAccountResp, error) {
	if err := authverify.CheckAdmin(ctx); err != nil {
		return nil, err
	}
	if req.AppMangerLevel < constant.AppNotificationAdmin {
		return nil, errs.ErrArgs.WithDetail("app level not supported")
	}
	if req.UserID == "" {
		for i := 0; i < 20; i++ {
			userId := s.genUserID()
			_, err := s.db.FindWithError(ctx, []string{userId})
			if err == nil {
				continue
			}
			req.UserID = userId
			break
		}
		if req.UserID == "" {
			return nil, errs.ErrInternalServer.WrapMsg("gen user id failed")
		}
	} else {
		_, err := s.db.FindWithError(ctx, []string{req.UserID})
		if err == nil {
			return nil, errs.ErrArgs.WrapMsg("userID is used")
		}
	}

	user := &tablerelation.User{
		UserID:         req.UserID,
		Nickname:       req.NickName,
		FaceURL:        req.FaceURL,
		CreateTime:     time.Now(),
		AppMangerLevel: req.AppMangerLevel,
	}
	if err := s.db.Create(ctx, []*tablerelation.User{user}); err != nil {
		return nil, err
	}

	return &pbuser.AddNotificationAccountResp{
		UserID:         req.UserID,
		NickName:       req.NickName,
		FaceURL:        req.FaceURL,
		AppMangerLevel: req.AppMangerLevel,
	}, nil
}

func (s *userServer) UpdateNotificationAccountInfo(ctx context.Context, req *pbuser.UpdateNotificationAccountInfoReq) (*pbuser.UpdateNotificationAccountInfoResp, error) {
	if err := authverify.CheckAdmin(ctx); err != nil {
		return nil, err
	}

	if _, err := s.db.FindWithError(ctx, []string{req.UserID}); err != nil {
		return nil, errs.ErrArgs.Wrap()
	}

	user := map[string]interface{}{}

	if req.NickName != "" {
		user["nickname"] = req.NickName
	}

	if req.FaceURL != "" {
		user["face_url"] = req.FaceURL
	}

	if err := s.db.UpdateByMap(ctx, req.UserID, user); err != nil {
		return nil, err
	}

	return &pbuser.UpdateNotificationAccountInfoResp{}, nil
}

func (s *userServer) SearchNotificationAccount(ctx context.Context, req *pbuser.SearchNotificationAccountReq) (*pbuser.SearchNotificationAccountResp, error) {
	// Check if user is an admin
	if err := authverify.CheckAdmin(ctx); err != nil {
		return nil, err
	}

	var users []*tablerelation.User
	var err error

	// If a keyword is provided in the request
	if req.Keyword != "" {
		// Find users by keyword
		users, err = s.db.Find(ctx, []string{req.Keyword})
		if err != nil {
			return nil, err
		}

		// Convert users to response format
		resp := s.userModelToResp(users, req.Pagination, req.AppManagerLevel)
		if resp.Total != 0 {
			return resp, nil
		}

		// Find users by nickname if no users found by keyword
		users, err = s.db.FindByNickname(ctx, req.Keyword)
		if err != nil {
			return nil, err
		}
		resp = s.userModelToResp(users, req.Pagination, req.AppManagerLevel)
		return resp, nil
	}

	// If no keyword, find users with notification settings
	if req.AppManagerLevel != nil {
		users, err = s.db.FindNotification(ctx, int64(*req.AppManagerLevel))
		if err != nil {
			return nil, err
		}
	} else {
		users, err = s.db.FindSystemAccount(ctx)
		if err != nil {
			return nil, err
		}
	}

	resp := s.userModelToResp(users, req.Pagination, req.AppManagerLevel)
	return resp, nil
}

func (s *userServer) GetNotificationAccount(ctx context.Context, req *pbuser.GetNotificationAccountReq) (*pbuser.GetNotificationAccountResp, error) {
	if req.UserID == "" {
		return nil, errs.ErrArgs.WrapMsg("userID is empty")
	}
	user, err := s.db.GetUserByID(ctx, req.UserID)
	if err != nil {
		return nil, servererrs.ErrUserIDNotFound.Wrap()
	}
	if user.AppMangerLevel >= constant.AppAdmin {
		return &pbuser.GetNotificationAccountResp{Account: &pbuser.NotificationAccountInfo{
			UserID:         user.UserID,
			FaceURL:        user.FaceURL,
			NickName:       user.Nickname,
			AppMangerLevel: user.AppMangerLevel,
		}}, nil
	}

	return nil, errs.ErrNoPermission.WrapMsg("notification messages cannot be sent for this ID")
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

func (s *userServer) userModelToResp(users []*tablerelation.User, pagination pagination.Pagination, appManagerLevel *int32) *pbuser.SearchNotificationAccountResp {
	accounts := make([]*pbuser.NotificationAccountInfo, 0)
	var total int64
	for _, v := range users {
		if v.AppMangerLevel >= constant.AppNotificationAdmin && !datautil.Contain(v.UserID, s.config.Share.IMAdminUserID...) {
			if appManagerLevel != nil {
				if v.AppMangerLevel != *appManagerLevel {
					continue
				}
			}
			temp := &pbuser.NotificationAccountInfo{
				UserID:         v.UserID,
				FaceURL:        v.FaceURL,
				NickName:       v.Nickname,
				AppMangerLevel: v.AppMangerLevel,
			}
			accounts = append(accounts, temp)
			total += 1
		}
	}

	notificationAccounts := datautil.Paginate(accounts, int(pagination.GetPageNumber()), int(pagination.GetShowNumber()))

	return &pbuser.SearchNotificationAccountResp{Total: total, NotificationAccounts: notificationAccounts}
}

func (s *userServer) NotificationUserInfoUpdate(ctx context.Context, userID string, oldUser *tablerelation.User) error {
	user, err := s.db.GetUserByID(ctx, userID)
	if err != nil {
		return err
	}
	if user.Nickname == oldUser.Nickname && user.FaceURL == oldUser.FaceURL {
		return nil
	}
	oldUserInfo := convert.UserDB2Pb(oldUser)
	newUserInfo := convert.UserDB2Pb(user)
	var wg sync.WaitGroup
	var es [2]error
	wg.Add(len(es))
	go func() {
		defer wg.Done()
		_, es[0] = s.groupClient.NotificationUserInfoUpdate(ctx, &group.NotificationUserInfoUpdateReq{
			UserID:      userID,
			OldUserInfo: oldUserInfo,
			NewUserInfo: newUserInfo,
		})
	}()

	go func() {
		defer wg.Done()
		_, es[1] = s.relationClient.NotificationUserInfoUpdate(ctx, &friendpb.NotificationUserInfoUpdateReq{
			UserID:      userID,
			OldUserInfo: oldUserInfo,
			NewUserInfo: newUserInfo,
		})
	}()
	wg.Wait()
	return errors.Join(es[:]...)
}

func (s *userServer) SortQuery(ctx context.Context, req *pbuser.SortQueryReq) (*pbuser.SortQueryResp, error) {
	users, err := s.db.SortQuery(ctx, req.UserIDName, req.Asc)
	if err != nil {
		return nil, err
	}
	return &pbuser.SortQueryResp{Users: convert.UsersDB2Pb(users)}, nil
}
