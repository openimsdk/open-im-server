package check

import (
	"OpenIM/pkg/common/config"
	"OpenIM/pkg/common/log"
	"OpenIM/pkg/discoveryregistry"
	"OpenIM/pkg/errs"
	"OpenIM/pkg/proto/sdkws"
	"OpenIM/pkg/proto/user"
	"OpenIM/pkg/utils"
	"context"
	"google.golang.org/grpc"
	"strings"
)

func NewUserCheck(client discoveryregistry.SvcDiscoveryRegistry) *UserCheck {
	return &UserCheck{
		client: client,
	}
}

type UserCheck struct {
	client discoveryregistry.SvcDiscoveryRegistry
}

func (u *UserCheck) getConn() (*grpc.ClientConn, error) {
	return u.client.GetConn(config.Config.RpcRegisterName.OpenImUserName)
}

func (u *UserCheck) GetUsersInfos(ctx context.Context, userIDs []string, complete bool) ([]*sdkws.UserInfo, error) {
	cc, err := u.getConn()
	if err != nil {
		return nil, err
	}
	resp, err := user.NewUserClient(cc).GetDesignateUsers(ctx, &user.GetDesignateUsersReq{
		UserIDs: userIDs,
	})
	if err != nil {
		log.Error("", "call GetDesignateUsers err", err.Error())
		return nil, err
	}
	if complete {
		if ids := utils.Single(userIDs, utils.Slice(resp.UsersInfo, func(e *sdkws.UserInfo) string {
			return e.UserID
		})); len(ids) > 0 {
			return nil, errs.ErrUserIDNotFound.Wrap(strings.Join(ids, ","))
		}
	}
	return resp.UsersInfo, nil
}

func (u *UserCheck) GetUserInfo(ctx context.Context, userID string) (*sdkws.UserInfo, error) {
	users, err := u.GetUsersInfos(ctx, []string{userID}, true)
	if err != nil {
		return nil, err
	}
	return users[0], nil
}

func (u *UserCheck) GetUsersInfoMap(ctx context.Context, userIDs []string, complete bool) (map[string]*sdkws.UserInfo, error) {
	users, err := u.GetUsersInfos(ctx, userIDs, complete)
	if err != nil {
		return nil, err
	}
	return utils.SliceToMap(users, func(e *sdkws.UserInfo) string {
		return e.UserID
	}), nil
}

func (u *UserCheck) GetPublicUserInfos(ctx context.Context, userIDs []string, complete bool) ([]*sdkws.PublicUserInfo, error) {
	users, err := u.GetUsersInfos(ctx, userIDs, complete)
	if err != nil {
		return nil, err
	}
	return utils.Slice(users, func(e *sdkws.UserInfo) *sdkws.PublicUserInfo {
		return &sdkws.PublicUserInfo{
			UserID:   e.UserID,
			Nickname: e.Nickname,
			FaceURL:  e.FaceURL,
			Gender:   e.Gender,
			Ex:       e.Ex,
		}
	}), nil
}

func (u *UserCheck) GetPublicUserInfo(ctx context.Context, userID string) (*sdkws.PublicUserInfo, error) {
	users, err := u.GetPublicUserInfos(ctx, []string{userID}, true)
	if err != nil {
		return nil, err
	}
	return users[0], nil
}

func (u *UserCheck) GetPublicUserInfoMap(ctx context.Context, userIDs []string, complete bool) (map[string]*sdkws.PublicUserInfo, error) {
	users, err := u.GetPublicUserInfos(ctx, userIDs, complete)
	if err != nil {
		return nil, err
	}
	return utils.SliceToMap(users, func(e *sdkws.PublicUserInfo) string {
		return e.UserID
	}), nil
}

func (u *UserCheck) GetUserGlobalMsgRecvOpt(ctx context.Context, userID string) (int32, error) {
	cc, err := u.getConn()
	if err != nil {
		return 0, err
	}
	resp, err := user.NewUserClient(cc).GetGlobalRecvMessageOpt(ctx, &user.GetGlobalRecvMessageOptReq{
		UserID: userID,
	})
	if err != nil {
		return 0, err
	}
	return resp.GlobalRecvMsgOpt, err
}
