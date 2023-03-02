package user

import (
	"OpenIM/internal/common/check"
	"OpenIM/internal/common/convert"
	"OpenIM/internal/common/notification"
	"OpenIM/pkg/common/config"
	"OpenIM/pkg/common/constant"
	"OpenIM/pkg/common/db/controller"
	"OpenIM/pkg/common/db/relation"
	tablerelation "OpenIM/pkg/common/db/table/relation"
	"OpenIM/pkg/common/tokenverify"
	"OpenIM/pkg/common/tracelog"
	registry "OpenIM/pkg/discoveryregistry"
	"OpenIM/pkg/proto/sdkws"
	pbuser "OpenIM/pkg/proto/user"
	"OpenIM/pkg/utils"
	"context"
	"github.com/OpenIMSDK/openKeeper"
	"google.golang.org/grpc"
)

type userServer struct {
	controller.UserDatabase
	notification        *notification.Check
	userCheck           *check.UserCheck
	ConversationChecker *check.ConversationChecker
	RegisterCenter      registry.SvcDiscoveryRegistry
	friendCheck         *check.FriendChecker
}

func Start(client registry.SvcDiscoveryRegistry, server *grpc.Server) error {
	gormDB, err := relation.NewGormDB()
	if err != nil {
		return err
	}
	if err := gormDB.AutoMigrate(&tablerelation.UserModel{}); err != nil {
		return err
	}
	users := make([]*tablerelation.UserModel, 0)
	if len(config.Config.Manager.AppManagerUid) != len(config.Config.Manager.Nickname) {
		return constant.ErrConfig.Wrap("len(config.Config.Manager.AppManagerUid) != len(config.Config.Manager.Nickname)")
	}
	for k, v := range config.Config.Manager.AppManagerUid {
		users = append(users, &tablerelation.UserModel{UserID: v, Nickname: config.Config.Manager.Nickname[k]})
	}
	u := &userServer{
		UserDatabase:   controller.NewUserDatabase(relation.NewUserGorm(gormDB)),
		notification:   notification.NewCheck(client),
		userCheck:      check.NewUserCheck(client),
		RegisterCenter: client,
	}
	pbuser.RegisterUserServer(server, u)
	u.UserDatabase.InitOnce(context.Background(), users)
	return nil
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
	err = s.Update(ctx, []*tablerelation.UserModel{user})
	if err != nil {
		return nil, err
	}
	friends, err := s.friendCheck.GetAllPageFriends(ctx, req.UserInfo.UserID)
	if err != nil {
		return nil, err
	}
	go func() {
		for _, v := range friends {
			s.notification.FriendInfoUpdatedNotification(ctx, req.UserInfo.UserID, v.FriendUser.UserID, tracelog.GetOpUserID(ctx))
		}
	}()
	s.notification.UserInfoUpdatedNotification(ctx, tracelog.GetOpUserID(ctx), req.UserInfo.UserID)
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
		return nil, constant.ErrArgs.Wrap("userID repeated")
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
	return resp, nil
}

// ok
func (s *userServer) UserRegister(ctx context.Context, req *pbuser.UserRegisterReq) (resp *pbuser.UserRegisterResp, err error) {
	resp = &pbuser.UserRegisterResp{}
	if utils.DuplicateAny(req.Users, func(e *sdkws.UserInfo) string { return e.UserID }) {
		return nil, constant.ErrArgs.Wrap("userID repeated")
	}
	userIDs := make([]string, 0)
	for _, v := range req.Users {
		userIDs = append(userIDs, v.UserID)
	}
	exist, err := s.IsExist(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, constant.ErrRegisteredAlready.Wrap("userID registered already")
	}
	users, err := (*convert.PBUser)(nil).PB2DB(req.Users)
	if err != nil {
		return nil, err
	}
	err = s.Create(ctx, users)
	if err != nil {
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
