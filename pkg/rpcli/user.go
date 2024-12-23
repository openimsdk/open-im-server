package rpcli

import (
	"context"
	"github.com/openimsdk/protocol/sdkws"
	"github.com/openimsdk/protocol/user"
	"github.com/openimsdk/tools/errs"
	"github.com/openimsdk/tools/utils/datautil"
	"google.golang.org/grpc"
)

func NewUserClient(cc grpc.ClientConnInterface) *UserClient {
	return &UserClient{user.NewUserClient(cc)}
}

type UserClient struct {
	user.UserClient
}

func (x *UserClient) GetUsersInfo(ctx context.Context, userIDs []string) ([]*sdkws.UserInfo, error) {
	req := &user.GetDesignateUsersReq{UserIDs: userIDs}
	return extractField(ctx, x.UserClient.GetDesignateUsers, req, (*user.GetDesignateUsersResp).GetUsersInfo)
}

func (x *UserClient) GetUserInfo(ctx context.Context, userID string) (*sdkws.UserInfo, error) {
	return firstValue(x.GetUsersInfo(ctx, []string{userID}))
}

func (x *UserClient) CheckUser(ctx context.Context, userIDs []string) error {
	users, err := x.GetUsersInfo(ctx, userIDs)
	if err != nil {
		return err
	}
	if len(users) != len(userIDs) {
		return errs.ErrRecordNotFound.WrapMsg("user not found")
	}
	return nil
}

func (x *UserClient) GetUsersInfoMap(ctx context.Context, userIDs []string) (map[string]*sdkws.UserInfo, error) {
	users, err := x.GetUsersInfo(ctx, userIDs)
	if err != nil {
		return nil, err
	}
	return datautil.SliceToMap(users, func(e *sdkws.UserInfo) string {
		return e.UserID
	}), nil
}

func (x *UserClient) GetAllOnlineUsers(ctx context.Context, cursor uint64) (*user.GetAllOnlineUsersResp, error) {
	req := &user.GetAllOnlineUsersReq{Cursor: cursor}
	return x.UserClient.GetAllOnlineUsers(ctx, req)
}

func (x *UserClient) GetUsersOnlinePlatform(ctx context.Context, userIDs []string) ([]*user.OnlineStatus, error) {
	if len(userIDs) == 0 {
		return nil, nil
	}
	req := &user.GetUserStatusReq{UserIDs: userIDs}
	return extractField(ctx, x.UserClient.GetUserStatus, req, (*user.GetUserStatusResp).GetStatusList)

}

func (x *UserClient) GetUserOnlinePlatform(ctx context.Context, userID string) ([]int32, error) {
	status, err := x.GetUsersOnlinePlatform(ctx, []string{userID})
	if err != nil {
		return nil, err
	}
	if len(status) == 0 {
		return nil, nil
	}
	return status[0].PlatformIDs, nil
}

func (x *UserClient) SetUserOnlineStatus(ctx context.Context, req *user.SetUserOnlineStatusReq) error {
	return ignoreResp(x.UserClient.SetUserOnlineStatus(ctx, req))
}
