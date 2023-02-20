package check

import (
	"Open_IM/pkg/common/config"
	"Open_IM/pkg/common/constant"
	discoveryRegistry "Open_IM/pkg/discoveryregistry"
	sdkws "Open_IM/pkg/proto/sdkws"
	"Open_IM/pkg/proto/user"
	"Open_IM/pkg/utils"
	"context"
	"google.golang.org/grpc"
	"strings"
)

func NewUserCheck(zk discoveryRegistry.SvcDiscoveryRegistry) *UserCheck {
	return &UserCheck{
		zk: zk,
	}
}

type UserCheck struct {
	zk discoveryRegistry.SvcDiscoveryRegistry
}

func (u *UserCheck) getConn() (*grpc.ClientConn, error) {
	return u.zk.GetConn(config.Config.RpcRegisterName.OpenImUserName)
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
		return nil, err
	}
	if complete {
		if ids := utils.Single(userIDs, utils.Slice(resp.UsersInfo, func(e *sdkws.UserInfo) string {
			return e.UserID
		})); len(ids) > 0 {
			return nil, constant.ErrUserIDNotFound.Wrap(strings.Join(ids, ","))
		}
	}
	return resp.UsersInfo, nil
}

func (u *UserCheck) GetUsersInfo(ctx context.Context, userID string) (*sdkws.UserInfo, error) {
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
	return 0, nil
}
