package user

import (
	"context"
	"errors"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/log"
	"strings"
	"time"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/convert"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/cache"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/controller"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/relation"
	tablerelation "github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/tx"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/tokenverify"
	registry "github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	pbuser "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/user"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/rpcclient"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/rpcclient/notification"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"google.golang.org/grpc"
)

type userServer struct {
	controller.UserDatabase
	notificationSender *notification.FriendNotificationSender
	friendRpcClient    *rpcclient.FriendRpcClient
	RegisterCenter     registry.SvcDiscoveryRegistry
}

func Start(client registry.SvcDiscoveryRegistry, server *grpc.Server) error {
	db, err := relation.NewGormDB()
	if err != nil {
		return err
	}
	rdb, err := cache.NewRedis()
	if err != nil {
		return err
	}
	if err := db.AutoMigrate(&tablerelation.UserModel{}); err != nil {
		return err
	}
	users := make([]*tablerelation.UserModel, 0)
	if len(config.Config.Manager.UserID) != len(config.Config.Manager.Nickname) {
		return errors.New("len(config.Config.Manager.AppManagerUid) != len(config.Config.Manager.Nickname)")
	}
	for k, v := range config.Config.Manager.UserID {
		users = append(users, &tablerelation.UserModel{UserID: v, Nickname: config.Config.Manager.Nickname[k]})
	}
	userDB := relation.NewUserGorm(db)
	cache := cache.NewUserCacheRedis(rdb, userDB, cache.GetDefaultOpt())
	database := controller.NewUserDatabase(userDB, cache, tx.NewGorm(db))
	friendRpcClient := rpcclient.NewFriendRpcClient(client)
	msgRpcClient := rpcclient.NewMessageRpcClient(client)
	u := &userServer{
		UserDatabase:       database,
		RegisterCenter:     client,
		friendRpcClient:    &friendRpcClient,
		notificationSender: notification.NewFriendNotificationSender(&msgRpcClient, notification.WithDBFunc(database.FindWithError)),
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
	err = tokenverify.CheckAccessV3(ctx, req.UserInfo.UserID)
	if err != nil {
		return nil, err
	}
	user := convert.UserPb2DB(req.UserInfo)
	if err != nil {
		return nil, err
	}
	err = s.Update(ctx, user)
	if err != nil {
		return nil, err
	}
	_ = s.notificationSender.UserInfoUpdatedNotification(ctx, req.UserInfo.UserID)
	friends, err := s.friendRpcClient.GetFriendIDs(ctx, req.UserInfo.UserID)
	if err != nil {
		return nil, err
	}
	for _, friendID := range friends {
		s.notificationSender.FriendInfoUpdatedNotification(ctx, req.UserInfo.UserID, friendID)
	}
	return resp, nil
}

func (s *userServer) SetGlobalRecvMessageOpt(ctx context.Context, req *pbuser.SetGlobalRecvMessageOptReq) (resp *pbuser.SetGlobalRecvMessageOptResp, err error) {
	resp = &pbuser.SetGlobalRecvMessageOptResp{}
	if _, err := s.FindWithError(ctx, []string{req.UserID}); err != nil {
		return nil, err
	}
	m := make(map[string]interface{}, 1)
	m["global_recv_msg_opt"] = req.GlobalRecvMsgOpt
	if err := s.UpdateByMap(ctx, req.UserID, m); err != nil {
		return nil, err
	}
	s.notificationSender.UserInfoUpdatedNotification(ctx, req.UserID)
	return resp, nil
}

func (s *userServer) AccountCheck(ctx context.Context, req *pbuser.AccountCheckReq) (resp *pbuser.AccountCheckResp, err error) {
	resp = &pbuser.AccountCheckResp{}
	if utils.Duplicate(req.CheckUserIDs) {
		return nil, errs.ErrArgs.Wrap("userID repeated")
	}
	err = tokenverify.CheckAdmin(ctx)
	if err != nil {
		return nil, err
	}
	users, err := s.Find(ctx, req.CheckUserIDs)
	if err != nil {
		return nil, err
	}
	userIDs := make(map[string]interface{}, 0)
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
	var pageNumber, showNumber int32
	if req.Pagination != nil {
		pageNumber = req.Pagination.PageNumber
		showNumber = req.Pagination.ShowNumber
	}
	users, total, err := s.Page(ctx, pageNumber, showNumber)
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
	return resp, nil
}

func (s *userServer) GetGlobalRecvMessageOpt(ctx context.Context, req *pbuser.GetGlobalRecvMessageOptReq) (resp *pbuser.GetGlobalRecvMessageOptResp, err error) {
	user, err := s.FindWithError(ctx, []string{req.UserID})
	if err != nil {
		return nil, err
	}
	return &pbuser.GetGlobalRecvMessageOptResp{GlobalRecvMsgOpt: user[0].GlobalRecvMsgOpt}, nil
}

func (s *userServer) GetAllUserID(ctx context.Context, req *pbuser.GetAllUserIDReq) (resp *pbuser.GetAllUserIDResp, err error) {
	userIDs, err := s.UserDatabase.GetAllUserID(ctx)
	if err != nil {
		return nil, err
	}
	return &pbuser.GetAllUserIDResp{UserIDs: userIDs}, nil
}
