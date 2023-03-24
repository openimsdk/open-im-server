package user

import (
	"context"
	"time"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/constant"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/controller"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/relation"
	tablerelation "github.com/OpenIMSDK/Open-IM-Server/pkg/common/db/table/relation"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/mcontext"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/tokenverify"
	registry "github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	pbuser "github.com/OpenIMSDK/Open-IM-Server/pkg/proto/user"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/rpcclient/check"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/rpcclient/convert"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/rpcclient/notification"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"google.golang.org/grpc"
)

type userServer struct {
	controller.UserDatabase
	notification        *notification.Check
	userCheck           *check.UserCheck
	conversationChecker *check.ConversationChecker
	RegisterCenter      registry.SvcDiscoveryRegistry
	friendCheck         *check.FriendChecker
}

func Start(client registry.SvcDiscoveryRegistry, server *grpc.Server) error {
	db, err := relation.NewGormDB()
	if err != nil {
		return err
	}
	if err := db.AutoMigrate(&tablerelation.UserModel{}); err != nil {
		return err
	}
	users := make([]*tablerelation.UserModel, 0)
	if len(config.Config.Manager.AppManagerUid) != len(config.Config.Manager.Nickname) {
		return errs.ErrConfig.Wrap("len(config.Config.Manager.AppManagerUid) != len(config.Config.Manager.Nickname)")
	}
	for k, v := range config.Config.Manager.AppManagerUid {
		users = append(users, &tablerelation.UserModel{UserID: v, Nickname: config.Config.Manager.Nickname[k]})
	}
	u := &userServer{
		UserDatabase:        controller.NewUserDatabase(relation.NewUserGorm(db)),
		notification:        notification.NewCheck(client),
		userCheck:           check.NewUserCheck(client),
		friendCheck:         check.NewFriendChecker(client),
		conversationChecker: check.NewConversationChecker(client),
		RegisterCenter:      client,
	}
	pbuser.RegisterUserServer(server, u)
	return u.UserDatabase.InitOnce(context.Background(), users)
}

// ok
func (s *userServer) GetDesignateUsers(ctx context.Context, req *pbuser.GetDesignateUsersReq) (resp *pbuser.GetDesignateUsersResp, err error) {
	resp = &pbuser.GetDesignateUsersResp{}
	users, err := s.FindWithError(ctx, req.UserIDs)
	if err != nil {
		return nil, err
	}
	resp.UsersInfo, err = (*convert.DBUser)(nil).DB2PB(users)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

// ok
func (s *userServer) UpdateUserInfo(ctx context.Context, req *pbuser.UpdateUserInfoReq) (resp *pbuser.UpdateUserInfoResp, err error) {
	resp = &pbuser.UpdateUserInfoResp{}
	err = tokenverify.CheckAccessV3(ctx, req.UserInfo.UserID)
	if err != nil {
		return nil, err
	}
	user, err := convert.NewPBUser(req.UserInfo).Convert()
	if err != nil {
		return nil, err
	}
	err = s.Update(ctx, user)
	if err != nil {
		return nil, err
	}
	friends, err := s.friendCheck.GetFriendIDs(ctx, req.UserInfo.UserID)
	if err != nil {
		return nil, err
	}
	go func() {
		for _, v := range friends {
			s.notification.FriendInfoUpdatedNotification(ctx, req.UserInfo.UserID, v, mcontext.GetOpUserID(ctx))
		}
	}()
	s.notification.UserInfoUpdatedNotification(ctx, mcontext.GetOpUserID(ctx), req.UserInfo.UserID)
	return resp, nil
}

// ok
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
	s.notification.UserInfoUpdatedNotification(ctx, req.UserID, req.UserID)
	return resp, nil
}

// ok
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

// ok
func (s *userServer) GetPaginationUsers(ctx context.Context, req *pbuser.GetPaginationUsersReq) (resp *pbuser.GetPaginationUsersResp, err error) {
	resp = &pbuser.GetPaginationUsersResp{}
	usersDB, total, err := s.Page(ctx, req.Pagination.PageNumber, req.Pagination.ShowNumber)
	if err != nil {
		return nil, err
	}
	resp.Total = int32(total)
	resp.Users, err = (*convert.DBUser)(nil).DB2PB(usersDB)
	return resp, err
}

// ok
func (s *userServer) UserRegister(ctx context.Context, req *pbuser.UserRegisterReq) (resp *pbuser.UserRegisterResp, err error) {
	resp = &pbuser.UserRegisterResp{}
	if len(req.Users) == 0 {
		return nil, errs.ErrArgs.Wrap("users is empty")
	}
	if utils.DuplicateAny(req.Users, func(e *sdkws.UserInfo) string { return e.UserID }) {
		return nil, errs.ErrArgs.Wrap("userID repeated")
	}
	userIDs := make([]string, 0)
	for _, user := range req.Users {
		if user.UserID == "" {
			return nil, errs.ErrArgs.Wrap("userID is empty")
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
	resp = &pbuser.GetGlobalRecvMessageOptResp{}
	user, err := s.FindWithError(ctx, []string{req.UserID})
	if err != nil {
		return nil, err
	}
	resp.GlobalRecvMsgOpt = user[0].GlobalRecvMsgOpt
	return resp, nil
}

func (s *userServer) GetAllUserID(ctx context.Context, req *pbuser.GetAllUserIDReq) (resp *pbuser.GetAllUserIDResp, err error) {
	userIDs, err := s.UserDatabase.GetAllUserID(ctx)
	if err != nil {
		return nil, err
	}
	resp = &pbuser.GetAllUserIDResp{UserIDs: userIDs}
	return resp, nil
}
