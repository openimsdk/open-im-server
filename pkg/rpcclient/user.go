package rpcclient

import (
	"context"
	"strings"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/config"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/common/tokenverify"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/discoveryregistry"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/errs"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/user"
	"github.com/OpenIMSDK/Open-IM-Server/pkg/utils"
	"google.golang.org/grpc"
)

type UserClient struct {
	conn *grpc.ClientConn
}

func NewUserClient(discov discoveryregistry.SvcDiscoveryRegistry) *UserClient {
	conn, err := discov.GetConn(context.Background(), config.Config.RpcRegisterName.OpenImUserName)
	if err != nil {
		panic(err)
	}
	return &UserClient{conn: conn}
}

func (u *UserClient) GetUsersInfo(ctx context.Context, userIDs []string) ([]*sdkws.UserInfo, error) {
	resp, err := user.NewUserClient(u.conn).GetDesignateUsers(ctx, &user.GetDesignateUsersReq{
		UserIDs: userIDs,
	})
	if err != nil {
		return nil, err
	}
	if ids := utils.Single(userIDs, utils.Slice(resp.UsersInfo, func(e *sdkws.UserInfo) string {
		return e.UserID
	})); len(ids) > 0 {
		return nil, errs.ErrUserIDNotFound.Wrap(strings.Join(ids, ","))
	}
	return resp.UsersInfo, nil
}

func (u *UserClient) GetUserInfo(ctx context.Context, userID string) (*sdkws.UserInfo, error) {
	users, err := u.GetUsersInfo(ctx, []string{userID})
	if err != nil {
		return nil, err
	}
	return users[0], nil
}

func (u *UserClient) GetUsersInfoMap(ctx context.Context, userIDs []string) (map[string]*sdkws.UserInfo, error) {
	users, err := u.GetUsersInfo(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	return utils.SliceToMap(users, func(e *sdkws.UserInfo) string {
		return e.UserID
	}), nil
}

func (u *UserClient) GetPublicUserInfos(ctx context.Context, userIDs []string, complete bool) ([]*sdkws.PublicUserInfo, error) {
	users, err := u.GetUsersInfo(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	return utils.Slice(users, func(e *sdkws.UserInfo) *sdkws.PublicUserInfo {
		return &sdkws.PublicUserInfo{
			UserID:   e.UserID,
			Nickname: e.Nickname,
			FaceURL:  e.FaceURL,
			Ex:       e.Ex,
		}
	}), nil
}

func (u *UserClient) GetPublicUserInfo(ctx context.Context, userID string) (*sdkws.PublicUserInfo, error) {
	users, err := u.GetPublicUserInfos(ctx, []string{userID}, true)
	if err != nil {
		return nil, err
	}
	return users[0], nil
}

func (u *UserClient) GetPublicUserInfoMap(ctx context.Context, userIDs []string, complete bool) (map[string]*sdkws.PublicUserInfo, error) {
	users, err := u.GetPublicUserInfos(ctx, userIDs, complete)
	if err != nil {
		return nil, err
	}
	return utils.SliceToMap(users, func(e *sdkws.PublicUserInfo) string {
		return e.UserID
	}), nil
}

func (u *UserClient) GetUserGlobalMsgRecvOpt(ctx context.Context, userID string) (int32, error) {
	resp, err := user.NewUserClient(u.conn).GetGlobalRecvMessageOpt(ctx, &user.GetGlobalRecvMessageOptReq{
		UserID: userID,
	})
	if err != nil {
		return 0, err
	}
	return resp.GlobalRecvMsgOpt, err
}

func (u *UserClient) Access(ctx context.Context, ownerUserID string) error {
	_, err := u.GetUserInfo(ctx, ownerUserID)
	if err != nil {
		return err
	}
	return tokenverify.CheckAccessV3(ctx, ownerUserID)
}
